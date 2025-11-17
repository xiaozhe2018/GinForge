# GinForge 配置文件说明

## 📁 目录结构

```
configs/
├── dev/                    # 开发环境配置
│   ├── base.yaml          # 基础配置（所有服务共用）
│   ├── admin-api.yaml     # admin-api 服务配置
│   ├── user-api.yaml      # user-api 服务配置
│   ├── merchant-api.yaml  # merchant-api 服务配置
│   ├── gateway.yaml       # gateway 服务配置
│   ├── gateway-worker.yaml # gateway-worker 服务配置
│   ├── file-api.yaml      # file-api 服务配置
│   ├── websocket-gateway.yaml # websocket-gateway 服务配置
│   └── demo.yaml          # demo 服务配置
├── test/                   # 测试环境配置
│   ├── base.yaml
│   └── ...
├── prod/                   # 生产环境配置
│   ├── base.yaml
│   └── ...
├── config.yaml            # 向后兼容：扁平化配置
├── config.prod.yaml       # 向后兼容：生产环境配置
├── config.test.yaml       # 向后兼容：测试环境配置
├── file-api.yaml          # 向后兼容：file-api 配置
├── notification.yaml      # 向后兼容：通知服务配置
├── README.md              # 本文件
└── SERVICES.md            # 服务配置清单
```

## 配置组织方式

GinForge 使用按环境目录组织配置：

```
configs/
├── dev/              # 开发环境
│   ├── base.yaml     # 基础配置（所有服务共用）
│   ├── admin-api.yaml
│   ├── user-api.yaml
│   └── ...
├── test/             # 测试环境
│   ├── base.yaml
│   └── ...
└── prod/             # 生产环境
    ├── base.yaml
    └── ...
```

**优点：**
- 环境隔离清晰
- 便于版本控制
- 支持环境级别的配置覆盖
- 配置结构统一

## 快速开始

### 1. 配置环境变量（超简单）

**方式一：使用 .env 文件（推荐）**
```bash
# 复制 env.example 为 .env
cp env.example .env

# 编辑 .env 文件，只配置环境变量（其他使用 YAML 默认值）
GOEASE_APP_ENV=dev  # 或 test, prod

# 如果需要覆盖某些配置，添加对应的环境变量
# GOEASE_DATABASE_PASSWORD=my_password

# 直接运行，.env 会自动加载
make run
```

**关键：不需要在 .env 中列出所有配置！YAML 文件已包含所有默认值。**

**方式二：直接设置环境变量**
```bash
# 设置环境（必需）
export GOEASE_APP_ENV=dev  # 或 test, prod

# 设置服务名称（可选，用于加载服务特定配置）
export SERVICE_NAME=admin-api
```

### 2. 运行服务

```bash
# 开发环境
export GOEASE_APP_ENV=dev
export SERVICE_NAME=admin-api
go run services/admin-api/cmd/server/main.go

# 测试环境
export GOEASE_APP_ENV=test
export SERVICE_NAME=admin-api
go run services/admin-api/cmd/server/main.go

# 生产环境
export GOEASE_APP_ENV=prod
export SERVICE_NAME=admin-api
./bin/admin-api
```

### 3. 环境说明

| 环境变量值 | 对应目录 | 说明 |
|-----------|---------|------|
| `dev` | `configs/dev/` | 开发环境（默认） |
| `test` | `configs/test/` | 测试环境 |
| `prod` | `configs/prod/` | 生产环境 |

## 配置优先级

配置读取优先级（从高到低）：

1. **环境变量**（`GOEASE_*` 前缀）- 最高优先级
2. **服务配置**（`configs/{env}/{service-name}.yaml`）
3. **基础配置**（`configs/{env}/base.yaml`）
4. **默认值**（代码中定义）

## 环境变量说明

### 必需环境变量

- `GOEASE_APP_ENV`: 环境名称（`dev`/`test`/`prod`），默认为 `dev`

### 可选环境变量

- `SERVICE_NAME`: 服务名称（如果不设置，会尝试从程序名自动推断）

### 配置覆盖环境变量

所有配置项都可以通过环境变量覆盖，格式：`GOEASE_{KEY}`

例如：
- `GOEASE_DATABASE_HOST` 覆盖 `database.host`
- `GOEASE_SERVICES_ADMIN_API_PORT` 覆盖 `services.admin_api.port`

## 配置文件示例

### 开发环境基础配置

`configs/dev/base.yaml` - 包含所有服务共用的配置（数据库、Redis、日志等）

### 服务特定配置

`configs/dev/admin-api.yaml` - admin-api 服务的特定配置（端口、日志级别等）

## 添加新服务配置

1. 在对应环境目录下创建服务配置文件：
   ```bash
   configs/dev/your-service.yaml
   configs/prod/your-service.yaml
   ```

2. 设置环境变量：
   ```bash
   export SERVICE_NAME=your-service
   ```

3. 配置会自动合并到基础配置之上

## 服务配置清单

所有服务的配置文件已创建完成，详情请参考：[SERVICES.md](./SERVICES.md)

## 配置文件说明

以下扁平化配置文件已废弃，仅保留用于参考：
- `config.yaml` - 已废弃，请使用 `dev/base.yaml`
- `config.prod.yaml` - 已废弃，请使用 `prod/base.yaml`
- `config.test.yaml` - 已废弃，请使用 `test/base.yaml`
- `file-api.yaml` - 已废弃，请使用 `{env}/file-api.yaml`
- `notification.yaml` - 已废弃，通知配置已整合到 `{env}/base.yaml`

**注意：** 系统只会从环境目录（`dev/`, `test/`, `prod/`）读取配置，扁平化配置文件将被忽略。

## 配置兼容性说明

### YAML vs 环境变量

**YAML 文件的作用：**
- ✅ 提供默认配置值（开发环境可直接使用）
- ✅ 作为配置结构文档（展示所有可用配置项）
- ✅ 展示配置格式和类型

**环境变量的作用：**
- ✅ 覆盖 YAML 中的配置（优先级最高）
- ✅ 管理敏感信息（密码、密钥等）
- ✅ 不同环境使用不同配置

**自动加载 .env 文件：**
- ✅ 代码会自动加载 `.env` 文件到环境变量
- ✅ `make run` 无需手动加载环境变量
- ✅ `.env` 文件不会被提交到代码库

**配置优先级**：
1. 环境变量（`GOEASE_*`）
2. .env 文件
3. YAML 配置文件
4. 代码默认值

### 简化使用建议

**开发环境（最简单）：**
```bash
# .env - 只需要配置环境变量，其他使用 YAML 默认值
GOEASE_APP_ENV=dev

# 不需要在 .env 中列出所有配置！
# YAML 文件已包含所有默认值，直接使用即可
```

**如果需要覆盖某些配置：**
```bash
# .env - 只配置需要覆盖的值
GOEASE_APP_ENV=dev
GOEASE_DATABASE_PASSWORD=my_password  # 只覆盖这一个
# 其他配置使用 YAML 默认值，无需同步
```

**生产环境：**
```bash
# .env 或 Docker/K8s 环境变量 - 只配置需要覆盖的值
GOEASE_APP_ENV=prod
GOEASE_DATABASE_HOST=prod-db
GOEASE_DATABASE_PASSWORD=secure_password
# 其他配置使用 YAML 默认值或通过容器注入
```

**关键点：**
- ✅ **无需同步**：.env 和 YAML 不需要保持一致
- ✅ **YAML 作为文档**：保留所有配置项的默认值
- ✅ **.env 只覆盖**：只配置需要覆盖的值（特别是敏感信息）
- ✅ **自动合并**：系统自动合并，环境变量优先级最高

## 详细文档

更多信息请参考：
- [环境配置指南](../docs/CONFIG_ENV_GUIDE.md) - 详细配置说明
- [配置最佳实践](../docs/CONFIG_BEST_PRACTICES.md) - 配置管理最佳实践
- [服务配置清单](./SERVICES.md) - 所有服务配置文件清单

