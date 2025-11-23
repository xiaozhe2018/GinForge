package handler

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/logger"
	"goweb/services/admin-api/internal/service"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationService *service.NotificationService
	logger              logger.Logger
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger.New("notification-handler", "info", "stdout", ""),
	}
}

// SetLogger 设置日志器
func (h *NotificationHandler) SetLogger(log logger.Logger) {
	h.logger = log
}

// SendSystemNotification 发送系统通知
// @Summary 发送系统通知
// @Description 发送系统通知给所有在线用户
// @Tags 通知
// @Accept json
// @Produce json
// @Param notification body SystemNotificationRequest true "通知内容"
// @Success 200 {object} response.Response "发送成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/notifications/system [post]
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req SystemNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.notificationService.SendSystemNotification(ctx, req.Title, req.Body, req.Icon, req.Link); err != nil {
		h.logger.Error("发送系统通知失败", err)
		c.JSON(500, gin.H{"error": "发送通知失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "通知已发送",
	})
}

// SendUserNotification 发送用户通知
// @Summary 发送用户通知
// @Description 发送通知给指定用户
// @Tags 通知
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param notification body UserNotificationRequest true "通知内容"
// @Success 200 {object} response.Response "发送成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/notifications/users/{user_id} [post]
func (h *NotificationHandler) SendUserNotification(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "用户ID不能为空"})
		return
	}

	var req UserNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.notificationService.SendUserNotification(ctx, userID, req.Title, req.Body, req.Icon, req.Link); err != nil {
		h.logger.Error("发送用户通知失败", err)
		c.JSON(500, gin.H{"error": "发送通知失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "通知已发送",
	})
}

// SendOrderNotification 发送订单通知
// @Summary 发送订单通知
// @Description 发送订单状态变更通知
// @Tags 通知
// @Accept json
// @Produce json
// @Param notification body OrderNotificationRequest true "订单通知内容"
// @Success 200 {object} response.Response "发送成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/notifications/orders [post]
func (h *NotificationHandler) SendOrderNotification(c *gin.Context) {
	var req OrderNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.notificationService.SendOrderNotification(ctx, req.UserID, req.OrderID, req.OrderStatus); err != nil {
		h.logger.Error("发送订单通知失败", err)
		c.JSON(500, gin.H{"error": "发送通知失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "订单通知已发送",
	})
}

// SystemNotificationRequest 系统通知请求
type SystemNotificationRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Icon  string `json:"icon"`
	Link  string `json:"link"`
}

// UserNotificationRequest 用户通知请求
type UserNotificationRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Icon  string `json:"icon"`
	Link  string `json:"link"`
}

// OrderNotificationRequest 订单通知请求
type OrderNotificationRequest struct {
	UserID      string `json:"user_id"`
	OrderID     string `json:"order_id" binding:"required"`
	OrderStatus string `json:"order_status" binding:"required"`
}
