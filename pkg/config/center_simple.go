package config

import (
	"sync"
	"time"
)

// SimpleConfigCenter 简化配置中心
type SimpleConfigCenter struct {
	config   *Config
	watchers map[string][]ConfigWatcher
	mu       sync.RWMutex
}

// ConfigWatcher 配置监听器
type ConfigWatcher interface {
	OnConfigChange(key string, oldValue, newValue interface{})
}

// ConfigChangeEvent 配置变更事件
type ConfigChangeEvent struct {
	Key       string      `json:"key"`
	OldValue  interface{} `json:"old_value"`
	NewValue  interface{} `json:"new_value"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewSimpleConfigCenter 创建简化配置中心
func NewSimpleConfigCenter(cfg *Config) *SimpleConfigCenter {
	return &SimpleConfigCenter{
		config:   cfg,
		watchers: make(map[string][]ConfigWatcher),
	}
}

// RegisterWatcher 注册配置监听器
func (cc *SimpleConfigCenter) RegisterWatcher(key string, watcher ConfigWatcher) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.watchers[key] = append(cc.watchers[key], watcher)
}

// SetConfig 设置配置值
func (cc *SimpleConfigCenter) SetConfig(key string, value interface{}) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// 获取旧值
	oldValue := cc.config.Get(key)

	// 设置新值
	cc.config.Set(key, value)

	// 通知监听器
	cc.notifyWatchers(key, oldValue, value)

	return nil
}

// GetConfig 获取配置值
func (cc *SimpleConfigCenter) GetConfig(key string) interface{} {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.config.Get(key)
}

// notifyWatchers 通知监听器
func (cc *SimpleConfigCenter) notifyWatchers(key string, oldValue, newValue interface{}) {
	watchers := cc.watchers[key]
	for _, watcher := range watchers {
		go func(w ConfigWatcher) {
			defer func() {
				if r := recover(); r != nil {
					// 处理panic
				}
			}()

			w.OnConfigChange(key, oldValue, newValue)
		}(watcher)
	}
}

// SimpleConfigWatcher 简单配置监听器
type SimpleConfigWatcher struct {
	callback func(key string, oldValue, newValue interface{})
}

// OnConfigChange 配置变更回调
func (w *SimpleConfigWatcher) OnConfigChange(key string, oldValue, newValue interface{}) {
	if w.callback != nil {
		w.callback(key, oldValue, newValue)
	}
}

// WatchConfig 监听配置变更
func (cc *SimpleConfigCenter) WatchConfig(key string, callback func(key string, oldValue, newValue interface{})) {
	watcher := &SimpleConfigWatcher{
		callback: callback,
	}

	cc.RegisterWatcher(key, watcher)
}
