package mqpopularity

import (
	"testing"
	"time"

	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/mq/rabbitmq"
)

func TestPopularityMQ_PublishView(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	receivedChan := make(chan PopularityMessage, 10)
	err = popularityMQ.Consume(func(msg PopularityMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	err = popularityMQ.PublishView(100, 1)
	if err != nil {
		t.Fatalf("Failed to publish view message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "view" {
			t.Errorf("Action不匹配: 期望 view, 实际 %s", msg.Action)
		}
		if msg.Score != 1 {
			t.Errorf("Score不匹配: 期望 1, 实际 %d", msg.Score)
		}
		if msg.VideoID != 100 || msg.AccountID != 1 {
			t.Errorf("消息内容不匹配: 期望 VideoID=100, AccountID=1, 实际 VideoID=%d, AccountID=%d", msg.VideoID, msg.AccountID)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestPopularityMQ_PublishLike(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	receivedChan := make(chan PopularityMessage, 10)
	err = popularityMQ.Consume(func(msg PopularityMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	err = popularityMQ.PublishLike(200, 2)
	if err != nil {
		t.Fatalf("Failed to publish like message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "like" {
			t.Errorf("Action不匹配: 期望 like, 实际 %s", msg.Action)
		}
		if msg.Score != 5 {
			t.Errorf("Score不匹配: 期望 5, 实际 %d", msg.Score)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestPopularityMQ_PublishComment(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	receivedChan := make(chan PopularityMessage, 10)
	err = popularityMQ.Consume(func(msg PopularityMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	err = popularityMQ.PublishComment(300, 3)
	if err != nil {
		t.Fatalf("Failed to publish comment message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "comment" {
			t.Errorf("Action不匹配: 期望 comment, 实际 %s", msg.Action)
		}
		if msg.Score != 10 {
			t.Errorf("Score不匹配: 期望 10, 实际 %d", msg.Score)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestPopularityMQ_PublishCustom(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	receivedChan := make(chan PopularityMessage, 10)
	err = popularityMQ.Consume(func(msg PopularityMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	customMsg := PopularityMessage{
		VideoID:   400,
		AccountID: 4,
		Action:    "share",
		Score:     20,
	}

	err = popularityMQ.Publish(customMsg)
	if err != nil {
		t.Fatalf("Failed to publish custom message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "share" || msg.Score != 20 {
			t.Errorf("自定义消息不匹配: 期望 Action=share, Score=20, 实际 Action=%s, Score=%d", msg.Action, msg.Score)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestPopularityMQ_ScoreWeights(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	receivedChan := make(chan PopularityMessage, 10)
	popularityMQ.Consume(func(msg PopularityMessage) {
		receivedChan <- msg
	})

	testCases := []struct {
		name          string
		videoID       uint
		accountID     uint
		action        string
		expectedScore int
	}{
		{"view", 1, 1, "view", 1},
		{"like", 2, 2, "like", 5},
		{"comment", 3, 3, "comment", 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch tc.action {
			case "view":
				err = popularityMQ.PublishView(tc.videoID, tc.accountID)
			case "like":
				err = popularityMQ.PublishLike(tc.videoID, tc.accountID)
			case "comment":
				err = popularityMQ.PublishComment(tc.videoID, tc.accountID)
			}

			if err != nil {
				t.Fatalf("Failed to publish %s message: %v", tc.action, err)
			}

			select {
			case msg := <-receivedChan:
				if msg.VideoID != tc.videoID || msg.Action != tc.action {
					t.Logf("收到其他消息: %+v, 等待正确消息", msg)
					select {
					case msg := <-receivedChan:
						if msg.Score != tc.expectedScore {
							t.Errorf("Score不匹配: 期望 %d, 实际 %d", tc.expectedScore, msg.Score)
						}
					case <-time.After(3 * time.Second):
						t.Error("超时未收到正确消息")
					}
				} else if msg.Score != tc.expectedScore {
					t.Errorf("Score不匹配: 期望 %d, 实际 %d", tc.expectedScore, msg.Score)
				}
			case <-time.After(3 * time.Second):
				t.Error("超时未收到消息")
			}
		})
	}
}

func TestPopularityMQ_Performance(t *testing.T) {
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

	popularityMQ, err := NewPopularityMQ(client)
	if err != nil {
		t.Fatalf("Failed to create PopularityMQ: %v", err)
	}

	msgCount := 0
	start := time.Now()
	popularityMQ.Consume(func(msg PopularityMessage) {
		msgCount++
	})

	go func() {
		for i := 0; i < 1000; i++ {
			switch i % 3 {
			case 0:
				popularityMQ.PublishView(uint(i), uint(i%100))
			case 1:
				popularityMQ.PublishLike(uint(i), uint(i%100))
			case 2:
				popularityMQ.PublishComment(uint(i), uint(i%100))
			}
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
