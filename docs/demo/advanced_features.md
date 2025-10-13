# 高级功能使用示例

## 1. 监控指标使用

```go
package main

import (
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/monitor"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.New()
    log := logger.New("my-service", cfg.GetString("log.level"))
    
    // 创建监控指标
    metrics := monitor.NewMetrics("my-service", log)
    
    // 创建健康检查器
    healthChecker := monitor.NewHealthChecker(log)
    
    // 注册健康检查
    healthChecker.Register(monitor.NewServiceHealthCheck(
        "custom-service",
        "自定义服务检查",
        func(ctx context.Context) error {
            // 自定义健康检查逻辑
            return nil
        },
    ))
    
    r := gin.New()
    
    // 添加监控中间件
    r.Use(metrics.HTTPMiddleware("my-service"))
    
    // 健康检查端点
    r.GET("/healthz", healthChecker.GinHandler())
    r.GET("/readyz", healthChecker.GinHandler())
    
    // 业务路由
    r.GET("/api/data", func(c *gin.Context) {
        // 记录业务指标
        metrics.RecordBusinessOperation("get_data", "success", "my-service")
        
        c.JSON(200, gin.H{"data": "hello"})
    })
}
```

## 2. 文件存储使用

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/storage"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type FileHandler struct {
    *base.BaseHandler
    storage storage.Storage
}

func NewFileHandler(storage storage.Storage, log logger.Logger) *FileHandler {
    return &FileHandler{
        BaseHandler: base.NewBaseHandler(log),
        storage:     storage,
    }
}

// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} response.Response{data=object}
// @Router /file/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        h.BadRequest(c, "文件上传失败")
        return
    }
    
    // 上传文件
    fileInfo, err := h.storage.UploadFile(file, "uploads")
    if err != nil {
        h.LogError("文件上传失败", err)
        h.InternalError(c, "文件上传失败")
        return
    }
    
    h.Success(c, gin.H{
        "file_info": fileInfo,
        "url":       "/files/" + fileInfo.RelativePath,
    })
}

// @Summary 获取文件
// @Description 根据路径获取文件信息
// @Tags file
// @Produce json
// @Param path path string true "文件路径"
// @Success 200 {object} response.Response{data=object}
// @Router /file/{path} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
    path := c.Param("path")
    
    fileInfo, err := h.storage.GetFile(path)
    if err != nil {
        h.NotFound(c, "文件不存在")
        return
    }
    
    h.Success(c, fileInfo)
}
```

## 3. 消息队列使用

```go
package service

import (
    "context"
    "goweb/pkg/queue"
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
    // ...
    
    return nil
}

// 启动消息消费者
func (s *OrderService) StartConsumer(ctx context.Context) error {
    return s.mq.Subscribe(ctx, "order.paid", s.HandleOrderPaid)
}
```

## 4. 分布式锁使用

```go
package service

import (
    "context"
    "time"
    "goweb/pkg/lock"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type InventoryService struct {
    *base.BaseService
    lock *lock.RedisLock
}

func NewInventoryService(lock *lock.RedisLock, log logger.Logger) *InventoryService {
    return &InventoryService{
        BaseService: base.NewBaseService(log),
        lock:        lock,
    }
}

// 扣减库存
func (s *InventoryService) DeductStock(ctx context.Context, productID string, quantity int) error {
    lockKey := "inventory:" + productID
    
    return s.lock.WithLock(ctx, lockKey, &lock.LockOptions{
        Expiration: 30 * time.Second,
        RetryDelay: 100 * time.Millisecond,
        MaxRetries: 10,
    }, func(ctx context.Context) error {
        // 检查库存
        currentStock, err := s.getCurrentStock(ctx, productID)
        if err != nil {
            return err
        }
        
        if currentStock < quantity {
            return fmt.Errorf("库存不足")
        }
        
        // 扣减库存
        return s.updateStock(ctx, productID, currentStock-quantity)
    })
}

// 使用互斥锁
func (s *InventoryService) UpdatePrice(ctx context.Context, productID string, price float64) error {
    mutex := lock.NewMutex(s.lock, "price:"+productID)
    
    return mutex.WithLock(ctx, func(ctx context.Context) error {
        // 更新价格逻辑
        return s.updatePrice(ctx, productID, price)
    })
}
```

## 5. 熔断器使用

```go
package service

import (
    "context"
    "goweb/pkg/circuit"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type ExternalService struct {
    *base.BaseService
    breaker *circuit.Breaker
}

func NewExternalService(log logger.Logger) *ExternalService {
    // 创建熔断器
    cfg := circuit.DefaultConfig("external-api")
    cfg.ReadyToTrip = func(counts circuit.Counts) bool {
        return counts.ConsecutiveFailures >= 3
    }
    cfg.OnStateChange = func(name string, from, to circuit.State) {
        log.Info("熔断器状态变化", "name", name, "from", from, "to", to)
    }
    
    breaker := circuit.NewBreaker(cfg, log)
    
    return &ExternalService{
        BaseService: base.NewBaseService(log),
        breaker:     breaker,
    }
}

// 调用外部 API
func (s *ExternalService) CallExternalAPI(ctx context.Context, data map[string]interface{}) (interface{}, error) {
    return s.breaker.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
        // 调用外部 API
        return s.doExternalAPICall(ctx, data)
    })
}

// 带熔断保护的批量调用
func (s *ExternalService) BatchCallExternalAPI(ctx context.Context, requests []map[string]interface{}) ([]interface{}, error) {
    var results []interface{}
    
    for _, req := range requests {
        result, err := s.breaker.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
            return s.doExternalAPICall(ctx, req)
        })
        
        if err != nil {
            s.LogError("外部 API 调用失败", err, "request", req)
            // 根据业务需求决定是否继续
            continue
        }
        
        results = append(results, result)
    }
    
    return results, nil
}
```

## 6. 完整服务示例

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/monitor"
    "goweb/pkg/storage"
    "goweb/pkg/redis"
    "goweb/pkg/circuit"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.New()
    log := logger.New("advanced-service", cfg.GetString("log.level"))
    
    // 初始化各种组件
    metrics := monitor.NewMetrics("advanced-service", log)
    healthChecker := monitor.NewHealthChecker(log)
    
    // 文件存储
    fileStorage := storage.NewLocalStorage("./uploads", log)
    
    // Redis 管理器（包含缓存、消息队列、分布式锁）
    redisManager := redis.NewManager(redisConfig, log)
    mqClient := redisManager.GetQueue()
    lockClient := redisManager.NewLock("example-lock", 30*time.Second)
    
    // 熔断器
    breaker := circuit.NewBreaker(circuit.DefaultConfig("external-api"), log)
    
    // 创建服务
    orderService := NewOrderService(mqClient, log)
    inventoryService := NewInventoryService(lockClient, log)
    externalService := NewExternalService(log)
    
    // 创建处理器
    fileHandler := NewFileHandler(fileStorage, log)
    orderHandler := NewOrderHandler(orderService, log)
    
    // 设置路由
    r := gin.New()
    
    // 监控中间件
    r.Use(metrics.HTTPMiddleware("advanced-service"))
    
    // 健康检查
    r.GET("/healthz", healthChecker.GinHandler())
    r.GET("/readyz", healthChecker.GinHandler())
    
    // API 路由
    api := r.Group("/api/v1")
    {
        // 文件相关
        api.POST("/files/upload", fileHandler.UploadFile)
        api.GET("/files/:path", fileHandler.GetFile)
        
        // 订单相关
        api.POST("/orders", orderHandler.CreateOrder)
        api.GET("/orders/:id", orderHandler.GetOrder)
    }
    
    // 启动消息消费者
    go func() {
        if err := orderService.StartConsumer(context.Background()); err != nil {
            log.Error("消息消费者启动失败", err)
        }
    }()
    
    // 启动服务
    r.Run(":8080")
}
```

## 7. 配置示例

```yaml
# configs/config.yaml
app:
  name: "advanced-service"
  port: 8080
  env: "development"

# 监控配置
monitor:
  enabled: true
  metrics_path: "/metrics"
  health_path: "/healthz"
  ready_path: "/readyz"

# 文件存储配置
storage:
  type: "local"
  base_path: "./uploads"
  max_size: 10485760  # 10MB
  allowed_types: ["image/jpeg", "image/png", "application/pdf"]

# 消息队列配置
mq:
  type: "redis"
  redis:
    host: "localhost"
    port: 6379
    database: 0

# 熔断器配置
circuit:
  external_api:
    max_requests: 3
    interval: "1m"
    timeout: "1m"
    ready_to_trip: "consecutive_failures >= 3"
```

这些示例展示了 GinForge 框架的高级功能，包括监控、文件存储、消息队列、分布式锁、熔断器等企业级功能。通过这些功能，开发者可以构建高可用、高性能的微服务应用。
