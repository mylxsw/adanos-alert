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

type RuleRepo struct {
	col *mongo.Collection
}

func NewRuleRepo(db *mongo.Database) repository.RuleRepo {
	return &RuleRepo{col: db.Collection("rule")}
}

func (r RuleRepo) Add(rule repository.Rule) (id primitive.ObjectID, err error) {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = rule.CreatedAt

	for i, t := range rule.Triggers {
		if t.ID.IsZero() {
			rule.Triggers[i].ID = primitive.NewObjectID()
		}
	}

	rs, err := r.col.InsertOne(context.TODO(), rule)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (r RuleRepo) Get(id primitive.ObjectID) (rule repository.Rule, err error) {
	err = r.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&rule)
	return
}

func (r RuleRepo) Paginate(filter interface{}, offset, limit int64) (rules []repository.Rule, next int64, err error) {
	rules = make([]repository.Rule, 0)
	cur, err := r.col.Find(context.TODO(), filter, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}).SetSkip(offset))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var rule repository.Rule
		if err = cur.Decode(&rule); err != nil {
			return
		}

		rules = append(rules, rule)
	}

	if int64(len(rules)) == limit {
		next = offset + limit
	}

	return rules, next, err
}

func (r RuleRepo) Find(filter bson.M) (rules []repository.Rule, err error) {
	rules = make([]repository.Rule, 0)
	cur, err := r.col.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "priority", Value: -1}}))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var rule repository.Rule
		if err = cur.Decode(&rule); err != nil {
			return
		}

		rules = append(rules, rule)
	}

	return
}

func (r RuleRepo) Traverse(filter bson.M, cb func(rule repository.Rule) error) error {
	cur, err := r.col.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "priority", Value: -1}}))
	if err != nil {
		return err
	}

	for cur.Next(context.TODO()) {
		var rule repository.Rule
		if err = cur.Decode(&rule); err != nil {
			return err
		}

		if err = cb(rule); err != nil {
			return err
		}
	}

	return nil
}

func (r RuleRepo) UpdateID(id primitive.ObjectID, rule repository.Rule) error {
	for i, t := range rule.Triggers {
		if t.ID.IsZero() {
			rule.Triggers[i].ID = primitive.NewObjectID()
		}
	}

	rule.UpdatedAt = time.Now()

	_, err := r.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, rule)
	return err
}

func (r RuleRepo) Count(filter bson.M) (int64, error) {
	return r.col.CountDocuments(context.TODO(), filter)
}

func (r RuleRepo) Delete(filter bson.M) error {
	_, err := r.col.DeleteMany(context.TODO(), filter)
	return err
}

func (r RuleRepo) DeleteID(id primitive.ObjectID) error {
	return r.Delete(bson.M{"_id": id})
}
