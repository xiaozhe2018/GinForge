package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Duration time.Duration
	KeyFunc  func(*gin.Context) string
	SkipFunc func(*gin.Context) bool
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Duration: 5 * time.Minute,
		KeyFunc:  defaultCacheKeyFunc,
		SkipFunc: func(c *gin.Context) bool { return false },
	}
}

// defaultCacheKeyFunc 默认缓存键生成函数
func defaultCacheKeyFunc(c *gin.Context) string {
	// 使用请求方法和路径作为缓存键
	return fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
}

// CacheEntry 缓存条目
type CacheEntry struct {
	Data      []byte      `json:"data"`
	Headers   http.Header `json:"headers"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
}

// MemoryCache 内存缓存
type MemoryCache struct {
	entries map[string]*CacheEntry
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		entries: make(map[string]*CacheEntry),
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (*CacheEntry, bool) {
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	return entry, true
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, entry *CacheEntry) {
	c.entries[key] = entry
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	delete(c.entries, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.entries = make(map[string]*CacheEntry)
}

// 全局内存缓存实例
var globalCache = NewMemoryCache()

// Cache 缓存中间件
func Cache(config *CacheConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCacheConfig()
	}

	return func(c *gin.Context) {
		// 检查是否跳过缓存
		if config.SkipFunc(c) {
			c.Next()
			return
		}

		// 生成缓存键
		key := config.KeyFunc(c)

		// 尝试从缓存获取
		if entry, exists := globalCache.Get(key); exists {
			// 检查是否过期
			if time.Since(entry.Timestamp) < config.Duration {
				// 设置响应头
				for k, v := range entry.Headers {
					for _, val := range v {
						c.Header(k, val)
					}
				}
				c.Header("X-Cache", "HIT")
				c.Data(entry.Status, "application/json", entry.Data)
				c.Abort()
				return
			}
		}

		// 使用自定义响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		// 只缓存成功的响应
		if writer.Status() >= 200 && writer.Status() < 300 {
			// 创建缓存条目
			entry := &CacheEntry{
				Data:      writer.body.Bytes(),
				Headers:   writer.Header().Clone(),
				Status:    writer.Status(),
				Timestamp: time.Now(),
			}

			// 存储到缓存
			globalCache.Set(key, entry)
		}

		// 设置缓存状态头
		c.Header("X-Cache", "MISS")
	}
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应体
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// WriteString 写入字符串
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// CacheByQuery 基于查询参数的缓存中间件
func CacheByQuery(duration time.Duration) gin.HandlerFunc {
	return Cache(&CacheConfig{
		Duration: duration,
		KeyFunc: func(c *gin.Context) string {
			// 包含查询参数的缓存键
			query := c.Request.URL.RawQuery
			if query != "" {
				return fmt.Sprintf("%s:%s?%s", c.Request.Method, c.Request.URL.Path, query)
			}
			return fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		},
	})
}

// CacheByUser 基于用户的缓存中间件
func CacheByUser(duration time.Duration) gin.HandlerFunc {
	return Cache(&CacheConfig{
		Duration: duration,
		KeyFunc: func(c *gin.Context) string {
			// 包含用户ID的缓存键
			userID, exists := c.Get("user_id")
			if exists {
				return fmt.Sprintf("%s:%s:user:%v", c.Request.Method, c.Request.URL.Path, userID)
			}
			return fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		},
	})
}

// CacheControl 缓存控制中间件
func CacheControl(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Next()
	}
}

// NoCache 禁用缓存中间件
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// ETag 生成ETag
func ETag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只对GET请求生成ETag
		if c.Request.Method != "GET" {
			return
		}

		// 生成ETag（基于请求路径）
		etag := fmt.Sprintf(`"%x"`, md5.Sum([]byte(c.Request.URL.Path)))

		// 设置ETag头
		c.Header("ETag", etag)

		// 检查If-None-Match头
		if match := c.GetHeader("If-None-Match"); match == etag {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}
	}
}

// ClearCache 清空缓存中间件
func ClearCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		globalCache.Clear()
		c.JSON(http.StatusOK, gin.H{"message": "Cache cleared"})
	}
}
