package worker

import (
	"context"
	"feedsystem_video_go/internal/social"
	"feedsystem_video_go/internal/sse"
	"log"
)

type SocialWorker struct {
	socialService *social.SocialService
	sseHub        *sse.Hub
}

func NewSocialWorker(socialService *social.SocialService, sseHub *sse.Hub) *SocialWorker {
	return &SocialWorker{
		socialService: socialService,
		sseHub:        sseHub,
	}
}

func (sw *SocialWorker) Run(ctx context.Context) {
	log.Println("[SocialWorker] starting...")
	log.Println("[SocialWorker] polling for follow events...")
}
