package feed

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/like"
	"feedsystem_video_go/internal/middleware/redis"
	"feedsystem_video_go/internal/pkg/logger"
	"feedsystem_video_go/internal/video"
)

const (
	CacheTTL = 5 * time.Minute
	HotTTL   = 30 * time.Minute
)

type FeedService struct {
	videoRepo   *video.VideoRepository
	accountRepo *account.AccountRepository
	likeRepo    *like.LikeRepository
	cache       *redis.Client
	localCache  sync.Map
	baseURL     string
}

func NewFeedService(videoRepo *video.VideoRepository, accountRepo *account.AccountRepository, likeRepo *like.LikeRepository, cache *redis.Client, baseURL string) *FeedService {
	return &FeedService{
		videoRepo:   videoRepo,
		accountRepo: accountRepo,
		likeRepo:    likeRepo,
		cache:       cache,
		baseURL:     baseURL,
	}
}

func (fs *FeedService) GetFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	switch req.Type {
	case "hot":
		return fs.getHotFeed(ctx, req)
	case "following":
		return fs.getFollowingFeed(ctx, req)
	case "tag":
		return fs.getTagFeed(ctx, req)
	case "search":
		return fs.searchFeed(ctx, req)
	default:
		return fs.getLatestFeed(ctx, req)
	}
}

func (fs *FeedService) getLatestFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	cacheKey := fs.cache.Key("feed:latest:%d", req.Cursor)

	if cached, ok := fs.localCache.Load(cacheKey); ok {
		if items, ok := cached.([]FeedItem); ok {
			return &FeedResponse{Items: items, HasMore: true, Next: req.Cursor + int64(req.Limit)}, nil
		}
	}

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if cached, err := fs.cache.Get(cacheCtx, cacheKey); err == nil && cached != "" {
			var items []FeedItem
			if err := decodeFeedItems(cached, &items); err == nil {
				fs.localCache.Store(cacheKey, items)
				return &FeedResponse{Items: items, HasMore: true, Next: req.Cursor + int64(req.Limit)}, nil
			}
		}
	}

	videos, err := fs.videoRepo.FindLatest(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	items, err := fs.buildFeedItems(ctx, videos, req.AccountID)
	if err != nil {
		return nil, err
	}

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if encoded, err := encodeFeedItems(items); err == nil {
			go fs.cache.Set(cacheCtx, cacheKey, encoded, CacheTTL)
		}
	}
	fs.localCache.Store(cacheKey, items)

	return &FeedResponse{
		Items:   items,
		HasMore: len(items) == req.Limit,
		Next:    req.Cursor + int64(req.Limit),
	}, nil
}

func (fs *FeedService) getHotFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	cacheKey := "feed:hot"

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if cached, err := fs.cache.Get(cacheCtx, cacheKey); err == nil && cached != "" {
			var items []FeedItem
			if err := decodeFeedItems(cached, &items); err == nil {
				return &FeedResponse{Items: items, HasMore: false, Next: 0}, nil
			}
		}
	}

	videos, err := fs.videoRepo.FindLatest(ctx, req.Limit*2)
	if err != nil {
		return nil, err
	}

	items, err := fs.buildFeedItems(ctx, videos, req.AccountID)
	if err != nil {
		return nil, err
	}

	items = fs.sortByHotScore(items)[:min(req.Limit, len(items))]

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if encoded, err := encodeFeedItems(items); err == nil {
			go fs.cache.Set(cacheCtx, cacheKey, encoded, HotTTL)
		}
	}

	return &FeedResponse{Items: items, HasMore: false, Next: 0}, nil
}

func (fs *FeedService) getFollowingFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	if req.AccountID == 0 {
		return &FeedResponse{Items: []FeedItem{}, HasMore: false, Next: 0}, nil
	}

	cacheKey := fs.cache.Key("feed:following:%d:%d", req.AccountID, req.Cursor)

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if cached, err := fs.cache.Get(cacheCtx, cacheKey); err == nil && cached != "" {
			var items []FeedItem
			if err := decodeFeedItems(cached, &items); err == nil {
				return &FeedResponse{Items: items, HasMore: true, Next: req.Cursor + int64(req.Limit)}, nil
			}
		}
	}

	followingIDs, err := fs.accountRepo.GetFollowing(ctx, req.AccountID)
	if err != nil {
		log.Printf("[WARN] [FeedService] failed to get following: %v", err)
		return &FeedResponse{Items: []FeedItem{}, HasMore: false, Next: 0}, nil
	}

	if len(followingIDs) == 0 {
		return &FeedResponse{Items: []FeedItem{}, HasMore: false, Next: 0}, nil
	}

	videos, err := fs.videoRepo.FindByAccountIDs(ctx, followingIDs, req.Limit, int(req.Cursor))
	if err != nil {
		return nil, err
	}

	items, err := fs.buildFeedItems(ctx, videos, req.AccountID)
	if err != nil {
		return nil, err
	}

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if encoded, err := encodeFeedItems(items); err == nil {
			go fs.cache.Set(cacheCtx, cacheKey, encoded, CacheTTL)
		}
	}

	return &FeedResponse{
		Items:   items,
		HasMore: len(items) == req.Limit,
		Next:    req.Cursor + int64(req.Limit),
	}, nil
}

func (fs *FeedService) getTagFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	if req.Tag == "" {
		return fs.getLatestFeed(ctx, req)
	}

	cacheKey := fs.cache.Key("feed:tag:%s:%d", req.Tag, req.Cursor)

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if cached, err := fs.cache.Get(cacheCtx, cacheKey); err == nil && cached != "" {
			var items []FeedItem
			if err := decodeFeedItems(cached, &items); err == nil {
				return &FeedResponse{Items: items, HasMore: true, Next: req.Cursor + int64(req.Limit)}, nil
			}
		}
	}

	videos, err := fs.videoRepo.FindByTag(ctx, req.Tag, req.Limit, int(req.Cursor))
	if err != nil {
		logger.Warn("FeedService", "failed to find by tag: %v", err)
		return nil, err
	}

	items, err := fs.buildFeedItems(ctx, videos, req.AccountID)
	if err != nil {
		return nil, err
	}

	if fs.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		if encoded, err := encodeFeedItems(items); err == nil {
			go fs.cache.Set(cacheCtx, cacheKey, encoded, CacheTTL)
		}
	}

	return &FeedResponse{
		Items:   items,
		HasMore: len(items) == req.Limit,
		Next:    req.Cursor + int64(req.Limit),
	}, nil
}

func (fs *FeedService) searchFeed(ctx context.Context, req FeedRequest) (*FeedResponse, error) {
	if req.Tag == "" {
		return &FeedResponse{Items: []FeedItem{}, HasMore: false, Next: 0}, nil
	}

	videos, err := fs.videoRepo.Search(ctx, req.Tag, req.Limit, int(req.Cursor))
	if err != nil {
		logger.Warn("FeedService", "failed to search videos: %v", err)
		return nil, err
	}

	items, err := fs.buildFeedItems(ctx, videos, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &FeedResponse{
		Items:   items,
		HasMore: len(items) == req.Limit,
		Next:    req.Cursor + int64(req.Limit),
	}, nil
}

func (fs *FeedService) buildFeedItems(ctx context.Context, videos []*video.Video, accountID uint) ([]FeedItem, error) {
	if len(videos) == 0 {
		return []FeedItem{}, nil
	}

	accountIDMap := make(map[uint]bool)
	for _, v := range videos {
		accountIDMap[v.AccountID] = true
	}

	accountIDs := make([]uint, 0, len(accountIDMap))
	for id := range accountIDMap {
		accountIDs = append(accountIDs, id)
	}

	accounts, err := fs.accountRepo.FindByIDs(ctx, accountIDs)
	if err != nil {
		logger.Warn("FeedService", "failed to find accounts: %v", err)
		return []FeedItem{}, nil
	}

	accountMap := make(map[uint]*account.Account)
	for _, acc := range accounts {
		accountMap[acc.ID] = acc
	}

	items := make([]FeedItem, 0, len(videos))
	for _, v := range videos {
		acc, exists := accountMap[v.AccountID]
		if !exists {
			continue
		}

		isLiked := false
		if accountID > 0 {
			isLiked, _ = fs.likeRepo.ExistsLike(ctx, accountID, v.ID)
		}

		items = append(items, FeedItem{
			VideoID:      v.ID,
			AccountID:    v.AccountID,
			Username:     acc.Username,
			AvatarURL:    acc.AvatarURL,
			Title:        v.Title,
			VideoURL:     fs.baseURL + "/" + v.VideoPath,
			CoverURL:     fs.baseURL + "/" + v.CoverPath,
			Duration:     v.Duration,
			Description:  v.Description,
			Tags:         v.Tags,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			CreatedAt:    v.CreatedAt,
			IsLiked:      isLiked,
		})
	}
	return items, nil
}

func (fs *FeedService) calculateHotScore(item FeedItem) int64 {
	age := time.Since(item.CreatedAt).Hours()
	decay := 1 / (1 + age/24)
	return int64(float64(item.LikeCount*2+item.CommentCount*5+item.ViewCount) * decay)
}

func (fs *FeedService) sortByHotScore(items []FeedItem) []FeedItem {
	sort.Slice(items, func(i, j int) bool {
		return fs.calculateHotScore(items[i]) > fs.calculateHotScore(items[j])
	})
	return items
}

func containsTag(tags, tag string) bool {
	if tags == "" || tag == "" {
		return false
	}
	for _, t := range strings.Split(tags, ",") {
		if t == tag {
			return true
		}
	}
	return false
}

func encodeFeedItems(items []FeedItem) (string, error) {
	data, err := json.Marshal(items)
	return string(data), err
}

func decodeFeedItems(data string, items *[]FeedItem) error {
	return json.Unmarshal([]byte(data), items)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}