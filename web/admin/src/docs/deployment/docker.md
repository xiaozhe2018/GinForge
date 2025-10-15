# Docker 部署

使用 Docker 容器化部署 GinForge 框架。

## 🐳 为什么使用 Docker？

- ✅ **环境一致**：开发、测试、生产环境完全一致
- ✅ **快速部署**：一键启动所有服务
- ✅ **资源隔离**：每个服务独立运行
- ✅ **易于扩展**：水平扩展很简单
- ✅ **版本管理**：通过镜像tag管理版本

## 🚀 快速开始

### 使用 docker-compose

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

## 📁 Docker 配置文件

### docker-compose.yml

```yaml
version: '3.8'

services:
  # MySQL 数据库
  mysql:
    image: mysql:8.0
    container_name: ginforge-mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: gin_forge
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./database/migrations:/docker-entrypoint-initdb.d
    networks:
      - ginforge-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: ginforge-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - ginforge-network
    command: redis-server --appendonly yes

  # Admin API
  admin-api:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
      args:
        SERVICE: admin-api
    container_name: ginforge-admin-api
    ports:
      - "8083:8083"
    environment:
      APP_ENV: production
      DB_TYPE: mysql
      DB_HOST: mysql
      DB_PASSWORD: 123456
      REDIS_ENABLED: "true"
      REDIS_HOST: redis
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - mysql
      - redis
    networks:
      - ginforge-network
    restart: unless-stopped

  # WebSocket Gateway
  websocket-gateway:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
      args:
        SERVICE: websocket-gateway
    container_name: ginforge-websocket
    ports:
      - "8087:8087"
    environment:
      REDIS_HOST: redis
    depends_on:
      - redis
    networks:
      - ginforge-network
    restart: unless-stopped

networks:
  ginforge-network:
    driver: bridge

volumes:
  mysql_data:
  redis_data:
```

## 📦 Dockerfile

### 多阶段构建

```dockerfile
# deployments/docker/Dockerfile

# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建参数
ARG SERVICE=admin-api

# 编译二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/${SERVICE} \
    ./services/${SERVICE}/cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
ARG SERVICE=admin-api
COPY --from=builder /app/bin/${SERVICE} /app/app
COPY --from=builder /app/configs /app/configs

# 暴露端口
EXPOSE 8083

# 运行服务
CMD ["/app/app"]
```

## 🔧 构建和运行

### 构建镜像

```bash
# 构建 admin-api
docker build \
  -t ginforge/admin-api:1.0.0 \
  -f deployments/docker/Dockerfile \
  --build-arg SERVICE=admin-api \
  .

# 构建 websocket-gateway
docker build \
  -t ginforge/websocket-gateway:1.0.0 \
  -f deployments/docker/Dockerfile \
  --build-arg SERVICE=websocket-gateway \
  .
```

### 运行容器

```bash
# 运行 admin-api
docker run -d \
  --name ginforge-admin \
  -p 8083:8083 \
  -e APP_ENV=production \
  -e DB_HOST=mysql \
  -e REDIS_HOST=redis \
  --network ginforge-network \
  ginforge/admin-api:1.0.0

# 查看日志
docker logs -f ginforge-admin

# 进入容器
docker exec -it ginforge-admin sh
```

## 📊 监控和维护

### 查看资源使用

```bash
# 查看容器资源使用情况
docker stats

# 查看特定容器
docker stats ginforge-admin
```

### 容器健康检查

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8083/health || exit 1
```

### 日志管理

```bash
# 查看日志
docker logs ginforge-admin

# 实时跟踪日志
docker logs -f ginforge-admin

# 查看最近 100 行
docker logs --tail 100 ginforge-admin

# 导出日志
docker logs ginforge-admin > admin-api.log 2>&1
```

## 🔄 更新和回滚

### 更新服务

```bash
# 1. 构建新版本
docker build -t ginforge/admin-api:1.1.0 .

# 2. 停止旧容器
docker stop ginforge-admin
docker rm ginforge-admin

# 3. 启动新容器
docker run -d --name ginforge-admin ginforge/admin-api:1.1.0

# 或使用 docker-compose
docker-compose pull
docker-compose up -d
```

### 回滚

```bash
# 回滚到旧版本
docker stop ginforge-admin
docker rm ginforge-admin
docker run -d --name ginforge-admin ginforge/admin-api:1.0.0
```

## 💾 数据持久化

### 挂载卷

```yaml
volumes:
  # 数据库数据
  - mysql_data:/var/lib/mysql
  
  # Redis 数据
  - redis_data:/data
  
  # 日志文件
  - ./logs:/var/log/ginforge
  
  # 上传文件
  - ./uploads:/app/uploads
  
  # 配置文件
  - ./configs:/app/configs:ro  # 只读
```

## 🌐 Nginx 反向代理

### 配置示例

```nginx
upstream admin_api {
    server 127.0.0.1:8083;
    # 多实例负载均衡
    # server 127.0.0.1:8084;
}

server {
    listen 80;
    server_name api.ginforge.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.ginforge.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # API 代理
    location /api/ {
        proxy_pass http://admin_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # 静态文件
    location /uploads/ {
        alias /opt/ginforge/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

## 📚 完整示例

查看详细配置：

- **Docker Compose**: `deployments/docker-compose.yml`
- **Dockerfile**: `deployments/docker/Dockerfile`
- **Nginx 配置**: `deployments/nginx.conf`
- **部署文档**: `docs/DEPLOYMENT.md`

## 🎯 下一步

- [生产部署](./production) - 生产环境配置
- [Kubernetes 部署](../advanced/kubernetes) - K8s 集群部署

---

**提示**: Docker 极大简化了部署流程，强烈推荐用于生产环境！

