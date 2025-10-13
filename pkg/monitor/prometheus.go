package monitor

import (
	"time"

	"goweb/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 监控指标
type Metrics struct {
	// HTTP 请求指标
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	// 业务指标
	businessCounter   *prometheus.CounterVec
	businessGauge     *prometheus.GaugeVec
	businessHistogram *prometheus.HistogramVec

	// 数据库指标
	dbConnections   prometheus.Gauge
	dbQueryDuration *prometheus.HistogramVec
	dbQueryTotal    *prometheus.CounterVec

	// 缓存指标
	cacheHits       *prometheus.CounterVec
	cacheMisses     *prometheus.CounterVec
	cacheOperations *prometheus.CounterVec

	// 自定义指标
	customCounters   map[string]*prometheus.CounterVec
	customGauges     map[string]*prometheus.GaugeVec
	customHistograms map[string]*prometheus.HistogramVec

	logger logger.Logger
}

// NewMetrics 创建监控指标
func NewMetrics(serviceName string, log logger.Logger) *Metrics {
	return &Metrics{
		// HTTP 请求指标
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "service"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status", "service"},
		),
		httpRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path", "service"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path", "status", "service"},
		),

		// 业务指标
		businessCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_operations_total",
				Help: "Total number of business operations",
			},
			[]string{"operation", "status", "service"},
		),
		businessGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "business_metrics",
				Help: "Business metrics",
			},
			[]string{"metric", "service"},
		),
		businessHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "business_operation_duration_seconds",
				Help:    "Business operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "service"},
		),

		// 数据库指标
		dbConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table", "service"},
		),
		dbQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table", "status", "service"},
		),

		// 缓存指标
		cacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type", "key_pattern", "service"},
		),
		cacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type", "key_pattern", "service"},
		),
		cacheOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_operations_total",
				Help: "Total number of cache operations",
			},
			[]string{"operation", "cache_type", "status", "service"},
		),

		customCounters:   make(map[string]*prometheus.CounterVec),
		customGauges:     make(map[string]*prometheus.GaugeVec),
		customHistograms: make(map[string]*prometheus.HistogramVec),
		logger:           log,
	}
}

// HTTPMiddleware HTTP 监控中间件
func (m *Metrics) HTTPMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method

		// 记录请求大小
		requestSize := c.Request.ContentLength
		if requestSize > 0 {
			m.httpRequestSize.WithLabelValues(method, path, serviceName).Observe(float64(requestSize))
		}

		c.Next()

		// 记录响应指标
		status := c.Writer.Status()
		duration := time.Since(start).Seconds()
		responseSize := c.Writer.Size()

		m.httpRequestsTotal.WithLabelValues(method, path, statusCode(status), serviceName).Inc()
		m.httpRequestDuration.WithLabelValues(method, path, statusCode(status), serviceName).Observe(duration)
		m.httpResponseSize.WithLabelValues(method, path, statusCode(status), serviceName).Observe(float64(responseSize))
	}
}

// RecordBusinessOperation 记录业务操作
func (m *Metrics) RecordBusinessOperation(operation, status, serviceName string) {
	m.businessCounter.WithLabelValues(operation, status, serviceName).Inc()
}

// RecordBusinessDuration 记录业务操作耗时
func (m *Metrics) RecordBusinessDuration(operation, serviceName string, duration time.Duration) {
	m.businessHistogram.WithLabelValues(operation, serviceName).Observe(duration.Seconds())
}

// SetBusinessGauge 设置业务指标
func (m *Metrics) SetBusinessGauge(metric, serviceName string, value float64) {
	m.businessGauge.WithLabelValues(metric, serviceName).Set(value)
}

// RecordDBCall 记录数据库调用
func (m *Metrics) RecordDBCall(operation, table, status, serviceName string, duration time.Duration) {
	m.dbQueryTotal.WithLabelValues(operation, table, status, serviceName).Inc()
	m.dbQueryDuration.WithLabelValues(operation, table, serviceName).Observe(duration.Seconds())
}

// SetDBConnections 设置数据库连接数
func (m *Metrics) SetDBConnections(count float64) {
	m.dbConnections.Set(count)
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit(cacheType, keyPattern, serviceName string) {
	m.cacheHits.WithLabelValues(cacheType, keyPattern, serviceName).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss(cacheType, keyPattern, serviceName string) {
	m.cacheMisses.WithLabelValues(cacheType, keyPattern, serviceName).Inc()
}

// RecordCacheOperation 记录缓存操作
func (m *Metrics) RecordCacheOperation(operation, cacheType, status, serviceName string) {
	m.cacheOperations.WithLabelValues(operation, cacheType, status, serviceName).Inc()
}

// CreateCustomCounter 创建自定义计数器
func (m *Metrics) CreateCustomCounter(name, help string, labels []string) *prometheus.CounterVec {
	if counter, exists := m.customCounters[name]; exists {
		return counter
	}

	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	m.customCounters[name] = counter
	return counter
}

// CreateCustomGauge 创建自定义仪表
func (m *Metrics) CreateCustomGauge(name, help string, labels []string) *prometheus.GaugeVec {
	if gauge, exists := m.customGauges[name]; exists {
		return gauge
	}

	gauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	m.customGauges[name] = gauge
	return gauge
}

// CreateCustomHistogram 创建自定义直方图
func (m *Metrics) CreateCustomHistogram(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	if histogram, exists := m.customHistograms[name]; exists {
		return histogram
	}

	opts := prometheus.HistogramOpts{
		Name: name,
		Help: help,
	}
	if buckets != nil {
		opts.Buckets = buckets
	}

	histogram := promauto.NewHistogramVec(opts, labels)
	m.customHistograms[name] = histogram
	return histogram
}

// statusCode 转换状态码为字符串
func statusCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "2xx"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	case code >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
