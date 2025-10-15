package redis

import (
	"context"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"
)

// Manager Redis 管理器，整合缓存、锁和消息队列
type Manager struct {
	client *Client
	cache  *Cache
	queue  *RedisQueue
}

// NewManager 创建 Redis 管理器
func NewManager(cfg *config.RedisConfig, log logger.Logger) *Manager {
	client := NewClient(cfg, log)

	return &Manager{
		client: client,
		cache:  NewCache(client, "cache:"),
		queue:  NewQueue(client, ""),
	}
}

// GetClient 获取 Redis 客户端
func (m *Manager) GetClient() *Client {
	return m.client
}

// GetCache 获取缓存实例
func (m *Manager) GetCache() *Cache {
	return m.cache
}

// GetQueue 获取消息队列实例
func (m *Manager) GetQueue() *RedisQueue {
	return m.queue
}

// NewLock 创建分布式锁
func (m *Manager) NewLock(key string, ttl time.Duration) *Lock {
	return NewLock(m.client, key, ttl)
}

// Ping 测试连接
func (m *Manager) Ping(ctx context.Context) error {
	return m.client.Ping(ctx).Err()
}

// Close 关闭连接
func (m *Manager) Close() error {
	return m.client.Close()
}

// IsEnabled 检查 Redis 是否启用
func (m *Manager) IsEnabled() bool {
	return m.client.IsEnabled()
}

// 缓存方法快捷访问
func (m *Manager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.cache.Set(ctx, key, value, expiration)
}

func (m *Manager) Get(ctx context.Context, key string, dest interface{}) error {
	return m.cache.Get(ctx, key, dest)
}

func (m *Manager) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	return m.cache.Exists(ctx, key)
}

func (m *Manager) Clear(ctx context.Context) error {
	return m.cache.Clear(ctx)
}

// 消息队列方法快捷访问
func (m *Manager) Publish(ctx context.Context, topic string, data map[string]interface{}) error {
	return m.queue.Publish(ctx, topic, data)
}

func (m *Manager) PublishWithDelay(ctx context.Context, topic string, data map[string]interface{}, delay time.Duration) error {
	return m.queue.PublishWithDelay(ctx, topic, data, delay)
}

func (m *Manager) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	return m.queue.Subscribe(ctx, topic, handler)
}

func (m *Manager) GetQueueLength(ctx context.Context, topic string) (int64, error) {
	return m.queue.GetQueueLength(ctx, topic)
}

func (m *Manager) PurgeQueue(ctx context.Context, topic string) error {
	return m.queue.PurgeQueue(ctx, topic)
}

// 分布式锁方法快捷访问
func (m *Manager) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	lock := m.NewLock(key, ttl)
	return lock.WithLock(ctx, fn)
}

func (m *Manager) TryLock(ctx context.Context, key string, ttl time.Duration, timeout time.Duration) (*Lock, bool, error) {
	lock := m.NewLock(key, ttl)
	acquired, err := lock.TryLock(ctx, timeout)
	if err != nil {
		return nil, false, err
	}
	if !acquired {
		return nil, false, nil
	}
	return lock, true, nil
}