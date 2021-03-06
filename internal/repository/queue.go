package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueItemStatus string

const (
	QueueItemStatusWait     QueueItemStatus = "wait"
	QueueItemStatusRunning  QueueItemStatus = "running"
	QueueItemStatusFailed   QueueItemStatus = "failed"
	QueueItemStatusSucceed  QueueItemStatus = "succeed"
	QueueItemStatusCanceled QueueItemStatus = "canceled"
)

type QueueJob struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Payload      string             `bson:"payload" json:"payload"`
	Status       QueueItemStatus    `bson:"status" json:"status"`
	LastError    string             `bson:"last_error" json:"last_error"`
	RequeueTimes int                `bson:"requeue_times" json:"requeue_times"`

	NextExecuteAt time.Time `bson:"next_execute_at" json:"next_execute_at"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type QueueRepo interface {
	// Enqueue add a item to queue
	// if the item is new(id is empty), add it to queue
	// if the item is already existed, replace it
	Enqueue(ctx context.Context, item QueueJob) (primitive.ObjectID, error)
	Dequeue(ctx context.Context) (QueueJob, error)
	UpdateID(ctx context.Context, id primitive.ObjectID, jobItem QueueJob) error
	Update(ctx context.Context, filter bson.M, item QueueJob) error
	Paginate(filter bson.M, offset, limit int64) (items []QueueJob, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Get(id primitive.ObjectID) (QueueJob, error)
	Count(filter bson.M) (int64, error)
}
