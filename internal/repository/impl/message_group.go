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

type MessageGroupRepo struct {
	col *mongo.Collection
}

func NewMessageGroupRepo(db *mongo.Database) repository.MessageGroupRepo {
	return &MessageGroupRepo{col: db.Collection("message_group")}
}

func (m MessageGroupRepo) Add(grp repository.MessageGroup) (id primitive.ObjectID, err error) {
	grp.CreatedAt = time.Now()
	grp.UpdatedAt = grp.CreatedAt
	rs, err := m.col.InsertOne(context.TODO(), grp)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (m MessageGroupRepo) Get(id primitive.ObjectID) (grp repository.MessageGroup, err error) {
	err = m.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&grp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}
	return
}

func (m MessageGroupRepo) Find(filter bson.M) (grps []repository.MessageGroup, err error) {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var grp repository.MessageGroup
		if err = cur.Decode(&grp); err != nil {
			return
		}

		grps = append(grps, grp)
	}

	return
}

func (m MessageGroupRepo) Paginate(filter bson.M, offset, limit int64) (grps []repository.MessageGroup, next int64, err error) {
	cur, err := m.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var grp repository.MessageGroup
		if err = cur.Decode(&grp); err != nil {
			return
		}

		grps = append(grps, grp)
	}

	if int64(len(grps)) == limit {
		next = offset + limit
	}

	return
}

func (m MessageGroupRepo) Traverse(filter bson.M, cb func(grp repository.MessageGroup) error) error {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	for cur.Next(context.TODO()) {
		var grp repository.MessageGroup
		if err = cur.Decode(&grp); err != nil {
			return err
		}

		if err = cb(grp); err != nil {
			return err
		}
	}

	return nil
}

func (m MessageGroupRepo) Update(id primitive.ObjectID, grp repository.MessageGroup) error {
	grp.UpdatedAt = time.Now()
	_, err := m.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, grp)
	return err
}

func (m MessageGroupRepo) Delete(filter bson.M) error {
	_, err := m.col.DeleteMany(context.TODO(), filter)
	return err
}

func (m MessageGroupRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m MessageGroupRepo) Count(filter bson.M) (int64, error) {
	return m.col.CountDocuments(context.TODO(), filter)
}
