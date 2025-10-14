<template>
  <div class="websocket-test">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>WebSocket 测试工具</span>
          <el-tag :type="connectionStatus.type">{{ connectionStatus.text }}</el-tag>
        </div>
      </template>

      <el-space direction="vertical" :size="20" style="width: 100%">
        <!-- 连接控制 -->
        <el-row :gutter="10">
          <el-col :span="12">
            <el-button type="primary" @click="connect" :disabled="isConnected">连接</el-button>
            <el-button type="danger" @click="disconnect" :disabled="!isConnected">断开</el-button>
            <el-button @click="clearLogs">清空日志</el-button>
          </el-col>
          <el-col :span="12" style="text-align: right">
            <span>在线用户: {{ onlineUsers }}</span>
          </el-col>
        </el-row>

        <!-- 发送消息 -->
        <el-card shadow="never">
          <template #header>发送测试消息</template>
          
          <el-form :model="messageForm" label-width="100px">
            <el-form-item label="消息类型">
              <el-select v-model="messageForm.type">
                <el-option label="系统消息" value="system" />
                <el-option label="通知消息" value="notification" />
                <el-option label="聊天消息" value="chat" />
                <el-option label="心跳" value="ping" />
              </el-select>
            </el-form-item>

            <el-form-item label="标题" v-if="messageForm.type === 'notification'">
              <el-input v-model="messageForm.title" placeholder="通知标题" />
            </el-form-item>

            <el-form-item label="内容">
              <el-input
                v-model="messageForm.content"
                type="textarea"
                :rows="3"
                placeholder="消息内容"
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="sendMessage" :disabled="!isConnected">
                发送消息
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 消息日志 -->
        <el-card shadow="never">
          <template #header>
            <div style="display: flex; justify-content: space-between">
              <span>消息日志 ({{ logs.length }})</span>
              <el-switch v-model="autoScroll" active-text="自动滚动" />
            </div>
          </template>
          
          <el-scrollbar height="400px" ref="scrollbarRef">
            <div class="log-list">
              <div v-for="(log, index) in logs" :key="index" class="log-item" :class="`log-${log.type}`">
                <div class="log-time">{{ formatTime(log.timestamp) }}</div>
                <div class="log-type">[{{ log.type }}]</div>
                <div class="log-content">{{ log.message }}</div>
              </div>
              
              <div v-if="logs.length === 0" class="empty-log">
                暂无日志
              </div>
            </div>
          </el-scrollbar>
        </el-card>
      </el-space>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import WebSocketClient from '@/utils/websocket'

interface LogItem {
  type: 'send' | 'receive' | 'info' | 'error'
  message: string
  timestamp: number
}

const ws = new WebSocketClient({ debug: true })

const isConnected = ref(false)
const onlineUsers = ref(0)
const logs = ref<LogItem[]>([])
const autoScroll = ref(true)
const scrollbarRef = ref()

const messageForm = ref({
  type: 'notification',
  title: 'WebSocket 测试通知',
  content: '这是一条测试消息'
})

const connectionStatus = computed(() => {
  if (isConnected.value) {
    return { text: '已连接', type: 'success' }
  }
  return { text: '未连接', type: 'info' }
})

// 连接
const connect = () => {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    addLog('error', '未登录，无法连接')
    return
  }
  
  ws.connect(token)
  
  ws.onOpen(() => {
    isConnected.value = true
    addLog('info', 'WebSocket 连接成功')
    getOnlineUsers()
  })
  
  ws.onClose(() => {
    isConnected.value = false
    addLog('info', 'WebSocket 连接已关闭')
  })
  
  ws.onError((error) => {
    addLog('error', `WebSocket 错误: ${error}`)
  })
  
  // 监听所有消息
  ws.on('*', (message) => {
    addLog('receive', `收到消息: ${message.type} - ${JSON.stringify(message.content)}`)
  })
}

// 断开
const disconnect = () => {
  ws.close()
  isConnected.value = false
}

// 发送消息
const sendMessage = () => {
  let content: any = messageForm.value.content
  
  if (messageForm.value.type === 'notification') {
    content = {
      title: messageForm.value.title,
      body: messageForm.value.content,
      icon: 'Bell',
      category: 'info'
    }
  }
  
  ws.send({
    type: messageForm.value.type,
    content
  })
  
  addLog('send', `发送 ${messageForm.value.type} 消息`)
}

// 获取在线用户数
const getOnlineUsers = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch('http://localhost:8083/api/v1/admin/ws/online-users', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const result = await response.json()
    if (result.code === 0) {
      onlineUsers.value = result.data.total
    }
  } catch (error) {
    console.error('获取在线用户失败:', error)
  }
}

// 添加日志
const addLog = (type: LogItem['type'], message: string) => {
  logs.value.push({
    type,
    message,
    timestamp: Date.now()
  })
  
  if (autoScroll.value) {
    nextTick(() => {
      scrollbarRef.value?.setScrollTop(9999)
    })
  }
}

// 清空日志
const clearLogs = () => {
  logs.value = []
}

// 格式化时间
const formatTime = (timestamp: number): string => {
  const date = new Date(timestamp)
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
}

onMounted(() => {
  addLog('info', 'WebSocket 测试工具已就绪')
})

onBeforeUnmount(() => {
  ws.close()
})
</script>

<style scoped>
.websocket-test {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-list {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.log-item {
  display: flex;
  gap: 10px;
  padding: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.log-item:last-child {
  border-bottom: none;
}

.log-time {
  color: #909399;
  flex-shrink: 0;
}

.log-type {
  font-weight: bold;
  flex-shrink: 0;
  width: 80px;
}

.log-content {
  flex: 1;
  word-break: break-all;
}

.log-send .log-type {
  color: #409eff;
}

.log-receive .log-type {
  color: #67c23a;
}

.log-info .log-type {
  color: #909399;
}

.log-error .log-type {
  color: #f56c6c;
}

.empty-log {
  text-align: center;
  padding: 40px;
  color: #909399;
}
</style>

