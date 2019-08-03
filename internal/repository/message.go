package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageTag struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content   string             `bson:"content" json:"content"`
	Tags      []MessageTag       `bson:"tags" json:"tags"`
	Origin    string             `bson:"origin" json:"origin"`
	GroupID   primitive.ObjectID `bson:"group_id" json:"group_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type MessageRepo interface {
	Add(msg Message) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (msg Message, err error)
	Find(filter bson.M) (messages []Message, err error)
	Paginate(filter bson.M, offset, limit int64) (messages []Message, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Traverse(filter bson.M, cb func(msg Message) error) error
	UpdateID(id primitive.ObjectID, update Message) error
	Count(filter bson.M) (int64, error)
}
