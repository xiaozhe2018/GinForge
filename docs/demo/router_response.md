# 路由与统一响应示例（router + response）

```go
package demo

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/response"
)

type DemoHandler struct{}

func NewDemoHandler()*DemoHandler{ return &DemoHandler{} }

func (h *DemoHandler) GetData(c *gin.Context){
    response.Success(c, gin.H{"message":"ok"})
}

func (h *DemoHandler) Bad(c *gin.Context){
    response.BadRequest(c, "参数错误")
}

func Setup(r *gin.Engine, h *DemoHandler){
    api := r.Group("/api/v1")
    api.GET("/data", h.GetData)
    api.GET("/bad", h.Bad)
}
```
