# 分布式锁

使用 Redis 实现分布式锁，解决并发问题。

## 🎯 为什么需要分布式锁？

### 解决的问题

在分布式系统中，多个服务实例可能同时访问共享资源：

```
Instance 1                Instance 2
    ↓                        ↓
读取库存: 10              读取库存: 10
    ↓                        ↓
库存 - 1 = 9             库存 - 1 = 9
    ↓                        ↓
保存: 9                   保存: 9
    ↓                        ↓
结果：库存应该是 8，但实际是 9（数据错误！）
```

### 使用分布式锁后

```
Instance 1                Instance 2
    ↓                        ↓
获取锁成功                获取锁失败（等待）
    ↓                        ↓
读取库存: 10              等待...
    ↓ 
库存 - 1 = 9
    ↓
保存: 9
    ↓
释放锁
    ↓                        ↓
                        获取锁成功
                             ↓
                        读取库存: 9
                             ↓
                        库存 - 1 = 8
                             ↓
                        保存: 8
                             ↓
                        释放锁

结果：库存正确为 8 ✅
```

## 📝 基础使用

### 简单锁

```go
import (
    "context"
    "time"
    "goweb/pkg/redis"
)

func updateInventory(redisClient *redis.Client, productID string) error {
    ctx := context.Background()
    lockKey := fmt.Sprintf("lock:inventory:%s", productID)
    
    // 获取锁（10秒超时）
    locked, err := redisClient.Lock(ctx, lockKey, 10*time.Second)
    if err != nil {
        return err
    }
    if !locked {
        return errors.New("获取锁失败，请稍后重试")
    }
    
    // 确保释放锁
    defer redisClient.Unlock(ctx, lockKey)
    
    // 执行业务逻辑
    inventory := getInventory(productID)
    inventory--
    saveInventory(productID, inventory)
    
    return nil
}
```

### 带重试的锁

```go
func updateInventoryWithRetry(redisClient *redis.Client, productID string) error {
    ctx := context.Background()
    lockKey := fmt.Sprintf("lock:inventory:%s", productID)
    maxRetries := 3
    retryDelay := 100 * time.Millisecond
    
    for i := 0; i < maxRetries; i++ {
        // 尝试获取锁
        locked, err := redisClient.Lock(ctx, lockKey, 10*time.Second)
        if err != nil {
            return err
        }
        
        if locked {
            defer redisClient.Unlock(ctx, lockKey)
            
            // 执行业务逻辑
            inventory := getInventory(productID)
            inventory--
            saveInventory(productID, inventory)
            
            return nil
        }
        
        // 没有获取到锁，等待后重试
        if i < maxRetries-1 {
            time.Sleep(retryDelay * time.Duration(i+1))
        }
    }
    
    return errors.New("获取锁失败，请稍后重试")
}
```

## 🔐 Lock 包使用

GinForge 提供了 `pkg/redis/lock.go` 封装：

```go
// 创建分布式锁
func (c *Client) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    // 使用 SET NX EX 命令
    // NX: 只在键不存在时设置
    // EX: 设置过期时间
    result, err := c.client.SetNX(ctx, key, "1", ttl).Result()
    return result, err
}

// 释放锁
func (c *Client) Unlock(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Error()
}

// 检查锁是否存在
func (c *Client) IsLocked(ctx context.Context, key string) (bool, error) {
    return c.Exists(ctx, key)
}
```

## 🎨 实战案例

### 案例 1：秒杀场景

```go
func (s *SeckillService) BuyProduct(userID, productID string) error {
    ctx := context.Background()
    lockKey := fmt.Sprintf("lock:seckill:%s", productID)
    
    // 获取锁（3秒超时）
    locked, err := s.redisClient.Lock(ctx, lockKey, 3*time.Second)
    if err != nil {
        return err
    }
    if !locked {
        return errors.New("商品太火爆，请稍后再试")
    }
    defer s.redisClient.Unlock(ctx, lockKey)
    
    // 1. 检查库存
    stock := s.getStock(productID)
    if stock <= 0 {
        return errors.New("商品已售罄")
    }
    
    // 2. 检查用户是否已购买
    if s.hasPurchased(userID, productID) {
        return errors.New("您已经购买过了")
    }
    
    // 3. 扣减库存
    s.reduceStock(productID, 1)
    
    // 4. 创建订单
    order := s.createOrder(userID, productID)
    
    s.logger.Info("秒杀成功", "user_id", userID, "product_id", productID, "order_id", order.ID)
    return nil
}
```

### 案例 2：防止重复提交

```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }
    
    userID := c.GetString("user_id")
    
    // 防重复提交锁（使用请求ID）
    requestID := c.GetString("request_id")
    lockKey := fmt.Sprintf("lock:order:create:%s", requestID)
    
    locked, _ := h.redisClient.Lock(c.Request.Context(), lockKey, 30*time.Second)
    if !locked {
        response.Error(c, 400, "请勿重复提交")
        return
    }
    defer h.redisClient.Unlock(c.Request.Context(), lockKey)
    
    // 创建订单
    order, err := h.orderService.CreateOrder(userID, &req)
    if err != nil {
        response.InternalError(c, "创建订单失败")
        return
    }
    
    response.Success(c, order)
}
```

### 案例 3：定时任务防并发

```go
// 定时任务：每分钟统计数据
func (s *StatsService) CalculateDailyStats() error {
    ctx := context.Background()
    lockKey := "lock:stats:daily"
    
    // 获取锁（5分钟超时）
    locked, err := s.redisClient.Lock(ctx, lockKey, 5*time.Minute)
    if err != nil {
        return err
    }
    if !locked {
        s.logger.Info("其他实例正在执行统计任务，跳过")
        return nil
    }
    defer s.redisClient.Unlock(ctx, lockKey)
    
    // 执行统计
    s.logger.Info("开始统计每日数据")
    stats := s.calculate()
    s.saveStats(stats)
    s.logger.Info("每日数据统计完成")
    
    return nil
}
```

## ⚠️ 注意事项

### 1. 避免死锁

```go
// ✅ 正确：使用 defer 确保锁被释放
func updateData(redisClient *redis.Client, key string) error {
    locked, _ := redisClient.Lock(ctx, lockKey, 10*time.Second)
    if locked {
        defer redisClient.Unlock(ctx, lockKey)  // 确保释放
        // 业务逻辑...
    }
    return nil
}

// ❌ 错误：忘记释放锁
func updateDataBad(redisClient *redis.Client, key string) error {
    locked, _ := redisClient.Lock(ctx, lockKey, 10*time.Second)
    if locked {
        // 业务逻辑...
        // 忘记释放锁，导致死锁！
    }
    return nil
}
```

### 2. 设置合理的超时时间

```go
// ✅ 根据业务设置合理的超时
redisClient.Lock(ctx, lockKey, 3*time.Second)   // 快速操作
redisClient.Lock(ctx, lockKey, 30*time.Second)  // 普通操作
redisClient.Lock(ctx, lockKey, 5*time.Minute)   // 复杂操作

// ❌ 超时时间过短
redisClient.Lock(ctx, lockKey, 100*time.Millisecond)  // 可能业务还没执行完就过期

// ❌ 超时时间过长
redisClient.Lock(ctx, lockKey, 1*time.Hour)  // 占用资源太久
```

### 3. 锁的粒度

```go
// ✅ 细粒度锁（推荐）
lockKey := fmt.Sprintf("lock:inventory:%s", productID)  // 每个商品一个锁

// ❌ 粗粒度锁（性能差）
lockKey := "lock:inventory"  // 所有商品共用一个锁
```

## 🔧 高级功能

### 可重入锁

```go
type ReentrantLock struct {
    redis     *redis.Client
    lockID    string
    holdCount int
}

func (l *ReentrantLock) Lock(ctx context.Context, key string, ttl time.Duration) error {
    if l.holdCount > 0 {
        // 已经持有锁，增加计数
        l.holdCount++
        return nil
    }
    
    // 获取锁
    locked, err := l.redis.Lock(ctx, key, ttl)
    if err != nil || !locked {
        return errors.New("获取锁失败")
    }
    
    l.lockID = key
    l.holdCount = 1
    return nil
}

func (l *ReentrantLock) Unlock(ctx context.Context) error {
    if l.holdCount <= 0 {
        return errors.New("未持有锁")
    }
    
    l.holdCount--
    
    // 计数归零才真正释放锁
    if l.holdCount == 0 {
        return l.redis.Unlock(ctx, l.lockID)
    }
    
    return nil
}
```

### 读写锁

```go
// 读锁（共享锁）
func (c *Client) RLock(ctx context.Context, key string, ttl time.Duration) error {
    rkey := fmt.Sprintf("%s:read", key)
    // 读锁允许多个持有者
    return c.client.Incr(ctx, rkey).Err()
}

// 写锁（排他锁）
func (c *Client) WLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    rkey := fmt.Sprintf("%s:read", key)
    wkey := fmt.Sprintf("%s:write", key)
    
    // 检查是否有读锁
    readers, _ := c.client.Get(ctx, rkey).Int()
    if readers > 0 {
        return false, nil
    }
    
    // 获取写锁
    return c.SetNX(ctx, wkey, "1", ttl).Result()
}
```

## 📚 完整示例

查看完整实现：

- **分布式锁**: `pkg/redis/lock.go`
- **Lock 使用**: `pkg/redis/client.go`

## 🎯 下一步

- [消息队列](./message-queue) - 异步任务处理
- [性能优化](./performance) - 系统优化技巧

---

**提示**: 分布式锁务必注意超时时间和异常处理，避免死锁！

