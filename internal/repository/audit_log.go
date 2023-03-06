package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SyslogType 系统日志类型
type SyslogType string

const (
	// SyslogTypeError 错误日志类型
	SyslogTypeError SyslogType = "ERROR"
	// SyslogTypeAction 行为动作类型
	SyslogTypeAction SyslogType = "ACTION"
	// SyslogTypeSystem 系统类型
	SyslogTypeSystem SyslogType = "SYSTEM"
)

// Syslog 系统日志存储
type Syslog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Type    SyslogType             `bson:"type" json:"type"`
	Context map[string]interface{} `bson:"context" json:"context"`
	Body    string                 `bson:"body" json:"body"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// SyslogRepo 系统日志仓库
type SyslogRepo interface {
	Add(al Syslog) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (al Syslog, err error)
	Paginate(filter bson.M, offset, limit int64) (logs []Syslog, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
}
