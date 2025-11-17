# GinForge 环境配置方案

## 配置方案设计

## 配置组织方式：按环境目录组织

```
configs/
├── dev/                    # 开发环境
│   ├── base.yaml          # 基础配置（所有服务共用）
│   ├── admin-api.yaml     # admin-api 服务配置
│   ├── user-api.yaml      # user-api 服务配置
│   ├── merchant-api.yaml  # merchant-api 服务配置
│   ├── gateway.yaml       # gateway 服务配置
│   ├── gateway-worker.yaml # gateway-worker 服务配置
│   ├── file-api.yaml      # file-api 服务配置
│   ├── websocket-gateway.yaml # websocket-gateway 服务配置
│   └── demo.yaml          # demo 服务配置
├── test/                   # 测试环境
│   ├── base.yaml
│   ├── admin-api.yaml
│   └── ...
└── prod/                   # 生产环境
    ├── base.yaml
    ├── admin-api.yaml
    └── ...
```

### 优点
- 环境隔离清晰，每个环境一个目录
- 便于版本控制和部署
- 支持环境级别的配置覆盖
- 便于 CI/CD 流程
- 配置结构统一，易于维护

### 配置读取逻辑

配置读取优先级（从高到低）：
1. **环境变量**（`GOEASE_*` 前缀）- 最高优先级
2. **服务配置**（`configs/{env}/{service-name}.yaml`）
3. **基础配置**（`configs/{env}/base.yaml`）
4. **默认值**（代码中定义）

**配置要求：**
- 必须存在 `configs/{env}/base.yaml` 文件
- 如果文件不存在，系统会报错并提示设置 `GOEASE_APP_ENV` 环境变量

### 环境变量设置

**方式一：使用 .env 文件（推荐）**
```bash
# 1. 复制 env.example 为 .env
cp env.example .env

# 2. 编辑 .env 文件，设置环境
GOEASE_APP_ENV=dev  # 或 test, prod

# 3. 设置服务名称（可选）
SERVICE_NAME=admin-api

# 4. 加载环境变量（如果需要）
export $(cat .env | xargs)
```

**方式二：直接设置环境变量**
```bash
# 设置环境（必需）
export GOEASE_APP_ENV=dev  # 或 test, prod

# 设置服务名称（可选，用于加载服务特定配置）
# 如果不设置，会尝试从程序名自动推断
export SERVICE_NAME=admin-api
```

### 环境值说明

| 环境变量值 | 对应目录 | 说明 |
|-----------|---------|------|
| `dev` | `configs/dev/` | 开发环境（默认值） |
| `test` | `configs/test/` | 测试环境 |
| `prod` | `configs/prod/` | 生产环境 |

### 使用示例

#### 开发环境运行 admin-api
```bash
export GOEASE_APP_ENV=dev
export SERVICE_NAME=admin-api
go run services/admin-api/cmd/server/main.go
```

#### 测试环境运行 admin-api
```bash
export GOEASE_APP_ENV=test
export SERVICE_NAME=admin-api
go run services/admin-api/cmd/server/main.go
```

#### 生产环境运行 admin-api
```bash
export GOEASE_APP_ENV=prod
export SERVICE_NAME=admin-api
./bin/admin-api
```

#### Docker 环境
```yaml
# docker-compose.yml
services:
  admin-api:
    environment:
      - GOEASE_APP_ENV=prod
      - SERVICE_NAME=admin-api
      - DB_HOST=mysql
      - DB_PASSWORD=secret
```

## 实现状态

### ✅ 已实现功能

1. **配置读取逻辑**（`pkg/config/config.go`）
   - ✅ 支持从 `configs/{env}/base.yaml` 读取基础配置
   - ✅ 支持从 `configs/{env}/{service-name}.yaml` 读取服务配置
   - ✅ 支持环境变量 `GOEASE_APP_ENV` 指定环境
   - ✅ 支持环境变量 `SERVICE_NAME` 指定服务
   - ✅ 如果配置文件不存在，会给出明确的错误提示

2. **示例配置文件**
   - ✅ `configs/dev/base.yaml` - 开发环境基础配置
   - ✅ `configs/dev/admin-api.yaml` - admin-api 开发环境配置
   - ✅ `configs/prod/base.yaml` - 生产环境基础配置
   - ✅ `configs/prod/admin-api.yaml` - admin-api 生产环境配置
   - ✅ `configs/test/base.yaml` - 测试环境基础配置

### 2. 配置文件示例

#### `configs/dev/base.yaml`
```yaml
# 开发环境基础配置
app:
  env: development
  debug: true

database:
  driver: mysql
  host: localhost
  port: 3306
  database: gin_forge
  username: root
  password: 123456

redis:
  enabled: false
  host: localhost
  port: 6379
```

#### `configs/dev/admin-api.yaml`
```yaml
# admin-api 服务开发环境配置
services:
  admin_api:
    port: 8083

log:
  level: debug
  output: stdout
```

#### `configs/prod/base.yaml`
```yaml
# 生产环境基础配置
app:
  env: production
  debug: false

database:
  driver: mysql
  host: ${DB_HOST}
  port: ${DB_PORT}
  database: ${DB_NAME}
  username: ${DB_USER}
  password: ${DB_PASSWORD}

redis:
  enabled: true
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
```

#### `configs/prod/admin-api.yaml`
```yaml
# admin-api 服务生产环境配置
services:
  admin_api:
    port: 8083

log:
  level: info
  output: file
  file_path: /var/log/ginforge/admin-api.log
```

## 服务列表

当前支持的服务配置：

| 服务名称 | 配置文件 | 默认端口 |
|---------|---------|---------|
| admin-api | `admin-api.yaml` | 8083 |
| user-api | `user-api.yaml` | 8081 |
| merchant-api | `merchant-api.yaml` | 8082 |
| gateway | `gateway.yaml` | 8080 |
| gateway-worker | `gateway-worker.yaml` | 8084 |
| file-api | `file-api.yaml` | 8086 |
| websocket-gateway | `websocket-gateway.yaml` | 8087 |
| demo | `demo.yaml` | 8085 |

## 添加新配置文件

### 自动发现机制

系统会自动扫描 `configs/{env}/` 目录下的所有 `.yaml` 文件并加载，**无需修改代码**。

### 添加新服务配置

1. 在对应环境目录下创建服务配置文件：
   ```bash
   # 开发环境
   configs/dev/your-service.yaml
   
   # 生产环境
   configs/prod/your-service.yaml
   ```

2. 编写配置内容：
   ```yaml
   # configs/dev/your-service.yaml
   services:
     your_service:
       port: 8088
   
   log:
     level: "debug"
   ```

3. 运行服务（自动加载）：
   ```bash
   export SERVICE_NAME=your-service
   go run services/your-service/cmd/server/main.go
   ```

**就这么简单！无需修改 `pkg/config/config.go`。**

### 添加功能模块配置

同样只需创建文件即可：

```bash
# 创建配置文件
touch configs/dev/cache.yaml
touch configs/prod/cache.yaml
```

详细说明请参考：[配置文件自动加载机制](./CONFIG_FILES.md)

## 配置要求

系统要求必须使用按环境目录组织的方式：
- `configs/{env}/base.yaml` - 基础配置（必需）
- `configs/{env}/{service-name}.yaml` - 服务配置（可选）

如果 `configs/{env}/base.yaml` 不存在，系统会报错并提示设置正确的 `GOEASE_APP_ENV` 环境变量。

**注意：** 扁平化配置文件（`config.yaml`, `config.prod.yaml` 等）已被废弃，系统不会读取这些文件。

## 环境变量优先级

1. **环境变量**（`GOEASE_*`）- 最高优先级
2. **服务配置**（`configs/{env}/{service}.yaml`）
3. **基础配置**（`configs/{env}/base.yaml`）
4. **默认值**（代码中定义）

