package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrClientNotFound 客户端不存在
	ErrClientNotFound = errors.New("client not found")
)

// SessionManager 会话管理器
type SessionManager struct {
	store      SessionStore
	clientSess map[string]string // 客户端ID -> 会话ID
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewSessionManager 创建会话管理器
func NewSessionManager(store SessionStore) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &SessionManager{
		store:      store,
		clientSess: make(map[string]string),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// NewMemorySessionManager 创建内存会话管理器
func NewMemorySessionManager() *SessionManager {
	return NewSessionManager(NewMemorySessionStore())
}

// NewRedisSessionManager 创建Redis会话管理器
func NewRedisSessionManager(client *redis.Client, prefix string, ttl time.Duration) *SessionManager {
	return NewSessionManager(NewRedisSessionStore(client, prefix, ttl))
}

// CreateSession 创建会话
func (m *SessionManager) CreateSession(clientID, userID string) (*Session, error) {
	session := NewSession(clientID, userID)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 保存会话
	if err := m.store.Save(m.ctx, session); err != nil {
		return nil, err
	}
	
	// 关联客户端与会话
	m.clientSess[clientID] = session.ID
	
	return session, nil
}

// GetSession 获取会话
func (m *SessionManager) GetSession(clientID string) (*Session, error) {
	m.mu.RLock()
	sessionID, ok := m.clientSess[clientID]
	m.mu.RUnlock()
	
	if !ok {
		return nil, ErrClientNotFound
	}
	
	return m.store.Get(m.ctx, sessionID)
}

// RemoveSession 移除会话
func (m *SessionManager) RemoveSession(clientID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	sessionID, ok := m.clientSess[clientID]
	if !ok {
		return nil // 客户端不存在，视为成功
	}
	
	// 删除会话
	err := m.store.Delete(m.ctx, sessionID)
	
	// 无论是否成功，都移除客户端关联
	delete(m.clientSess, clientID)
	
	return err
}

// GetUserSessions 获取用户的所有会话
func (m *SessionManager) GetUserSessions(userID string) ([]*Session, error) {
	return m.store.GetByUserID(m.ctx, userID)
}

// SetSessionData 设置会话数据
func (m *SessionManager) SetSessionData(clientID, key string, value interface{}) error {
	session, err := m.GetSession(clientID)
	if err != nil {
		return err
	}
	
	session.Set(key, value)
	return m.store.Save(m.ctx, session)
}

// GetSessionData 获取会话数据
func (m *SessionManager) GetSessionData(clientID, key string) (interface{}, error) {
	session, err := m.GetSession(clientID)
	if err != nil {
		return nil, err
	}
	
	value, ok := session.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}
	
	return value, nil
}

// UpdateSessionData 更新会话数据
func (m *SessionManager) UpdateSessionData(clientID string, data map[string]interface{}) error {
	session, err := m.GetSession(clientID)
	if err != nil {
		return err
	}
	
	session.Update(data)
	return m.store.Save(m.ctx, session)
}

// DeleteSessionData 删除会话数据
func (m *SessionManager) DeleteSessionData(clientID, key string) error {
	session, err := m.GetSession(clientID)
	if err != nil {
		return err
	}
	
	session.Delete(key)
	return m.store.Save(m.ctx, session)
}

// ClearSessionData 清空会话数据
func (m *SessionManager) ClearSessionData(clientID string) error {
	session, err := m.GetSession(clientID)
	if err != nil {
		return err
	}
	
	session.Clear()
	return m.store.Save(m.ctx, session)
}

// Close 关闭会话管理器
func (m *SessionManager) Close() {
	m.cancel()
}
