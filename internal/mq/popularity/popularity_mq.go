package mqpopularity

import (
	"encoding/json"
	"log"

	"feedsystem_video_go/internal/mq/rabbitmq"
)

const PopularityExchange = "popularity_exchange"
const PopularityQueue = "popularity_queue"

type PopularityMessage struct {
	VideoID   uint   `json:"video_id"`
	AccountID uint   `json:"account_id"`
	Action    string `json:"action"` // "view", "like", "comment", "share"
	Score     int    `json:"score"`  // 权重分数
}

type PopularityMQ struct {
	client *rabbitmq.Client
}

func NewPopularityMQ(client *rabbitmq.Client) (*PopularityMQ, error) {
	err := client.DeclareExchange(PopularityExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = client.DeclareQueue(PopularityQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = client.BindQueue(PopularityQueue, "popularity", PopularityExchange, false, nil)
	if err != nil {
		return nil, err
	}

	return &PopularityMQ{client: client}, nil
}

func (mq *PopularityMQ) Publish(msg PopularityMessage) error {
	return mq.client.PublishJSON(PopularityExchange, "popularity", msg)
}

func (mq *PopularityMQ) PublishView(videoID, accountID uint) error {
	return mq.Publish(PopularityMessage{
		VideoID:   videoID,
		AccountID: accountID,
		Action:    "view",
		Score:     1,
	})
}

func (mq *PopularityMQ) PublishLike(videoID, accountID uint) error {
	return mq.Publish(PopularityMessage{
		VideoID:   videoID,
		AccountID: accountID,
		Action:    "like",
		Score:     5,
	})
}

func (mq *PopularityMQ) PublishComment(videoID, accountID uint) error {
	return mq.Publish(PopularityMessage{
		VideoID:   videoID,
		AccountID: accountID,
		Action:    "comment",
		Score:     10,
	})
}

func (mq *PopularityMQ) Consume(handler func(PopularityMessage)) error {
	msgs, err := mq.client.Consume(PopularityQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg PopularityMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal popularity message: %v", err)
				d.Ack(false)
				continue
			}

			handler(msg)
			d.Ack(false)
		}
	}()

	return nil
}
