package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	// DefaultUpgrader 默认的 WebSocket 升级器
	DefaultUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有跨域请求（生产环境应该限制）
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Upgrader 配置选项
type UpgraderConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     func(r *http.Request) bool
}

// NewUpgrader 创建自定义升级器
func NewUpgrader(config *UpgraderConfig) *websocket.Upgrader {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  config.ReadBufferSize,
		WriteBufferSize: config.WriteBufferSize,
		CheckOrigin:     config.CheckOrigin,
	}
	
	if upgrader.ReadBufferSize == 0 {
		upgrader.ReadBufferSize = 1024
	}
	if upgrader.WriteBufferSize == 0 {
		upgrader.WriteBufferSize = 1024
	}
	if upgrader.CheckOrigin == nil {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	
	return upgrader
}

