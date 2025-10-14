import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Notification {
  id: string
  title: string
  body: string
  icon?: string
  link?: string
  category?: string
  read: boolean
  timestamp: number
  data?: Record<string, any>
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  
  // 计算属性
  const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)
  
  const notificationsByCategory = computed(() => {
    const result: Record<string, Notification[]> = {}
    
    notifications.value.forEach(notification => {
      const category = notification.category || 'other'
      if (!result[category]) {
        result[category] = []
      }
      result[category].push(notification)
    })
    
    return result
  })
  
  // 加载通知
  const loadNotifications = async () => {
    try {
      isLoading.value = true
      error.value = null
      
      // 从本地存储加载通知
      const storedNotifications = localStorage.getItem('notifications')
      if (storedNotifications) {
        notifications.value = JSON.parse(storedNotifications)
      }
    } catch (err) {
      error.value = '加载通知失败'
      console.error('加载通知失败:', err)
    } finally {
      isLoading.value = false
    }
  }
  
  // 添加通知
  const addNotification = (notification: Notification) => {
    // 检查是否已存在相同ID的通知
    const existingIndex = notifications.value.findIndex(n => n.id === notification.id)
    
    if (existingIndex >= 0) {
      // 更新现有通知
      notifications.value[existingIndex] = notification
    } else {
      // 添加新通知
      notifications.value.unshift(notification)
    }
    
    // 限制最多保存100条通知
    if (notifications.value.length > 100) {
      notifications.value = notifications.value.slice(0, 100)
    }
    
    // 保存到本地存储
    saveToLocalStorage()
  }
  
  // 标记通知为已读
  const markAsRead = (id: string) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.read = true
      saveToLocalStorage()
    }
  }
  
  // 标记所有通知为已读
  const markAllAsRead = () => {
    notifications.value.forEach(notification => {
      notification.read = true
    })
    saveToLocalStorage()
  }
  
  // 删除通知
  const removeNotification = (id: string) => {
    notifications.value = notifications.value.filter(n => n.id !== id)
    saveToLocalStorage()
  }
  
  // 清空所有通知
  const clearAll = () => {
    notifications.value = []
    saveToLocalStorage()
  }
  
  // 保存到本地存储
  const saveToLocalStorage = () => {
    localStorage.setItem('notifications', JSON.stringify(notifications.value))
  }
  
  return {
    notifications,
    unreadCount,
    notificationsByCategory,
    isLoading,
    error,
    loadNotifications,
    addNotification,
    markAsRead,
    markAllAsRead,
    removeNotification,
    clearAll
  }
})
