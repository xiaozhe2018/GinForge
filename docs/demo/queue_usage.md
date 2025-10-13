# 消息队列使用示例

## 1. 基础使用

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
    "github.com/redis/go-redis/v9"
)

func main() {
    cfg := config.New()
    log := logger.New("mq-service", cfg.GetString("log.level"))
    
    // 创建 Redis 客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // 创建 Redis 消息队列
    redisQueue := queue.NewRedisQueue(redisClient, cfg.GetRedisConfig(), log)
    
    // 创建消息队列管理器
    mq := queue.NewManager(redisQueue)
    
    // 发布消息
    err := mq.Publish(context.Background(), "user.created", map[string]interface{}{
        "user_id": "123",
        "username": "john",
        "email": "john@example.com",
    })
    if err != nil {
        log.Error("发布消息失败", err)
    }
    
    // 延迟发布消息
    err = mq.PublishWithDelay(context.Background(), "order.reminder", map[string]interface{}{
        "order_id": "456",
        "user_id": "123",
    }, 24*time.Hour) // 24小时后发送
    if err != nil {
        log.Error("延迟发布消息失败", err)
    }
}
```

## 2. 消息订阅

```go
package service

import (
    "context"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type UserService struct {
    *base.BaseService
    mq *redis.RedisQueue
}

func NewUserService(mq *queue.Manager, log logger.Logger) *UserService {
    return &UserService{
        BaseService: base.NewBaseService(log),
        mq:          mq,
    }
}

// 处理用户创建事件
func (s *UserService) HandleUserCreated(ctx context.Context, msg *queue.Message) error {
    userID := msg.Data["user_id"].(string)
    username := msg.Data["username"].(string)
    
    s.LogInfo("处理用户创建事件", "user_id", userID, "username", username)
    
    // 发送欢迎邮件
    // 创建用户配置
    // 初始化用户数据
    // ...
    
    return nil
}

// 启动消息消费者
func (s *UserService) StartConsumer(ctx context.Context) error {
    return s.mq.Subscribe(ctx, "user.created", s.HandleUserCreated)
}
```

## 3. 订单服务示例

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
    mq *redis.RedisQueue
}

func NewOrderService(mq *redis.RedisQueue, log logger.Logger) *OrderService {
    return &OrderService{
        BaseService: base.NewBaseService(log),
        mq:          mq,
    }
}

// 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, orderData map[string]interface{}) error {
    // 发布订单创建事件
    err := s.mq.Publish(ctx, "order.created", map[string]interface{}{
        "order_id": orderData["id"],
        "user_id":  orderData["user_id"],
        "amount":   orderData["amount"],
        "items":    orderData["items"],
    })
    
    if err != nil {
        s.LogError("发布订单创建事件失败", err)
        return err
    }
    
    s.LogInfo("订单创建事件已发布", "order_id", orderData["id"])
    return nil
}

// 处理订单支付事件
func (s *OrderService) HandleOrderPaid(ctx context.Context, msg *queue.Message) error {
    orderID := msg.Data["order_id"].(string)
    
    s.LogInfo("处理订单支付事件", "order_id", orderID)
    
    // 更新订单状态
    // 发送确认邮件
    // 更新库存
    // ...
    
    return nil
}

// 处理订单取消事件
func (s *OrderService) HandleOrderCancelled(ctx context.Context, msg *queue.Message) error {
    orderID := msg.Data["order_id"].(string)
    
    s.LogInfo("处理订单取消事件", "order_id", orderID)
    
    // 恢复库存
    // 处理退款
    // 发送通知
    // ...
    
    return nil
}

// 启动所有消息消费者
func (s *OrderService) StartConsumers(ctx context.Context) error {
    // 启动订单支付消费者
    go func() {
        if err := s.mq.Subscribe(ctx, "order.paid", s.HandleOrderPaid); err != nil {
            s.LogError("订单支付消费者启动失败", err)
        }
    }()
    
    // 启动订单取消消费者
    go func() {
        if err := s.mq.Subscribe(ctx, "order.cancelled", s.HandleOrderCancelled); err != nil {
            s.LogError("订单取消消费者启动失败", err)
        }
    }()
    
    return nil
}
```

## 4. 在 Handler 中使用

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/redis"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type OrderHandler struct {
    *base.BaseHandler
    mq *redis.RedisQueue
}

func NewOrderHandler(mq *queue.Manager, log logger.Logger) *OrderHandler {
    return &OrderHandler{
        BaseHandler: base.NewBaseHandler(log),
        mq:          mq,
    }
}

// @Summary 创建订单
// @Description 创建新订单
// @Tags order
// @Accept json
// @Produce json
// @Param order body map[string]interface{} true "订单数据"
// @Success 200 {object} response.Response{data=object}
// @Router /order [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var orderData map[string]interface{}
    if err := c.ShouldBindJSON(&orderData); err != nil {
        h.BadRequest(c, "参数错误")
        return
    }
    
    // 发布订单创建事件
    err := h.mq.Publish(c.Request.Context(), "order.created", orderData)
    if err != nil {
        h.LogError("发布订单事件失败", err)
        h.InternalError(c, "创建订单失败")
        return
    }
    
    h.Success(c, gin.H{
        "message": "订单创建成功",
        "order_id": orderData["id"],
    })
}
```

## 5. 配置示例

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
  max_conn_age: "1h"
  pool_timeout: "4s"
  idle_timeout: "5m"
  idle_check_frequency: "1m"
```

## 6. 完整服务示例

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

func main() {
    cfg := config.New()
    log := logger.New("mq-service", cfg.GetString("log.level"))
    
    // 创建 Redis 客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr:     cfg.GetString("redis.host") + ":" + cfg.GetString("redis.port"),
        Password: cfg.GetString("redis.password"),
        DB:       cfg.GetInt("redis.database"),
    })
    
    // 创建消息队列
    redisManager := redis.NewManager(cfg.GetRedisConfig(), log)
    mq := redisManager.GetQueue()
    
    // 创建服务
    orderService := NewOrderService(mq, log)
    
    // 启动消息消费者
    go func() {
        if err := orderService.StartConsumers(context.Background()); err != nil {
            log.Error("消息消费者启动失败", err)
        }
    }()
    
    // 创建路由
    r := gin.New()
    
    // 创建处理器
    orderHandler := NewOrderHandler(mq, log)
    
    // 注册路由
    api := r.Group("/api/v1")
    {
        api.POST("/orders", orderHandler.CreateOrder)
    }
    
    // 启动服务
    r.Run(":8080")
}
```

## 7. 消息队列管理

```go
// 获取队列长度
length, err := mq.GetQueueLength(ctx, "user.created")
if err != nil {
    log.Error("获取队列长度失败", err)
} else {
    log.Info("队列长度", "topic", "user.created", "length", length)
}

// 清空队列
err = mq.PurgeQueue(ctx, "user.created")
if err != nil {
    log.Error("清空队列失败", err)
}
```

## 8. 错误处理

消息队列支持自动重试机制：

- 消息处理失败时，会自动重试
- 达到最大重试次数后，消息会进入死信队列
- 死信队列的键格式：`mq:dead-letter:{topic}`

## 9. 延迟消息

```go
// 发送延迟消息（24小时后发送）
err := mq.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
    "order_id": "123",
    "user_id": "456",
}, 24*time.Hour)
```

## 10. 消息格式

```go
type Message struct {
    ID        string                 `json:"id"`        // 消息ID
    Topic     string                 `json:"topic"`     // 主题
    Data      map[string]interface{} `json:"data"`      // 消息数据
    Timestamp time.Time              `json:"timestamp"` // 时间戳
    Retry     int                    `json:"retry"`     // 重试次数
    MaxRetry  int                    `json:"max_retry"` // 最大重试次数
}
```

这个简化的消息队列包结构更加清晰，遵循了 cache 包的设计模式，使用起来更加简单直观。
