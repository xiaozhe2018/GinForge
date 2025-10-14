package websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"goweb/pkg/logger"
)

const (
	// 写超时
	writeWait = 10 * time.Second

	// 读取ping的等待时间
	pongWait = 60 * time.Second

	// ping周期，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512 * 1024 // 512KB
)

// Client WebSocket 客户端
type Client struct {
	ID       string                 // 客户端ID
	UserID   string                 // 用户ID
	UserName string                 // 用户名
	Conn     *websocket.Conn        // WebSocket 连接
	Send     chan []byte            // 发送通道
	Manager  *Manager               // 管理器
	Rooms    map[string]bool        // 客户端加入的房间
	MetaData map[string]interface{} // 元数据
	mu       sync.RWMutex           // 读写锁
}

// NewClient 创建新的客户端
func NewClient(id, userID, userName string, conn *websocket.Conn, manager *Manager, log logger.Logger) *Client {
	return &Client{
		ID:       id,
		UserID:   userID,
		UserName: userName,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Manager:  manager,
		Rooms:    make(map[string]bool),
		MetaData: make(map[string]interface{}),
	}
}

// ReadPump 读取泵
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Manager.logger.Error("unexpected close error", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(data, &message); err != nil {
			c.Manager.logger.Error("failed to unmarshal message", err)
			continue
		}

		// 使用事件系统处理消息
		c.Manager.HandleMessageWithEvents(c, &message)
	}
}

// WritePump 写入泵
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 管理器关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送单条消息（不批量）
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.Send <- data:
		return nil
	default:
		close(c.Send)
		c.Manager.Unregister(c)
		return errors.New("client send buffer is full")
	}
}

// JoinRoom 加入房间
func (c *Client) JoinRoom(room string) {
	c.mu.Lock()
	c.Rooms[room] = true
	c.mu.Unlock()

	c.Manager.AddClientToRoomWithEvents(room, c)
}

// LeaveRoom 离开房间
func (c *Client) LeaveRoom(room string) {
	c.mu.Lock()
	delete(c.Rooms, room)
	c.mu.Unlock()

	c.Manager.RemoveClientFromRoomWithEvents(room, c)
}

// GetRooms 获取客户端加入的所有房间
func (c *Client) GetRooms() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rooms := make([]string, 0, len(c.Rooms))
	for room := range c.Rooms {
		rooms = append(rooms, room)
	}

	return rooms
}

// InRoom 检查客户端是否在房间中
func (c *Client) InRoom(room string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.Rooms[room]
	return ok
}

// SetMetaData 设置元数据
func (c *Client) SetMetaData(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.MetaData[key] = value
}

// GetMetaData 获取元数据
func (c *Client) GetMetaData(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.MetaData[key]
	return value, ok
}

// GetAllMetaData 获取所有元数据
func (c *Client) GetAllMetaData() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 创建副本
	result := make(map[string]interface{})
	for k, v := range c.MetaData {
		result[k] = v
	}

	return result
}