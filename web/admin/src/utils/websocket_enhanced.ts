import { ref } from 'vue'

// 消息类型定义
export enum MessageType {
  SYSTEM = 'system',
  PING = 'ping',
  PONG = 'pong',
  WELCOME = 'welcome',
  ERROR = 'error',
  CHAT = 'chat',
  NOTIFICATION = 'notification',
  BROADCAST = 'broadcast',
  JOIN_ROOM = 'join_room',
  LEAVE_ROOM = 'leave_room',
  ROOM_MESSAGE = 'room_message',
  USER_ONLINE = 'user_online',
  USER_OFFLINE = 'user_offline',
  USER_STATUS = 'user_status',
  DATA_UPDATE = 'data_update',
  REFRESH = 'refresh'
}

// 消息处理器类型
export type MessageHandler = (message: any) => void

// WebSocket 选项
export interface WebSocketOptions {
  url?: string
  reconnect?: boolean
  reconnectInterval?: number
  maxReconnectAttempts?: number
  pingInterval?: number
  debug?: boolean
}

// 会话数据
export interface SessionData {
  [key: string]: any
}

// 分组数据
export interface GroupData {
  id: string
  name?: string
  metadata?: Record<string, any>
}

/**
 * 增强版 WebSocket 客户端
 */
export class WebSocketClient {
  private ws: WebSocket | null = null
  private url: string
  private reconnect: boolean
  private reconnectInterval: number
  private maxReconnectAttempts: number
  private pingInterval: number
  private debug: boolean
  private reconnectAttempts: number = 0
  private pingTimer: number | null = null
  private messageHandlers: Map<string, MessageHandler[]> = new Map()
  private isConnecting: boolean = false
  private token: string = ''
  private clientId: string = ''
  private sessionData = ref<SessionData>({})
  private groups = ref<GroupData[]>([])
  
  // 连接状态
  public isConnected = ref<boolean>(false)

  /**
   * 构造函数
   */
  constructor(options: WebSocketOptions = {}) {
    this.url = options.url || 'ws://localhost:8087/ws'
    this.reconnect = options.reconnect !== undefined ? options.reconnect : true
    this.reconnectInterval = options.reconnectInterval || 3000
    this.maxReconnectAttempts = options.maxReconnectAttempts || 10
    this.pingInterval = options.pingInterval || 30000
    this.debug = options.debug || false
  }

  /**
   * 连接 WebSocket
   */
  connect(token?: string): void {
    if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
      this.log('WebSocket 已连接或正在连接')
      return
    }

    if (this.isConnecting) {
      this.log('WebSocket 正在连接中...')
      return
    }

    this.isConnecting = true
    
    if (token) {
      this.token = token
    }

    const wsUrl = this.token ? `${this.url}?token=${this.token}` : this.url
    this.log(`连接 WebSocket: ${wsUrl}`)

    try {
      this.ws = new WebSocket(wsUrl)
      this.setupEventListeners()
    } catch (error) {
      this.log(`WebSocket 连接错误: ${error}`)
      this.handleReconnect()
    }
  }

  /**
   * 设置 WebSocket 事件监听器
   */
  private setupEventListeners(): void {
    if (!this.ws) return

    this.ws.onopen = () => {
      this.log('WebSocket 连接已建立')
      this.isConnected.value = true
      this.isConnecting = false
      this.reconnectAttempts = 0
      this.startPingInterval()
      
      // 触发所有 open 处理器
      this.triggerOpenHandlers()
    }

    this.ws.onclose = (event) => {
      this.log(`WebSocket 连接已关闭: ${event.code} ${event.reason}`)
      this.isConnected.value = false
      this.isConnecting = false
      this.stopPingInterval()
      
      // 触发所有 close 处理器
      this.triggerCloseHandlers(event)
      
      // 尝试重新连接
      this.handleReconnect()
    }

    this.ws.onerror = (error) => {
      this.log(`WebSocket 错误: ${error}`)
      
      // 触发所有 error 处理器
      this.triggerErrorHandlers(error)
    }

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        this.log(`收到消息: ${message.type}`, message)
        
        // 处理欢迎消息，获取客户端ID
        if (message.type === MessageType.WELCOME && message.data && message.data.client_id) {
          this.clientId = message.data.client_id
          this.log(`客户端ID: ${this.clientId}`)
          
          // 加载会话数据
          this.loadSessionData()
        }
        
        // 处理 PONG 消息
        if (message.type === MessageType.PONG) {
          this.log('收到 PONG 响应')
          return
        }
        
        // 触发消息处理器
        this.triggerMessageHandlers(message)
      } catch (error) {
        this.log(`消息解析错误: ${error}`, event.data)
      }
    }
  }

  /**
   * 处理重新连接
   */
  private handleReconnect(): void {
    if (!this.reconnect || this.reconnectAttempts >= this.maxReconnectAttempts) {
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        this.log(`达到最大重连次数 (${this.maxReconnectAttempts})`)
      }
      return
    }

    this.reconnectAttempts++
    this.log(`尝试重新连接 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)

    setTimeout(() => {
      this.connect()
    }, this.reconnectInterval)
  }

  /**
   * 开始 PING 定时器
   */
  private startPingInterval(): void {
    this.stopPingInterval()
    
    this.pingTimer = window.setInterval(() => {
      this.sendPing()
    }, this.pingInterval)
  }

  /**
   * 停止 PING 定时器
   */
  private stopPingInterval(): void {
    if (this.pingTimer !== null) {
      clearInterval(this.pingTimer)
      this.pingTimer = null
    }
  }

  /**
   * 发送 PING 消息
   */
  private sendPing(): void {
    this.send({
      type: MessageType.PING,
      timestamp: Date.now()
    })
  }

  /**
   * 发送消息
   */
  send(data: any): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.log('WebSocket 未连接，无法发送消息')
      return false
    }

    try {
      const message = typeof data === 'string' ? data : JSON.stringify(data)
      this.ws.send(message)
      return true
    } catch (error) {
      this.log(`发送消息错误: ${error}`)
      return false
    }
  }

  /**
   * 关闭连接
   */
  close(): void {
    this.stopPingInterval()
    
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    
    this.isConnected.value = false
  }

  /**
   * 注册消息处理器
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
   * 触发消息处理器
   */
  private triggerMessageHandlers(message: any): void {
    const type = message.type
    
    // 调用特定类型的处理器
    const handlers = this.messageHandlers.get(type)
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(message)
        } catch (error) {
          this.log(`处理器错误 (${type}): ${error}`)
        }
      })
    }
    
    // 调用通配符处理器
    const wildcardHandlers = this.messageHandlers.get('*')
    if (wildcardHandlers) {
      wildcardHandlers.forEach(handler => {
        try {
          handler(message)
        } catch (error) {
          this.log(`通配符处理器错误: ${error}`)
        }
      })
    }
  }
  
  /**
   * 注册连接打开处理器
   */
  onOpen(handler: () => void): void {
    this.on('__open__', handler)
  }
  
  /**
   * 触发连接打开处理器
   */
  private triggerOpenHandlers(): void {
    const handlers = this.messageHandlers.get('__open__')
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler({})
        } catch (error) {
          this.log(`连接打开处理器错误: ${error}`)
        }
      })
    }
  }
  
  /**
   * 注册连接关闭处理器
   */
  onClose(handler: (event?: CloseEvent) => void): void {
    this.on('__close__', handler)
  }
  
  /**
   * 触发连接关闭处理器
   */
  private triggerCloseHandlers(event?: CloseEvent): void {
    const handlers = this.messageHandlers.get('__close__')
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(event)
        } catch (error) {
          this.log(`连接关闭处理器错误: ${error}`)
        }
      })
    }
  }
  
  /**
   * 注册错误处理器
   */
  onError(handler: (error?: Event) => void): void {
    this.on('__error__', handler)
  }
  
  /**
   * 触发错误处理器
   */
  private triggerErrorHandlers(error?: Event): void {
    const handlers = this.messageHandlers.get('__error__')
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(error)
        } catch (err) {
          this.log(`错误处理器错误: ${err}`)
        }
      })
    }
  }
  
  /**
   * 发送通知
   */
  sendNotification(userId: string, title: string, body: string, options?: { icon?: string, link?: string }): boolean {
    return this.send({
      type: MessageType.NOTIFICATION,
      to: userId,
      content: {
        title,
        body,
        icon: options?.icon,
        link: options?.link
      },
      timestamp: Date.now()
    })
  }
  
  /**
   * 加入房间/分组
   */
  joinRoom(roomId: string): boolean {
    return this.send({
      type: MessageType.JOIN_ROOM,
      room: roomId,
      timestamp: Date.now()
    })
  }
  
  /**
   * 离开房间/分组
   */
  leaveRoom(roomId: string): boolean {
    return this.send({
      type: MessageType.LEAVE_ROOM,
      room: roomId,
      timestamp: Date.now()
    })
  }
  
  /**
   * 发送房间消息
   */
  sendRoomMessage(roomId: string, content: any): boolean {
    return this.send({
      type: MessageType.ROOM_MESSAGE,
      room: roomId,
      content,
      timestamp: Date.now()
    })
  }
  
  /**
   * 发送私聊消息
   */
  sendChatMessage(userId: string, text: string, options?: { mediaType?: string, mediaUrl?: string }): boolean {
    return this.send({
      type: MessageType.CHAT,
      to: userId,
      content: {
        text,
        mediaType: options?.mediaType,
        mediaUrl: options?.mediaUrl
      },
      timestamp: Date.now()
    })
  }
  
  /**
   * 加载会话数据
   */
  async loadSessionData(): Promise<void> {
    if (!this.clientId || !this.token) {
      this.log('无法加载会话数据：客户端ID或令牌缺失')
      return
    }
    
    try {
      const response = await fetch(`/ws/session/client/${this.clientId}`, {
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      })
      
      if (response.ok) {
        const result = await response.json()
        if (result.data) {
          this.sessionData.value = result.data
          this.log('会话数据加载成功', this.sessionData.value)
        }
      } else {
        this.log(`加载会话数据失败: ${response.status} ${response.statusText}`)
      }
    } catch (error) {
      this.log(`加载会话数据错误: ${error}`)
    }
  }
  
  /**
   * 保存会话数据
   */
  async saveSessionData(data: SessionData): Promise<boolean> {
    if (!this.clientId || !this.token) {
      this.log('无法保存会话数据：客户端ID或令牌缺失')
      return false
    }
    
    try {
      const response = await fetch(`/ws/session/client/${this.clientId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`
        },
        body: JSON.stringify(data)
      })
      
      if (response.ok) {
        // 更新本地会话数据
        this.sessionData.value = { ...this.sessionData.value, ...data }
        this.log('会话数据保存成功')
        return true
      } else {
        this.log(`保存会话数据失败: ${response.status} ${response.statusText}`)
        return false
      }
    } catch (error) {
      this.log(`保存会话数据错误: ${error}`)
      return false
    }
  }
  
  /**
   * 删除会话数据
   */
  async deleteSessionData(key: string): Promise<boolean> {
    if (!this.clientId || !this.token) {
      this.log('无法删除会话数据：客户端ID或令牌缺失')
      return false
    }
    
    try {
      const response = await fetch(`/ws/session/client/${this.clientId}/${key}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      })
      
      if (response.ok) {
        // 更新本地会话数据
        const newSessionData = { ...this.sessionData.value }
        delete newSessionData[key]
        this.sessionData.value = newSessionData
        this.log(`会话数据删除成功: ${key}`)
        return true
      } else {
        this.log(`删除会话数据失败: ${response.status} ${response.statusText}`)
        return false
      }
    } catch (error) {
      this.log(`删除会话数据错误: ${error}`)
      return false
    }
  }
  
  /**
   * 加载用户分组
   */
  async loadGroups(): Promise<void> {
    if (!this.clientId || !this.token) {
      this.log('无法加载分组：客户端ID或令牌缺失')
      return
    }
    
    try {
      const response = await fetch(`/ws/group/client/${this.clientId}`, {
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      })
      
      if (response.ok) {
        const result = await response.json()
        if (result.groups) {
          // 转换为分组数据
          this.groups.value = result.groups.map((groupId: string) => ({ id: groupId }))
          this.log('分组加载成功', this.groups.value)
          
          // 加载每个分组的元数据
          await this.loadGroupsMetadata()
        }
      } else {
        this.log(`加载分组失败: ${response.status} ${response.statusText}`)
      }
    } catch (error) {
      this.log(`加载分组错误: ${error}`)
    }
  }
  
  /**
   * 加载分组元数据
   */
  private async loadGroupsMetadata(): Promise<void> {
    if (!this.token) return
    
    for (const group of this.groups.value) {
      try {
        const response = await fetch(`/ws/group/${group.id}/metadata`, {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        })
        
        if (response.ok) {
          const result = await response.json()
          if (result.metadata) {
            group.metadata = result.metadata
            if (result.metadata.name) {
              group.name = result.metadata.name
            }
          }
        }
      } catch (error) {
        this.log(`加载分组元数据错误: ${error}`)
      }
    }
  }
  
  /**
   * 获取会话数据
   */
  getSessionData(): SessionData {
    return this.sessionData.value
  }
  
  /**
   * 获取分组
   */
  getGroups(): GroupData[] {
    return this.groups.value
  }
  
  /**
   * 日志输出
   */
  private log(message: string, ...args: any[]): void {
    if (this.debug) {
      console.log(`[WebSocket] ${message}`, ...args)
    }
  }
}

// 创建全局 WebSocket 客户端实例
export const createWebSocketClient = (options?: WebSocketOptions) => {
  return new WebSocketClient(options)
}

// 默认导出
export default createWebSocketClient
