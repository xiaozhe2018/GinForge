import api from './index'

// 商户相关API
export interface Merchant {
  id: string
  name: string
  email: string
  phone: string
  address: string
  status: number
  created_at: string
  updated_at: string
}

export interface MerchantListParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
}

export interface MerchantListResponse {
  list: Merchant[]
  total: number
  page: number
  page_size: number
}

// 获取商户列表
export const getMerchantList = (params: MerchantListParams = {}) => {
  return api.get<MerchantListResponse>('/v1/admin/merchant/list', { params })
}

// 更新商户状态
export const updateMerchantStatus = (id: string, status: number) => {
  return api.put(`/v1/admin/merchant/${id}/status`, { status })
}

// 获取商户详情
export const getMerchantDetail = (id: string) => {
  return api.get<Merchant>(`/v1/admin/merchant/${id}`)
}

// 创建商户
export const createMerchant = (data: Partial<Merchant>) => {
  return api.post<Merchant>('/v1/admin/merchant', data)
}

// 更新商户
export const updateMerchant = (id: string, data: Partial<Merchant>) => {
  return api.put<Merchant>(`/v1/admin/merchant/${id}`, data)
}

// 删除商户
export const deleteMerchant = (id: string) => {
  return api.delete(`/v1/admin/merchant/${id}`)
}
