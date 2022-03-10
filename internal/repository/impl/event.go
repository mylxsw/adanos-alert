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

func (m EventRepo) Has(filter interface{}) (bool, error) {
	var msg repository.Event
	err := m.col.FindOne(context.TODO(), filter).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (m EventRepo) Find(filter interface{}) (messages []repository.Event, err error) {
	messages = make([]repository.Event, 0)
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var msg repository.Event
		if err = cur.Decode(&msg); err != nil {
			return
		}

		messages = append(messages, msg)
	}

	return
}

func (m EventRepo) FindIDs(ctx context.Context, filter interface{}, limit int64) ([]primitive.ObjectID, error) {
	ids := make([]primitive.ObjectID, 0)
	cur, err := m.col.Find(ctx, filter, options.Find().SetLimit(limit).SetProjection(bson.D{{"_id", 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var evt repository.Event
		if err := cur.Decode(&evt); err != nil {
			return nil, err
		}

		ids = append(ids, evt.ID)
	}

	return ids, nil
}

func (m EventRepo) Paginate(filter interface{}, offset, limit int64) (messages []repository.Event, next int64, err error) {
	messages = make([]repository.Event, 0)
	cur, err := m.col.Find(context.TODO(), filter, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}).SetSkip(offset))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

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
	defer cur.Close(context.TODO())

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

func (m EventRepo) CountByDatetime(ctx context.Context, filter bson.M, startTime, endTime time.Time, hour int64) ([]repository.EventByDatetimeCount, error) {
	if filter == nil {
		filter = bson.M{}
	}

	filter["created_at"] = bson.M{"$gt": startTime, "$lte": endTime}

	aggregate, err := m.col.Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", filter}},
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"$toDate": bson.M{
					"$subtract": bson.A{
						bson.M{"$toLong": "$created_at"},
						bson.M{"$mod": bson.A{
							bson.M{"$toLong": "$created_at"},
							1000 * 60 * 60 * hour,
						}},
					},
				},
			},
			"count": bson.M{"$sum": 1},
		}}},
		bson.D{{"$project", bson.M{
			"datetime": "$_id",
			"total":    "$count",
			"_id":      0,
		}}},
		bson.D{{"$sort", bson.M{"datetime": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer aggregate.Close(ctx)

	results := make([]repository.EventByDatetimeCount, 0)
	for aggregate.Next(ctx) {
		var res repository.EventByDatetimeCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}
