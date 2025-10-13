package handler

import (
	"github.com/gin-gonic/gin"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/merchant-api/internal/service"
)

// MerchantHandler 商户处理器
type MerchantHandler struct {
	merchantService *service.MerchantService
	productService  *service.ProductService
	orderService    *service.OrderService
	logger          logger.Logger
}

func NewMerchantHandler(merchantService *service.MerchantService, productService *service.ProductService, orderService *service.OrderService) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
		productService:  productService,
		orderService:    orderService,
	}
}

// SetLogger 设置日志器
func (h *MerchantHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// GetMerchantInfo 获取商户信息
// @Summary      获取商户信息
// @Description  获取当前登录商户的详细信息
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=model.Merchant}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /merchant/info [get]
func (h *MerchantHandler) GetMerchantInfo(c *gin.Context) {
	merchantID := c.GetString("merchant_id") // 从JWT中获取
	if merchantID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	merchant, err := h.merchantService.GetMerchantInfo(merchantID)
	if err != nil {
		h.logger.Error("get merchant info error", err)
		response.InternalError(c, "获取商户信息失败")
		return
	}

	response.Success(c, merchant)
}

// UpdateMerchantInfo 更新商户信息
// @Summary      更新商户信息
// @Description  更新当前登录商户的信息
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  service.UpdateMerchantRequest  true  "更新商户信息请求"
// @Success      200     {object}  response.Response{data=object}
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /merchant/info [put]
func (h *MerchantHandler) UpdateMerchantInfo(c *gin.Context) {
	merchantID := c.GetString("merchant_id")
	if merchantID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.merchantService.UpdateMerchantInfo(merchantID, &req); err != nil {
		h.logger.Error("update merchant info error", err)
		response.InternalError(c, "更新商户信息失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetProducts 获取商品列表
// @Summary      获取商品列表
// @Description  获取当前商户的商品列表
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=object{list=[]model.Product,total=int}}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /product/list [get]
func (h *MerchantHandler) GetProducts(c *gin.Context) {
	merchantID := c.GetString("merchant_id")
	if merchantID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	products, total, err := h.productService.GetProducts(merchantID, 1, 10)
	if err != nil {
		h.logger.Error("get products error", err)
		response.InternalError(c, "获取商品列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  products,
		"total": total,
	})
}

// CreateProduct 创建商品
// @Summary      创建商品
// @Description  创建新商品
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  service.CreateProductRequest  true  "创建商品请求"
// @Success      200     {object}  response.Response{data=model.Product}
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /product/create [post]
func (h *MerchantHandler) CreateProduct(c *gin.Context) {
	merchantID := c.GetString("merchant_id")
	if merchantID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	product, err := h.productService.CreateProduct(merchantID, &req)
	if err != nil {
		h.logger.Error("create product error", err)
		response.InternalError(c, "创建商品失败")
		return
	}

	response.Success(c, product)
}

// GetOrders 获取订单列表
// @Summary      获取订单列表
// @Description  获取当前商户的订单列表
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=object{list=[]model.Order,total=int}}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /order/list [get]
func (h *MerchantHandler) GetOrders(c *gin.Context) {
	merchantID := c.GetString("merchant_id")
	if merchantID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	orders, total, err := h.orderService.GetMerchantOrders(merchantID, 1, 10)
	if err != nil {
		h.logger.Error("get orders error", err)
		response.InternalError(c, "获取订单列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
	})
}
