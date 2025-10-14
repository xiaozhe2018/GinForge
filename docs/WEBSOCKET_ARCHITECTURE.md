# 🌐 GinForge WebSocket 架构设计

## 🏗️ 架构概述

GinForge 采用**集中式 WebSocket Gateway 架构**，将 WebSocket 功能集成到 API Gateway 中，实现统一的实时通信入口。

### 设计原则

1. **职责单一** - Gateway 负责所有实时通信（HTTP + WebSocket）
2. **服务解耦** - 通过 Redis PubSub 解耦服务间通信
3. **水平扩展** - 支持多个 Gateway 实例负载均衡
4. **统一入口** - 对外只暴露一个端口（8080）

---

## 📐 架构图

```
┌──────────────────────────────────────────────────────────────┐
│                         客户端层                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │
│  │ Web前端  │  │ Mobile  │  │ Desktop │  │ 其他客户端│       │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘       │
│       │            │             │            │             │
│       │ WS://      │ WS://       │ WS://      │ WS://      │
│       └────────────┴─────────────┴────────────┘             │
└──────────────────────┬───────────────────────────────────────┘
                       │
                       ↓
┌──────────────────────────────────────────────────────────────┐
│                    Gateway (8080)                             │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │             WebSocket 管理器                              │ │
│ │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │ │
│ │  │ 连接管理    │  │ 房间管理    │  │ 消息路由    │     │ │
│ │  │ (Clients)   │  │ (Rooms)     │  │ (Broadcast) │     │ │
│ │  └─────────────┘  └─────────────┘  └─────────────┘     │ │
│ └──────────────────────┬───────────────────────────────────┘ │
│                        │                                      │
│ ┌──────────────────────┴───────────────────────────────────┐ │
│ │             Redis PubSub 订阅                             │ │
│ │  订阅频道:                                                │ │
│ │  - websocket:broadcast       (广播消息)                  │ │
│ │  - websocket:user:*          (用户消息)                  │ │
│ │  - websocket:room:*          (房间消息)                  │ │
│ │  - websocket:notification:*  (通知消息)                  │ │
│ └──────────────────────────────────────────────────────────┘ │
└──────────────────────┬───────────────────────────────────────┘
                       │
                       ↓
               ┌───────────────┐
               │ Redis PubSub  │ (消息总线)
               └───────┬───────┘
                       │
        ┌──────────────┼──────────────┬──────────────┐
        │              │              │              │
        ↓              ↓              ↓              ↓
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  admin-api  │ │  user-api   │ │  file-api   │ │gateway-worker│
│  (8083)     │ │  (8081)     │ │  (8086)     │ │  (8084)     │
│             │ │             │ │             │ │             │
│ 发布消息:   │ │ 发布消息:   │ │ 发布消息:   │ │ 发布消息:   │
│ - 用户变更  │ │ - 订单状态  │ │ - 上传完成  │ │ - 异步任务  │
│ - 权限变更  │ │ - 用户通知  │ │ - 处理进度  │ │ - 定时提醒  │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘
```

---

## 🔄 消息流转

### 场景 1: 用户登录后接收实时通知

```
1. 用户登录 admin-api
   ↓
2. 前端携带 Token 连接 Gateway WebSocket
   WS://localhost:8080/ws?token=xxx
   ↓
3. Gateway 验证 Token，建立连接
   ↓
4. admin-api 创建用户时发布消息:
   Redis PUBLISH websocket:broadcast {
     "type": "data_update",
     "content": {"entity": "user", "action": "create"}
   }
   ↓
5. Gateway 接收 Redis 消息并广播
   ↓
6. 所有连接的客户端实时收到更新通知
```

### 场景 2: 发送通知给特定用户

```
1. admin-api 需要通知用户 ID=123
   ↓
2. admin-api 发布消息:
   Redis PUBLISH websocket:notification:123 {
     "type": "notification",
     "content": {"title": "系统通知", "body": "您的订单已发货"}
   }
   ↓
3. Gateway 接收消息
   ↓
4. Gateway 查找用户 123 的所有连接
   ↓
5. 发送给用户的所有客户端（支持多端登录）
```

### 场景 3: 房间广播

```
1. 用户加入聊天室 "chat:room001"
   客户端发送: {"type": "join_room", "content": "chat:room001"}
   ↓
2. Gateway 将用户添加到房间
   ↓
3. 其他服务发布房间消息:
   Redis PUBLISH websocket:room:chat:room001 {
     "type": "room_message",
     "content": {"text": "有新用户加入"}
   }
   ↓
4. Gateway 广播给房间内所有用户
```

---

## 🔑 核心组件

### 1. WebSocket 管理器 (`pkg/websocket/manager.go`)

**职责**:
- 管理所有 WebSocket 连接
- 维护用户与连接的映射关系
- 管理房间和频道
- 消息路由和广播

**核心方法**:
```go
// 注册/注销客户端
Register(client *Client)
Unregister(client *Client)

// 消息发送
Broadcast(message *Message)                    // 广播给所有人
SendToUser(userID string, message *Message)    // 发送给用户
BroadcastToRoom(room string, message *Message) // 发送给房间
SendNotification(userID string, notification *NotificationMessage) // 发送通知

// 状态查询
IsUserOnline(userID string) bool
GetOnlineUsers() []string
GetStats() map[string]interface{}
```

### 2. Redis PubSub 服务 (`services/gateway/internal/service/pubsub_service.go`)

**职责**:
- 订阅 Redis 频道
- 接收其他服务发布的消息
- 将消息转发给 WebSocket 管理器

**订阅频道**:
- `websocket:broadcast` - 广播消息
- `websocket:user:*` - 用户消息（支持通配符）
- `websocket:room:*` - 房间消息
- `websocket:notification:*` - 通知消息

### 3. 通知客户端 (`pkg/notification/client.go`)

**职责**:
- 供其他服务调用
- 发布消息到 Redis
- 简化 WebSocket 消息发送

**使用示例**:
```go
// 在任何服务中使用
import "goweb/pkg/notification"

notifyClient := notification.NewClient(redisClient)

// 发送通知给用户
notifyClient.SendNotification(ctx, "user123", &websocket.NotificationMessage{
    Title: "订单通知",
    Body: "您的订单已发货",
    Icon: "Truck",
    Link: "/orders/12345",
})

// 广播通知
notifyClient.BroadcastNotification(ctx, &websocket.NotificationMessage{
    Title: "系统维护通知",
    Body: "系统将于今晚22:00进行维护",
})
```

---

## 🔧 服务职责划分

### Gateway (8080) - 统一通信网关

**HTTP 代理**:
- ✅ 路由 HTTP 请求到各个微服务
- ✅ 负载均衡
- ✅ 限流和熔断

**WebSocket 管理** (NEW):
- ✅ WebSocket 连接管理
- ✅ 实时消息路由
- ✅ Redis PubSub 订阅
- ✅ 在线状态维护

### Gateway Worker (8084) - 异步任务处理器

**消息队列消费**:
- ✅ 订单提醒
- ✅ 用户通知
- ✅ 系统清理
- ✅ 支付重试
- ✅ 库存告警

**与 WebSocket 的协作**:
- 通过 Redis PubSub 发送实时通知
- 不直接管理 WebSocket 连接

### 其他微服务 (admin-api, user-api, file-api, etc.)

**通过 notification.Client 发送 WebSocket 消息**:
- ✅ 数据变更通知
- ✅ 操作结果通知
- ✅ 进度更新
- ✅ 系统消息

---

## 🎯 优势

### vs 独立 WebSocket 服务

| 特性 | 独立服务 | 集成到 Gateway |
|------|---------|---------------|
| 服务数量 | +1 个服务 | 0（使用现有）|
| 端口数量 | +1 个端口 | 0（共用8080）|
| 部署复杂度 | 高 | 低 |
| 维护成本 | 高 | 低 |
| 统一性 | 差（多个入口）| 好（统一入口）|
| 职责清晰度 | 中 | 高 |

### 扩展性

**水平扩展**:
```
                     负载均衡
                        │
          ┌─────────────┼─────────────┐
          ↓             ↓             ↓
      Gateway-1     Gateway-2     Gateway-3
      (WebSocket)   (WebSocket)   (WebSocket)
          │             │             │
          └─────────────┼─────────────┘
                        ↓
                   Redis PubSub
                   (消息同步)
```

- 多个 Gateway 实例通过 Redis PubSub 同步消息
- 客户端可以连接到任意 Gateway
- 消息可以到达所有 Gateway 上的客户端

---

## 📊 性能考虑

### 连接数

- 单个 Gateway: 1万+ 并发连接
- 多个 Gateway: 线性扩展

### 消息延迟

- 内存消息: < 1ms
- Redis PubSub: < 10ms
- 端到端: < 50ms

### 资源消耗

- 每个连接: ~10KB 内存
- 1万连接: ~100MB
- CPU: 主要在消息序列化

---

## 🔐 安全考虑

### 认证

- ✅ JWT Token 认证
- ✅ Token 过期自动断开
- ✅ Token 黑名单支持

### 权限

- ✅ 基于用户 ID 的消息隔离
- ✅ 房间权限控制
- ✅ 消息来源验证

### 防护

- ✅ 消息大小限制（512KB）
- ✅ 连接心跳检测（60秒超时）
- ✅ 自动断线重连
- ✅ 消息缓冲区限制（防止内存溢出）

---

## 🚀 使用指南

### 后端：发送 WebSocket 消息

```go
// services/admin-api/internal/service/xxx_service.go
import "goweb/pkg/notification"

type SomeService struct {
    notifyClient *notification.Client
}

func (s *SomeService) CreateUser(user *User) error {
    // 1. 创建用户
    // ...
    
    // 2. 发送 WebSocket 通知
    notification := &websocket.NotificationMessage{
        Title: "用户创建成功",
        Body:  fmt.Sprintf("用户 %s 已创建", user.Username),
        Icon:  "UserPlus",
        Link:  "/users/" + user.ID,
    }
    
    // 发送给特定用户
    s.notifyClient.SendNotification(ctx, adminUserID, notification)
    
    // 或广播给所有在线管理员
    s.notifyClient.BroadcastNotification(ctx, notification)
    
    return nil
}
```

### 前端：接收 WebSocket 消息

```typescript
// web/admin/src/utils/websocket.ts
import WebSocketClient from '@/utils/websocket'

// 创建客户端
const ws = new WebSocketClient()

// 连接（携带 JWT Token）
ws.connect(token)

// 监听通知
ws.on('notification', (message) => {
  ElNotification({
    title: message.content.title,
    message: message.content.body,
    type: 'info'
  })
})

// 监听数据更新
ws.on('data_update', (message) => {
  // 刷新列表
  refreshUserList()
})

// 发送消息
ws.send({
  type: 'chat',
  to: 'user123',
  content: { text: 'Hello' }
})
```

---

## 📝 消息协议

### 消息格式

```typescript
interface WebSocketMessage {
  type: string           // 消息类型
  id?: string            // 消息ID（可选）
  from?: string          // 发送者ID
  from_name?: string     // 发送者名称
  to?: string            // 接收者ID（私聊）
  room?: string          // 房间名称
  content: any           // 消息内容
  data?: object          // 额外数据
  timestamp: number      // 时间戳
}
```

### 消息类型

| 类型 | 说明 | 方向 |
|------|------|------|
| `welcome` | 欢迎消息 | 服务器 → 客户端 |
| `ping/pong` | 心跳 | 双向 |
| `notification` | 通知消息 | 服务器 → 客户端 |
| `chat` | 聊天消息 | 双向 |
| `broadcast` | 广播消息 | 服务器 → 客户端 |
| `join_room` | 加入房间 | 客户端 → 服务器 |
| `leave_room` | 离开房间 | 客户端 → 服务器 |
| `room_message` | 房间消息 | 双向 |
| `user_status` | 用户状态 | 服务器 → 客户端 |
| `data_update` | 数据更新 | 服务器 → 客户端 |
| `error` | 错误消息 | 服务器 → 客户端 |

---

## 🎯 应用场景

### 1. 实时通知系统

**场景**: 管理员在后台操作，需要实时通知相关用户

```go
// admin-api 创建用户后
notifyClient.BroadcastDataUpdate(ctx, "user", "create", user.ID, userData)

// 前端自动刷新用户列表
ws.on('data_update', (msg) => {
  if (msg.content.entity === 'user' && msg.content.action === 'create') {
    refreshUserList()
  }
})
```

### 2. 在线状态监控

**场景**: 查看当前在线的管理员

```vue
<template>
  <el-tag v-for="user in onlineUsers" :key="user">
    {{ user }} 在线
  </el-tag>
</template>

<script setup>
// 监听用户上线/下线
ws.on('user_status', (msg) => {
  if (msg.content.status === 'online') {
    onlineUsers.value.push(msg.content.user_id)
  } else {
    const index = onlineUsers.value.indexOf(msg.content.user_id)
    if (index > -1) onlineUsers.value.splice(index, 1)
  }
})
</script>
```

### 3. 文件上传进度

**场景**: 大文件上传时显示实时进度

```go
// file-api 上传过程中
for progress := 0; progress <= 100; progress += 10 {
    notifyClient.SendNotification(ctx, userID, &websocket.NotificationMessage{
        Title: "上传进度",
        Body:  fmt.Sprintf("已完成 %d%%", progress),
    })
    time.Sleep(time.Second)
}
```

### 4. 协同编辑

**场景**: 多人同时编辑同一个文档

```typescript
// 用户加入编辑房间
ws.joinRoom('doc:12345')

// 监听房间消息
ws.on('room_message', (msg) => {
  if (msg.room === 'doc:12345') {
    // 应用其他用户的编辑
    applyChanges(msg.content.changes)
  }
})

// 发送编辑内容
ws.sendToRoom('doc:12345', {
  type: 'edit',
  changes: [...edits]
})
```

---

## ⚡ 最佳实践

### 1. 消息去重

```go
// 使用消息 ID 防止重复处理
message.SetData("message_id", uuid.New().String())
```

### 2. 消息持久化

```go
// 离线消息存储到数据库
if !wsManager.IsUserOnline(userID) {
    saveOfflineMessage(userID, message)
}
```

### 3. 错误处理

```typescript
ws.onError((error) => {
  ElMessage.error('连接异常，正在重连...')
})

ws.onClose(() => {
  // 自动重连（内置）
})
```

### 4. 性能优化

```go
// 批量发送
messages := []*websocket.Message{msg1, msg2, msg3}
for _, msg := range messages {
    wsManager.Broadcast(msg)
}
```

---

## 📈 监控指标

### WebSocket 指标

```bash
# 获取统计信息
curl http://localhost:8080/ws/stats

# 响应
{
  "total_clients": 123,   # 总连接数
  "online_users": 89,     # 在线用户数
  "total_rooms": 5,       # 房间数
  "uptime": 3600          # 运行时间（秒）
}
```

### Prometheus 指标

- `websocket_connections_total` - 总连接数
- `websocket_messages_sent_total` - 发送消息数
- `websocket_messages_received_total` - 接收消息数
- `websocket_errors_total` - 错误数

---

## 🔄 部署策略

### 单实例部署

```yaml
# docker-compose.yml
gateway:
  ports:
    - "8080:8080"
  environment:
    - REDIS_ENABLED=true
```

### 多实例部署（高可用）

```yaml
gateway:
  deploy:
    replicas: 3          # 3个实例
  ports:
    - "8080-8082:8080"   # 映射多个端口
  
nginx:
  # Nginx 负载均衡
  upstream gateway_backend {
    server gateway1:8080
    server gateway2:8080
    server gateway3:8080
  }
```

**注意**: 多实例部署时，必须启用 Redis，否则消息无法在实例间同步。

---

## 🆚 对比其他方案

### 方案对比

| 特性 | 独立 WebSocket 服务 | 集成到 Gateway | 集成到业务服务 |
|------|-------------------|---------------|---------------|
| 服务数量 | +1 | 0 | 0 |
| 职责清晰 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| 部署复杂度 | 高 | 低 | 中 |
| 扩展性 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| 维护成本 | 高 | 低 | 中 |
| 推荐度 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |

**结论**: 集成到 Gateway 是最佳方案！

---

## 🎓 总结

GinForge 的 WebSocket 架构设计特点：

1. **统一网关** - HTTP + WebSocket 统一入口
2. **服务解耦** - 通过 Redis PubSub 解耦
3. **职责清晰** - Gateway 管理连接，其他服务发布消息
4. **易于扩展** - 支持水平扩展和多实例部署
5. **开发友好** - 简单的 API，易于集成

**适用场景**:
- ✅ 中小型应用（<10万并发连接）
- ✅ 需要实时通知的管理系统
- ✅ 协同办公应用
- ✅ 实时监控系统

**大规模场景**:
- 考虑使用专业的消息推送服务（如 Firebase, Pusher）
- 或使用消息队列（Kafka + WebSocket）

---

**更新时间**: 2025-10-14  
**架构版本**: v1.0

