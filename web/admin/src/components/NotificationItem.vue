<template>
  <div 
    class="notification-item" 
    :class="{ 'notification-unread': !notification.read }"
    @click="handleClick"
  >
    <!-- 图标 -->
    <div class="notification-icon">
      <el-icon :size="24" :color="iconColor">
        <component :is="getIcon(notification.icon)" />
      </el-icon>
    </div>
    
    <!-- 内容 -->
    <div class="notification-content">
      <div class="notification-header">
        <h4 class="notification-title">{{ notification.title }}</h4>
        <span class="notification-time">{{ formatTime(notification.timestamp) }}</span>
      </div>
      <p class="notification-body">{{ notification.body }}</p>
      
      <!-- 操作按钮 -->
      <div class="notification-actions">
        <el-button 
          v-if="!notification.read" 
          type="primary" 
          link 
          size="small"
          @click.stop="markAsRead"
        >
          标记为已读
        </el-button>
        <el-button 
          type="danger" 
          link 
          size="small"
          @click.stop="remove"
        >
          删除
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import type { Notification } from '@/stores/notification'

// 属性
const props = defineProps<{
  notification: Notification
}>()

// 事件
const emit = defineEmits<{
  'mark-as-read': [id: string]
  'remove': [id: string]
}>()

// 路由
const router = useRouter()

// 计算属性
const iconColor = computed(() => {
  const category = props.notification.category || 'other'
  
  switch (category) {
    case 'system':
      return '#409EFF' // 蓝色
    case 'order':
      return '#67C23A' // 绿色
    case 'user':
      return '#E6A23C' // 黄色
    default:
      return '#909399' // 灰色
  }
})

// 方法
const getIcon = (iconName?: string) => {
  if (!iconName) return ElementPlusIconsVue.Bell
  
  // 查找对应的图标
  const icon = (ElementPlusIconsVue as any)[iconName]
  return icon || ElementPlusIconsVue.Bell
}

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp * 1000)
  const now = new Date()
  
  // 今天的通知显示时间
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 昨天的通知
  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 一周内的通知
  const oneWeekAgo = new Date(now)
  oneWeekAgo.setDate(now.getDate() - 7)
  if (date > oneWeekAgo) {
    const days = ['日', '一', '二', '三', '四', '五', '六']
    return '周' + days[date.getDay()] + ' ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 更早的通知
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const handleClick = () => {
  // 如果未读，标记为已读
  if (!props.notification.read) {
    markAsRead()
  }
  
  // 如果有链接，跳转
  if (props.notification.link) {
    router.push(props.notification.link)
  }
}

const markAsRead = () => {
  emit('mark-as-read', props.notification.id)
}

const remove = () => {
  emit('remove', props.notification.id)
}
</script>

<style scoped>
.notification-item {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-unread {
  background-color: #ecf5ff;
}

.notification-unread:hover {
  background-color: #e4efff;
}

.notification-icon {
  margin-right: 12px;
  display: flex;
  align-items: flex-start;
  padding-top: 2px;
}

.notification-content {
  flex: 1;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.notification-title {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}

.notification-time {
  font-size: 12px;
  color: #909399;
}

.notification-body {
  margin: 0 0 8px;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
}

.notification-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
