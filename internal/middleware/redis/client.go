package redis

import (
	"context"
	"feedsystem_video_go/internal/config"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	client *goredis.Client
	prefix string
}

// NewClient 创建 Redis 客户端实例
func NewClient(cfg config.RedisConfig) *Client {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{client: client, prefix: "feedsystem:"}
}

// Key 生成带前缀的 Redis key
func (c *Client) Key(format string, args ...interface{}) string {
	return c.prefix + fmt.Sprintf(format, args...)
}

// Set 设置字符串类型的键值对
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// SetBytes 设置字节数组类型的键值对
func (c *Client) SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get 获取字符串类型的键值对
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// GetBytes 获取字节数组类型的键值对
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

// Del 删除键值对
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Incr 增加键值对的数值
func (c *Client) Incr(ctx context.Context, key string) error {
	return c.client.Incr(ctx, key).Err()
}

// Decr 减少键值对的数值
func (c *Client) Decr(ctx context.Context, key string) error {
	return c.client.Decr(ctx, key).Err()
}

// ZAdd 添加有序集合元素/goredis.ZRedis：有序集合的成员（带分数）
func (c *Client) ZAdd(ctx context.Context, key string, members ...goredis.Z) error {
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合元素（按分数升序）
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合元素（按分数降序）
func (c *Client) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRem 删除有序集合元素
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Expire 设置键的过期时间
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// Ping 发送 Ping 请求
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close 关闭 Redis 连接
func (c *Client) Close() error {
	return c.client.Close()
}
