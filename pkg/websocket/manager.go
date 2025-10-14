package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"goweb/pkg/logger"
)

var (
	// ErrClientNotFound 客户端不存在
	ErrClientNotFound = errors.New("client not found")
	
	// ErrBufferFull 发送缓冲区已满
	ErrBufferFull = errors.New("send buffer is full")
	
	// ErrRoomNotFound 房间不存在
	ErrRoomNotFound = errors.New("room not found")
)

// MessageHandler 消息处理器
type MessageHandler func(client *Client, message *Message) error

// Manager WebSocket 连接管理器
type Manager struct {
	clients     map[string]*Client           // 客户端映射 (clientID -> client)
	userClients map[string]map[string]*Client // 用户客户端映射 (userID -> map[clientID]client)
	rooms       map[string]map[string]*Client // 房间映射 (room -> map[clientID]client)
	handlers    map[MessageType]MessageHandler // 消息处理器
	register    chan *Client                  // 注册通道
	unregister  chan *Client                  // 注销通道
	broadcast   chan *Message                 // 广播通道
	mu          sync.RWMutex                  // 读写锁
	logger      logger.Logger                 // 日志器
	ctx         context.Context               // 上下文
	cancel      context.CancelFunc            // 取消函数
	eventBus    interface{}                  // 事件总线
}

// NewManager 创建 WebSocket 管理器
func NewManager(log logger.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Manager{
		clients:     make(map[string]*Client),
		userClients: make(map[string]map[string]*Client),
		rooms:       make(map[string]map[string]*Client),
		handlers:    make(map[MessageType]MessageHandler),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *Message, 256),
		logger:      log,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Run 运行管理器
func (m *Manager) Run() {
	m.logger.Info("WebSocket manager started")
	
	// 初始化事件总线
	m.InitEventBus()
	
	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("WebSocket manager stopped")
			return
			
		case client := <-m.register:
			m.registerClientWithEvents(client)
			
		case client := <-m.unregister:
			m.unregisterClientWithEvents(client)
			
		case message := <-m.broadcast:
			m.broadcastMessage(message)
		}
	}
}

// Stop 停止管理器
func (m *Manager) Stop() {
	m.cancel()
	
	// 关闭所有客户端
	m.mu.Lock()
	for _, client := range m.clients {
		close(client.Send)
		client.Conn.Close()
	}
	m.mu.Unlock()
}

// Register 注册客户端
func (m *Manager) Register(client *Client) {
	m.register <- client
}

// Unregister 注销客户端
func (m *Manager) Unregister(client *Client) {
	m.unregister <- client
}

// registerClient 注册客户端（内部）
func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 添加到客户端映射
	m.clients[client.ID] = client
	
	// 添加到用户客户端映射
	if _, ok := m.userClients[client.UserID]; !ok {
		m.userClients[client.UserID] = make(map[string]*Client)
	}
	m.userClients[client.UserID][client.ID] = client
	
	m.logger.Info("client registered",
		"client_id", client.ID,
		"user_id", client.UserID,
		"user_name", client.UserName,
		"total_clients", len(m.clients))
	
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
}

// unregisterClient 注销客户端（内部）
func (m *Manager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
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
		}
	}
	
	m.logger.Info("client unregistered",
		"client_id", client.ID,
		"user_id", client.UserID,
		"total_clients", len(m.clients))
}

// broadcastMessage 广播消息（内部）
func (m *Manager) broadcastMessage(message *Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, client := range m.clients {
		select {
		case client.Send <- mustMarshal(message):
		default:
			m.logger.Warn("failed to send broadcast message to client",
				"client_id", client.ID,
				"user_id", client.UserID)
		}
	}
}

// broadcastUserOffline 广播用户下线（内部，需要持有锁）
func (m *Manager) broadcastUserOffline(userID, userName string) {
	msg := NewMessage(MessageTypeUserOffline, UserStatusMessage{
		UserID:   userID,
		UserName: userName,
		Status:   "offline",
	})
	
	data := mustMarshal(msg)
	for _, client := range m.clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

// Broadcast 广播消息给所有客户端
func (m *Manager) Broadcast(message *Message) {
	m.broadcast <- message
}

// BroadcastToRoom 广播消息到房间
func (m *Manager) BroadcastToRoom(room string, message *Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	roomClients, ok := m.rooms[room]
	if !ok {
		return ErrRoomNotFound
	}
	
	message.Room = room
	data := mustMarshal(message)
	
	for _, client := range roomClients {
		select {
		case client.Send <- data:
		default:
			m.logger.Warn("failed to send room message", "room", room, "client_id", client.ID)
		}
	}
	
	return nil
}

// SendToUser 发送消息给特定用户（所有连接）
func (m *Manager) SendToUser(userID string, message *Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	userClientsMap, ok := m.userClients[userID]
	if !ok || len(userClientsMap) == 0 {
		return ErrClientNotFound
	}
	
	message.To = userID
	data := mustMarshal(message)
	
	for _, client := range userClientsMap {
		select {
		case client.Send <- data:
		default:
			m.logger.Warn("failed to send message to user", "user_id", userID, "client_id", client.ID)
		}
	}
	
	return nil
}

// SendToClient 发送消息给特定客户端
func (m *Manager) SendToClient(clientID string, message *Message) error {
	m.mu.RLock()
	client, ok := m.clients[clientID]
	m.mu.RUnlock()
	
	if !ok {
		return ErrClientNotFound
	}
	
	return client.SendMessage(message)
}

// AddClientToRoom 添加客户端到房间
func (m *Manager) AddClientToRoom(room string, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.rooms[room]; !ok {
		m.rooms[room] = make(map[string]*Client)
	}
	m.rooms[room][client.ID] = client
}

// RemoveClientFromRoom 从房间移除客户端
func (m *Manager) RemoveClientFromRoom(room string, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if roomClients, ok := m.rooms[room]; ok {
		delete(roomClients, client.ID)
		if len(roomClients) == 0 {
			delete(m.rooms, room)
		}
	}
}

// GetClient 获取客户端
func (m *Manager) GetClient(clientID string) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[clientID]
	return client, ok
}

// GetUserClients 获取用户的所有客户端
func (m *Manager) GetUserClients(userID string) []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	userClientsMap, ok := m.userClients[userID]
	if !ok {
		return nil
	}
	
	clients := make([]*Client, 0, len(userClientsMap))
	for _, client := range userClientsMap {
		clients = append(clients, client)
	}
	return clients
}

// IsUserOnline 检查用户是否在线
func (m *Manager) IsUserOnline(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	userClientsMap, ok := m.userClients[userID]
	return ok && len(userClientsMap) > 0
}

// GetOnlineUsers 获取所有在线用户
func (m *Manager) GetOnlineUsers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	users := make([]string, 0, len(m.userClients))
	for userID := range m.userClients {
		users = append(users, userID)
	}
	return users
}

// GetStats 获取统计信息
func (m *Manager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return map[string]interface{}{
		"total_clients": len(m.clients),
		"online_users":  len(m.userClients),
		"total_rooms":   len(m.rooms),
		"uptime":        time.Since(time.Now()), // 简化版，实际应记录启动时间
	}
}

// RegisterHandler 注册消息处理器
func (m *Manager) RegisterHandler(msgType MessageType, handler MessageHandler) {
	m.handlers[msgType] = handler
}

// HandleMessage 处理消息
func (m *Manager) HandleMessage(client *Client, message *Message) {
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

// defaultHandler 默认消息处理器
func (m *Manager) defaultHandler(client *Client, message *Message) {
	switch message.Type {
	case MessageTypePing:
		// 心跳响应
		pong := NewMessage(MessageTypePong, nil)
		client.SendMessage(pong)
		
	case MessageTypeJoinRoom:
		// 加入房间
		if room, ok := message.Content.(string); ok {
			client.JoinRoom(room)
		}
		
	case MessageTypeLeaveRoom:
		// 离开房间
		if room, ok := message.Content.(string); ok {
			client.LeaveRoom(room)
		}
		
	case MessageTypeRoomMessage:
		// 房间消息
		if message.Room != "" {
			m.BroadcastToRoom(message.Room, message)
		}
		
	case MessageTypeChat:
		// 私聊消息
		if message.To != "" {
			m.SendToUser(message.To, message)
		}
		
	default:
		m.logger.Warn("unhandled message type",
			"type", message.Type,
			"client_id", client.ID)
	}
}

// BroadcastUserStatus 广播用户状态
func (m *Manager) BroadcastUserStatus(userID, userName, status string) {
	msg := NewMessage(MessageTypeUserStatus, UserStatusMessage{
		UserID:   userID,
		UserName: userName,
		Status:   status,
	})
	m.Broadcast(msg)
}

// BroadcastDataUpdate 广播数据更新
func (m *Manager) BroadcastDataUpdate(entity, action, id string, data map[string]interface{}) {
	msg := NewMessage(MessageTypeDataUpdate, DataUpdateMessage{
		Entity: entity,
		Action: action,
		ID:     id,
		Data:   data,
	})
	m.Broadcast(msg)
}

// SendNotification 发送通知给用户
func (m *Manager) SendNotification(userID string, notification *NotificationMessage) error {
	msg := NewMessage(MessageTypeNotification, notification)
	msg.SetData("notification_id", uuid.New().String())
	msg.SetData("read", false)
	return m.SendToUser(userID, msg)
}

// mustMarshal JSON 序列化（panic on error）
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// GetRoomClients 获取房间内的所有客户端
func (m *Manager) GetRoomClients(room string) []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	roomClients, ok := m.rooms[room]
	if !ok {
		return nil
	}
	
	clients := make([]*Client, 0, len(roomClients))
	for _, client := range roomClients {
		clients = append(clients, client)
	}
	return clients
}

// GetRoomSize 获取房间人数
func (m *Manager) GetRoomSize(room string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if roomClients, ok := m.rooms[room]; ok {
		return len(roomClients)
	}
	return 0
}

