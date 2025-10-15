package session

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrSessionNotFound 会话不存在
	ErrSessionNotFound = errors.New("session not found")

	// ErrKeyNotFound 键不存在
	ErrKeyNotFound = errors.New("key not found in session")

	// ErrInvalidSessionData 无效的会话数据
	ErrInvalidSessionData = errors.New("invalid session data")
)

// SessionData 会话数据
type SessionData map[string]interface{}

// Session 会话
type Session struct {
	ID        string      // 会话ID
	UserID    string      // 用户ID
	Data      SessionData // 会话数据
	CreatedAt time.Time   // 创建时间
	UpdatedAt time.Time   // 更新时间
	mu        sync.RWMutex
}

// NewSession 创建新会话
func NewSession(id, userID string) *Session {
	now := time.Now()
	return &Session{
		ID:        id,
		UserID:    userID,
		Data:      make(SessionData),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Set 设置会话数据
func (s *Session) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data[key] = value
	s.UpdatedAt = time.Now()
}

// Get 获取会话数据
func (s *Session) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.Data[key]
	return value, ok
}

// Delete 删除会话数据
func (s *Session) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Data, key)
	s.UpdatedAt = time.Now()
}

// Clear 清空会话数据
func (s *Session) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data = make(SessionData)
	s.UpdatedAt = time.Now()
}

// GetAll 获取所有会话数据
func (s *Session) GetAll() SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建副本
	data := make(SessionData, len(s.Data))
	for k, v := range s.Data {
		data[k] = v
	}
	return data
}

// Update 批量更新会话数据
func (s *Session) Update(data SessionData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range data {
		s.Data[k] = v
	}
	s.UpdatedAt = time.Now()
}

// MarshalJSON 序列化为JSON
func (s *Session) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return json.Marshal(map[string]interface{}{
		"id":         s.ID,
		"user_id":    s.UserID,
		"data":       s.Data,
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
	})
}

// UnmarshalJSON 从JSON反序列化
func (s *Session) UnmarshalJSON(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var temp struct {
		ID        string          `json:"id"`
		UserID    string          `json:"user_id"`
		Data      json.RawMessage `json:"data"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	s.ID = temp.ID
	s.UserID = temp.UserID
	s.CreatedAt = temp.CreatedAt
	s.UpdatedAt = temp.UpdatedAt

	// 解析会话数据
	var sessionData SessionData
	if err := json.Unmarshal(temp.Data, &sessionData); err != nil {
		return err
	}
	s.Data = sessionData

	return nil
}

// SessionStore 会话存储接口
type SessionStore interface {
	// Get 获取会话
	Get(ctx context.Context, id string) (*Session, error)

	// Save 保存会话
	Save(ctx context.Context, session *Session) error

	// Delete 删除会话
	Delete(ctx context.Context, id string) error

	// GetByUserID 根据用户ID获取会话
	GetByUserID(ctx context.Context, userID string) ([]*Session, error)
}

// MemorySessionStore 内存会话存储
type MemorySessionStore struct {
	sessions  map[string]*Session        // 会话ID -> 会话
	userIndex map[string]map[string]bool // 用户ID -> 会话ID集合
	mu        sync.RWMutex
}

// NewMemorySessionStore 创建内存会话存储
func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions:  make(map[string]*Session),
		userIndex: make(map[string]map[string]bool),
	}
}

// Get 获取会话
func (m *MemorySessionStore) Get(ctx context.Context, id string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

// Save 保存会话
func (m *MemorySessionStore) Save(ctx context.Context, session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 保存会话
	m.sessions[session.ID] = session

	// 更新用户索引
	if session.UserID != "" {
		if _, ok := m.userIndex[session.UserID]; !ok {
			m.userIndex[session.UserID] = make(map[string]bool)
		}
		m.userIndex[session.UserID][session.ID] = true
	}

	return nil
}

// Delete 删除会话
func (m *MemorySessionStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		return nil // 已经不存在，视为成功
	}

	// 从用户索引中删除
	if session.UserID != "" {
		if userSessions, ok := m.userIndex[session.UserID]; ok {
			delete(userSessions, id)
			if len(userSessions) == 0 {
				delete(m.userIndex, session.UserID)
			}
		}
	}

	// 删除会话
	delete(m.sessions, id)
	return nil
}

// GetByUserID 根据用户ID获取会话
func (m *MemorySessionStore) GetByUserID(ctx context.Context, userID string) ([]*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessionIDs, ok := m.userIndex[userID]
	if !ok {
		return nil, nil // 用户没有会话
	}

	sessions := make([]*Session, 0, len(sessionIDs))
	for id := range sessionIDs {
		if session, ok := m.sessions[id]; ok {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

// RedisSessionStore Redis会话存储
type RedisSessionStore struct {
	client redis.UniversalClient
	prefix string
	ttl    time.Duration
}

// NewRedisSessionStore 创建Redis会话存储
func NewRedisSessionStore(client redis.UniversalClient, prefix string, ttl time.Duration) *RedisSessionStore {
	return &RedisSessionStore{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}
}

// sessionKey 生成会话键
func (r *RedisSessionStore) sessionKey(id string) string {
	return r.prefix + ":session:" + id
}

// userSessionsKey 生成用户会话索引键
func (r *RedisSessionStore) userSessionsKey(userID string) string {
	return r.prefix + ":user_sessions:" + userID
}

// Get 获取会话
func (r *RedisSessionStore) Get(ctx context.Context, id string) (*Session, error) {
	data, err := r.client.Get(ctx, r.sessionKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Save 保存会话
func (r *RedisSessionStore) Save(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// 使用管道执行多个命令
	pipe := r.client.Pipeline()

	// 保存会话数据
	pipe.Set(ctx, r.sessionKey(session.ID), data, r.ttl)

	// 更新用户会话索引
	if session.UserID != "" {
		pipe.SAdd(ctx, r.userSessionsKey(session.UserID), session.ID)
		pipe.Expire(ctx, r.userSessionsKey(session.UserID), r.ttl)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// Delete 删除会话
func (r *RedisSessionStore) Delete(ctx context.Context, id string) error {
	// 先获取会话，以便从用户索引中删除
	session, err := r.Get(ctx, id)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil // 已经不存在，视为成功
		}
		return err
	}

	// 使用管道执行多个命令
	pipe := r.client.Pipeline()

	// 删除会话
	pipe.Del(ctx, r.sessionKey(id))

	// 从用户会话索引中删除
	if session.UserID != "" {
		pipe.SRem(ctx, r.userSessionsKey(session.UserID), id)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// GetByUserID 根据用户ID获取会话
func (r *RedisSessionStore) GetByUserID(ctx context.Context, userID string) ([]*Session, error) {
	// 获取用户的所有会话ID
	sessionIDs, err := r.client.SMembers(ctx, r.userSessionsKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	if len(sessionIDs) == 0 {
		return nil, nil
	}

	// 使用管道批量获取会话
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(sessionIDs))

	for i, id := range sessionIDs {
		cmds[i] = pipe.Get(ctx, r.sessionKey(id))
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 解析会话数据
	sessions := make([]*Session, 0, len(sessionIDs))
	for _, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				continue // 跳过不存在的会话
			}
			return nil, err
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}
