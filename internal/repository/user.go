package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusDisabled UserStatus = "disabled"
)

type UserMeta struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	Name  string          `bson:"name" json:"name"`
	Metas []UserMeta      `bson:"metas" json:"metas"`

	Status UserStatus     `bson:"status" json:"status"`

	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}

type UserRepo interface {
	Add(user User) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (user User, err error)
	Find(filter bson.M) (users []User, err error)
	Paginate(filter bson.M, offset, limit int) (users []User, next int, err error)
	Delete(id primitive.ObjectID) error
	Update(id primitive.ObjectID, user User) error
}
