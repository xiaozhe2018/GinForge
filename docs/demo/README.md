# GinForge Demo 使用示例索引

本目录提供"可复制即用"的代码片段与说明，覆盖配置、中间件、数据库、缓存、Swagger、校验、路由与统一响应等核心能力。

## 基础功能示例
- config.md：配置加载与覆盖
- middleware.md：常用中间件注册与链路追踪
- db.md：GORM 管理器、迁移与事务
- cache.md：内存/Redis 缓存
- router_response.md：路由与统一响应
- swagger.md：注解与文档生成
- validation.md：参数校验

## 高级功能示例
- base_classes_usage.md：基类体系使用（BaseService、BaseHandler、BaseRepository、BaseController）
- gateway_usage.md：Gateway 客户端与服务间通信
- redis_usage.md：Redis 统一包使用（缓存、消息队列、分布式锁）
- queue_usage.md：消息队列使用（发布订阅、延迟消息、死信队列）
- advanced_features.md：高级功能使用（监控、文件存储、熔断器）

## 运行示例

```bash
# 启动 demo 服务
go run ./services/demo/cmd/server

# 测试 API
curl http://localhost:8080/api/v1/data
curl http://localhost:8080/api/v1/user/123
```

建议：按顺序阅读，从 config → middleware → router_response → db/cache → swagger → validation。


