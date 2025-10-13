# GinForge 框架文档（面向新手与团队）

GinForge 是一个基于 Go + Gin 的微服务开发框架与最佳实践集合，目标是“让开发更加简单”。

你可以把 GinForge 理解为“贴近业务、拿来即用”的工程化脚手架：统一配置、统一日志、统一中间件、统一响应、统一目录与代码生成，让团队更快地从 0 到 1。

— 如果你喜欢 Gin 的极致性能与简洁 API，那么你会爱上 GinForge 的工程化体验。

## Blog（版本动态与路线）

- v0.1.0（当前）：多端分离微服务架构雏形；统一配置/日志/中间件/响应；模板与生成器上线
- v0.2.0（计划）：集成 Swagger/OpenAPI、AutoMigrate、分页/校验/错误码规范、Prometheus/OTEL
- v0.3.0（计划）：脚手架 CLI、服务发现（Consul/Etcd 可选）、灰度与熔断限流示例

## 介绍（What is GinForge?）

GinForge 是一套围绕 Gin 打造的“工程化微服务框架”：
- 提供统一的基础设施：配置、日志、CORS/Recovery/访问日志、JWT 认证、统一响应
- 提供标准化目录与代码生成：一行命令拉起服务骨架；复制模板即可新增服务
- 提供容器友好部署：Dockerfile、Compose、K8s 清单结构清晰

为什么不是又一个“轮子”？
- 我们不重新实现 Web 框架，直接拥抱 Gin
- 我们补齐的是“工程实践层”的最后一公里：目录、规范、模板、示例、文档

## 快速入门（5 分钟）

前置依赖：Go 1.20+、Git、Docker（可选）

1) 获取代码并安装依赖
```bash
git clone <your-repo> && cd goweb
go mod tidy
```

2) 生成你的第一个服务（示例：catalog）
```bash
go run ./cmd/generator -command=service -name=catalog
```

3) 启动服务
```bash
go run ./services/catalog/cmd/server
```

4) 验证健康检查与示例 API（不同服务路由略有差异）
```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/data
```

5)（可选）一次性拉起网关与内置示例服务
```bash
docker-compose -f deployments/docker-compose.yml up
```

## 基准测试（如何简单压测）

GinForge 使用 Gin 作为核心 HTTP 框架，常见场景下可稳定支撑高并发。你可以用 wrk/ab/vegeta 简单压测：

wrk 示例：
```bash
wrk -t4 -c200 -d30s http://localhost:8080/healthz
```

ab 示例：
```bash
ab -n 50000 -c 200 http://localhost:8080/healthz
```

vegeta 示例：
```bash
echo "GET http://localhost:8080/healthz" | vegeta attack -duration=30s -rate=500 | vegeta report
```

提示：
- 开启 release 模式（生产）可获得更好延迟：`GIN_MODE=release`
- 关闭多余日志、确保中间件轻量、使用复用的 http.Client
- 生产环境请结合 Prometheus + Grafana 做时序指标与告警

## 特性（Features）

- 多端分离 + 微服务模板：user-api / merchant-api / admin-api / gateway
- 统一配置：`.env` 环境变量 + 默认值，云原生友好
- 统一日志：zap JSON 日志，支持结构化字段与链路 ID
- 通用中间件：Recovery、RequestID、访问日志、CORS、JWT 认证
- 统一响应：Success/Error 标准化输出，便于前后端协同
- **Swagger/OpenAPI 文档**：自动生成 API 文档，支持在线测试
- 代码生成：`cmd/generator` 与 `templates/` 快速复制最佳实践
- 部署即用：Docker / Compose / K8s 清晰布局

## 编码（Coding Guide）

目录约定：
```
pkg/            # 共享基础库（config/logger/middleware/response/swagger/...）
services/
  <svc>/
    cmd/server  # 服务入口
    internal/   # handler/service/router（仅服务内部可见）
    docs/       # Swagger 生成的文档
templates/      # 代码生成模板
deployments/    # Docker/K8s
scripts/        # 工具脚本（如 swagger.sh）
```

统一配置：
```go
cfg := config.New()                         // 新版配置API
log := logger.New("service-name", cfg.GetString("log.level"))

// 获取配置值
port := cfg.GetInt("app.port")
dbHost := cfg.GetString("database.host")
jwtSecret := cfg.GetString("jwt.secret")
```

统一中间件：
```go
r := gin.New()
r.Use(middleware.Recovery(log))
r.Use(middleware.RequestID())
r.Use(middleware.AccessLogger(log))
// JWT：
api := r.Group("/api/v1", middleware.JWTAuth(cfg.GetString("jwt.secret")))
```

统一响应：
```go
response.Success(c, gin.H{"message":"ok"})
response.BadRequest(c, "参数错误")
```

服务最小实现：
```go
// service
type DemoService struct{}
func (s *DemoService) GetData()(any,error){ return gin.H{"hello":"world"}, nil }

// handler
type DemoHandler struct{ svc *DemoService }
func (h *DemoHandler) GetData(c *gin.Context){ d,_ := h.svc.GetData(); response.Success(c,d) }
```

Swagger 注解示例：
```go
// @Summary      获取数据
// @Description  获取示例数据
// @Tags         demo
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=object}
// @Router       /demo/data [get]
func (h *DemoHandler) GetData(c *gin.Context) {
    // 实现逻辑
}
```

Swagger 文档生成：
```bash
# 生成所有服务的 Swagger 文档
make swagger

# 访问文档
# 用户端: http://localhost:8081/swagger/index.html
# 商户端: http://localhost:8082/swagger/index.html
```

代码风格要点：
 - 函数与变量使用完整词义、避免缩写；错误优先返回；使用守卫语句
 - 日志使用结构化字段；不可吞异常；避免深层嵌套
 - 组件“高内聚、低耦合”：handler 仅处理协议层，service 聚焦业务

## 新手入门（从 0 到 1）

目标：在 30~60 分钟内，完成“环境准备 → 启动模板服务 → 调用 API → 生成 Swagger → 写第一个接口”。

1) 环境准备
```bash
git clone <your-repo> && cd goweb
go mod tidy
cp env.example .env  # 可选：调整端口、JWT 等
```

2) 启动一个现有服务（示例 user-api）
```bash
go run ./services/user-api/cmd/server
# 浏览器访问： http://localhost:8081/healthz
```

3) 生成并查看 Swagger 文档
```bash
make swagger
# 浏览器访问： http://localhost:8081/swagger/index.html
```

4) 新建一个服务骨架（使用生成器）
```bash
go run ./cmd/generator -command=service -name=catalog
go run ./services/catalo g/cmd/server
```

5) 给新服务添加第一个 API
```go
// services/<svc>/internal/handler/demo_handler.go
type DemoHandler struct{}
func NewDemoHandler() *DemoHandler { return &DemoHandler{} }
func (h *DemoHandler) Ping(c *gin.Context) { response.Success(c, gin.H{"pong": true}) }

// services/<svc>/internal/router/router.go （注册路由）
api := r.Group("/api/v1")
api.GET("/ping", handler.Ping)
```

6) 提交前本地验证
```bash
go build ./...
go test ./...
```

提示：遇到问题先看“常见问题与排错”。

## 配置系统详解（config）

核心目标：一致、可覆盖、可热更（可选）。框架采用 viper + YAML + 环境变量。

- 加载入口：
```go
cfg := config.New()                  // 自动加载 configs/config.yaml + .env + 环境变量
log := logger.New("svc", cfg.GetString("log.level"))
```

- 配置优先级：命令行 > 环境变量 > YAML > 默认值
- 默认值查看：`pkg/config/config.go`（如 app/read_timeout、jwt/issuer 等）
- 常用键位：
  - app.port, app.read_timeout, app.write_timeout, app.idle_timeout
  - log.level, log.format, log.output
  - database.driver/host/port/database/username/password
  - redis.enabled/host/port/password/database/pool_size
  - jwt.secret/expire/issuer/audience
  - external_services.user_api_url 等

示例：
```yaml
# configs/config.yaml
app:
  port: 8081
jwt:
  secret: "dev-secret"
```

```env
# .env 覆盖 YAML
APP_PORT=8088
JWT_SECRET=override-secret
```

```go
port := cfg.GetInt("app.port")           // 8088（被 .env 覆盖）
secret := cfg.GetString("jwt.secret")    // override-secret
```

## 中间件清单与用法（middleware）

- Recovery(log)：全局 panic 保护，输出结构化日志
- RequestID()：为每个请求注入 X-Request-Id，便于链路追踪
- AccessLogger(log)：记录 method/path/status/latency 等关键信息
- CORS：`github.com/gin-contrib/cors`，默认放开，生产建议收敛域名
- JWTAuth(secret)：基于 Authorization: Bearer <token> 校验，失败返回 401
- RateLimit：令牌桶限流（可选）
- Validation：统一参数校验（依赖 validator）
- Cache：返回头缓存/ETag（可选）

注册示例：
```go
r := gin.New()
r.Use(middleware.Recovery(log))
r.Use(middleware.RequestID())
r.Use(middleware.AccessLogger(log))
r.Use(cors.New(cors.Config{ /* ... */ }))
api := r.Group("/api/v1", middleware.JWTAuth(cfg.GetString("jwt.secret")))
```

## 数据库与缓存使用（db, cache）

数据库（GORM 管理器）：`pkg/db/manager.go`
- 支持 sqlite/mysql/postgres；连接池、日志级别、超时等来源于配置
- 典型用法：
```go
dbm := db.NewManager(cfg, log)
gormDB, err := dbm.Open()
if err != nil { log.Fatal("db open error", err) }
// AutoMigrate 示例
gormDB.AutoMigrate(&model.User{})
// 事务示例
err = gormDB.Transaction(func(tx *gorm.DB) error {
  // tx.Create(&user)
  return nil
})
```

缓存（内存 + Redis）：`pkg/cache/cache.go`, `pkg/cache/redis.go`
- 通过简单工厂选择实现：
```go
cm := cache.NewManager(cfg)
ctx := context.Background()
_ = cm.Set(ctx, "key", []byte("value"), time.Minute)
val, _ := cm.Get(ctx, "key")
```
- 场景建议：
  - 读多写少：优先本地内存缓存，命中率高、延迟低
  - 多实例共享：使用 Redis，并设置合理 ttl 与命名空间

## 代码生成与首个 API 教程（generator）

1) 生成服务骨架
```bash
go run ./cmd/generator -command=service -name=catalog
```

2) 实现 Service 与 Handler
```go
// internal/service/catalog_service.go
type CatalogService struct{}
func NewCatalogService()*CatalogService{ return &CatalogService{} }
func (s *CatalogService) List(){ /* ... */ }

// internal/handler/catalog_handler.go
type CatalogHandler struct{ svc *CatalogService }
func NewCatalogHandler(s *CatalogService)*CatalogHandler{ return &CatalogHandler{svc:s} }
func (h *CatalogHandler) List(c *gin.Context){ response.Success(c, gin.H{"items": []any{}}) }
```

3) 注册路由并加注解（Swagger）
```go
// @Summary 列表
// @Tags catalog
// @Success 200 {object} response.Response{data=object}
// @Router /catalog/list [get]
api := r.Group("/api/v1")
api.GET("/catalog/list", handler.List)
```

4) 生成 Swagger 并验证
```bash
make swagger
open http://localhost:<port>/swagger/index.html
```

## Swagger 使用与常见问题

生成
```bash
make swagger
```

接入要点
- 在服务 main.go 文件保留顶层注释（title/version/host/BasePath/security 定义）
- 在 handler 方法上方写注解（Summary/Tags/Router 等）
- 每个服务会生成到 `services/<svc>/docs`

常见问题
- 访问 404：确认对应服务已启动；路径是否 `/swagger/index.html`
- 注解不生效：检查注解格式、是否在目标函数上方；重新执行 `make swagger`
- 跨域：网关或服务侧开启 CORS 中间件；前端使用正确的 Origin

## 常见问题与排错（FAQ）

- 启动端口被占用？
  - 修改 `configs/config.yaml` 的 `app.port` 或 `.env` 中 `APP_PORT`
- JWT 验证失败？
  - 确认 `Authorization: Bearer <token>`；`jwt.secret` 一致；时间未过期
- 数据库连接失败？
  - 驱动与 DSN 是否正确；网络连通；账号权限；`conn_max_lifetime` 是否过短
- Redis 连接失败？
  - `redis.enabled=true`；端口、密码、数据库号正确；检查超时配置
- Swagger 生成失败？
  - `go install github.com/swaggo/swag/cmd/swag@latest`；检查 handler 注解合法性
- 如何区分开发/生产？
  - 通过 `.env` 或容器环境变量设置 `APP_ENV=production`，代码中 `cfg.IsProduction()` 做分支

## 部署（Deploy）

Docker：
```bash
docker build -f deployments/docker/Dockerfile -t ginforge:latest .
docker-compose -f deployments/docker-compose.yml up -d
```

Kubernetes（示例）：
```bash
kubectl apply -f deployments/k8s/
```

上线建议：
- 按服务独立容器化与发布；滚动升级；健康检查与就绪探针
- 外挂配置（ConfigMap/Env）；只读镜像；最小权限
- 采集日志与指标，联动告警

## 示例与测试（Examples & Tests）

示例 API（路由）：
```go
api := r.Group("/api/v1")
api.GET("/data", handler.GetData)
api.POST("/data", handler.CreateData)
```

处理器示例：
```go
func (h *DemoHandler) CreateData(c *gin.Context){
  var req struct{ Name string `json:"name" binding:"required"` }
  if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, err.Error()); return }
  response.Success(c, gin.H{"id":"123", "name": req.Name})
}
```

单元测试（httptest）：
```go
func TestHealthz(t *testing.T){
  r := gin.New(); r.GET("/healthz", func(c *gin.Context){ c.JSON(200, gin.H{"status":"ok"}) })
  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/healthz", nil)
  r.ServeHTTP(w, req)
  if w.Code != http.StatusOK { t.Fatalf("expect 200 got %d", w.Code) }
}
```

集成测试建议：
- 使用 docker-compose 启动依赖（如 DB/Redis），go test -tags=integration
- 隔离测试数据；测试前后清理；关注外部系统超时与降级

更多使用示例请见 `docs/demo/`：

- [demo/config.md](./demo/config.md)：配置加载与覆盖示例
- [demo/middleware.md](./demo/middleware.md)：中间件注册与链路追踪示例
- [demo/db.md](./demo/db.md)：GORM 管理器、事务、迁移示例
- [demo/cache.md](./demo/cache.md)：内存/Redis 缓存示例
- [demo/router_response.md](./demo/router_response.md)：路由与统一响应示例
- [demo/swagger.md](./demo/swagger.md)：Swagger 注解与生成示例
- [demo/validation.md](./demo/validation.md)：参数校验与中间件示例
- [demo/delayed_queue_usage.md](./demo/delayed_queue_usage.md)：延时队列详细使用指南
- [demo/gateway_worker_usage.md](./demo/gateway_worker_usage.md)：Gateway Worker 工作服务使用

---

许可证：MIT
