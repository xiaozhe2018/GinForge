<template>
  <div class="notification-center">
    <!-- 通知图标和徽章 -->
    <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
      <el-button
        class="notification-button"
        :icon="Bell"
        circle
        @click="toggleNotificationPanel"
      />
    </el-badge>

    <!-- 通知面板 -->
    <el-drawer
      v-model="showNotifications"
      title="通知中心"
      direction="rtl"
      size="350px"
      :with-header="true"
    >
      <template #header>
        <div class="notification-drawer-header">
          <span>通知中心</span>
          <div class="notification-actions">
            <el-button
              v-if="unreadCount > 0"
              type="primary"
              link
              @click="markAllAsRead"
            >
              全部标记为已读
            </el-button>
            <el-button
              type="danger"
              link
              @click="clearAll"
            >
              清空
            </el-button>
          </div>
        </div>
      </template>

      <!-- 通知列表 -->
      <div v-if="notifications.length > 0" class="notification-list">
        <el-tabs v-model="activeTab" class="notification-tabs">
          <el-tab-pane label="全部" name="all">
            <notification-item
              v-for="notification in notifications"
              :key="notification.id"
              :notification="notification"
              @mark-as-read="markAsRead"
              @remove="removeNotification"
            />
          </el-tab-pane>

          <el-tab-pane label="系统" name="system">
            <notification-item
              v-for="notification in systemNotifications"
              :key="notification.id"
              :notification="notification"
              @mark-as-read="markAsRead"
              @remove="removeNotification"
            />
          </el-tab-pane>

          <el-tab-pane label="订单" name="order">
            <notification-item
              v-for="notification in orderNotifications"
              :key="notification.id"
              :notification="notification"
              @mark-as-read="markAsRead"
              @remove="removeNotification"
            />
          </el-tab-pane>

          <el-tab-pane label="用户" name="user">
            <notification-item
              v-for="notification in userNotifications"
              :key="notification.id"
              :notification="notification"
              @mark-as-read="markAsRead"
              @remove="removeNotification"
            />
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 空状态 -->
      <div v-else class="notification-empty">
        <el-empty description="暂无通知" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { Bell } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import { useNotificationStore, type Notification } from '@/stores/notification'
import { useWebSocketStore } from '@/stores/websocket'
import NotificationItem from './NotificationItem.vue'
import { markRaw } from 'vue'

// 状态
const showNotifications = ref(false)
const activeTab = ref('all')

// Store
const notificationStore = useNotificationStore()
const wsStore = useWebSocketStore()

// 计算属性
const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const systemNotifications = computed(() => 
  notifications.value.filter(n => n.category === 'system')
)

const orderNotifications = computed(() => 
  notifications.value.filter(n => n.category === 'order')
)

const userNotifications = computed(() => 
  notifications.value.filter(n => n.category === 'user')
)

// 方法
const toggleNotificationPanel = () => {
  showNotifications.value = !showNotifications.value
}

const markAsRead = (id: string) => {
  notificationStore.markAsRead(id)
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const removeNotification = (id: string) => {
  notificationStore.removeNotification(id)
}

const clearAll = () => {
  notificationStore.clearAll()
}

// 处理通知消息
const handleNotification = (message: any) => {
  console.log('收到通知消息:', message)
  const content = message.content
  
  // 添加到通知中心
  notificationStore.addNotification({
    id: message.data?.notification_id || String(Date.now()),
    title: content.title,
    body: content.body,
    icon: content.icon,
    link: content.link,
    category: content.category || 'other',
    read: false,
    timestamp: message.timestamp,
    data: content.data
  })
  
  // 显示弹窗通知
  ElNotification({
    title: content.title,
    message: content.body,
    type: 'info',
    duration: 4500,
    onClick: () => {
      if (content.link) {
        window.location.href = content.link
      }
    }
  })
}

// 生命周期钩子
onMounted(() => {
  // 加载通知
  notificationStore.loadNotifications()
  
  // 注册 WebSocket 通知处理器
  wsStore.on('notification', handleNotification)
})

onBeforeUnmount(() => {
  // 移除 WebSocket 通知处理器
  wsStore.off('notification', handleNotification)
})
</script>

<style scoped>
.notification-center {
  display: inline-block;
}

.notification-badge {
  margin-right: 16px;
}

.notification-button {
  font-size: 18px;
}

.notification-drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.notification-actions {
  display: flex;
  gap: 8px;
}

.notification-list {
  height: 100%;
  overflow-y: auto;
}

.notification-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-tabs {
  height: 100%;
}

:deep(.el-tabs__content) {
  padding: 0;
  height: calc(100vh - 120px);
  overflow-y: auto;
}
</style>