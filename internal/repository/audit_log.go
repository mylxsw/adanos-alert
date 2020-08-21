package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLogType 审计日志类型
type AuditLogType string

const (
	// AuditLogTypeError 错误日志类型
	AuditLogTypeError AuditLogType = "ERROR"
	// AuditLogTypeAction 行为动作类型
	AuditLogTypeAction AuditLogType = "ACTION"
	// AuditLogTypeSystem 系统类型
	AuditLogTypeSystem AuditLogType = "SYSTEM"
)

// AuditLog 审计日志存储
type AuditLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Type    AuditLogType           `bson:"type" json:"type"`
	Context map[string]interface{} `bson:"context" json:"context"`
	Body    string                 `bson:"body" json:"body"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// AuditLogRepo 审计日志仓库
type AuditLogRepo interface {
	Add(al AuditLog) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (al AuditLog, err error)
	Paginate(filter bson.M, offset, limit int64) (logs []AuditLog, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
}
