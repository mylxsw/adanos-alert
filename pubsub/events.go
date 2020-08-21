package pubsub

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
)

// EventType 事件类型
type EventType string

const (
	// EventTypeAdd 新增事件
	EventTypeAdd EventType = "added"
	// EventTypeUpdate 更新事件
	EventTypeUpdate EventType = "updated"
	// EventTypeDelete 删除事件
	EventTypeDelete EventType = "deleted"
)

// RuleChangedEvent 规则变更事件
type RuleChangedEvent struct {
	Rule      repository.Rule
	Type      EventType
	CreatedAt time.Time
}

// DingdingRobotEvent 钉钉机器人变更事件
type DingdingRobotEvent struct {
	DingDingRobot repository.DingdingRobot
	Type          EventType
	CreatedAt     time.Time
}

// UserChangedEvent 用户变更事件
type UserChangedEvent struct {
	User      repository.User
	Type      EventType
	CreatedAt time.Time
}

// SystemUpDownEvent 系统启停事件
type SystemUpDownEvent struct {
	Up        bool
	CreatedAt time.Time
}
