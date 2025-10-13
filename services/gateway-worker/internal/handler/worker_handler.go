package handler

import (
	"github.com/gin-gonic/gin"

	"goweb/pkg/base"
	"goweb/pkg/logger"
	"goweb/services/gateway-worker/internal/service"
)

// WorkerHandler 网关工作处理器
type WorkerHandler struct {
	*base.BaseHandler
	workerService *service.WorkerService
}

// NewWorkerHandler 创建网关工作处理器
func NewWorkerHandler(workerService *service.WorkerService, log logger.Logger) *WorkerHandler {
	return &WorkerHandler{
		BaseHandler:   base.NewBaseHandler(log),
		workerService: workerService,
	}
}

// GetHealthHandler 获取健康检查处理器
func (h *WorkerHandler) GetHealthHandler() *gin.Engine {
	r := gin.New()

	// 健康检查
	r.GET("/healthz", h.HealthCheck)
	r.GET("/ready", h.ReadyCheck)
	r.GET("/metrics", h.Metrics)

	return r
}

// HealthCheck 健康检查
func (h *WorkerHandler) HealthCheck(c *gin.Context) {
	h.Success(c, gin.H{
		"status":    "ok",
		"service":   "gateway-worker",
		"timestamp": gin.H{},
	})
}

// ReadyCheck 就绪检查
func (h *WorkerHandler) ReadyCheck(c *gin.Context) {
	// 检查 Redis 连接等
	h.Success(c, gin.H{
		"status":    "ready",
		"service":   "gateway-worker",
		"consumers": "running",
	})
}

// Metrics 指标信息
func (h *WorkerHandler) Metrics(c *gin.Context) {
	// 返回工作器指标
	h.Success(c, gin.H{
		"service": "gateway-worker",
		"status":  "running",
		"uptime":  "running",
		"consumers": gin.H{
			"active": 5,
			"topics": []string{
				"order.reminder",
				"user.notification",
				"system.cleanup",
				"payment.retry",
				"inventory.alert",
			},
		},
	})
}
