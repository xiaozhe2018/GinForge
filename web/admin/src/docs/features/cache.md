# 缓存系统

GinForge 提供了强大的缓存系统，支持内存缓存和 Redis 缓存。

## 🎯 为什么需要缓存？

- ⚡ **提升性能**：减少数据库查询
- 💰 **降低成本**：减少 CPU 和 I/O 消耗
- 🚀 **提高并发**：支持更多并发请求
- 📈 **改善体验**：更快的响应时间

## 🔧 配置 Redis

### 配置文件

```yaml
# configs/config.yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
```

### 启动 Redis

```bash
# 使用 Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 验证运行
docker exec redis redis-cli ping
# 输出：PONG
```

## 📝 基础用法

### 初始化 Redis 客户端

```go
package main

import (
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
)

func main() {
    cfg := config.New()
    log := logger.New("app", cfg.GetString("log.level"))
    
    // 初始化 Redis 客户端
    redisConfig := cfg.GetRedisConfig()
    if redisConfig.Enabled {
        redisClient := redis.NewClient(&redisConfig, log)
        log.Info("redis client initialized")
    }
}
```

### 基本操作

```go
import (
    "context"
    "time"
)

func cacheExample(redisClient *redis.Client) {
    ctx := context.Background()
    
    // 1. Set - 设置缓存
    err := redisClient.Set(ctx, "user:123", `{"name":"John"}`, 5*time.Minute)
    
    // 2. Get - 获取缓存
    value, err := redisClient.Get(ctx, "user:123")
    
    // 3. Delete - 删除缓存
    err = redisClient.Delete(ctx, "user:123")
    
    // 4. Exists - 检查是否存在
    exists, err := redisClient.Exists(ctx, "user:123")
    
    // 5. Expire - 设置过期时间
    err = redisClient.Expire(ctx, "user:123", 10*time.Minute)
}
```

## 🎨 缓存模式

### 1. Cache-Aside（旁路缓存）

最常用的缓存模式：

```go
func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error) {
    // 1. 先查缓存
    cacheKey := fmt.Sprintf("user:%s", userID)
    var user model.User
    
    cached, err := s.redisClient.Get(ctx, cacheKey)
    if err == nil && cached != "" {
        // 缓存命中
        json.Unmarshal([]byte(cached), &user)
        s.logger.Info("cache hit", "user_id", userID)
        return &user, nil
    }
    
    // 2. 缓存未命中，查询数据库
    user, err = s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    userData, _ := json.Marshal(user)
    s.redisClient.Set(ctx, cacheKey, string(userData), 10*time.Minute)
    
    s.logger.Info("cache miss", "user_id", userID)
    return &user, nil
}
```

### 2. Write-Through（写穿）

更新数据时同时更新缓存：

```go
func (s *UserService) UpdateUser(ctx context.Context, userID string, data map[string]interface{}) error {
    // 1. 更新数据库
    err := s.userRepo.Update(userID, data)
    if err != nil {
        return err
    }
    
    // 2. 更新缓存
    user, _ := s.userRepo.GetByID(userID)
    cacheKey := fmt.Sprintf("user:%s", userID)
    userData, _ := json.Marshal(user)
    s.redisClient.Set(ctx, cacheKey, string(userData), 10*time.Minute)
    
    return nil
}
```

### 3. Write-Behind（写回）

先更新缓存，异步写入数据库：

```go
func (s *UserService) UpdateUserAsync(ctx context.Context, userID string, data map[string]interface{}) error {
    // 1. 先更新缓存
    cacheKey := fmt.Sprintf("user:%s", userID)
    s.redisClient.Set(ctx, cacheKey, data, 10*time.Minute)
    
    // 2. 异步写入数据库
    go func() {
        time.Sleep(100 * time.Millisecond)
        s.userRepo.Update(userID, data)
    }()
    
    return nil
}
```

## 🔑 缓存键设计

### 命名规范

```
格式：业务:实体:ID:字段
示例：
  user:123               # 用户基本信息
  user:123:profile       # 用户资料
  order:456              # 订单信息
  config:system          # 系统配置
  stat:daily:2025-10-15  # 每日统计
```

### 使用常量

```go
const (
    CacheKeyUser    = "user:%s"
    CacheKeyOrder   = "order:%s"
    CacheKeyConfig  = "config:%s"
)

// 使用
cacheKey := fmt.Sprintf(CacheKeyUser, userID)
```

## ⏱️ 过期时间策略

### 不同数据的过期时间

```go
const (
    CacheTTLShort  = 1 * time.Minute   // 短期：频繁变化的数据
    CacheTTLMedium = 10 * time.Minute  // 中期：一般数据
    CacheTTLLong   = 1 * time.Hour     // 长期：很少变化的数据
    CacheTTLDay    = 24 * time.Hour    // 一天：配置数据
)

// 用户信息：10分钟
redisClient.Set(ctx, "user:123", data, CacheTTLMedium)

// 系统配置：24小时
redisClient.Set(ctx, "config:system", data, CacheTTLDay)

// 实时数据：1分钟
redisClient.Set(ctx, "stats:online", data, CacheTTLShort)
```

## 🔄 缓存更新策略

### 主动更新

数据变更时立即更新缓存：

```go
func (s *UserService) UpdateUser(userID string, data map[string]interface{}) error {
    // 1. 更新数据库
    err := s.userRepo.Update(userID, data)
    if err != nil {
        return err
    }
    
    // 2. 删除旧缓存（让下次查询时重建）
    cacheKey := fmt.Sprintf("user:%s", userID)
    s.redisClient.Delete(context.Background(), cacheKey)
    
    return nil
}
```

### 延迟双删

避免并发问题的缓存更新策略：

```go
func (s *UserService) UpdateUserWithDoubleDelete(userID string, data map[string]interface{}) error {
    cacheKey := fmt.Sprintf("user:%s", userID)
    ctx := context.Background()
    
    // 第一次删除缓存
    s.redisClient.Delete(ctx, cacheKey)
    
    // 更新数据库
    err := s.userRepo.Update(userID, data)
    if err != nil {
        return err
    }
    
    // 延迟第二次删除缓存（500ms 后）
    go func() {
        time.Sleep(500 * time.Millisecond)
        s.redisClient.Delete(ctx, cacheKey)
    }()
    
    return nil
}
```

## 🎯 实际应用场景

### 场景 1：系统配置缓存

```go
func (s *SystemService) GetConfig(key string) (string, error) {
    cacheKey := fmt.Sprintf("config:%s", key)
    ctx := context.Background()
    
    // 查缓存
    if cached, err := s.redisClient.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }
    
    // 查数据库
    value, err := s.configRepo.GetByKey(key)
    if err != nil {
        return "", err
    }
    
    // 写缓存（24小时）
    s.redisClient.Set(ctx, cacheKey, value, 24*time.Hour)
    
    return value, nil
}
```

### 场景 2：用户会话缓存

```go
// 登录时缓存用户信息
func (s *AuthService) Login(username, password string) (string, error) {
    // 验证用户...
    
    // 缓存用户会话（24小时）
    sessionKey := fmt.Sprintf("session:%s", token)
    sessionData := map[string]interface{}{
        "user_id":  user.ID,
        "username": user.Username,
        "roles":    user.Roles,
    }
    
    data, _ := json.Marshal(sessionData)
    s.redisClient.Set(context.Background(), sessionKey, string(data), 24*time.Hour)
    
    return token, nil
}
```

### 场景 3：统计数据缓存

```go
// 缓存今日统计数据（5分钟）
func (s *StatsService) GetTodayStats() (*Stats, error) {
    cacheKey := "stats:today"
    ctx := context.Background()
    
    // 查缓存
    var stats Stats
    if cached, err := s.redisClient.Get(ctx, cacheKey); err == nil {
        json.Unmarshal([]byte(cached), &stats)
        return &stats, nil
    }
    
    // 计算统计
    stats = s.calculateTodayStats()
    
    // 缓存 5 分钟
    data, _ := json.Marshal(stats)
    s.redisClient.Set(ctx, cacheKey, string(data), 5*time.Minute)
    
    return &stats, nil
}
```

## 🛡️ 缓存穿透/雪崩/击穿

### 缓存穿透防护

查询不存在的数据导致缓存穿透：

```go
func (s *UserService) GetUser(userID string) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    ctx := context.Background()
    
    // 查缓存
    cached, err := s.redisClient.Get(ctx, cacheKey)
    if err == nil {
        if cached == "null" {
            // 缓存了空值，说明数据不存在
            return nil, errors.New("用户不存在")
        }
        var user model.User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 查数据库
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        // 数据不存在，缓存空值（5分钟）
        s.redisClient.Set(ctx, cacheKey, "null", 5*time.Minute)
        return nil, err
    }
    
    // 缓存数据
    data, _ := json.Marshal(user)
    s.redisClient.Set(ctx, cacheKey, string(data), 10*time.Minute)
    
    return &user, nil
}
```

### 缓存雪崩防护

大量缓存同时过期：

```go
// 给缓存时间添加随机偏移
func setWithJitter(key string, value interface{}, baseTTL time.Duration) error {
    // 添加 ±10% 的随机偏移
    jitter := time.Duration(rand.Intn(int(baseTTL / 10)))
    ttl := baseTTL + jitter
    
    data, _ := json.Marshal(value)
    return redisClient.Set(ctx, key, string(data), ttl)
}

// 使用
setWithJitter("user:123", user, 10*time.Minute)
```

## 📊 缓存监控

### 统计缓存命中率

```go
type CacheStats struct {
    Hits   int64
    Misses int64
}

func (s *UserService) GetUser(userID string) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    cached, err := s.redisClient.Get(ctx, cacheKey)
    if err == nil {
        // 缓存命中
        s.stats.Hits++
        // ...
    } else {
        // 缓存未命中
        s.stats.Misses++
        // ...
    }
    
    // 计算命中率
    hitRate := float64(s.stats.Hits) / float64(s.stats.Hits + s.stats.Misses) * 100
    log.Info("cache hit rate", "rate", hitRate)
}
```

## 🎯 实际示例

查看完整示例：

- **Redis 客户端**: `pkg/redis/client.go`
- **缓存操作**: `pkg/redis/cache.go`
- **使用示例**: `docs/demo/cache.md`
- **Redis 使用**: `docs/demo/redis_usage.md`

## 💡 最佳实践

### 1. 缓存什么数据？

✅ **适合缓存**：
- 读多写少的数据
- 计算复杂的数据
- 配置信息
- 会话信息

❌ **不适合缓存**：
- 实时性要求高的数据
- 频繁变化的数据
- 大量数据（占用内存）

### 2. 设置合理的 TTL

```go
// 配置数据：1天
config: 24 * time.Hour

// 用户信息：10分钟
user: 10 * time.Minute

// 列表数据：1分钟
list: 1 * time.Minute

// 实时统计：10秒
stats: 10 * time.Second
```

### 3. 使用命名空间

```go
const (
    NamespaceUser   = "user"
    NamespaceOrder  = "order"
    NamespaceConfig = "config"
)

// 生成缓存键
func getCacheKey(namespace, id string) string {
    return fmt.Sprintf("%s:%s", namespace, id)
}
```

## 🎯 下一步

- [WebSocket](./websocket) - 实时通信
- [消息队列](../api-reference/queue) - 异步任务处理
- [Redis 使用文档](../../demo/redis_usage.md) - 更多 Redis 功能

---

**提示**: 合理使用缓存可以显著提升系统性能，但要注意缓存一致性问题。

