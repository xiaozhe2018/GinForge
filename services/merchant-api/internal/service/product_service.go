package service

import (
	"goweb/pkg/model"
)

// ProductService 商品服务
type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

// GetProducts 获取商品列表
func (s *ProductService) GetProducts(merchantID string, page, size int) ([]*model.Product, int, error) {
	// 模拟数据
	products := []*model.Product{
		{
			ID:          "product1",
			MerchantID:  merchantID,
			Name:        "示例商品1",
			Description: "这是一个示例商品",
			Price:       99.99,
			Stock:       100,
			Images:      "https://example.com/product1.jpg",
			Status:      1,
		},
	}
	return products, 1, nil
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(merchantID string, req *CreateProductRequest) (*model.Product, error) {
	// 实际应该保存到数据库
	product := &model.Product{
		ID:          "new_product",
		MerchantID:  merchantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Images:      req.Images,
		Status:      1,
	}
	return product, nil
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"商品名称"`
	Description string  `json:"description" example:"商品描述"`
	Price       float64 `json:"price" binding:"required,min=0" example:"99.99"`
	Stock       int     `json:"stock" binding:"min=0" example:"100"`
	Images      string  `json:"images" example:"https://example.com/product.jpg"`
}
