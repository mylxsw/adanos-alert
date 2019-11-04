package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TemplateRepo struct {
	col *mongo.Collection
}

func NewTemplateRepo(db *mongo.Database) repository.TemplateRepo {
	return &TemplateRepo{col: db.Collection("template")}
}

func (t TemplateRepo) Add(temp repository.Template) (id primitive.ObjectID, err error) {
	temp.CreatedAt = time.Now()
	temp.UpdatedAt = temp.CreatedAt

	rs, err := t.col.InsertOne(context.TODO(), temp)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (t TemplateRepo) Get(id primitive.ObjectID) (temp repository.Template, err error) {
	err = t.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&temp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (t TemplateRepo) Find(filter bson.M) (templates []repository.Template, err error) {
	templates = make([]repository.Template, 0)
	cur, err := t.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var temp repository.Template
		if err = cur.Decode(&temp); err != nil {
			return
		}

		templates = append(templates, temp)
	}

	return
}

func (t TemplateRepo) Paginate(filter bson.M, offset, limit int64) (templates []repository.Template, next int64, err error) {
	templates = make([]repository.Template, 0)
	cur, err := t.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var temp repository.Template
		if err = cur.Decode(&temp); err != nil {
			return
		}

		templates = append(templates, temp)
	}

	if int64(len(templates)) == limit {
		next = offset + limit
	}

	return
}

func (t TemplateRepo) DeleteID(id primitive.ObjectID) error {
	return t.Delete(bson.M{"_id": id})
}

func (t TemplateRepo) Delete(filter bson.M) error {
	_, err := t.col.DeleteMany(context.TODO(), filter)
	return err
}

func (t TemplateRepo) Update(id primitive.ObjectID, temp repository.Template) error {
	temp.UpdatedAt = time.Now()
	_, err := t.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, temp)
	return err
}

func (t TemplateRepo) Count(filter bson.M) (int64, error) {
	return t.col.CountDocuments(context.TODO(), filter)
}
