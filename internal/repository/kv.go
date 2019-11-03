package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type KV struct {
	Key       string      `json:"key" bson:"key"`
	Value     interface{} `json:"value" bson:"value"`
	ExpiredAt time.Time   `json:"expired_at" bson:"expired_at"`
	WithTTL   bool        `json:"with_ttl" bson:"with_ttl"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
}

type KVRepo interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string) (pair KV, err error)
	Remove(key string) (removeCount int64, err error)
	All(filter bson.M) (pairs []KV, err error)
	GC() error
}
