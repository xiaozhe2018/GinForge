# Gateway 客户端使用示例

## 在其他服务中调用 Gateway

### 1. 基础使用

```go
package service

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/gateway"
    "goweb/pkg/logger"
)

type UserService struct {
    gatewayClient *gateway.Client
    logger        logger.Logger
}

func NewUserService(cfg *config.Config, log logger.Logger) *UserService {
    return &UserService{
        gatewayClient: gateway.NewClient(cfg, log),
        logger:        log,
    }
}

// 通过 Gateway 调用其他服务
func (s *UserService) GetMerchantInfo(ctx context.Context, merchantID string) (map[string]interface{}, error) {
    resp, err := s.gatewayClient.GetMerchant(ctx, merchantID)
    if err != nil {
        s.logger.Error("failed to get merchant", err, "merchant_id", merchantID)
        return nil, err
    }
    
    if resp.Code != 0 {
        return nil, fmt.Errorf("merchant service error: %s", resp.Message)
    }
    
    return resp.Data.(map[string]interface{}), nil
}
```

### 2. 使用服务注册表

```go
package service

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/service"
    "goweb/pkg/logger"
)

type OrderService struct {
    serviceRegistry *service.ServiceRegistry
    userClient      *service.ServiceClient
    merchantClient  *service.ServiceClient
    logger          logger.Logger
}

func NewOrderService(cfg *config.Config, log logger.Logger) *OrderService {
    registry := service.NewServiceRegistry(cfg, log)
    
    // 注册服务
    registry.Register(service.ServiceInfo{
        Name:    "user-api",
        Host:    "localhost",
        Port:    8081,
        Version: "1.0.0",
        Status:  "active",
    })
    
    registry.Register(service.ServiceInfo{
        Name:    "merchant-api", 
        Host:    "localhost",
        Port:    8082,
        Version: "1.0.0",
        Status:  "active",
    })
    
    return &OrderService{
        serviceRegistry: registry,
        userClient:      service.NewServiceClient(registry, "user-api"),
        merchantClient:  service.NewServiceClient(registry, "merchant-api"),
        logger:          log,
    }
}

// 创建订单时调用用户和商户服务
func (s *OrderService) CreateOrder(ctx context.Context, orderData map[string]interface{}) error {
    userID := orderData["user_id"].(string)
    merchantID := orderData["merchant_id"].(string)
    
    // 验证用户
    userResp, err := s.userClient.Get(ctx, "/api/v1/user/" + userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    // 验证商户
    merchantResp, err := s.merchantClient.Get(ctx, "/api/v1/merchant/" + merchantID)
    if err != nil {
        return fmt.Errorf("failed to get merchant: %w", err)
    }
    
    s.logger.Info("order created", "user_id", userID, "merchant_id", merchantID)
    return nil
}
```

### 3. 在 Handler 中使用

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/base"
    "goweb/pkg/gateway"
)

type OrderHandler struct {
    *base.BaseHandler
    gatewayClient *gateway.Client
}

func NewOrderHandler(gatewayClient *gateway.Client, log logger.Logger) *OrderHandler {
    return &OrderHandler{
        BaseHandler:   base.NewBaseHandler(log),
        gatewayClient: gatewayClient,
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
    
    // 通过 Gateway 调用用户服务验证用户
    userID := orderData["user_id"].(string)
    userResp, err := h.gatewayClient.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.LogError("failed to get user", err, "user_id", userID)
        h.InternalError(c, "获取用户信息失败")
        return
    }
    
    if userResp.Code != 0 {
        h.Error(c, userResp.Code, userResp.Message)
        return
    }
    
    h.Success(c, gin.H{
        "order_id": "12345",
        "user":     userResp.Data,
        "status":   "created",
    })
}
```

### 4. 配置示例

```yaml
# configs/config.yaml
gateway:
  base_url: "http://localhost:8080"
  timeout: "30s"
  retry_count: 3
  retry_delay: "1s"

# 服务发现配置
service_discovery:
  type: "static"  # static, consul, etcd
  services:
    user-api:
      host: "localhost"
      port: 8081
      version: "1.0.0"
      status: "active"
    merchant-api:
      host: "localhost" 
      port: 8082
      version: "1.0.0"
      status: "active"
```

### 5. 错误处理

```go
func (s *UserService) CallExternalService(ctx context.Context) error {
    resp, err := s.gatewayClient.Get(ctx, "/api/v1/some-service")
    if err != nil {
        // 网络错误
        s.logger.Error("network error", err)
        return fmt.Errorf("service unavailable: %w", err)
    }
    
    if resp.Code != 0 {
        // 业务错误
        s.logger.Error("business error", nil, "code", resp.Code, "message", resp.Message)
        return fmt.Errorf("service error %d: %s", resp.Code, resp.Message)
    }
    
    return nil
}
```
