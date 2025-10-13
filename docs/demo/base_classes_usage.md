# 基类使用示例

## 1. BaseService 使用示例

```go
package service

import (
    "context"
    "goweb/pkg/base"
    "goweb/pkg/gateway"
    "goweb/pkg/logger"
)

type UserService struct {
    *base.BaseService
    gatewayClient *gateway.Client
}

func NewUserService(gatewayClient *gateway.Client, log logger.Logger) *UserService {
    return &UserService{
        BaseService:   base.NewBaseService(log),
        gatewayClient: gatewayClient,
    }
}

func (s *UserService) GetUser(userID string) (map[string]interface{}, error) {
    s.LogInfo("getting user", "user_id", userID)
    
    // 通过 Gateway 调用其他服务
    resp, err := s.gatewayClient.GetUser(context.Background(), userID)
    if err != nil {
        s.LogError("failed to get user via gateway", err, "user_id", userID)
        return nil, err
    }
    
    s.LogInfo("user retrieved successfully", "user_id", userID)
    return resp.Data.(map[string]interface{}), nil
}

func (s *UserService) CreateUser(userData map[string]interface{}) error {
    s.LogInfo("creating user", "username", userData["username"])
    
    resp, err := s.gatewayClient.CreateUser(context.Background(), userData)
    if err != nil {
        s.LogError("failed to create user", err, "username", userData["username"])
        return err
    }
    
    s.LogInfo("user created successfully", "user_id", resp.Data)
    return nil
}
```

## 2. BaseHandler 使用示例

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/base"
    "goweb/pkg/gateway"
    "goweb/pkg/logger"
)

type UserHandler struct {
    *base.BaseHandler
    gatewayClient *gateway.Client
}

func NewUserHandler(gatewayClient *gateway.Client, log logger.Logger) *UserHandler {
    return &UserHandler{
        BaseHandler:   base.NewBaseHandler(log),
        gatewayClient: gatewayClient,
    }
}

// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} response.Response{data=object}
// @Router /user/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("user_id")
    if userID == "" {
        h.BadRequest(c, "用户ID不能为空")
        return
    }
    
    h.LogInfo("getting user", "user_id", userID, "trace_id", h.GetTraceID(c))
    
    // 通过 Gateway 调用用户服务
    resp, err := h.gatewayClient.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.LogError("failed to get user", err, "user_id", userID)
        h.InternalError(c, "获取用户信息失败")
        return
    }
    
    if resp.Code != 0 {
        h.Error(c, resp.Code, resp.Message)
        return
    }
    
    h.Success(c, resp.Data)
}

// @Summary 创建用户
// @Description 创建新用户
// @Tags user
// @Accept json
// @Produce json
// @Param user body map[string]interface{} true "用户数据"
// @Success 200 {object} response.Response{data=object}
// @Router /user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var userData map[string]interface{}
    if err := c.ShouldBindJSON(&userData); err != nil {
        h.BadRequest(c, "参数格式错误")
        return
    }
    
    h.LogInfo("creating user", "username", userData["username"])
    
    resp, err := h.gatewayClient.CreateUser(c.Request.Context(), userData)
    if err != nil {
        h.LogError("failed to create user", err, "username", userData["username"])
        h.InternalError(c, "创建用户失败")
        return
    }
    
    h.Success(c, resp.Data)
}
```

## 3. BaseRepository 使用示例

```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "goweb/pkg/base"
    "goweb/pkg/logger"
    "goweb/pkg/model"
)

type UserRepository struct {
    *base.BaseRepository
}

func NewUserRepository(db *gorm.DB, log logger.Logger) *UserRepository {
    return &UserRepository{
        BaseRepository: base.NewBaseRepository(db, log),
    }
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
    var user model.User
    err := r.FindByID(ctx, &user, userID)
    if err != nil {
        r.LogError("failed to find user by id", err, "user_id", userID)
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.FindOne(ctx, &user, "email = ?", email)
    if err != nil {
        r.LogError("failed to find user by email", err, "email", email)
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    err := r.Create(ctx, user)
    if err != nil {
        r.LogError("failed to create user", err, "username", user.Username)
        return err
    }
    return nil
}

func (r *UserRepository) List(ctx context.Context, pagination *model.Pagination) (*model.PaginationResult, error) {
    var users []model.User
    result, err := r.Paginate(ctx, &users, pagination)
    if err != nil {
        r.LogError("failed to list users", err)
        return nil, err
    }
    return result, nil
}
```

## 4. BaseController 使用示例

```go
package controller

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/base"
    "goweb/pkg/logger"
    "goweb/pkg/model"
)

type UserController struct {
    *base.BaseController
}

func NewUserController(log logger.Logger) *UserController {
    return &UserController{
        BaseController: base.NewBaseController(log),
    }
}

func (c *UserController) GetUser(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    if userID == "" {
        c.BadRequest(ctx, "用户ID不能为空")
        return
    }
    
    c.LogInfo("getting user", "user_id", userID, "client_ip", c.GetClientIP(ctx))
    
    // 业务逻辑...
    ctx.JSON(200, gin.H{"user_id": userID})
}

func (c *UserController) ListUsers(ctx *gin.Context) {
    pagination := c.GetPagination(ctx)
    c.LogInfo("listing users", "page", pagination.Page, "page_size", pagination.PageSize)
    
    // 业务逻辑...
    ctx.JSON(200, gin.H{"users": []interface{}{}})
}
```

## 5. 服务注册与发现使用示例

```go
package main

import (
    "context"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/service"
)

func main() {
    cfg := config.New()
    log := logger.New("my-service", cfg.GetString("log.level"))
    
    // 创建服务注册表
    registry := service.NewServiceRegistry(cfg, log)
    
    // 注册服务
    registry.Register(service.ServiceInfo{
        Name:    "user-api",
        Host:    "localhost",
        Port:    8081,
        Version: "1.0.0",
        Status:  "active",
    })
    
    // 创建服务客户端
    userClient := service.NewServiceClient(registry, "user-api")
    
    // 调用服务
    resp, err := userClient.Get(context.Background(), "/api/v1/user/123")
    if err != nil {
        log.Error("failed to call user service", err)
        return
    }
    
    log.Info("user service response", "code", resp.Code, "data", resp.Data)
}
```

## 6. 错误码使用示例

```go
package service

import (
    "goweb/pkg/errors"
    "goweb/pkg/response"
)

func (s *UserService) ValidateUser(userData map[string]interface{}) error {
    if userData["username"] == nil {
        return errors.New(errors.MissingParameter, "用户名不能为空")
    }
    
    if userData["email"] == nil {
        return errors.New(errors.MissingParameter, "邮箱不能为空")
    }
    
    // 检查用户是否已存在
    if s.userExists(userData["username"].(string)) {
        return errors.New(errors.UserExists, "用户已存在")
    }
    
    return nil
}

// 在 Handler 中使用
func (h *UserHandler) CreateUser(c *gin.Context) {
    var userData map[string]interface{}
    if err := c.ShouldBindJSON(&userData); err != nil {
        h.Error(c, errors.InvalidFormat, "参数格式错误")
        return
    }
    
    if err := h.userService.ValidateUser(userData); err != nil {
        h.Error(c, errors.GetCode(err), errors.GetMessage(err))
        return
    }
    
    h.Success(c, gin.H{"message": "用户创建成功"})
}
```

## 7. 常量使用示例

```go
package service

import (
    "goweb/pkg/constants"
    "goweb/pkg/model"
)

func (s *UserService) CreateUser(userData map[string]interface{}) (*model.User, error) {
    user := &model.User{
        Username: userData["username"].(string),
        Email:    userData["email"].(string),
        Status:   constants.UserStatusActive, // 使用常量
    }
    
    // 验证用户名长度
    if len(user.Username) < constants.UsernameMinLength {
        return nil, errors.New(errors.InvalidParameter, "用户名长度不能少于3位")
    }
    
    return user, nil
}
```

## 8. 完整的主函数示例

```go
package main

import (
    "goweb/pkg/config"
    "goweb/pkg/gateway"
    "goweb/pkg/logger"
    "goweb/services/my-service/internal/handler"
    "goweb/services/my-service/internal/router"
    "goweb/services/my-service/internal/service"
)

func main() {
    // 加载配置
    cfg := config.New()
    log := logger.New("my-service", cfg.GetString("log.level"))
    
    // 创建 Gateway 客户端
    gatewayClient := gateway.NewClient(cfg, log)
    
    // 创建服务
    userService := service.NewUserService(gatewayClient, log)
    
    // 创建处理器
    userHandler := handler.NewUserHandler(gatewayClient, log)
    
    // 创建路由
    r := router.NewRouter(cfg, log, userHandler)
    
    // 启动服务
    // ... 启动逻辑
}
```

这些基类提供了：

1. **统一的日志记录**：所有基类都支持结构化日志
2. **统一的错误处理**：标准化的错误响应
3. **统一的上下文传递**：支持链路追踪
4. **统一的数据库操作**：CRUD 操作封装
5. **统一的 Gateway 调用**：服务间通信
6. **统一的配置管理**：配置获取和验证
7. **统一的常量定义**：业务常量管理
8. **统一的错误码**：标准化错误码体系

使用这些基类可以大大简化开发工作，提高代码的一致性和可维护性。
