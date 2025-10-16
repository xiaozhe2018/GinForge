package model

import (
	"time"
)

// Articles Articles管理
type Articles struct {
	Id int64 `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;not null" json:"id"` // 文章ID
	Title string `gorm:"column:title;type:varchar(255);not null" json:"title"` // 标题
	Content string `gorm:"column:content;type:text;not null" json:"content"` // 内容
	Summary string `gorm:"column:summary;type:varchar(500)" json:"summary,omitempty"` // 摘要
	AuthorId int64 `gorm:"column:author_id;type:bigint unsigned;not null" json:"author_id"` // 作者ID
	AuthorName string `gorm:"column:author_name;type:varchar(100)" json:"author_name,omitempty"` // 作者名称
	CategoryId *int64 `gorm:"column:category_id;type:bigint unsigned" json:"category_id,omitempty"` // 分类ID
	CoverImage string `gorm:"column:cover_image;type:varchar(500)" json:"cover_image,omitempty"` // 封面图片
	Tags string `gorm:"column:tags;type:varchar(255)" json:"tags,omitempty"` // 标签（逗号分隔）
	Status int8 `gorm:"column:status;type:tinyint;not null" json:"status"` // 状态:0草稿,1已发布,2已下线
	IsTop int8 `gorm:"column:is_top;type:tinyint(1);not null" json:"is_top"` // 是否置顶
	ViewCount int `gorm:"column:view_count;type:int unsigned;not null" json:"view_count"` // 浏览次数
	LikeCount int `gorm:"column:like_count;type:int unsigned;not null" json:"like_count"` // 点赞次数
	CommentCount int `gorm:"column:comment_count;type:int unsigned;not null" json:"comment_count"` // 评论次数
	PublishedAt *time.Time `gorm:"column:published_at;type:datetime" json:"published_at,omitempty"` // 发布时间
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime" json:"created_at,omitempty"` // 创建时间
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at,omitempty"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"` // 删除时间
}

// TableName 指定表名
func (articles *Articles) TableName() string {
	return "articles"
}

// ArticlesListRequest Articles管理列表请求
type ArticlesListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`
	SortBy   string `form:"sort_by"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ArticlesCreateRequest 创建Articles管理请求
type ArticlesCreateRequest struct {
	Title string `json:"title" binding:"required,max:255"` // 标题
	Content string `json:"content" binding:"required"` // 内容
	Summary string `json:"summary" binding:"max:500"` // 摘要
	AuthorId int64 `json:"author_id" binding:"required"` // 作者ID
	AuthorName string `json:"author_name" binding:"max:100"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	CoverImage string `json:"cover_image" binding:"max:500"` // 封面图片
	Tags string `json:"tags" binding:"max:255"` // 标签（逗号分隔）
	Status int8 `json:"status" binding:"required"` // 状态:0草稿,1已发布,2已下线
	IsTop int8 `json:"is_top" binding:"required"` // 是否置顶
	ViewCount int `json:"view_count" binding:"required"` // 浏览次数
	LikeCount int `json:"like_count" binding:"required"` // 点赞次数
	CommentCount int `json:"comment_count" binding:"required"` // 评论次数
	PublishedAt *time.Time `json:"published_at"` // 发布时间
}

// ArticlesUpdateRequest 更新Articles管理请求
type ArticlesUpdateRequest struct {
	Title string `json:"title" binding:"required,max:255"` // 标题
	Content string `json:"content" binding:"required"` // 内容
	Summary string `json:"summary" binding:"max:500"` // 摘要
	AuthorId int64 `json:"author_id" binding:"required"` // 作者ID
	AuthorName string `json:"author_name" binding:"max:100"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	CoverImage string `json:"cover_image" binding:"max:500"` // 封面图片
	Tags string `json:"tags" binding:"max:255"` // 标签（逗号分隔）
	Status int8 `json:"status" binding:"required"` // 状态:0草稿,1已发布,2已下线
	IsTop int8 `json:"is_top" binding:"required"` // 是否置顶
	ViewCount int `json:"view_count" binding:"required"` // 浏览次数
	LikeCount int `json:"like_count" binding:"required"` // 点赞次数
	CommentCount int `json:"comment_count" binding:"required"` // 评论次数
	PublishedAt *time.Time `json:"published_at"` // 发布时间
}

// ArticlesResponse Articles管理响应
type ArticlesResponse struct {
	Id int64 `json:"id"` // 文章ID
	Title string `json:"title"` // 标题
	Summary string `json:"summary"` // 摘要
	AuthorId int64 `json:"author_id"` // 作者ID
	AuthorName string `json:"author_name"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	CoverImage string `json:"cover_image"` // 封面图片
	Tags string `json:"tags"` // 标签（逗号分隔）
	Status int8 `json:"status"` // 状态:0草稿,1已发布,2已下线
	IsTop int8 `json:"is_top"` // 是否置顶
	ViewCount int `json:"view_count"` // 浏览次数
	LikeCount int `json:"like_count"` // 点赞次数
	CommentCount int `json:"comment_count"` // 评论次数
	PublishedAt *time.Time `json:"published_at"` // 发布时间
	CreatedAt *time.Time `json:"created_at"` // 创建时间
	UpdatedAt *time.Time `json:"updated_at"` // 更新时间
}

// ToArticlesResponse 转换为响应对象
func ToArticlesResponse(articles *Articles) *ArticlesResponse {
	if articles == nil {
		return nil
	}
	
	return &ArticlesResponse{
		Id: articles.Id,
		Title: articles.Title,
		Summary: articles.Summary,
		AuthorId: articles.AuthorId,
		AuthorName: articles.AuthorName,
		CategoryId: articles.CategoryId,
		CoverImage: articles.CoverImage,
		Tags: articles.Tags,
		Status: articles.Status,
		IsTop: articles.IsTop,
		ViewCount: articles.ViewCount,
		LikeCount: articles.LikeCount,
		CommentCount: articles.CommentCount,
		PublishedAt: articles.PublishedAt,
		CreatedAt: articles.CreatedAt,
		UpdatedAt: articles.UpdatedAt,
	}
}

// ToArticlesResponseList 批量转换为响应对象
func ToArticlesResponseList(list []*Articles) []*ArticlesResponse {
	result := make([]*ArticlesResponse, 0, len(list))
	for _, item := range list {
		result = append(result, ToArticlesResponse(item))
	}
	return result
}
