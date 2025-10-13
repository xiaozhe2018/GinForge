import api from './index'

// 订单相关API
export interface Order {
  id: string
  order_no: string
  user_id: string
  merchant_id: string
  total_amount: number
  status: number
  payment_status: number
  shipping_status: number
  items: OrderItem[]
  created_at: string
  updated_at: string
}

export interface OrderItem {
  id: string
  product_id: string
  product_name: string
  price: number
  quantity: number
  total: number
}

export interface OrderListParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
  payment_status?: number
  shipping_status?: number
  start_date?: string
  end_date?: string
}

export interface OrderListResponse {
  list: Order[]
  total: number
  page: number
  page_size: number
}

// 获取订单列表
export const getOrderList = (params: OrderListParams = {}) => {
  return api.get<OrderListResponse>('/v1/admin/order/list', { params })
}

// 获取订单详情
export const getOrderDetail = (id: string) => {
  return api.get<Order>(`/v1/admin/order/${id}`)
}

// 更新订单状态
export const updateOrderStatus = (id: string, status: number) => {
  return api.put(`/v1/admin/order/${id}/status`, { status })
}

// 更新支付状态
export const updatePaymentStatus = (id: string, payment_status: number) => {
  return api.put(`/v1/admin/order/${id}/payment-status`, { payment_status })
}

// 更新发货状态
export const updateShippingStatus = (id: string, shipping_status: number) => {
  return api.put(`/v1/admin/order/${id}/shipping-status`, { shipping_status })
}

// 取消订单
export const cancelOrder = (id: string, reason: string) => {
  return api.put(`/v1/admin/order/${id}/cancel`, { reason })
}
