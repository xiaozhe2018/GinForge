# WebSocket Session 管理

WebSocket Session 管理模块提供了对 WebSocket 客户端会话数据的管理功能。

## 功能特性

- 会话数据的存储和获取
- 支持内存存储和 Redis 持久化
- 用户会话关联和查询
- 会话数据的序列化和反序列化

## 使用方法

### 创建会话管理器

```go
// 使用内存存储
sessionManager := session.NewMemorySessionManager()

// 或使用 Redis 存储
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
sessionManager := session.NewRedisSessionManager(redisClient, "app", 24*time.Hour)
```

### 会话操作

```go
// 创建会话
session, err := sessionManager.CreateSession(clientID, userID)

// 获取会话
session, err := sessionManager.GetSession(clientID)

// 设置会话数据
err := sessionManager.SetSessionData(clientID, "preferences", map[string]interface{}{
    "theme": "dark",
    "notifications": true,
})

// 获取会话数据
preferences, err := sessionManager.GetSessionData(clientID, "preferences")

// 更新会话数据
err := sessionManager.UpdateSessionData(clientID, map[string]interface{}{
    "lastActive": time.Now(),
})

// 删除会话数据
err := sessionManager.DeleteSessionData(clientID, "tempData")

// 清空会话数据
err := sessionManager.ClearSessionData(clientID)

// 移除会话
err := sessionManager.RemoveSession(clientID)

// 获取用户的所有会话
sessions, err := sessionManager.GetUserSessions(userID)
```

## 与 WebSocket Manager 集成

```go
// 在 WebSocket Manager 中添加会话管理器
type Manager struct {
    // 现有字段
    clients     map[string]*Client
    userClients map[string]map[string]*Client
    rooms       map[string]map[string]*Client
    
    // 新增字段
    sessionManager *session.SessionManager
}

// 初始化
func NewManager(log logger.Logger, redisClient *redis.Client) *Manager {
    m := &Manager{
        // 现有初始化
        clients:     make(map[string]*Client),
        userClients: make(map[string]map[string]*Client),
        rooms:       make(map[string]map[string]*Client),
        
        // 新增会话管理器
        sessionManager: session.NewRedisSessionManager(redisClient, "ws", 24*time.Hour),
    }
    
    return m
}

// 在客户端注册时创建会话
func (m *Manager) registerClient(client *Client) {
    // 现有代码
    
    // 创建会话
    _, err := m.sessionManager.CreateSession(client.ID, client.UserID)
    if err != nil {
        m.logger.Error("failed to create session", err)
    }
}

// 在客户端注销时移除会话
func (m *Manager) unregisterClient(client *Client) {
    // 现有代码
    
    // 移除会话
    if err := m.sessionManager.RemoveSession(client.ID); err != nil {
        m.logger.Error("failed to remove session", err)
    }
}
```

## 会话数据持久化

使用 Redis 作为会话存储时，会话数据会持久化到 Redis 中，支持以下特性：

- 会话数据 TTL 自动过期
- 用户会话关联索引
- 分布式环境下的会话共享
- 服务重启后会话恢复

Redis 存储结构：

- `{prefix}:session:{sessionID}`: 会话数据
- `{prefix}:user_sessions:{userID}`: 用户会话索引（Set 类型）
