# GinForge 配置最佳实践

## 🎯 核心原则：只需维护一套配置，无需同步！

**推荐方式：**
- **YAML** = 文档 + 默认值（开发环境直接使用）
- **.env** = 只配置需要覆盖的值（无需同步所有配置）
- **环境变量** = 自动覆盖（无需手动同步）

详见：[简化配置管理方案](./CONFIG_SIMPLE.md)

---

## 配置管理方案

GinForge 采用 **环境变量优先 + YAML 默认值** 的配置管理方案，这是业界最佳实践。

### 配置优先级

配置读取优先级（从高到低）：

1. **环境变量**（`GOEASE_*` 前缀）- 最高优先级
2. **.env 文件** - 自动加载的环境变量
3. **YAML 配置文件**（`configs/{env}/base.yaml`）- 默认配置和文档
4. **代码中的默认值** - 最终兜底

### 工作原理

1. **YAML 作为默认配置和文档**
   - 包含所有配置项的默认值
   - 作为配置结构说明文档
   - 开发环境可以直接使用

2. **.env 文件覆盖敏感配置**
   - 管理敏感信息（密码、密钥等）
   - 不会被提交到代码库
   - 自动加载到环境变量

3. **环境变量最终决定**
   - Viper 自动读取环境变量
   - 优先级最高，会覆盖 YAML 和 .env 中的值

## 使用方式

### 1. 开发环境（最简单）

**使用 YAML 默认值，.env 只配置环境名**

```bash
# 1. 复制 env.example 为 .env
cp env.example .env

# 2. 编辑 .env，只配置环境变量（不需要配置所有值！）
GOEASE_APP_ENV=dev

# 3. 其他配置使用 YAML 默认值，无需在 .env 中重复配置
# 直接运行（自动加载）
make run
```

**YAML 配置文件包含所有默认值**：
```yaml
# configs/dev/base.yaml - 包含所有默认配置
database:
  host: "localhost"      # 开发环境默认值
  port: 3306
  username: "root"
  password: "123456"     # 开发环境默认值，可以直接用
```

**关键：不需要在 .env 中列出所有配置！只需配置需要覆盖的值。**

### 2. 生产环境

**完全使用环境变量覆盖敏感配置**

```bash
# .env 文件（不提交到代码库）
GOEASE_APP_ENV=prod
GOEASE_DATABASE_HOST=prod-db-host
GOEASE_DATABASE_PASSWORD=secure_password
GOEASE_REDIS_PASSWORD=redis_password
GOEASE_JWT_SECRET=production-secret-key
GOEASE_SMTP_PASSWORD=smtp_password
```

**YAML 配置提供结构**：
```yaml
# configs/prod/base.yaml
database:
  host: "${DB_HOST}"      # 使用环境变量占位符（文档说明）
  password: "${DB_PASSWORD}"  # 实际由环境变量覆盖
```

## 配置兼容性说明

### YAML 中的默认值

YAML 文件保留所有配置项的默认值，用于：

1. **开发环境直接使用**：开发时不需要配置环境变量
2. **配置结构说明**：作为配置项的完整列表和说明
3. **类型和格式参考**：显示每个配置项的类型和格式

### 环境变量覆盖

所有 YAML 配置都可以通过环境变量覆盖：

| YAML 配置项 | 环境变量 |
|------------|---------|
| `database.host` | `GOEASE_DATABASE_HOST` |
| `database.password` | `GOEASE_DATABASE_PASSWORD` |
| `jwt.secret` | `GOEASE_JWT_SECRET` |
| `services.admin_api.port` | `GOEASE_SERVICES_ADMIN_API_PORT` |

**转换规则**：
- YAML 的层级用下划线连接：`database.host` → `GOEASE_DATABASE_HOST`
- 嵌套配置自动展开：`services.admin_api.port` → `GOEASE_SERVICES_ADMIN_API_PORT`

### 自动加载 .env

代码会自动加载 `.env` 文件：

```go
// pkg/config/config.go
func New() *Config {
    // 自动加载 .env 文件（如果存在）
    _ = godotenv.Load()
    
    // ... 后续配置读取
}
```

**无需手动加载**：
- ✅ `make run` - 自动加载
- ✅ `go run` - 自动加载
- ✅ `./bin/admin-api` - 自动加载

## 实际使用示例

### 开发环境设置

```bash
# .env
GOEASE_APP_ENV=dev
# 其他配置使用 YAML 默认值，无需设置
```

### 生产环境设置

```bash
# .env（不提交到 Git）
GOEASE_APP_ENV=prod
GOEASE_DATABASE_HOST=mysql.prod.internal
GOEASE_DATABASE_PASSWORD=prod_secret_password
GOEASE_REDIS_HOST=redis.prod.internal
GOEASE_REDIS_PASSWORD=redis_secret_password
GOEASE_JWT_SECRET=prod_jwt_secret_change_me
GOEASE_EMAIL_SMTP_PASSWORD=smtp_secret_password
```

### Docker 环境

```yaml
# docker-compose.yml
services:
  admin-api:
    environment:
      - GOEASE_APP_ENV=prod
      - GOEASE_DATABASE_HOST=mysql
      - GOEASE_DATABASE_PASSWORD=${DB_PASSWORD}
      - GOEASE_JWT_SECRET=${JWT_SECRET}
```

## 最佳实践总结

### 简化原则：只需维护一套配置

1. **开发环境**：
   - YAML 保留所有默认值（直接使用）
   - .env 只配置 `GOEASE_APP_ENV=dev`
   - 其他配置使用 YAML 默认值，无需在 .env 中重复配置

2. **生产环境**：
   - YAML 使用占位符或默认值（作为文档）
   - .env 或环境变量只配置需要覆盖的值（敏感信息）
   - 通过 Docker/K8s 环境变量注入敏感配置

3. **添加新配置**：
   - 在 YAML 中添加默认值（作为文档）
   - 如需覆盖，在 .env 中添加环境变量
   - **无需修改 config.go**（自动发现）

4. **核心优势**：
   - ✅ 无需同步：.env 和 YAML 不需要保持一致
   - ✅ 自动合并：环境变量自动覆盖 YAML 配置
   - ✅ 最小配置：.env 只配置需要覆盖的值

### 工作流示例

```bash
# 日常开发：改 YAML 即可
vim configs/dev/base.yaml
make run

# 环境切换：改环境变量
export GOEASE_APP_ENV=test
make run

# 生产部署：通过环境变量覆盖
docker run -e GOEASE_DATABASE_PASSWORD=xxx ...
```

## 常见问题

### Q: 为什么 YAML 中还有默认值？直接都从环境变量读取不行吗？

A: YAML 默认值的好处：
- 开发环境无需配置即可运行
- 作为配置项的完整文档
- 展示配置结构和类型
- 环境变量是可选的覆盖，不是必须的

### Q: .env 文件需要提交到 Git 吗？

A: **不需要**。`.env` 应该添加到 `.gitignore`：
```gitignore
.env
.env.local
.env.*.local
```

### Q: 如何为不同环境管理配置？

A: 使用不同的 `.env` 文件：
- `.env.dev` - 开发环境
- `.env.test` - 测试环境
- `.env.prod` - 生产环境（不提交）

或在 CI/CD 中直接设置环境变量，不使用 `.env` 文件。

### Q: 生产环境 YAML 中的占位符 `${DB_HOST}` 会被替换吗？

A: **不会**。这些只是文档说明，实际值由环境变量提供。Viper 不会替换 YAML 中的占位符，而是直接从环境变量读取。

