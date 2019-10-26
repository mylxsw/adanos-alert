package repository

import (
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

type QueueItem struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Payload      string             `bson:"payload" json:"payload"`
	Status       QueueItemStatus    `bson:"status" json:"status"`
	RequeueTimes int                `bson:"requeue_times" json:"requeue_times"`

	NextExecuteAt time.Time `bson:"next_execute_at" json:"next_execute_at"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type QueueRepo interface {
	// Enqueue add a item to queue
	// if the item is new(id is empty), add it to queue
	// if the item is already existed, replace it
	Enqueue(item QueueItem) (primitive.ObjectID, error)
	Dequeue() (QueueItem, error)
	UpdateID(id primitive.ObjectID, item QueueItem) error
	Update(filter bson.M, item QueueItem) error
	Paginate(filter bson.M, offset, limit int64) (items []QueueItem, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Get(id primitive.ObjectID) (QueueItem, error)
	Count(filter bson.M) (int64, error)
}
