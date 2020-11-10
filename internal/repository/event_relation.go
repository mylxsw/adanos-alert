package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventRelation 事件关联
type EventRelation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MatchedRuleID primitive.ObjectID `bson:"matched_rule_id,omitempty" json:"matched_rule_id,omitempty"`
	Summary       string             `bson:"summary" json:"summary,omitempty"`
	EventCount    int64              `bson:"event_count" json:"event_count,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

// EventRelationRepo 事件关联仓库接口
type EventRelationRepo interface {
	AddOrUpdateEventRelation(ctx context.Context, summary string, matchedRuleID primitive.ObjectID) (EventRelation, error)
	Get(ctx context.Context, id primitive.ObjectID) (eventRel EventRelation, err error)
	Paginate(ctx context.Context, filter interface{}, offset, limit int64) (eventRels []EventRelation, next int64, err error)
	Count(ctx context.Context, filter interface{}) (int64, error)
}
