# 🚀 GinForge 生产环境部署指南

## 📋 目录

- [部署架构](#部署架构)
- [快速部署](#快速部署)
- [详细配置](#详细配置)
- [安全加固](#安全加固)
- [监控运维](#监控运维)
- [故障排查](#故障排查)

---

## 🏗️ 部署架构

### 架构图

```
                    Internet
                       ↓
                 ┌──────────┐
                 │  Nginx   │ (80/443)
                 │  反向代理  │
                 └────┬─────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
     /api/*        /uploads/*      /*
        ↓             ↓             ↓
   ┌─────────┐  ┌─────────┐  ┌──────────┐
   │ Gateway │  │File-API │  │  静态文件  │
   │ (8080)  │  │ (8086)  │  │  (Vue3)   │
   └────┬────┘  └─────────┘  └──────────┘
        │
    ┌───┴────┬───────┬────────┐
    ↓        ↓       ↓        ↓
┌────────┐┌────────┐┌────────┐┌─────────┐
│user-api││merchant││admin   ││gateway  │
│ (8081) ││  (8082)││  (8083)││-worker  │
└───┬────┘└───┬────┘└───┬────┘└────┬────┘
    │         │         │          │
    └─────────┴─────────┴──────────┘
              ↓         ↓
        ┌──────────┐ ┌──────────┐
        │  MySQL   │ │  Redis   │
        │  (3306)  │ │  (6379)  │
        └──────────┘ └──────────┘
```

### 服务说明

| 服务 | 端口 | 对外 | 说明 |
|------|------|------|------|
| **Nginx** | 80/443 | ✅ | 唯一对外入口 |
| **Gateway** | 8080 | 🔒 | API 网关（内网）|
| **user-api** | 8081 | ❌ | 用户服务（内网）|
| **merchant-api** | 8082 | ❌ | 商户服务（内网）|
| **admin-api** | 8083 | ❌ | 管理后台（内网）|
| **file-api** | 8086 | ❌ | 文件服务（内网）|
| **MySQL** | 3306 | ❌ | 数据库（内网）|
| **Redis** | 6379 | ❌ | 缓存（内网）|

> 🔒 生产环境只需开放 **80/443** 端口

---

## 🚀 快速部署

### 步骤 1: 准备环境

```bash
# 克隆项目
git clone https://your-repo/GinForge.git
cd GinForge

# 检查 Docker
docker --version
docker-compose --version
```

### 步骤 2: 配置环境变量

```bash
cd deployments

# 复制环境变量配置
cp env.production.example .env.production

# 编辑配置（重要！）
vim .env.production
```

**必须修改的配置：**
- ✅ `MYSQL_PASSWORD` - 数据库密码
- ✅ `REDIS_PASSWORD` - Redis 密码
- ✅ `JWT_SECRET` - JWT 密钥（至少32位随机字符串）
- ✅ `CORS_ORIGINS` - 允许的域名

### 步骤 3: 构建前端

```bash
# 进入前端目录
cd ../web/admin

# 安装依赖
npm install

# 构建生产版本
npm run build

# 返回项目根目录
cd ../..
```

### 步骤 4: 启动服务

```bash
cd deployments

# 启动所有服务
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# 查看启动日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 步骤 5: 验证部署

```bash
# 等待所有服务健康
docker-compose -f docker-compose.prod.yml ps

# 健康检查
curl http://localhost/healthz

# 访问前端
# 浏览器打开: http://localhost
# 或: http://your-domain.com
```

---

## 🔧 详细配置

### 1. 数据库初始化

数据库会在首次启动时自动初始化：

```bash
# MySQL 容器启动时会自动执行:
# /docker-entrypoint-initdb.d/001_create_admin_tables.sql
```

如果需要手动初始化：

```bash
docker exec -i ginforge-mysql mysql -uginforge -p${MYSQL_PASSWORD} ginforge \
  < database/migrations/001_create_admin_tables.sql
```

### 2. SSL/HTTPS 配置

**生成自签名证书（测试用）：**

```bash
cd deployments
mkdir -p ssl

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/key.pem \
  -out ssl/cert.pem \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=GinForge/CN=yourdomain.com"
```

**使用 Let's Encrypt（生产环境）：**

```bash
# 使用 certbot 获取免费证书
docker run -it --rm \
  -v $(pwd)/ssl:/etc/letsencrypt \
  certbot/certbot certonly --standalone \
  -d yourdomain.com \
  -d www.yourdomain.com
```

然后取消注释 `nginx.conf` 中的 HTTPS 配置。

### 3. 配置文件说明

**configs/config.prod.yaml** - 生产环境主配置
- 使用 MySQL 而非 SQLite
- 启用 Redis
- 关闭 Debug 模式
- 严格的安全策略

**deployments/nginx.conf** - Nginx 配置
- 反向代理到 Gateway
- 静态资源服务
- Gzip 压缩
- 缓存策略

**deployments/redis.conf** - Redis 配置
- 持久化配置
- 安全加固
- 性能优化

---

## 🔐 安全加固

### 1. 修改默认密码

```bash
# 编辑 .env.production
JWT_SECRET=$(openssl rand -base64 32)
MYSQL_PASSWORD=$(openssl rand -base64 16)
REDIS_PASSWORD=$(openssl rand -base64 16)
```

### 2. 防火墙配置

```bash
# 只开放必要端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# 禁止直接访问内部服务端口
# 8080-8086 端口不对外开放
```

### 3. 限制容器权限

```yaml
# docker-compose.prod.yml 中添加
services:
  gateway:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
      - /app/logs
```

### 4. 网络隔离

```yaml
# 内部服务不暴露端口
admin-api:
  # 移除 ports 配置
  # 只在内网通过 Docker 网络访问
```

---

## 📊 监控运维

### 1. 健康检查

所有服务都配置了健康检查：

```bash
# 查看服务健康状态
docker-compose -f docker-compose.prod.yml ps

# Gateway 健康检查
curl http://localhost/healthz

# 各服务健康检查（容器内部）
docker exec ginforge-admin-api wget -q -O- http://localhost:8083/healthz
```

### 2. 日志管理

```bash
# 查看所有服务日志
docker-compose -f docker-compose.prod.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.prod.yml logs -f gateway

# 导出日志
docker-compose -f docker-compose.prod.yml logs --no-color > logs/all-services.log

# 清理旧日志
docker-compose -f docker-compose.prod.yml logs --tail=100
```

日志会自动轮转（最大 10MB，保留 3-5 个文件）

### 3. Prometheus 监控

所有服务暴露 `/metrics` 端点：

```bash
# Gateway 指标
curl http://localhost:8080/metrics

# Admin API 指标  
curl http://localhost:8083/metrics
```

可以添加 Prometheus + Grafana：

```yaml
# 在 docker-compose.prod.yml 中添加
prometheus:
  image: prom/prometheus:latest
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
    - prometheus_data:/prometheus
  ports:
    - "9090:9090"

grafana:
  image: grafana/grafana:latest
  ports:
    - "3001:3000"
  volumes:
    - grafana_data:/var/lib/grafana
```

### 4. 数据备份

**MySQL 备份：**

```bash
# 备份数据库
docker exec ginforge-mysql mysqldump \
  -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} \
  > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker exec -i ginforge-mysql mysql \
  -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} \
  < backup_20241013_120000.sql
```

**Redis 备份：**

```bash
# RDB 快照备份
docker exec ginforge-redis redis-cli -a ${REDIS_PASSWORD} SAVE
docker cp ginforge-redis:/data/dump.rdb ./backup/

# AOF 备份
docker cp ginforge-redis:/data/appendonly.aof ./backup/
```

---

## 🔄 更新部署

### 零停机更新

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建前端
cd web/admin
npm run build
cd ../..

# 3. 重新构建后端镜像
cd deployments
docker-compose -f docker-compose.prod.yml build

# 4. 滚动更新（一次更新一个服务）
docker-compose -f docker-compose.prod.yml up -d --no-deps --build gateway
docker-compose -f docker-compose.prod.yml up -d --no-deps --build admin-api
docker-compose -f docker-compose.prod.yml up -d --no-deps --build file-api
# ... 其他服务

# 5. 验证更新
curl http://localhost/healthz
```

### 蓝绿部署

使用 Docker Swarm 或 Kubernetes 实现蓝绿部署：

```bash
# 使用 Docker Stack
docker stack deploy -c docker-compose.prod.yml ginforge
```

---

## 🐛 故障排查

### 1. 服务无法启动

```bash
# 查看容器状态
docker-compose -f docker-compose.prod.yml ps

# 查看详细日志
docker-compose -f docker-compose.prod.yml logs <service-name>

# 检查容器内部
docker exec -it ginforge-<service> sh
```

### 2. 数据库连接失败

```bash
# 检查 MySQL 是否运行
docker exec ginforge-mysql mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} -e "SELECT 1"

# 检查网络连接
docker exec ginforge-admin-api ping mysql

# 查看数据库日志
docker logs ginforge-mysql
```

### 3. Redis 连接问题

```bash
# 测试 Redis 连接
docker exec ginforge-redis redis-cli -a ${REDIS_PASSWORD} PING

# 检查 Redis 配置
docker exec ginforge-redis redis-cli -a ${REDIS_PASSWORD} CONFIG GET requirepass
```

### 4. Nginx 502 错误

```bash
# 检查 Gateway 是否运行
docker exec ginforge-nginx ping gateway

# 查看 Nginx 错误日志
docker exec ginforge-nginx cat /var/log/nginx/error.log

# 重启 Nginx
docker-compose -f docker-compose.prod.yml restart nginx
```

---

## 📈 性能优化

### 1. 容器资源调整

根据实际负载调整 `docker-compose.prod.yml` 中的资源限制：

```yaml
deploy:
  resources:
    limits:
      cpus: '2.0'      # 根据服务器配置调整
      memory: 2G
    reservations:
      cpus: '1.0'
      memory: 1G
```

### 2. 数据库优化

```yaml
mysql:
  command:
    - --max_connections=2000
    - --innodb_buffer_pool_size=1G
    - --query_cache_size=64M
```

### 3. Redis 优化

修改 `redis.conf`：

```conf
maxmemory 512mb
maxmemory-policy allkeys-lru
```

### 4. Nginx 优化

```nginx
worker_processes auto;
worker_connections 2048;
keepalive_timeout 65;
client_max_body_size 100M;
```

---

## 🔄 扩展部署

### 水平扩展

```bash
# 扩展 Gateway 实例
docker-compose -f docker-compose.prod.yml up -d --scale gateway=3

# Nginx 自动负载均衡
# Gateway upstream 会自动发现多个实例
```

### 使用 Kubernetes

```bash
# 应用 K8s 配置
kubectl apply -f deployments/k8s/

# 使用 Istio 服务网格
kubectl apply -f deployments/k8s/istio/
```

---

## 📝 生产环境检查清单

### 部署前检查

- [ ] 修改所有默认密码
- [ ] 配置 JWT_SECRET
- [ ] 设置正确的 CORS 域名
- [ ] 配置 SSL 证书（HTTPS）
- [ ] 准备数据库初始化脚本
- [ ] 测试所有服务健康检查
- [ ] 配置日志轮转
- [ ] 设置资源限制
- [ ] 准备备份策略

### 部署后验证

- [ ] 所有容器状态为 healthy
- [ ] 可以访问前端页面
- [ ] API 接口正常响应
- [ ] 登录功能正常
- [ ] 文件上传功能正常
- [ ] 健康检查端点正常
- [ ] Prometheus 指标正常
- [ ] 日志正常输出

---

## 🆘 紧急处理

### 回滚到上一个版本

```bash
# 查看镜像历史
docker images | grep ginforge

# 回滚到特定版本
docker-compose -f docker-compose.prod.yml down
docker tag ginforge:backup ginforge:latest
docker-compose -f docker-compose.prod.yml up -d
```

### 紧急停止

```bash
# 停止所有服务
docker-compose -f docker-compose.prod.yml stop

# 紧急情况下强制停止
docker-compose -f docker-compose.prod.yml kill
```

### 数据恢复

```bash
# 从备份恢复 MySQL
docker exec -i ginforge-mysql mysql -uroot -p${MYSQL_ROOT_PASSWORD} \
  < backup/backup_20241013.sql

# 从备份恢复 Redis
docker cp backup/dump.rdb ginforge-redis:/data/
docker-compose -f docker-compose.prod.yml restart redis
```

---

## 📚 相关文档

- [部署指南](./DEPLOYMENT.md)
- [README](./README.md)
- [Nginx 配置说明](./nginx.conf)
- [Kubernetes 部署](./k8s/)

---

## 💡 最佳实践

1. **使用配置中心**：大规模部署时使用 Consul/Etcd
2. **服务网格**：使用 Istio 管理微服务通信
3. **CI/CD**：集成 Jenkins/GitLab CI 自动化部署
4. **监控告警**：Prometheus + Grafana + AlertManager
5. **日志聚合**：ELK Stack 或 Loki
6. **分布式追踪**：Jaeger 或 Zipkin
7. **负载均衡**：Nginx + Keepalived 高可用
8. **数据库高可用**：MySQL 主从复制或 Galera 集群
9. **Redis 集群**：Redis Sentinel 或 Cluster

---

## 🎯 总结

### 当前方案评分

| 项目 | 评分 | 说明 |
|------|------|------|
| **容器化** | ⭐⭐⭐⭐⭐ | 完整的 Docker 支持 |
| **安全性** | ⭐⭐⭐⭐ | 需要改默认密码 |
| **可扩展性** | ⭐⭐⭐⭐ | 支持水平扩展 |
| **监控** | ⭐⭐⭐⭐ | 有健康检查和指标 |
| **自动化** | ⭐⭐⭐⭐ | 一键部署 |
| **文档** | ⭐⭐⭐⭐⭐ | 完整的文档 |

**总体评价：⭐⭐⭐⭐ 4/5 星**

适合中小型生产环境部署。大规模部署建议迁移到 Kubernetes。

