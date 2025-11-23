# GinForge 高级功能

本文档介绍了 GinForge 框架的高级功能，包括测试框架、CLI工具和服务网格。

## 🧪 测试框架

### 功能特性

- **单元测试**: 完整的单元测试工具包
- **集成测试**: 支持服务间集成测试
- **基准测试**: 性能基准测试支持
- **模拟测试**: HTTP客户端和服务模拟
- **测试套件**: 结构化测试管理

### 快速开始

```go
// 创建测试套件
suite := testing.NewTestSuite(t)
suite.Setup()
defer suite.Teardown()

// 创建测试服务器
server := testing.NewTestServer(t, func(r *gin.Engine) {
    r.GET("/api/v1/data", handler.GetData)
})
defer server.Close()

// 运行测试
resp := server.Get("/api/v1/data")
resp.AssertStatus(200)
resp.AssertContains("success")
```

### 测试工具

- `TestServer`: HTTP测试服务器
- `TestResponse`: 响应断言工具
- `TestSuite`: 测试套件管理
- `MockHTTPClient`: HTTP客户端模拟
- `TestHelper`: 测试辅助函数

## 🛠️ CLI工具

### 功能特性

- **服务生成**: 一键生成微服务骨架
- **代码生成**: 自动生成处理器、模型等
- **测试运行**: 运行测试并生成报告
- **配置管理**: 配置查询和更新
- **部署支持**: 服务部署和管理

### 安装CLI

```bash
# 构建CLI工具
make build

# 安装到系统
make install-cli

# 使用CLI
ginforge --help
```

### 常用命令

```bash
# 创建新服务
ginforge service --name=payment --port=8086

# 创建处理器
ginforge handler --service=user --name=profile

# 创建数据模型
ginforge model --name=user --fields=name,email,age

# 运行测试
ginforge test --service=user --coverage

# 部署服务
ginforge deploy --env=production
```

## 🌐 服务网格 (Istio)

### 功能特性

- **流量管理**: 智能路由和负载均衡
- **安全策略**: mTLS和授权策略
- **可观测性**: 指标、日志和追踪
- **故障恢复**: 熔断和重试机制
- **配置管理**: 动态配置更新

### 快速开始

```go
// 创建Istio管理器
istioManager := mesh.NewIstioManager(cfg, log)

// 生成Sidecar配置
annotations := istioManager.GenerateSidecarAnnotation()

// 生成VirtualService
vs := istioManager.GenerateVirtualService("user-api", 
    []string{"user.ginforge.local"}, 
    []Route{...})

// 部署配置
err := istioManager.DeployConfig(ctx, vs)
```

### 配置示例

```yaml
# Gateway配置
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: ginforge-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*.ginforge.local"
```

## 📊 监控和可观测性

### 功能特性

- **指标收集**: Prometheus指标
- **日志聚合**: 结构化日志
- **分布式追踪**: Jaeger/Zipkin支持
- **健康检查**: 服务健康状态
- **告警通知**: 异常告警

### 配置示例

```yaml
# Prometheus配置
prometheus:
  enabled: true
  port: 9090
  path: /metrics

# Jaeger配置
jaeger:
  enabled: true
  endpoint: http://jaeger:14268/api/traces
  service_name: ginforge

# 健康检查
health:
  enabled: true
  port: 8080
  path: /healthz
```

## 🚀 部署和运维

### Docker支持

```bash
# 构建镜像
make docker

# 启动服务
make compose

# 停止服务
make compose-down
```

### Kubernetes支持

```bash
# 部署到K8s
kubectl apply -f deployments/k8s/

# 部署Istio配置
kubectl apply -f deployments/k8s/istio/
```

### 环境配置

```yaml
# 开发环境
app:
  env: development
  log_level: debug

# 生产环境
app:
  env: production
  log_level: warn
  istio:
    enabled: true
```

## 🔧 开发工具

### 代码生成

```bash
# 生成服务
ginforge service --name=payment

# 生成处理器
ginforge handler --service=payment --name=order

# 生成模型
ginforge model --name=order --fields=id,amount,status
```

### 测试工具

```bash
# 运行单元测试
make test

# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage

# 运行基准测试
make benchmark
```

### 代码质量

```bash
# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 依赖检查
go mod tidy
```

## 📚 最佳实践

### 1. 测试策略

- 单元测试覆盖核心业务逻辑
- 集成测试验证服务间交互
- 基准测试确保性能要求
- 模拟测试隔离外部依赖

### 2. 服务网格

- 启用mTLS确保安全
- 配置流量管理策略
- 监控服务健康状态
- 实现故障恢复机制

### 3. 监控告警

- 收集关键业务指标
- 设置合理的告警阈值
- 实现日志聚合和分析
- 建立故障响应流程

## 🎯 总结

GinForge 提供了完整的企业级微服务开发解决方案，包括：

- ✅ **测试框架**: 完整的测试工具和最佳实践
- ✅ **CLI工具**: 强大的代码生成和管理工具
- ✅ **服务网格**: Istio集成和配置管理
- ✅ **监控运维**: 完整的可观测性解决方案

通过这些高级功能，开发者可以快速构建、测试、部署和运维高质量的微服务系统。
