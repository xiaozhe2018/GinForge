package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// SoftDeleteModel 软删除模型
type SoftDeleteModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TimestampModel 时间戳模型
type TimestampModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IDModel ID模型
type IDModel struct {
	ID uint `json:"id" gorm:"primarykey"`
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// PaginationResult 分页结果
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// NewPagination 创建分页参数
func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset 计算偏移量
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// NewPaginationResult 创建分页结果
func NewPaginationResult(data interface{}, total int64, page, pageSize int) *PaginationResult {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &PaginationResult{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// Sort 排序参数
type Sort struct {
	Field string `json:"field" form:"sort_field"`
	Order string `json:"order" form:"sort_order" binding:"oneof=asc desc"`
}

// NewSort 创建排序参数
func NewSort(field, order string) *Sort {
	if order != "desc" {
		order = "asc"
	}
	return &Sort{
		Field: field,
		Order: order,
	}
}

// Search 搜索参数
type Search struct {
	Keyword string   `json:"keyword" form:"keyword"`
	Fields  []string `json:"fields" form:"fields"`
}

// NewSearch 创建搜索参数
func NewSearch(keyword string, fields ...string) *Search {
	return &Search{
		Keyword: keyword,
		Fields:  fields,
	}
}

// DateRange 日期范围
type DateRange struct {
	StartDate *time.Time `json:"start_date" form:"start_date"`
	EndDate   *time.Time `json:"end_date" form:"end_date"`
}

// NewDateRange 创建日期范围
func NewDateRange(start, end *time.Time) *DateRange {
	return &DateRange{
		StartDate: start,
		EndDate:   end,
	}
}

// IsValid 检查日期范围是否有效
func (d *DateRange) IsValid() bool {
	if d.StartDate == nil || d.EndDate == nil {
		return true
	}
	return d.StartDate.Before(*d.EndDate) || d.StartDate.Equal(*d.EndDate)
}

// Filter 过滤参数
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, ne, gt, gte, lt, lte, like, in, not_in
	Value    interface{} `json:"value"`
}

// NewFilter 创建过滤参数
func NewFilter(field, operator string, value interface{}) *Filter {
	return &Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// QueryParams 查询参数
type QueryParams struct {
	Pagination *Pagination `json:"pagination"`
	Sort       *Sort       `json:"sort"`
	Search     *Search     `json:"search"`
	DateRange  *DateRange  `json:"date_range"`
	Filters    []*Filter   `json:"filters"`
}

// NewQueryParams 创建查询参数
func NewQueryParams() *QueryParams {
	return &QueryParams{
		Pagination: NewPagination(1, 10),
		Sort:       NewSort("id", "desc"),
		Search:     NewSearch(""),
		DateRange:  NewDateRange(nil, nil),
		Filters:    make([]*Filter, 0),
	}
}

// AddFilter 添加过滤条件
func (q *QueryParams) AddFilter(field, operator string, value interface{}) {
	q.Filters = append(q.Filters, NewFilter(field, operator, value))
}

// SetPagination 设置分页
func (q *QueryParams) SetPagination(page, pageSize int) {
	q.Pagination = NewPagination(page, pageSize)
}

// SetSort 设置排序
func (q *QueryParams) SetSort(field, order string) {
	q.Sort = NewSort(field, order)
}

// SetSearch 设置搜索
func (q *QueryParams) SetSearch(keyword string, fields ...string) {
	q.Search = NewSearch(keyword, fields...)
}

// SetDateRange 设置日期范围
func (q *QueryParams) SetDateRange(start, end *time.Time) {
	q.DateRange = NewDateRange(start, end)
}
