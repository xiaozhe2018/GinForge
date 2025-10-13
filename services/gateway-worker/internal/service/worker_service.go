package service

import (
	"context"
	"sync"
	"time"

	"goweb/pkg/base"
	"goweb/pkg/logger"
	"goweb/pkg/redis"
)

// WorkerService 网关工作服务
type WorkerService struct {
	*base.BaseService
	redisManager *redis.Manager
	consumers    map[string]context.CancelFunc
	mutex        sync.RWMutex
}

// NewWorkerService 创建网关工作服务
func NewWorkerService(redisManager *redis.Manager, log logger.Logger) *WorkerService {
	return &WorkerService{
		BaseService:  base.NewBaseService(log),
		redisManager: redisManager,
		consumers:    make(map[string]context.CancelFunc),
	}
}

// StartConsumers 启动所有消息消费者
func (s *WorkerService) StartConsumers(ctx context.Context) error {
	queue := s.redisManager.GetQueue()

	// 启动延时消息处理器
	topics := []string{
		"order.reminder",
		"user.notification",
		"system.cleanup",
		"payment.retry",
		"inventory.alert",
	}

	for _, topic := range topics {
		if err := s.startConsumer(ctx, queue, topic); err != nil {
			s.LogError("failed to start consumer", err, "topic", topic)
			return err
		}
	}

	s.LogInfo("all consumers started", "count", len(topics))
	return nil
}

// startConsumer 启动单个消费者
func (s *WorkerService) startConsumer(ctx context.Context, queue *redis.RedisQueue, topic string) error {
	consumerCtx, cancel := context.WithCancel(ctx)

	s.mutex.Lock()
	s.consumers[topic] = cancel
	s.mutex.Unlock()

	go func() {
		defer cancel()

		s.LogInfo("starting consumer", "topic", topic)

		if err := queue.Subscribe(consumerCtx, topic, s.getMessageHandler(topic)); err != nil {
			s.LogError("consumer failed", err, "topic", topic)
		}

		s.LogInfo("consumer stopped", "topic", topic)
	}()

	return nil
}

// getMessageHandler 获取消息处理器
func (s *WorkerService) getMessageHandler(topic string) redis.MessageHandler {
	handlers := map[string]redis.MessageHandler{
		"order.reminder":    s.handleOrderReminder,
		"user.notification": s.handleUserNotification,
		"system.cleanup":    s.handleSystemCleanup,
		"payment.retry":     s.handlePaymentRetry,
		"inventory.alert":   s.handleInventoryAlert,
	}

	handler, exists := handlers[topic]
	if !exists {
		return s.handleDefault
	}

	return handler
}

// handleOrderReminder 处理订单提醒
func (s *WorkerService) handleOrderReminder(ctx context.Context, msg *redis.Message) error {
	orderID := msg.Data["order_id"].(string)
	userID := msg.Data["user_id"].(string)
	reminderType := msg.Data["type"].(string)

	s.LogInfo("processing order reminder", "order_id", orderID, "user_id", userID, "type", reminderType)

	// 发送提醒通知
	return s.sendOrderReminderNotification(ctx, orderID, userID, reminderType)
}

// handleUserNotification 处理用户通知
func (s *WorkerService) handleUserNotification(ctx context.Context, msg *redis.Message) error {
	userID := msg.Data["user_id"].(string)
	notificationType := msg.Data["type"].(string)
	content := msg.Data["content"].(string)

	s.LogInfo("processing user notification", "user_id", userID, "type", notificationType)

	// 发送用户通知
	return s.sendUserNotification(ctx, userID, notificationType, content)
}

// handleSystemCleanup 处理系统清理
func (s *WorkerService) handleSystemCleanup(ctx context.Context, msg *redis.Message) error {
	cleanupType := msg.Data["type"].(string)
	params := msg.Data["params"].(map[string]interface{})

	s.LogInfo("processing system cleanup", "type", cleanupType)

	// 执行系统清理
	return s.performSystemCleanup(ctx, cleanupType, params)
}

// handlePaymentRetry 处理支付重试
func (s *WorkerService) handlePaymentRetry(ctx context.Context, msg *redis.Message) error {
	paymentID := msg.Data["payment_id"].(string)
	orderID := msg.Data["order_id"].(string)
	retryCount := int(msg.Data["retry_count"].(float64))

	s.LogInfo("processing payment retry", "payment_id", paymentID, "order_id", orderID, "retry_count", retryCount)

	// 重试支付
	return s.retryPayment(ctx, paymentID, orderID, retryCount)
}

// handleInventoryAlert 处理库存告警
func (s *WorkerService) handleInventoryAlert(ctx context.Context, msg *redis.Message) error {
	productID := msg.Data["product_id"].(string)
	currentStock := int(msg.Data["current_stock"].(float64))
	threshold := int(msg.Data["threshold"].(float64))

	s.LogInfo("processing inventory alert", "product_id", productID, "current_stock", currentStock, "threshold", threshold)

	// 发送库存告警
	return s.sendInventoryAlert(ctx, productID, currentStock, threshold)
}

// handleDefault 默认处理器
func (s *WorkerService) handleDefault(ctx context.Context, msg *redis.Message) error {
	s.LogWarn("unhandled message", "topic", msg.Topic, "message_id", msg.ID)
	return nil
}

// StopConsumers 停止所有消费者
func (s *WorkerService) StopConsumers() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for topic, cancel := range s.consumers {
		cancel()
		s.LogInfo("consumer stopped", "topic", topic)
	}

	s.consumers = make(map[string]context.CancelFunc)
}

// 业务方法实现

func (s *WorkerService) sendOrderReminderNotification(ctx context.Context, orderID, userID, reminderType string) error {
	// 模拟发送订单提醒通知
	s.LogInfo("sending order reminder notification", "order_id", orderID, "user_id", userID, "type", reminderType)

	// 这里可以集成邮件、短信、推送等服务
	time.Sleep(100 * time.Millisecond) // 模拟处理时间

	return nil
}

func (s *WorkerService) sendUserNotification(ctx context.Context, userID, notificationType, content string) error {
	// 模拟发送用户通知
	s.LogInfo("sending user notification", "user_id", userID, "type", notificationType, "content", content)

	// 这里可以集成推送、站内信等服务
	time.Sleep(50 * time.Millisecond) // 模拟处理时间

	return nil
}

func (s *WorkerService) performSystemCleanup(ctx context.Context, cleanupType string, params map[string]interface{}) error {
	// 模拟系统清理
	s.LogInfo("performing system cleanup", "type", cleanupType, "params", params)

	// 这里可以执行日志清理、缓存清理、临时文件清理等
	time.Sleep(200 * time.Millisecond) // 模拟处理时间

	return nil
}

func (s *WorkerService) retryPayment(ctx context.Context, paymentID, orderID string, retryCount int) error {
	// 模拟支付重试
	s.LogInfo("retrying payment", "payment_id", paymentID, "order_id", orderID, "retry_count", retryCount)

	// 这里可以调用支付接口重试
	time.Sleep(300 * time.Millisecond) // 模拟处理时间

	return nil
}

func (s *WorkerService) sendInventoryAlert(ctx context.Context, productID string, currentStock, threshold int) error {
	// 模拟发送库存告警
	s.LogInfo("sending inventory alert", "product_id", productID, "current_stock", currentStock, "threshold", threshold)

	// 这里可以发送告警邮件、通知管理员等
	time.Sleep(100 * time.Millisecond) // 模拟处理时间

	return nil
}
