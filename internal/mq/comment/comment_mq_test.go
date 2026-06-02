package comment

import (
	"testing"
	"time"

	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/mq/rabbitmq"
)

func TestCommentMQ_PublishCreate(t *testing.T) {
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

	commentMQ, err := NewCommentMQ(client)
	if err != nil {
		t.Fatalf("Failed to create CommentMQ: %v", err)
	}

	receivedChan := make(chan CommentMessage, 10)
	err = commentMQ.Consume(func(msg CommentMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	testComment := CommentMessage{
		ID:        1,
		AccountID: 100,
		VideoID:   200,
		Content:   "测试评论内容",
		ReplyTo:   0,
		Mentioned: "",
	}

	err = commentMQ.PublishCreate(testComment)
	if err != nil {
		t.Fatalf("Failed to publish create message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "create" {
			t.Errorf("Action不匹配: 期望 create, 实际 %s", msg.Action)
		}
		if msg.ID != testComment.ID || msg.AccountID != testComment.AccountID {
			t.Errorf("消息内容不匹配: 期望 %+v, 实际 %+v", testComment, msg)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestCommentMQ_PublishDelete(t *testing.T) {
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

	commentMQ, err := NewCommentMQ(client)
	if err != nil {
		t.Fatalf("Failed to create CommentMQ: %v", err)
	}

	receivedChan := make(chan CommentMessage, 10)
	err = commentMQ.Consume(func(msg CommentMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	testComment := CommentMessage{
		ID:        2,
		AccountID: 101,
		VideoID:   201,
		Content:   "要删除的评论",
	}

	err = commentMQ.PublishDelete(testComment)
	if err != nil {
		t.Fatalf("Failed to publish delete message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.Action != "delete" {
			t.Errorf("Action不匹配: 期望 delete, 实际 %s", msg.Action)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestCommentMQ_ReplyToComment(t *testing.T) {
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

	commentMQ, err := NewCommentMQ(client)
	if err != nil {
		t.Fatalf("Failed to create CommentMQ: %v", err)
	}

	receivedChan := make(chan CommentMessage, 10)
	err = commentMQ.Consume(func(msg CommentMessage) {
		receivedChan <- msg
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	replyComment := CommentMessage{
		ID:        3,
		AccountID: 102,
		VideoID:   202,
		Content:   "回复评论",
		ReplyTo:   1,
		Mentioned: "@user1",
	}

	err = commentMQ.PublishCreate(replyComment)
	if err != nil {
		t.Fatalf("Failed to publish reply message: %v", err)
	}

	select {
	case msg := <-receivedChan:
		t.Logf("收到消息: %+v", msg)
		if msg.ReplyTo != 1 || msg.Mentioned != "@user1" {
			t.Errorf("回复信息不匹配: 期望 ReplyTo=1, Mentioned=@user1, 实际 ReplyTo=%d, Mentioned=%s", msg.ReplyTo, msg.Mentioned)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}

func TestCommentMQ_Performance(t *testing.T) {
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

	commentMQ, err := NewCommentMQ(client)
	if err != nil {
		t.Fatalf("Failed to create CommentMQ: %v", err)
	}

	msgCount := 0
	start := time.Now()
	commentMQ.Consume(func(msg CommentMessage) {
		msgCount++
	})

	go func() {
		for i := 0; i < 1000; i++ {
			comment := CommentMessage{
				ID:        uint(i),
				AccountID: uint(i % 100),
				VideoID:   uint(i % 50),
				Content:   "性能测试评论",
			}
			commentMQ.PublishCreate(comment)
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
