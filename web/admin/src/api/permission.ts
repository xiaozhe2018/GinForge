import request from './index'

// 权限相关API
export interface Permission {
  id: number
  name: string
  code: string
  type: string
  description?: string
  status: number  // 状态：0=禁用，1=启用
  created_at: string
  updated_at: string
}

export interface PermissionListParams {
  page?: number
  page_size?: number
  keyword?: string
  type?: string
}

export interface PermissionListResponse {
  list: Permission[]
  total: number
}

export interface CreatePermissionParams {
  name: string
  code: string
  type: string
  description?: string
  status?: number  // 状态：0=禁用，1=启用
}

export interface UpdatePermissionParams {
  name: string
  code: string
  type: string
  description?: string
  status?: number  // 状态：0=禁用，1=启用
}

// 获取权限列表
export const getPermissionList = (params: PermissionListParams) => {
  return request.get<PermissionListResponse>('/api/v1/admin/permissions', { params })
}

// 创建权限
export const createPermission = (data: CreatePermissionParams) => {
  return request.post('/api/v1/admin/permissions', data)
}

// 获取权限详情
export const getPermissionDetail = (id: number) => {
  return request.get<Permission>(`/api/v1/admin/permissions/${id}`)
}

// 更新权限
export const updatePermission = (id: number, data: UpdatePermissionParams) => {
  return request.put(`/api/v1/admin/permissions/${id}`, data)
}

// 删除权限
export const deletePermission = (id: number) => {
  return request.delete(`/api/v1/admin/permissions/${id}`)
}

// 更新权限状态
export const updatePermissionStatus = (id: number, status: number) => {
  return request.put(`/api/v1/admin/permissions/${id}/status?status=${status}`)
}