# Redis 统一包使用示例

## 概述

新的 `pkg/redis` 包整合了所有 Redis 相关功能：
- **缓存**：键值存储、TTL、批量操作
- **消息队列**：发布订阅、延迟消息、死信队列
- **分布式锁**：互斥锁、超时、自动续期

## 1. 基础使用

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
    log := logger.New("redis-service", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisManager := redis.NewManager(cfg.GetRedisConfig(), log)
    
    // 检查 Redis 是否启用
    if !redisManager.IsEnabled() {
        log.Warn("Redis is not enabled")
        return
    }
    
    // 测试连接
    ctx := context.Background()
    if err := redisManager.Ping(ctx); err != nil {
        log.Error("Redis ping failed", err)
        return
    }
    
    log.Info("Redis connected successfully")
}
```

## 2. 缓存使用

```go
package service

import (
    "context"
    "time"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type UserService struct {
    *base.BaseService
    cache *redis.Cache
}

func NewUserService(redisManager *redis.Manager, log logger.Logger) *UserService {
    return &UserService{
        BaseService: base.NewBaseService(log),
        cache:       redisManager.GetCache(),
    }
}

// 获取用户信息（带缓存）
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    cacheKey := "user:" + userID
    var user User
    
    // 尝试从缓存获取
    if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
        s.LogInfo("user found in cache", "user_id", userID)
        return &user, nil
    }
    
    // 缓存未命中，从数据库获取
    user, err := s.getUserFromDB(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存，TTL 5 分钟
    if err := s.cache.Set(ctx, cacheKey, user, 5*time.Minute); err != nil {
        s.LogWarn("failed to cache user", "user_id", userID, "error", err)
    }
    
    return &user, nil
}

// 更新用户信息（清除缓存）
func (s *UserService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
    // 更新数据库
    if err := s.updateUserInDB(ctx, userID, updates); err != nil {
        return err
    }
    
    // 清除缓存
    cacheKey := "user:" + userID
    if err := s.cache.Delete(ctx, cacheKey); err != nil {
        s.LogWarn("failed to clear user cache", "user_id", userID, "error", err)
    }
    
    return nil
}

// 批量获取用户
func (s *UserService) GetUsers(ctx context.Context, userIDs []string) ([]*User, error) {
    users := make([]*User, 0, len(userIDs))
    missingIDs := make([]string, 0)
    
    // 尝试从缓存批量获取
    for _, userID := range userIDs {
        var user User
        if err := s.cache.Get(ctx, "user:"+userID, &user); err == nil {
            users = append(users, &user)
        } else {
            missingIDs = append(missingIDs, userID)
        }
    }
    
    // 从数据库获取缺失的用户
    if len(missingIDs) > 0 {
        dbUsers, err := s.getUsersFromDB(ctx, missingIDs)
        if err != nil {
            return nil, err
        }
        
        // 存入缓存
        for _, user := range dbUsers {
            s.cache.Set(ctx, "user:"+user.ID, user, 5*time.Minute)
            users = append(users, user)
        }
    }
    
    return users, nil
}
```

## 3. 消息队列使用

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
    // 保存订单到数据库
    orderID, err := s.saveOrderToDB(ctx, orderData)
    if err != nil {
        return err
    }
    
    // 发布订单创建事件
    eventData := map[string]interface{}{
        "order_id": orderID,
        "user_id":  orderData["user_id"],
        "amount":   orderData["amount"],
        "items":    orderData["items"],
        "created_at": time.Now(),
    }
    
    if err := s.queue.Publish(ctx, "order.created", eventData); err != nil {
        s.LogError("failed to publish order created event", err)
        return err
    }
    
    s.LogInfo("order created and event published", "order_id", orderID)
    return nil
}

// 处理订单支付
func (s *OrderService) ProcessPayment(ctx context.Context, orderID string, paymentData map[string]interface{}) error {
    // 使用分布式锁确保订单支付的原子性
    lockKey := "order:payment:" + orderID
    return s.queue.WithLock(ctx, lockKey, 30*time.Second, func() error {
        // 检查订单状态
        order, err := s.getOrderFromDB(ctx, orderID)
        if err != nil {
            return err
        }
        
        if order.Status != "pending" {
            return fmt.Errorf("order is not in pending status")
        }
        
        // 处理支付
        if err := s.processPaymentInDB(ctx, orderID, paymentData); err != nil {
            return err
        }
        
        // 发布支付成功事件
        eventData := map[string]interface{}{
            "order_id": orderID,
            "payment_id": paymentData["payment_id"],
            "amount": order.Amount,
            "paid_at": time.Now(),
        }
        
        return s.queue.Publish(ctx, "order.paid", eventData)
    })
}

// 发送订单提醒（延迟消息）
func (s *OrderService) SendOrderReminder(ctx context.Context, orderID string, userID string) error {
    reminderData := map[string]interface{}{
        "order_id": orderID,
        "user_id":  userID,
        "type":     "payment_reminder",
    }
    
    // 24小时后发送提醒
    return s.queue.PublishWithDelay(ctx, "order.reminder", reminderData, 24*time.Hour)
}

// 启动消息消费者
func (s *OrderService) StartConsumers(ctx context.Context) error {
    // 启动订单支付消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.paid", s.handleOrderPaid); err != nil {
            s.LogError("order paid consumer failed", err)
        }
    }()
    
    // 启动订单取消消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.cancelled", s.handleOrderCancelled); err != nil {
            s.LogError("order cancelled consumer failed", err)
        }
    }()
    
    // 启动订单提醒消费者
    go func() {
        if err := s.queue.Subscribe(ctx, "order.reminder", s.handleOrderReminder); err != nil {
            s.LogError("order reminder consumer failed", err)
        }
    }()
    
    return nil
}

// 处理订单支付事件
func (s *OrderService) handleOrderPaid(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    paymentID := msg.Data["payment_id"].(string)
    
    s.LogInfo("processing order payment", "order_id", orderID, "payment_id", paymentID)
    
    // 更新订单状态
    if err := s.updateOrderStatus(ctx, orderID, "paid"); err != nil {
        return err
    }
    
    // 发送确认邮件
    if err := s.sendPaymentConfirmation(ctx, orderID, paymentID); err != nil {
        s.LogWarn("failed to send payment confirmation", "order_id", orderID, "error", err)
    }
    
    // 更新库存
    if err := s.updateInventory(ctx, orderID); err != nil {
        s.LogWarn("failed to update inventory", "order_id", orderID, "error", err)
    }
    
    return nil
}

// 处理订单取消事件
func (s *OrderService) handleOrderCancelled(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    
    s.LogInfo("processing order cancellation", "order_id", orderID)
    
    // 恢复库存
    if err := s.restoreInventory(ctx, orderID); err != nil {
        s.LogWarn("failed to restore inventory", "order_id", orderID, "error", err)
    }
    
    // 处理退款
    if err := s.processRefund(ctx, orderID); err != nil {
        s.LogWarn("failed to process refund", "order_id", orderID, "error", err)
    }
    
    return nil
}

// 处理订单提醒事件
func (s *OrderService) handleOrderReminder(ctx context.Context, msg *redis.Message) error {
    orderID := msg.Data["order_id"].(string)
    userID := msg.Data["user_id"].(string)
    reminderType := msg.Data["type"].(string)
    
    s.LogInfo("sending order reminder", "order_id", orderID, "user_id", userID, "type", reminderType)
    
    // 发送提醒通知
    return s.sendOrderReminderNotification(ctx, orderID, userID, reminderType)
}
```

## 4. 分布式锁使用

```go
package service

import (
    "context"
    "time"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type InventoryService struct {
    *base.BaseService
    redisManager *redis.Manager
}

func NewInventoryService(redisManager *redis.Manager, log logger.Logger) *InventoryService {
    return &InventoryService{
        BaseService:  base.NewBaseService(log),
        redisManager: redisManager,
    }
}

// 扣减库存（使用分布式锁）
func (s *InventoryService) DeductStock(ctx context.Context, productID string, quantity int) error {
    lockKey := "inventory:" + productID
    
    return s.redisManager.WithLock(ctx, lockKey, 10*time.Second, func() error {
        // 获取当前库存
        currentStock, err := s.getCurrentStock(ctx, productID)
        if err != nil {
            return err
        }
        
        // 检查库存是否足够
        if currentStock < quantity {
            return fmt.Errorf("insufficient stock: current=%d, required=%d", currentStock, quantity)
        }
        
        // 扣减库存
        return s.updateStock(ctx, productID, currentStock-quantity)
    })
}

// 批量扣减库存
func (s *InventoryService) DeductStockBatch(ctx context.Context, items map[string]int) error {
    // 按产品ID排序，避免死锁
    productIDs := make([]string, 0, len(items))
    for productID := range items {
        productIDs = append(productIDs, productID)
    }
    sort.Strings(productIDs)
    
    // 依次扣减每个产品的库存
    for _, productID := range productIDs {
        quantity := items[productID]
        if err := s.DeductStock(ctx, productID, quantity); err != nil {
            // 如果某个产品扣减失败，需要回滚已扣减的库存
            s.rollbackStock(ctx, productIDs[:len(productIDs)-1], items)
            return err
        }
    }
    
    return nil
}

// 尝试获取锁（带超时）
func (s *InventoryService) TryUpdatePrice(ctx context.Context, productID string, newPrice float64) error {
    lockKey := "price:" + productID
    lock, acquired, err := s.redisManager.TryLock(ctx, lockKey, 30*time.Second, 5*time.Second)
    if err != nil {
        return err
    }
    
    if !acquired {
        return fmt.Errorf("failed to acquire lock for price update")
    }
    
    defer func() {
        if releaseErr := lock.Release(ctx); releaseErr != nil {
            s.LogError("failed to release price lock", releaseErr)
        }
    }()
    
    // 更新价格
    return s.updateProductPrice(ctx, productID, newPrice)
}
```

## 5. 完整服务示例

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.New()
    log := logger.New("redis-service", cfg.GetString("log.level"))
    
    // 创建 Redis 管理器
    redisManager := redis.NewManager(cfg.GetRedisConfig(), log)
    
    // 检查 Redis 连接
    ctx := context.Background()
    if err := redisManager.Ping(ctx); err != nil {
        log.Fatal("Redis connection failed", err)
    }
    
    // 创建服务
    userService := NewUserService(redisManager, log)
    orderService := NewOrderService(redisManager, log)
    inventoryService := NewInventoryService(redisManager, log)
    
    // 启动消息消费者
    if err := orderService.StartConsumers(ctx); err != nil {
        log.Error("Failed to start message consumers", err)
    }
    
    // 创建路由
    r := gin.New()
    
    // 注册路由
    api := r.Group("/api/v1")
    {
        api.GET("/users/:id", userService.GetUserHandler)
        api.PUT("/users/:id", userService.UpdateUserHandler)
        api.POST("/orders", orderService.CreateOrderHandler)
        api.POST("/orders/:id/pay", orderService.ProcessPaymentHandler)
        api.POST("/inventory/deduct", inventoryService.DeductStockHandler)
    }
    
    // 启动服务
    r.Run(":8080")
}
```

## 6. 配置示例

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

## 7. 性能优化建议

1. **连接池配置**：根据并发量调整 `pool_size` 和 `min_idle_conns`
2. **超时设置**：合理设置各种超时时间，避免长时间阻塞
3. **批量操作**：使用 Pipeline 进行批量 Redis 操作
4. **缓存策略**：合理设置 TTL，避免缓存雪崩
5. **锁粒度**：尽量减小锁的粒度，避免死锁

## 8. 监控和调试

```go
// 获取队列长度
length, err := redisManager.GetQueueLength(ctx, "order.created")
if err != nil {
    log.Error("Failed to get queue length", err)
} else {
    log.Info("Queue length", "topic", "order.created", "length", length)
}

// 清空队列
if err := redisManager.PurgeQueue(ctx, "order.created"); err != nil {
    log.Error("Failed to purge queue", err)
}

// 检查缓存状态
exists, err := redisManager.Exists(ctx, "user:123")
if err != nil {
    log.Error("Failed to check cache", err)
} else {
    log.Info("Cache exists", "key", "user:123", "exists", exists)
}
```

这个统一的 Redis 包让所有 Redis 相关功能更加集中和易于管理！🚀
