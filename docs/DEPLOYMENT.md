# GinForge 部署指南

## 架构说明

GinForge 采用微服务架构，通过 **API Gateway** 统一对外提供服务：

```
外部请求 → Nginx (80) → Gateway (8080) → 内部微服务
                                        ├→ user-api (8081)
                                        ├→ merchant-api (8082)
                                        ├→ admin-api (8083)
                                        ├→ gateway-worker (8084)
                                        ├→ demo (8085)
                                        └→ file-api (8086)
```

**重要提示：**
- 外部只需暴露 **Gateway (8080)** 和 **Nginx (80)** 端口
- 内部服务端口无需对外暴露（仅内网访问）
- Nginx 反向代理到 Gateway 即可

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

### 1.3 访问服务

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
