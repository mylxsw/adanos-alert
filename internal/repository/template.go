package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateType string

const (
	TemplateTypeMatchRule   TemplateType = "match_rule"
	TemplateTypeTemplate    TemplateType = "template"
	TemplateTypeTriggerRule TemplateType = "trigger_rule"
)

type Template struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Content     string             `bson:"content" json:"content"`
	Type        TemplateType       `bson:"type" json:"type"`
	Predefined  bool               `bson:"predefined" json:"predefined"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type TemplateRepo interface {
	Add(temp Template) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (temp Template, err error)
	Find(filter bson.M) (templates []Template, err error)
	Paginate(filter bson.M, offset, limit int64) (templates []Template, next int64, err error)
	DeleteID(id primitive.ObjectID) error
	Delete(filter bson.M) error
	Update(id primitive.ObjectID, temp Template) error
	Count(filter bson.M) (int64, error)
}
