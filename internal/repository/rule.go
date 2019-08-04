package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleStatus string

const (
	RuleStatusEnabled  RuleStatus = "enabled"
	RuleStatusDisabled RuleStatus = "disabled"
)

type Rule struct {
	ID          primitive.ObjectID
	Name        string
	Description string

	Interval  int64
	Threshold int64
	Priority  int64

	Status RuleStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RuleRepo interface {
}
