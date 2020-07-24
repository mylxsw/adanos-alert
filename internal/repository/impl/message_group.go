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
	col     *mongo.Collection
	seqRepo repository.SequenceRepo
}

func NewMessageGroupRepo(db *mongo.Database, seqRepo repository.SequenceRepo) repository.MessageGroupRepo {
	return &MessageGroupRepo{col: db.Collection("message_group"), seqRepo: seqRepo}
}

func (m MessageGroupRepo) Add(grp repository.MessageGroup) (id primitive.ObjectID, err error) {
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

func (m MessageGroupRepo) Get(id primitive.ObjectID) (grp repository.MessageGroup, err error) {
	err = m.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&grp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}
	return
}

func (m MessageGroupRepo) Find(filter bson.M) (grps []repository.MessageGroup, err error) {
	grps = make([]repository.MessageGroup, 0)
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
	grps = make([]repository.MessageGroup, 0)
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

func (m MessageGroupRepo) UpdateID(id primitive.ObjectID, grp repository.MessageGroup) error {
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

func (m MessageGroupRepo) CollectingGroup(rule repository.MessageGroupRule) (group repository.MessageGroup, err error) {
	err = m.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{"rule._id": rule.ID, "status": repository.MessageGroupStatusCollecting},
		bson.M{"$set": bson.M{"status": repository.MessageGroupStatusCollecting}},
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
		_ = m.UpdateID(group.ID, group)
	}

	return
}

func (m MessageGroupRepo) LastGroup(filter bson.M) (grp repository.MessageGroup, err error) {
	rs := m.col.FindOne(context.TODO(), filter, options.FindOne().SetSort(bson.M{"updated_at": -1}))
	err = rs.Decode(&grp)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}
	return grp, err
}

func (m MessageGroupRepo) StatByRuleCount(ctx context.Context, startTime, endTime time.Time) ([]repository.MessageGroupByRuleCount, error) {
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

	results := make([]repository.MessageGroupByRuleCount, 0)
	for aggregate.Next(ctx) {
		var res repository.MessageGroupByRuleCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}

func (m MessageGroupRepo) StatByUserCount(ctx context.Context, startTime, endTime time.Time) ([]repository.MessageGroupByUserCount, error) {
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

	results := make([]repository.MessageGroupByUserCount, 0)
	for aggregate.Next(ctx) {
		var res repository.MessageGroupByUserCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}

func (m MessageGroupRepo) StatByDatetimeCount(ctx context.Context, startTime, endTime time.Time, hour int64) ([]repository.MessageGroupByDatetimeCount, error) {
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

	results := make([]repository.MessageGroupByDatetimeCount, 0)
	for aggregate.Next(ctx) {
		var res repository.MessageGroupByDatetimeCount
		if err := aggregate.Decode(&res); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}
