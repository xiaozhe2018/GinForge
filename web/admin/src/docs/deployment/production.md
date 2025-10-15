# 生产部署

将 GinForge 部署到生产环境的完整指南。

## ✅ 部署检查清单

### 配置检查

- [ ] 修改 JWT Secret 为强随机字符串
- [ ] 使用 MySQL 而非 SQLite
- [ ] 启用 Redis
- [ ] 设置合理的数据库连接池
- [ ] 配置日志输出到文件
- [ ] 关闭 Debug 模式
- [ ] 设置合理的超时时间

### 安全检查

- [ ] 使用 HTTPS
- [ ] 配置防火墙规则
- [ ] 设置强密码策略
- [ ] 启用 CORS 白名单
- [ ] 配置限流规则
- [ ] 定期备份数据库

### 性能检查

- [ ] 优化数据库查询
- [ ] 启用 Redis 缓存
- [ ] 配置 CDN（静态资源）
- [ ] 启用 Gzip 压缩
- [ ] 设置合理的缓存策略

## 🔧 生产环境配置

### 1. 环境变量

```bash
# /etc/environment 或 .env
APP_ENV=production
APP_PORT=8083
APP_DEBUG=false

# 数据库（使用环境变量保护敏感信息）
DB_TYPE=mysql
DB_HOST=prod-db-host.example.com
DB_PORT=3306
DB_DATABASE=gin_forge_prod
DB_USERNAME=app_user
DB_PASSWORD=your-strong-password

# Redis
REDIS_ENABLED=true
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# JWT
JWT_SECRET=your-very-long-and-random-secret-key
JWT_EXPIRE_HOURS=720

# 日志
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=file
```

### 2. 配置文件

```yaml
# configs/config.prod.yaml
app:
  env: "production"
  debug: false
  port: 8083
  read_timeout: 60s
  write_timeout: 60s

database:
  type: "mysql"
  host: "${DB_HOST}"
  port: 3306
  database: "${DB_DATABASE}"
  username: "${DB_USERNAME}"
  password: "${DB_PASSWORD}"
  max_idle_conns: 20
  max_open_conns: 200
  conn_max_lifetime: 3600

redis:
  enabled: true
  host: "${REDIS_HOST}"
  password: "${REDIS_PASSWORD}"
  pool_size: 50
  min_idle_conns: 10

log:
  level: "info"
  format: "json"
  output: "file"
  file_path: "/var/log/ginforge/admin-api.log"
```

## 🏗️ 部署方式

### 方式一：直接部署

```bash
# 1. 构建二进制文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o bin/admin-api \
  ./services/admin-api/cmd/server/main.go

# 2. 上传到服务器
scp bin/admin-api user@server:/opt/ginforge/

# 3. 在服务器上运行
ssh user@server
cd /opt/ginforge
./admin-api
```

### 方式二：使用 Systemd

创建服务文件 `/etc/systemd/system/ginforge-admin.service`：

```ini
[Unit]
Description=GinForge Admin API
After=network.target mysql.service redis.service

[Service]
Type=simple
User=ginforge
Group=ginforge
WorkingDirectory=/opt/ginforge
ExecStart=/opt/ginforge/bin/admin-api
Restart=always
RestartSec=10
StandardOutput=append:/var/log/ginforge/admin-api.log
StandardError=append:/var/log/ginforge/admin-api-error.log

# 环境变量
Environment="APP_ENV=production"
Environment="DB_PASSWORD=your-password"
Environment="JWT_SECRET=your-secret"

[Install]
WantedBy=multi-user.target
```

管理服务：

```bash
# 启动服务
sudo systemctl start ginforge-admin

# 开机自启
sudo systemctl enable ginforge-admin

# 查看状态
sudo systemctl status ginforge-admin

# 重启服务
sudo systemctl restart ginforge-admin

# 查看日志
sudo journalctl -u ginforge-admin -f
```

### 方式三：使用 Docker（推荐）

```bash
# 1. 构建镜像
docker build -t ginforge/admin-api:latest -f deployments/docker/Dockerfile .

# 2. 运行容器
docker run -d \
  --name ginforge-admin \
  -p 8083:8083 \
  -e APP_ENV=production \
  -e DB_HOST=mysql \
  -e DB_PASSWORD=password \
  -e REDIS_HOST=redis \
  -e JWT_SECRET=secret \
  --network ginforge-network \
  ginforge/admin-api:latest

# 3. 查看日志
docker logs -f ginforge-admin
```

## 🔒 安全配置

### 1. HTTPS 配置

使用 Nginx 反向代理：

```nginx
# /etc/nginx/sites-available/ginforge
server {
    listen 443 ssl http2;
    server_name api.ginforge.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    location / {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. 防火墙配置

```bash
# 只开放必要的端口
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 22/tcp    # SSH
sudo ufw enable
```

### 3. 限制数据库访问

```sql
-- 只允许应用服务器访问
CREATE USER 'app_user'@'app-server-ip' IDENTIFIED BY 'strong_password';
GRANT ALL PRIVILEGES ON gin_forge.* TO 'app_user'@'app-server-ip';
FLUSH PRIVILEGES;
```

## 📊 监控和日志

### 1. 日志管理

```bash
# 配置 logrotate
cat > /etc/logrotate.d/ginforge << EOF
/var/log/ginforge/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 ginforge ginforge
    sharedscripts
    postrotate
        systemctl reload ginforge-admin
    endscript
}
EOF
```

### 2. 性能监控

使用 Prometheus：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'ginforge-admin'
    static_configs:
      - targets: ['localhost:8083']
    metrics_path: '/metrics'
```

### 3. 健康检查

```bash
# 定期检查服务健康
*/5 * * * * curl -f http://localhost:8083/api/v1/admin/system/health || systemctl restart ginforge-admin
```

## 💾 数据备份

### MySQL 备份

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR=/opt/backups/mysql
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份
mysqldump -h localhost -u root -p123456 gin_forge > $BACKUP_DIR/backup_$DATE.sql

# 压缩备份
gzip $BACKUP_DIR/backup_$DATE.sql

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

设置定时任务：

```bash
# 每天凌晨 2 点备份
0 2 * * * /opt/scripts/backup.sh
```

## 🚀 性能优化

### 1. 编译优化

```bash
# 减小二进制文件大小
go build -ldflags="-w -s" -o bin/admin-api ./services/admin-api/cmd/server/main.go

# 使用 UPX 进一步压缩
upx --best bin/admin-api
```

### 2. 数据库优化

```sql
-- 添加索引
CREATE INDEX idx_username ON users(username);
CREATE INDEX idx_status_created ON users(status, created_at);

-- 查询慢查询
SELECT * FROM mysql.slow_log;

-- 分析查询
EXPLAIN SELECT * FROM users WHERE username = 'admin';
```

### 3. Redis 优化

```redis
# 设置最大内存
config set maxmemory 2gb

# 设置淘汰策略
config set maxmemory-policy allkeys-lru
```

## 📚 完整部署示例

查看详细部署文档：

- **Docker 部署**: `docs/DEPLOYMENT.md`
- **Kubernetes 部署**: `deployments/k8s/`
- **Nginx 配置**: `deployments/nginx.conf`

## 🎯 下一步

- [Docker 部署](./docker) - 容器化部署
- [监控和告警](../best-practices/monitoring) - 系统监控

---

**提示**: 生产环境部署要特别注意安全性和稳定性，建议先在测试环境验证！

