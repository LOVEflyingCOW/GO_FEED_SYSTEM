package worker

import (
	"context"
	"feedsystem_video_go/internal/like"
	mqlike "feedsystem_video_go/internal/mq/like"
	"log"
)

type LikeWorker struct {
	likeService *like.LikeService
	likeMQ      *mqlike.LikeMQ
}

func NewLikeWorker(likeService *like.LikeService, likeMQ *mqlike.LikeMQ) *LikeWorker {
	return &LikeWorker{
		likeService: likeService,
		likeMQ:      likeMQ,
	}
}

func (lw *LikeWorker) Run(ctx context.Context) {
	log.Println("[LikeWorker] starting...")
	err := lw.likeMQ.Consume(func(msg mqlike.LikeMessage) {
		log.Printf("[LikeWorker] received: account=%d, video=%d, action=%s", msg.AccountID, msg.VideoID, msg.Action)

		switch msg.Action {
		case "like":
			_, err := lw.likeService.LikeVideo(ctx, msg.AccountID, msg.VideoID)
			if err != nil {
				log.Printf("[LikeWorker] error liking video: %v", err)
			}
		case "unlike":
			_, err := lw.likeService.UnlikeVideo(ctx, msg.AccountID, msg.VideoID)
			if err != nil {
				log.Printf("[LikeWorker] error unliking video: %v", err)
			}
		}
	})

	if err != nil {
		log.Printf("[LikeWorker] error starting consumer: %v", err)
	}
}
