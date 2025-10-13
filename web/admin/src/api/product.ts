import api from './index'

// 商品相关API
export interface Product {
  id: string
  name: string
  description: string
  price: number
  stock: number
  category: string
  images: string[]
  status: number
  merchant_id: string
  created_at: string
  updated_at: string
}

export interface ProductListParams {
  page?: number
  page_size?: number
  keyword?: string
  category?: string
  status?: number
  merchant_id?: string
}

export interface ProductListResponse {
  list: Product[]
  total: number
  page: number
  page_size: number
}

// 获取商品列表
export const getProductList = (params: ProductListParams = {}) => {
  return api.get<ProductListResponse>('/v1/admin/product/list', { params })
}

// 获取商品详情
export const getProductDetail = (id: string) => {
  return api.get<Product>(`/v1/admin/product/${id}`)
}

// 创建商品
export const createProduct = (data: Partial<Product>) => {
  return api.post<Product>('/v1/admin/product', data)
}

// 更新商品
export const updateProduct = (id: string, data: Partial<Product>) => {
  return api.put<Product>(`/v1/admin/product/${id}`, data)
}

// 删除商品
export const deleteProduct = (id: string) => {
  return api.delete(`/v1/admin/product/${id}`)
}

// 更新商品状态
export const updateProductStatus = (id: string, status: number) => {
  return api.put(`/v1/admin/product/${id}/status`, { status })
}
