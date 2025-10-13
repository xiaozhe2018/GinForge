package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock Redis 分布式锁
type Lock struct {
	client *Client
	key    string
	value  string
	ttl    time.Duration
}

// NewLock 创建分布式锁
func NewLock(client *Client, key string, ttl time.Duration) *Lock {
	return &Lock{
		client: client,
		key:    key,
		value:  fmt.Sprintf("%d", time.Now().UnixNano()),
		ttl:    ttl,
	}
}

// Acquire 获取锁
func (l *Lock) Acquire(ctx context.Context) (bool, error) {
	if !l.client.IsEnabled() {
		return true, nil // 如果 Redis 未启用，直接返回成功
	}

	result, err := l.client.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

// Release 释放锁
func (l *Lock) Release(ctx context.Context) error {
	if !l.client.IsEnabled() {
		return nil
	}

	// 使用 Lua 脚本确保只有锁的持有者才能释放锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not owned by this client")
	}

	return nil
}

// Refresh 刷新锁的过期时间
func (l *Lock) Refresh(ctx context.Context) error {
	if !l.client.IsEnabled() {
		return nil
	}

	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.client.Eval(ctx, script, []string{l.key}, l.value, int(l.ttl.Seconds())).Result()
	if err != nil {
		return fmt.Errorf("failed to refresh lock: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not owned by this client")
	}

	return nil
}

// Status 检查锁状态
func (l *Lock) Status(ctx context.Context) (bool, error) {
	if !l.client.IsEnabled() {
		return true, nil
	}

	value, err := l.client.client.Get(ctx, l.key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return value == l.value, nil
}

// TryLock 尝试获取锁，带超时
func (l *Lock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
			acquired, err := l.Acquire(ctx)
			if err != nil {
				return false, err
			}
			if acquired {
				return true, nil
			}
		}
	}
}

// WithLock 执行带锁的操作
func (l *Lock) WithLock(ctx context.Context, fn func() error) error {
	acquired, err := l.Acquire(ctx)
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("failed to acquire lock")
	}

	defer func() {
		if releaseErr := l.Release(ctx); releaseErr != nil {
			l.client.logger.Error("failed to release lock", releaseErr, "key", l.key)
		}
	}()

	return fn()
}
