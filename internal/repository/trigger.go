package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TriggerStatus string

const (
	TriggerStatusOK     TriggerStatus = "ok"
	TriggerStatusFailed TriggerStatus = "failed"
)

// Trigger is a action trigger for matched rules
type Trigger struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
	// IsElseTrigger 是否是兜底的 Trigger，当所有的非 ElseTrigger 均未匹配时生效
	// IsElseTrigger 为 true 时，忽略 PreCondition 规则，全部匹配
	IsElseTrigger bool                 `bson:"is_else_trigger" json:"is_else_trigger"`
	PreCondition  string               `bson:"pre_condition" json:"pre_condition"`
	Action        string               `bson:"action" json:"action"`
	Meta          string               `bson:"meta" json:"meta"`
	UserRefs      []primitive.ObjectID `bson:"user_refs" json:"user_refs"`
	UserEvalFunc  string               `bson:"user_eval_func" json:"user_eval_func"`
	// for group actions
	Status       TriggerStatus `bson:"trigger_status,omitempty" json:"trigger_status,omitempty"`
	FailedCount  int           `bson:"failed_count" json:"failed_count"`
	FailedReason string        `bson:"failed_reason" json:"failed_reason"`
}
