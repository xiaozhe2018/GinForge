import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8083',
  timeout: 10000,
  withCredentials: true, // 重要：允许携带 Cookie（用于 CSRF Token）
  headers: {
    'Content-Type': 'application/json',
  },
})

// 是否正在跳转登录页（避免重复跳转）
let isRedirectingToLogin = false

// CSRF Token 缓存
let csrfToken: string | null = null

// 导出重置函数（用于登录成功后重置状态）
export const resetRedirectFlag = () => {
  isRedirectingToLogin = false
}

// 获取 CSRF Token（从 Cookie 或 API）
const getCSRFToken = async (): Promise<string | null> => {
  // 如果已有缓存的 Token，直接返回
  if (csrfToken) {
    return csrfToken
  }

  // 尝试从 Cookie 读取（如果浏览器允许）
  // 注意：如果后端设置了 HttpOnly，前端无法读取，需要通过 API 获取
  try {
    // 尝试从 Cookie 读取
    const cookies = document.cookie.split(';')
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split('=')
      if (name === 'csrf_token') {
        csrfToken = decodeURIComponent(value)
        return csrfToken
      }
    }
  } catch (e) {
    // Cookie 读取失败，忽略
  }

  // 如果 Cookie 中没有，尝试从 API 获取
  try {
    const response = await axios.get('/api/v1/admin/csrf-token', {
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8083',
      withCredentials: true, // 重要：允许携带 Cookie
    })
    if (response.data?.data?.csrf_token) {
      csrfToken = response.data.data.csrf_token
      return csrfToken
    }
  } catch (e) {
    // 如果获取失败，可能是未登录或网络问题，不影响 GET 请求
    console.warn('获取 CSRF Token 失败（这通常不影响 GET 请求）:', e)
  }

  return null
}

// 请求拦截器
request.interceptors.request.use(
  async (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 对于需要认证的修改操作（POST, PUT, DELETE, PATCH），添加 CSRF Token
    const method = config.method?.toUpperCase()
    if (method && ['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
      // 检查是否是认证路由（需要 CSRF Token）
      const url = config.url || ''
      if (url.startsWith('/api/v1/admin') && !url.includes('/login') && !url.includes('/csrf-token')) {
        const token = await getCSRFToken()
        if (token) {
          config.headers['X-CSRF-Token'] = token
        }
      }
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
          // CSRF Token 验证失败，尝试重新获取
          if (data?.message?.includes('CSRF')) {
            csrfToken = null // 清除缓存的 Token
            ElMessage.error('CSRF 令牌验证失败，请重试')
            // 可以在这里自动重试请求，但为了避免无限循环，暂时只提示
          } else {
            ElMessage.error('拒绝访问，权限不足')
          }
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