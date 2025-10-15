# 消息队列详解

深入了解 GinForge 的消息队列系统，掌握异步任务处理。

## 🎯 为什么需要消息队列？

### 解决的问题

1. **解耦服务** - 服务之间不直接调用
2. **异步处理** - 耗时任务异步执行
3. **削峰填谷** - 应对流量高峰
4. **可靠性** - 消息持久化，不丢失
5. **可扩展** - 水平扩展消费者

### 适用场景

- ✅ 发送邮件/短信
- ✅ 数据同步
- ✅ 日志处理
- ✅ 文件处理
- ✅ 订单处理
- ✅ 通知推送

## 🏗️ 架构设计

```
┌─────────────┐     发布消息     ┌─────────────┐
│  Admin API  │ ───────────────> │    Redis    │
│  (生产者)    │                  │   消息队列   │
└─────────────┘                  └──────┬──────┘
                                        │
                                 消费消息 │
                                        ↓
                              ┌──────────────────┐
                              │ Gateway Worker   │
                              │    (消费者)       │
                              └──────────────────┘
                                        │
                                        ↓
                              执行异步任务（发邮件、处理文件等）
```

## 📝 基础使用

### 1. 初始化消息队列

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
)

func main() {
    cfg := config.New()
    log := logger.New("mq-demo", cfg.GetString("log.level"))
    
    // 创建 Redis 客户端
    redisConfig := cfg.GetRedisConfig()
    redisClient := redis.NewClient(&redisConfig, log)
    
    // 获取消息队列管理器
    queue := redisClient.GetQueue()
    
    log.Info("消息队列初始化成功")
}
```

### 2. 发布消息

```go
ctx := context.Background()

// 发布普通消息
err := queue.Publish(ctx, "user.welcome", map[string]interface{}{
    "user_id": "123",
    "username": "john",
    "email": "john@example.com",
})

if err != nil {
    log.Error("发布消息失败", err)
}
```

### 3. 消费消息

```go
// 定义消息处理函数
func handleUserWelcome(ctx context.Context, data map[string]interface{}) error {
    userID := data["user_id"].(string)
    email := data["email"].(string)
    
    log.Info("处理欢迎消息", "user_id", userID)
    
    // 发送欢迎邮件
    return sendWelcomeEmail(email)
}

// 启动消费者
go func() {
    err := queue.Subscribe(ctx, "user.welcome", handleUserWelcome)
    if err != nil {
        log.Error("消费者启动失败", err)
    }
}()
```

## ⏰ 延时消息

### 使用延时队列

```go
// 发送延时消息（24小时后执行）
err := queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
    "order_id": "ORD123456",
    "user_id": "789",
    "message": "您的订单尚未支付，请及时完成支付",
}, 24*time.Hour)
```

### 延时任务示例

```go
// 订单超时自动取消
func (s *OrderService) CreateOrder(orderData map[string]interface{}) error {
    // 1. 创建订单
    order := createOrder(orderData)
    
    // 2. 发送延时消息（30分钟后检查支付状态）
    s.queue.PublishWithDelay(ctx, "order.payment.check", map[string]interface{}{
        "order_id": order.ID,
    }, 30*time.Minute)
    
    return nil
}

// 处理订单支付检查
func (s *OrderService) HandlePaymentCheck(ctx context.Context, data map[string]interface{}) error {
    orderID := data["order_id"].(string)
    
    // 查询订单状态
    order, err := s.orderRepo.GetByID(orderID)
    if err != nil {
        return err
    }
    
    // 如果still未支付，自动取消
    if order.PaymentStatus == "unpaid" {
        s.logger.Info("订单超时未支付，自动取消", "order_id", orderID)
        return s.CancelOrder(orderID, "超时未支付")
    }
    
    return nil
}
```

## 🔄 完整的消息处理流程

### 场景：用户注册流程

```go
// 1. 用户注册 Handler
func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }
    
    // 创建用户
    user, err := h.userService.CreateUser(&req)
    if err != nil {
        response.InternalError(c, "注册失败")
        return
    }
    
    // 发布用户注册事件
    h.queue.Publish(c.Request.Context(), "user.registered", map[string]interface{}{
        "user_id":  user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
    
    response.Success(c, user)
}

// 2. Gateway Worker 消费消息
func (w *Worker) handleUserRegistered(ctx context.Context, data map[string]interface{}) error {
    userID := data["user_id"].(string)
    email := data["email"].(string)
    
    // 任务 1：发送欢迎邮件
    w.sendWelcomeEmail(email)
    
    // 任务 2：创建用户配置
    w.createUserPreferences(userID)
    
    // 任务 3：添加到新用户群组
    w.addToNewUserGroup(userID)
    
    // 任务 4：发送欢迎礼包（延时1小时）
    w.queue.PublishWithDelay(ctx, "user.welcome.gift", map[string]interface{}{
        "user_id": userID,
    }, 1*time.Hour)
    
    w.logger.Info("用户注册事件处理完成", "user_id", userID)
    return nil
}
```

## 🎨 消息类型定义

### 推荐的消息命名

```
格式：<实体>.<动作>

示例：
user.created        # 用户创建
user.updated        # 用户更新
user.deleted        # 用户删除
order.paid          # 订单支付
order.shipped       # 订单发货
order.cancelled     # 订单取消
email.send          # 发送邮件
sms.send            # 发送短信
file.uploaded       # 文件上传
file.processed      # 文件处理完成
```

### 消息数据结构

```go
type Message struct {
    ID        string                 `json:"id"`
    Topic     string                 `json:"topic"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    Retry     int                    `json:"retry"`
}

// 示例消息
{
    "id": "msg_123456",
    "topic": "order.paid",
    "data": {
        "order_id": "ORD123",
        "user_id": "789",
        "amount": 299.99,
        "payment_method": "alipay"
    },
    "timestamp": "2025-10-15T10:30:00Z",
    "retry": 0
}
```

## 🛡️ 可靠性保证

### 1. 消息确认机制

```go
func (w *Worker) handleMessage(ctx context.Context, data map[string]interface{}) error {
    // 处理消息
    if err := w.processMessage(data); err != nil {
        // 返回错误，消息会重新入队
        return err
    }
    
    // 返回 nil，消息被确认删除
    return nil
}
```

### 2. 重试机制

```go
// 消息处理失败自动重试
func (q *RedisQueue) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
    maxRetries := 3
    
    for {
        msg, err := q.consume(ctx, topic)
        if err != nil {
            continue
        }
        
        // 处理消息
        if err := handler(ctx, msg.Data); err != nil {
            msg.Retry++
            
            if msg.Retry < maxRetries {
                // 重新入队
                q.requeueMessage(msg)
            } else {
                // 达到最大重试次数，移入死信队列
                q.moveToDeadLetter(msg)
            }
        }
    }
}
```

### 3. 死信队列

```go
// 查看死信队列
func (q *RedisQueue) GetDeadLetters(topic string) ([]Message, error) {
    key := fmt.Sprintf("mq:dead-letter:%s", topic)
    // 获取死信消息列表
    return q.getMessages(key)
}

// 重新处理死信消息
func (q *RedisQueue) ReprocessDeadLetter(messageID string) error {
    // 从死信队列移回正常队列
    return q.moveFromDeadLetter(messageID)
}
```

## 📊 监控和管理

### 查看队列状态

```go
// 获取队列长度
length, err := queue.GetQueueLength(ctx, "user.welcome")
log.Info("队列长度", "topic", "user.welcome", "length", length)

// 获取所有队列
queues, err := queue.GetAllQueues(ctx)
for _, q := range queues {
    log.Info("队列信息", "topic", q.Topic, "length", q.Length)
}
```

### Redis 命令查看

```bash
# 连接 Redis
docker exec -it redis redis-cli

# 查看所有队列
KEYS mq:*

# 查看队列长度
LLEN mq:user.welcome

# 查看队列内容（不弹出）
LRANGE mq:user.welcome 0 10

# 清空队列
DEL mq:user.welcome

# 查看死信队列
LLEN mq:dead-letter:user.welcome
```

## 💡 最佳实践

### 1. 消息幂等性

确保消息被重复消费时不会产生副作用：

```go
func (w *Worker) handleEmailSend(ctx context.Context, data map[string]interface{}) error {
    messageID := data["message_id"].(string)
    email := data["email"].(string)
    
    // 检查是否已处理
    key := fmt.Sprintf("processed:%s", messageID)
    exists, _ := w.redis.Exists(ctx, key)
    if exists {
        w.logger.Info("消息已处理，跳过", "message_id", messageID)
        return nil
    }
    
    // 发送邮件
    if err := w.sendEmail(email); err != nil {
        return err
    }
    
    // 标记为已处理（24小时过期）
    w.redis.Set(ctx, key, "1", 24*time.Hour)
    
    return nil
}
```

### 2. 批量处理

提高处理效率：

```go
func (w *Worker) handleBatchEmails(ctx context.Context) error {
    batchSize := 100
    messages := make([]Message, 0, batchSize)
    
    // 批量获取消息
    for i := 0; i < batchSize; i++ {
        msg, err := w.queue.PopMessage(ctx, "email.send")
        if err != nil {
            break
        }
        messages = append(messages, msg)
    }
    
    if len(messages) == 0 {
        return nil
    }
    
    // 批量发送
    w.logger.Info("批量处理邮件", "count", len(messages))
    return w.sendBatchEmails(messages)
}
```

### 3. 优先级队列

```go
// 高优先级消息
func (s *OrderService) SendUrgentNotification(orderID string) error {
    return s.queue.PublishToQueue(ctx, "notification.urgent", map[string]interface{}{
        "order_id": orderID,
        "priority": "high",
    })
}

// 普通优先级消息
func (s *OrderService) SendNormalNotification(orderID string) error {
    return s.queue.PublishToQueue(ctx, "notification.normal", map[string]interface{}{
        "order_id": orderID,
        "priority": "normal",
    })
}

// 消费时优先处理高优先级队列
func (w *Worker) StartConsumers() {
    go w.consumeUrgent()  // 高优先级
    go w.consumeNormal()  // 普通优先级
}
```

## 🔧 高级功能

### 延时队列实现

查看完整实现：`pkg/redis/delayed_worker.go`

```go
// 启动延时任务处理器
func (w *DelayedWorker) Start(ctx context.Context) error {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // 检查到期的消息
            w.processExpiredMessages(ctx)
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// 处理到期消息
func (w *DelayedWorker) processExpiredMessages(ctx context.Context) {
    now := time.Now().Unix()
    
    // 查询到期的消息（使用 Redis Sorted Set）
    messages, err := w.redis.ZRangeByScore(ctx, "delayed:messages", 0, float64(now))
    if err != nil {
        return
    }
    
    for _, msgData := range messages {
        var msg Message
        json.Unmarshal([]byte(msgData), &msg)
        
        // 发送到正常队列
        w.queue.Publish(ctx, msg.Topic, msg.Data)
        
        // 从延时队列移除
        w.redis.ZRem(ctx, "delayed:messages", msgData)
    }
}
```

## 📚 实战案例

### 案例 1：订单系统

```go
// 订单创建
func (s *OrderService) CreateOrder(ctx context.Context, orderData map[string]interface{}) error {
    // 1. 创建订单
    order := s.createOrder(orderData)
    
    // 2. 发布订单创建事件
    s.queue.Publish(ctx, "order.created", map[string]interface{}{
        "order_id": order.ID,
        "user_id":  order.UserID,
        "amount":   order.Amount,
    })
    
    // 3. 30分钟后检查支付状态
    s.queue.PublishWithDelay(ctx, "order.payment.check", map[string]interface{}{
        "order_id": order.ID,
    }, 30*time.Minute)
    
    // 4. 24小时后发送提醒
    s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
        "order_id": order.ID,
        "type":     "payment",
    }, 24*time.Hour)
    
    return nil
}

// 处理订单创建事件
func (w *Worker) handleOrderCreated(ctx context.Context, data map[string]interface{}) error {
    orderID := data["order_id"].(string)
    
    // 任务 1：扣减库存
    w.reduceInventory(orderID)
    
    // 任务 2：发送确认短信
    w.sendOrderConfirmSMS(data["user_id"].(string), orderID)
    
    // 任务 3：通知商家
    w.notifyMerchant(orderID)
    
    return nil
}

// 处理支付状态检查
func (w *Worker) handlePaymentCheck(ctx context.Context, data map[string]interface{}) error {
    orderID := data["order_id"].(string)
    
    order, _ := w.orderRepo.GetByID(orderID)
    
    // 如果未支付，自动取消订单
    if order.PaymentStatus == "unpaid" {
        w.logger.Info("订单超时未支付，自动取消", "order_id", orderID)
        w.cancelOrder(orderID)
        
        // 恢复库存
        w.restoreInventory(orderID)
    }
    
    return nil
}
```

### 案例 2：文件处理

```go
// 文件上传后异步处理
func (s *FileService) HandleUpload(file *multipart.FileHeader) error {
    // 1. 保存原始文件
    savedFile := s.saveFile(file)
    
    // 2. 发布文件处理消息
    s.queue.Publish(ctx, "file.process", map[string]interface{}{
        "file_id":   savedFile.ID,
        "file_path": savedFile.Path,
        "file_type": savedFile.Type,
    })
    
    return nil
}

// 处理文件（生成缩略图、转码等）
func (w *Worker) handleFileProcess(ctx context.Context, data map[string]interface{}) error {
    fileID := data["file_id"].(string)
    filePath := data["file_path"].(string)
    fileType := data["file_type"].(string)
    
    w.logger.Info("开始处理文件", "file_id", fileID, "type", fileType)
    
    switch fileType {
    case "image":
        // 生成缩略图
        w.generateThumbnail(filePath)
        // 压缩图片
        w.compressImage(filePath)
    case "video":
        // 视频转码
        w.transcodeVideo(filePath)
        // 生成预览图
        w.generateVideoThumbnail(filePath)
    }
    
    // 更新文件处理状态
    w.updateFileStatus(fileID, "processed")
    
    w.logger.Info("文件处理完成", "file_id", fileID)
    return nil
}
```

### 案例 3：数据同步

```go
// 数据变更时发布同步事件
func (s *UserService) UpdateUser(userID string, data map[string]interface{}) error {
    // 更新数据库
    if err := s.userRepo.Update(userID, data); err != nil {
        return err
    }
    
    // 发布数据同步事件
    s.queue.Publish(ctx, "user.sync", map[string]interface{}{
        "user_id": userID,
        "action":  "update",
        "data":    data,
    })
    
    return nil
}

// 处理数据同步
func (w *Worker) handleUserSync(ctx context.Context, data map[string]interface{}) error {
    userID := data["user_id"].(string)
    action := data["action"].(string)
    
    switch action {
    case "update":
        // 同步到缓存
        w.syncToCache(userID, data["data"])
        
        // 同步到搜索引擎
        w.syncToElasticsearch(userID, data["data"])
        
        // 同步到数据仓库
        w.syncToDataWarehouse(userID, data["data"])
    }
    
    return nil
}
```

## 🚀 性能优化

### 1. 并发消费

```go
// 启动多个消费者
func (w *Worker) StartMultipleConsumers(topic string, count int) {
    for i := 0; i < count; i++ {
        go func(index int) {
            w.logger.Info("启动消费者", "topic", topic, "index", index)
            w.queue.Subscribe(ctx, topic, w.handleMessage)
        }(i)
    }
}

// 使用
w.StartMultipleConsumers("email.send", 5)  // 5个并发消费者
```

### 2. 消息批处理

```go
func (w *Worker) batchConsume(ctx context.Context, topic string, batchSize int) {
    ticker := time.NewTicker(1 * time.Second)
    batch := make([]Message, 0, batchSize)
    
    for {
        select {
        case msg := <-w.msgChan:
            batch = append(batch, msg)
            
            // 达到批量大小，批量处理
            if len(batch) >= batchSize {
                w.processBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            // 定时处理未满的批次
            if len(batch) > 0 {
                w.processBatch(batch)
                batch = batch[:0]
            }
        }
    }
}
```

## 📚 完整示例

查看完整实现：

- **消息队列**: `pkg/redis/queue.go`
- **延时队列**: `pkg/redis/delayed_worker.go`
- **Gateway Worker**: `services/gateway-worker/`
- **使用示例**: `docs/demo/queue_usage.md`
- **延时队列示例**: `docs/demo/delayed_queue_usage.md`

## 🎯 下一步

- [分布式锁](./distributed-lock) - 并发控制
- [WebSocket 通知](../features/websocket) - 实时推送
- [性能优化](./performance) - 系统优化

---

**提示**: 消息队列是微服务架构的核心组件，务必掌握其使用方法！

