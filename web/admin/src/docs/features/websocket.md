# WebSocket 实时通信

GinForge 提供了完整的 WebSocket 解决方案，支持实时通知、消息推送、在线状态等功能。

## 🎯 功能特性

- ✅ **实时通信**：双向实时消息传输
- ✅ **连接管理**：自动重连、心跳检测
- ✅ **房间管理**：支持多房间广播
- ✅ **用户管理**：单用户多端登录
- ✅ **消息类型**：系统消息、用户消息、通知、数据更新
- ✅ **Redis PubSub**：支持分布式部署

## 🏗️ 架构设计

```
Client (浏览器)
    ↓ WebSocket
Gateway (8087) - WebSocket 服务
    ↓ Redis PubSub
Admin-API (8083) - 发布消息
User-API (8081) - 发布消息
```

## 🔧 启动 WebSocket 服务

```bash
# 启动 WebSocket Gateway
go run ./services/websocket-gateway/cmd/server/main.go

# 服务运行在 8087 端口
```

## 📝 前端连接

### Vue 3 使用示例

```typescript
// stores/websocket.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWebSocketStore = defineStore('websocket', () => {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const messages = ref<any[]>([])
  
  // 连接 WebSocket
  const connect = (token: string) => {
    const wsUrl = `ws://localhost:8087/ws?token=${token}`
    ws.value = new WebSocket(wsUrl)
    
    // 连接成功
    ws.value.onopen = () => {
      console.log('WebSocket connected')
      connected.value = true
    }
    
    // 接收消息
    ws.value.onmessage = (event) => {
      const data = JSON.parse(event.data)
      messages.value.push(data)
      
      // 处理不同类型的消息
      handleMessage(data)
    }
    
    // 连接关闭
    ws.value.onclose = () => {
      console.log('WebSocket closed')
      connected.value = false
      
      // 自动重连（5秒后）
      setTimeout(() => connect(token), 5000)
    }
    
    // 错误处理
    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }
  
  // 发送消息
  const send = (message: any) => {
    if (ws.value && connected.value) {
      ws.value.send(JSON.stringify(message))
    }
  }
  
  // 断开连接
  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
      connected.value = false
    }
  }
  
  return {
    connected,
    messages,
    connect,
    send,
    disconnect
  }
})
```

### 在组件中使用

```vue
<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'

const wsStore = useWebSocketStore()

onMounted(() => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    wsStore.connect(token)
  }
})

onUnmounted(() => {
  wsStore.disconnect()
})
</script>
```

## 📨 消息类型

### 1. 系统消息

```json
{
  "type": "system",
  "content": {
    "message": "系统将于 30 分钟后维护",
    "level": "warning",
    "code": 1001
  },
  "timestamp": 1697356800
}
```

### 2. 通知消息

```json
{
  "type": "notification",
  "content": {
    "title": "用户登录提醒",
    "body": "admin 于 2025-10-15 14:30:00 登录了系统",
    "icon": "User",
    "link": "/system/users"
  },
  "timestamp": 1697356800
}
```

### 3. 数据更新消息

```json
{
  "type": "data_update",
  "content": {
    "entity": "user",
    "action": "create",
    "id": "123",
    "data": {
      "username": "new_user"
    }
  },
  "timestamp": 1697356800
}
```

### 4. 用户消息

```json
{
  "type": "user_message",
  "content": {
    "from": "user_123",
    "to": "user_456",
    "text": "Hello!"
  },
  "timestamp": 1697356800
}
```

## 🔔 发送通知

### 从后端发送系统通知

```go
// 使用通知客户端
notifyClient := notification.NewClient(redisClient)

// 发送系统通知
notifyClient.SendSystemNotification(ctx, &notification.SystemNotificationRequest{
    Title: "系统维护通知",
    Body:  "系统将于今晚 22:00 进行维护，预计 2 小时",
    Icon:  "Setting",
    Link:  "/system/maintenance",
})
```

### 发送给特定用户

```go
// 发送通知给用户 ID=123
notifyClient.SendUserNotification(ctx, &notification.UserNotificationRequest{
    UserID: "123",
    Title:  "订单通知",
    Body:   "您的订单已发货",
    Icon:   "ShoppingCart",
    Link:   "/orders/456",
})
```

## 🏠 房间管理

### 加入房间

```typescript
// 前端加入房间
wsStore.send({
  type: 'join_room',
  content: 'chat:room001'
})
```

### 房间广播

```go
// 后端向房间广播消息
manager.BroadcastToRoom("chat:room001", &websocket.Message{
    Type: websocket.MessageTypeRoomMessage,
    Content: map[string]interface{}{
        "text": "有新用户加入",
    },
})
```

### 离开房间

```typescript
wsStore.send({
  type: 'leave_room',
  content: 'chat:room001'
})
```

## 💬 聊天示例

### 前端聊天组件

```vue
<template>
  <div class="chat-container">
    <!-- 消息列表 -->
    <div class="messages">
      <div v-for="msg in messages" :key="msg.id" class="message">
        <span class="username">{{ msg.username }}:</span>
        <span class="text">{{ msg.text }}</span>
      </div>
    </div>
    
    <!-- 输入框 -->
    <el-input
      v-model="inputMessage"
      placeholder="输入消息..."
      @keyup.enter="sendMessage"
    >
      <template #append>
        <el-button @click="sendMessage">发送</el-button>
      </template>
    </el-input>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'

const wsStore = useWebSocketStore()
const messages = ref<any[]>([])
const inputMessage = ref('')

// 发送消息
const sendMessage = () => {
  if (!inputMessage.value.trim()) return
  
  wsStore.send({
    type: 'room_message',
    content: {
      room: 'chat:room001',
      text: inputMessage.value
    }
  })
  
  inputMessage.value = ''
}

// 监听消息
onMounted(() => {
  // 加入房间
  wsStore.send({
    type: 'join_room',
    content: 'chat:room001'
  })
})
</script>
```

## 📡 处理实时更新

### 监听数据更新

```vue
<script setup lang="ts">
import { watch } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'

const wsStore = useWebSocketStore()

// 监听 WebSocket 消息
watch(() => wsStore.messages, (messages) => {
  const lastMessage = messages[messages.length - 1]
  
  if (lastMessage?.type === 'data_update') {
    const { entity, action, data } = lastMessage.content
    
    // 处理数据更新
    if (entity === 'user' && action === 'create') {
      ElMessage.success('有新用户加入')
      // 刷新用户列表
      loadUsers()
    }
    
    if (entity === 'order' && action === 'update') {
      ElMessage.info('订单状态已更新')
      // 更新订单状态
      updateOrderStatus(data)
    }
  }
  
  if (lastMessage?.type === 'notification') {
    const { title, body } = lastMessage.content
    ElNotification({
      title,
      message: body,
      type: 'info'
    })
  }
}, { deep: true })
</script>
```

## 🔒 连接认证

### 后端 JWT 验证

```go
// WebSocket 连接时验证 Token
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    // 从查询参数获取 Token
    token := c.Query("token")
    if token == "" {
        c.JSON(401, gin.H{"error": "未提供令牌"})
        return
    }
    
    // 验证 Token
    claims, err := validateJWT(token, jwtSecret)
    if err != nil {
        c.JSON(401, gin.H{"error": "令牌无效"})
        return
    }
    
    // 升级为 WebSocket 连接
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    
    // 创建客户端
    client := websocket.NewClient(
        conn,
        claims["user_id"].(string),
        claims["username"].(string),
        h.manager,
    )
    
    // 注册客户端
    h.manager.Register(client)
    
    // 开始处理消息
    go client.ReadPump()
    go client.WritePump()
}
```

## 💡 使用技巧

### 1. 心跳检测

```typescript
// 前端定时发送心跳
setInterval(() => {
  wsStore.send({ type: 'ping' })
}, 30000)  // 每 30 秒
```

### 2. 自动重连

```typescript
let reconnectAttempts = 0
const maxReconnectAttempts = 5

ws.onclose = () => {
  if (reconnectAttempts < maxReconnectAttempts) {
    setTimeout(() => {
      reconnectAttempts++
      connect(token)
    }, 5000)
  }
}
```

### 3. 消息队列

```typescript
// 离线消息队列
const messageQueue: any[] = []

const send = (message: any) => {
  if (connected.value) {
    ws.value?.send(JSON.stringify(message))
  } else {
    // 未连接时放入队列
    messageQueue.push(message)
  }
}

// 连接成功后发送队列中的消息
ws.onopen = () => {
  while (messageQueue.length > 0) {
    const msg = messageQueue.shift()
    ws.value?.send(JSON.stringify(msg))
  }
}
```

## 📚 完整示例

查看完整实现：

- **WebSocket 管理器**: `pkg/websocket/manager.go`
- **WebSocket 客户端**: `pkg/websocket/client.go`
- **WebSocket 服务**: `services/websocket-gateway/`
- **详细架构**: `docs/WEBSOCKET_ARCHITECTURE.md`
- **通知服务**: `docs/NOTIFICATION_SERVICE.md`

## 🎯 典型应用场景

### 1. 实时通知

- 系统维护通知
- 用户行为提醒
- 订单状态更新

### 2. 即时聊天

- 客服系统
- 团队协作
- 评论互动

### 3. 数据同步

- 多端数据同步
- 协同编辑
- 实时监控

### 4. 在线状态

- 用户在线状态
- 设备连接状态
- 服务健康状态

## 🎯 下一步

- [消息队列](../api-reference/queue) - 异步任务处理
- [Redis 使用](../../demo/redis_usage.md) - Redis 高级功能
- [通知服务](../../NOTIFICATION_SERVICE.md) - 邮件和短信通知

---

**提示**: WebSocket 非常适合需要实时性的场景，但要注意连接数和服务器资源消耗。

