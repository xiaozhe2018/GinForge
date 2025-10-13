package service

import (
	"goweb/pkg/base"
	"goweb/pkg/gateway"
	"goweb/pkg/logger"
)

type DemoService struct {
	*base.BaseService
	gatewayClient *gateway.Client
}

func NewDemoService(gatewayClient *gateway.Client, log logger.Logger) *DemoService {
	return &DemoService{
		BaseService:   base.NewBaseService(log),
		gatewayClient: gatewayClient,
	}
}

func (s *DemoService) GetData() (any, error) {
	s.LogInfo("getting demo data")
	return map[string]any{"hello": "world"}, nil
}

// 通过 Gateway 调用其他服务
func (s *DemoService) GetUserInfo(userID string) (any, error) {
	s.LogInfo("getting user info via gateway", "user_id", userID)

	// 这里演示如何通过 Gateway 调用其他服务
	// 实际使用时需要确保 Gateway 服务正在运行
	resp, err := s.gatewayClient.GetUser(nil, userID)
	if err != nil {
		s.LogError("failed to get user via gateway", err, "user_id", userID)
		return nil, err
	}

	return resp.Data, nil
}
