package model

import (
	"time"
)

// Articles Articles管理
type Articles struct {
	Id int64 `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;not null" json:"id"` // 文章ID
	Title string `gorm:"column:title;type:varchar(200);not null" json:"title"` // 文章标题
	Slug string `gorm:"column:slug;type:varchar(200)" json:"slug,omitempty"` // URL别名
	AuthorId int64 `gorm:"column:author_id;type:bigint unsigned;not null" json:"author_id"` // 作者ID
	AuthorName string `gorm:"column:author_name;type:varchar(50)" json:"author_name,omitempty"` // 作者名称
	CategoryId *int64 `gorm:"column:category_id;type:bigint unsigned" json:"category_id,omitempty"` // 分类ID
	Summary string `gorm:"column:summary;type:varchar(500)" json:"summary,omitempty"` // 文章摘要
	Content string `gorm:"column:content;type:longtext;not null" json:"content"` // 文章内容
	CoverImage string `gorm:"column:cover_image;type:varchar(255)" json:"cover_image,omitempty"` // 封面图片
	ViewCount int `gorm:"column:view_count;type:int;not null" json:"view_count"` // 浏览次数
	LikeCount int `gorm:"column:like_count;type:int;not null" json:"like_count"` // 点赞次数
	CommentCount int `gorm:"column:comment_count;type:int;not null" json:"comment_count"` // 评论次数
	IsPublished int8 `gorm:"column:is_published;type:tinyint(1);not null" json:"is_published"` // 是否发布: 1-已发布, 0-草稿
	IsTop int8 `gorm:"column:is_top;type:tinyint(1);not null" json:"is_top"` // 是否置顶: 1-是, 0-否
	IsFeatured int8 `gorm:"column:is_featured;type:tinyint(1);not null" json:"is_featured"` // 是否推荐: 1-是, 0-否
	PublishedAt *time.Time `gorm:"column:published_at;type:timestamp" json:"published_at,omitempty"` // 发布时间
	Tags string `gorm:"column:tags;type:varchar(500)" json:"tags,omitempty"` // 标签(逗号分隔)
	SeoTitle string `gorm:"column:seo_title;type:varchar(200)" json:"seo_title,omitempty"` // SEO标题
	SeoKeywords string `gorm:"column:seo_keywords;type:varchar(500)" json:"seo_keywords,omitempty"` // SEO关键词
	SeoDescription string `gorm:"column:seo_description;type:varchar(500)" json:"seo_description,omitempty"` // SEO描述
	Status int8 `gorm:"column:status;type:tinyint(1);not null" json:"status"` // 状态: 1-正常, 0-禁用
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deleted_at,omitempty"` // 删除时间
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
	Title string `json:"title" binding:"required,max:200"` // 文章标题
	Slug string `json:"slug" binding:"max:200"` // URL别名
	AuthorId int64 `json:"author_id" binding:"required"` // 作者ID
	AuthorName string `json:"author_name" binding:"max:50"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	Summary string `json:"summary" binding:"max:500"` // 文章摘要
	Content string `json:"content" binding:"required"` // 文章内容
	CoverImage string `json:"cover_image" binding:"max:255"` // 封面图片
	ViewCount int `json:"view_count" binding:"required"` // 浏览次数
	LikeCount int `json:"like_count" binding:"required"` // 点赞次数
	CommentCount int `json:"comment_count" binding:"required"` // 评论次数
	IsPublished int8 `json:"is_published" binding:"required"` // 是否发布: 1-已发布, 0-草稿
	IsTop int8 `json:"is_top" binding:"required"` // 是否置顶: 1-是, 0-否
	IsFeatured int8 `json:"is_featured" binding:"required"` // 是否推荐: 1-是, 0-否
	PublishedAt *time.Time `json:"published_at"` // 发布时间
	Tags string `json:"tags" binding:"max:500"` // 标签(逗号分隔)
	SeoTitle string `json:"seo_title" binding:"max:200"` // SEO标题
	SeoKeywords string `json:"seo_keywords" binding:"max:500"` // SEO关键词
	SeoDescription string `json:"seo_description" binding:"max:500"` // SEO描述
	Status int8 `json:"status" binding:"required"` // 状态: 1-正常, 0-禁用
}

// ArticlesUpdateRequest 更新Articles管理请求
type ArticlesUpdateRequest struct {
	Title string `json:"title" binding:"required,max:200"` // 文章标题
	Slug string `json:"slug" binding:"max:200"` // URL别名
	AuthorId int64 `json:"author_id" binding:"required"` // 作者ID
	AuthorName string `json:"author_name" binding:"max:50"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	Summary string `json:"summary" binding:"max:500"` // 文章摘要
	Content string `json:"content" binding:"required"` // 文章内容
	CoverImage string `json:"cover_image" binding:"max:255"` // 封面图片
	ViewCount int `json:"view_count" binding:"required"` // 浏览次数
	LikeCount int `json:"like_count" binding:"required"` // 点赞次数
	CommentCount int `json:"comment_count" binding:"required"` // 评论次数
	IsPublished int8 `json:"is_published" binding:"required"` // 是否发布: 1-已发布, 0-草稿
	IsTop int8 `json:"is_top" binding:"required"` // 是否置顶: 1-是, 0-否
	IsFeatured int8 `json:"is_featured" binding:"required"` // 是否推荐: 1-是, 0-否
	PublishedAt *time.Time `json:"published_at"` // 发布时间
	Tags string `json:"tags" binding:"max:500"` // 标签(逗号分隔)
	SeoTitle string `json:"seo_title" binding:"max:200"` // SEO标题
	SeoKeywords string `json:"seo_keywords" binding:"max:500"` // SEO关键词
	SeoDescription string `json:"seo_description" binding:"max:500"` // SEO描述
	Status int8 `json:"status" binding:"required"` // 状态: 1-正常, 0-禁用
}

// ArticlesResponse Articles管理响应
type ArticlesResponse struct {
	Id int64 `json:"id"` // 文章ID
	Title string `json:"title"` // 文章标题
	Slug string `json:"slug"` // URL别名
	AuthorId int64 `json:"author_id"` // 作者ID
	AuthorName string `json:"author_name"` // 作者名称
	CategoryId *int64 `json:"category_id"` // 分类ID
	Summary string `json:"summary"` // 文章摘要
	CoverImage string `json:"cover_image"` // 封面图片
	ViewCount int `json:"view_count"` // 浏览次数
	LikeCount int `json:"like_count"` // 点赞次数
	CommentCount int `json:"comment_count"` // 评论次数
	IsPublished int8 `json:"is_published"` // 是否发布: 1-已发布, 0-草稿
	IsTop int8 `json:"is_top"` // 是否置顶: 1-是, 0-否
	IsFeatured int8 `json:"is_featured"` // 是否推荐: 1-是, 0-否
	PublishedAt *time.Time `json:"published_at"` // 发布时间
	Tags string `json:"tags"` // 标签(逗号分隔)
	SeoTitle string `json:"seo_title"` // SEO标题
	SeoKeywords string `json:"seo_keywords"` // SEO关键词
	SeoDescription string `json:"seo_description"` // SEO描述
	Status int8 `json:"status"` // 状态: 1-正常, 0-禁用
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// ToArticlesResponse 转换为响应对象
func ToArticlesResponse(articles *Articles) *ArticlesResponse {
	if articles == nil {
		return nil
	}
	
	return &ArticlesResponse{
		Id: articles.Id,
		Title: articles.Title,
		Slug: articles.Slug,
		AuthorId: articles.AuthorId,
		AuthorName: articles.AuthorName,
		CategoryId: articles.CategoryId,
		Summary: articles.Summary,
		CoverImage: articles.CoverImage,
		ViewCount: articles.ViewCount,
		LikeCount: articles.LikeCount,
		CommentCount: articles.CommentCount,
		IsPublished: articles.IsPublished,
		IsTop: articles.IsTop,
		IsFeatured: articles.IsFeatured,
		PublishedAt: articles.PublishedAt,
		Tags: articles.Tags,
		SeoTitle: articles.SeoTitle,
		SeoKeywords: articles.SeoKeywords,
		SeoDescription: articles.SeoDescription,
		Status: articles.Status,
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
