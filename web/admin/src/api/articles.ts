import request from '@/utils/request'

// ========== 类型定义 ==========

/**
 * Articles管理
 */
export interface Articles {
  id: number // 文章ID
  title: string // 文章标题
  slug: string // URL别名
  author_id: number // 作者ID
  author_name: string // 作者名称
  category_id: number // 分类ID
  summary: string // 文章摘要
  cover_image: string // 封面图片
  view_count: number // 浏览次数
  like_count: number // 点赞次数
  comment_count: number // 评论次数
  is_published: number // 是否发布: 1-已发布, 0-草稿
  is_top: number // 是否置顶: 1-是, 0-否
  is_featured: number // 是否推荐: 1-是, 0-否
  published_at: string // 发布时间
  tags: string // 标签(逗号分隔)
  seo_title: string // SEO标题
  seo_keywords: string // SEO关键词
  seo_description: string // SEO描述
  status: number // 状态: 1-正常, 0-禁用
  created_at: string // 创建时间
  updated_at: string // 更新时间
}

/**
 * Articles管理列表请求参数
 */
export interface ArticlesListParams {
  page?: number
  page_size?: number
  keyword?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/**
 * Articles管理列表响应
 */
export interface ArticlesListResponse {
  list: Articles[]
  total: number
  page: number
  page_size: number
}

/**
 * 创建Articles管理请求参数
 */
export interface ArticlesCreateParams {
  title: string // 文章标题
  slug?: string // URL别名
  author_id: number // 作者ID
  author_name?: string // 作者名称
  category_id?: number // 分类ID
  summary?: string // 文章摘要
  content: string // 文章内容
  cover_image?: string // 封面图片
  view_count: number // 浏览次数
  like_count: number // 点赞次数
  comment_count: number // 评论次数
  is_published: number // 是否发布: 1-已发布, 0-草稿
  is_top: number // 是否置顶: 1-是, 0-否
  is_featured: number // 是否推荐: 1-是, 0-否
  published_at?: string // 发布时间
  tags?: string // 标签(逗号分隔)
  seo_title?: string // SEO标题
  seo_keywords?: string // SEO关键词
  seo_description?: string // SEO描述
  status: number // 状态: 1-正常, 0-禁用
}

/**
 * 更新Articles管理请求参数
 */
export interface ArticlesUpdateParams {
  title?: string // 文章标题
  slug?: string // URL别名
  author_id?: number // 作者ID
  author_name?: string // 作者名称
  category_id?: number // 分类ID
  summary?: string // 文章摘要
  content?: string // 文章内容
  cover_image?: string // 封面图片
  view_count?: number // 浏览次数
  like_count?: number // 点赞次数
  comment_count?: number // 评论次数
  is_published?: number // 是否发布: 1-已发布, 0-草稿
  is_top?: number // 是否置顶: 1-是, 0-否
  is_featured?: number // 是否推荐: 1-是, 0-否
  published_at?: string // 发布时间
  tags?: string // 标签(逗号分隔)
  seo_title?: string // SEO标题
  seo_keywords?: string // SEO关键词
  seo_description?: string // SEO描述
  status?: number // 状态: 1-正常, 0-禁用
}

// ========== API 方法 ==========

/**
 * 获取Articles管理列表
 */
export const getArticlesList = (params?: ArticlesListParams) => {
  return request.get<ArticlesListResponse>('/api/v1/admin/articleses', { params })
}

/**
 * 获取Articles管理详情
 */
export const getArticles = (id: number) => {
  return request.get<Articles>(`/api/v1/admin/articleses/${id}`)
}

/**
 * 创建Articles管理
 */
export const createArticles = (data: ArticlesCreateParams) => {
  return request.post<Articles>('/api/v1/admin/articleses', data)
}

/**
 * 更新Articles管理
 */
export const updateArticles = (id: number, data: ArticlesUpdateParams) => {
  return request.put(`/api/v1/admin/articleses/${id}`, data)
}

/**
 * 删除Articles管理
 */
export const deleteArticles = (id: number) => {
  return request.delete(`/api/v1/admin/articleses/${id}`)
}
