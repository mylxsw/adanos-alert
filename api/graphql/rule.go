package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/mylxsw/adanos-alert/api/view"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	gql "github.com/mylxsw/adanos-alert/pkg/graphql"
	"github.com/mylxsw/go-toolkit/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ruleType = gql.Object(view.Rule{})

type RuleObject struct {
	cc *container.Container
}

func (r *RuleObject) Register(builder gql.GraphQL) {
	r.cc.MustResolve(func(ruleRepo repository.RuleRepo) {
		builder.Query("rules", r.GetAllRules(ruleRepo))
		builder.Mutation("addRule", r.AddRule(ruleRepo))
		builder.Mutation("deleteRule", r.DeleteRule(ruleRepo))
	})
}

func NewRuleObject(cc *container.Container) *RuleObject {
	return &RuleObject{cc: cc}
}

func (r *RuleObject) GetAllRules(ruleRepo repository.RuleRepo) *graphql.Field {
	return gql.CreateField(
		graphql.NewList(ruleType),
		gql.BindArgs(struct{ ID string `json:"id"` }{}),
		func(arg struct{ ID string `json:"id"` }) (interface{}, error) {
			filter := bson.M{}

			if arg.ID != "" {
				id, err := primitive.ObjectIDFromHex(arg.ID)
				if err != nil {
					return nil, err
				}

				filter["_id"] = id
			}

			rules, err := ruleRepo.Find(filter)
			if err != nil {
				return nil, err
			}

			return view.RulesFromRepos(rules), nil
		},
	)
}

// func (r *RuleObject) GetAllRules(ruleRepo repository.RuleRepo) *graphql.Field {
// 	return &graphql.Field{
// 		Type: graphql.NewList(ruleType),
// 		Args: graphql.FieldConfigArgument{
// 			"id": &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "",},
// 		},
// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			filter := bson.M{}
//
// 			idHex := p.Args["id"].(string)
// 			if idHex != "" {
// 				id, err := primitive.ObjectIDFromHex(idHex)
// 				if err != nil {
// 					return nil, err
// 				}
//
// 				filter["_id"] = id
// 			}
//
// 			rules, err := ruleRepo.Find(filter)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			return view.RulesFromRepos(rules), nil
// 		},
// 	}
// }

func (r *RuleObject) AddRule(repo repository.RuleRepo) *graphql.Field {
	args := gql.BindArgs(view.Rule{}, "name", "description", "interval", "threshold", "priority", "triggers", "rule", "template", "summary_template", "status")
	for name, arg := range args {
		fmt.Println(name, arg.Type, arg.Type.Name())
	}
	return gql.CreateField(
		ruleType,
		gql.BindArgs(view.Rule{}, "name", "description", "interval", "threshold", "priority", "triggers", "rule", "template", "summary_template", "status"),
		func(rule view.Rule) (interface{}, error) {
			value := view.RuleToRepo(rule)

			_, err := matcher.NewMessageMatcher(value)
			if err != nil {
				return nil, err
			}

			id, err := repo.Add(value)
			if err != nil {
				return nil, err
			}

			rs, err := repo.Get(id)
			if err != nil {
				return nil, err
			}

			return view.RuleFromRepo(rs), nil
		},
	)
}

// func (r *RuleObject) AddRule(repo repository.RuleRepo) *graphql.Field {
// 	return &graphql.Field{
// 		Type: ruleType,
// 		// Args: graphql.FieldConfigArgument{
// 		// 	"name":             &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "Graphql",},
// 		// 	"description":      &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "",},
// 		// 	"threshold":        &graphql.ArgumentConfig{Type: Int64, DefaultValue: int64(0),},
// 		// 	"priority":         &graphql.ArgumentConfig{Type: Int64, DefaultValue: int64(0),},
// 		// 	"interval":         &graphql.ArgumentConfig{Type: Int64, DefaultValue: int64(0),},
// 		// 	"rule":             &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "",},
// 		// 	"template":         &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "",},
// 		// 	"summary_template": &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "",},
// 		// 	"status":           &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: repository.RuleStatusEnabled,},
// 		// },
// 		Args: graphql.BindArg(view.Rule{}, "name", "description", "interval", "threshold", "priority", "rule", "template", "summary_template", "status"),
// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			var rule view.Rule
// 			Parse(p, &rule)
// 			value := view.RuleToRepo(rule)
// 			// value := view.RuleToRepo(view.Rule{
// 			// 	Name:            p.Args["name"].(string),
// 			// 	Description:     p.Args["description"].(string),
// 			// 	Interval:        p.Args["interval"].(int64),
// 			// 	Threshold:       p.Args["threshold"].(int64),
// 			// 	Priority:        p.Args["priority"].(int64),
// 			// 	Rule:            p.Args["rule"].(string),
// 			// 	Template:        p.Args["template"].(string),
// 			// 	SummaryTemplate: p.Args["summary_template"].(string),
// 			// 	Status:          p.Args["status"].(repository.RuleStatus),
// 			// })
//
// 			_, err := ruleMatcher.NewMessageMatcher(value)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			id, err := repo.Add(value)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			rs, err := repo.Get(id)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			return view.RuleFromRepo(rs), nil
// 		},
// 	}
// }

func (r *RuleObject) DeleteRule(repo repository.RuleRepo) *graphql.Field {
	return gql.CreateField(
		ruleType,
		gql.BindArgs(struct{ ID string `json:"id"` }{}),
		func(arg struct{ ID string `json:"id"` }) (interface{}, error) {
			id, err := primitive.ObjectIDFromHex(arg.ID)
			if err != nil {
				return nil, err
			}

			rule, err := repo.Get(id)
			if err != nil {
				return nil, err
			}

			return rule, repo.DeleteID(id)
		},
	)
}

//
// func (r *RuleObject) DeleteRule(repo repository.RuleRepo) *graphql.Field {
// 	return &graphql.Field{
// 		Type: ruleType,
// 		Args: graphql.FieldConfigArgument{
// 			"id": &graphql.ArgumentConfig{Type: graphql.String,},
// 		},
// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			id, err := primitive.ObjectIDFromHex(p.Args["id"].(string))
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			rule, err := repo.Get(id)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			return rule, repo.DeleteID(id)
// 		},
// 	}
// }
