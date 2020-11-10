package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventRelationNote 事件关联备注
type EventRelationNote struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RelationID  primitive.ObjectID `bson:"relation_id,omitempty" json:"relation_id"`
	EventID     primitive.ObjectID `bson:"event_id,omitempty" json:"event_id"`
	Note        string             `bson:"note" json:"note"`
	CreatorID   primitive.ObjectID `bson:"creator_id,omitempty" json:"creator_id"`
	CreatorName string             `bson:"creator_name" json:"creator_name"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	DeletedAt   time.Time          `bson:"deleted_at" json:"deleted_at"`
}

// EventRelationNoteRepo 事件关联备注仓库接口
type EventRelationNoteRepo interface {
	AddNote(ctx context.Context, note EventRelationNote) (ID, error)
	PaginateNotes(ctx context.Context, relID primitive.ObjectID, filter bson.M, offset, limit int64) (notes []EventRelationNote, next int64, err error)
	DeleteNote(ctx context.Context, relID primitive.ObjectID, filter bson.M) error
}
