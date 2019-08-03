package impl

import (
	"context"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SequenceRepo struct {
	col *mongo.Collection
}

func NewSequenceRepo(db *mongo.Database) repository.SequenceRepo {
	return &SequenceRepo{col: db.Collection("sequence")}
}

func (s SequenceRepo) Next(name string) (seq repository.Sequence, err error) {
	rs := s.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{"name": name},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if err := rs.Decode(&seq); err != nil {
		return seq, err
	}

	return
}

func (s SequenceRepo) Truncate(name string) error {
	_, err := s.col.DeleteOne(context.TODO(), bson.M{"name": name})
	return err
}
