package comment

import (
	"encoding/json"
	"log"

	"feedsystem_video_go/internal/mq/rabbitmq"
)

const CommentExchange = "comment_exchange"
const CommentQueue = "comment_queue"

type CommentMessage struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id"`
	VideoID   uint   `json:"video_id"`
	Content   string `json:"content"`
	ReplyTo   uint   `json:"reply_to,omitempty"`
	Mentioned string `json:"mentioned,omitempty"`
	Action    string `json:"action"` // "create" or "delete"
}

type CommentMQ struct {
	client *rabbitmq.Client
}

func NewCommentMQ(client *rabbitmq.Client) (*CommentMQ, error) {
	err := client.DeclareExchange(CommentExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = client.DeclareQueue(CommentQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = client.BindQueue(CommentQueue, "comment", CommentExchange, false, nil)
	if err != nil {
		return nil, err
	}

	return &CommentMQ{client: client}, nil
}

func (mq *CommentMQ) PublishCreate(comment CommentMessage) error {
	comment.Action = "create"
	return mq.client.PublishJSON(CommentExchange, "comment", comment)
}

func (mq *CommentMQ) PublishDelete(comment CommentMessage) error {
	comment.Action = "delete"
	return mq.client.PublishJSON(CommentExchange, "comment", comment)
}

func (mq *CommentMQ) Consume(handler func(CommentMessage)) error {
	msgs, err := mq.client.Consume(CommentQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg CommentMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal comment message: %v", err)
				d.Ack(false)
				continue
			}

			handler(msg)
			d.Ack(false)
		}
	}()

	return nil
}
