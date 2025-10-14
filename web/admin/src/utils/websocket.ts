/**
 * WebSocket 客户端工具类
 * 提供 WebSocket 连接管理、消息收发、自动重连等功能
 */

export interface WebSocketMessage {
  type: string
  id?: string
  from?: string
  from_name?: string
  to?: string
  room?: string
  content: any
  data?: Record<string, any>
  timestamp: number
  created_at?: string
}

export interface WebSocketOptions {
  url?: string
  reconnect?: boolean
  reconnectInterval?: number
  maxReconnectAttempts?: number
  heartbeatInterval?: number
  debug?: boolean
}

export type MessageHandler = (message: WebSocketMessage) => void

class WebSocketClient {
  private ws: WebSocket | null = null
  private url = 'ws://localhost:8087/ws' // WebSocket Gateway 地址
  private token: string = ''
  private reconnect: boolean
  private reconnectInterval: number
  private maxReconnectAttempts: number
  private heartbeatInterval: number
  private debug: boolean

  private reconnectAttempts: number = 0
  private reconnectTimer: number | null = null
  private heartbeatTimer: number | null = null
  private isManualClose: boolean = false

  private messageHandlers: Map<string, MessageHandler[]> = new Map()
  private onOpenCallback: (() => void) | null = null
  private onCloseCallback: (() => void) | null = null
  private onErrorCallback: ((error: Event) => void) | null = null

  constructor(options: WebSocketOptions = {}) {
    // 默认使用 WebSocket Gateway (8087)
    this.url = options.url || 'ws://localhost:8087/ws'
    this.reconnect = options.reconnect !== false
    this.reconnectInterval = options.reconnectInterval || 3000
    this.maxReconnectAttempts = options.maxReconnectAttempts || 10
    this.heartbeatInterval = options.heartbeatInterval || 30000
    this.debug = options.debug || false
  }

  /**
   * 连接 WebSocket
   */
  connect(token: string): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.log('WebSocket 已连接')
      return
    }

    this.token = token
    this.isManualClose = false
    
    const wsUrl = `${this.url}?token=${token}`
    this.log('正在连接 WebSocket:', wsUrl)

    try {
      this.ws = new WebSocket(wsUrl)
      this.setupEventHandlers()
    } catch (error) {
      this.log('WebSocket 连接失败:', error)
      this.scheduleReconnect()
    }
  }

  /**
   * 设置事件处理器
   */
  private setupEventHandlers(): void {
    if (!this.ws) return

    this.ws.onopen = () => {
      this.log('✅ WebSocket 连接成功')
      this.reconnectAttempts = 0
      this.startHeartbeat()
      this.onOpenCallback?.()
    }

    this.ws.onmessage = (event: MessageEvent) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.log('📨 收到消息:', message)
        this.handleMessage(message)
      } catch (error) {
        this.log('消息解析失败:', error, event.data)
      }
    }

    this.ws.onerror = (event: Event) => {
      this.log('❌ WebSocket 错误:', event)
      this.onErrorCallback?.(event)
    }

    this.ws.onclose = (event: CloseEvent) => {
      this.log('🔌 WebSocket 连接关闭:', event.code, event.reason)
      this.stopHeartbeat()
      this.onCloseCallback?.()

      if (!this.isManualClose && this.reconnect) {
        this.scheduleReconnect()
      }
    }
  }

  /**
   * 处理收到的消息
   */
  private handleMessage(message: WebSocketMessage): void {
    // 触发全局处理器
    const globalHandlers = this.messageHandlers.get('*') || []
    globalHandlers.forEach(handler => handler(message))

    // 触发特定类型的处理器
    const typeHandlers = this.messageHandlers.get(message.type) || []
    typeHandlers.forEach(handler => handler(message))
  }

  /**
   * 发送消息
   */
  send(message: Partial<WebSocketMessage>): void {
    if (this.ws?.readyState !== WebSocket.OPEN) {
      this.log('⚠️ WebSocket 未连接，无法发送消息')
      return
    }

    const fullMessage: WebSocketMessage = {
      type: message.type || 'chat',
      content: message.content,
      to: message.to,
      room: message.room,
      data: message.data,
      timestamp: Date.now(),
      ...message
    }

    this.ws.send(JSON.stringify(fullMessage))
    this.log('📤 发送消息:', fullMessage)
  }

  /**
   * 发送心跳
   */
  private sendHeartbeat(): void {
    this.send({ type: 'ping' })
  }

  /**
   * 启动心跳
   */
  private startHeartbeat(): void {
    this.stopHeartbeat()
    this.heartbeatTimer = window.setInterval(() => {
      this.sendHeartbeat()
    }, this.heartbeatInterval)
  }

  /**
   * 停止心跳
   */
  private stopHeartbeat(): void {
    if (this.heartbeatTimer !== null) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  /**
   * 调度重连
   */
  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.log(`❌ 重连失败，已达到最大重连次数 (${this.maxReconnectAttempts})`)
      return
    }

    this.reconnectAttempts++
    this.log(`🔄 将在 ${this.reconnectInterval}ms 后尝试第 ${this.reconnectAttempts} 次重连...`)

    this.reconnectTimer = window.setTimeout(() => {
      this.log(`正在尝试第 ${this.reconnectAttempts} 次重连...`)
      this.connect(this.token)
    }, this.reconnectInterval)
  }

  /**
   * 取消重连
   */
  private cancelReconnect(): void {
    if (this.reconnectTimer !== null) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }

  /**
   * 关闭连接
   */
  close(): void {
    this.isManualClose = true
    this.cancelReconnect()
    this.stopHeartbeat()
    
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    
    this.log('🔌 WebSocket 已主动关闭')
  }

  /**
   * 注册消息处理器
   * @param type 消息类型，'*' 表示所有消息
   * @param handler 处理函数
   */
  on(type: string, handler: MessageHandler): void {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, [])
    }
    
    const handlers = this.messageHandlers.get(type)!
    
    // 防止重复注册同一个处理器
    if (handlers.includes(handler)) {
      this.log(`⚠️ 处理器已存在，跳过注册: ${type}`)
      return
    }
    
    handlers.push(handler)
    this.log(`✅ 注册消息处理器: ${type} (当前数量: ${handlers.length})`)
  }

  /**
   * 移除消息处理器
   */
  off(type: string, handler?: MessageHandler): void {
    if (!handler) {
      // 移除该类型的所有处理器
      this.messageHandlers.delete(type)
      this.log(`🧹 清理所有 ${type} 处理器`)
      return
    }

    const handlers = this.messageHandlers.get(type)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
        this.log(`🧹 清理 ${type} 处理器 (剩余: ${handlers.length})`)
      } else {
        this.log(`⚠️ 未找到要清理的 ${type} 处理器`)
      }
    }
  }
  
  /**
   * 重置所有处理器
   */
  resetAllHandlers(): void {
    this.messageHandlers.clear()
    this.log(`🧹 已重置所有消息处理器`)
  }

  /**
   * 设置连接打开回调
   */
  onOpen(callback: () => void): void {
    this.onOpenCallback = callback
  }

  /**
   * 设置连接关闭回调
   */
  onClose(callback: () => void): void {
    this.onCloseCallback = callback
  }

  /**
   * 设置错误回调
   */
  onError(callback: (error: Event) => void): void {
    this.onErrorCallback = callback
  }

  /**
   * 获取连接状态
   */
  getState(): number {
    return this.ws?.readyState ?? WebSocket.CLOSED
  }

  /**
   * 是否已连接
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  /**
   * 日志输出
   */
  private log(...args: any[]): void {
    if (this.debug) {
      console.log('[WebSocket]', ...args)
    }
  }

  /**
   * 加入房间
   */
  joinRoom(room: string): void {
    this.send({
      type: 'join_room',
      content: room
    })
  }

  /**
   * 离开房间
   */
  leaveRoom(room: string): void {
    this.send({
      type: 'leave_room',
      content: room
    })
  }

  /**
   * 发送聊天消息
   */
  sendChat(to: string, text: string): void {
    this.send({
      type: 'chat',
      to,
      content: {
        text,
        media_type: 'text'
      }
    })
  }

  /**
   * 发送房间消息
   */
  sendToRoom(room: string, content: any): void {
    this.send({
      type: 'room_message',
      room,
      content
    })
  }
}

// 导出单例实例（可选）
export const wsClient = new WebSocketClient({
  debug: import.meta.env.DEV
})

export default WebSocketClient

