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
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name         string               `bson:"name" json:"name"`
	PreCondition string               `bson:"pre_condition" json:"pre_condition"`
	Action       string               `bson:"action" json:"action"`
	Meta         string               `bson:"meta" json:"meta"`
	UserRefs     []primitive.ObjectID `bson:"user_refs" json:"user_refs"`
	// for group actions
	Status       TriggerStatus `bson:"trigger_status,omitempty" json:"trigger_status,omitempty"`
	FailedCount  int           `bson:"failed_count" json:"failed_count"`
	FailedReason string        `bson:"failed_reason" json:"failed_reason"`
}
