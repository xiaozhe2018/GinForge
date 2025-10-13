package errors

// 业务错误码定义
const (
	// 通用错误码 (1000-1999)
	Success            = 0
	UnknownError       = 1000
	InvalidParameter   = 1001
	MissingParameter   = 1002
	InvalidFormat      = 1003
	OperationFailed    = 1004
	ResourceNotFound   = 1005
	ResourceExists     = 1006
	PermissionDenied   = 1007
	RateLimitExceeded  = 1008
	ServiceUnavailable = 1009

	// 认证授权错误码 (2000-2999)
	Unauthorized    = 2000
	TokenExpired    = 2001
	TokenInvalid    = 2002
	TokenMissing    = 2003
	PasswordError   = 2004
	AccountDisabled = 2005
	AccountLocked   = 2006
	PermissionError = 2007
	RoleError       = 2008

	// 用户相关错误码 (3000-3999)
	UserNotFound      = 3000
	UserExists        = 3001
	UserDisabled      = 3002
	UserLocked        = 3003
	EmailExists       = 3004
	PhoneExists       = 3005
	UsernameExists    = 3006
	PasswordWeak      = 3007
	ProfileIncomplete = 3008

	// 商户相关错误码 (4000-4999)
	MerchantNotFound            = 4000
	MerchantExists              = 4001
	MerchantDisabled            = 4002
	MerchantPending             = 4003
	ShopNameExists              = 4004
	BusinessLicenseInvalid      = 4005
	MerchantBalanceInsufficient = 4006

	// 商品相关错误码 (5000-5999)
	ProductNotFound         = 5000
	ProductExists           = 5001
	ProductDisabled         = 5002
	ProductOutOfStock       = 5003
	ProductCategoryNotFound = 5004
	ProductPriceInvalid     = 5005
	ProductImageInvalid     = 5006

	// 订单相关错误码 (6000-6999)
	OrderNotFound      = 6000
	OrderExists        = 6001
	OrderCancelled     = 6002
	OrderPaid          = 6003
	OrderShipped       = 6004
	OrderCompleted     = 6005
	OrderRefunded      = 6006
	OrderExpired       = 6007
	OrderStatusInvalid = 6008
	PaymentFailed      = 6009
	RefundFailed       = 6010

	// 系统相关错误码 (9000-9999)
	DatabaseError     = 9000
	CacheError        = 9001
	NetworkError      = 9002
	FileUploadError   = 9003
	FileDownloadError = 9004
	ConfigError       = 9005
	ServiceError      = 9006
	GatewayError      = 9007
	TimeoutError      = 9008
	InternalError     = 9999
)

// 错误码映射
var ErrorMessages = map[int]string{
	Success:            "操作成功",
	UnknownError:       "未知错误",
	InvalidParameter:   "参数无效",
	MissingParameter:   "缺少必要参数",
	InvalidFormat:      "格式错误",
	OperationFailed:    "操作失败",
	ResourceNotFound:   "资源不存在",
	ResourceExists:     "资源已存在",
	PermissionDenied:   "权限不足",
	RateLimitExceeded:  "请求过于频繁",
	ServiceUnavailable: "服务不可用",

	Unauthorized:    "未授权",
	TokenExpired:    "令牌已过期",
	TokenInvalid:    "令牌无效",
	TokenMissing:    "缺少令牌",
	PasswordError:   "密码错误",
	AccountDisabled: "账户已禁用",
	AccountLocked:   "账户已锁定",
	PermissionError: "权限错误",
	RoleError:       "角色错误",

	UserNotFound:      "用户不存在",
	UserExists:        "用户已存在",
	UserDisabled:      "用户已禁用",
	UserLocked:        "用户已锁定",
	EmailExists:       "邮箱已存在",
	PhoneExists:       "手机号已存在",
	UsernameExists:    "用户名已存在",
	PasswordWeak:      "密码强度不够",
	ProfileIncomplete: "资料不完整",

	MerchantNotFound:            "商户不存在",
	MerchantExists:              "商户已存在",
	MerchantDisabled:            "商户已禁用",
	MerchantPending:             "商户待审核",
	ShopNameExists:              "店铺名已存在",
	BusinessLicenseInvalid:      "营业执照无效",
	MerchantBalanceInsufficient: "商户余额不足",

	ProductNotFound:         "商品不存在",
	ProductExists:           "商品已存在",
	ProductDisabled:         "商品已下架",
	ProductOutOfStock:       "商品库存不足",
	ProductCategoryNotFound: "商品分类不存在",
	ProductPriceInvalid:     "商品价格无效",
	ProductImageInvalid:     "商品图片无效",

	OrderNotFound:      "订单不存在",
	OrderExists:        "订单已存在",
	OrderCancelled:     "订单已取消",
	OrderPaid:          "订单已支付",
	OrderShipped:       "订单已发货",
	OrderCompleted:     "订单已完成",
	OrderRefunded:      "订单已退款",
	OrderExpired:       "订单已过期",
	OrderStatusInvalid: "订单状态无效",
	PaymentFailed:      "支付失败",
	RefundFailed:       "退款失败",

	DatabaseError:     "数据库错误",
	CacheError:        "缓存错误",
	NetworkError:      "网络错误",
	FileUploadError:   "文件上传失败",
	FileDownloadError: "文件下载失败",
	ConfigError:       "配置错误",
	ServiceError:      "服务错误",
	GatewayError:      "网关错误",
	TimeoutError:      "请求超时",
	InternalError:     "内部错误",
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}
