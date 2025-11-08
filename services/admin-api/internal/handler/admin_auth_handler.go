package handler

import (
	"fmt"
	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminAuthHandler 管理后台认证处理器
type AdminAuthHandler struct {
	userService *service.UserService
	logger      logger.Logger
}

// NewAdminAuthHandler 创建管理后台认证处理器实例
func NewAdminAuthHandler(userService *service.UserService) *AdminAuthHandler {
	return &AdminAuthHandler{
		userService: userService,
	}
}

// SetLogger 设置日志器
func (h *AdminAuthHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// Login 用户登录
// @Summary 管理员登录
// @Description 管理员用户登录接口
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body model.AdminUserLoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=model.AdminUserLoginResponse} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/auth/login [post]
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req model.AdminUserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind login request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 获取客户端IP和User-Agent
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 将User-Agent传递给context（如果需要的话，可以通过context传递）
	// 这里我们直接传递IP，User-Agent可以在服务层从请求中获取

	// 调用服务层登录
	result, err := h.userService.Login(&req, clientIP, userAgent)
	if err != nil {
		h.logger.Error("login error", err, "ip", clientIP, "username", req.Username)
		response.Unauthorized(c, err.Error())
		return
	}

	h.logger.Info("user login", "user_id", result.User.ID, "username", result.User.Username, "ip", clientIP)
	response.Success(c, result)
}

// Logout 用户登出
// @Summary 管理员登出
// @Description 管理员用户登出接口
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "登出成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /admin/auth/logout [post]
func (h *AdminAuthHandler) Logout(c *gin.Context) {
	// 从上下文获取用户ID和用户名
	userID, exists := c.Get("user_id")
	if !exists {
		// 即使没有用户ID，也允许登出（清除前端token）
		response.Success(c, gin.H{"message": "登出成功"})
		return
	}

	username, _ := c.Get("username")
	usernameStr := ""
	if username != nil {
		usernameStr = username.(string)
	}

	// 从Authorization header获取token
	authHeader := c.GetHeader("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 调用服务层处理登出逻辑（将token加入黑名单）
	if err := h.userService.Logout(userID.(string), usernameStr, token, clientIP); err != nil {
		h.logger.Error("logout error", err, "user_id", userID, "username", usernameStr, "ip", clientIP)
		// 即使服务层失败，也返回成功（登出操作应该是容错的）
	}

	h.logger.Info("user logout", "user_id", userID, "username", usernameStr, "ip", clientIP)
	response.Success(c, gin.H{"message": "登出成功"})
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的基本信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.AdminUser} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/auth/profile [get]
func (h *AdminAuthHandler) GetProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 转换user_id为uint64
	var uid uint64
	if userIDStr, ok := userID.(string); ok {
		fmt.Sscanf(userIDStr, "%d", &uid)
	} else {
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUser(uid)
	if err != nil {
		h.logger.Error("get user profile error", err)
		response.InternalError(c, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前用户信息
// @Description 更新当前登录用户的基本信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AdminUserUpdateRequest true "更新请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/auth/profile [put]
func (h *AdminAuthHandler) UpdateProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 转换user_id为uint64
	var uid uint64
	if userIDStr, ok := userID.(string); ok {
		fmt.Sscanf(userIDStr, "%d", &uid)
	} else {
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	var req model.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind update profile request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 更新用户信息
	if err := h.userService.UpdateUser(uid, &req); err != nil {
		h.logger.Error("update user profile error", err)
		response.InternalError(c, "更新用户信息失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/auth/change-password [put]
func (h *AdminAuthHandler) ChangePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 转换user_id为uint64
	var uid uint64
	if userIDStr, ok := userID.(string); ok {
		fmt.Sscanf(userIDStr, "%d", &uid)
	} else {
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind change password request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 修改密码
	if err := h.userService.ChangePassword(uid, &req); err != nil {
		h.logger.Error("change password error", err)
		response.InternalError(c, "修改密码失败")
		return
	}

	response.Success(c, gin.H{"message": "修改成功"})
}
