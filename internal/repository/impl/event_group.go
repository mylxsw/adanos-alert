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

type EventGroupRepo struct {
	col     *mongo.Collection
	seqRepo repository.SequenceRepo
}

func NewEventGroupRepo(db *mongo.Database, seqRepo repository.SequenceRepo) repository.EventGroupRepo {
	grp := db.Collection("message_group")
	_, err := grp.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"created_at": 1},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		log.Errorf("can not create index for message_group.created_at: %v", err)
	}

	return &EventGroupRepo{col: grp, seqRepo: seqRepo}
}

func (m EventGroupRepo) Add(grp repository.EventGroup) (id primitive.ObjectID, err error) {
	grp.CreatedAt = time.Now()
	grp.UpdatedAt = grp.CreatedAt
	seq, err := m.seqRepo.Next("group_seq")
	if err == nil {
		grp.SeqNum = seq.Value
	}

	rs, err := m.col.InsertOne(context.TODO(), grp)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (m EventGroupRepo) Get(id primitive.ObjectID) (grp repository.EventGroup, err error) {
	err = m.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&grp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}
	return
}

func (m EventGroupRepo) Find(filter bson.M) (grps []repository.EventGroup, err error) {
	grps = make([]repository.EventGroup, 0)
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var grp repository.EventGroup
		if err = cur.Decode(&grp); err != nil {
			return
		}

		grps = append(grps, grp)
	}

	return
}

func (m EventGroupRepo) Paginate(filter bson.M, offset, limit int64) (grps []repository.EventGroup, next int64, err error) {
	grps = make([]repository.EventGroup, 0)
	cur, err := m.col.Find(
		context.TODO(),
		filter,
		options.Find().
			SetSkip(offset).
			SetLimit(limit).
			SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var grp repository.EventGroup
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

func (m EventGroupRepo) Traverse(filter bson.M, cb func(grp repository.EventGroup) error) error {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var grp repository.EventGroup
		if err = cur.Decode(&grp); err != nil {
			return err
		}

		if err = cb(grp); err != nil {
			return err
		}
	}

	return nil
}

func (m EventGroupRepo) UpdateID(id primitive.ObjectID, grp repository.EventGroup) error {
	grp.UpdatedAt = time.Now()
	_, err := m.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, grp)
	return err
}

func (m EventGroupRepo) Delete(filter bson.M) error {
	_, err := m.col.DeleteMany(context.TODO(), filter)
	return err
}

func (m EventGroupRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m EventGroupRepo) Count(filter bson.M) (int64, error) {
	return m.col.CountDocuments(context.TODO(), filter)
}

func (m EventGroupRepo) CollectingGroup(rule repository.EventGroupRule) (group repository.EventGroup, err error) {
	err = m.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{"rule._id": rule.ID, "rule.aggregate_key": rule.AggregateKey, "rule.type": rule.Type, "status": repository.EventGroupStatusCollecting},
		bson.M{"$set": bson.M{"status": repository.EventGroupStatusCollecting}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&group)

	// since we create a group automatically, we need update created_at fields manually
	if err == nil && group.CreatedAt.IsZero() {
		group.CreatedAt = time.Now()
		group.UpdatedAt = group.CreatedAt
		seq, err := m.seqRepo.Next("group_seq")
		if err == nil {
			group.SeqNum = seq.Value
		}

		group.Rule = rule
		group.AggregateKey = rule.AggregateKey
		group.Type = rule.Type

		_ = m.UpdateID(group.ID, group)
	}

	return
}

func (m EventGroupRepo) LastGroup(filter bson.M) (grp repository.EventGroup, err error) {
	rs := m.col.FindOne(context.TODO(), filter, options.FindOne().SetSort(bson.M{"updated_at": -1}))
	err = rs.Decode(&grp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}
	return grp, err
}

func (m EventGroupRepo) StatByRuleCount(ctx context.Context, startTime, endTime time.Time) ([]repository.EventGroupByRuleCount, error) {
	aggregate, err := m.col.Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"updated_at": bson.M{"$gt": startTime, "$lte": endTime}}}},
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"rule_id":   "$rule._id",
				"rule_name": "$rule.name",
			},
			"count":         bson.M{"$sum": 1},
			"message_count": bson.M{"$sum": "$message_count"},
		}}},
		bson.D{{"$project", bson.M{
			"rule_id":        "$_id.rule_id",
			"rule_name":      "$_id.rule_name",
			"total":          "$count",
			"total_messages": "$message_count",
			"_id":            0,
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer aggregate.Close(ctx)

	results := make([]repository.EventGroupByRuleCount, 0)
	for aggregate.Next(ctx) {
		var res repository.EventGroupByRuleCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}

func (m EventGroupRepo) StatByUserCount(ctx context.Context, startTime, endTime time.Time) ([]repository.EventGroupByUserCount, error) {
	aggregate, err := m.col.Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"updated_at": bson.M{"$gt": startTime, "$lte": endTime}}}},
		bson.D{{"$unwind", "$actions"}},
		bson.D{{"$unwind", "$actions.user_refs"}},
		bson.D{{"$group", bson.M{
			"_id":           "$actions.user_refs",
			"count":         bson.M{"$sum": 1},
			"message_count": bson.M{"$sum": "$message_count"},
		}}},
		bson.D{{"$lookup", bson.M{
			"localField":   "_id",
			"from":         "user",
			"foreignField": "_id",
			"as":           "user",
		}}},
		bson.D{{"$replaceRoot", bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					bson.M{
						"_id":           "$user._id",
						"name":          "$user.name",
						"count":         "$count",
						"message_count": "$message_count",
					},
					"$$ROOT",
				},
			},
		}}},
		bson.D{{"$unwind", "$name"}},
		bson.D{{"$project", bson.M{
			"user_id":        "$_id",
			"user_name":      "$name",
			"total":          "$count",
			"total_messages": "$message_count",
			"_id":            0,
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer aggregate.Close(ctx)

	results := make([]repository.EventGroupByUserCount, 0)
	for aggregate.Next(ctx) {
		var res repository.EventGroupByUserCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}

func (m EventGroupRepo) StatByDatetimeCount(ctx context.Context, startTime, endTime time.Time, hour int64) ([]repository.EventGroupByDatetimeCount, error) {
	unixTime, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z00:00")
	aggregate, err := m.col.Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"updated_at": bson.M{"$gt": startTime, "$lte": endTime}}}},
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"$subtract": bson.A{
					bson.M{"$subtract": bson.A{"$updated_at", unixTime}},
					bson.M{"$mod": bson.A{
						bson.M{"$subtract": bson.A{"$updated_at", unixTime}},
						1000 * 60 * 60 * hour,
					}},
				},
			},
			"count":         bson.M{"$sum": 1},
			"message_count": bson.M{"$sum": "$message_count"},
		}}},
		bson.D{{"$project", bson.M{
			"datetime":       bson.M{"$add": bson.A{unixTime, "$_id"}},
			"total":          "$count",
			"total_messages": "$message_count",
			"_id":            0,
		}}},
		bson.D{{"$sort", bson.M{"datetime": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer aggregate.Close(ctx)

	results := make([]repository.EventGroupByDatetimeCount, 0)
	for aggregate.Next(ctx) {
		var res repository.EventGroupByDatetimeCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}
