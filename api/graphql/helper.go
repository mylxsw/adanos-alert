package graphql

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
)

func UserToRepo(original repository.User, user NewUser) repository.User {
	var metas []repository.UserMeta
	_ = coll.Map(user.Metas, &metas, func(meta *NewUserMeta) repository.UserMeta {
		return repository.UserMeta{
			Key:   meta.Key,
			Value: meta.Value,
		}
	})

	original.Name = user.Name
	original.Metas = metas
	original.Status = repository.UserStatus(user.Status)

	return original
}

func RepoToUser(user repository.User) *User {
	var metas []*UserMeta
	_ = coll.Map(user.Metas, &metas, func(meta repository.UserMeta) *UserMeta {
		return &UserMeta{
			Key:   meta.Key,
			Value: meta.Value,
		}
	})
	return &User{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Metas:     metas,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func RuleToRepo(rule NewRule) repository.Rule {
	var triggers []repository.Trigger
	_ = coll.Map(rule.Triggers, &triggers, func(r *NewTrigger) repository.Trigger {
		return repository.Trigger{
			PreCondition: r.PreCondition,
			Action:       r.Action,
		}
	})

	return repository.Rule{
		Name:            rule.Name,
		Description:     nilString(rule.Description),
		Interval:        int64(rule.Interval),
		Threshold:       int64(rule.Threshold),
		Priority:        int64(rule.Priority),
		Rule:            rule.Rule,
		Template:        nilString(rule.Template),
		SummaryTemplate: nilString(rule.SummaryTemplate),
		Triggers:        triggers,
		Status:          repository.RuleStatus(nilString(rule.Status)),
	}
}

func RepoToRule(r repository.Rule) *Rule {
	var triggers []*Trigger
	if len(r.Triggers) > 0 {
		_ = coll.Map(r.Triggers, &triggers, func(tr repository.Trigger) *Trigger {
			return &Trigger{
				ID:           tr.ID.Hex(),
				PreCondition: tr.PreCondition,
				Action:       tr.Action,
			}
		})
	}

	return &Rule{
		ID:              r.ID.Hex(),
		Name:            r.Name,
		Description:     r.Description,
		Interval:        int(r.Interval),
		Threshold:       int(r.Threshold),
		Priority:        int(r.Priority),
		Rule:            r.Rule,
		Template:        r.Template,
		SummaryTemplate: r.SummaryTemplate,
		Triggers:        triggers,
		Status:          string(r.Status),
		CreatedAt:       r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC3339),
	}
}

func nilString(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}
