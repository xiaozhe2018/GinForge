# 中间件示例（middleware）

```go
package demo

import (
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/middleware"
)

func ExampleMiddleware() *gin.Engine {
    cfg := config.New()
    log := logger.New("demo", cfg.GetString("log.level"))

    r := gin.New()
    r.Use(middleware.Recovery(log))
    r.Use(middleware.RequestID())
    r.Use(middleware.AccessLogger(log))
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
        ExposeHeaders:    []string{"X-Request-Id"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    api := r.Group("/api/v1", middleware.JWTAuth(cfg.GetString("jwt.secret")))
    api.GET("/ping", func(c *gin.Context){ c.JSON(200, gin.H{"pong": true}) })

    return r
}
```
