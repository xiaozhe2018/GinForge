package constants

// 用户状态
const (
	UserStatusActive   = 1 // 正常
	UserStatusDisabled = 2 // 禁用
	UserStatusLocked   = 3 // 锁定
)

// 商户状态
const (
	MerchantStatusPending  = 1 // 待审核
	MerchantStatusApproved = 2 // 已审核
	MerchantStatusRejected = 3 // 已拒绝
	MerchantStatusDisabled = 4 // 已禁用
)

// 商品状态
const (
	ProductStatusDraft    = 1 // 草稿
	ProductStatusActive   = 2 // 上架
	ProductStatusInactive = 3 // 下架
	ProductStatusDeleted  = 4 // 删除
)

// 订单状态
const (
	OrderStatusPending   = 1 // 待支付
	OrderStatusPaid      = 2 // 已支付
	OrderStatusShipped   = 3 // 已发货
	OrderStatusCompleted = 4 // 已完成
	OrderStatusCancelled = 5 // 已取消
	OrderStatusRefunded  = 6 // 已退款
)

// 支付状态
const (
	PaymentStatusPending  = 1 // 待支付
	PaymentStatusPaid     = 2 // 已支付
	PaymentStatusFailed   = 3 // 支付失败
	PaymentStatusRefunded = 4 // 已退款
)

// 性别
const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)

// 文件类型
const (
	FileTypeImage = "image"
	FileTypeVideo = "video"
	FileTypeAudio = "audio"
	FileTypeDoc   = "document"
	FileTypeOther = "other"
)

// 缓存键前缀
const (
	CacheKeyUser     = "user:"
	CacheKeyMerchant = "merchant:"
	CacheKeyProduct  = "product:"
	CacheKeyOrder    = "order:"
	CacheKeySession  = "session:"
	CacheKeyToken    = "token:"
)

// 缓存过期时间（秒）
const (
	CacheExpireShort  = 300   // 5分钟
	CacheExpireMedium = 1800  // 30分钟
	CacheExpireLong   = 3600  // 1小时
	CacheExpireDay    = 86400 // 1天
)

// 分页默认值
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// 密码强度
const (
	PasswordMinLength = 8
	PasswordMaxLength = 32
)

// 用户名规则
const (
	UsernameMinLength = 3
	UsernameMaxLength = 20
)

// 手机号格式
const (
	PhonePattern = `^1[3-9]\d{9}$`
)

// 邮箱格式
const (
	EmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// 文件大小限制（字节）
const (
	MaxImageSize = 5 * 1024 * 1024   // 5MB
	MaxVideoSize = 100 * 1024 * 1024 // 100MB
	MaxDocSize   = 10 * 1024 * 1024  // 10MB
)

// 允许的图片格式
var AllowedImageFormats = []string{
	"jpg", "jpeg", "png", "gif", "webp",
}

// 允许的视频格式
var AllowedVideoFormats = []string{
	"mp4", "avi", "mov", "wmv", "flv",
}

// 允许的文档格式
var AllowedDocFormats = []string{
	"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx",
}
