# GinForge 部署指南

## 架构说明

GinForge 采用微服务架构，通过 **API Gateway** 统一对外提供服务：

```
外部请求 → Nginx (80) → Gateway (8080) → 内部微服务
                        ├→ WebSocket Gateway (8087) (WebSocket 连接)
                        │                        ├→ user-api (8081)
                        │                        ├→ merchant-api (8082)
                        │                        ├→ admin-api (8083)
                        │                        ├→ gateway-worker (8084)
                        │                        ├→ demo (8085)
                        │                        └→ file-api (8086)
                        │
                        └─────── Redis PubSub ─────┘
                                 (消息总线)
```

**重要提示：**
- 外部只需暴露 **Nginx (80)** 端口
- **Gateway (8080)** - HTTP API 路由
- **WebSocket Gateway (8087)** - WebSocket 实时通信
- 内部服务端口无需对外暴露（仅内网访问）
- Nginx 反向代理 HTTP 请求到 Gateway，WebSocket 请求到 WebSocket Gateway

## 1. 本地开发部署

### 1.1 快速启动

```bash
# 方式一：使用启动脚本（推荐）
./scripts/start-services.sh

# 方式二：使用Makefile
make build && make run

# 方式三：手动启动
go run ./services/gateway/cmd/server &
go run ./services/admin-api/cmd/server &
go run ./services/file-api/cmd/server &
# ... 其他服务
```

### 1.2 停止服务

```bash
# 方式一：使用停止脚本
./scripts/stop-services.sh

# 方式二：使用Makefile
make stop
```

### 1.3 服务端口说明

| 服务 | 端口 | 说明 |
|------|------|------|
| Gateway | 8080 | HTTP API 网关 |
| User API | 8081 | 用户端服务 |
| Merchant API | 8082 | 商户端服务 |
| Admin API | 8083 | 管理后台服务 |
| Gateway Worker | 8084 | 异步任务处理 |
| Demo | 8085 | 演示服务 |
| File API | 8086 | 文件服务 |
| **WebSocket Gateway** | **8087** | **WebSocket 实时通信** |

### 1.4 访问服务

```bash
# 通过 Gateway 访问（生产环境方式）
curl http://localhost:8080/api/v1/admin/users

# 直接访问服务（仅开发调试）
curl http://localhost:8083/api/v1/admin/users
curl http://localhost:8086/api/v1/files/statistics
```

## 2. Docker 部署（推荐）

### 2.1 使用 Docker Compose

```bash
# 启动所有服务
docker-compose -f deployments/docker-compose.yml up -d

# 查看服务状态
docker-compose -f deployments/docker-compose.yml ps

# 查看日志
docker-compose -f deployments/docker-compose.yml logs -f gateway
docker-compose -f deployments/docker-compose.yml logs -f file-api

# 停止服务
docker-compose -f deployments/docker-compose.yml down
```

### 2.2 服务访问

Docker Compose 启动后：

**外部访问（通过Nginx）：**
```bash
# 前端页面
http://localhost

# API请求（通过Gateway）
curl http://localhost/api/v1/admin/users
curl http://localhost/api/v1/files/upload
```

**内部服务通信：**
```
gateway → http://admin-api:8083
gateway → http://file-api:8086
gateway → http://user-api:8081
```

### 2.3 端口映射

Docker Compose 端口映射：

| 服务 | 容器端口 | 宿主机端口 | 是否对外 |
|------|----------|-----------|---------|
| nginx | 80 | 80 | ✅ 是 |
| gateway | 8080 | 8080 | ✅ 是（开发用） |
| admin-api | 8083 | - | ❌ 否（内网） |
| file-api | 8086 | - | ❌ 否（内网） |
| mysql | 3306 | 3306 | 🔒 可选 |
| redis | 6379 | 6379 | 🔒 可选 |

**说明：**
- 生产环境只需暴露 80(Nginx) 和 8080(Gateway)
- 内部服务通过Docker网络通信，无需暴露端口
- MySQL/Redis端口可选暴露（用于外部管理工具）

## 3. 生产环境部署

### 3.1 服务器部署（二进制）

```bash
# 1. 编译所有服务
make build

# 2. 部署到服务器
scp -r bin/ user@server:/opt/ginforge/
scp -r configs/ user@server:/opt/ginforge/

# 3. 启动服务（使用systemd或supervisor）
# 只需要暴露Gateway给外部
ssh user@server "cd /opt/ginforge && nohup ./bin/gateway > logs/gateway.log 2>&1 &"
ssh user@server "cd /opt/ginforge && nohup ./bin/file-api > logs/file-api.log 2>&1 &"
# ... 其他内部服务
```

### 3.2 Nginx 配置

Nginx只需反向代理到Gateway：

```nginx
upstream ginforge_gateway {
    server localhost:8080;
    # 如果有多个Gateway实例
    # server localhost:8080;
    # server localhost:8081;
}

server {
    listen 80;
    server_name your-domain.com;

    # API请求转发到Gateway
    location /api/ {
        proxy_pass http://ginforge_gateway;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 前端静态文件
    location / {
        root /var/www/ginforge/dist;
        try_files $uri $uri/ /index.html;
    }

    # 文件访问（如果需要直接访问上传的文件）
    location /uploads/ {
        alias /opt/ginforge/uploads/;
        expires 7d;
        add_header Cache-Control "public, immutable";
    }
}
```

## 4. Kubernetes 部署

### 4.1 服务暴露策略

**对外服务（LoadBalancer/Ingress）：**
- gateway - API网关（唯一对外入口）

**内部服务（ClusterIP）：**
- admin-api
- user-api
- merchant-api
- file-api
- gateway-worker
- demo

### 4.2 部署配置

```bash
# 应用K8s配置
kubectl apply -f deployments/k8s/

# Gateway通过Ingress对外暴露
kubectl apply -f deployments/k8s/ingress.yaml
```

示例Ingress配置：
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ginforge-ingress
spec:
  rules:
  - host: api.your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gateway
            port:
              number: 8080
```

## 5. 服务通信

### 5.1 Gateway路由配置

Gateway需要配置到各个内部服务的路由：

```yaml
# configs/config.yaml
services:
  user_api:
    addr: "http://user-api:8081"      # Docker/K8s内网地址
  merchant_api:
    addr: "http://merchant-api:8082"
  admin_api:
    addr: "http://admin-api:8083"
  file_api:
    addr: "http://file-api:8086"      # 文件服务
```

### 5.2 请求路由

Gateway根据路径前缀转发请求：

```
/api/v1/users/*     → user-api:8081
/api/v1/merchants/* → merchant-api:8082
/api/v1/admin/*     → admin-api:8083
/api/v1/files/*     → file-api:8086
```

## 6. 扩展性

### 6.1 水平扩展

**扩展Gateway：**
```bash
# Docker Compose
docker-compose up -d --scale gateway=3

# Kubernetes
kubectl scale deployment gateway --replicas=3
```

**扩展文件服务：**
```bash
# Docker Compose
docker-compose up -d --scale file-api=2

# Kubernetes
kubectl scale deployment file-api --replicas=2
```

### 6.2 负载均衡

- Nginx负载均衡到多个Gateway实例
- Gateway负载均衡到多个后端服务实例
- 使用Istio/Envoy实现更高级的流量管理

## 7. 监控和日志

### 7.1 集中式日志

所有服务通过Gateway的Request ID进行链路追踪：

```
Request → Gateway (生成request_id) → file-api (携带request_id)
```

### 7.2 Prometheus监控

```bash
# Gateway指标
curl http://localhost:8080/metrics

# 文件服务指标
curl http://localhost:8086/metrics
```

## 8. 安全配置

### 8.1 防火墙规则

生产环境防火墙配置：

```bash
# 只开放必要端口
firewall-cmd --add-port=80/tcp --permanent    # Nginx
firewall-cmd --add-port=443/tcp --permanent   # HTTPS
# 内部服务端口不对外开放
```

### 8.2 网络隔离

Docker网络配置：
```yaml
networks:
  ginforge-network:
    driver: bridge
    internal: false  # Gateway需要对外
```

## 9. 常用命令

```bash
# 编译
make build

# 启动所有服务
./scripts/start-services.sh

# 停止所有服务
./scripts/stop-services.sh

# Docker部署
make compose

# 停止Docker服务
make compose-down

# 查看服务状态
make status

# 生成API文档
make swagger
```

## 10. 故障排查

### 问题1：无法访问服务

检查Gateway是否正常运行：
```bash
curl http://localhost:8080/healthz
```

检查内部服务是否正常：
```bash
curl http://localhost:8086/healthz  # 开发环境
```

### 问题2：文件上传失败

通过Gateway上传：
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@test.jpg"
```

如果失败，检查：
1. Gateway是否正确路由到file-api
2. file-api服务是否正常运行
3. 上传目录权限是否正确

### 问题3：端口冲突

```bash
# 查看端口占用
lsof -i :8080  # Gateway
lsof -i :8086  # file-api

# 停止冲突进程
./scripts/stop-services.sh
```

## 总结

**核心要点：**
1. ✅ Gateway是唯一对外入口
2. ✅ Nginx反向代理到Gateway
3. ✅ 内部服务通过Gateway访问
4. ✅ 内部服务端口无需对外暴露
5. ✅ 统一的二进制文件管理（bin/目录）

## 8. WebSocket 实时通信部署

### 8.1 架构说明

WebSocket Gateway 是独立的实时通信服务：

```
客户端
  │
  ├─ HTTP 请求 → Gateway (8080) → 各个 API 服务
  │
  └─ WebSocket → WebSocket Gateway (8087)
                      │
                      ↓
                 Redis PubSub ← 其他服务发布消息
                      │
                      ↓
                 推送给 WebSocket 客户端
```

**设计优势**：
- ✅ 职责分离：HTTP 和 WebSocket 独立扩展
- ✅ 故障隔离：WebSocket 崩溃不影响 HTTP
- ✅ 独立扩展：可针对性扩展 WebSocket 服务
- ✅ 易于维护：代码职责清晰

### 8.2 开发环境

```bash
# 1. 确保 Redis 运行
docker run -d -p 6379:6379 redis:alpine

# 2. 启动 WebSocket Gateway
go run ./services/websocket-gateway/cmd/server

# 3. 测试连接（需要 JWT Token）
ws://localhost:8087/ws?token=YOUR_JWT_TOKEN

# 4. 查看统计
curl http://localhost:8087/ws/stats
```

### 8.3 生产环境 Nginx 配置

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    # SSL 配置
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # HTTP API → Gateway
    location /api/ {
        proxy_pass http://gateway:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # WebSocket → WebSocket Gateway
    location /ws {
        proxy_pass http://websocket-gateway:8087;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 86400;  # 24小时超时
    }
}
```

### 8.4 前端使用

```typescript
// 连接 WebSocket
const ws = new WebSocket('wss://yourdomain.com/ws?token=' + jwtToken)

// 接收消息
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data)
  if (msg.type === 'notification') {
    showNotification(msg.content)
  }
}

// 发送消息
ws.send(JSON.stringify({
  type: 'chat',
  content: { text: 'Hello' }
}))
```

### 8.5 后端发送 WebSocket 消息

```go
import "goweb/pkg/notification"

// 在任何服务中使用
notifyClient := notification.NewClient(redisClient)

// 发送通知给用户
notifyClient.SendNotification(ctx, userID, &websocket.NotificationMessage{
    Title: "订单通知",
    Body:  "您的订单已发货",
    Icon:  "Truck",
    Link:  "/orders/12345",
})

// 广播给所有在线用户
notifyClient.BroadcastNotification(ctx, notification)
```

### 8.6 Docker Compose 部署

WebSocket Gateway 已包含在 `docker-compose.yml` 中：

```yaml
websocket-gateway:
  build:
    context: ..
    dockerfile: deployments/docker/Dockerfile
  command: ["./bin/websocket-gateway"]
  ports:
    - "8087:8087"
  environment:
    - REDIS_ENABLED=true
    - REDIS_HOST=redis
    - REDIS_PORT=6379
  depends_on:
    - redis
```

启动：
```bash
docker-compose up -d
```

### 8.7 监控和调试

```bash
# 查看 WebSocket 统计
curl http://localhost:8087/ws/stats

# 查看在线用户（需要 JWT）
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8087/ws/online-users

# 查看健康状态
curl http://localhost:8087/healthz
```

### 8.8 性能优化

**单实例性能**：
- 支持 10,000+ 并发连接
- 消息延迟 < 50ms
- 内存占用：~100MB (1万连接)

**多实例部署**（负载均衡）：
```yaml
# docker-compose.yml
websocket-gateway:
  deploy:
    replicas: 3  # 3个实例
```

**注意**：多实例必须使用 Redis，消息才能同步到所有实例。

