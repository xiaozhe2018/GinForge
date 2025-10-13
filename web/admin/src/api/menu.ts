import request from './index'

// 菜单相关API
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
  created_at: string
  updated_at: string
  parent?: Menu
  children?: Menu[]
}

export interface MenuListParams {
  page?: number
  page_size?: number
  keyword?: string
  type?: string
  status?: number
  parent_id?: number
}

export interface MenuListResponse {
  list: Menu[]
  total: number
}

export interface MenuTreeResponse {
  list: Menu[]
}

export interface CreateMenuParams {
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

export interface UpdateMenuParams {
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

// 获取菜单列表
export const getMenuList = (params: MenuListParams) => {
  return request.get<MenuListResponse>('/api/v1/admin/menus', { params })
}

// 获取菜单树
export const getMenuTree = () => {
  return request.get<MenuTreeResponse>('/api/v1/admin/menus/tree')
}

// 创建菜单
export const createMenu = (data: CreateMenuParams) => {
  return request.post('/api/v1/admin/menus', data)
}

// 获取菜单详情
export const getMenuDetail = (id: number) => {
  return request.get<Menu>(`/api/v1/admin/menus/${id}`)
}

// 更新菜单
export const updateMenu = (id: number, data: UpdateMenuParams) => {
  return request.put(`/api/v1/admin/menus/${id}`, data)
}

// 删除菜单
export const deleteMenu = (id: number) => {
  return request.delete(`/api/v1/admin/menus/${id}`)
}