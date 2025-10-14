package group

import (
	"context"
	"errors"
	"sync"
	"time"
	
	"github.com/redis/go-redis/v9"
)

var (
	// ErrGroupNotFound 分组不存在
	ErrGroupNotFound = errors.New("group not found")
	
	// ErrClientNotFound 客户端不存在
	ErrClientNotFound = errors.New("client not found")
)

// GroupStore 分组存储接口
type GroupStore interface {
	// AddClientToGroup 添加客户端到分组
	AddClientToGroup(ctx context.Context, groupID, clientID string, clientInfo map[string]interface{}) error
	
	// RemoveClientFromGroup 从分组移除客户端
	RemoveClientFromGroup(ctx context.Context, groupID, clientID string) error
	
	// GetGroupClients 获取分组中的所有客户端
	GetGroupClients(ctx context.Context, groupID string) (map[string]map[string]interface{}, error)
	
	// GetClientGroups 获取客户端加入的所有分组
	GetClientGroups(ctx context.Context, clientID string) ([]string, error)
	
	// SetGroupMetadata 设置分组元数据
	SetGroupMetadata(ctx context.Context, groupID string, metadata map[string]interface{}) error
	
	// GetGroupMetadata 获取分组元数据
	GetGroupMetadata(ctx context.Context, groupID string) (map[string]interface{}, error)
	
	// GetAllGroups 获取所有分组
	GetAllGroups(ctx context.Context) ([]string, error)
	
	// GetGroupSize 获取分组大小
	GetGroupSize(ctx context.Context, groupID string) (int64, error)
	
	// GroupExists 检查分组是否存在
	GroupExists(ctx context.Context, groupID string) (bool, error)
	
	// ClientInGroup 检查客户端是否在分组中
	ClientInGroup(ctx context.Context, groupID, clientID string) (bool, error)
	
	// DeleteGroup 删除分组
	DeleteGroup(ctx context.Context, groupID string) error
}

// MemoryGroupStore 内存分组存储
type MemoryGroupStore struct {
	groups      map[string]map[string]map[string]interface{} // groupID -> clientID -> clientInfo
	clientGroups map[string]map[string]bool                  // clientID -> groupID -> true
	groupMeta   map[string]map[string]interface{}            // groupID -> metadata
	mu          sync.RWMutex
}

// NewMemoryGroupStore 创建内存分组存储
func NewMemoryGroupStore() *MemoryGroupStore {
	return &MemoryGroupStore{
		groups:      make(map[string]map[string]map[string]interface{}),
		clientGroups: make(map[string]map[string]bool),
		groupMeta:   make(map[string]map[string]interface{}),
	}
}

// AddClientToGroup 添加客户端到分组
func (m *MemoryGroupStore) AddClientToGroup(ctx context.Context, groupID, clientID string, clientInfo map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 初始化分组（如果不存在）
	if _, ok := m.groups[groupID]; !ok {
		m.groups[groupID] = make(map[string]map[string]interface{})
	}
	
	// 添加客户端到分组
	m.groups[groupID][clientID] = clientInfo
	
	// 初始化客户端分组（如果不存在）
	if _, ok := m.clientGroups[clientID]; !ok {
		m.clientGroups[clientID] = make(map[string]bool)
	}
	
	// 添加分组到客户端分组
	m.clientGroups[clientID][groupID] = true
	
	return nil
}

// RemoveClientFromGroup 从分组移除客户端
func (m *MemoryGroupStore) RemoveClientFromGroup(ctx context.Context, groupID, clientID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查分组是否存在
	groupClients, ok := m.groups[groupID]
	if !ok {
		return nil // 分组不存在，视为成功
	}
	
	// 从分组中移除客户端
	delete(groupClients, clientID)
	
	// 如果分组为空，删除分组
	if len(groupClients) == 0 {
		delete(m.groups, groupID)
		delete(m.groupMeta, groupID)
	}
	
	// 从客户端分组中移除分组
	if clientGroups, ok := m.clientGroups[clientID]; ok {
		delete(clientGroups, groupID)
		
		// 如果客户端没有分组，删除客户端
		if len(clientGroups) == 0 {
			delete(m.clientGroups, clientID)
		}
	}
	
	return nil
}

// GetGroupClients 获取分组中的所有客户端
func (m *MemoryGroupStore) GetGroupClients(ctx context.Context, groupID string) (map[string]map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	groupClients, ok := m.groups[groupID]
	if !ok {
		return nil, ErrGroupNotFound
	}
	
	// 创建副本
	result := make(map[string]map[string]interface{})
	for clientID, clientInfo := range groupClients {
		clientInfoCopy := make(map[string]interface{})
		for k, v := range clientInfo {
			clientInfoCopy[k] = v
		}
		result[clientID] = clientInfoCopy
	}
	
	return result, nil
}

// GetClientGroups 获取客户端加入的所有分组
func (m *MemoryGroupStore) GetClientGroups(ctx context.Context, clientID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	clientGroups, ok := m.clientGroups[clientID]
	if !ok {
		return nil, nil // 客户端没有分组，返回空列表
	}
	
	// 创建分组列表
	groups := make([]string, 0, len(clientGroups))
	for groupID := range clientGroups {
		groups = append(groups, groupID)
	}
	
	return groups, nil
}

// SetGroupMetadata 设置分组元数据
func (m *MemoryGroupStore) SetGroupMetadata(ctx context.Context, groupID string, metadata map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查分组是否存在
	if _, ok := m.groups[groupID]; !ok {
		return ErrGroupNotFound
	}
	
	// 设置分组元数据
	m.groupMeta[groupID] = metadata
	
	return nil
}

// GetGroupMetadata 获取分组元数据
func (m *MemoryGroupStore) GetGroupMetadata(ctx context.Context, groupID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	metadata, ok := m.groupMeta[groupID]
	if !ok {
		return nil, nil // 分组没有元数据，返回空映射
	}
	
	// 创建副本
	result := make(map[string]interface{})
	for k, v := range metadata {
		result[k] = v
	}
	
	return result, nil
}

// GetAllGroups 获取所有分组
func (m *MemoryGroupStore) GetAllGroups(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	groups := make([]string, 0, len(m.groups))
	for groupID := range m.groups {
		groups = append(groups, groupID)
	}
	
	return groups, nil
}

// GetGroupSize 获取分组大小
func (m *MemoryGroupStore) GetGroupSize(ctx context.Context, groupID string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	groupClients, ok := m.groups[groupID]
	if !ok {
		return 0, nil // 分组不存在，返回 0
	}
	
	return int64(len(groupClients)), nil
}

// GroupExists 检查分组是否存在
func (m *MemoryGroupStore) GroupExists(ctx context.Context, groupID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, ok := m.groups[groupID]
	return ok, nil
}

// ClientInGroup 检查客户端是否在分组中
func (m *MemoryGroupStore) ClientInGroup(ctx context.Context, groupID, clientID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	groupClients, ok := m.groups[groupID]
	if !ok {
		return false, nil // 分组不存在，返回 false
	}
	
	_, ok = groupClients[clientID]
	return ok, nil
}

// DeleteGroup 删除分组
func (m *MemoryGroupStore) DeleteGroup(ctx context.Context, groupID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查分组是否存在
	groupClients, ok := m.groups[groupID]
	if !ok {
		return nil // 分组不存在，视为成功
	}
	
	// 从每个客户端的分组列表中移除分组
	for clientID := range groupClients {
		if clientGroups, ok := m.clientGroups[clientID]; ok {
			delete(clientGroups, groupID)
			
			// 如果客户端没有分组，删除客户端
			if len(clientGroups) == 0 {
				delete(m.clientGroups, clientID)
			}
		}
	}
	
	// 删除分组和元数据
	delete(m.groups, groupID)
	delete(m.groupMeta, groupID)
	
	return nil
}

// GroupManager 分组管理器
type GroupManager struct {
	store  GroupStore
	ctx    context.Context
	cancel context.CancelFunc
}

// NewGroupManager 创建分组管理器
func NewGroupManager(store GroupStore) *GroupManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &GroupManager{
		store:  store,
		ctx:    ctx,
		cancel: cancel,
	}
}

// NewMemoryGroupManager 创建内存分组管理器
func NewMemoryGroupManager() *GroupManager {
	return NewGroupManager(NewMemoryGroupStore())
}

// NewRedisGroupManager 创建 Redis 分组管理器
func NewRedisGroupManager(redisClient *redis.Client, prefix string, ttl time.Duration) *GroupManager {
	return NewGroupManager(NewRedisGroupStore(redisClient, prefix, ttl))
}

// JoinGroup 加入分组
func (m *GroupManager) JoinGroup(groupID, clientID string, clientInfo map[string]interface{}) error {
	return m.store.AddClientToGroup(m.ctx, groupID, clientID, clientInfo)
}

// LeaveGroup 离开分组
func (m *GroupManager) LeaveGroup(groupID, clientID string) error {
	return m.store.RemoveClientFromGroup(m.ctx, groupID, clientID)
}

// GetGroupMembers 获取分组成员
func (m *GroupManager) GetGroupMembers(groupID string) (map[string]map[string]interface{}, error) {
	return m.store.GetGroupClients(m.ctx, groupID)
}

// GetClientGroups 获取客户端分组
func (m *GroupManager) GetClientGroups(clientID string) ([]string, error) {
	return m.store.GetClientGroups(m.ctx, clientID)
}

// SetGroupMetadata 设置分组元数据
func (m *GroupManager) SetGroupMetadata(groupID string, metadata map[string]interface{}) error {
	return m.store.SetGroupMetadata(m.ctx, groupID, metadata)
}

// GetGroupMetadata 获取分组元数据
func (m *GroupManager) GetGroupMetadata(groupID string) (map[string]interface{}, error) {
	return m.store.GetGroupMetadata(m.ctx, groupID)
}

// GetAllGroups 获取所有分组
func (m *GroupManager) GetAllGroups() ([]string, error) {
	return m.store.GetAllGroups(m.ctx)
}

// GetGroupSize 获取分组大小
func (m *GroupManager) GetGroupSize(groupID string) (int64, error) {
	return m.store.GetGroupSize(m.ctx, groupID)
}

// GroupExists 检查分组是否存在
func (m *GroupManager) GroupExists(groupID string) (bool, error) {
	return m.store.GroupExists(m.ctx, groupID)
}

// ClientInGroup 检查客户端是否在分组中
func (m *GroupManager) ClientInGroup(groupID, clientID string) (bool, error) {
	return m.store.ClientInGroup(m.ctx, groupID, clientID)
}

// DeleteGroup 删除分组
func (m *GroupManager) DeleteGroup(groupID string) error {
	return m.store.DeleteGroup(m.ctx, groupID)
}

// Close 关闭分组管理器
func (m *GroupManager) Close() {
	m.cancel()
}
