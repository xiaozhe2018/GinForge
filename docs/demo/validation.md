# 参数校验示例（validation）

```go
package demo

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/response"
)

type CreateReq struct {
    Name  string `json:"name" binding:"required,min=2,max=32"`
    Email string `json:"email" binding:"required,email"`
}

func Create(c *gin.Context){
    var req CreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    response.Success(c, gin.H{"id":"123", "name": req.Name})
}

func Register(r *gin.Engine){
    api := r.Group("/api/v1")
    api.POST("/create", Create)
}
```
