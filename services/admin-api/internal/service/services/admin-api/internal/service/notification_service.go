package service

import (
	"context"
	"fmt"
	"time"

	"goweb/pkg/logger"
	"goweb/pkg/notification"
	pkgRedis "goweb/pkg/redis"
	"goweb/pkg/websocket"

	"gorm.io/gorm"
)

// NotificationService 通知服务
type NotificationService struct {
	db           *gorm.DB
	notifyClient *notification.Client
	logger       logger.Logger
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB, redisClient *pkgRedis.Client, log logger.Logger) *NotificationService {
	return &NotificationService{
		db:           db,
		notifyClient: notification.NewClient(redisClient),
		logger:       log,
	}
}

// SendLoginNotification 发送登录通知
func (s *NotificationService) SendLoginNotification(ctx context.Context, userID, username string) error {
	notification := &websocket.NotificationMessage{
		Title:    "用户登录提醒",
		Body:     fmt.Sprintf("%s 于 %s 登录了系统", username, time.Now().Format("2006-01-02 15:04:05")),
		Icon:     "User",
		Link:     "/system/users",
		Category: "system",
	}

	s.logger.Info("发送登录通知", "username", username)
	return s.notifyClient.BroadcastNotification(ctx, notification)
}

// SendOrderNotification 发送订单通知
func (s *NotificationService) SendOrderNotification(ctx context.Context, userID, orderID, orderStatus string) error {
	var title, body, icon string

	switch orderStatus {
	case "created":
		title = "新订单通知"
		body = fmt.Sprintf("订单 #%s 已创建，等待处理", orderID)
		icon = "ShoppingCart"
	case "paid":
		title = "订单支付通知"
		body = fmt.Sprintf("订单 #%s 已支付，等待发货", orderID)
		icon = "CreditCard"
	case "shipped":
		title = "订单发货通知"
		body = fmt.Sprintf("订单 #%s 已发货", orderID)
		icon = "Truck"
	case "completed":
		title = "订单完成通知"
		body = fmt.Sprintf("订单 #%s 已完成", orderID)
		icon = "Check"
	default:
		title = "订单状态更新"
		body = fmt.Sprintf("订单 #%s 状态已更新为: %s", orderID, orderStatus)
		icon = "Bell"
	}

	notification := &websocket.NotificationMessage{
		Title:    title,
		Body:     body,
		Icon:     icon,
		Link:     "/orders/" + orderID,
		Category: "order",
		Data: map[string]interface{}{
			"order_id":     orderID,
			"order_status": orderStatus,
		},
	}

	s.logger.Info("发送订单通知", "order_id", orderID, "status", orderStatus, "user_id", userID)

	// 发送给指定用户
	if userID != "" {
		return s.notifyClient.SendNotification(ctx, userID, notification)
	}

	// 广播给所有管理员
	return s.notifyClient.BroadcastNotification(ctx, notification)
}

// SendSystemNotification 发送系统通知
func (s *NotificationService) SendSystemNotification(ctx context.Context, title, body, icon, link string) error {
	notification := &websocket.NotificationMessage{
		Title:    title,
		Body:     body,
		Icon:     icon,
		Link:     link,
		Category: "system",
	}

	s.logger.Info("发送系统通知", "title", title)
	return s.notifyClient.BroadcastNotification(ctx, notification)
}

// SendUserNotification 发送用户通知
func (s *NotificationService) SendUserNotification(ctx context.Context, userID, title, body, icon, link string) error {
	notification := &websocket.NotificationMessage{
		Title:    title,
		Body:     body,
		Icon:     icon,
		Link:     link,
		Category: "user",
	}

	s.logger.Info("发送用户通知", "user_id", userID, "title", title)
	return s.notifyClient.SendNotification(ctx, userID, notification)
}