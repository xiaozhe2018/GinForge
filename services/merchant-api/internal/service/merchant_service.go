package service

import (
	"goweb/pkg/model"
)

// MerchantService 商户服务
type MerchantService struct{}

func NewMerchantService() *MerchantService {
	return &MerchantService{}
}

// GetMerchantInfo 获取商户信息
func (s *MerchantService) GetMerchantInfo(merchantID string) (*model.Merchant, error) {
	// 模拟数据
	return &model.Merchant{
		ID:          merchantID,
		UserID:      "user1",
		ShopName:    "示例商店",
		Description: "这是一个示例商店",
		Logo:        "https://example.com/logo.jpg",
		Status:      1,
	}, nil
}

// UpdateMerchantInfo 更新商户信息
func (s *MerchantService) UpdateMerchantInfo(merchantID string, req *UpdateMerchantRequest) error {
	// 实际应该更新数据库
	return nil
}

// UpdateMerchantRequest 更新商户请求
type UpdateMerchantRequest struct {
	ShopName    string `json:"shop_name" example:"新商店名称"`
	Description string `json:"description" example:"商店描述"`
	Logo        string `json:"logo" example:"https://example.com/new-logo.jpg"`
}
