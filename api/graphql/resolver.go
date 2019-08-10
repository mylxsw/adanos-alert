package graphql

import (
	"context"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	cc       *container.Container
	ruleRepo repository.RuleRepo
}

func NewResolver(cc *container.Container) *Resolver {
	res := Resolver{cc: cc}
	cc.MustResolve(func(ruleRepo repository.RuleRepo) {
		res.ruleRepo = ruleRepo
	})

	return &res
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) AddRule(ctx context.Context, rule *NewRule) (*Rule, error) {
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

type queryResolver struct{ *Resolver }

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
