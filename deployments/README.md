# GinForge Docker 部署指南

## 📦 部署架构

```
┌─────────────┐
│  客户端/浏览器 │
└──────┬──────┘
       │ HTTP/HTTPS
       ↓
┌─────────────────────────┐
│  Nginx (80端口)          │  ← 反向代理 + 静态资源服务器
│  ├─ /api/*  → Gateway   │
│  └─ /*      → 前端静态    │
└──────────┬──────────────┘
           │
           ↓
┌─────────────────────────┐
│  Gateway 网关 (8080)     │  ← API 网关/路由转发
│  ├─ /api/v1/user/*      │
│  ├─ /api/v1/merchant/*  │
│  └─ /api/v1/admin/*     │
└──────────┬──────────────┘
           │
    ┌──────┴──────┬──────────┐
    ↓             ↓          ↓
┌─────────┐  ┌─────────┐  ┌─────────┐
│user-api │  │merchant │  │admin-api│
│ (8081)  │  │  (8082) │  │ (8083)  │
└─────────┘  └─────────┘  └─────────┘
```

## 🚀 快速部署

### 1. 前置要求

- Docker 20.10+
- Docker Compose 1.29+

### 2. 一键启动

```bash
# 进入部署目录
cd deployments

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 3. 访问服务

| 服务 | 地址 | 说明 |
|-----|------|------|
| 🌐 **前端管理后台** | http://localhost | Nginx 反向代理 |
| 🚪 **API 网关** | http://localhost/api | 所有 API 请求入口 |
| 📚 **API 文档** | http://localhost/swagger/index.html | Swagger 文档 |
| 💓 **健康检查** | http://localhost/healthz | 服务健康状态 |

### 4. 直接访问后端服务（调试用）

| 服务 | 端口 | 直接访问地址 |
|-----|------|-------------|
| Gateway | 8080 | http://localhost:8080 |
| user-api | 8081 | http://localhost:8081 |
| merchant-api | 8082 | http://localhost:8082 |
| admin-api | 8083 | http://localhost:8083 |

## 🔧 Nginx 反向代理配置

### 核心配置说明

```nginx
# 1. API 请求代理到 Gateway
location /api/ {
    proxy_pass http://gateway:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

# 2. 前端静态资源
location / {
    root /usr/share/nginx/html/admin;
    try_files $uri $uri/ /index.html;
}
```

### 工作流程

1. **前端请求 API**：
   ```
   浏览器 → http://localhost/api/v1/admin/users
   ```

2. **Nginx 接收并转发**：
   ```
   Nginx → http://gateway:8080/api/v1/admin/users
   ```

3. **Gateway 路由到后端**：
   ```
   Gateway → http://admin-api:8080/api/v1/admin/users
   ```

4. **响应返回**：
   ```
   admin-api → Gateway → Nginx → 浏览器
   ```

## 🔍 常用命令

```bash
# 查看所有容器
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f gateway
docker-compose logs -f nginx

# 重启所有服务
docker-compose restart

# 重启单个服务
docker-compose restart gateway

# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v

# 重新构建镜像
docker-compose build --no-cache

# 进入容器
docker-compose exec gateway sh
docker-compose exec nginx sh
```

## 🐛 故障排查

### 1. Nginx 启动失败

**症状**：`nginx` 容器启动失败

**检查**：
```bash
# 查看 Nginx 错误日志
docker-compose logs nginx

# 验证配置文件
docker run --rm -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro nginx:alpine nginx -t
```

**常见问题**：
- 配置文件路径错误
- 配置语法错误
- 端口已被占用

### 2. Gateway 无法连接后端服务

**症状**：API 返回 502 Bad Gateway

**检查**：
```bash
# 检查服务状态
docker-compose ps

# 检查网络连通性
docker-compose exec gateway ping user-api
docker-compose exec gateway ping merchant-api
docker-compose exec gateway ping admin-api
```

**解决方案**：
- 确保所有后端服务都在运行
- 检查 docker-compose.yml 中的网络配置
- 检查环境变量配置

### 3. 前端无法访问

**症状**：访问 http://localhost 显示 404

**检查**：
```bash
# 检查静态文件挂载
docker-compose exec nginx ls -la /usr/share/nginx/html/admin

# 检查 Nginx 配置
docker-compose exec nginx cat /etc/nginx/nginx.conf
```

### 4. CORS 跨域问题

**症状**：浏览器控制台显示 CORS 错误

**解决方案**：
- Gateway 已配置 CORS，检查 `services/gateway/internal/router/router.go`
- 确保前端请求通过 Nginx 代理（相同域名）

## 🔐 生产环境配置

### 1. 启用 HTTPS

在 `nginx.conf` 中取消注释 HTTPS 配置：

```nginx
server {
    listen       443 ssl http2;
    server_name  your-domain.com;
    
    ssl_certificate      /etc/nginx/ssl/cert.pem;
    ssl_certificate_key  /etc/nginx/ssl/key.pem;
    
    # ... 其他配置
}
```

### 2. 环境变量配置

修改 `docker-compose.yml`：

```yaml
environment:
  - DB_DRIVER=mysql
  - DB_DSN=user:pass@tcp(mysql:3306)/dbname
  - REDIS_ADDR=redis:6379
  - JWT_SECRET=your-production-secret-key
```

### 3. 添加数据库和缓存

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your-password
      MYSQL_DATABASE: ginforge
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - goweb-network

  redis:
    image: redis:alpine
    volumes:
      - redis_data:/data
    networks:
      - goweb-network
```

### 4. 日志管理

```yaml
services:
  gateway:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 📊 监控和维护

### 健康检查

```bash
# Gateway 健康检查
curl http://localhost/healthz

# 各服务健康检查
curl http://localhost:8080/healthz  # Gateway
curl http://localhost:8081/healthz  # user-api
curl http://localhost:8082/healthz  # merchant-api
curl http://localhost:8083/healthz  # admin-api
```

### 性能监控

添加 Prometheus 和 Grafana：

```yaml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
```

## 🔄 更新部署

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建镜像
docker-compose build

# 3. 滚动更新（零停机）
docker-compose up -d --no-deps --build gateway

# 4. 验证更新
docker-compose ps
curl http://localhost/healthz
```

## 📝 注意事项

1. **数据持久化**：生产环境建议使用外部数据库和 Redis
2. **安全配置**：修改默认密钥和密码
3. **资源限制**：为容器设置内存和 CPU 限制
4. **备份策略**：定期备份数据卷
5. **日志轮转**：配置日志大小和保留策略

## 🆘 获取帮助

- 查看项目文档：`docs/` 目录
- 提交 Issue：GitHub Issues
- 查看日志：`docker-compose logs -f`

