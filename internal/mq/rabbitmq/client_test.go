package rabbitmq

import (
	"testing"
	"time"

	"feedsystem_video_go/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

func TestClient_NewClient(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	if client == nil {
		t.Fatal("Client should not be nil")
	}

	if client.conn == nil {
		t.Error("Connection should not be nil")
	}

	if client.channel == nil {
		t.Error("Channel should not be nil")
	}

	t.Log("Client created successfully")
}

func TestClient_DeclareQueue(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	queueName := "test_queue"
	queue, err := client.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	if queue.Name != queueName {
		t.Errorf("Queue name mismatch: expected %s, got %s", queueName, queue.Name)
	}

	t.Logf("Queue declared: %+v", queue)
}

func TestClient_DeclareExchange(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	exchangeName := "test_exchange"
	err = client.DeclareExchange(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	t.Logf("Exchange declared: %s", exchangeName)
}

func TestClient_BindQueue(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	exchangeName := "test_bind_exchange"
	queueName := "test_bind_queue"
	routingKey := "test_routing_key"

	err = client.DeclareExchange(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = client.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	err = client.BindQueue(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		t.Fatalf("Failed to bind queue: %v", err)
	}

	t.Logf("Queue %s bound to exchange %s with routing key %s", queueName, exchangeName, routingKey)
}

func TestClient_PublishAndConsume(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	exchangeName := "test_pubsub_exchange"
	queueName := "test_pubsub_queue"
	routingKey := "test_pubsub_routing"

	err = client.DeclareExchange(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = client.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	err = client.BindQueue(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		t.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := client.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	testMessage := []byte("test message content")
	err = client.Publish(exchangeName, routingKey, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        testMessage,
	})
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	select {
	case msg := <-msgs:
		if string(msg.Body) != string(testMessage) {
			t.Errorf("Message body mismatch: expected %s, got %s", string(testMessage), string(msg.Body))
		}
		msg.Ack(false)
		t.Logf("Message received: %s", string(msg.Body))
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestClient_PublishJSON(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	exchangeName := "test_json_exchange"
	queueName := "test_json_queue"
	routingKey := "test_json_routing"

	err = client.DeclareExchange(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = client.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	err = client.BindQueue(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		t.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := client.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	type TestMessage struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testMsg := TestMessage{Name: "test", Value: 123}
	err = client.PublishJSON(exchangeName, routingKey, testMsg)
	if err != nil {
		t.Fatalf("Failed to publish JSON message: %v", err)
	}

	select {
	case msg := <-msgs:
		if msg.ContentType != "application/json" {
			t.Errorf("Content-Type mismatch: expected application/json, got %s", msg.ContentType)
		}
		t.Logf("JSON message received: %s", string(msg.Body))
		msg.Ack(false)
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestClient_Ping(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	err = client.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	t.Log("Ping successful")

	client.Close()

	err = client.Ping()
	if err == nil {
		t.Error("Ping should fail after connection closed")
	}

	t.Log("Ping correctly failed after connection closed")
}

func TestClient_Close(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Failed to close client: %v", err)
	}

	t.Log("Client closed successfully")
}

func TestClient_Performance(t *testing.T) {
	cfg := config.RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "password123",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	exchangeName := "test_perf_exchange"
	queueName := "test_perf_queue"
	routingKey := "test_perf_routing"

	err = client.DeclareExchange(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = client.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	err = client.BindQueue(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		t.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := client.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	msgCount := 0
	start := time.Now()

	go func() {
		for range msgs {
			msgCount++
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			client.Publish(exchangeName, routingKey, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte("test"),
			})
		}
	}()

	for msgCount < 10000 {
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	tps := float64(msgCount) / elapsed.Seconds()
	t.Logf("吞吐量: %.2f msg/sec", tps)
	t.Logf("总耗时: %v", elapsed)
}
