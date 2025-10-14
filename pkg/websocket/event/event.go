package event

import (
	"context"
	"sync"
)

// Event 事件接口
type Event interface {
	// Type 返回事件类型
	Type() string
	
	// Data 返回事件数据
	Data() interface{}
	
	// Context 返回事件上下文
	Context() context.Context
}

// BasicEvent 基础事件实现
type BasicEvent struct {
	EventType string
	EventData interface{}
	Ctx       context.Context
}

// NewEvent 创建新事件
func NewEvent(eventType string, data interface{}) *BasicEvent {
	return &BasicEvent{
		EventType: eventType,
		EventData: data,
		Ctx:       context.Background(),
	}
}

// NewEventWithContext 创建带上下文的新事件
func NewEventWithContext(ctx context.Context, eventType string, data interface{}) *BasicEvent {
	return &BasicEvent{
		EventType: eventType,
		EventData: data,
		Ctx:       ctx,
	}
}

// Type 返回事件类型
func (e *BasicEvent) Type() string {
	return e.EventType
}

// Data 返回事件数据
func (e *BasicEvent) Data() interface{} {
	return e.EventData
}

// Context 返回事件上下文
func (e *BasicEvent) Context() context.Context {
	return e.Ctx
}

// WithContext 设置事件上下文
func (e *BasicEvent) WithContext(ctx context.Context) *BasicEvent {
	e.Ctx = ctx
	return e
}

// Handler 事件处理器
type Handler func(Event) error

// Middleware 事件中间件
type Middleware func(Handler) Handler

// EventBus 事件总线
type EventBus struct {
	handlers   map[string][]Handler
	middleware []Middleware
	mu         sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers:   make(map[string][]Handler),
		middleware: []Middleware{},
	}
}

// On 注册事件处理器
func (b *EventBus) On(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, ok := b.handlers[eventType]; !ok {
		b.handlers[eventType] = []Handler{}
	}
	
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Off 移除事件处理器
func (b *EventBus) Off(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if handlers, ok := b.handlers[eventType]; ok {
		for i, h := range handlers {
			if &h == &handler {
				b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Use 添加中间件
func (b *EventBus) Use(middleware ...Middleware) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.middleware = append(b.middleware, middleware...)
}

// Emit 触发事件
func (b *EventBus) Emit(event Event) []error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.Type()]
	middleware := b.middleware
	b.mu.RUnlock()
	
	if !ok {
		return nil
	}
	
	var errors []error
	
	// 应用中间件
	for _, h := range handlers {
		handler := h
		
		// 应用所有中间件
		for i := len(middleware) - 1; i >= 0; i-- {
			handler = middleware[i](handler)
		}
		
		// 执行处理器
		if err := handler(event); err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}

// HasHandlers 检查是否有事件处理器
func (b *EventBus) HasHandlers(eventType string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	handlers, ok := b.handlers[eventType]
	return ok && len(handlers) > 0
}

// Clear 清除所有事件处理器
func (b *EventBus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.handlers = make(map[string][]Handler)
}

// ClearType 清除特定类型的事件处理器
func (b *EventBus) ClearType(eventType string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	delete(b.handlers, eventType)
}

// EventTypes 获取所有已注册的事件类型
func (b *EventBus) EventTypes() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	types := make([]string, 0, len(b.handlers))
	for t := range b.handlers {
		types = append(types, t)
	}
	
	return types
}
