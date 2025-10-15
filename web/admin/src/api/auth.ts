import request from './index'

// 登录参数
interface LoginParams {
  username: string
  password: string
  captcha?: string
}

// 登录响应
interface LoginResponse {
  token: string
  user: User
  menus: Menu[]
  permissions: string[]
}

// 用户信息
interface User {
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

// 角色信息
interface Role {
  id: number
  name: string
  code: string
  description?: string
  sort: number
  status: number
}

// 菜单信息
interface Menu {
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
  children?: Menu[]
}

// 权限信息
interface Permission {
  id: number
  name: string
  code: string
  type: string
  description?: string
  created_at: string
  updated_at: string
}

// 更新个人信息参数
interface UpdateProfileParams {
  email: string
  phone?: string
  name: string
  status: number
  role_ids: number[]
  avatar?: string
}

// 修改密码参数
interface ChangePasswordParams {
  old_password: string
  new_password: string
}

// 导出所有类型
export type {
  LoginParams,
  LoginResponse,
  User,
  Role,
  Menu,
  Permission,
  UpdateProfileParams,
  ChangePasswordParams
}

// 用户登录
export const login = (data: LoginParams) => {
  return request.post<LoginResponse>('/api/v1/admin/login', data)
}

// 用户登出
export const logout = () => {
  return request.post('/api/v1/admin/logout')
}

// 获取当前用户信息
export const getProfile = () => {
  return request.get<User>('/api/v1/admin/profile')
}

// 更新当前用户信息
export const updateProfile = (data: UpdateProfileParams) => {
  return request.put('/api/v1/admin/profile', data)
}

// 修改密码
export const changePassword = (data: ChangePasswordParams) => {
  return request.put('/api/v1/admin/change-password', data)
}
