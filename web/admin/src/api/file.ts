import request from './index'

/**
 * 上传文件响应
 */
export interface UploadFileResponse {
  id: number
  original_name: string
  file_name: string
  file_size: number
  file_type: string
  mime_type: string
  url: string
  hash: string
  upload_time: string
}

/**
 * 文件记录
 */
export interface FileRecord {
  id: number
  original_name: string
  file_name: string
  relative_path: string
  file_size: number
  mime_type: string
  file_type: string
  hash: string
  storage_type: string
  url: string
  uploaded_by: number
  user_type: string
  upload_ip: string
  upload_time: string
  download_count: number
  status: number
  description: string
  tags: string
  created_at: string
  updated_at: string
}

/**
 * 上传文件
 * @param file 文件对象
 * @param options 上传选项
 * @returns 上传结果
 */
export const uploadFile = (file: File, options?: {
  description?: string
  tags?: string
  sub_path?: string
  user_id?: number
  user_type?: string
}) => {
  const formData = new FormData()
  formData.append('file', file)
  
  // 添加可选参数
  if (options?.description) {
    formData.append('description', options.description)
  }
  if (options?.tags) {
    formData.append('tags', options.tags)
  }
  if (options?.sub_path) {
    formData.append('sub_path', options.sub_path)
  }
  if (options?.user_id) {
    formData.append('user_id', options.user_id.toString())
  }
  if (options?.user_type) {
    formData.append('user_type', options.user_type)
  }
  
  return request.post<UploadFileResponse>('/api/v1/files/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    baseURL: 'http://localhost:8080' // 使用网关地址
  })
}

/**
 * 获取文件列表
 * @param params 查询参数
 * @returns 文件列表
 */
export const getFileList = (params: {
  page: number
  page_size: number
  user_id?: number
  file_type?: string
}) => {
  return request.get<{
    total: number
    list: FileRecord[]
  }>('/api/v1/files', { params })
}

/**
 * 获取文件详情
 * @param id 文件ID
 * @returns 文件详情
 */
export const getFileDetail = (id: number) => {
  return request.get<FileRecord>(`/api/v1/files/${id}`)
}

/**
 * 删除文件
 * @param id 文件ID
 * @param userId 用户ID（可选）
 * @returns 操作结果
 */
export const deleteFile = (id: number, userId?: number) => {
  const params: any = {}
  if (userId) {
    params.user_id = userId
  }
  
  return request.delete(`/api/v1/files/${id}`, { params })
}

/**
 * 获取文件下载链接
 * @param id 文件ID
 * @param userId 用户ID（可选）
 * @returns 下载链接
 */
export const getFileDownloadUrl = (id: number, userId?: number) => {
  let url = `/api/v1/files/${id}/download`
  if (userId) {
    url += `?user_id=${userId}`
  }
  return url
}
