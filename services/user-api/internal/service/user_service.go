package service

import (
	"goweb/pkg/model"
)

// UserService 用户服务
type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUserProfile 获取用户资料
func (s *UserService) GetUserProfile(userID string) (*model.User, error) {
	// 模拟数据，实际应该从数据库获取
	return &model.User{
		ID:       userID,
		Username: "user123",
		Email:    "user@example.com",
		Nickname: "用户昵称",
		Avatar:   "https://example.com/avatar.jpg",
		Phone:    "13800138000",
		Status:   1,
	}, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(userID string, req *UpdateProfileRequest) error {
	// 实际应该更新数据库
	return nil
}

// GetUserOrders 获取用户订单
func (s *UserService) GetUserOrders(userID string, page, size int) ([]*model.Order, int, error) {
	// 模拟数据
	orders := []*model.Order{
		{
			ID:         "order1",
			UserID:     userID,
			MerchantID: "merchant1",
			Total:      99.99,
			Status:     1,
		},
	}
	return orders, 1, nil
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"新昵称"`
	Avatar   string `json:"avatar" example:"https://example.com/new-avatar.jpg"`
	Phone    string `json:"phone" example:"13800138001"`
}
