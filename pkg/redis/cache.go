package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache Redis 缓存实现
type Cache struct {
	client *Client
	prefix string
}

// NewCache 创建 Redis 缓存
func NewCache(client *Client, prefix string) *Cache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if !c.client.IsEnabled() {
		return nil
	}

	key = c.prefix + key
	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	return c.client.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	if !c.client.IsEnabled() {
		return fmt.Errorf("redis not enabled")
	}

	key = c.prefix + key
	data, err := c.client.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return err
	}

	switch dest := dest.(type) {
	case *[]byte:
		*dest = []byte(data)
	case *string:
		*dest = data
	default:
		return json.Unmarshal([]byte(data), dest)
	}

	return nil
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, key string) error {
	if !c.client.IsEnabled() {
		return nil
	}

	key = c.prefix + key
	return c.client.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	if !c.client.IsEnabled() {
		return false, nil
	}

	key = c.prefix + key
	result, err := c.client.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Clear 清空缓存
func (c *Cache) Clear(ctx context.Context) error {
	if !c.client.IsEnabled() {
		return nil
	}

	pattern := c.prefix + "*"
	keys, err := c.client.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.client.Del(ctx, keys...).Err()
	}

	return nil
}

// GetTTL 获取键的过期时间
func (c *Cache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	if !c.client.IsEnabled() {
		return 0, fmt.Errorf("redis not enabled")
	}

	key = c.prefix + key
	return c.client.client.TTL(ctx, key).Result()
}

// SetNX 设置键值，仅当键不存在时
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if !c.client.IsEnabled() {
		return false, nil
	}

	key = c.prefix + key
	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	return c.client.client.SetNX(ctx, key, data, expiration).Result()
}
