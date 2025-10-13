package circuit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"goweb/pkg/logger"
)

// State 熔断器状态
type State int

const (
	StateClosed   State = iota // 关闭状态（正常）
	StateOpen                  // 开启状态（熔断）
	StateHalfOpen              // 半开状态（尝试恢复）
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker 熔断器
type Breaker struct {
	name          string
	maxRequests   uint32        // 半开状态下的最大请求数
	interval      time.Duration // 统计时间窗口
	timeout       time.Duration // 熔断超时时间
	readyToTrip   func(counts Counts) bool
	onStateChange func(name string, from State, to State)

	mu         sync.Mutex
	state      State
	generation uint64
	counts     Counts
	expiry     time.Time
}

// Counts 计数器
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// Request 请求结果
type Request struct {
	Success bool
	Error   error
}

// Config 熔断器配置
type Config struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts Counts) bool
	OnStateChange func(name string, from State, to State)
}

// DefaultConfig 默认配置
func DefaultConfig(name string) *Config {
	return &Config{
		Name:        name,
		MaxRequests: 3,
		Interval:    time.Minute,
		Timeout:     time.Minute,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from State, to State) {
			// 默认不处理状态变化
		},
	}
}

// NewBreaker 创建熔断器
func NewBreaker(cfg *Config, log logger.Logger) *Breaker {
	if cfg.ReadyToTrip == nil {
		cfg.ReadyToTrip = func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 5
		}
	}

	cb := &Breaker{
		name:          cfg.Name,
		maxRequests:   cfg.MaxRequests,
		interval:      cfg.Interval,
		timeout:       cfg.Timeout,
		readyToTrip:   cfg.ReadyToTrip,
		onStateChange: cfg.OnStateChange,
		state:         StateClosed,
		expiry:        time.Now().Add(cfg.Interval),
	}

	return cb
}

// Execute 执行函数
func (cb *Breaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := req()
	cb.afterRequest(generation, err == nil)
	return result, err
}

// ExecuteWithContext 带上下文的执行函数
func (cb *Breaker) ExecuteWithContext(ctx context.Context, req func(context.Context) (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := req(ctx)
	cb.afterRequest(generation, err == nil)
	return result, err
}

// beforeRequest 请求前检查
func (cb *Breaker) beforeRequest() (uint64, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	generation, state := cb.currentState(now)

	switch state {
	case StateOpen:
		return generation, fmt.Errorf("circuit breaker is open")
	case StateHalfOpen:
		if cb.counts.Requests >= cb.maxRequests {
			return generation, fmt.Errorf("circuit breaker is half-open and max requests reached")
		}
	}

	cb.counts.onRequest()
	return generation, nil
}

// afterRequest 请求后处理
func (cb *Breaker) afterRequest(before uint64, success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	generation, state := cb.currentState(now)

	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

// currentState 获取当前状态
func (cb *Breaker) currentState(now time.Time) (uint64, State) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && now.After(cb.expiry) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if !cb.expiry.IsZero() && now.After(cb.expiry) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.generation, cb.state
}

// onSuccess 成功处理
func (cb *Breaker) onSuccess(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onSuccess()
	case StateHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
			cb.setState(StateClosed, now)
		}
	}
}

// onFailure 失败处理
func (cb *Breaker) onFailure(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onFailure()
		if cb.readyToTrip(cb.counts) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

// setState 设置状态
func (cb *Breaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}
}

// toNewGeneration 创建新代
func (cb *Breaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts = Counts{}

	var zero time.Time
	switch cb.state {
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	default: // StateHalfOpen
		cb.expiry = zero
	}
}

// State 获取当前状态
func (cb *Breaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	_, state := cb.currentState(time.Now())
	return state
}

// Counts 获取计数器
func (cb *Breaker) Counts() Counts {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	_, _ = cb.currentState(time.Now())
	return cb.counts
}

// Reset 重置熔断器
func (cb *Breaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateClosed, time.Now())
}

// onRequest 请求计数
func (c *Counts) onRequest() {
	c.Requests++
}

// onSuccess 成功计数
func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

// onFailure 失败计数
func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

// SuccessRate 成功率
func (c *Counts) SuccessRate() float64 {
	if c.Requests == 0 {
		return 0
	}
	return float64(c.TotalSuccesses) / float64(c.Requests)
}

// FailureRate 失败率
func (c *Counts) FailureRate() float64 {
	if c.Requests == 0 {
		return 0
	}
	return float64(c.TotalFailures) / float64(c.Requests)
}

// BreakerManager 熔断器管理器
type BreakerManager struct {
	breakers map[string]*Breaker
	mu       sync.RWMutex
	logger   logger.Logger
}

// NewBreakerManager 创建熔断器管理器
func NewBreakerManager(log logger.Logger) *BreakerManager {
	return &BreakerManager{
		breakers: make(map[string]*Breaker),
		logger:   log,
	}
}

// GetBreaker 获取熔断器
func (bm *BreakerManager) GetBreaker(name string) *Breaker {
	bm.mu.RLock()
	breaker, exists := bm.breakers[name]
	bm.mu.RUnlock()

	if exists {
		return breaker
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 双重检查
	if breaker, exists := bm.breakers[name]; exists {
		return breaker
	}

	// 创建新的熔断器
	cfg := DefaultConfig(name)
	breaker = NewBreaker(cfg, bm.logger)
	bm.breakers[name] = breaker

	bm.logger.Info("circuit breaker created", "name", name)
	return breaker
}

// ListBreakers 列出所有熔断器
func (bm *BreakerManager) ListBreakers() map[string]*Breaker {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	breakers := make(map[string]*Breaker)
	for name, breaker := range bm.breakers {
		breakers[name] = breaker
	}
	return breakers
}

// ResetBreaker 重置熔断器
func (bm *BreakerManager) ResetBreaker(name string) {
	bm.mu.RLock()
	breaker, exists := bm.breakers[name]
	bm.mu.RUnlock()

	if exists {
		breaker.Reset()
		bm.logger.Info("circuit breaker reset", "name", name)
	}
}
