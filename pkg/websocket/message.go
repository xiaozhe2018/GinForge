package websocket

import (
	"encoding/json"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	// 系统消息
	MessageTypeSystem       MessageType = "system"        // 系统消息
	MessageTypePing         MessageType = "ping"          // 心跳请求
	MessageTypePong         MessageType = "pong"          // 心跳响应
	MessageTypeWelcome      MessageType = "welcome"       // 欢迎消息
	MessageTypeError        MessageType = "error"         // 错误消息
	
	// 用户消息
	MessageTypeChat         MessageType = "chat"          // 聊天消息
	MessageTypeNotification MessageType = "notification"  // 通知消息
	MessageTypeBroadcast    MessageType = "broadcast"     // 广播消息
	
	// 房间消息
	MessageTypeJoinRoom     MessageType = "join_room"     // 加入房间
	MessageTypeLeaveRoom    MessageType = "leave_room"    // 离开房间
	MessageTypeRoomMessage  MessageType = "room_message"  // 房间消息
	
	// 用户状态
	MessageTypeUserOnline   MessageType = "user_online"   // 用户上线
	MessageTypeUserOffline  MessageType = "user_offline"  // 用户下线
	MessageTypeUserStatus   MessageType = "user_status"   // 用户状态变更
	
	// 数据更新
	MessageTypeDataUpdate   MessageType = "data_update"   // 数据更新通知
	MessageTypeRefresh      MessageType = "refresh"       // 刷新请求
)

// Message WebSocket 消息
type Message struct {
	Type      MessageType            `json:"type"`                 // 消息类型
	ID        string                 `json:"id,omitempty"`         // 消息ID
	From      string                 `json:"from,omitempty"`       // 发送者ID
	FromName  string                 `json:"from_name,omitempty"`  // 发送者名称
	To        string                 `json:"to,omitempty"`         // 接收者ID（私聊）
	Room      string                 `json:"room,omitempty"`       // 房间名称
	Content   interface{}            `json:"content"`              // 消息内容
	Data      map[string]interface{} `json:"data,omitempty"`       // 额外数据
	Timestamp int64                  `json:"timestamp"`            // 时间戳
	CreatedAt time.Time              `json:"created_at,omitempty"` // 创建时间
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, content interface{}) *Message {
	return &Message{
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now().Unix(),
		CreatedAt: time.Now(),
		Data:      make(map[string]interface{}),
	}
}

// ToJSON 转换为 JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从 JSON 解析
func (m *Message) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}

// SetData 设置额外数据
func (m *Message) SetData(key string, value interface{}) *Message {
	if m.Data == nil {
		m.Data = make(map[string]interface{})
	}
	m.Data[key] = value
	return m
}

// GetData 获取额外数据
func (m *Message) GetData(key string) (interface{}, bool) {
	if m.Data == nil {
		return nil, false
	}
	val, ok := m.Data[key]
	return val, ok
}

// NotificationMessage 通知消息内容
type NotificationMessage struct {
	Title    string                 `json:"title"`              // 标题
	Body     string                 `json:"body"`               // 内容
	Icon     string                 `json:"icon,omitempty"`     // 图标
	Link     string                 `json:"link,omitempty"`     // 跳转链接
	Category string                 `json:"category,omitempty"` // 分类
	Data     map[string]interface{} `json:"data,omitempty"`     // 额外数据
}

// NewNotificationMessage 创建通知消息
func NewNotificationMessage(title, body string) *NotificationMessage {
	return &NotificationMessage{
		Title: title,
		Body:  body,
		Data:  make(map[string]interface{}),
	}
}

// ChatMessage 聊天消息内容
type ChatMessage struct {
	Text      string                 `json:"text"`               // 文本内容
	MediaType string                 `json:"media_type,omitempty"` // 媒体类型: text, image, video, file
	MediaURL  string                 `json:"media_url,omitempty"`  // 媒体URL
	Reply     *ChatMessage           `json:"reply,omitempty"`      // 回复的消息
	Data      map[string]interface{} `json:"data,omitempty"`       // 额外数据
}

// SystemMessage 系统消息内容
type SystemMessage struct {
	Code    int                    `json:"code"`             // 消息代码
	Message string                 `json:"message"`          // 消息内容
	Level   string                 `json:"level,omitempty"`  // 级别: info, warning, error
	Data    map[string]interface{} `json:"data,omitempty"`   // 额外数据
}

// UserStatusMessage 用户状态消息
type UserStatusMessage struct {
	UserID   string `json:"user_id"`   // 用户ID
	UserName string `json:"user_name"` // 用户名
	Status   string `json:"status"`    // 状态: online, offline, busy, away
	Avatar   string `json:"avatar,omitempty"` // 头像
}

// DataUpdateMessage 数据更新消息
type DataUpdateMessage struct {
	Entity string                 `json:"entity"`           // 实体类型: user, role, menu, etc.
	Action string                 `json:"action"`           // 操作: create, update, delete
	ID     string                 `json:"id,omitempty"`     // 实体ID
	Data   map[string]interface{} `json:"data,omitempty"`   // 更新数据
}

