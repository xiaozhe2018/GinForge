import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createWebSocketClient } from '@/utils/websocket_enhanced'

export const useWebSocketStore = defineStore('websocket', () => {
  // WebSocket 客户端实例
  const wsClient = createWebSocketClient({
    url: 'ws://localhost:8087/ws',
    reconnect: true,
    debug: true
  })
  
  // 连接状态
  const isConnected = ref(false)
  const isConnecting = ref(false)
  
  // 监听连接状态
  wsClient.onOpen(() => {
    isConnected.value = true
    isConnecting.value = false
    console.log('WebSocket 已连接')
  })
  
  wsClient.onClose(() => {
    isConnected.value = false
    isConnecting.value = false
    console.log('WebSocket 已断开')
  })
  
  wsClient.onError((error) => {
    isConnecting.value = false
    console.error('WebSocket 错误:', error)
  })
  
  // 连接方法
  const connect = (token?: string) => {
    if (isConnected.value || isConnecting.value) return
    
    isConnecting.value = true
    wsClient.connect(token)
  }
  
  // 关闭连接
  const disconnect = () => {
    wsClient.close()
    isConnected.value = false
  }
  
  // 注册消息处理器
  const on = (type: string, handler: (message: any) => void) => {
    wsClient.on(type, handler)
  }
  
  // 移除消息处理器
  const off = (type: string, handler?: (message: any) => void) => {
    wsClient.off(type, handler)
  }
  
  // 发送消息
  const send = (data: any) => {
    return wsClient.send(data)
  }
  
  // 发送通知
  const sendNotification = (userId: string, title: string, body: string, options?: { icon?: string, link?: string }) => {
    return wsClient.sendNotification(userId, title, body, options)
  }
  
  // 加入房间
  const joinRoom = (roomId: string) => {
    return wsClient.joinRoom(roomId)
  }
  
  // 离开房间
  const leaveRoom = (roomId: string) => {
    return wsClient.leaveRoom(roomId)
  }
  
  // 发送房间消息
  const sendRoomMessage = (roomId: string, content: any) => {
    return wsClient.sendRoomMessage(roomId, content)
  }
  
  return {
    isConnected,
    isConnecting,
    connect,
    disconnect,
    on,
    off,
    send,
    sendNotification,
    joinRoom,
    leaveRoom,
    sendRoomMessage
  }
})
