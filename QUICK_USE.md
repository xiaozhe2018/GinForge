# 🚀 GinForge 快速使用指南

> 本指南提供常用操作的快速命令参考，适合已经完成初始化的开发者使用。

## 📑 目录
- [快速启动](#一快速启动)
- [开发调试](#二开发调试)
- [API测试](#三api-测试)
- [数据库操作](#四数据库操作)
- [代码生成](#五代码生成)
- [测试运行](#六测试运行)
- [部署发布](#七部署发布)

## 一、快速启动

### 🔥 方式1：启动所有服务

```bash
# 启动所有后端微服务
make run

# 启动前端开发服务器（新终端）
cd web/admin && npm run dev
```

### 🎯 方式2：只启动管理后台

```bash
# 后端（管理后台API）
go run ./services/admin-api/cmd/server

# 前端（管理后台界面）
cd web/admin && npm run dev
```

### 🛑 停止服务

```bash
# 停止所有后端服务
make stop

# 或者单独停止某个端口
lsof -ti :8083 | xargs kill -9

# 前端按 Ctrl+C 停止
```

### 🔄 重启服务

```bash
# 重启所有服务
make restart

# 或分别重启
make stop && make run
```

### 👀 查看服务状态

```bash
# 查看所有服务端口占用
make status

# 查看具体进程
ps aux | grep "services/"
```

## 二、访问地址

### 🌐 服务端点

| 服务 | 地址 | 说明 |
|------|------|------|
| **前端管理后台** | http://localhost:3000 | Vue3 + Element Plus |
| **管理后台API** | http://localhost:8083 | Admin RESTful API |
| **Swagger文档** | http://localhost:8083/swagger/index.html | 在线API文档 |
| 用户端API | http://localhost:8081 | User API |
| 商户端API | http://localhost:8082 | Merchant API |
| API网关 | http://localhost:8080 | Gateway |
| 网关工作器 | http://localhost:8084 | Gateway Worker |
| 演示服务 | http://localhost:8085 | Demo Service |

### 🔑 默认账号

```
管理后台登录：
用户名: admin
密码: admin123
```

## 三、开发调试

### 🔍 查看日志

```bash
# 后端日志（如果使用文件日志）
tail -f server.log

# 前端开发日志（控制台输出）
# 已在终端显示

# 数据库日志
# MySQL
docker logs -f mysql

# 查看 Redis 日志
docker logs -f redis
```

### 🐛 调试模式

```bash
# 开启 Debug 模式
export APP_DEBUG=true
go run ./services/admin-api/cmd/server

# 使用 Delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./services/admin-api/cmd/server

# 在 VS Code 中调试
# 使用 .vscode/launch.json 配置
```

### 📊 性能分析

```bash
# 启用 pprof 性能分析
# 访问 http://localhost:8083/debug/pprof/

# CPU 分析
curl http://localhost:8083/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# 内存分析
curl http://localhost:8083/debug/pprof/heap > mem.prof
go tool pprof mem.prof
```

## 四、API 测试

### 🧪 使用 curl 测试

#### 1. 登录获取 Token
```bash
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 返回示例：
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "token": "eyJhbGc...",
#     "user": {...},
#     "menus": [...],
#     "permissions": [...]
#   }
# }
```

#### 2. 使用 Token 访问受保护的接口
```bash
# 设置 Token 变量（替换为实际 Token）
TOKEN="eyJhbGc..."

# 获取用户列表
curl -X GET "http://localhost:8083/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 创建用户
curl -X POST http://localhost:8083/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "real_name": "测试用户",
    "phone": "13800138000",
    "password": "123456",
    "role_ids": [2]
  }'

# 更新用户
curl -X PUT http://localhost:8083/api/v1/admin/users/2 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "real_name": "新名字"
  }'

# 删除用户
curl -X DELETE http://localhost:8083/api/v1/admin/users/2 \
  -H "Authorization: Bearer $TOKEN"
```

#### 3. 角色管理
```bash
# 获取角色列表
curl -X GET "http://localhost:8083/api/v1/admin/roles" \
  -H "Authorization: Bearer $TOKEN"

# 创建角色
curl -X POST http://localhost:8083/api/v1/admin/roles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "客服",
    "code": "customer_service",
    "description": "客服角色",
    "menu_ids": [1, 2, 3]
  }'

# 分配权限
curl -X POST http://localhost:8083/api/v1/admin/roles/2/permissions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3, 4, 5]
  }'
```

#### 4. 登出测试
```bash
# 登出
curl -X POST http://localhost:8083/api/v1/admin/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# 再次使用 Token（应该返回 401）
curl -X GET http://localhost:8083/api/v1/admin/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 🎯 使用 Postman/Insomnia

1. **导入 Swagger 文档**
   - URL: http://localhost:8083/swagger/doc.json
   - 自动生成所有接口

2. **配置环境变量**
   ```
   base_url: http://localhost:8083
   token: (登录后获取)
   ```

3. **设置请求头**
   ```
   Authorization: Bearer {{token}}
   Content-Type: application/json
   ```

## 五、数据库操作

### 📊 SQLite 操作

```bash
# 进入数据库
sqlite3 goweb.db

# 查看所有表
.tables

# 查看用户列表
SELECT * FROM admin_users;

# 查看角色列表
SELECT * FROM admin_roles;

# 查看操作日志
SELECT * FROM admin_operation_logs ORDER BY created_at DESC LIMIT 10;

# 退出
.quit
```

### 🐬 MySQL 操作

```bash
# 进入数据库
docker exec -it mysql mysql -uroot -p123456 gin_forge

# 或直接执行命令
docker exec mysql mysql -uroot -p123456 gin_forge -e "SELECT * FROM admin_users;"

# 常用查询
# 查看用户及其角色
SELECT u.username, u.email, r.name as role 
FROM admin_users u 
LEFT JOIN admin_user_roles ur ON u.id = ur.user_id 
LEFT JOIN admin_roles r ON ur.role_id = r.id;

# 查看最近的操作日志
SELECT username, method, path, status_code, created_at 
FROM admin_operation_logs 
ORDER BY created_at DESC 
LIMIT 10;

# 查看角色的权限
SELECT r.name as role, p.name as permission 
FROM admin_roles r 
LEFT JOIN admin_role_permissions rp ON r.id = rp.role_id 
LEFT JOIN admin_permissions p ON rp.permission_id = p.id;
```

### 🔴 Redis 操作

```bash
# 进入 Redis CLI
docker exec -it redis redis-cli

# 查看所有键
KEYS *

# 查看 Token 黑名单
KEYS "token:blacklist:*"

# 查看缓存
KEYS "cache:*"

# 获取某个键的值
GET "token:blacklist:xxx"

# 查看键的过期时间
TTL "token:blacklist:xxx"

# 清空所有数据（慎用！）
FLUSHDB

# 退出
exit
```

## 六、代码生成

### 🛠️ 生成新服务

```bash
# 生成支付服务
go run ./cmd/generator -command=service -name=payment

# 生成订单服务
go run ./cmd/generator -command=service -name=order

# 查看生成的文件
ls -la services/payment/
```

### 📝 生成 Handler

```bash
# 生成用户 Handler
go run ./cmd/cli -command=handler -name=user

# 生成产品 Handler
go run ./cmd/cli -command=handler -name=product
```

### 📄 生成 Swagger 文档

```bash
# 生成所有服务的文档
make swagger

# 单独生成某个服务的文档
swag init -g services/admin-api/cmd/server/main.go \
  -o services/admin-api/docs \
  --parseDependency --parseInternal

# 查看生成的文档
ls -la services/admin-api/docs/
```

## 七、测试运行

### ✅ 运行所有测试

```bash
# 运行所有测试
make test

# 显示详细输出
go test -v ./...

# 只测试某个包
go test -v ./pkg/middleware/...
go test -v ./services/admin-api/internal/service/...
```

### 📊 测试覆盖率

```bash
# 生成覆盖率报告
make test-coverage

# 查看 HTML 报告
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows

# 命令行查看覆盖率
go tool cover -func=coverage.out
```

### 🚀 基准测试

```bash
# 运行基准测试
make benchmark

# 或直接使用 go test
go test -bench=. -benchmem ./...

# 只测试某个包的基准
go test -bench=. -benchmem ./pkg/utils/
```

### 🧪 集成测试

```bash
# 运行集成测试
make test-integration

# 或使用标签
go test -tags=integration ./...
```

## 八、部署发布

### 🏗️ 构建

```bash
# 构建所有服务
make build

# 查看构建结果
ls -la bin/

# 运行构建后的程序
./bin/admin-api
./bin/user-api
```

### 🐳 Docker 部署

```bash
# 构建 Docker 镜像
make docker

# 查看镜像
docker images | grep ginforge

# 使用 Docker Compose 启动
make compose

# 查看运行状态
docker-compose -f deployments/docker-compose.yml ps

# 查看日志
docker-compose -f deployments/docker-compose.yml logs -f

# 停止服务
make compose-down
```

### ☸️ Kubernetes 部署

```bash
# 应用配置
kubectl apply -f deployments/k8s/

# 查看 Pod 状态
kubectl get pods -w

# 查看服务
kubectl get svc

# 查看日志
kubectl logs -f <pod-name>

# 进入容器
kubectl exec -it <pod-name> -- /bin/sh

# 删除部署
kubectl delete -f deployments/k8s/
```

### 🌐 Istio 服务网格

```bash
# 部署 Istio 配置
kubectl apply -f deployments/k8s/istio/

# 查看网关
kubectl get gateway

# 查看虚拟服务
kubectl get virtualservice

# 查看目标规则
kubectl get destinationrule

# 查看流量
istioctl dashboard kiali
```

## 九、常用命令速查

### 📋 Make 命令

```bash
make help          # 查看所有可用命令
make build         # 构建所有服务
make run           # 启动所有服务
make stop          # 停止所有服务
make restart       # 重启所有服务
make status        # 查看服务状态
make test          # 运行测试
make swagger       # 生成 API 文档
make clean         # 清理构建文件
make docker        # 构建 Docker 镜像
make compose       # 启动 Docker Compose
make web-dev       # 启动前端开发服务器
make web-build     # 构建前端生产版本
```

### 🔧 Go 命令

```bash
go run ./services/admin-api/cmd/server  # 运行服务
go build -o bin/admin-api ./services/admin-api/cmd/server  # 构建
go test ./...      # 运行测试
go mod tidy        # 整理依赖
go mod download    # 下载依赖
go fmt ./...       # 格式化代码
go vet ./...       # 代码检查
go clean           # 清理
```

### 📦 NPM 命令

```bash
cd web/admin
npm install        # 安装依赖
npm run dev        # 开发模式
npm run build      # 生产构建
npm run preview    # 预览构建结果
npm run lint       # 代码检查
npm run type-check # 类型检查
```

## 十、故障排查

### 🐛 常见问题

| 问题 | 解决方案 |
|------|---------|
| 端口被占用 | `lsof -ti :8083 \| xargs kill -9` |
| 数据库连接失败 | 检查配置文件和数据库状态 |
| Token 失效 | 重新登录获取新 Token |
| 前端 CORS 错误 | 检查后端 CORS 中间件配置 |
| npm 安装失败 | 清理缓存：`npm cache clean --force` |
| Go 模块找不到 | 运行：`go mod tidy` |

### 📞 获取帮助

- **详细文档**: [docs/INDEX.md](./docs/INDEX.md)
- **快速开始**: [GETTING_STARTED.md](./GETTING_STARTED.md)
- **API 文档**: http://localhost:8083/swagger/index.html
- **示例代码**: [docs/demo/](./docs/demo/)
- **GitHub Issues**: 提交问题

---

**GinForge - 让开发更加简单** 🚀

**提示**: 建议将本文档加入浏览器书签，方便随时查阅！
