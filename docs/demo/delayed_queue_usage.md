# 延时队列使用示例

## 概述

延时队列是消息队列的重要功能，允许在指定时间后发送消息。GinForge 框架提供了完整的延时队列解决方案，包括消息存储、定时处理和可靠性保证。

## 1. 基础概念

### 延时队列工作原理

```
发送延时消息 → Redis ZSet 存储 → 定时扫描 → 到期发布 → 正常消费
```

1. **消息存储**：使用 Redis ZSet 存储，score 为到期时间戳
2. **定时扫描**：DelayedWorker 每秒扫描到期的消息
3. **消息发布**：到期消息自动发布到正常队列
4. **正常消费**：通过 Subscribe 方法消费消息

### 关键组件

- **DelayedWorker**：延时消息处理器，管理所有延时消息
- **RedisQueue**：消息队列实现，支持延时发布
- **Redis ZSet**：存储延时消息，按时间排序

## 2. 基础使用

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
    log := logger.New("delayed-queue-demo", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisConfig := cfg.GetRedisConfig()
    redisManager := redis.NewManager(&redisConfig, log)
    
    // 获取消息队列
    queue := redisManager.GetQueue()
    
    ctx := context.Background()
    
    // 发送延时消息
    err := queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
        "order_id": "12345",
        "user_id":  "67890",
        "type":     "payment_reminder",
    }, 24*time.Hour) // 24小时后发送
    
    if err != nil {
        log.Error("发送延时消息失败", err)
        return
    }
    
    log.Info("延时消息发送成功")
}
```

## 3. 完整示例

### 订单提醒系统

```go
package service

import (
    "context"
    "time"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type OrderService struct {
    *base.BaseService
    queue *redis.RedisQueue
}

func NewOrderService(redisManager *redis.Manager, log logger.Logger) *OrderService {
    return &OrderService{
        BaseService: base.NewBaseService(log),
        queue:       redisManager.GetQueue(),
    }
}

// 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, orderData map[string]interface{}) error {
    orderID := orderData["id"].(string)
    userID := orderData["user_id"].(string)
    
    // 保存订单到数据库
    if err := s.saveOrderToDB(ctx, orderData); err != nil {
        return err
    }
    
    // 发送订单创建事件
    if err := s.queue.Publish(ctx, "order.created", map[string]interface{}{
        "order_id": orderID,
        "user_id":  userID,
        "amount":   orderData["amount"],
        "created_at": time.Now(),
    }); err != nil {
        s.LogError("发送订单创建事件失败", err)
        return err
    }
    
    // 设置支付提醒（30分钟后）
    if err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
        "order_id": orderID,
        "user_id":  userID,
        "type":     "payment_reminder",
        "message":  "请及时完成支付",
    }, 30*time.Minute); err != nil {
        s.LogWarn("设置支付提醒失败", "order_id", orderID, "error", err)
    }
    
    // 设置订单超时提醒（24小时后）
    if err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
        "order_id": orderID,
        "user_id":  userID,
        "type":     "order_timeout",
        "message":  "订单即将超时，请尽快支付",
    }, 24*time.Hour); err != nil {
        s.LogWarn("设置订单超时提醒失败", "order_id", orderID, "error", err)
    }
    
    s.LogInfo("订单创建成功", "order_id", orderID)
    return nil
}

// 处理订单支付
func (s *OrderService) ProcessPayment(ctx context.Context, orderID string, paymentData map[string]interface{}) error {
    // 处理支付逻辑
    if err := s.processPaymentInDB(ctx, orderID, paymentData); err != nil {
        return err
    }
    
    // 取消未到期的提醒
    if err := s.cancelOrderReminders(ctx, orderID); err != nil {
        s.LogWarn("取消订单提醒失败", "order_id", orderID, "error", err)
    }
    
    // 发送支付成功事件
    if err := s.queue.Publish(ctx, "order.paid", map[string]interface{}{
        "order_id": orderID,
        "payment_id": paymentData["payment_id"],
        "paid_at": time.Now(),
    }); err != nil {
        s.LogError("发送支付成功事件失败", err)
        return err
    }
    
    s.LogInfo("订单支付成功", "order_id", orderID)
    return nil
}

// 取消订单提醒（示例实现）
func (s *OrderService) cancelOrderReminders(ctx context.Context, orderID string) error {
    // 这里可以实现取消延时消息的逻辑
    // 由于 Redis ZSet 的特性，可以通过删除特定消息来实现
    s.LogInfo("取消订单提醒", "order_id", orderID)
    return nil
}
```

### 消息消费者

```go
package service

import (
    "context"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type NotificationService struct {
    *base.BaseService
    queue *redis.RedisQueue
}

func NewNotificationService(redisManager *redis.Manager, log logger.Logger) *NotificationService {
    return &NotificationService{
        BaseService: base.NewBaseService(log),
        queue:       redisManager.GetQueue(),
    }
}

// 启动消息消费者
func (s *NotificationService) StartConsumers(ctx context.Context) error {
    // 启动订单提醒消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.reminder", s.handleOrderReminder); err != nil {
            s.LogError("订单提醒消费者启动失败", err)
        }
    }()
    
    // 启动订单创建消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.created", s.handleOrderCreated); err != nil {
            s.LogError("订单创建消费者启动失败", err)
        }
    }()
    
    // 启动订单支付消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.paid", s.handleOrderPaid); err != nil {
            s.LogError("订单支付消费者启动失败", err)
        }
    }()
    
    return nil
}

// 处理订单提醒
func (s *NotificationService) handleOrderReminder(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    userID := msg.Data["user_id"].(string)
    reminderType := msg.Data["type"].(string)
    message := msg.Data["message"].(string)
    
    s.LogInfo("处理订单提醒", "order_id", orderID, "user_id", userID, "type", reminderType)
    
    // 发送通知
    switch reminderType {
    case "payment_reminder":
        return s.sendPaymentReminder(ctx, orderID, userID, message)
    case "order_timeout":
        return s.sendOrderTimeoutReminder(ctx, orderID, userID, message)
    default:
        s.LogWarn("未知的提醒类型", "type", reminderType)
        return nil
    }
}

// 处理订单创建
func (s *NotificationService) handleOrderCreated(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    userID := msg.Data["user_id"].(string)
    
    s.LogInfo("处理订单创建通知", "order_id", orderID, "user_id", userID)
    
    // 发送订单创建确认邮件
    return s.sendOrderCreatedEmail(ctx, orderID, userID)
}

// 处理订单支付
func (s *NotificationService) handleOrderPaid(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    paymentID := msg.Data["payment_id"].(string)
    
    s.LogInfo("处理订单支付通知", "order_id", orderID, "payment_id", paymentID)
    
    // 发送支付成功邮件
    return s.sendPaymentSuccessEmail(ctx, orderID, paymentID)
}

// 发送支付提醒
func (s *NotificationService) sendPaymentReminder(ctx context.Context, orderID, userID, message string) error {
    s.LogInfo("发送支付提醒", "order_id", orderID, "user_id", userID, "message", message)
    
    // 这里可以集成邮件、短信、推送等服务
    // 例如：发送邮件、发送短信、发送推送通知等
    
    return nil
}

// 发送订单超时提醒
func (s *NotificationService) sendOrderTimeoutReminder(ctx context.Context, orderID, userID, message string) error {
    s.LogInfo("发送订单超时提醒", "order_id", orderID, "user_id", userID, "message", message)
    
    // 这里可以发送超时提醒，并可能自动取消订单
    
    return nil
}

// 发送订单创建邮件
func (s *NotificationService) sendOrderCreatedEmail(ctx context.Context, orderID, userID string) error {
    s.LogInfo("发送订单创建邮件", "order_id", orderID, "user_id", userID)
    return nil
}

// 发送支付成功邮件
func (s *NotificationService) sendPaymentSuccessEmail(ctx context.Context, orderID, paymentID string) error {
    s.LogInfo("发送支付成功邮件", "order_id", orderID, "payment_id", paymentID)
    return nil
}
```

## 4. 高级用法

### 批量延时消息

```go
// 批量设置多个提醒
func (s *OrderService) SetMultipleReminders(ctx context.Context, orderID, userID string) error {
    reminders := []struct {
        delay   time.Duration
        msgType string
        message string
    }{
        {30 * time.Minute, "payment_reminder", "请及时完成支付"},
        {2 * time.Hour, "payment_reminder", "支付提醒：订单即将超时"},
        {24 * time.Hour, "order_timeout", "订单已超时，将被自动取消"},
    }
    
    for _, reminder := range reminders {
        if err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
            "order_id": orderID,
            "user_id":  userID,
            "type":     reminder.msgType,
            "message":  reminder.message,
        }, reminder.delay); err != nil {
            s.LogError("设置提醒失败", err, "delay", reminder.delay, "type", reminder.msgType)
            return err
        }
    }
    
    return nil
}
```

### 条件延时消息

```go
// 根据订单金额设置不同的提醒策略
func (s *OrderService) SetSmartReminders(ctx context.Context, orderID, userID string, amount float64) error {
    var reminders []struct {
        delay   time.Duration
        msgType string
        message string
    }
    
    if amount >= 1000 {
        // 大额订单：更频繁的提醒
        reminders = []struct {
            delay   time.Duration
            msgType string
            message string
        }{
            {15 * time.Minute, "payment_reminder", "大额订单支付提醒"},
            {1 * time.Hour, "payment_reminder", "大额订单支付提醒"},
            {6 * time.Hour, "payment_reminder", "大额订单支付提醒"},
            {24 * time.Hour, "order_timeout", "大额订单即将超时"},
        }
    } else {
        // 小额订单：标准提醒
        reminders = []struct {
            delay   time.Duration
            msgType string
            message string
        }{
            {30 * time.Minute, "payment_reminder", "请及时完成支付"},
            {24 * time.Hour, "order_timeout", "订单即将超时"},
        }
    }
    
    for _, reminder := range reminders {
        if err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
            "order_id": orderID,
            "user_id":  userID,
            "type":     reminder.msgType,
            "message":  reminder.message,
            "amount":   amount,
        }, reminder.delay); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 5. 监控和调试

### 检查延时队列状态

```go
// 获取延时队列长度
func (s *OrderService) GetDelayedQueueStatus(ctx context.Context, topic string) (int64, error) {
    return s.queue.GetQueueLength(ctx, topic)
}

// 清空延时队列（谨慎使用）
func (s *OrderService) ClearDelayedQueue(ctx context.Context, topic string) error {
    return s.queue.PurgeQueue(ctx, topic)
}
```

### 日志监控

```go
// 在服务启动时记录延时队列状态
func (s *OrderService) LogQueueStatus(ctx context.Context) {
    topics := []string{"order.reminder", "user.notification", "system.cleanup"}
    
    for _, topic := range topics {
        length, err := s.queue.GetQueueLength(ctx, topic)
        if err != nil {
            s.LogError("获取队列长度失败", err, "topic", topic)
            continue
        }
        
        s.LogInfo("队列状态", "topic", topic, "length", length)
    }
}
```

## 6. 配置说明

### Redis 配置

```yaml
# configs/config.yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"
  idle_timeout: "5m"
  idle_check_frequency: "1m"
```

### 延时队列配置

```go
// 创建自定义间隔的延时处理器
delayedWorker := redis.NewDelayedWorker(redisClient, 2*time.Second) // 2秒扫描一次
```

## 7. 最佳实践

### 1. 消息设计

```go
// 好的消息设计
type OrderReminderMessage struct {
    OrderID     string    `json:"order_id"`
    UserID      string    `json:"user_id"`
    Type        string    `json:"type"`
    Message     string    `json:"message"`
    Amount      float64   `json:"amount"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   time.Time `json:"expires_at"`
}

// 发送结构化的消息
func (s *OrderService) SendStructuredReminder(ctx context.Context, orderID, userID string) error {
    message := OrderReminderMessage{
        OrderID:   orderID,
        UserID:    userID,
        Type:      "payment_reminder",
        Message:   "请及时完成支付",
        Amount:    100.0,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    
    return s.queue.PublishWithDelay(ctx, "order.reminder", message, 30*time.Minute)
}
```

### 2. 错误处理

```go
// 带重试的延时消息发送
func (s *OrderService) SendReminderWithRetry(ctx context.Context, orderID, userID string, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
            "order_id": orderID,
            "user_id":  userID,
            "type":     "payment_reminder",
            "retry":    i,
        }, 30*time.Minute)
        
        if err == nil {
            return nil
        }
        
        s.LogWarn("发送延时消息失败，重试中", "retry", i+1, "error", err)
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return fmt.Errorf("发送延时消息失败，已重试 %d 次", maxRetries)
}
```

### 3. 性能优化

```go
// 批量发送延时消息
func (s *OrderService) BatchSendReminders(ctx context.Context, orders []Order) error {
    for _, order := range orders {
        // 使用 goroutine 并发发送
        go func(o Order) {
            if err := s.queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
                "order_id": o.ID,
                "user_id":  o.UserID,
                "type":     "payment_reminder",
            }, 30*time.Minute); err != nil {
                s.LogError("批量发送提醒失败", err, "order_id", o.ID)
            }
        }(order)
    }
    
    return nil
}
```

## 8. 常见问题

### Q: 延时消息会丢失吗？
A: 不会。消息存储在 Redis ZSet 中，具有持久化特性。即使服务重启，延时消息也会被正确处理。

### Q: 如何取消已发送的延时消息？
A: 由于 Redis ZSet 的特性，可以通过删除特定消息来实现。建议在消息中包含唯一标识符。

### Q: 延时消息的精度如何？
A: 默认每秒扫描一次，精度为秒级。可以通过调整 `DelayedWorker` 的间隔来提高精度。

### Q: 如何处理大量延时消息？
A: 建议使用多个 `DelayedWorker` 实例，或者调整扫描间隔。同时注意 Redis 内存使用情况。

## 9. 完整示例

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
    log := logger.New("delayed-queue-demo", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisConfig := cfg.GetRedisConfig()
    redisManager := redis.NewManager(&redisConfig, log)
    
    // 创建服务
    orderService := NewOrderService(redisManager, log)
    notificationService := NewNotificationService(redisManager, log)
    
    ctx := context.Background()
    
    // 启动消息消费者
    if err := notificationService.StartConsumers(ctx); err != nil {
        log.Fatal("启动消费者失败", err)
    }
    
    // 创建订单
    orderData := map[string]interface{}{
        "id":      "12345",
        "user_id": "67890",
        "amount":  299.99,
    }
    
    if err := orderService.CreateOrder(ctx, orderData); err != nil {
        log.Error("创建订单失败", err)
        return
    }
    
    log.Info("延时队列示例运行中...")
    
    // 保持程序运行
    select {}
}
```

这个延时队列系统提供了完整的解决方案，支持各种业务场景的延时消息需求！🚀
