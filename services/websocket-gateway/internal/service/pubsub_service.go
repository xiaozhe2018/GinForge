package service

import (
	"context"
	"encoding/json"

	"goweb/pkg/logger"
	"goweb/pkg/redis"
	"goweb/pkg/websocket"
)

// PubSubService Redis 发布订阅服务
// 用于接收其他服务发送的 WebSocket 消息
type PubSubService struct {
	wsManager   *websocket.Manager
	redisClient *redis.Client
	logger      logger.Logger
	cancel      context.CancelFunc
}

// NewPubSubService 创建发布订阅服务
func NewPubSubService(wsManager *websocket.Manager, redisClient *redis.Client, log logger.Logger) *PubSubService {
	return &PubSubService{
		wsManager:   wsManager,
		redisClient: redisClient,
		logger:      log,
	}
}

// Start 启动订阅
func (s *PubSubService) Start(ctx context.Context) error {
	if s.redisClient == nil {
		s.logger.Warn("Redis is not enabled, pubsub will not work")
		return nil
	}

	// 订阅频道（支持通配符）
	channels := []string{
		"websocket:broadcast",      // 广播消息
		"websocket:user:*",         // 用户消息
		"websocket:room:*",         // 房间消息
		"websocket:notification:*", // 通知消息
	}

	pubsubCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.logger.Info("starting redis pubsub", "channels", channels)

	// 订阅
	pubsub := s.redisClient.GetClient().PSubscribe(pubsubCtx, channels...)

	// 在协程中接收消息
	go func() {
		defer pubsub.Close()

		for {
			select {
			case <-pubsubCtx.Done():
				s.logger.Info("redis pubsub stopped")
				return

			case msg := <-pubsub.Channel():
				s.handleMessage(msg.Channel, msg.Payload)
			}
		}
	}()

	return nil
}

// Stop 停止订阅
func (s *PubSubService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// handleMessage 处理接收到的消息
func (s *PubSubService) handleMessage(channel, payload string) {
	s.logger.Info("received redis message", "channel", channel, "payload_size", len(payload), "payload", payload)

	var msg RedisMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.logger.Error("failed to parse redis message", err, "channel", channel, "payload", payload)
		return
	}
	
	s.logger.Info("parsed redis message", "message_type", msg.Message.Type)

	// 根据频道类型路由消息
	switch {
	case channel == "websocket:broadcast":
		// 广播消息到所有客户端
		s.logger.Info("broadcasting message to all clients", "message_type", msg.Message.Type)
		s.wsManager.Broadcast(msg.Message)
		s.logger.Info("broadcast completed")

	case len(channel) > 18 && channel[:18] == "websocket:user:":
		// 发送消息到指定用户
		userID := channel[18:]
		if err := s.wsManager.SendToUser(userID, msg.Message); err != nil {
			s.logger.Warn("failed to send to user", "user_id", userID, "error", err)
		}

	case len(channel) > 18 && channel[:18] == "websocket:room:":
		// 广播消息到房间
		room := channel[18:]
		if err := s.wsManager.BroadcastToRoom(room, msg.Message); err != nil {
			s.logger.Warn("failed to broadcast to room", "room", room, "error", err)
		}

	case len(channel) > 25 && channel[:25] == "websocket:notification:":
		// 发送通知消息
		userID := channel[25:]
		if notification, ok := msg.Message.Content.(*websocket.NotificationMessage); ok {
			if err := s.wsManager.SendNotification(userID, notification); err != nil {
				s.logger.Warn("failed to send notification", "user_id", userID, "error", err)
			}
		} else {
			// 尝试重新解析通知内容
			if content, ok := msg.Message.Content.(map[string]interface{}); ok {
				notification := &websocket.NotificationMessage{
					Title: getStringValue(content, "title"),
					Body:  getStringValue(content, "body"),
					Icon:  getStringValue(content, "icon"),
					Link:  getStringValue(content, "link"),
				}
				if err := s.wsManager.SendNotification(userID, notification); err != nil {
					s.logger.Warn("failed to send notification", "user_id", userID, "error", err)
				}
			}
		}
	}
}

// RedisMessage Redis 传输的消息格式
type RedisMessage struct {
	Message *websocket.Message `json:"message"`
}

// getStringValue 安全获取 map 中的字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

