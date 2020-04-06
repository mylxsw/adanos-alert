package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DingdingRobot struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`

	Token  string `bson:"token" json:"token"`
	Secret string `bson:"secret" json:"secret"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type DingdingRobotRepo interface {
	Add(robot DingdingRobot) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (robot DingdingRobot, err error)
	Find(filter bson.M) (robots []DingdingRobot, err error)
	Paginate(filter bson.M, offset, limit int64) (robots []DingdingRobot, next int64, err error)
	DeleteID(id primitive.ObjectID) error
	Delete(filter bson.M) error
	Update(id primitive.ObjectID, robot DingdingRobot) error
	Count(filter bson.M) (int64, error)
}
