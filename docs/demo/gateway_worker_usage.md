# Gateway Worker 使用示例

## 概述

Gateway Worker 是专门处理消息队列的工作服务，与 API Gateway 分离，负责消费各种业务消息。它提供了高可用、可扩展的消息处理能力。

## 1. 架构说明

### 服务分离

```
┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │ Gateway Worker  │
│   (HTTP 代理)    │    │  (消息队列消费)  │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     │
            ┌─────────────────┐
            │   Redis 队列    │
            └─────────────────┘
```

- **API Gateway**：处理 HTTP 请求，代理到后端服务
- **Gateway Worker**：消费消息队列，处理异步任务
- **Redis 队列**：消息存储和传递

### 消息类型

Gateway Worker 预定义了以下消息类型：

- `order.reminder`：订单提醒
- `user.notification`：用户通知
- `system.cleanup`：系统清理
- `payment.retry`：支付重试
- `inventory.alert`：库存告警

## 2. 基础使用

### 启动 Gateway Worker

```bash
# 直接运行
go run ./services/gateway-worker/cmd/server

# 或者构建后运行
go build -o bin/gateway-worker ./services/gateway-worker/cmd/server
./bin/gateway-worker
```

### 健康检查

```bash
# 健康检查
curl http://localhost:8084/healthz

# 就绪检查
curl http://localhost:8084/ready

# 指标信息
curl http://localhost:8084/metrics
```

## 3. 消息发送示例

### 发送订单提醒

```go
package main

import (
    "context"
    "time"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
)

func main() {
    cfg := config.New()
    log := logger.New("message-sender", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisConfig := cfg.GetRedisConfig()
    redisManager := redis.NewManager(&redisConfig, log)
    
    // 获取消息队列
    queue := redisManager.GetQueue()
    
    ctx := context.Background()
    
    // 发送订单提醒消息
    err := queue.Publish(ctx, "order.reminder", map[string]interface{}{
        "order_id": "12345",
        "user_id":  "67890",
        "type":     "payment_reminder",
        "message":  "请及时完成支付",
        "amount":   299.99,
    })
    
    if err != nil {
        log.Error("发送消息失败", err)
        return
    }
    
    log.Info("消息发送成功")
}
```

### 发送延时消息

```go
// 发送延时消息（24小时后）
err := queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
    "order_id": "12345",
    "user_id":  "67890",
    "type":     "order_timeout",
    "message":  "订单即将超时，请尽快支付",
}, 24*time.Hour)
```

## 4. 自定义消息处理器

### 创建自定义 Worker 服务

```go
package service

import (
    "context"
    "sync"
    "time"
    "goweb/pkg/base"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
)

type CustomWorkerService struct {
    *base.BaseService
    redisManager *redis.Manager
    consumers    map[string]context.CancelFunc
    mutex        sync.RWMutex
}

func NewCustomWorkerService(redisManager *redis.Manager, log logger.Logger) *CustomWorkerService {
    return &CustomWorkerService{
        BaseService:  base.NewBaseService(log),
        redisManager: redisManager,
        consumers:    make(map[string]context.CancelFunc),
    }
}

// 启动自定义消费者
func (s *CustomWorkerService) StartCustomConsumers(ctx context.Context) error {
    queue := s.redisManager.GetQueue()
    
    // 自定义消息类型
    customTopics := []string{
        "email.send",
        "sms.send",
        "push.notification",
        "file.process",
        "data.sync",
    }
    
    for _, topic := range customTopics {
        if err := s.startConsumer(ctx, queue, topic); err != nil {
            s.LogError("启动消费者失败", err, "topic", topic)
            return err
        }
    }
    
    s.LogInfo("自定义消费者启动成功", "count", len(customTopics))
    return nil
}

// 启动单个消费者
func (s *CustomWorkerService) startConsumer(ctx context.Context, queue *redis.RedisQueue, topic string) error {
    consumerCtx, cancel := context.WithCancel(ctx)
    
    s.mutex.Lock()
    s.consumers[topic] = cancel
    s.mutex.Unlock()
    
    go func() {
        defer cancel()
        
        s.LogInfo("启动消费者", "topic", topic)
        
        if err := queue.Subscribe(consumerCtx, topic, s.getMessageHandler(topic)); err != nil {
            s.LogError("消费者运行失败", err, "topic", topic)
        }
        
        s.LogInfo("消费者停止", "topic", topic)
    }()
    
    return nil
}

// 获取消息处理器
func (s *CustomWorkerService) getMessageHandler(topic string) redis.MessageHandler {
    handlers := map[string]redis.MessageHandler{
        "email.send":         s.handleEmailSend,
        "sms.send":           s.handleSMSSend,
        "push.notification":  s.handlePushNotification,
        "file.process":       s.handleFileProcess,
        "data.sync":          s.handleDataSync,
    }
    
    handler, exists := handlers[topic]
    if !exists {
        return s.handleDefault
    }
    
    return handler
}

// 处理邮件发送
func (s *CustomWorkerService) handleEmailSend(ctx context.Context, msg *redis.Message) error {
    to := msg.Data["to"].(string)
    subject := msg.Data["subject"].(string)
    content := msg.Data["content"].(string)
    
    s.LogInfo("处理邮件发送", "to", to, "subject", subject)
    
    // 发送邮件逻辑
    return s.sendEmail(ctx, to, subject, content)
}

// 处理短信发送
func (s *CustomWorkerService) handleSMSSend(ctx context.Context, msg *redis.Message) error {
    phone := msg.Data["phone"].(string)
    content := msg.Data["content"].(string)
    
    s.LogInfo("处理短信发送", "phone", phone)
    
    // 发送短信逻辑
    return s.sendSMS(ctx, phone, content)
}

// 处理推送通知
func (s *CustomWorkerService) handlePushNotification(ctx context.Context, msg *redis.Message) error {
    userID := msg.Data["user_id"].(string)
    title := msg.Data["title"].(string)
    body := msg.Data["body"].(string)
    
    s.LogInfo("处理推送通知", "user_id", userID, "title", title)
    
    // 发送推送通知逻辑
    return s.sendPushNotification(ctx, userID, title, body)
}

// 处理文件处理
func (s *CustomWorkerService) handleFileProcess(ctx context.Context, msg *redis.Message) error {
    fileID := msg.Data["file_id"].(string)
    processType := msg.Data["process_type"].(string)
    
    s.LogInfo("处理文件", "file_id", fileID, "type", processType)
    
    // 文件处理逻辑
    return s.processFile(ctx, fileID, processType)
}

// 处理数据同步
func (s *CustomWorkerService) handleDataSync(ctx context.Context, msg *redis.Message) error {
    syncType := msg.Data["sync_type"].(string)
    source := msg.Data["source"].(string)
    target := msg.Data["target"].(string)
    
    s.LogInfo("处理数据同步", "type", syncType, "source", source, "target", target)
    
    // 数据同步逻辑
    return s.syncData(ctx, syncType, source, target)
}

// 默认处理器
func (s *CustomWorkerService) handleDefault(ctx context.Context, msg *redis.Message) error {
    s.LogWarn("未处理的消息", "topic", msg.Topic, "message_id", msg.ID)
    return nil
}

// 停止所有消费者
func (s *CustomWorkerService) StopConsumers() {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    for topic, cancel := range s.consumers {
        cancel()
        s.LogInfo("消费者停止", "topic", topic)
    }
    
    s.consumers = make(map[string]context.CancelFunc)
}

// 业务方法实现（示例）

func (s *CustomWorkerService) sendEmail(ctx context.Context, to, subject, content string) error {
    s.LogInfo("发送邮件", "to", to, "subject", subject)
    time.Sleep(100 * time.Millisecond) // 模拟发送时间
    return nil
}

func (s *CustomWorkerService) sendSMS(ctx context.Context, phone, content string) error {
    s.LogInfo("发送短信", "phone", phone)
    time.Sleep(50 * time.Millisecond) // 模拟发送时间
    return nil
}

func (s *CustomWorkerService) sendPushNotification(ctx context.Context, userID, title, body string) error {
    s.LogInfo("发送推送通知", "user_id", userID, "title", title)
    time.Sleep(30 * time.Millisecond) // 模拟发送时间
    return nil
}

func (s *CustomWorkerService) processFile(ctx context.Context, fileID, processType string) error {
    s.LogInfo("处理文件", "file_id", fileID, "type", processType)
    time.Sleep(500 * time.Millisecond) // 模拟处理时间
    return nil
}

func (s *CustomWorkerService) syncData(ctx context.Context, syncType, source, target string) error {
    s.LogInfo("同步数据", "type", syncType, "source", source, "target", target)
    time.Sleep(1 * time.Second) // 模拟同步时间
    return nil
}
```

## 5. 监控和运维

### 健康检查接口

```bash
# 基础健康检查
curl http://localhost:8084/healthz
# 返回: {"status":"ok","service":"gateway-worker","timestamp":{}}

# 就绪检查
curl http://localhost:8084/ready
# 返回: {"status":"ready","service":"gateway-worker","consumers":"running"}

# 指标信息
curl http://localhost:8084/metrics
# 返回: {"service":"gateway-worker","status":"running","uptime":"running","consumers":{"active":5,"topics":["order.reminder","user.notification","system.cleanup","payment.retry","inventory.alert"]}}
```

### 日志监控

```go
// 在服务中添加监控日志
func (s *CustomWorkerService) LogConsumerStatus() {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    s.LogInfo("消费者状态", "active_count", len(s.consumers))
    
    for topic := range s.consumers {
        s.LogInfo("活跃消费者", "topic", topic)
    }
}
```

### 性能监控

```go
// 添加性能监控
func (s *CustomWorkerService) handleEmailSendWithMetrics(ctx context.Context, msg *redis.Message) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        s.LogInfo("邮件发送完成", "duration", duration.String())
    }()
    
    // 处理邮件发送
    return s.handleEmailSend(ctx, msg)
}
```

## 6. 部署配置

### Docker 部署

```dockerfile
# Dockerfile for gateway-worker
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o gateway-worker ./services/gateway-worker/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway-worker .
CMD ["./gateway-worker"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes

  gateway-worker:
    build: .
    ports:
      - "8084:8084"
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - APP_PORT=8084
    depends_on:
      - redis
    restart: unless-stopped
```

### Kubernetes 部署

```yaml
# gateway-worker-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-worker
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway-worker
  template:
    metadata:
      labels:
        app: gateway-worker
    spec:
      containers:
      - name: gateway-worker
        image: ginforge/gateway-worker:latest
        ports:
        - containerPort: 8084
        env:
        - name: REDIS_HOST
          value: "redis-service"
        - name: REDIS_PORT
          value: "6379"
        - name: APP_PORT
          value: "8084"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8084
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8084
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 7. 最佳实践

### 1. 消息幂等性

```go
// 确保消息处理的幂等性
func (s *CustomWorkerService) handleEmailSendIdempotent(ctx context.Context, msg *redis.Message) error {
    messageID := msg.ID
    to := msg.Data["to"].(string)
    
    // 检查是否已经处理过
    if s.isMessageProcessed(messageID) {
        s.LogInfo("消息已处理，跳过", "message_id", messageID)
        return nil
    }
    
    // 处理消息
    err := s.sendEmail(ctx, to, msg.Data["subject"].(string), msg.Data["content"].(string))
    
    if err == nil {
        // 标记消息已处理
        s.markMessageProcessed(messageID)
    }
    
    return err
}
```

### 2. 错误处理和重试

```go
// 带重试的消息处理
func (s *CustomWorkerService) handleWithRetry(ctx context.Context, msg *redis.Message, handler redis.MessageHandler) error {
    maxRetries := 3
    retryDelay := time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := handler(ctx, msg)
        if err == nil {
            return nil
        }
        
        s.LogWarn("消息处理失败，重试中", "retry", i+1, "error", err)
        
        if i < maxRetries-1 {
            time.Sleep(retryDelay * time.Duration(i+1))
        }
    }
    
    s.LogError("消息处理失败，已达最大重试次数", "message_id", msg.ID)
    return fmt.Errorf("消息处理失败，已重试 %d 次", maxRetries)
}
```

### 3. 批量处理

```go
// 批量处理消息
func (s *CustomWorkerService) handleBatchEmail(ctx context.Context, messages []redis.Message) error {
    if len(messages) == 0 {
        return nil
    }
    
    s.LogInfo("批量处理邮件", "count", len(messages))
    
    // 批量发送邮件
    for _, msg := range messages {
        go func(m redis.Message) {
            if err := s.handleEmailSend(ctx, &m); err != nil {
                s.LogError("批量邮件发送失败", err, "message_id", m.ID)
            }
        }(msg)
    }
    
    return nil
}
```

## 8. 故障排查

### 常见问题

1. **消费者无法启动**
   - 检查 Redis 连接
   - 检查配置文件
   - 查看日志错误信息

2. **消息处理失败**
   - 检查消息格式
   - 检查处理器逻辑
   - 查看重试机制

3. **性能问题**
   - 调整消费者数量
   - 优化处理逻辑
   - 监控资源使用

### 调试命令

```bash
# 查看 Redis 队列状态
redis-cli
> LLEN mq:order.reminder
> LLEN mq:user.notification

# 查看服务日志
docker logs gateway-worker

# 检查服务状态
curl http://localhost:8084/metrics
```

## 9. 完整示例

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
    "goweb/services/gateway-worker/internal/service"
)

func main() {
    cfg := config.New()
    log := logger.New("gateway-worker-demo", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisConfig := cfg.GetRedisConfig()
    redisManager := redis.NewManager(&redisConfig, log)
    
    // 创建自定义 Worker 服务
    workerService := service.NewCustomWorkerService(redisManager, log)
    
    ctx := context.Background()
    
    // 启动消费者
    if err := workerService.StartCustomConsumers(ctx); err != nil {
        log.Fatal("启动消费者失败", err)
    }
    
    log.Info("Gateway Worker 启动成功")
    
    // 保持程序运行
    select {}
}
```

Gateway Worker 提供了完整的消息队列处理解决方案，支持各种业务场景的异步任务处理！🚀
