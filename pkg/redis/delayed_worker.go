package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// DelayedWorker 延时消息处理器
type DelayedWorker struct {
	client   *Client
	topics   map[string]bool
	workers  map[string]context.CancelFunc
	mutex    sync.RWMutex
	interval time.Duration
}

// NewDelayedWorker 创建延时消息处理器
func NewDelayedWorker(client *Client, interval time.Duration) *DelayedWorker {
	return &DelayedWorker{
		client:   client,
		topics:   make(map[string]bool),
		workers:  make(map[string]context.CancelFunc),
		interval: interval,
	}
}

// StartTopicWorker 启动指定 topic 的延时消息处理器
func (dw *DelayedWorker) StartTopicWorker(ctx context.Context, topic string) error {
	if !dw.client.IsEnabled() {
		return nil
	}

	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	// 检查是否已经启动
	if dw.topics[topic] {
		return nil
	}

	// 创建子上下文
	workerCtx, cancel := context.WithCancel(ctx)
	dw.workers[topic] = cancel
	dw.topics[topic] = true

	// 启动处理器
	go dw.processDelayedMessages(workerCtx, topic)

	dw.client.logger.Info("delayed worker started", "topic", topic)
	return nil
}

// StopTopicWorker 停止指定 topic 的延时消息处理器
func (dw *DelayedWorker) StopTopicWorker(topic string) {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if cancel, exists := dw.workers[topic]; exists {
		cancel()
		delete(dw.workers, topic)
		delete(dw.topics, topic)
		dw.client.logger.Info("delayed worker stopped", "topic", topic)
	}
}

// StopAllWorkers 停止所有延时消息处理器
func (dw *DelayedWorker) StopAllWorkers() {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	for topic, cancel := range dw.workers {
		cancel()
		delete(dw.workers, topic)
		delete(dw.topics, topic)
		dw.client.logger.Info("delayed worker stopped", "topic", topic)
	}
}

// processDelayedMessages 处理延时消息
func (dw *DelayedWorker) processDelayedMessages(ctx context.Context, topic string) {
	ticker := time.NewTicker(dw.interval)
	defer ticker.Stop()

	delayKey := dw.client.prefix + "mq:delay:" + topic

	for {
		select {
		case <-ctx.Done():
			dw.client.logger.Info("delayed worker stopped", "topic", topic)
			return
		case <-ticker.C:
			now := float64(time.Now().Unix())

			// 获取到期的消息
			messages, err := dw.client.client.ZRangeByScoreWithScores(ctx, delayKey, &redis.ZRangeBy{
				Min: "0",
				Max: fmt.Sprintf("%f", now),
			}).Result()

			if err != nil {
				dw.client.logger.Error("failed to get delayed messages", err, "topic", topic)
				continue
			}

			if len(messages) == 0 {
				continue
			}

			dw.client.logger.Info("processing delayed messages", "topic", topic, "count", len(messages))

			for _, msg := range messages {
				// 发布消息到正常队列
				streamKey := dw.client.prefix + "mq:" + topic
				_, err := dw.client.client.XAdd(ctx, &redis.XAddArgs{
					Stream: streamKey,
					Values: map[string]interface{}{
						"message": msg.Member,
					},
				}).Result()

				if err != nil {
					dw.client.logger.Error("failed to publish delayed message", err, "topic", topic)
					continue
				}

				// 从延时队列中移除
				if err := dw.client.client.ZRem(ctx, delayKey, msg.Member).Err(); err != nil {
					dw.client.logger.Error("failed to remove delayed message", err, "topic", topic)
				}
			}
		}
	}
}

// GetActiveTopics 获取活跃的 topic 列表
func (dw *DelayedWorker) GetActiveTopics() []string {
	dw.mutex.RLock()
	defer dw.mutex.RUnlock()

	topics := make([]string, 0, len(dw.topics))
	for topic := range dw.topics {
		topics = append(topics, topic)
	}
	return topics
}

// IsTopicActive 检查 topic 是否活跃
func (dw *DelayedWorker) IsTopicActive(topic string) bool {
	dw.mutex.RLock()
	defer dw.mutex.RUnlock()

	return dw.topics[topic]
}
