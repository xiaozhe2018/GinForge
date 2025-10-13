package monitor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/redis"

	"github.com/gin-gonic/gin"
	redisClient "github.com/redis/go-redis/v9"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	checks map[string]HealthCheck
	logger logger.Logger
}

// HealthCheck 健康检查接口
type HealthCheck interface {
	Name() string
	Check(ctx context.Context) HealthStatus
}

// HealthStatus 健康状态
type HealthStatus struct {
	Name      string         `json:"name"`
	Status    string         `json:"status"` // healthy, unhealthy, degraded
	Message   string         `json:"message,omitempty"`
	Duration  time.Duration  `json:"duration"`
	Details   map[string]any `json:"details,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(log logger.Logger) *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
		logger: log,
	}
}

// Register 注册健康检查
func (hc *HealthChecker) Register(check HealthCheck) {
	hc.checks[check.Name()] = check
	hc.logger.Info("health check registered", "name", check.Name())
}

// CheckAll 检查所有健康状态
func (hc *HealthChecker) CheckAll(ctx context.Context) map[string]HealthStatus {
	results := make(map[string]HealthStatus)

	for name, check := range hc.checks {
		start := time.Now()
		status := check.Check(ctx)
		status.Duration = time.Since(start)
		status.Timestamp = time.Now()
		results[name] = status
	}

	return results
}

// GetOverallStatus 获取整体健康状态
func (hc *HealthChecker) GetOverallStatus(ctx context.Context) HealthStatus {
	results := hc.CheckAll(ctx)

	overall := HealthStatus{
		Name:      "overall",
		Status:    "healthy",
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	// 转换 results 为 any 类型
	for k, v := range results {
		overall.Details[k] = v
	}

	unhealthyCount := 0
	degradedCount := 0

	for _, status := range results {
		switch status.Status {
		case "unhealthy":
			unhealthyCount++
		case "degraded":
			degradedCount++
		}
	}

	if unhealthyCount > 0 {
		overall.Status = "unhealthy"
		overall.Message = fmt.Sprintf("%d services unhealthy", unhealthyCount)
	} else if degradedCount > 0 {
		overall.Status = "degraded"
		overall.Message = fmt.Sprintf("%d services degraded", degradedCount)
	}

	return overall
}

// GinHandler 返回 Gin 处理函数
func (hc *HealthChecker) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		overall := hc.GetOverallStatus(ctx)

		statusCode := 200
		if overall.Status == "unhealthy" {
			statusCode = 503
		} else if overall.Status == "degraded" {
			statusCode = 200 // 降级但可用
		}

		c.JSON(statusCode, overall)
	}
}

// 内置健康检查实现

// DatabaseHealthCheck 数据库健康检查
type DatabaseHealthCheck struct {
	db     *sql.DB
	name   string
	config *config.DatabaseConfig
}

// NewDatabaseHealthCheck 创建数据库健康检查
func NewDatabaseHealthCheck(db *sql.DB, name string, config *config.DatabaseConfig) *DatabaseHealthCheck {
	return &DatabaseHealthCheck{
		db:     db,
		name:   name,
		config: config,
	}
}

func (d *DatabaseHealthCheck) Name() string {
	return d.name
}

func (d *DatabaseHealthCheck) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Name:   d.name,
		Status: "healthy",
		Details: map[string]any{
			"driver":   d.config.Driver,
			"host":     d.config.Host,
			"port":     d.config.Port,
			"database": d.config.Database,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
		return status
	}

	// 检查连接池状态
	stats := d.db.Stats()
	status.Details["open_connections"] = stats.OpenConnections
	status.Details["in_use"] = stats.InUse
	status.Details["idle"] = stats.Idle
	status.Details["wait_count"] = stats.WaitCount
	status.Details["wait_duration"] = stats.WaitDuration.String()

	// 如果等待连接过多，标记为降级
	if stats.WaitCount > 10 {
		status.Status = "degraded"
		status.Message = "high connection wait count"
	}

	return status
}

// RedisHealthCheck Redis 健康检查
type RedisHealthCheck struct {
	client redisClient.UniversalClient
	name   string
	config *config.RedisConfig
}

// NewRedisHealthCheck 创建 Redis 健康检查
func NewRedisHealthCheck(client redisClient.UniversalClient, name string, config *config.RedisConfig) *RedisHealthCheck {
	return &RedisHealthCheck{
		client: client,
		name:   name,
		config: config,
	}
}

func (r *RedisHealthCheck) Name() string {
	return r.name
}

func (r *RedisHealthCheck) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Name:   r.name,
		Status: "healthy",
		Details: map[string]any{
			"host":     r.config.Host,
			"port":     r.config.Port,
			"database": r.config.Database,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	start := time.Now()
	pong, err := r.client.Ping(ctx).Result()
	status.Duration = time.Since(start)

	if err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
		return status
	}

	status.Details["ping"] = pong
	status.Details["response_time"] = status.Duration.String()

	// 如果响应时间过长，标记为降级
	if status.Duration > 100*time.Millisecond {
		status.Status = "degraded"
		status.Message = "slow response time"
	}

	return status
}

// CacheHealthCheck 缓存健康检查
type CacheHealthCheck struct {
	cache  *redis.Manager
	name   string
	config *config.RedisConfig
}

// NewCacheHealthCheck 创建缓存健康检查
func NewCacheHealthCheck(cache *redis.Manager, name string, config *config.RedisConfig) *CacheHealthCheck {
	return &CacheHealthCheck{
		cache:  cache,
		name:   name,
		config: config,
	}
}

func (c *CacheHealthCheck) Name() string {
	return c.name
}

func (c *CacheHealthCheck) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Name:   c.name,
		Status: "healthy",
		Details: map[string]any{
			"type": "cache",
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 测试缓存操作
	testKey := "health_check_test"
	testValue := []byte("test")

	start := time.Now()
	err := c.cache.Set(ctx, testKey, testValue, time.Second)
	status.Duration = time.Since(start)

	if err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
		return status
	}

	// 测试读取
	var result []byte
	err = c.cache.Get(ctx, testKey, &result)
	if err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
		return status
	}

	// 清理测试数据
	c.cache.Delete(ctx, testKey)

	status.Details["response_time"] = status.Duration.String()

	// 如果响应时间过长，标记为降级
	if status.Duration > 50*time.Millisecond {
		status.Status = "degraded"
		status.Message = "slow cache response time"
	}

	return status
}

// ServiceHealthCheck 服务健康检查
type ServiceHealthCheck struct {
	name        string
	checkFunc   func(ctx context.Context) error
	description string
}

// NewServiceHealthCheck 创建服务健康检查
func NewServiceHealthCheck(name, description string, checkFunc func(ctx context.Context) error) *ServiceHealthCheck {
	return &ServiceHealthCheck{
		name:        name,
		checkFunc:   checkFunc,
		description: description,
	}
}

func (s *ServiceHealthCheck) Name() string {
	return s.name
}

func (s *ServiceHealthCheck) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Name:   s.name,
		Status: "healthy",
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()
	err := s.checkFunc(ctx)
	status.Duration = time.Since(start)

	if err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
	} else {
		status.Details = map[string]any{
			"response_time": status.Duration.String(),
		}
	}

	return status
}
