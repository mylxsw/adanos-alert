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
	Key   string `bson:"key" json:"key" schema:"key"`
	Value string `bson:"value" json:"value" schema:"value"`
}

type UserMetas []UserMeta

func (ums UserMetas) Get(key string) string {
	for _, v := range ums {
		if v.Key == key {
			return v.Value
		}
	}

	return ""
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Phone string `bson:"phone" json:"phone"`

	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`

	Metas UserMetas `bson:"metas" json:"metas"`

	Status UserStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserRepo interface {
	Add(user User) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (user User, err error)
	GetByEmail(email string) (user User, err error)
	Find(filter bson.M) (users []User, err error)
	Paginate(filter bson.M, offset, limit int64) (users []User, next int64, err error)
	DeleteID(id primitive.ObjectID) error
	Delete(filter bson.M) error
	Update(id primitive.ObjectID, user User) error
	Count(filter bson.M) (int64, error)

	GetUserMetas(queryK, queryV, field string) ([]string, error)
}
