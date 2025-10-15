package redis

import (
	"context"
	"fmt"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端封装
type Client struct {
	client redis.UniversalClient
	logger logger.Logger
	config *config.RedisConfig
	prefix string
}

// NewClient 创建 Redis 客户端
func NewClient(cfg *config.RedisConfig, log logger.Logger) *Client {
	if !cfg.Enabled {
		return &Client{
			client: nil,
			logger: log,
			config: cfg,
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	return &Client{
		client: client,
		logger: log,
		config: cfg,
		prefix: "",
	}
}

// GetClient 获取原始 Redis 客户端
func (c *Client) GetClient() redis.UniversalClient {
	return c.client
}

// Ping 测试连接
func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	if c.client == nil {
		return redis.NewStatusCmd(ctx)
	}
	return c.client.Ping(ctx)
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// IsEnabled 检查 Redis 是否启用
func (c *Client) IsEnabled() bool {
	return c.config.Enabled && c.client != nil
}

// Set 设置键值对
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}
	return c.client.Get(ctx, key).Result()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}
	result, err := c.client.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if c.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return c.client.Del(ctx, keys...).Err()
}

// FlushDB 清空当前数据库
func (c *Client) FlushDB(ctx context.Context) *redis.StatusCmd {
	if c.client == nil {
		return redis.NewStatusCmd(ctx)
	}
	return c.client.FlushDB(ctx)
}

// SCard 获取集合的基数(元素数量)
func (c *Client) SCard(ctx context.Context, key string) *redis.IntCmd {
	if c.client == nil {
		return redis.NewIntCmd(ctx)
	}
	return c.client.SCard(ctx, key)
}

// Incr 增加键的整数值
func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	if c.client == nil {
		return redis.NewIntCmd(ctx)
	}
	return c.client.Incr(ctx, key)
}

// Expire 设置键的过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	if c.client == nil {
		return redis.NewBoolCmd(ctx)
	}
	return c.client.Expire(ctx, key, expiration)
}

// TTL 获取键的剩余生存时间
func (c *Client) TTL(ctx context.Context, key string) *redis.DurationCmd {
	if c.client == nil {
		return redis.NewDurationCmd(ctx, time.Second, "ttl", key)
	}
	return c.client.TTL(ctx, key)
}