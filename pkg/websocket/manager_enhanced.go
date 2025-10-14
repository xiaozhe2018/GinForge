package websocket

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"goweb/pkg/websocket/session"
)

// 增强 Manager 结构体，添加会话管理器
type ManagerEnhanced struct {
	*Manager
	sessionManager *session.SessionManager
	redisClient    *redis.Client
}

// NewEnhancedManager 创建增强型管理器
func NewEnhancedManager(manager *Manager, redisClient *redis.Client) *ManagerEnhanced {
	var sessionMgr *session.SessionManager
	if redisClient != nil {
		// 使用 Redis 会话管理
		sessionMgr = session.NewRedisSessionManager(redisClient, "ws", 24*time.Hour)
	} else {
		// 使用内存会话管理
		sessionMgr = session.NewMemorySessionManager()
	}
	
	return &ManagerEnhanced{
		Manager:        manager,
		sessionManager: sessionMgr,
		redisClient:    redisClient,
	}
}

// 重写 registerClient 方法
func (m *ManagerEnhanced) RegisterClient(client *Client) {
	// 调用原始注册方法
	m.Manager.Register(client)
	
	// 创建或恢复会话
	sess, err := m.sessionManager.CreateSession(client.ID, client.UserID)
	if err != nil {
		m.logger.Error("failed to create session", err)
		return
	}
	
	// 恢复之前的房间
	if rooms, ok := sess.Get("rooms"); ok {
		if roomList, ok := rooms.([]string); ok {
			for _, room := range roomList {
				client.JoinRoom(room)
			}
		}
	}
	
	// 恢复元数据
	if metadata, ok := sess.Get("metadata"); ok {
		if metaMap, ok := metadata.(map[string]interface{}); ok {
			for k, v := range metaMap {
				client.SetMetaData(k, v)
			}
		}
	}
}

// 重写 unregisterClient 方法
func (m *ManagerEnhanced) UnregisterClient(client *Client) {
	// 保存会话数据
	rooms := client.GetRooms()
	if len(rooms) > 0 {
		m.sessionManager.SetSessionData(client.ID, "rooms", rooms)
	}
	
	// 保存元数据
	m.mu.RLock()
	if len(client.MetaData) > 0 {
		m.sessionManager.SetSessionData(client.ID, "metadata", client.MetaData)
	}
	m.mu.RUnlock()
	
	// 调用原始注销方法
	m.Manager.Unregister(client)
}

// 会话管理方法

// GetClientSession 获取客户端会话
func (m *ManagerEnhanced) GetClientSession(clientID string) (*session.Session, error) {
	return m.sessionManager.GetSession(clientID)
}

// SetClientSessionData 设置客户端会话数据
func (m *ManagerEnhanced) SetClientSessionData(clientID string, key string, value interface{}) error {
	return m.sessionManager.SetSessionData(clientID, key, value)
}

// GetClientSessionData 获取客户端会话数据
func (m *ManagerEnhanced) GetClientSessionData(clientID string, key string) (interface{}, error) {
	return m.sessionManager.GetSessionData(clientID, key)
}

// UpdateClientSessionData 更新客户端会话数据
func (m *ManagerEnhanced) UpdateClientSessionData(clientID string, data map[string]interface{}) error {
	return m.sessionManager.UpdateSessionData(clientID, data)
}

// 统一消息发送方法

// SendMessage 发送消息
func (m *ManagerEnhanced) SendMessage(target Target, msgType MessageType, content interface{}, options ...MessageOption) error {
	if err := target.Validate(); err != nil {
		return err
	}
	
	msg := NewMessage(msgType, content)
	
	// 应用选项
	for _, option := range options {
		option(msg)
	}
	
	// 根据目标类型发送
	switch target.Type {
	case TargetTypeUser:
		return m.SendToUser(target.ID, msg)
	case TargetTypeGroup:
		return m.BroadcastToRoom(target.ID, msg)
	case TargetTypeClient:
		return m.SendToClient(target.ID, msg)
	case TargetTypeAll:
		m.Broadcast(msg)
		return nil
	default:
		return errors.New("unknown target type")
	}
}

// SendNotification 发送通知消息
func (m *ManagerEnhanced) SendNotification(target Target, title, body string, options ...NotificationOption) error {
	notification := NewNotificationMessage(title, body)
	
	// 应用选项
	for _, option := range options {
		option(notification)
	}
	
	msg := NewMessage(MessageTypeNotification, notification)
	msg.SetData("notification_id", uuid.New().String())
	
	// 根据目标类型发送
	switch target.Type {
	case TargetTypeUser:
		return m.SendToUser(target.ID, msg)
	case TargetTypeGroup:
		return m.BroadcastToRoom(target.ID, msg)
	case TargetTypeClient:
		return m.SendToClient(target.ID, msg)
	case TargetTypeAll:
		m.Broadcast(msg)
		return nil
	default:
		return errors.New("unknown target type")
	}
}

// SendSystemMessage 发送系统消息
func (m *ManagerEnhanced) SendSystemMessage(target Target, message string, options ...SystemMessageOption) error {
	sysMsg := &SystemMessage{
		Code:    0,
		Message: message,
		Level:   "info",
	}
	
	// 应用选项
	for _, option := range options {
		option(sysMsg)
	}
	
	return m.SendMessage(target, MessageTypeSystem, sysMsg)
}

// SendDataUpdate 发送数据更新消息
func (m *ManagerEnhanced) SendDataUpdate(target Target, entity, action string, options ...DataUpdateOption) error {
	updateMsg := &DataUpdateMessage{
		Entity: entity,
		Action: action,
	}
	
	// 应用选项
	for _, option := range options {
		option(updateMsg)
	}
	
	return m.SendMessage(target, MessageTypeDataUpdate, updateMsg)
}
