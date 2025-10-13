import request from './index'

// 角色相关API
export interface Role {
  id: number
  name: string
  code: string
  description?: string
  sort: number
  status: number
  created_at: string
  updated_at: string
  permissions?: Permission[]
  menus?: Menu[]
}

export interface RoleListParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
}

export interface RoleListResponse {
  list: Role[]
  total: number
}

export interface CreateRoleParams {
  name: string
  code: string
  description?: string
  sort: number
  status: number
  permission_ids: number[]
  menu_ids: number[]
}

export interface UpdateRoleParams {
  name: string
  code: string
  description?: string
  sort: number
  status: number
  permission_ids: number[]
  menu_ids: number[]
}

export interface Permission {
  id: number
  name: string
  code: string
  type: string
  description?: string
}

export interface Menu {
  id: number
  parent_id: number
  name: string
  code: string
  type: string
  path?: string
  component?: string
  icon?: string
  sort: number
  visible: number
  status: number
  description?: string
}

// 获取角色列表
export const getRoleList = (params: RoleListParams) => {
  return request.get<RoleListResponse>('/api/v1/admin/roles', { params })
}

// 创建角色
export const createRole = (data: CreateRoleParams) => {
  return request.post('/api/v1/admin/roles', data)
}

// 获取角色详情
export const getRoleDetail = (id: number) => {
  return request.get<Role>(`/api/v1/admin/roles/${id}`)
}

// 更新角色
export const updateRole = (id: number, data: UpdateRoleParams) => {
  return request.put(`/api/v1/admin/roles/${id}`, data)
}

// 删除角色
export const deleteRole = (id: number) => {
  return request.delete(`/api/v1/admin/roles/${id}`)
}