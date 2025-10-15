# 基础类

GinForge 提供了四大基础类，简化开发工作，提高代码一致性。

## 🎯 基类体系

```
pkg/base/
├── handler.go      # BaseHandler    - HTTP 处理器基类
├── service.go      # BaseService    - 业务服务基类
├── repository.go   # BaseRepository - 数据仓储基类
└── controller.go   # BaseController - RESTful 控制器基类
```

## 📦 1. BaseHandler

HTTP 请求处理器的基类，提供统一的响应方法。

### 核心方法

```go
type BaseHandler struct {
    logger logger.Logger
}

// 成功响应
Success(c *gin.Context, data interface{})

// 错误响应
BadRequest(c *gin.Context, message string)
Unauthorized(c *gin.Context, message string)
Forbidden(c *gin.Context, message string)
NotFound(c *gin.Context, message string)
InternalError(c *gin.Context, message string)
Error(c *gin.Context, code int, message string)

// 工具方法
GetTraceID(c *gin.Context) string      // 获取请求追踪 ID
GetUserID(c *gin.Context) string       // 获取当前用户 ID
GetClientIP(c *gin.Context) string     // 获取客户端 IP
LogInfo(message string, args ...interface{})   // 记录信息日志
LogError(message string, err error, args ...interface{})  // 记录错误日志
```

### 使用示例

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type UserHandler struct {
    *base.BaseHandler
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService, log logger.Logger) *UserHandler {
    return &UserHandler{
        BaseHandler: base.NewBaseHandler(log),
        userService: userService,
    }
}

func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    // 使用基类的日志方法
    h.LogInfo("getting user", "user_id", userID, "trace_id", h.GetTraceID(c))
    
    // 调用服务层
    user, err := h.userService.GetUserByID(userID)
    if err != nil {
        h.LogError("failed to get user", err, "user_id", userID)
        h.InternalError(c, "获取用户失败")
        return
    }
    
    // 使用基类的成功响应方法
    h.Success(c, user)
}
```

## 📝 2. BaseService

业务逻辑层的基类，提供日志记录功能。

### 核心方法

```go
type BaseService struct {
    logger logger.Logger
}

// 日志方法
LogInfo(message string, args ...interface{})
LogWarn(message string, args ...interface{})
LogError(message string, err error, args ...interface{})
LogDebug(message string, args ...interface{})
```

### 使用示例

```go
package service

import (
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type UserService struct {
    *base.BaseService
    userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository, log logger.Logger) *UserService {
    return &UserService{
        BaseService: base.NewBaseService(log),
        userRepo:    userRepo,
    }
}

func (s *UserService) CreateUser(username, email string) (*model.User, error) {
    // 使用基类的日志方法
    s.LogInfo("creating user", "username", username)
    
    // 检查用户是否已存在
    existing, _ := s.userRepo.GetByUsername(username)
    if existing != nil {
        s.LogWarn("user already exists", "username", username)
        return nil, errors.New("用户已存在")
    }
    
    // 创建用户
    user := &model.User{
        Username: username,
        Email:    email,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        s.LogError("failed to create user", err, "username", username)
        return nil, err
    }
    
    s.LogInfo("user created successfully", "user_id", user.ID)
    return user, nil
}
```

## 🗄️ 3. BaseRepository

数据访问层的基类，提供通用的 CRUD 操作。

### 核心方法

```go
type BaseRepository struct {
    db     *gorm.DB
    logger logger.Logger
}

// CRUD 方法
Create(ctx context.Context, model interface{}) error
Update(ctx context.Context, model interface{}) error
Delete(ctx context.Context, model interface{}) error
FindByID(ctx context.Context, model interface{}, id interface{}) error
FindOne(ctx context.Context, model interface{}, query string, args ...interface{}) error
FindAll(ctx context.Context, models interface{}, query string, args ...interface{}) error

// 分页查询
Paginate(ctx context.Context, models interface{}, pagination *model.Pagination) (*model.PaginationResult, error)

// 日志方法
LogInfo(message string, args ...interface{})
LogError(message string, err error, args ...interface{})
```

### 使用示例

```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "goweb/pkg/base"
    "goweb/pkg/model"
    "goweb/pkg/logger"
)

type UserRepository struct {
    *base.BaseRepository
}

func NewUserRepository(db *gorm.DB, log logger.Logger) *UserRepository {
    return &UserRepository{
        BaseRepository: base.NewBaseRepository(db, log),
    }
}

// 使用基类的方法
func (r *UserRepository) Create(user *model.User) error {
    return r.BaseRepository.Create(context.Background(), user)
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
    var user model.User
    err := r.BaseRepository.FindByID(context.Background(), &user, id)
    return &user, err
}

// 自定义方法
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.BaseRepository.FindOne(context.Background(), &user, "username = ?", username)
    if err != nil {
        r.LogError("failed to find user by username", err, "username", username)
        return nil, err
    }
    return &user, nil
}

// 分页查询
func (r *UserRepository) List(req *model.PaginationRequest) (*model.PaginationResult, error) {
    var users []model.User
    return r.BaseRepository.Paginate(context.Background(), &users, req)
}
```

## 🎮 4. BaseController

RESTful 控制器的基类，适用于标准的 RESTful API。

### 核心方法

```go
type BaseController struct {
    *BaseHandler
}

// 获取分页参数
GetPagination(c *gin.Context) *model.Pagination

// 绑定参数
BindJSON(c *gin.Context, obj interface{}) error
BindQuery(c *gin.Context, obj interface{}) error
```

### 使用示例

```go
package controller

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/base"
    "goweb/pkg/logger"
)

type UserController struct {
    *base.BaseController
    userService *service.UserService
}

func NewUserController(userService *service.UserService, log logger.Logger) *UserController {
    return &UserController{
        BaseController: base.NewBaseController(log),
        userService:    userService,
    }
}

func (c *UserController) List(ctx *gin.Context) {
    // 使用基类方法获取分页参数
    pagination := c.GetPagination(ctx)
    
    c.LogInfo("listing users", "page", pagination.Page, "page_size", pagination.PageSize)
    
    // 调用服务层
    users, total, err := c.userService.ListUsers(pagination)
    if err != nil {
        c.InternalError(ctx, "获取用户列表失败")
        return
    }
    
    c.Success(ctx, gin.H{
        "list":  users,
        "total": total,
        "page":  pagination.Page,
    })
}
```

## 🔗 完整使用流程

### 创建完整的用户模块

```go
// 1. 定义模型 (model/user.go)
type User struct {
    ID       uint64 `json:"id" gorm:"primaryKey"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 2. 创建 Repository (repository/user_repository.go)
type UserRepository struct {
    *base.BaseRepository
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
    var user model.User
    err := r.FindByID(context.Background(), &user, id)
    return &user, err
}

// 3. 创建 Service (service/user_service.go)
type UserService struct {
    *base.BaseService
    userRepo *UserRepository
}

func (s *UserService) GetUser(id uint64) (*model.User, error) {
    s.LogInfo("getting user", "id", id)
    return s.userRepo.GetByID(id)
}

// 4. 创建 Handler (handler/user_handler.go)
type UserHandler struct {
    *base.BaseHandler
    userService *UserService
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.userService.GetUser(id)
    if err != nil {
        h.NotFound(c, "用户不存在")
        return
    }
    h.Success(c, user)
}

// 5. 注册路由 (router/router.go)
r.GET("/users/:id", userHandler.GetUser)
```

## 💡 基类的优势

### 1. 代码复用

✅ 避免重复编写日志代码  
✅ 避免重复编写响应代码  
✅ 避免重复编写 CRUD 代码  

### 2. 统一规范

✅ 统一的日志格式  
✅ 统一的响应格式  
✅ 统一的错误处理  

### 3. 易于维护

✅ 修改一处，全局生效  
✅ 代码结构清晰  
✅ 新人容易上手  

## 📚 完整示例

查看完整实现：

- **基类定义**: `pkg/base/`
- **使用示例**: `docs/demo/base_classes_usage.md`
- **实际应用**: `services/admin-api/internal/`

## 🎯 下一步

- [工具函数](./utilities) - 常用工具函数参考
- [配置选项](./config-options) - 完整的配置项说明
- [最佳实践](../best-practices/code-style) - 代码规范

---

**提示**: 所有业务代码都应该继承这些基类，保持代码的一致性！

