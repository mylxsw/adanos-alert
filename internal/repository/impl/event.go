package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepo struct {
	col     *mongo.Collection
	seqRepo repository.SequenceRepo
}

func NewEventRepo(db *mongo.Database, seqRepo repository.SequenceRepo) repository.EventRepo {
	col := db.Collection("message")

	if _, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"created_at": 1},
		Options: options.Index().SetUnique(false),
	}); err != nil {
		log.Errorf("can not create index for message.created_at: %v", err)
	}

	if _, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"group_ids": 1},
		Options: options.Index().SetUnique(false),
	}); err != nil {
		log.Errorf("can not create index for message.group_ids: %v", err)
	}

	if _, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"relation_ids": 1},
		Options: options.Index().SetUnique(false),
	}); err != nil {
		log.Errorf("can not create index for message.relation_ids: %v", err)
	}

	return &EventRepo{col: col, seqRepo: seqRepo}
}

func (m EventRepo) AddWithContext(ctx context.Context, msg repository.Event) (id primitive.ObjectID, err error) {
	msg.CreatedAt = time.Now()
	if msg.Status == "" {
		msg.Status = repository.EventStatusPending
	}

	seq, err := m.seqRepo.Next("message_seq")
	if err == nil {
		msg.SeqNum = seq.Value
	}

	if msg.Type == "" {
		msg.Type = repository.EventTypePlain
	}

	rs, err := m.col.InsertOne(ctx, msg)
	if err != nil {
		return id, err
	}

	return rs.InsertedID.(primitive.ObjectID), err
}

func (m EventRepo) Add(msg repository.Event) (id primitive.ObjectID, err error) {
	return m.AddWithContext(context.TODO(), msg)
}

func (m EventRepo) Get(id primitive.ObjectID) (msg repository.Event, err error) {
	err = m.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return msg, repository.ErrNotFound
	}

	return msg, err
}

func (m EventRepo) Find(filter interface{}) (messages []repository.Event, err error) {
	messages = make([]repository.Event, 0)
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var msg repository.Event
		if err = cur.Decode(&msg); err != nil {
			return
		}

		messages = append(messages, msg)
	}

	return
}

func (m EventRepo) Paginate(filter interface{}, offset, limit int64) (messages []repository.Event, next int64, err error) {
	messages = make([]repository.Event, 0)
	cur, err := m.col.Find(context.TODO(), filter, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}).SetSkip(offset))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var msg repository.Event
		if err = cur.Decode(&msg); err != nil {
			return
		}

		messages = append(messages, msg)
	}

	if int64(len(messages)) == limit {
		next = offset + limit
	}

	return messages, next, err
}

func (m EventRepo) Delete(filter interface{}) error {
	_, err := m.col.DeleteMany(context.TODO(), filter)
	return err
}

func (m EventRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m EventRepo) Traverse(filter interface{}, cb func(msg repository.Event) error) error {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	for cur.Next(context.TODO()) {
		var msg repository.Event
		if err = cur.Decode(&msg); err != nil {
			return err
		}

		if err = cb(msg); err != nil {
			return err
		}
	}

	return nil
}

func (m EventRepo) UpdateID(id primitive.ObjectID, update repository.Event) error {
	_, err := m.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (m EventRepo) Count(filter interface{}) (int64, error) {
	return m.col.CountDocuments(context.TODO(), filter)
}
