package mqlike

import (
	"encoding/json"
	"log"

	"feedsystem_video_go/internal/mq/rabbitmq"
)

const LikeExchange = "like_exchange"
const LikeQueue = "like_queue"

type LikeMessage struct {
	AccountID uint   `json:"account_id"`
	VideoID   uint   `json:"video_id"`
	Action    string `json:"action"` // "like" or "unlike"
}

type LikeMQ struct {
	client *rabbitmq.Client
}

func NewLikeMQ(client *rabbitmq.Client) (*LikeMQ, error) {
	err := client.DeclareExchange(LikeExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = client.DeclareQueue(LikeQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = client.BindQueue(LikeQueue, "like", LikeExchange, false, nil)
	if err != nil {
		return nil, err
	}

	return &LikeMQ{client: client}, nil
}

func (mq *LikeMQ) PublishLike(accountID, videoID uint) error {
	msg := LikeMessage{
		AccountID: accountID,
		VideoID:   videoID,
		Action:    "like",
	}

	return mq.client.PublishJSON(LikeExchange, "like", msg)
}

func (mq *LikeMQ) PublishUnlike(accountID, videoID uint) error {
	msg := LikeMessage{
		AccountID: accountID,
		VideoID:   videoID,
		Action:    "unlike",
	}

	return mq.client.PublishJSON(LikeExchange, "like", msg)
}

func (mq *LikeMQ) Consume(handler func(LikeMessage)) error {
	msgs, err := mq.client.Consume(LikeQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg LikeMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal like message: %v", err)
				d.Ack(false)
				continue
			}

			handler(msg)
			d.Ack(false)
		}
	}()

	return nil
}
