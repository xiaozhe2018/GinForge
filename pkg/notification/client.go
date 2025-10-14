package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"goweb/pkg/redis"
	"goweb/pkg/websocket"
)

// Client 通知客户端
// 其他服务使用此客户端发送 WebSocket 通知（通过 Redis PubSub）
type Client struct {
	redisClient *redis.Client
}

// NewClient 创建通知客户端
func NewClient(redisClient *redis.Client) *Client {
	return &Client{
		redisClient: redisClient,
	}
}

// SendNotification 发送通知给指定用户
func (c *Client) SendNotification(ctx context.Context, userID string, notification *websocket.NotificationMessage) error {
	if c.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	message := websocket.NewMessage(websocket.MessageTypeNotification, notification)
	
	redisMsg := RedisMessage{
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	channel := fmt.Sprintf("websocket:notification:%s", userID)
	return c.redisClient.GetClient().Publish(ctx, channel, data).Err()
}

// BroadcastNotification 广播通知给所有用户
func (c *Client) BroadcastNotification(ctx context.Context, notification *websocket.NotificationMessage) error {
	if c.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	message := websocket.NewMessage(websocket.MessageTypeNotification, notification)
	
	redisMsg := RedisMessage{
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	return c.redisClient.GetClient().Publish(ctx, "websocket:broadcast", data).Err()
}

// SendMessage 发送自定义消息给用户
func (c *Client) SendMessage(ctx context.Context, userID string, message *websocket.Message) error {
	if c.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	redisMsg := RedisMessage{
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("websocket:user:%s", userID)
	return c.redisClient.GetClient().Publish(ctx, channel, data).Err()
}

// BroadcastMessage 广播自定义消息
func (c *Client) BroadcastMessage(ctx context.Context, message *websocket.Message) error {
	if c.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	redisMsg := RedisMessage{
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.redisClient.GetClient().Publish(ctx, "websocket:broadcast", data).Err()
}

// SendToRoom 发送消息到房间
func (c *Client) SendToRoom(ctx context.Context, room string, message *websocket.Message) error {
	if c.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	redisMsg := RedisMessage{
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("websocket:room:%s", room)
	return c.redisClient.GetClient().Publish(ctx, channel, data).Err()
}

// BroadcastDataUpdate 广播数据更新消息
func (c *Client) BroadcastDataUpdate(ctx context.Context, entity, action, id string, data map[string]interface{}) error {
	updateMsg := &websocket.DataUpdateMessage{
		Entity: entity,
		Action: action,
		ID:     id,
		Data:   data,
	}
	
	message := websocket.NewMessage(websocket.MessageTypeDataUpdate, updateMsg)
	return c.BroadcastMessage(ctx, message)
}

// RedisMessage Redis 传输的消息格式
type RedisMessage struct {
	Message *websocket.Message `json:"message"`
}

