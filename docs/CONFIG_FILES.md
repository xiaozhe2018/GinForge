# 配置文件自动加载机制

## 自动发现机制

GinForge 的配置系统采用**自动发现**机制，无需修改代码即可添加新配置文件。

### 工作原理

1. **base.yaml 必须存在**：作为基础配置，最先加载
2. **自动扫描目录**：自动加载 `configs/{env}/` 目录下的所有 `.yaml` 文件
3. **按文件名排序**：确保加载顺序一致性
4. **服务配置优先**：如果设置了 `SERVICE_NAME`，对应的服务配置文件会优先加载

### 加载顺序

配置文件按以下顺序加载：

1. `base.yaml` - 基础配置（必须，最先加载）
2. 如果设置了 `SERVICE_NAME`，对应的 `{service-name}.yaml` 优先加载
3. 其他 `.yaml` 文件按文件名排序加载

### 配置合并规则

- 后加载的配置会覆盖先加载的配置（如果键相同）
- 最终环境变量会覆盖所有配置文件的值

## 添加新配置文件

### 方式一：添加服务配置（推荐）

```bash
# 1. 在对应环境目录下创建服务配置文件
# 例如：为 new-service 创建配置
configs/dev/new-service.yaml
configs/prod/new-service.yaml
```

```yaml
# configs/dev/new-service.yaml
services:
  new_service:
    port: 8088

log:
  level: "debug"
  output: "stdout"

# 服务特定配置
new_service:
  custom_setting: "value"
```

**无需修改代码**，系统会自动加载！

### 方式二：添加功能模块配置

```bash
# 例如：添加缓存模块配置
configs/dev/cache.yaml
configs/prod/cache.yaml
```

```yaml
# configs/dev/cache.yaml
cache:
  enabled: true
  provider: "redis"
  ttl: "1h"
```

### 方式三：控制加载顺序

如果需要控制配置文件的加载顺序，可以使用文件名前缀：

```bash
configs/dev/01-base-extras.yaml    # 紧跟在 base.yaml 之后加载
configs/dev/02-service-config.yaml
configs/dev/99-final-overrides.yaml  # 最后加载，覆盖前面的配置
```

## 示例场景

### 场景1：添加新的微服务

```bash
# 1. 创建服务配置
touch configs/dev/payment-api.yaml
touch configs/prod/payment-api.yaml

# 2. 编写配置
# configs/dev/payment-api.yaml
services:
  payment_api:
    port: 8089

# 3. 运行服务（自动加载）
export SERVICE_NAME=payment-api
go run services/payment-api/cmd/server/main.go
```

**无需修改 `pkg/config/config.go`！**

### 场景2：添加第三方服务配置

```bash
# 添加 Elasticsearch 配置
touch configs/dev/elasticsearch.yaml
touch configs/prod/elasticsearch.yaml
```

```yaml
# configs/dev/elasticsearch.yaml
elasticsearch:
  enabled: true
  hosts:
    - "http://localhost:9200"
  index_prefix: "dev_"
```

### 场景3：环境特定覆盖

```bash
# 开发环境：使用本地配置
configs/dev/local-storage.yaml

# 生产环境：使用云存储配置
configs/prod/cloud-storage.yaml
```

## 配置文件命名规范

### 推荐命名方式

- **服务配置**：`{service-name}.yaml`（如：`admin-api.yaml`）
- **功能模块**：`{module-name}.yaml`（如：`cache.yaml`、`storage.yaml`）
- **第三方服务**：`{service-name}.yaml`（如：`elasticsearch.yaml`、`kafka.yaml`）

### 命名约定

- 使用小写字母和连字符（kebab-case）
- 避免使用特殊字符
- 文件名清晰表达配置用途

## 注意事项

1. **base.yaml 不能删除**：这是必须的基础配置文件
2. **配置文件顺序**：如果需要控制加载顺序，使用文件名前缀
3. **配置冲突**：后加载的配置会覆盖先加载的配置
4. **环境变量优先**：环境变量始终会覆盖配置文件的值

## 最佳实践

1. **服务配置独立文件**：每个服务的配置放在独立的文件中
2. **模块化配置**：功能模块的配置独立管理
3. **环境隔离**：不同环境使用不同的配置目录
4. **敏感信息使用环境变量**：密码、密钥等通过 `.env` 或环境变量设置

## 常见问题

### Q: 添加新配置文件后需要重启服务吗？

A: 需要重启服务，因为配置文件在启动时加载。如果启用了配置热加载（`WatchConfig`），修改配置文件后会自动重新加载。

### Q: 如何删除不需要的配置文件？

A: 直接删除文件即可，系统会自动跳过不存在的文件。

### Q: 配置文件太多会影响性能吗？

A: 不会。配置文件只在启动时加载一次，加载顺序对性能影响很小。建议保持配置文件的合理数量和大小。

### Q: 如何调试配置加载顺序？

A: 可以在代码中添加日志，或者使用 `config.AllSettings()` 查看最终合并后的配置。

