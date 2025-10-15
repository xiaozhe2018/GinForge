# 配置系统

GinForge 提供了强大灵活的配置管理系统，支持多种配置源和环境。

## 🎯 配置特点

- ✅ **多配置源**：YAML + 环境变量 + 默认值
- ✅ **配置优先级**：环境变量 > .env > config.yaml > 默认值
- ✅ **类型安全**：强类型配置读取
- ✅ **多环境支持**：dev/test/prod 环境切换
- ✅ **热重载**：配置变更自动生效（开发中）

## 📁 配置文件结构

```
configs/
├── config.yaml         # 默认配置（开发环境）
├── config.test.yaml    # 测试环境配置
├── config.prod.yaml    # 生产环境配置
├── file-api.yaml       # 文件服务配置
└── notification.yaml   # 通知服务配置
```

## 📝 配置文件示例

### 基础配置 (`config.yaml`)

```yaml
# 应用配置
app:
  name: "GinForge"
  version: "1.0.0"
  env: "development"
  debug: true
  port: 8080
  read_timeout: 60s
  write_timeout: 60s
  idle_timeout: 120s

# 服务端口配置
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

# 数据库配置
database:
  type: "sqlite"          # sqlite, mysql, postgres
  host: "localhost"
  port: 3306
  database: "goweb.db"    # SQLite 文件路径
  username: "root"
  password: ""
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

# Redis 配置
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  pool_timeout: 4s

# JWT 配置
jwt:
  secret: "your-secret-key-change-in-production"
  expire_hours: 24

# 日志配置
log:
  level: "debug"          # debug, info, warn, error
  format: "json"          # json, text
  output: "stdout"        # stdout, file
  file_path: "logs/app.log"
```

## 🔧 使用配置

### 1. 基础用法

```go
package main

import (
    "goweb/pkg/config"
    "fmt"
)

func main() {
    // 加载配置
    cfg := config.New()
    
    // 读取字符串
    appName := cfg.GetString("app.name")
    fmt.Println("App Name:", appName)
    
    // 读取整数
    port := cfg.GetInt("app.port")
    fmt.Println("Port:", port)
    
    // 读取布尔值
    debug := cfg.GetBool("app.debug")
    fmt.Println("Debug:", debug)
    
    // 读取时间间隔
    timeout := cfg.GetDuration("app.read_timeout")
    fmt.Println("Timeout:", timeout)
}
```

### 2. 读取嵌套配置

```go
// 读取数据库配置
dbType := cfg.GetString("database.type")
dbHost := cfg.GetString("database.host")
dbPort := cfg.GetInt("database.port")

// 读取 Redis 配置
redisEnabled := cfg.GetBool("redis.enabled")
redisHost := cfg.GetString("redis.host")
```

### 3. 获取完整配置对象

```go
// 获取 Redis 配置对象
redisConfig := cfg.GetRedisConfig()
// 返回 config.RedisConfig 结构体

// 使用配置对象
if redisConfig.Enabled {
    fmt.Printf("Redis: %s:%d\n", redisConfig.Host, redisConfig.Port)
}
```

### 4. 判断运行环境

```go
// 检查是否为生产环境
if cfg.IsProduction() {
    // 生产环境逻辑
    log.Info("running in production mode")
}

// 检查是否为调试模式
if cfg.GetBool("app.debug") {
    // 调试模式逻辑
    log.Debug("debug mode enabled")
}
```

## 🌍 环境变量覆盖

### 创建 `.env` 文件

```bash
# 复制示例文件
cp env.example .env

# 编辑配置
vim .env
```

### `.env` 文件示例

```bash
# 应用配置
APP_ENV=development
APP_PORT=8083
APP_DEBUG=true

# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=gin_forge
DB_USERNAME=root
DB_PASSWORD=123456

# Redis 配置
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
```

### 环境变量命名规则

将 YAML 路径转换为大写，用下划线分隔：

```
app.port          → APP_PORT
database.host     → DB_HOST
redis.enabled     → REDIS_ENABLED
jwt.secret        → JWT_SECRET
```

## 🔐 生产环境配置

### 使用 `config.prod.yaml`

```yaml
# configs/config.prod.yaml
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
  password: "${DB_PASSWORD}"  # 从环境变量读取

redis:
  enabled: true
  host: "prod-redis.example.com"
  password: "${REDIS_PASSWORD}"

jwt:
  secret: "${JWT_SECRET}"
  expire_hours: 720  # 30 天

log:
  level: "info"
  format: "json"
  output: "file"
  file_path: "/var/log/ginforge/app.log"
```

### 启动生产环境

```bash
# 方式 1：设置环境变量
export APP_ENV=production
go run ./services/admin-api/cmd/server/main.go

# 方式 2：使用配置文件
APP_ENV=production DB_PASSWORD=xxx JWT_SECRET=xxx go run ./services/admin-api/cmd/server/main.go
```

## 💡 最佳实践

### 1. 敏感信息管理

❌ **不要** 将敏感信息写在配置文件中：

```yaml
# ❌ 错误示例
jwt:
  secret: "my-secret-key-123"
```

✅ **应该** 使用环境变量：

```yaml
# ✅ 正确示例
jwt:
  secret: "${JWT_SECRET}"
```

```bash
# 在 .env 或系统环境变量中设置
JWT_SECRET=randomly-generated-secret-key
```

### 2. 多环境配置

不同环境使用不同配置文件：

```bash
# 开发环境（默认）
go run main.go

# 测试环境
APP_ENV=test go run main.go

# 生产环境
APP_ENV=production go run main.go
```

### 3. 配置验证

启动时验证必需的配置：

```go
func validateConfig(cfg *config.Config) error {
    if cfg.GetString("jwt.secret") == "" {
        return errors.New("JWT secret is required")
    }
    
    if cfg.GetString("database.host") == "" {
        return errors.New("database host is required")
    }
    
    return nil
}
```

### 4. 默认值设置

为可选配置提供合理的默认值：

```go
// 读取端口，默认 8080
port := cfg.GetInt("app.port")
if port == 0 {
    port = 8080
}

// 读取日志级别，默认 info
logLevel := cfg.GetString("log.level")
if logLevel == "" {
    logLevel = "info"
}
```

## 🔄 配置热重载（开发中）

在开发环境，配置文件变更后自动重载：

```go
// TODO: 实现配置热重载
cfg.Watch(func(event config.Event) {
    log.Info("config changed", "key", event.Key, "value", event.Value)
    // 重新加载配置
})
```

## 📚 相关文档

- [环境变量完整列表](../api-reference/config-options)
- [生产部署配置](../deployment/production)
- [配置中心集成](../../demo/config_center_usage.md)

---

**下一步**: [学习路由管理](./routing)

