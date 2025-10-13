import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8083',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 是否正在跳转登录页（避免重复跳转）
let isRedirectingToLogin = false

// 导出重置函数（用于登录成功后重置状态）
export const resetRedirectFlag = () => {
  isRedirectingToLogin = false
}

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { data } = response
    
    // 如果返回的数据结构是 { code, message, data }
    if (data.code !== undefined) {
      if (data.code === 200 || data.code === 0) {
        return data.data || data
      } else if (data.code === 401) {
        // 处理token过期（避免重复跳转）
        if (!isRedirectingToLogin) {
          isRedirectingToLogin = true
          ElMessage.error('登录已过期，请重新登录')
          // 清除所有本地存储数据
          localStorage.removeItem('admin_token')
          localStorage.removeItem('admin_user_info')
          localStorage.removeItem('admin_permissions')
          localStorage.removeItem('admin_menus')
          // 使用Vue Router跳转到登录页
          router.push('/login').catch(() => {
            // 如果路由跳转失败，使用window.location
            window.location.href = '/login'
          })
        }
        return Promise.reject(new Error(data.message || '认证失败'))
      } else {
        ElMessage.error(data.message || '请求失败')
        return Promise.reject(new Error(data.message || '请求失败'))
      }
    }
    
    return data
  },
  (error) => {
    const { response } = error
    
    if (response) {
      const { status, data } = response
      
      switch (status) {
        case 401:
          // 处理HTTP 401状态码（避免重复跳转）
          if (!isRedirectingToLogin) {
            isRedirectingToLogin = true
            ElMessage.error('登录已过期，请重新登录')
            // 清除所有本地存储数据
            localStorage.removeItem('admin_token')
            localStorage.removeItem('admin_user_info')
            localStorage.removeItem('admin_permissions')
            localStorage.removeItem('admin_menus')
            // 使用Vue Router跳转到登录页
            router.push('/login').catch(() => {
              // 如果路由跳转失败，使用window.location
              window.location.href = '/login'
            })
          }
          break
        case 403:
          ElMessage.error('拒绝访问，权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    
    return Promise.reject(error)
  }
)

export default request