package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Queue 消息队列接口
type Queue interface {
	// 发布消息
	Publish(ctx context.Context, topic string, data map[string]interface{}) error
	// 延迟发布消息
	PublishWithDelay(ctx context.Context, topic string, data map[string]interface{}, delay time.Duration) error
	// 订阅消息
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	// 获取队列长度
	GetQueueLength(ctx context.Context, topic string) (int64, error)
	// 清空队列
	PurgeQueue(ctx context.Context, topic string) error
}

// Message 消息结构
type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Retry     int                    `json:"retry"`
	MaxRetry  int                    `json:"max_retry"`
}

// MessageHandler 消息处理器
type MessageHandler func(ctx context.Context, msg *Message) error

// RedisQueue Redis 消息队列实现
type RedisQueue struct {
	client        *Client
	prefix        string
	delayedWorker *DelayedWorker
}

// NewQueue 创建 Redis 消息队列
func NewQueue(client *Client, prefix string) *RedisQueue {
	return &RedisQueue{
		client:        client,
		prefix:        prefix,
		delayedWorker: NewDelayedWorker(client, time.Second),
	}
}

// Publish 发布消息
func (q *RedisQueue) Publish(ctx context.Context, topic string, data map[string]interface{}) error {
	if !q.client.IsEnabled() {
		return nil
	}

	message := &Message{
		ID:        q.generateMessageID(),
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now(),
		Retry:     0,
		MaxRetry:  3,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 发布到 Redis Stream
	streamKey := q.prefix + "mq:" + topic
	_, err = q.client.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"message": string(messageBytes),
		},
	}).Result()

	if err != nil {
		q.client.logger.Error("failed to publish message", err, "topic", topic, "message_id", message.ID)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	q.client.logger.Info("message published", "topic", topic, "message_id", message.ID)
	return nil
}

// Subscribe 订阅消息
func (q *RedisQueue) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	if !q.client.IsEnabled() {
		return nil
	}

	streamKey := q.prefix + "mq:" + topic
	groupName := fmt.Sprintf("consumer-group-%s", topic)
	consumerName := fmt.Sprintf("consumer-%d", time.Now().UnixNano())

	// 创建消费者组
	err := q.client.client.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	q.client.logger.Info("started consuming messages", "topic", topic, "group", groupName, "consumer", consumerName)

	for {
		select {
		case <-ctx.Done():
			q.client.logger.Info("stopped consuming messages", "topic", topic)
			return ctx.Err()
		default:
			// 读取消息
			streams, err := q.client.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamKey, ">"},
				Count:    1,
				Block:    time.Second * 5,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // 没有消息，继续等待
				}
				q.client.logger.Error("failed to read message", err, "topic", topic)
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					if err := q.processMessage(ctx, message, handler); err != nil {
						q.client.logger.Error("failed to process message", err, "topic", topic, "message_id", message.ID)
					}
				}
			}
		}
	}
}

// processMessage 处理消息
func (q *RedisQueue) processMessage(ctx context.Context, redisMsg redis.XMessage, handler MessageHandler) error {
	messageData, exists := redisMsg.Values["message"]
	if !exists {
		return fmt.Errorf("message data not found")
	}

	var message Message
	if err := json.Unmarshal([]byte(messageData.(string)), &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// 处理消息
	if err := handler(ctx, &message); err != nil {
		// 处理失败，增加重试次数
		message.Retry++
		if message.Retry < message.MaxRetry {
			// 重新发布消息
			messageBytes, _ := json.Marshal(message)
			q.client.client.XAdd(ctx, &redis.XAddArgs{
				Stream: q.prefix + "mq:" + message.Topic,
				Values: map[string]interface{}{
					"message": string(messageBytes),
				},
			})
			q.client.logger.Warn("message processing failed, retrying", "topic", message.Topic, "message_id", message.ID, "retry", message.Retry)
		} else {
			// 达到最大重试次数，记录到死信队列
			q.sendToDeadLetterQueue(ctx, &message, err)
		}
		return err
	}

	// 处理成功，确认消息
	streamKey := q.prefix + "mq:" + message.Topic
	groupName := fmt.Sprintf("consumer-group-%s", message.Topic)
	q.client.client.XAck(ctx, streamKey, groupName, redisMsg.ID)

	q.client.logger.Info("message processed successfully", "topic", message.Topic, "message_id", message.ID)
	return nil
}

// sendToDeadLetterQueue 发送到死信队列
func (q *RedisQueue) sendToDeadLetterQueue(ctx context.Context, message *Message, err error) {
	deadLetterData := map[string]interface{}{
		"original_message": message,
		"error":            err.Error(),
		"failed_at":        time.Now(),
	}

	deadLetterBytes, _ := json.Marshal(deadLetterData)
	q.client.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.prefix + "mq:dead-letter:" + message.Topic,
		Values: map[string]interface{}{
			"message": string(deadLetterBytes),
		},
	})

	q.client.logger.Error("message sent to dead letter queue", err, "topic", message.Topic, "message_id", message.ID)
}

// generateMessageID 生成消息ID
func (q *RedisQueue) generateMessageID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

// GetQueueLength 获取队列长度
func (q *RedisQueue) GetQueueLength(ctx context.Context, topic string) (int64, error) {
	if !q.client.IsEnabled() {
		return 0, nil
	}

	streamKey := q.prefix + "mq:" + topic
	return q.client.client.XLen(ctx, streamKey).Result()
}

// PurgeQueue 清空队列
func (q *RedisQueue) PurgeQueue(ctx context.Context, topic string) error {
	if !q.client.IsEnabled() {
		return nil
	}

	streamKey := q.prefix + "mq:" + topic
	return q.client.client.Del(ctx, streamKey).Err()
}

// PublishWithDelay 延迟发布消息
func (q *RedisQueue) PublishWithDelay(ctx context.Context, topic string, data map[string]interface{}, delay time.Duration) error {
	if !q.client.IsEnabled() {
		return nil
	}

	message := &Message{
		ID:        q.generateMessageID(),
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now().Add(delay),
		Retry:     0,
		MaxRetry:  3,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 使用 Redis 的延迟队列功能
	delayKey := q.prefix + "mq:delay:" + topic
	score := float64(time.Now().Add(delay).Unix())

	err = q.client.client.ZAdd(ctx, delayKey, redis.Z{
		Score:  score,
		Member: string(messageBytes),
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to schedule delayed message: %w", err)
	}

	// 启动延迟消息处理器（如果未启动）
	if err := q.delayedWorker.StartTopicWorker(ctx, topic); err != nil {
		q.client.logger.Error("failed to start delayed worker", err, "topic", topic)
	}

	q.client.logger.Info("delayed message scheduled", "topic", topic, "message_id", message.ID, "delay", delay)
	return nil
}
