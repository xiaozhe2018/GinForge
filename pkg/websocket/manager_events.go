package websocket

import (
	"context"
	"time"
)

// 事件类型常量
const (
	// 连接事件
	EventClientConnected    = "client.connected"
	EventClientDisconnected = "client.disconnected"
	
	// 消息事件
	EventMessageReceived = "message.received"
	EventMessageSent     = "message.sent"
	
	// 分组事件
	EventClientJoinedGroup  = "group.joined"
	EventClientLeftGroup    = "group.left"
	
	// 会话事件
	EventSessionCreated = "session.created"
	EventSessionUpdated = "session.updated"
	EventSessionDeleted = "session.deleted"
)

// ClientEventData 客户端事件数据
type ClientEventData struct {
	ClientID string
	UserID   string
	UserName string
	Time     time.Time
	Metadata map[string]interface{}
}

// MessageEventData 消息事件数据
type MessageEventData struct {
	Message  *Message
	ClientID string
	UserID   string
	Time     time.Time
}

// GroupEventData 分组事件数据
type GroupEventData struct {
	GroupID  string
	ClientID string
	UserID   string
	UserName string
	Time     time.Time
}

// SessionEventData 会话事件数据
type SessionEventData struct {
	SessionID string
	ClientID  string
	UserID    string
	Data      map[string]interface{}
	Time      time.Time
}

// InitEventBus 初始化事件总线
func (m *Manager) InitEventBus() {
	// 在实际实现中初始化事件总线
}

// On 注册事件处理器
func (m *Manager) On(eventType string, handler interface{}) {
	// 在实际实现中注册事件处理器
}

// Off 移除事件处理器
func (m *Manager) Off(eventType string, handler interface{}) {
	// 在实际实现中移除事件处理器
}

// UseMiddleware 添加事件中间件
func (m *Manager) UseMiddleware(middleware ...interface{}) {
	// 在实际实现中添加中间件
}

// emitEvent 触发事件
func (m *Manager) emitEvent(eventType string, data interface{}) {
	// 在实际实现中触发事件
}

// emitEventWithContext 触发带上下文的事件
func (m *Manager) emitEventWithContext(ctx context.Context, eventType string, data interface{}) {
	// 在实际实现中触发带上下文的事件
}

// 修改 registerClient 方法，添加事件触发
func (m *Manager) registerClientWithEvents(client *Client) {
	m.mu.Lock()
	
	// 添加到客户端映射
	m.clients[client.ID] = client
	
	// 添加到用户客户端映射
	if _, ok := m.userClients[client.UserID]; !ok {
		m.userClients[client.UserID] = make(map[string]*Client)
	}
	m.userClients[client.UserID][client.ID] = client
	
	m.mu.Unlock()
	
	m.logger.Info("client registered",
		"client_id", client.ID,
		"user_id", client.UserID,
		"user_name", client.UserName)
	
	// 发送欢迎消息
	welcome := NewMessage(MessageTypeWelcome, SystemMessage{
		Code:    0,
		Message: "欢迎连接到 GinForge WebSocket 服务",
		Level:   "info",
	})
	welcome.SetData("client_id", client.ID)
	welcome.SetData("server_time", time.Now())
	client.SendMessage(welcome)
	
	// 广播用户上线消息
	m.BroadcastUserStatus(client.UserID, client.UserName, "online")
	
	// 触发连接事件
	m.emitEvent(EventClientConnected, ClientEventData{
		ClientID: client.ID,
		UserID:   client.UserID,
		UserName: client.UserName,
		Time:     time.Now(),
		Metadata: client.MetaData,
	})
}

// 修改 unregisterClient 方法，添加事件触发
func (m *Manager) unregisterClientWithEvents(client *Client) {
	m.mu.Lock()
	
	// 从客户端映射中删除
	if _, ok := m.clients[client.ID]; ok {
		delete(m.clients, client.ID)
		close(client.Send)
	}
	
	// 从用户客户端映射中删除
	if userClientsMap, ok := m.userClients[client.UserID]; ok {
		delete(userClientsMap, client.ID)
		if len(userClientsMap) == 0 {
			delete(m.userClients, client.UserID)
			// 用户完全下线，广播下线消息
			m.broadcastUserOffline(client.UserID, client.UserName)
		}
	}
	
	// 从所有房间中删除
	for room := range client.Rooms {
		if roomClients, ok := m.rooms[room]; ok {
			delete(roomClients, client.ID)
			if len(roomClients) == 0 {
				delete(m.rooms, room)
			}
			
			// 触发离开房间事件
			m.emitEvent(EventClientLeftGroup, GroupEventData{
				GroupID:  room,
				ClientID: client.ID,
				UserID:   client.UserID,
				UserName: client.UserName,
				Time:     time.Now(),
			})
		}
	}
	
	m.mu.Unlock()
	
	m.logger.Info("client unregistered",
		"client_id", client.ID,
		"user_id", client.UserID)
	
	// 触发断开连接事件
	m.emitEvent(EventClientDisconnected, ClientEventData{
		ClientID: client.ID,
		UserID:   client.UserID,
		UserName: client.UserName,
		Time:     time.Now(),
	})
}

// 修改 HandleMessage 方法，添加事件触发
func (m *Manager) HandleMessageWithEvents(client *Client, message *Message) {
	// 触发消息接收事件
	m.emitEvent(EventMessageReceived, MessageEventData{
		Message:  message,
		ClientID: client.ID,
		UserID:   client.UserID,
		Time:     time.Now(),
	})
	
	// 查找处理器
	handler, ok := m.handlers[message.Type]
	if !ok {
		// 默认处理器
		m.defaultHandler(client, message)
		return
	}
	
	// 执行处理器
	if err := handler(client, message); err != nil {
		m.logger.Error("message handler error",
			err,
			"type", message.Type,
			"client_id", client.ID,
			"user_id", client.UserID)
		
		// 发送错误消息给客户端
		errMsg := NewMessage(MessageTypeError, SystemMessage{
			Code:    500,
			Message: "消息处理失败",
			Level:   "error",
		})
		client.SendMessage(errMsg)
	}
}

// 修改 AddClientToRoom 方法，添加事件触发
func (m *Manager) AddClientToRoomWithEvents(room string, client *Client) {
	m.mu.Lock()
	
	if _, ok := m.rooms[room]; !ok {
		m.rooms[room] = make(map[string]*Client)
	}
	m.rooms[room][client.ID] = client
	
	m.mu.Unlock()
	
	// 触发加入房间事件
	m.emitEvent(EventClientJoinedGroup, GroupEventData{
		GroupID:  room,
		ClientID: client.ID,
		UserID:   client.UserID,
		UserName: client.UserName,
		Time:     time.Now(),
	})
}

// 修改 RemoveClientFromRoom 方法，添加事件触发
func (m *Manager) RemoveClientFromRoomWithEvents(room string, client *Client) {
	m.mu.Lock()
	
	if roomClients, ok := m.rooms[room]; ok {
		delete(roomClients, client.ID)
		if len(roomClients) == 0 {
			delete(m.rooms, room)
		}
	}
	
	m.mu.Unlock()
	
	// 触发离开房间事件
	m.emitEvent(EventClientLeftGroup, GroupEventData{
		GroupID:  room,
		ClientID: client.ID,
		UserID:   client.UserID,
		UserName: client.UserName,
		Time:     time.Now(),
	})
}
