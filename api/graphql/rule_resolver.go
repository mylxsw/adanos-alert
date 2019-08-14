package graphql

import (
	"context"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Rule(ctx context.Context, id string) (*Rule, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	rule, err := r.ruleRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	return RepoToRule(rule), nil
}

func (r *queryResolver) Rules(ctx context.Context) ([]*Rule, error) {
	rules, err := r.ruleRepo.Find(bson.M{})
	if err != nil {
		return nil, err
	}

	var results []*Rule
	return results, coll.Map(rules, &results, func(rule repository.Rule) *Rule {
		return RepoToRule(rule)
	})
}

func (r *mutationResolver) AddRule(ctx context.Context, rule NewRule) (*Rule, error) {
	id, err := r.ruleRepo.Add(RuleToRepo(rule))
	if err != nil {
		return nil, err
	}

	res, err := r.ruleRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return RepoToRule(res), nil
}

func (r *mutationResolver) UpdateRule(ctx context.Context, id string, rule NewRule) (*Rule, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.ruleRepo.UpdateID(mid, RuleToRepo(rule))
	if err != nil {
		return nil, err
	}

	res, err := r.ruleRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	return RepoToRule(res), nil
}

func (r *mutationResolver) DeleteRule(ctx context.Context, id string) (*Rule, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	rule, err := r.ruleRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	if err := r.ruleRepo.DeleteID(mid); err != nil {
		return nil, err
	}

	return RepoToRule(rule), nil
}
