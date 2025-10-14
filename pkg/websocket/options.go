package websocket

// NotificationOption 通知选项函数
type NotificationOption func(*NotificationMessage)

// WithIcon 设置通知图标
func WithIcon(icon string) NotificationOption {
	return func(n *NotificationMessage) {
		n.Icon = icon
	}
}

// WithLink 设置通知链接
func WithLink(link string) NotificationOption {
	return func(n *NotificationMessage) {
		n.Link = link
	}
}

// WithCategory 设置通知分类
func WithCategory(category string) NotificationOption {
	return func(n *NotificationMessage) {
		n.Category = category
	}
}

// WithNotificationData 设置通知额外数据
func WithNotificationData(key string, value interface{}) NotificationOption {
	return func(n *NotificationMessage) {
		if n.Data == nil {
			n.Data = make(map[string]interface{})
		}
		n.Data[key] = value
	}
}

// MessageOption 消息选项函数
type MessageOption func(*Message)

// WithMessageID 设置消息ID
func WithMessageID(id string) MessageOption {
	return func(m *Message) {
		m.ID = id
	}
}

// WithRoom 设置消息房间
func WithRoom(room string) MessageOption {
	return func(m *Message) {
		m.Room = room
	}
}

// WithMessageData 设置消息额外数据
func WithMessageData(key string, value interface{}) MessageOption {
	return func(m *Message) {
		m.SetData(key, value)
	}
}

// SystemMessageOption 系统消息选项函数
type SystemMessageOption func(*SystemMessage)

// WithCode 设置系统消息代码
func WithCode(code int) SystemMessageOption {
	return func(s *SystemMessage) {
		s.Code = code
	}
}

// WithLevel 设置系统消息级别
func WithLevel(level string) SystemMessageOption {
	return func(s *SystemMessage) {
		s.Level = level
	}
}

// WithSystemData 设置系统消息额外数据
func WithSystemData(key string, value interface{}) SystemMessageOption {
	return func(s *SystemMessage) {
		if s.Data == nil {
			s.Data = make(map[string]interface{})
		}
		s.Data[key] = value
	}
}

// DataUpdateOption 数据更新选项函数
type DataUpdateOption func(*DataUpdateMessage)

// WithEntityID 设置实体ID
func WithEntityID(id string) DataUpdateOption {
	return func(d *DataUpdateMessage) {
		d.ID = id
	}
}

// WithUpdateData 设置更新数据
func WithUpdateData(data map[string]interface{}) DataUpdateOption {
	return func(d *DataUpdateMessage) {
		d.Data = data
	}
}
