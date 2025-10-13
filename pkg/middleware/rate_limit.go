package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	limiter *rate.Limiter
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rps int, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Allow 检查是否允许请求
func (t *TokenBucketLimiter) Allow() bool {
	return t.limiter.Allow()
}

// Wait 等待直到允许请求
func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
	return t.limiter.Wait(ctx)
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	window   time.Duration
	limit    int
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(window time.Duration, limit int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

// Allow 检查是否允许请求
func (s *SlidingWindowLimiter) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	// 清理过期请求
	if requests, exists := s.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		s.requests[key] = validRequests
	}

	// 检查是否超过限制
	if len(s.requests[key]) >= s.limit {
		return false
	}

	// 记录当前请求
	s.requests[key] = append(s.requests[key], now)
	return true
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Type    string                    // "token_bucket" 或 "sliding_window"
	RPS     int                       // 每秒请求数
	Burst   int                       // 突发请求数
	Window  time.Duration             // 滑动窗口大小
	Limit   int                       // 窗口内最大请求数
	KeyFunc func(*gin.Context) string // 限流键生成函数
	Message string                    // 限流时的错误消息
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Type:    "token_bucket",
		RPS:     100,
		Burst:   200,
		KeyFunc: func(c *gin.Context) string { return c.ClientIP() },
		Message: "Too many requests",
	}
}

// RateLimit 限流中间件
func RateLimit(config *RateLimitConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	var limiter RateLimiter
	var slidingLimiter *SlidingWindowLimiter

	switch config.Type {
	case "sliding_window":
		slidingLimiter = NewSlidingWindowLimiter(config.Window, config.Limit)
	default:
		limiter = NewTokenBucketLimiter(config.RPS, config.Burst)
	}

	return func(c *gin.Context) {
		var allowed bool

		if slidingLimiter != nil {
			key := config.KeyFunc(c)
			allowed = slidingLimiter.Allow(key)
		} else {
			allowed = limiter.Allow()
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": config.Message,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimit IP限流中间件
func IPRateLimit(rps int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	mu := sync.RWMutex{}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests from this IP",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimit 用户限流中间件（需要JWT中间件先执行）
func UserRateLimit(rps int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	mu := sync.RWMutex{}

	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userIDStr := fmt.Sprintf("%v", userID)

		mu.Lock()
		limiter, exists := limiters[userIDStr]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[userIDStr] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests from this user",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
