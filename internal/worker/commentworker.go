package worker

import (
	"context"
	"feedsystem_video_go/internal/mq/comment"
	"feedsystem_video_go/internal/sse"
	"log"
)

type CommentWorker struct {
	sseHub    *sse.Hub
	commentMQ *comment.CommentMQ
}

func NewCommentWorker(sseHub *sse.Hub, commentMQ *comment.CommentMQ) *CommentWorker {
	return &CommentWorker{
		sseHub:    sseHub,
		commentMQ: commentMQ,
	}
}

func (cw *CommentWorker) Run(ctx context.Context) {
	log.Println("[CommentWorker] starting...")
	err := cw.commentMQ.Consume(func(msg comment.CommentMessage) {
		log.Printf("[CommentWorker] received: account=%d, video=%d, action=%s", msg.AccountID, msg.VideoID, msg.Action)

		notification := sse.Notification{
			Type:   "comment",
			FromID: msg.AccountID,
			From:   "",
			Content: map[string]interface{}{
				"video_id":   msg.VideoID,
				"comment_id": msg.ID,
				"content":    msg.Content,
			},
		}

		cw.sseHub.Broadcast(notification)
	})

	if err != nil {
		log.Printf("[CommentWorker] error starting consumer: %v", err)
	}
}
