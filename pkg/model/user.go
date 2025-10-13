package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Status    int            `json:"status" gorm:"type:tinyint;default:1;comment:1正常2禁用"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Merchant 商户模型
type Merchant struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	UserID      string         `json:"user_id" gorm:"type:char(36);not null;index"`
	ShopName    string         `json:"shop_name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Logo        string         `json:"logo" gorm:"type:varchar(255)"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1;comment:1正常2禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product 商品模型
type Product struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	MerchantID  string         `json:"merchant_id" gorm:"type:char(36);not null;index"`
	Name        string         `json:"name" gorm:"type:varchar(200);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int            `json:"stock" gorm:"type:int;default:0"`
	Images      string         `json:"images" gorm:"type:text"`
	CategoryID  string         `json:"category_id" gorm:"type:char(36);index"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1;comment:1上架2下架"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Order 订单模型
type Order struct {
	ID         string         `json:"id" gorm:"type:char(36);primaryKey"`
	UserID     string         `json:"user_id" gorm:"type:char(36);not null;index"`
	MerchantID string         `json:"merchant_id" gorm:"type:char(36);not null;index"`
	Total      float64        `json:"total" gorm:"type:decimal(10,2);not null"`
	Status     int            `json:"status" gorm:"type:tinyint;default:1;comment:1待付款2已付款3已发货4已完成5已取消"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        string  `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID   string  `json:"order_id" gorm:"type:char(36);not null;index"`
	ProductID string  `json:"product_id" gorm:"type:char(36);not null"`
	Quantity  int     `json:"quantity" gorm:"type:int;not null"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}
