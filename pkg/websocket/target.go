package websocket

import "errors"

// TargetType 消息目标类型
type TargetType string

const (
	// TargetTypeUser 用户目标
	TargetTypeUser TargetType = "user"
	
	// TargetTypeGroup 分组目标
	TargetTypeGroup TargetType = "group"
	
	// TargetTypeClient 客户端目标
	TargetTypeClient TargetType = "client"
	
	// TargetTypeAll 所有客户端目标
	TargetTypeAll TargetType = "all"
)

// Target 消息目标
type Target struct {
	Type TargetType // 目标类型
	ID   string     // 目标ID
}

// NewUserTarget 创建用户目标
func NewUserTarget(userID string) Target {
	return Target{
		Type: TargetTypeUser,
		ID:   userID,
	}
}

// NewGroupTarget 创建分组目标
func NewGroupTarget(groupID string) Target {
	return Target{
		Type: TargetTypeGroup,
		ID:   groupID,
	}
}

// NewClientTarget 创建客户端目标
func NewClientTarget(clientID string) Target {
	return Target{
		Type: TargetTypeClient,
		ID:   clientID,
	}
}

// NewBroadcastTarget 创建广播目标
func NewBroadcastTarget() Target {
	return Target{
		Type: TargetTypeAll,
	}
}

// Validate 验证目标是否有效
func (t Target) Validate() error {
	switch t.Type {
	case TargetTypeUser, TargetTypeGroup, TargetTypeClient:
		if t.ID == "" {
			return errors.New("target ID cannot be empty")
		}
	case TargetTypeAll:
		// 广播目标不需要ID
	default:
		return errors.New("invalid target type")
	}
	
	return nil
}
