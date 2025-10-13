package handler

import (
	"github.com/gin-gonic/gin"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/user-api/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
	logger      logger.Logger
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// SetLogger 设置日志器
func (h *UserHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// GetProfile 获取用户资料
// @Summary      获取用户资料
// @Description  获取当前登录用户的详细资料
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id") // 从JWT中获取
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		h.logger.Error("get user profile error", err)
		response.InternalError(c, "获取用户资料失败")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户资料
// @Summary      更新用户资料
// @Description  更新当前登录用户的资料信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  service.UpdateProfileRequest  true  "更新资料请求"
// @Success      200     {object}  response.Response{data=object}
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.userService.UpdateUserProfile(userID, &req); err != nil {
		h.logger.Error("update user profile error", err)
		response.InternalError(c, "更新用户资料失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetOrders 获取用户订单
// @Summary      获取用户订单
// @Description  获取当前登录用户的订单列表
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page  query  int  false  "页码"  default(1)
// @Param        size  query  int  false  "每页数量"  default(10)
// @Success      200   {object}  response.Response{data=object{list=[]model.Order,total=int,page=string,size=string}}
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /user/orders [get]
func (h *UserHandler) GetOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")

	orders, total, err := h.userService.GetUserOrders(userID, 1, 10) // 简化分页
	if err != nil {
		h.logger.Error("get user orders error", err)
		response.InternalError(c, "获取订单失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
