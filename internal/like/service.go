package like

import (
	"context"
	"errors"
	"strconv"
	"time"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware/redis"
	"feedsystem_video_go/internal/pkg/logger"
	"feedsystem_video_go/internal/video"

	"gorm.io/gorm"
)

type LikeService struct {
	likeRepository  *LikeRepository
	videoRepository *video.VideoRepository
	cache           *redis.Client
}

func NewLikeService(likeRepository *LikeRepository, videoRepository *video.VideoRepository, cache *redis.Client) *LikeService {
	return &LikeService{
		likeRepository:  likeRepository,
		videoRepository: videoRepository,
		cache:           cache,
	}
}

// LikeVideo 点赞视频
func (ls *LikeService) LikeVideo(ctx context.Context, accountID, videoID uint) (*LikeResponse, error) {
	// 校验视频是否存在
	_, err := ls.videoRepository.FindByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrVideoNotFound
		}
		return nil, err
	}

	// 判断用户是否已经点过赞
	exists, err := ls.likeRepository.ExistsLike(ctx, accountID, videoID)
	if err != nil {
		return nil, err
	}

	// 已经点过赞：直接返回
	if exists {
		likeCount, _ := ls.likeRepository.CountByVideoID(ctx, videoID)
		return &LikeResponse{
			VideoID:   videoID,
			LikeCount: likeCount,
			IsLiked:   true,
		}, nil
	}

	// 没点过赞：走正常点赞流程
	like := &Like{
		AccountID: accountID,
		VideoID:   videoID,
	}

	// 插入点赞记录
	if err := ls.likeRepository.CreateLike(ctx, like); err != nil {
		return nil, err
	}

	// 视频数据库点赞数 +1
	if err := ls.videoRepository.IncreaseLikeCount(ctx, videoID); err != nil {
		logger.Warn("LikeService", "failed to increase like count: %v", err)
	}

	// Redis 缓存点赞数 +1
	if ls.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		_ = ls.cache.Incr(cacheCtx, ls.cache.Key("video:%d:likes", videoID))
	}

	// 返回最新点赞数
	likeCount, _ := ls.likeRepository.CountByVideoID(ctx, videoID)
	return &LikeResponse{
		VideoID:   videoID,
		LikeCount: likeCount,
		IsLiked:   true,
	}, nil
}

// UnlikeVideo 取消点赞
func (ls *LikeService) UnlikeVideo(ctx context.Context, accountID, videoID uint) (*LikeResponse, error) {
	if err := ls.likeRepository.DeleteLike(ctx, accountID, videoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrLikeNotFound
		}
		return nil, err
	}

	if err := ls.videoRepository.DecreaseLikeCount(ctx, videoID); err != nil {
		logger.Warn("LikeService", "failed to decrease like count: %v", err)
	}

	if ls.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		_ = ls.cache.Decr(cacheCtx, ls.cache.Key("video:%d:likes", videoID))
	}

	likeCount, _ := ls.likeRepository.CountByVideoID(ctx, videoID)

	return &LikeResponse{
		VideoID:   videoID,
		LikeCount: likeCount,
		IsLiked:   false,
	}, nil
}

// GetLikeStatus 获取点赞状态
func (ls *LikeService) GetLikeStatus(ctx context.Context, accountID, videoID uint) (*LikeResponse, error) {
	isLiked, err := ls.likeRepository.ExistsLike(ctx, accountID, videoID)
	if err != nil {
		return nil, err
	}

	var likeCount int64
	cacheHit := false

	if ls.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		countStr, err := ls.cache.Get(cacheCtx, ls.cache.Key("video:%d:likes", videoID))
		if err == nil {
			likeCount, _ = strconv.ParseInt(countStr, 10, 64)
			cacheHit = true
		}
	}

	if !cacheHit {
		likeCount, _ = ls.likeRepository.CountByVideoID(ctx, videoID)
	}

	return &LikeResponse{
		VideoID:   videoID,
		LikeCount: likeCount,
		IsLiked:   isLiked,
	}, nil
}

// ListLikes 获取用户点赞列表
func (ls *LikeService) ListLikes(ctx context.Context, accountID uint, cursor, limit int) (*LikeListResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}

	likes, total, err := ls.likeRepository.FindByAccountID(ctx, accountID, cursor, limit)
	if err != nil {
		return nil, err
	}

	videoIDs := make([]uint, 0, len(likes))
	for _, like := range likes {
		videoIDs = append(videoIDs, like.VideoID)
	}

	videos, err := ls.videoRepository.FindByIDs(ctx, videoIDs)
	if err != nil {
		return nil, err
	}

	videoMap := make(map[uint]*video.Video)
	for _, v := range videos {
		videoMap[v.ID] = v
	}

	items := make([]*VideoItem, 0, len(likes))
	for _, like := range likes {
		if v, ok := videoMap[like.VideoID]; ok {
			items = append(items, &VideoItem{
				ID:           v.ID,
				AccountID:    v.AccountID,
				Username:     v.Username,
				Title:        v.Title,
				PlayURL:      v.PlayURL,
				CoverURL:     v.CoverURL,
				Duration:     v.Duration,
				Description:  v.Description,
				Tags:         v.Tags,
				ViewCount:    v.ViewCount,
				LikeCount:    v.LikeCount,
				CommentCount: v.CommentCount,
				CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	return &LikeListResponse{
		Videos: items,
		Total:  total,
	}, nil
}
