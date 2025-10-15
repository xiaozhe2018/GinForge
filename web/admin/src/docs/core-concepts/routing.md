# 路由管理

学习如何在 GinForge 中定义和管理 API 路由。

## 🎯 路由基础

GinForge 基于 Gin 框架，提供了简洁而强大的路由系统。

### 基本路由定义

```go
package router

import (
    "github.com/gin-gonic/gin"
    "goweb/services/admin-api/internal/handler"
)

func NewRouter(h *handler.UserHandler) *gin.Engine {
    r := gin.Default()
    
    // 定义路由
    r.GET("/users", h.GetUsers)          // GET 请求
    r.POST("/users", h.CreateUser)       // POST 请求
    r.PUT("/users/:id", h.UpdateUser)    // PUT 请求
    r.DELETE("/users/:id", h.DeleteUser) // DELETE 请求
    
    return r
}
```

## 📁 路由组织

### 1. 使用路由组

```go
func NewRouter() *gin.Engine {
    r := gin.Default()
    
    // API v1 路由组
    v1 := r.Group("/api/v1")
    {
        // 公开路由
        v1.GET("/health", healthCheck)
        v1.POST("/login", login)
        
        // 需要认证的路由
        auth := v1.Group("/admin")
        auth.Use(middleware.JWT())  // 添加 JWT 中间件
        {
            // 用户相关
            auth.GET("/users", getUserList)
            auth.GET("/users/:id", getUser)
            auth.POST("/users", createUser)
            auth.PUT("/users/:id", updateUser)
            auth.DELETE("/users/:id", deleteUser)
            
            // 角色相关
            auth.GET("/roles", getRoleList)
            auth.POST("/roles", createRole)
        }
    }
    
    return r
}
```

### 2. 实际项目示例

```go
// services/admin-api/internal/router/router.go
func NewRouter(db *gorm.DB, redisClient *redis.Client, log logger.Logger, cfg *config.Config) *gin.Engine {
    r := gin.New()
    
    // 添加全局中间件
    r.Use(middleware.Recovery(log))
    r.Use(middleware.RequestID())
    r.Use(middleware.AccessLogger(log))
    r.Use(middleware.CORS())
    
    // API 版本 1
    api := r.Group("/api/v1/admin")
    
    // 公开路由（无需认证）
    api.GET("/system/basic", adminSystemHandler.GetBasicInfo)
    api.POST("/login", adminAuthHandler.Login)
    
    // 需要认证的路由
    auth := api.Group("")
    auth.Use(middleware.JWTAuthWithRedis(cfg.GetString("jwt.secret"), redisClient))
    {
        // 用户管理路由
        auth.GET("/users", adminUserHandler.GetUsers)
        auth.GET("/users/:id", adminUserHandler.GetUser)
        auth.POST("/users", adminUserHandler.CreateUser)
        auth.PUT("/users/:id", adminUserHandler.UpdateUser)
        auth.PUT("/users/:id/status", adminUserHandler.UpdateUserStatus)
        auth.DELETE("/users/:id", adminUserHandler.DeleteUser)
        
        // 角色管理路由
        auth.GET("/roles", adminRoleHandler.GetRoles)
        auth.POST("/roles", adminRoleHandler.CreateRole)
        auth.PUT("/roles/:id", adminRoleHandler.UpdateRole)
        auth.DELETE("/roles/:id", adminRoleHandler.DeleteRole)
        
        // 系统管理路由
        auth.GET("/system/info", adminSystemHandler.GetSystemInfo)
        auth.GET("/system/configs", adminSystemHandler.GetConfigs)
        auth.PUT("/system/configs/:key", adminSystemHandler.UpdateConfig)
    }
    
    return r
}
```

## 🔐 路由中间件

### 全局中间件

```go
r := gin.New()

// 全局中间件（应用于所有路由）
r.Use(middleware.Recovery(log))       // Panic 恢复
r.Use(middleware.RequestID())         // 请求 ID
r.Use(middleware.AccessLogger(log))   // 访问日志
r.Use(middleware.CORS())              // 跨域处理
```

### 路由组中间件

```go
// 只对特定路由组使用中间件
auth := r.Group("/api/v1/admin")
auth.Use(middleware.JWT())            // JWT 认证
auth.Use(middleware.RateLimit(100))   // 限流
auth.Use(middleware.OperationLog())   // 操作日志
{
    auth.GET("/users", getUsers)
}
```

### 单个路由中间件

```go
// 只对单个路由使用中间件
r.GET("/admin/sensitive", 
    middleware.RateLimit(10),          // 限流
    middleware.AuditLog(),             // 审计日志
    sensitiveHandler,                  // 处理函数
)
```

## 📋 路由参数

### 1. 路径参数

```go
// 定义带参数的路由
r.GET("/users/:id", func(c *gin.Context) {
    // 获取路径参数
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})

// 访问：GET /users/123
// 结果：{"user_id": "123"}
```

### 2. 查询参数

```go
r.GET("/users", func(c *gin.Context) {
    // 获取查询参数
    page := c.DefaultQuery("page", "1")      // 带默认值
    pageSize := c.Query("page_size")         // 不带默认值
    keyword := c.Query("keyword")
    
    c.JSON(200, gin.H{
        "page": page,
        "page_size": pageSize,
        "keyword": keyword,
    })
})

// 访问：GET /users?page=2&page_size=10&keyword=admin
```

### 3. 请求体参数

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

r.POST("/users", func(c *gin.Context) {
    var req CreateUserRequest
    
    // 绑定 JSON 参数
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"username": req.Username})
})
```

## 🎨 RESTful API 设计

GinForge 推荐使用 RESTful 风格的 API 设计：

| HTTP 方法 | 路径 | 说明 | 示例 |
|-----------|------|------|------|
| GET | `/users` | 获取用户列表 | 分页查询 |
| GET | `/users/:id` | 获取单个用户 | 查看详情 |
| POST | `/users` | 创建用户 | 新增 |
| PUT | `/users/:id` | 更新用户 | 完整更新 |
| PATCH | `/users/:id` | 部分更新用户 | 字段更新 |
| DELETE | `/users/:id` | 删除用户 | 删除 |

### 完整示例

```go
// 用户资源的 RESTful 路由
users := r.Group("/api/v1/users")
users.Use(middleware.JWT())
{
    users.GET("", userHandler.List)           // 列表
    users.GET("/:id", userHandler.Get)        // 详情
    users.POST("", userHandler.Create)        // 创建
    users.PUT("/:id", userHandler.Update)     // 更新
    users.PATCH("/:id", userHandler.Patch)    // 部分更新
    users.DELETE("/:id", userHandler.Delete)  // 删除
    
    // 子资源
    users.GET("/:id/orders", userHandler.GetOrders)  // 用户的订单
    users.GET("/:id/profile", userHandler.GetProfile) // 用户资料
}
```

## 🔄 路由器实例方法

### 常用路由方法

```go
// HTTP 方法路由
r.GET("/path", handler)
r.POST("/path", handler)
r.PUT("/path", handler)
r.DELETE("/path", handler)
r.PATCH("/path", handler)
r.HEAD("/path", handler)
r.OPTIONS("/path", handler)

// 处理所有 HTTP 方法
r.Any("/path", handler)

// 处理多个 HTTP 方法
r.Match([]string{"GET", "POST"}, "/path", handler)

// 静态文件
r.Static("/assets", "./public/assets")
r.StaticFile("/favicon.ico", "./public/favicon.ico")
r.StaticFS("/static", gin.Dir("./public", false))
```

## 📚 路由注释（Swagger）

使用注释自动生成 API 文档：

```go
// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户 ID 获取用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response{data=model.AdminUser}
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /api/v1/admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // 处理逻辑
}
```

生成 Swagger 文档：

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g services/admin-api/cmd/server/main.go -o services/admin-api/docs

# 访问文档
# http://localhost:8083/swagger/index.html
```

## 🛡️ 路由安全

### 1. 认证保护

```go
// 需要登录才能访问
auth := r.Group("/api/v1")
auth.Use(middleware.JWTAuthWithRedis(jwtSecret, redisClient))
{
    auth.GET("/profile", getProfile)
}
```

### 2. 权限控制

```go
// 需要特定权限
admin := r.Group("/api/v1/admin")
admin.Use(middleware.JWT())
admin.Use(middleware.Permission("admin:read"))  // 检查权限
{
    admin.GET("/users", getUsers)
}
```

### 3. 限流保护

```go
// 限制访问频率
r.GET("/api/v1/public", 
    middleware.RateLimit(100),  // 每分钟最多 100 次
    publicHandler,
)
```

## 📖 完整示例

查看实际的路由配置：

- **Admin API**: `services/admin-api/internal/router/router.go`
- **User API**: `services/user-api/internal/router/router.go`
- **Gateway**: `services/gateway/internal/router/router.go`

## 🎯 下一步

- [学习中间件](./middleware) - 掌握中间件的使用
- [数据库操作](./database) - 学习 GORM 数据库操作
- [统一响应](../../demo/router_response.md) - 了解响应格式

---

**提示**: 遵循 RESTful 设计原则可以让你的 API 更规范、更易于理解和维护。

