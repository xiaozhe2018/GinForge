package handler

import (
	"goweb/pkg/base"
	"goweb/pkg/gateway"
	"goweb/pkg/logger"
	"goweb/services/demo/internal/service"

	"github.com/gin-gonic/gin"
)

type DemoHandler struct {
	*base.BaseHandler
	svc *service.DemoService
}

func NewDemoHandler(s *service.DemoService, gatewayClient *gateway.Client, log logger.Logger) *DemoHandler {
	return &DemoHandler{
		BaseHandler: base.NewBaseHandler(log),
		svc:         s,
	}
}

// GetData 示例 API
func (h *DemoHandler) GetData(c *gin.Context) {
	data, err := h.svc.GetData()
	if err != nil {
		h.LogError("get data error", err)
		h.InternalError(c, "获取数据失败")
		return
	}
	h.Success(c, data)
}

// GetUserInfo 通过 Gateway 获取用户信息
func (h *DemoHandler) GetUserInfo(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		h.BadRequest(c, "用户ID不能为空")
		return
	}

	data, err := h.svc.GetUserInfo(userID)
	if err != nil {
		h.LogError("get user info error", err, "user_id", userID)
		h.InternalError(c, "获取用户信息失败")
		return
	}

	h.Success(c, data)
}
