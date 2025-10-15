# 中间件

中间件是 GinForge 中处理横切关注点的核心机制，用于处理认证、日志、限流等通用功能。

## 🎯 什么是中间件？

中间件是一个函数，它可以在请求到达最终处理函数之前或之后执行特定的逻辑。

```
Client Request
    ↓
[中间件 1: RequestID] → 生成请求 ID
    ↓
[中间件 2: Logger] → 记录请求日志
    ↓
[中间件 3: JWT] → 验证 Token
    ↓
[Handler] → 业务处理
    ↓
Response
```

## 📦 内置中间件

GinForge 提供了丰富的内置中间件：

| 中间件 | 功能 | 位置 |
|--------|------|------|
| `Recovery` | Panic 恢复 | `pkg/middleware/recovery.go` |
| `RequestID` | 请求追踪 | `pkg/middleware/request_id.go` |
| `AccessLogger` | 访问日志 | `pkg/middleware/logger.go` |
| `CORS` | 跨域处理 | `pkg/middleware/cors.go` |
| `JWT` | JWT 认证 | `pkg/middleware/jwt.go` |
| `RateLimit` | 限流 | `pkg/middleware/rate_limit.go` |
| `Cache` | HTTP 缓存 | `pkg/middleware/cache.go` |
| `OperationLog` | 操作日志 | `pkg/middleware/operation_log.go` |
| `Validation` | 参数验证 | `pkg/middleware/validation.go` |

## 🔧 使用中间件

### 1. 全局中间件

应用于所有路由：

```go
r := gin.New()

// 添加全局中间件
r.Use(middleware.Recovery(log))       // Panic 恢复
r.Use(middleware.RequestID())         // 生成请求 ID
r.Use(middleware.AccessLogger(log))   // 访问日志
r.Use(middleware.CORS())              // CORS 处理
```

### 2. 路由组中间件

只应用于特定路由组：

```go
// 公开路由（无需认证）
public := r.Group("/api/v1/public")
{
    public.GET("/health", healthCheck)
}

// 需要认证的路由
auth := r.Group("/api/v1/admin")
auth.Use(middleware.JWTAuthWithRedis(jwtSecret, redisClient))  // 添加 JWT 认证
auth.Use(middleware.OperationLog())                             // 添加操作日志
{
    auth.GET("/users", getUsers)
    auth.POST("/users", createUser)
}
```

### 3. 单个路由中间件

```go
// 只对特定路由使用中间件
r.GET("/api/v1/sensitive",
    middleware.RateLimit(10),       // 限流：每分钟10次
    middleware.AuditLog(),          // 审计日志
    sensitiveHandler,               // 处理函数
)
```

## 🔐 认证中间件

### JWT 认证

```go
// 基础 JWT 认证
r.Use(middleware.JWTAuth(jwtSecret))

// 带 Redis 的 JWT 认证（支持 Token 黑名单）
r.Use(middleware.JWTAuthWithRedis(jwtSecret, redisClient))
```

**使用示例**：

```go
package router

func NewRouter(cfg *config.Config, redisClient *redis.Client) *gin.Engine {
    r := gin.New()
    
    api := r.Group("/api/v1")
    
    // 公开接口
    api.POST("/login", loginHandler)
    
    // 需要认证的接口
    auth := api.Group("")
    auth.Use(middleware.JWTAuthWithRedis(cfg.GetString("jwt.secret"), redisClient))
    {
        auth.GET("/profile", getProfile)
        auth.POST("/logout", logout)
    }
    
    return r
}
```

**在 Handler 中获取用户信息**：

```go
func (h *UserHandler) GetProfile(c *gin.Context) {
    // 从 Context 中获取用户 ID
    userID := c.GetString("user_id")
    username := c.GetString("username")
    
    log.Info("user accessing profile", "user_id", userID, "username", username)
    
    // ... 业务逻辑
}
```

## 📝 日志中间件

### 访问日志

自动记录所有 HTTP 请求：

```go
r.Use(middleware.AccessLogger(log))
```

**日志输出示例**：

```json
{
  "level": "info",
  "ts": "2025-10-15T12:00:00+0700",
  "msg": "access",
  "status": 200,
  "method": "GET",
  "path": "/api/v1/users",
  "client_ip": "::1",
  "latency": 0.0123,
  "request_id": "xxx-xxx-xxx"
}
```

### 操作日志

记录用户操作到数据库：

```go
auth.Use(middleware.OperationLog())
```

**功能**：
- 记录用户 ID 和用户名
- 记录请求方法和路径
- 记录 IP 地址和 User-Agent
- 记录请求参数和响应数据
- 记录响应状态和耗时

## 🚦 限流中间件

### 使用限流

```go
// 全局限流：每分钟 1000 次
r.Use(middleware.RateLimit(1000))

// API 限流：每分钟 100 次
api.Use(middleware.RateLimit(100))

// 单个路由限流：每分钟 10 次
r.GET("/sensitive", middleware.RateLimit(10), handler)
```

**限流算法**：基于 Redis 的滑动窗口算法

**超出限制时的响应**：

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "trace_id": "xxx"
}
```

## 💾 缓存中间件

### HTTP 缓存

```go
// 缓存 GET 请求的响应，5 分钟有效
r.GET("/api/v1/config",
    middleware.Cache(5 * time.Minute, redisClient),
    getConfig,
)
```

**工作原理**：
1. 首次请求：执行 Handler，将响应缓存到 Redis
2. 后续请求：直接从 Redis 返回缓存的响应
3. 过期后：重新执行 Handler 并更新缓存

**缓存键**：基于请求路径和查询参数生成

## 🌍 CORS 中间件

### 跨域配置

```go
r.Use(middleware.CORS())
```

**`pkg/middleware/cors.go` 配置**：

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
        c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, X-CSRF-Token, Token, session, X-Requested-With")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")
        
        // 处理 OPTIONS 预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

## 🛠️ 自定义中间件

### 创建自定义中间件

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "time"
)

// CustomMiddleware 自定义中间件示例
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求前处理
        startTime := time.Now()
        
        // 设置自定义变量
        c.Set("custom_key", "custom_value")
        
        // 继续处理请求
        c.Next()
        
        // 请求后处理
        latency := time.Since(startTime)
        log.Printf("Request took %v", latency)
        
        // 可以修改响应
        // c.Header("X-Custom-Header", "value")
    }
}
```

### 使用自定义中间件

```go
r.Use(CustomMiddleware())
```

### 中间件最佳实践

```go
func AuthorizationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Token
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未提供认证令牌"})
            c.Abort()  // 终止请求
            return
        }
        
        // 2. 验证 Token
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "令牌无效"})
            c.Abort()
            return
        }
        
        // 3. 设置用户信息到 Context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        
        // 4. 继续处理
        c.Next()
    }
}
```

## 📊 中间件执行顺序

中间件按照注册顺序执行：

```go
r.Use(A())  // 第 1 个执行
r.Use(B())  // 第 2 个执行
r.Use(C())  // 第 3 个执行

r.GET("/path", handler)

// 执行顺序：
// A (before) → B (before) → C (before) → Handler → C (after) → B (after) → A (after)
```

### 示例流程

```go
r.Use(middleware.RequestID())      // 1. 生成请求 ID
r.Use(middleware.AccessLogger(log)) // 2. 开始记录日志
r.Use(middleware.CORS())            // 3. 处理跨域
r.Use(middleware.JWT())             // 4. 验证 Token

api.GET("/users", handler)          // 5. 业务处理

// 执行后：
// 6. 记录响应日志
// 7. 返回客户端
```

## 🎯 中间件组合

### 推荐的中间件组合

```go
func NewRouter(cfg *config.Config, log logger.Logger, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
    r := gin.New()
    
    // ==== 第一层：基础中间件 ====
    r.Use(middleware.Recovery(log))       // Panic 恢复
    r.Use(middleware.RequestID())         // 请求 ID
    r.Use(middleware.AccessLogger(log))   // 访问日志
    r.Use(middleware.CORS())              // 跨域
    
    // ==== 第二层：API 路由 ====
    api := r.Group("/api/v1")
    
    // 公开路由
    api.POST("/login", loginHandler)
    
    // ==== 第三层：认证路由 ====
    auth := api.Group("")
    auth.Use(middleware.JWTAuthWithRedis(cfg.GetString("jwt.secret"), redisClient))
    auth.Use(middleware.OperationLog())    // 操作日志（需要用户信息）
    {
        // ==== 第四层：特定路由的限流 ====
        auth.GET("/users", 
            middleware.RateLimit(100),     // 限流
            getUserList,
        )
        
        // ==== 第五层：带缓存的路由 ====
        auth.GET("/config",
            middleware.Cache(5*time.Minute, redisClient),  // 5分钟缓存
            getConfig,
        )
    }
    
    return r
}
```

## 📚 实际示例

### Admin API 中间件配置

查看完整配置：`services/admin-api/internal/router/router.go`

```go
r := gin.New()

// 基础中间件
r.Use(middleware.Recovery(log))
r.Use(middleware.RequestID())
r.Use(middleware.AccessLogger(log))
r.Use(middleware.CORS())

// API 路由
api := r.Group("/api/v1/admin")

// 公开路由
api.GET("/system/basic", adminSystemHandler.GetBasicInfo)
api.POST("/login", adminAuthHandler.Login)

// 认证路由
auth := api.Group("")
auth.Use(middleware.JWTAuthWithRedis(cfg.GetString("jwt.secret"), redisClient))
{
    // 用户管理
    auth.GET("/users", adminUserHandler.GetUsers)
    auth.POST("/users", adminUserHandler.CreateUser)
    
    // 系统管理
    auth.GET("/system/info", adminSystemHandler.GetSystemInfo)
}
```

## 🎨 中间件开发指南

### 中间件模板

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/logger"
)

// MyMiddleware 自定义中间件
func MyMiddleware(log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ===== 请求前处理 =====
        log.Info("before request")
        
        // 可以设置变量到 Context
        c.Set("custom_key", "value")
        
        // 可以修改请求
        // c.Request.Header.Set("X-Custom", "value")
        
        // ===== 继续处理请求 =====
        c.Next()
        
        // ===== 请求后处理 =====
        log.Info("after request", "status", c.Writer.Status())
        
        // 可以修改响应头
        // c.Header("X-Custom-Response", "value")
    }
}
```

### 中断请求

如果验证失败，可以中止请求：

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            // 返回错误并中止
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()  // 中止后续中间件和 Handler
            return
        }
        
        c.Next()  // 通过验证，继续
    }
}
```

## 💡 使用技巧

### 1. 从 Context 获取数据

```go
// 在中间件中设置
c.Set("user_id", "123")

// 在 Handler 中获取
userID := c.GetString("user_id")
```

### 2. 获取请求 ID

```go
// RequestID 中间件会自动生成
requestID := c.GetString("request_id")
```

### 3. 检查用户权限

```go
func PermissionMiddleware(required string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取用户权限列表
        permissions, _ := c.Get("permissions")
        
        // 检查权限
        if !hasPermission(permissions, required) {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 使用
auth.GET("/users", 
    PermissionMiddleware("user:read"),
    getUsers,
)
```

## 🚀 性能优化

### 1. 限流保护

防止 API 被滥用：

```go
// 公开 API：严格限流
public.Use(middleware.RateLimit(50))

// 认证 API：适度限流
auth.Use(middleware.RateLimit(200))

// 高频 API：宽松限流
auth.GET("/heartbeat", middleware.RateLimit(1000), heartbeat)
```

### 2. 缓存加速

对不常变化的数据使用缓存：

```go
// 系统配置：缓存 10 分钟
auth.GET("/system/config",
    middleware.Cache(10*time.Minute, redisClient),
    getSystemConfig,
)

// 用户列表：缓存 1 分钟
auth.GET("/users",
    middleware.Cache(1*time.Minute, redisClient),
    getUsers,
)
```

## 📖 实际示例

查看更多示例：

- **中间件使用**: `docs/demo/middleware.md`
- **JWT 认证示例**: `services/admin-api/internal/router/router.go`
- **限流示例**: `pkg/middleware/rate_limit.go`
- **缓存示例**: `pkg/middleware/cache.go`

## 🎯 下一步

- [数据库操作](./database) - 学习 GORM 数据库操作
- [统一响应](../../demo/router_response.md) - 了解响应格式
- [认证授权](../features/authentication) - 深入了解认证系统

---

**提示**: 合理使用中间件可以大大简化代码，提高系统的可维护性和安全性。

