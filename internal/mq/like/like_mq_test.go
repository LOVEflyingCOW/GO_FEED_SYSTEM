package mqlike

import (
	"testing"
	"time"

	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/mq/rabbitmq"
)

func TestLikeMQ_PublishAndConsume(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	likeMQ, err := NewLikeMQ(client)
	if err != nil {
		t.Fatalf("Failed to create LikeMQ: %v", err)
	}

	receivedChan := make(chan LikeMessage, 10)
	err = likeMQ.Consume(func(msg LikeMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	testCases := []struct {
		name      string
		accountID uint
		videoID   uint
		action    string
	}{
		{"like", 1, 100, "like"},
		{"unlike", 2, 200, "unlike"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.action == "like" {
				err = likeMQ.PublishLike(tc.accountID, tc.videoID)
			} else {
				err = likeMQ.PublishUnlike(tc.accountID, tc.videoID)
			}
			if err != nil {
				t.Fatalf("Failed to publish message: %v", err)
			}

			select {
			case msg := <-receivedChan:
				t.Logf("收到消息: %+v", msg)
				if msg.AccountID != tc.accountID || msg.VideoID != tc.videoID || msg.Action != tc.action {
					t.Errorf("消息内容不匹配: 期望 %+v, 实际 %+v", tc, msg)
				}
			case <-time.After(3 * time.Second):
				t.Error("超时未收到消息")
			}
		})
	}
}

func TestLikeMQ_Performance(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	likeMQ, err := NewLikeMQ(client)
	if err != nil {
		t.Fatalf("Failed to create LikeMQ: %v", err)
	}

	msgCount := 0
	start := time.Now()
	likeMQ.Consume(func(msg LikeMessage) {
		msgCount++
	})

	go func() {
		for i := 0; i < 1000; i++ {
			likeMQ.PublishLike(uint(i), uint(i*10))
		}
	}()

	for msgCount < 1000 {
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	tps := float64(msgCount) / elapsed.Seconds()
	t.Logf("吞吐量: %.2f msg/sec", tps)
	t.Logf("总耗时: %v", elapsed)
}
