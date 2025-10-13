# 配置中心使用指南

GinForge 提供了强大的配置中心功能，支持动态配置更新、配置监听、配置验证等特性。

## 功能特性

- 🔄 **动态配置更新**: 支持运行时配置更新，无需重启服务
- 👂 **配置监听**: 监听配置变更事件，自动通知相关组件
- ✅ **配置验证**: 支持配置值验证，确保配置正确性
- 💾 **配置备份**: 自动备份配置，支持配置恢复
- 🔍 **配置查询**: 支持配置查询和管理
- 📊 **配置统计**: 提供配置使用统计信息

## 快速开始

### 1. 初始化配置中心

```go
package main

import (
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/redis"
)

func main() {
    // 加载配置
    cfg := config.New()
    log := logger.New("config-center", "info")
    
    // 初始化Redis
    redisMgr, err := redis.NewManager(cfg)
    if err != nil {
        log.Fatal("Redis初始化失败", err)
    }
    
    // 创建配置中心
    configCenter := config.NewConfigCenter(cfg, log, redisMgr)
    
    // 使用配置中心
    port := configCenter.GetInt("app.port")
    log.Info("应用端口", "port", port)
}
```

### 2. 动态配置更新

```go
// 更新单个配置
err := configCenter.SetConfig("app.port", 8081)
if err != nil {
    log.Error("配置更新失败", err)
}

// 批量更新配置
configs := map[string]interface{}{
    "app.port": 8081,
    "log.level": "debug",
    "database.host": "localhost",
}
err := configCenter.UpdateConfigs(configs)
if err != nil {
    log.Error("批量配置更新失败", err)
}
```

### 3. 配置监听

```go
// 注册配置监听器
configCenter.RegisterWatcher("app.port", &PortWatcher{})

// 实现配置监听器
type PortWatcher struct{}

func (w *PortWatcher) OnConfigChange(key string, oldValue, newValue interface{}) {
    log.Info("配置已变更", 
        "key", key, 
        "old_value", oldValue, 
        "new_value", newValue)
    
    // 处理配置变更
    if key == "app.port" {
        // 重启HTTP服务器
        restartHTTPServer(newValue.(int))
    }
}

// 简单配置监听
configCenter.WatchConfig("log.level", func(key string, oldValue, newValue interface{}) {
    log.Info("日志级别已变更", "old", oldValue, "new", newValue)
    // 更新日志级别
    updateLogLevel(newValue.(string))
})
```

### 4. 配置验证

```go
// 创建带验证的配置中心
configCenter := config.NewConfigCenterWithValidation(cfg, log, redisMgr)

// 注册配置验证器
configCenter.RegisterValidator("app.port", &PortValidator{})
configCenter.RegisterValidator("log.level", &LogLevelValidator{})

// 实现配置验证器
type PortValidator struct{}

func (v *PortValidator) Validate(key string, value interface{}) error {
    port, ok := value.(int)
    if !ok {
        return fmt.Errorf("端口必须是整数")
    }
    
    if port < 1 || port > 65535 {
        return fmt.Errorf("端口必须在1-65535范围内")
    }
    
    return nil
}

type LogLevelValidator struct{}

func (v *LogLevelValidator) Validate(key string, value interface{}) error {
    level, ok := value.(string)
    if !ok {
        return fmt.Errorf("日志级别必须是字符串")
    }
    
    validLevels := []string{"debug", "info", "warn", "error", "fatal"}
    for _, validLevel := range validLevels {
        if level == validLevel {
            return nil
        }
    }
    
    return fmt.Errorf("无效的日志级别: %s", level)
}
```

### 5. 配置备份和恢复

```go
// 备份配置
backup, err := configCenter.BackupConfigs()
if err != nil {
    log.Error("配置备份失败", err)
} else {
    log.Info("配置已备份", "timestamp", backup.Timestamp)
}

// 恢复配置
err := configCenter.RestoreConfigs(backup)
if err != nil {
    log.Error("配置恢复失败", err)
} else {
    log.Info("配置已恢复")
}
```

## 高级用法

### 1. 配置变更事件处理

```go
// 监听配置变更事件
go func() {
    for {
        select {
        case event := <-configChangeEvents:
            handleConfigChange(event)
        }
    }
}()

func handleConfigChange(event config.ConfigChangeEvent) {
    log.Info("配置变更事件", 
        "key", event.Key,
        "old_value", event.OldValue,
        "new_value", event.NewValue,
        "timestamp", event.Timestamp)
    
    // 根据配置键处理不同逻辑
    switch event.Key {
    case "app.port":
        handlePortChange(event.OldValue, event.NewValue)
    case "log.level":
        handleLogLevelChange(event.OldValue, event.NewValue)
    case "database.host":
        handleDatabaseHostChange(event.OldValue, event.NewValue)
    }
}
```

### 2. 配置热重载

```go
// 实现配置热重载
type ConfigReloader struct {
    configCenter *config.ConfigCenter
    services     map[string]ReloadableService
}

type ReloadableService interface {
    ReloadConfig(configs map[string]interface{}) error
}

func (cr *ConfigReloader) Start() {
    // 监听配置变更
    cr.configCenter.WatchConfig("", func(key string, oldValue, newValue interface{}) {
        // 获取所有配置
        allConfigs := cr.configCenter.GetAllConfigs()
        
        // 通知所有服务重新加载配置
        for name, service := range cr.services {
            if err := service.ReloadConfig(allConfigs); err != nil {
                log.Error("服务配置重载失败", "service", name, "error", err)
            }
        }
    })
}
```

### 3. 配置管理API

```go
// 配置管理API
func setupConfigAPI(r *gin.Engine, configCenter *config.ConfigCenter) {
    api := r.Group("/api/v1/config")
    
    // 获取配置
    api.GET("/:key", func(c *gin.Context) {
        key := c.Param("key")
        value := configCenter.GetConfig(key)
        response.Success(c, gin.H{"key": key, "value": value})
    })
    
    // 设置配置
    api.POST("/:key", func(c *gin.Context) {
        key := c.Param("key")
        var req struct {
            Value interface{} `json:"value"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            response.BadRequest(c, err.Error())
            return
        }
        
        if err := configCenter.SetConfig(key, req.Value); err != nil {
            response.InternalError(c, err.Error())
            return
        }
        
        response.Success(c, gin.H{"message": "配置已更新"})
    })
    
    // 获取所有配置
    api.GET("/", func(c *gin.Context) {
        configs := configCenter.GetAllConfigs()
        response.Success(c, configs)
    })
    
    // 批量更新配置
    api.PUT("/", func(c *gin.Context) {
        var req map[string]interface{}
        if err := c.ShouldBindJSON(&req); err != nil {
            response.BadRequest(c, err.Error())
            return
        }
        
        if err := configCenter.UpdateConfigs(req); err != nil {
            response.InternalError(c, err.Error())
            return
        }
        
        response.Success(c, gin.H{"message": "配置已批量更新"})
    })
}
```

## 最佳实践

### 1. 配置分层

```go
// 配置分层管理
type ConfigLayer struct {
    Defaults map[string]interface{}
    Overrides map[string]interface{}
}

func (cl *ConfigLayer) Get(key string) interface{} {
    // 优先使用覆盖值
    if value, exists := cl.Overrides[key]; exists {
        return value
    }
    
    // 使用默认值
    if value, exists := cl.Defaults[key]; exists {
        return value
    }
    
    return nil
}
```

### 2. 配置缓存

```go
// 配置缓存
type ConfigCache struct {
    cache map[string]interface{}
    mu    sync.RWMutex
    ttl   time.Duration
}

func (cc *ConfigCache) Get(key string) (interface{}, bool) {
    cc.mu.RLock()
	defer cc.mu.RUnlock()
	
	value, exists := cc.cache[key]
	return value, exists
}

func (cc *ConfigCache) Set(key string, value interface{}) {
    cc.mu.Lock()
	defer cc.mu.Unlock()
	
	cc.cache[key] = value
}
```

### 3. 配置监控

```go
// 配置监控
type ConfigMonitor struct {
    configCenter *config.ConfigCenter
    metrics      *prometheus.CounterVec
}

func (cm *ConfigMonitor) Start() {
    // 监听配置变更
    cm.configCenter.WatchConfig("", func(key string, oldValue, newValue interface{}) {
        // 记录配置变更指标
        cm.metrics.WithLabelValues(key, "changed").Inc()
    })
}
```

## 故障排除

### 1. 配置更新失败

```go
// 检查Redis连接
if err := configCenter.SetConfig("test.key", "test.value"); err != nil {
    log.Error("配置更新失败", "error", err)
    
    // 检查Redis连接
    if err := redisMgr.Client.Ping(context.Background()).Err(); err != nil {
        log.Error("Redis连接失败", "error", err)
    }
}
```

### 2. 配置监听器不工作

```go
// 检查监听器注册
watchers := configCenter.GetWatchers("app.port")
if len(watchers) == 0 {
    log.Warn("没有注册的监听器", "key", "app.port")
}
```

### 3. 配置验证失败

```go
// 检查验证器
if err := configCenter.SetConfig("app.port", "invalid"); err != nil {
    log.Error("配置验证失败", "error", err)
    
    // 检查验证器是否注册
    if validator, exists := configCenter.GetValidator("app.port"); exists {
        log.Info("验证器已注册", "validator", validator)
    }
}
```

## 总结

GinForge 的配置中心提供了完整的配置管理解决方案，支持动态更新、监听、验证等高级功能。通过合理使用这些功能，可以构建更加灵活和可维护的微服务系统。
