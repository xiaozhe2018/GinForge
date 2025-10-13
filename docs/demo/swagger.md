# Swagger 示例（注解与生成）

1) 在服务 main.go 顶部添加基本信息注解：

```go
// @title          示例服务 API
// @version        1.0
// @description    示例服务接口文档
// @host           localhost:8081
// @BasePath       /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

2) 在 handler 方法上添加注解：

```go
// @Summary 获取数据
// @Tags demo
// @Success 200 {object} response.Response{data=object}
// @Router /demo/data [get]
func (h *DemoHandler) GetData(c *gin.Context){
    // ...
}
```

3) 生成文档：

```bash
make swagger
# 浏览器访问：http://localhost:<port>/swagger/index.html
```
