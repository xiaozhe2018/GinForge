import request from './index'

// ========== 类型定义 ==========

/**
 * Articles管理
 */
export interface Articles {
  id: number // 文章ID
  title: string // 标题
  summary: string // 摘要
  author_id: number // 作者ID
  author_name: string // 作者名称
  category_id: number // 分类ID
  cover_image: string // 封面图片
  tags: string // 标签（逗号分隔）
  status: number // 状态:0草稿,1已发布,2已下线
  is_top: number // 是否置顶
  view_count: number // 浏览次数
  like_count: number // 点赞次数
  comment_count: number // 评论次数
  published_at: string // 发布时间
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
  title: string // 标题
  content: string // 内容
  summary?: string // 摘要
  author_id: number // 作者ID
  author_name?: string // 作者名称
  category_id?: number // 分类ID
  cover_image?: string // 封面图片
  tags?: string // 标签（逗号分隔）
  status: number // 状态:0草稿,1已发布,2已下线
  is_top: number // 是否置顶
  view_count: number // 浏览次数
  like_count: number // 点赞次数
  comment_count: number // 评论次数
  published_at?: string // 发布时间
}

/**
 * 更新Articles管理请求参数
 */
export interface ArticlesUpdateParams {
  title?: string // 标题
  content?: string // 内容
  summary?: string // 摘要
  author_id?: number // 作者ID
  author_name?: string // 作者名称
  category_id?: number // 分类ID
  cover_image?: string // 封面图片
  tags?: string // 标签（逗号分隔）
  status?: number // 状态:0草稿,1已发布,2已下线
  is_top?: number // 是否置顶
  view_count?: number // 浏览次数
  like_count?: number // 点赞次数
  comment_count?: number // 评论次数
  published_at?: string // 发布时间
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
