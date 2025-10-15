# 实战：创建完整的业务模块

从零到一，完整演示如何在 GinForge 中创建一个文章管理模块。

## 🎯 目标

创建一个包含以下功能的文章管理系统：

- ✅ 文章 CRUD（创建、读取、更新、删除）
- ✅ 文章分类
- ✅ 文章标签
- ✅ 文章发布/下线
- ✅ 文章浏览统计
- ✅ 评论功能

## 📋 步骤概览

```
1. 定义数据模型
2. 创建 Repository（数据访问层）
3. 创建 Service（业务逻辑层）
4. 创建 Handler（HTTP 处理层）
5. 注册路由
6. 测试 API
7. 集成前端
```

## 📦 步骤 1：定义数据模型

创建文件：`services/admin-api/internal/model/article.go`

```go
package model

import (
    "time"
)

// Article 文章模型
type Article struct {
    ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
    Title       string     `json:"title" gorm:"type:varchar(200);not null;index"`
    Content     string     `json:"content" gorm:"type:longtext"`
    Summary     string     `json:"summary" gorm:"type:varchar(500)"`
    Cover       string     `json:"cover" gorm:"type:varchar(255)"`
    CategoryID  uint64     `json:"category_id" gorm:"index"`
    AuthorID    uint64     `json:"author_id" gorm:"index"`
    Status      int8       `json:"status" gorm:"type:tinyint(1);default:0;index;comment:状态:0=草稿,1=已发布,2=已下线"`
    ViewCount   uint64     `json:"view_count" gorm:"default:0"`
    LikeCount   uint64     `json:"like_count" gorm:"default:0"`
    PublishedAt *time.Time `json:"published_at"`
    CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
    
    // 关联
    Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    Author   *AdminUser `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
    Tags     []Tag `json:"tags,omitempty" gorm:"many2many:article_tags;"`
}

// TableName 指定表名
func (Article) TableName() string {
    return "articles"
}

// Category 分类模型
type Category struct {
    ID        uint64     `json:"id" gorm:"primaryKey"`
    Name      string     `json:"name" gorm:"type:varchar(50);not null"`
    Sort      int        `json:"sort" gorm:"default:0"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

func (Category) TableName() string {
    return "article_categories"
}

// Tag 标签模型
type Tag struct {
    ID        uint64     `json:"id" gorm:"primaryKey"`
    Name      string     `json:"name" gorm:"type:varchar(50);uniqueIndex"`
    CreatedAt time.Time  `json:"created_at"`
}

func (Tag) TableName() string {
    return "article_tags"
}

// 请求和响应模型

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
    Title      string   `json:"title" binding:"required,min=2,max=200"`
    Content    string   `json:"content" binding:"required"`
    Summary    string   `json:"summary" binding:"max=500"`
    Cover      string   `json:"cover"`
    CategoryID uint64   `json:"category_id" binding:"required"`
    TagIDs     []uint64 `json:"tag_ids"`
    Status     int8     `json:"status" binding:"oneof=0 1 2"`
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
    Title      string   `json:"title" binding:"required,min=2,max=200"`
    Content    string   `json:"content" binding:"required"`
    Summary    string   `json:"summary" binding:"max=500"`
    Cover      string   `json:"cover"`
    CategoryID uint64   `json:"category_id" binding:"required"`
    TagIDs     []uint64 `json:"tag_ids"`
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
    Page       int    `form:"page"`
    PageSize   int    `form:"page_size" binding:"max=100"`
    Keyword    string `form:"keyword"`
    CategoryID uint64 `form:"category_id"`
    Status     int8   `form:"status"`
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
    List  []Article `json:"list"`
    Total int64     `json:"total"`
}
```

## 🗄️ 步骤 2：创建 Repository

创建文件：`services/admin-api/internal/repository/article_repository.go`

```go
package repository

import (
    "gorm.io/gorm"
    "goweb/services/admin-api/internal/model"
)

type ArticleRepository struct {
    db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
    return &ArticleRepository{db: db}
}

// Create 创建文章
func (r *ArticleRepository) Create(article *model.Article) error {
    return r.db.Create(article).Error
}

// GetByID 根据 ID 获取文章
func (r *ArticleRepository) GetByID(id uint64) (*model.Article, error) {
    var article model.Article
    err := r.db.Preload("Category").
        Preload("Author").
        Preload("Tags").
        First(&article, id).Error
    return &article, err
}

// List 获取文章列表
func (r *ArticleRepository) List(req *model.ArticleListRequest) ([]model.Article, int64, error) {
    var articles []model.Article
    var total int64
    
    query := r.db.Model(&model.Article{})
    
    // 搜索条件
    if req.Keyword != "" {
        query = query.Where("title LIKE ? OR content LIKE ?", 
            "%"+req.Keyword+"%", "%"+req.Keyword+"%")
    }
    
    if req.CategoryID > 0 {
        query = query.Where("category_id = ?", req.CategoryID)
    }
    
    if req.Status >= 0 {
        query = query.Where("status = ?", req.Status)
    }
    
    // 计数
    query.Count(&total)
    
    // 分页
    offset := (req.Page - 1) * req.PageSize
    err := query.Preload("Category").
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "username", "name")
        }).
        Preload("Tags").
        Offset(offset).
        Limit(req.PageSize).
        Order("created_at DESC").
        Find(&articles).Error
    
    return articles, total, err
}

// Update 更新文章
func (r *ArticleRepository) Update(article *model.Article) error {
    return r.db.Save(article).Error
}

// UpdateStatus 更新文章状态
func (r *ArticleRepository) UpdateStatus(id uint64, status int8) error {
    return r.db.Model(&model.Article{}).
        Where("id = ?", id).
        Update("status", status).Error
}

// IncrementViewCount 增加浏览次数
func (r *ArticleRepository) IncrementViewCount(id uint64) error {
    return r.db.Model(&model.Article{}).
        Where("id = ?", id).
        Update("view_count", gorm.Expr("view_count + 1")).Error
}

// Delete 删除文章
func (r *ArticleRepository) Delete(id uint64) error {
    return r.db.Delete(&model.Article{}, id).Error
}

// UpdateTags 更新文章标签
func (r *ArticleRepository) UpdateTags(articleID uint64, tagIDs []uint64) error {
    var article model.Article
    if err := r.db.First(&article, articleID).Error; err != nil {
        return err
    }
    
    var tags []model.Tag
    if len(tagIDs) > 0 {
        r.db.Find(&tags, tagIDs)
    }
    
    return r.db.Model(&article).Association("Tags").Replace(tags)
}
```

## 💼 步骤 3：创建 Service

创建文件：`services/admin-api/internal/service/article_service.go`

```go
package service

import (
    "context"
    "errors"
    "time"
    
    "gorm.io/gorm"
    "goweb/pkg/config"
    pkgRedis "goweb/pkg/redis"
    "goweb/services/admin-api/internal/model"
    "goweb/services/admin-api/internal/repository"
)

type ArticleService struct {
    *AdminService
    articleRepo *repository.ArticleRepository
}

func NewArticleService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *ArticleService {
    return &ArticleService{
        AdminService: NewAdminService(db, cfg, redisClient),
        articleRepo:  repository.NewArticleRepository(db),
    }
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req *model.ArticleCreateRequest, authorID uint64) (*model.Article, error) {
    article := &model.Article{
        Title:      req.Title,
        Content:    req.Content,
        Summary:    req.Summary,
        Cover:      req.Cover,
        CategoryID: req.CategoryID,
        AuthorID:   authorID,
        Status:     req.Status,
    }
    
    // 如果是发布状态，设置发布时间
    if article.Status == 1 {
        now := time.Now()
        article.PublishedAt = &now
    }
    
    // 创建文章
    if err := s.articleRepo.Create(article); err != nil {
        s.logger.Error("create article failed", err)
        return nil, err
    }
    
    // 更新标签
    if len(req.TagIDs) > 0 {
        if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
            s.logger.Error("update article tags failed", err)
        }
    }
    
    s.logger.Info("article created", "id", article.ID, "title", article.Title)
    return article, nil
}

// GetArticle 获取文章详情
func (s *ArticleService) GetArticle(id uint64, incrementView bool) (*model.Article, error) {
    article, err := s.articleRepo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("文章不存在")
        }
        return nil, err
    }
    
    // 增加浏览次数
    if incrementView {
        go func() {
            s.articleRepo.IncrementViewCount(id)
        }()
    }
    
    return article, nil
}

// ListArticles 获取文章列表
func (s *ArticleService) ListArticles(req *model.ArticleListRequest) (*model.ArticleListResponse, error) {
    // 设置默认值
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    
    articles, total, err := s.articleRepo.List(req)
    if err != nil {
        return nil, err
    }
    
    return &model.ArticleListResponse{
        List:  articles,
        Total: total,
    }, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id uint64, req *model.ArticleUpdateRequest) error {
    article, err := s.articleRepo.GetByID(id)
    if err != nil {
        return errors.New("文章不存在")
    }
    
    // 更新字段
    article.Title = req.Title
    article.Content = req.Content
    article.Summary = req.Summary
    article.Cover = req.Cover
    article.CategoryID = req.CategoryID
    
    if err := s.articleRepo.Update(article); err != nil {
        s.logger.Error("update article failed", err)
        return err
    }
    
    // 更新标签
    if err := s.articleRepo.UpdateTags(id, req.TagIDs); err != nil {
        s.logger.Error("update article tags failed", err)
    }
    
    s.logger.Info("article updated", "id", id)
    return nil
}

// PublishArticle 发布文章
func (s *ArticleService) PublishArticle(id uint64) error {
    article, err := s.articleRepo.GetByID(id)
    if err != nil {
        return errors.New("文章不存在")
    }
    
    article.Status = 1
    now := time.Now()
    article.PublishedAt = &now
    
    if err := s.articleRepo.Update(article); err != nil {
        s.logger.Error("publish article failed", err)
        return err
    }
    
    s.logger.Info("article published", "id", id)
    return nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id uint64) error {
    return s.articleRepo.Delete(id)
}
```

## 🎮 步骤 4：创建 Handler

创建文件：`services/admin-api/internal/handler/article_handler.go`

```go
package handler

import (
    "strconv"
    
    "github.com/gin-gonic/gin"
    "goweb/pkg/logger"
    "goweb/pkg/response"
    "goweb/services/admin-api/internal/model"
    "goweb/services/admin-api/internal/service"
)

type ArticleHandler struct {
    articleService *service.ArticleService
    logger         logger.Logger
}

func NewArticleHandler(articleService *service.ArticleService, log logger.Logger) *ArticleHandler {
    return &ArticleHandler{
        articleService: articleService,
        logger:         log,
    }
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Description 创建新的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body model.ArticleCreateRequest true "创建文章请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req model.ArticleCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("bind request failed", err)
        response.BadRequest(c, "请求参数错误")
        return
    }
    
    // 获取当前用户 ID
    authorID, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
    
    article, err := h.articleService.CreateArticle(&req, authorID)
    if err != nil {
        h.logger.Error("create article failed", err)
        response.InternalError(c, err.Error())
        return
    }
    
    response.Success(c, article)
}

// GetArticles 获取文章列表
// @Summary 获取文章列表
// @Description 获取文章列表，支持搜索和过滤
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=model.ArticleListResponse}
// @Router /api/v1/admin/articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
    var req model.ArticleListRequest
    c.ShouldBindQuery(&req)
    
    result, err := h.articleService.ListArticles(&req)
    if err != nil {
        h.logger.Error("get articles failed", err)
        response.InternalError(c, "获取文章列表失败")
        return
    }
    
    response.Success(c, result)
}

// GetArticle 获取文章详情
// @Summary 获取文章详情
// @Description 根据 ID 获取文章详细信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=model.Article}
// @Router /api/v1/admin/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        response.BadRequest(c, "无效的文章ID")
        return
    }
    
    article, err := h.articleService.GetArticle(id, true)
    if err != nil {
        h.logger.Error("get article failed", err)
        response.InternalError(c, err.Error())
        return
    }
    
    response.Success(c, article)
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body model.ArticleUpdateRequest true "更新文章请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        response.BadRequest(c, "无效的文章ID")
        return
    }
    
    var req model.ArticleUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("bind request failed", err)
        response.BadRequest(c, "请求参数错误")
        return
    }
    
    if err := h.articleService.UpdateArticle(id, &req); err != nil {
        h.logger.Error("update article failed", err)
        response.InternalError(c, err.Error())
        return
    }
    
    response.Success(c, gin.H{"message": "更新成功"})
}

// PublishArticle 发布文章
// @Summary 发布文章
// @Description 将文章状态改为已发布
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articles/{id}/publish [put]
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        response.BadRequest(c, "无效的文章ID")
        return
    }
    
    if err := h.articleService.PublishArticle(id); err != nil {
        h.logger.Error("publish article failed", err)
        response.InternalError(c, err.Error())
        return
    }
    
    response.Success(c, gin.H{"message": "发布成功"})
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 删除指定的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        response.BadRequest(c, "无效的文章ID")
        return
    }
    
    if err := h.articleService.DeleteArticle(id); err != nil {
        h.logger.Error("delete article failed", err)
        response.InternalError(c, "删除文章失败")
        return
    }
    
    response.Success(c, gin.H{"message": "删除成功"})
}
```

## 🛣️ 步骤 5：注册路由

编辑文件：`services/admin-api/internal/router/router.go`

```go
// 创建文章 Handler
articleRepo := repository.NewArticleRepository(database)
articleService := service.NewArticleService(database, cfg, redisClient)
articleHandler := handler.NewArticleHandler(articleService, log)
articleHandler.SetLogger(log)

// 注册文章路由
auth.GET("/articles", articleHandler.GetArticles)
auth.GET("/articles/:id", articleHandler.GetArticle)
auth.POST("/articles", articleHandler.CreateArticle)
auth.PUT("/articles/:id", articleHandler.UpdateArticle)
auth.PUT("/articles/:id/publish", articleHandler.PublishArticle)
auth.DELETE("/articles/:id", articleHandler.DeleteArticle)
```

## 🗄️ 步骤 6：数据库迁移

创建 SQL 文件：`database/migrations/003_create_article_tables.sql`

```sql
-- 创建分类表
CREATE TABLE IF NOT EXISTS `article_categories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL,
    `sort` int DEFAULT 0,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建标签表
CREATE TABLE IF NOT EXISTS `article_tags` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建文章表
CREATE TABLE IF NOT EXISTS `articles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(200) NOT NULL,
    `content` longtext,
    `summary` varchar(500) DEFAULT NULL,
    `cover` varchar(255) DEFAULT NULL,
    `category_id` bigint unsigned DEFAULT NULL,
    `author_id` bigint unsigned DEFAULT NULL,
    `status` tinyint(1) DEFAULT 0 COMMENT '状态:0=草稿,1=已发布,2=已下线',
    `view_count` bigint unsigned DEFAULT 0,
    `like_count` bigint unsigned DEFAULT 0,
    `published_at` datetime(3) DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_title` (`title`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建文章标签关联表
CREATE TABLE IF NOT EXISTS `article_tag_relations` (
    `article_id` bigint unsigned NOT NULL,
    `tag_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`article_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入示例分类
INSERT INTO `article_categories` (`name`, `sort`, `created_at`, `updated_at`) VALUES
('技术文章', 1, NOW(), NOW()),
('产品动态', 2, NOW(), NOW()),
('公司新闻', 3, NOW(), NOW());

-- 插入示例标签
INSERT INTO `article_tags` (`name`, `created_at`) VALUES
('Go', NOW()),
('Gin', NOW()),
('微服务', NOW()),
('教程', NOW());
```

执行迁移：

```bash
# MySQL
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/003_create_article_tables.sql

# 或在代码中使用 AutoMigrate
db.AutoMigrate(
    &model.Article{},
    &model.Category{},
    &model.Tag{},
)
```

## 🧪 步骤 7：测试 API

```bash
# 获取 Token
TOKEN=$(curl -s -X POST http://localhost:8083/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# 创建文章
curl -X POST http://localhost:8083/api/v1/admin/articles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "GinForge 快速入门",
    "content": "这是一篇关于 GinForge 框架的教程...",
    "summary": "学习如何使用 GinForge 框架",
    "category_id": 1,
    "tag_ids": [1, 2, 3],
    "status": 0
  }'

# 获取文章列表
curl -X GET "http://localhost:8083/api/v1/admin/articles?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 获取文章详情
curl -X GET "http://localhost:8083/api/v1/admin/articles/1" \
  -H "Authorization: Bearer $TOKEN"

# 发布文章
curl -X PUT "http://localhost:8083/api/v1/admin/articles/1/publish" \
  -H "Authorization: Bearer $TOKEN"

# 删除文章
curl -X DELETE "http://localhost:8083/api/v1/admin/articles/1" \
  -H "Authorization: Bearer $TOKEN"
```

## 🎨 步骤 8：前端集成

### 创建 API 接口

`web/admin/src/api/article.ts`:

```typescript
import request from './index'

export interface Article {
  id: number
  title: string
  content: string
  summary?: string
  cover?: string
  category_id: number
  status: number
  view_count: number
  created_at: string
  category?: { id: number; name: string }
  author?: { id: number; name: string }
  tags?: Array<{ id: number; name: string }>
}

export interface CreateArticleParams {
  title: string
  content: string
  summary?: string
  cover?: string
  category_id: number
  tag_ids?: number[]
  status: number
}

export const getArticleList = (params: any) => {
  return request.get<{ list: Article[]; total: number }>('/api/v1/admin/articles', { params })
}

export const getArticleDetail = (id: number) => {
  return request.get<Article>(`/api/v1/admin/articles/${id}`)
}

export const createArticle = (data: CreateArticleParams) => {
  return request.post('/api/v1/admin/articles', data)
}

export const updateArticle = (id: number, data: any) => {
  return request.put(`/api/v1/admin/articles/${id}`, data)
}

export const publishArticle = (id: number) => {
  return request.put(`/api/v1/admin/articles/${id}/publish`)
}

export const deleteArticle = (id: number) => {
  return request.delete(`/api/v1/admin/articles/${id}`)
}
```

### 创建文章列表页面

`web/admin/src/views/Articles.vue`:

```vue
<template>
  <div class="articles-page">
    <el-card>
      <!-- 搜索栏 -->
      <el-form :inline="true">
        <el-form-item>
          <el-input v-model="searchForm.keyword" placeholder="搜索文章" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadArticles">搜索</el-button>
          <el-button @click="handleCreate">创建文章</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 文章列表 -->
      <el-table :data="articles" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" width="300" />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column prop="author.name" label="作者" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="info">草稿</el-tag>
            <el-tag v-else-if="row.status === 1" type="success">已发布</el-tag>
            <el-tag v-else type="warning">已下线</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="success" @click="handlePublish(row)" v-if="row.status === 0">发布</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="total"
        @current-change="loadArticles"
        layout="total, prev, pager, next"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as articleApi from '@/api/article'

const loading = ref(false)
const articles = ref<any[]>([])
const total = ref(0)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

// 加载文章列表
const loadArticles = async () => {
  loading.value = true
  try {
    const result = await articleApi.getArticleList({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword
    })
    articles.value = result.list
    total.value = result.total
  } catch (error) {
    ElMessage.error('加载文章列表失败')
  } finally {
    loading.value = false
  }
}

// 创建文章
const handleCreate = () => {
  // 跳转到创建页面
  router.push('/dashboard/articles/create')
}

// 编辑文章
const handleEdit = (row: any) => {
  router.push(`/dashboard/articles/${row.id}/edit`)
}

// 发布文章
const handlePublish = async (row: any) => {
  try {
    await articleApi.publishArticle(row.id)
    ElMessage.success('发布成功')
    loadArticles()
  } catch (error) {
    ElMessage.error('发布失败')
  }
}

// 删除文章
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning'
    })
    
    await articleApi.deleteArticle(row.id)
    ElMessage.success('删除成功')
    loadArticles()
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  loadArticles()
})
</script>
```

## 🎯 完整功能清单

### 已实现的功能

- ✅ 文章创建、编辑、删除
- ✅ 文章发布、下线
- ✅ 文章分类管理
- ✅ 文章标签管理
- ✅ 浏览次数统计
- ✅ 搜索和过滤
- ✅ 分页查询
- ✅ 软删除

### 可扩展的功能

- 📝 富文本编辑器集成
- 🖼️ 图片上传和管理
- 💬 评论系统
- 👍 点赞功能
- 📊 数据统计
- 🔔 发布通知（WebSocket）

## 💡 最佳实践

### 1. 使用缓存

```go
// 缓存热门文章
func (s *ArticleService) GetHotArticles() ([]model.Article, error) {
    cacheKey := "articles:hot"
    ctx := context.Background()
    
    // 查缓存
    var articles []model.Article
    cached, err := s.redisClient.Get(ctx, cacheKey)
    if err == nil {
        json.Unmarshal([]byte(cached), &articles)
        return articles, nil
    }
    
    // 查数据库
    articles, _, err = s.articleRepo.List(&model.ArticleListRequest{
        Page:     1,
        PageSize: 10,
        Status:   1,
    })
    
    // 写缓存（10分钟）
    data, _ := json.Marshal(articles)
    s.redisClient.Set(ctx, cacheKey, string(data), 10*time.Minute)
    
    return articles, err
}
```

### 2. 使用消息队列

```go
// 发布文章时发送通知
func (s *ArticleService) PublishArticle(id uint64) error {
    // 更新状态
    if err := s.articleRepo.UpdateStatus(id, 1); err != nil {
        return err
    }
    
    // 发送通知（异步）
    go func() {
        s.notificationClient.SendSystemNotification(context.Background(), &notification.SystemNotificationRequest{
            Title: "新文章发布",
            Body:  "有新文章发布，快来查看吧！",
            Link:  fmt.Sprintf("/articles/%d", id),
        })
    }()
    
    return nil
}
```

## 🎯 总结

通过这个完整的实战案例，你学会了：

1. ✅ 如何设计数据模型和关联关系
2. ✅ 如何使用 Repository 模式封装数据访问
3. ✅ 如何在 Service 层实现业务逻辑
4. ✅ 如何创建 RESTful API
5. ✅ 如何编写 Swagger 文档
6. ✅ 如何集成前端页面
7. ✅ 如何使用缓存和消息队列

现在你可以用同样的方式创建任何业务模块！🚀

