# 配置选项

GinForge 框架的完整配置项参考。

## 📋 应用配置 (`app`)

```yaml
app:
  name: "GinForge"          # 应用名称
  version: "1.0.0"          # 应用版本
  env: "development"        # 环境：development, test, production
  debug: true               # 调试模式
  port: 8080                # HTTP 端口
  read_timeout: 60s         # 读取超时
  write_timeout: 60s        # 写入超时
  idle_timeout: 120s        # 空闲超时
```

## 🗄️ 数据库配置 (`database`)

```yaml
database:
  type: "mysql"                 # 数据库类型：sqlite, mysql, postgres
  host: "localhost"             # 主机地址
  port: 3306                    # 端口号
  database: "gin_forge"         # 数据库名
  username: "root"              # 用户名
  password: "password"          # 密码
  charset: "utf8mb4"            # 字符集
  max_idle_conns: 10            # 最大空闲连接数
  max_open_conns: 100           # 最大打开连接数
  conn_max_lifetime: 3600       # 连接最大生命周期（秒）
  log_level: "error"            # 日志级别：silent, error, warn, info
```

## 📮 Redis 配置 (`redis`)

```yaml
redis:
  enabled: true            # 是否启用 Redis
  host: "localhost"        # 主机地址
  port: 6379               # 端口号
  password: ""             # 密码
  db: 0                    # 数据库编号 (0-15)
  pool_size: 10            # 连接池大小
  min_idle_conns: 5        # 最小空闲连接数
  dial_timeout: 5s         # 连接超时
  read_timeout: 3s         # 读取超时
  write_timeout: 3s        # 写入超时
  pool_timeout: 4s         # 连接池超时
```

## 🔑 JWT 配置 (`jwt`)

```yaml
jwt:
  secret: "your-secret-key"  # 密钥（生产环境必须修改）
  expire_hours: 24           # Token 过期时间（小时）
```

## 📝 日志配置 (`log`)

```yaml
log:
  level: "info"              # 日志级别：debug, info, warn, error, fatal
  format: "json"             # 日志格式：json, text
  output: "stdout"           # 输出：stdout, file
  file_path: "logs/app.log"  # 日志文件路径（output=file 时）
```

## 🌐 服务端口配置 (`services`)

```yaml
services:
  user_api:
    port: 8081
    name: "user-api"
  merchant_api:
    port: 8082
    name: "merchant-api"
  admin_api:
    port: 8083
    name: "admin-api"
  gateway:
    port: 8080
    name: "gateway"
  websocket_gateway:
    port: 8087
    name: "websocket-gateway"
  file_api:
    port: 8086
    name: "file-api"
```

## 📦 文件存储配置 (`storage`)

```yaml
storage:
  type: "local"               # 存储类型：local, oss, s3, minio
  max_file_size: 104857600    # 最大文件大小（字节），100MB
  
  # 本地存储
  local:
    base_path: "./uploads"
    url_prefix: "http://localhost:8086/uploads"
  
  # 阿里云 OSS
  oss:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    access_key_id: "${OSS_ACCESS_KEY}"
    access_key_secret: "${OSS_ACCESS_SECRET}"
    bucket_name: "my-bucket"
```

## 🎯 完整配置示例

### 开发环境 (`configs/config.yaml`)

```yaml
app:
  env: "development"
  debug: true
  port: 8083

database:
  type: "sqlite"
  database: "goweb.db"

redis:
  enabled: false

log:
  level: "debug"
  format: "text"
  output: "stdout"
```

### 生产环境 (`configs/config.prod.yaml`)

```yaml
app:
  env: "production"
  debug: false
  port: 8083

database:
  type: "mysql"
  host: "prod-db.example.com"
  port: 3306
  database: "gin_forge_prod"
  username: "app_user"
  password: "${DB_PASSWORD}"
  max_open_conns: 200

redis:
  enabled: true
  host: "prod-redis.example.com"
  password: "${REDIS_PASSWORD}"
  pool_size: 50

jwt:
  secret: "${JWT_SECRET}"
  expire_hours: 720  # 30天

log:
  level: "info"
  format: "json"
  output: "file"
  file_path: "/var/log/ginforge/app.log"
```

## 🌍 环境变量映射

| 配置项 | 环境变量 | 示例 |
|--------|----------|------|
| `app.port` | `APP_PORT` | `8083` |
| `database.host` | `DB_HOST` | `localhost` |
| `database.password` | `DB_PASSWORD` | `secret` |
| `redis.host` | `REDIS_HOST` | `localhost` |
| `redis.password` | `REDIS_PASSWORD` | `secret` |
| `jwt.secret` | `JWT_SECRET` | `random-key` |
| `log.level` | `LOG_LEVEL` | `info` |

## 📚 相关文档

- [配置系统](../core-concepts/configuration) - 配置使用详解
- [部署指南](../deployment/production) - 生产环境配置

---

**提示**: 生产环境的敏感配置务必使用环境变量，不要写在配置文件中！

