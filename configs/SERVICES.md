# 服务配置文件清单

## 服务列表

| 服务名称 | 配置文件 | 默认端口 | 说明 |
|---------|---------|---------|------|
| admin-api | `admin-api.yaml` | 8083 | 管理后台API服务 |
| user-api | `user-api.yaml` | 8081 | 用户API服务 |
| merchant-api | `merchant-api.yaml` | 8082 | 商户API服务 |
| gateway | `gateway.yaml` | 8080 | API网关服务 |
| gateway-worker | `gateway-worker.yaml` | 8084 | 网关工作服务 |
| file-api | `file-api.yaml` | 8086 | 文件服务 |
| websocket-gateway | `websocket-gateway.yaml` | 8087 | WebSocket网关 |
| demo | `demo.yaml` | 8085 | 演示服务 |

## 环境配置状态

### 开发环境 (dev/)
- ✅ base.yaml
- ✅ admin-api.yaml
- ✅ user-api.yaml
- ✅ merchant-api.yaml
- ✅ gateway.yaml
- ✅ gateway-worker.yaml
- ✅ file-api.yaml
- ✅ websocket-gateway.yaml
- ✅ demo.yaml

### 测试环境 (test/)
- ✅ base.yaml
- ✅ admin-api.yaml
- ✅ user-api.yaml
- ✅ merchant-api.yaml
- ✅ gateway.yaml
- ✅ file-api.yaml

### 生产环境 (prod/)
- ✅ base.yaml
- ✅ admin-api.yaml
- ✅ user-api.yaml
- ✅ merchant-api.yaml
- ✅ gateway.yaml
- ✅ gateway-worker.yaml
- ✅ file-api.yaml
- ✅ websocket-gateway.yaml
- ✅ demo.yaml

## 使用说明

1. 设置环境变量：
   ```bash
   export GOEASE_APP_ENV=dev  # 或 test, prod
   export SERVICE_NAME=admin-api
   ```

2. 运行服务：
   ```bash
   go run services/admin-api/cmd/server/main.go
   ```

3. 配置会自动从 `configs/{env}/base.yaml` 和 `configs/{env}/{service}.yaml` 加载

## 添加新服务

1. 在对应环境目录下创建服务配置文件
2. 设置 `SERVICE_NAME` 环境变量
3. 配置会自动合并到基础配置之上

