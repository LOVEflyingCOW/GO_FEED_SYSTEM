package rabbitmq

import (
	"context"
	"encoding/json"
	"feedsystem_video_go/internal/config"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewClient(cfg config.RabbitMQConfig) (*Client, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, channel: channel}, nil
}

func (c *Client) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	return c.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *Client) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	return c.channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *Client) BindQueue(queueName, routingKey, exchangeName string, noWait bool, args amqp091.Table) error {
	return c.channel.QueueBind(queueName, routingKey, exchangeName, noWait, args)
}

func (c *Client) Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp091.Publishing) error {
	return c.channel.PublishWithContext(context.Background(), exchange, routingKey, mandatory, immediate, msg)
}

func (c *Client) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	return c.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *Client) PublishJSON(exchange, routingKey string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return c.Publish(exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        data,
		Expiration:  "60000",
	})
}

func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		log.Printf("warning: failed to close channel: %v", err)
	}
	return c.conn.Close()
}

func (c *Client) Ping() error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	if c.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}
	return nil
}
