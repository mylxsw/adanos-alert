package repository

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/go-toolkit/collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleRepo struct {
	Rules []repository.Rule
}

func NewRuleRepo() repository.RuleRepo {
	return &RuleRepo{Rules: make([]repository.Rule, 0)}
}

func (r *RuleRepo) Add(rule repository.Rule) (id primitive.ObjectID, err error) {
	rule.ID = primitive.NewObjectID()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = rule.CreatedAt
	for i, tr := range rule.Triggers {
		if tr.ID.IsZero() {
			rule.Triggers[i].ID = primitive.NewObjectID()
		}
	}

	r.Rules = append(r.Rules, rule)
	return rule.ID, nil
}

func (r *RuleRepo) Get(id primitive.ObjectID) (rule repository.Rule, err error) {
	panic("implement me")
}

func (r *RuleRepo) Find(filter bson.M) (rules []repository.Rule, err error) {
	return r.filter(filter), nil
}

func (r *RuleRepo) Traverse(filter bson.M, cb func(rule repository.Rule) error) error {
	panic("implement me")
}

func (r *RuleRepo) UpdateID(id primitive.ObjectID, rule repository.Rule) error {
	panic("implement me")
}

func (r *RuleRepo) Count(filter bson.M) (int64, error) {
	return int64(len(r.filter(filter))), nil
}

func (r *RuleRepo) Delete(filter bson.M) error {
	r.Rules = r.filter(filter)
	return nil
}

func (r *RuleRepo) DeleteID(id primitive.ObjectID) error {
	return r.Delete(bson.M{"_id": id})
}

func (r *RuleRepo) filter(filter bson.M) (rules []repository.Rule) {
	err := collection.MustNew(r.Rules).Filter(func(rule repository.Rule) bool {
		if status, ok := filter["status"]; ok && rule.Status != status {
			return false
		}

		if id, ok := filter["_id"]; ok && id != rule.ID {
			return false
		}

		return true
	}).All(&rules)

	if err != nil {
		panic(err)
	}

	return
}
