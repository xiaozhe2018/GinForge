import request from './index'

// 用户相关API
export interface User {
  id: number
  username: string
  email: string
  phone?: string
  name?: string
  avatar?: string
  status: number
  last_login_at?: string
  last_login_ip?: string
  created_at: string
  updated_at: string
  roles?: Role[]
}

export interface UserListParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
  role_id?: number
}

export interface UserListResponse {
  list: User[]
  total: number
}

export interface CreateUserParams {
  username: string
  email: string
  phone?: string
  password: string
  name: string
  role_ids: number[]
}

export interface UpdateUserParams {
  email: string
  phone?: string
  name: string
  status: number
  role_ids: number[]
}

export interface Role {
  id: number
  name: string
  code: string
  description?: string
  sort: number
  status: number
}

// 获取用户列表
export const getUserList = (params: UserListParams) => {
  return request.get<UserListResponse>('/api/v1/admin/users', { params })
}

// 创建用户
export const createUser = (data: CreateUserParams) => {
  return request.post('/api/v1/admin/users', data)
}

// 获取用户详情
export const getUserDetail = (id: number) => {
  return request.get<User>(`/api/v1/admin/users/${id}`)
}

// 更新用户
export const updateUser = (id: number, data: UpdateUserParams) => {
  return request.put(`/api/v1/admin/users/${id}`, data)
}

// 更新用户状态
export const updateUserStatus = (id: number, status: number) => {
  return request.put(`/api/v1/admin/users/${id}/status`, null, { params: { status } })
}

// 删除用户
export const deleteUser = (id: number) => {
  return request.delete(`/api/v1/admin/users/${id}`)
}

// 重置用户密码
export const resetUserPassword = (id: number, password: string) => {
  return request.put(`/api/v1/admin/users/${id}/reset-password`, { password })
}