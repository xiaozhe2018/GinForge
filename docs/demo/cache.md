# 缓存示例（cache）

```go
package demo

import (
    "context"
    "time"
    "goweb/pkg/config"
    "goweb/pkg/redis"
)

func ExampleCache() error {
    cfg := config.New()
    redisManager := redis.NewManager(cfg.GetRedisConfig(), log)
    cm := redisManager.GetCache()

    ctx := context.Background()
    if err := cm.Set(ctx, "k1", []byte("v1"), time.Minute); err != nil { return err }
    b, err := cm.Get(ctx, "k1")
    if err != nil { return err }
    _ = b

    // 删除
    _ = cm.Delete(ctx, "k1")
    return nil
}
```

提示：多实例部署时建议启用 Redis 并设置合理的过期时间与 key 前缀。
