import request from './index'

// 系统相关API
export interface SystemInfo {
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  online_users: number
  total_users: number
  total_requests: number
  error_count: number
}

export interface SystemConfig {
  id: number
  key: string
  value: string
  type: string
  description?: string
  group: string
  sort: number
  created_at: string
  updated_at: string
}

export interface SystemConfigListParams {
  page?: number
  page_size?: number
  keyword?: string
  group?: string
}

export interface SystemConfigListResponse {
  list: SystemConfig[]
  total: number
}

export interface UpdateSystemConfigParams {
  value: string
}

export interface SystemLog {
  id: number
  user_id?: number
  username?: string
  method: string
  path: string
  ip: string
  user_agent: string
  status_code: number
  duration: number
  created_at: string
}

export interface SystemLogListParams {
  page?: number
  page_size?: number
  username?: string
  method?: string
  status_code?: number
  start_time?: string
  end_time?: string
}

export interface SystemLogListResponse {
  list: SystemLog[]
  total: number
}

// 获取系统基本信息（公开接口，不需要登录）
export const getSystemBasicInfo = () => {
  return request.get<Record<string, string>>('/api/v1/admin/system/basic-info')
}

// 获取系统信息
export const getSystemInfo = () => {
  return request.get<SystemInfo>('/api/v1/admin/system/info')
}

// 获取系统配置列表
export const getSystemConfigList = (params: SystemConfigListParams) => {
  return request.get<SystemConfigListResponse>('/api/v1/admin/system/configs', { params })
}

// 更新系统配置
export const updateSystemConfig = (key: string, data: UpdateSystemConfigParams) => {
  return request.put(`/api/v1/admin/system/configs/${key}`, data)
}

// 获取系统日志列表
export const getSystemLogList = (params: SystemLogListParams) => {
  return request.get<SystemLogListResponse>('/api/v1/admin/system/logs', { params })
}

// 清空系统日志
export const clearSystemLogs = () => {
  return request.delete('/api/v1/admin/system/logs')
}

