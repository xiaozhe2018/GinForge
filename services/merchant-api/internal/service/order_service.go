package service

import (
	"goweb/pkg/model"
)

// OrderService 订单服务
type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

// GetMerchantOrders 获取商户订单
func (s *OrderService) GetMerchantOrders(merchantID string, page, size int) ([]*model.Order, int, error) {
	// 模拟数据
	orders := []*model.Order{
		{
			ID:         "order1",
			UserID:     "user1",
			MerchantID: merchantID,
			Total:      99.99,
			Status:     1,
		},
	}
	return orders, 1, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID string, status int) error {
	// 实际应该更新数据库
	return nil
}
