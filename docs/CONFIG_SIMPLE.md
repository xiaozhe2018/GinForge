# 简化配置管理方案

## 核心理念

**只需要维护一套配置，无需同步！**

## 推荐方案：YAML 作为文档 + .env 覆盖

### 工作方式

1. **YAML 文件保留默认值**（开发环境直接使用）
2. **.env 只配置需要覆盖的值**（无需同步所有配置）
3. **环境变量自动覆盖**（无需手动同步）

### 实际操作

#### 开发环境（最简单）

```bash
# .env - 只需要配置环境变量
GOEASE_APP_ENV=dev

# 其他配置使用 YAML 默认值，无需配置
# 不需要在 .env 中重复 YAML 中的所有配置
```

**YAML 文件作为：**
- ✅ 默认配置（开发环境直接使用）
- ✅ 配置文档（展示所有可用配置项）
- ✅ 类型和格式参考

#### 生产环境（只覆盖必要项）

```bash
# .env - 只配置需要覆盖的值
GOEASE_APP_ENV=prod
GOEASE_DATABASE_HOST=prod-db-host
GOEASE_DATABASE_PASSWORD=secure_password
GOEASE_REDIS_PASSWORD=redis_password
GOEASE_JWT_SECRET=prod-secret-key

# 其他配置使用 YAML 默认值或通过 Docker/K8s 环境变量注入
# 不需要在 .env 中列出所有配置
```

### 关键点

1. **无需同步**：.env 和 YAML 不需要保持一致
2. **YAML 作为文档**：保留所有配置项的默认值和说明
3. **.env 只覆盖**：只配置需要覆盖的值（特别是敏感信息）
4. **自动合并**：系统自动合并，环境变量优先级最高

## 三种简化方案对比

### 方案一：YAML 为主，.env 为辅（推荐）

**适用场景**：大多数配置是静态的，只有少量敏感信息需要覆盖

```yaml
# configs/dev/base.yaml - 保留所有默认值
database:
  host: "localhost"      # 默认值
  port: 3306
  username: "root"
  password: "123456"     # 默认值，可以通过 .env 覆盖
```

```bash
# .env - 只覆盖需要改的值
GOEASE_APP_ENV=dev
GOEASE_DATABASE_PASSWORD=my_password  # 只覆盖这一个
```

**优点：**
- ✅ YAML 保留完整配置结构（作为文档）
- ✅ .env 只需配置少量覆盖值
- ✅ 开发环境开箱即用（使用 YAML 默认值）

### 方案二：.env 为主，YAML 为辅

**适用场景**：需要频繁修改配置，或者配置主要在环境变量中管理

```bash
# .env - 所有配置都在这里
GOEASE_APP_ENV=dev
GOEASE_DATABASE_HOST=localhost
GOEASE_DATABASE_PORT=3306
GOEASE_DATABASE_PASSWORD=123456
# ... 所有配置
```

```yaml
# configs/dev/base.yaml - 保留最小默认值或占位符
database:
  host: "${DB_HOST}"      # 占位符，实际由环境变量提供
  password: "${DB_PASSWORD}"
```

**优点：**
- ✅ 所有配置集中在一个地方
- ✅ 便于环境管理

**缺点：**
- ❌ .env 文件会很大
- ❌ 失去 YAML 的配置文档作用

### 方案三：完全环境变量（Docker/K8s 场景）

**适用场景**：生产环境，通过容器编排工具管理

```yaml
# configs/prod/base.yaml - 完全使用占位符
database:
  host: "${DB_HOST}"
  password: "${DB_PASSWORD}"
  # 实际值由环境变量或容器注入
```

```yaml
# docker-compose.yml 或 K8s ConfigMap
environment:
  - GOEASE_DATABASE_HOST=mysql
  - GOEASE_DATABASE_PASSWORD=${DB_PASSWORD}
```

**优点：**
- ✅ 完全分离配置和代码
- ✅ 适合云原生环境

## 推荐实践：混合方案

### 开发环境

```yaml
# configs/dev/base.yaml - 保留所有默认值
database:
  host: "localhost"
  port: 3306
  password: "123456"  # 开发默认值
```

```bash
# .env - 最小配置
GOEASE_APP_ENV=dev
# 其他使用 YAML 默认值
```

### 生产环境

```yaml
# configs/prod/base.yaml - 使用占位符或默认值
database:
  host: "mysql.prod.internal"
  password: "${DB_PASSWORD}"  # 占位符，实际由环境变量提供
```

```bash
# .env 或 Docker/K8s 环境变量
GOEASE_APP_ENV=prod
GOEASE_DATABASE_PASSWORD=secure_password
```

## 最佳实践总结

### 开发时

1. **使用 YAML 默认值**：大多数配置用 YAML 的默认值
2. **.env 只配置必要项**：比如 `GOEASE_APP_ENV` 和需要覆盖的敏感信息
3. **无需同步**：不需要在 .env 中列出所有配置

### 生产时

1. **YAML 使用占位符**：作为文档说明，实际值由环境变量提供
2. **通过环境变量注入**：Docker/K8s 通过环境变量注入敏感配置
3. **集中管理**：在 CI/CD 或容器编排工具中管理环境变量

### 添加新配置

1. **在 YAML 中添加默认值**：作为文档和默认值
2. **如果需要在 .env 覆盖**：添加对应的环境变量
3. **无需修改 config.go**：系统自动发现和加载

## 常见误解

### ❌ 误解1：需要在 .env 中列出所有配置

**事实**：不需要！只需配置需要覆盖的值。

```bash
# ❌ 错误：列出所有配置
GOEASE_DATABASE_HOST=localhost
GOEASE_DATABASE_PORT=3306
GOEASE_DATABASE_USERNAME=root
GOEASE_DATABASE_PASSWORD=123456
# ... 所有配置

# ✅ 正确：只配置需要覆盖的
GOEASE_APP_ENV=dev
GOEASE_DATABASE_PASSWORD=my_password  # 只覆盖这一个
```

### ❌ 误解2：.env 和 YAML 必须保持一致

**事实**：不需要！环境变量会自动覆盖 YAML 的值。

```yaml
# YAML 保留默认值
database:
  host: "localhost"      # 如果 .env 没设置，使用这个
  password: "123456"     # 如果 .env 设置了，使用 .env 的值
```

```bash
# .env 只覆盖需要的
GOEASE_DATABASE_PASSWORD=my_password  # 覆盖 YAML 中的值
# host 使用 YAML 的默认值 "localhost"
```

### ❌ 误解3：修改配置需要改多个地方

**事实**：
- **开发环境**：改 YAML 即可（.env 只配置环境名）
- **生产环境**：改环境变量即可（YAML 作为文档）

## 简化工作流

### 日常开发

```bash
# 1. 改配置
vim configs/dev/base.yaml

# 2. 重启服务
make run

# .env 基本不需要动
```

### 环境切换

```bash
# 1. 改环境变量
export GOEASE_APP_ENV=test

# 2. 如果测试环境需要特殊配置，在 .env 中覆盖
# 或在 configs/test/base.yaml 中设置

# 3. 重启服务
make run
```

### 添加新配置项

```bash
# 1. 在 YAML 中添加（作为文档和默认值）
# configs/dev/base.yaml
new_feature:
  enabled: true
  setting: "default"

# 2. 如果需要覆盖，在 .env 中添加
GOEASE_NEW_FEATURE_ENABLED=false

# 3. 无需修改 config.go（自动发现）
```

## 总结

**核心原则：**
1. **YAML = 文档 + 默认值**
2. **.env = 只配置需要覆盖的值**
3. **无需同步，系统自动合并**

这样就不需要同时维护两套配置了！只需要：
- 开发时改 YAML
- 生产时通过环境变量覆盖
- .env 只配置环境名和敏感信息

