package feed

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const TimelineKey = "timeline:global"

func (fs *FeedService) AddToTimeline(ctx context.Context, videoID uint) error {
	if fs.cache == nil {
		return nil
	}

	timestamp := time.Now().Unix()

	cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return fs.cache.ZAdd(cacheCtx, TimelineKey, redis.Z{
		Score:  float64(timestamp),
		Member: strconv.FormatUint(uint64(videoID), 10),
	})
}

func (fs *FeedService) GetFromTimeline(ctx context.Context, offset, limit int64) ([]uint, error) {
	if fs.cache == nil {
		return nil, nil
	}

	cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	members, err := fs.cache.ZRevRange(cacheCtx, TimelineKey, offset, offset+limit-1)
	if err != nil {
		return nil, err
	}

	videoIDs := make([]uint, 0, len(members))
	for _, m := range members {
		id, err := strconv.ParseUint(m, 10, 64)
		if err == nil {
			videoIDs = append(videoIDs, uint(id))
		}
	}

	return videoIDs, nil
}

func (fs *FeedService) RemoveFromTimeline(ctx context.Context, videoID uint) error {
	if fs.cache == nil {
		return nil
	}

	cacheCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return fs.cache.ZRem(cacheCtx, TimelineKey, strconv.FormatUint(uint64(videoID), 10))
}
