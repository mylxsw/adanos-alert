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

type EventRelationRepo struct {
	col *mongo.Collection
}

func NewEventRelationRepo(db *mongo.Database) repository.EventRelationRepo {
	col := db.Collection("event_relation")
	if _, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"matched_rule_id": 1},
		Options: options.Index().SetUnique(false),
	}); err != nil {
		log.Errorf("can not create index for event_relation.matched_rule_id: %v", err)
	}

	return &EventRelationRepo{col: col}
}

func (m *EventRelationRepo) AddOrUpdateEventRelation(ctx context.Context, summary string, matchedRuleID primitive.ObjectID) (eventRel repository.EventRelation, err error) {
	err = m.col.FindOneAndUpdate(
		ctx,
		bson.M{"matched_rule_id": matchedRuleID, "summary": summary},
		bson.M{"$inc": bson.M{"event_count": 1}, "updated_at": time.Now()},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&eventRel)

	if err == nil && eventRel.CreatedAt.IsZero() {
		eventRel.CreatedAt = time.Now()
		eventRel.UpdatedAt = eventRel.CreatedAt
		eventRel.EventCount = 1

		_, err = m.col.ReplaceOne(ctx, bson.M{"_id": eventRel.ID}, eventRel)
	}

	return
}

func (m *EventRelationRepo) Get(ctx context.Context, id primitive.ObjectID) (eventRel repository.EventRelation, err error) {
	err = m.col.FindOne(ctx, bson.M{"_id": id}).Decode(&eventRel)
	if err == mongo.ErrNoDocuments {
		return eventRel, repository.ErrNotFound
	}

	return eventRel, err
}

func (m *EventRelationRepo) Paginate(ctx context.Context, filter interface{}, offset, limit int64) (eventRels []repository.EventRelation, next int64, err error) {
	eventRels = make([]repository.EventRelation, 0)
	cur, err := m.col.Find(ctx, filter, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}).SetSkip(offset))
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var evtRel repository.EventRelation
		if err = cur.Decode(&evtRel); err != nil {
			return
		}

		eventRels = append(eventRels, evtRel)
	}

	if int64(len(eventRels)) == limit {
		next = offset + limit
	}

	return eventRels, next, err
}

func (m *EventRelationRepo) Count(ctx context.Context, filter interface{}) (int64, error) {
	return m.col.CountDocuments(ctx, filter)
}
