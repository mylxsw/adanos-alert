package view

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Trigger is a action trigger for matched rules
type Trigger struct {
	ID           string                   `json:"id"`
	PreCondition string                   `json:"pre_condition"`
	Action       string                   `json:"action"`
	Meta         string                   `json:"meta"`
	Status       repository.TriggerStatus `json:"status"`
	FailedCount  int                      `json:"failed_count"`
	FailedReason string                   `json:"failed_reason"`
}

// Rule is a rule definition
type Rule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Interval  int64 `json:"interval"`
	Threshold int64 `json:"threshold"`
	Priority  int64 `json:"priority"`

	Rule            string     `json:"rule"`
	Template        string     `json:"template"`
	SummaryTemplate string     `json:"summary_template"`
	Triggers        []*Trigger `json:"triggers"`

	Status repository.RuleStatus `json:"status"`

	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	CreatedTime time.Time `json:"created_time"`
}

func (r Rule) TriggersResolver() ([]Trigger, error) {
	return []Trigger{
		{
			ID:           "000001",
			PreCondition: "a == b",
		},
		{
			ID:           "000002",
			PreCondition: "a == b and c == d",
		},
	}, nil
}

func (r Rule) CreatedTimeResolver() (string, error) {
	log.Info("query created_time")
	return r.CreatedTime.Format(time.RFC3339), nil
}

func RuleFromRepo(r repository.Rule) Rule {
	return Rule{
		ID:              r.ID.Hex(),
		Name:            r.Name,
		Description:     r.Description,
		Interval:        r.Interval,
		Threshold:       r.Threshold,
		Priority:        r.Priority,
		Rule:            r.Rule,
		Template:        r.Template,
		SummaryTemplate: r.SummaryTemplate,
		Triggers:        TriggersFromRepos(r.Triggers),
		Status:          r.Status,
		CreatedAt:       r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC3339),
		CreatedTime:     r.CreatedAt,
	}
}

func RulesFromRepos(rs []repository.Rule) []Rule {
	var rules = make([]Rule, len(rs))
	for i, r := range rs {
		rules[i] = RuleFromRepo(r)
	}

	return rules
}

func RuleToRepo(r Rule) repository.Rule {
	id, _ := primitive.ObjectIDFromHex(r.ID)
	createdAt, _ := time.Parse(time.RFC3339, r.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, r.UpdatedAt)
	return repository.Rule{
		ID:              id,
		Name:            r.Name,
		Description:     r.Description,
		Interval:        r.Interval,
		Threshold:       r.Threshold,
		Priority:        r.Priority,
		Rule:            r.Rule,
		Template:        r.Template,
		SummaryTemplate: r.SummaryTemplate,
		Triggers:        TriggersToRepo(r.Triggers),
		Status:          r.Status,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func TriggerFromRepo(tr repository.Trigger) *Trigger {
	return &Trigger{
		ID:           tr.ID.Hex(),
		PreCondition: tr.PreCondition,
		Action:       tr.Action,
		Meta:         tr.Meta,
		Status:       tr.Status,
		FailedCount:  tr.FailedCount,
		FailedReason: tr.FailedReason,
	}
}

func TriggersFromRepos(trs []repository.Trigger) []*Trigger {
	triggers := make([]*Trigger, len(trs))
	for i, tr := range trs {
		triggers[i] = TriggerFromRepo(tr)
	}

	return triggers
}

func TriggerToRepo(tr Trigger) repository.Trigger {
	id, _ := primitive.ObjectIDFromHex(tr.ID)
	return repository.Trigger{
		ID:           id,
		PreCondition: tr.PreCondition,
		Action:       tr.Action,
		Meta:         tr.Meta,
		Status:       tr.Status,
		FailedCount:  tr.FailedCount,
		FailedReason: tr.FailedReason,
	}
}

func TriggersToRepo(trs []*Trigger) []repository.Trigger {
	var triggers = make([]repository.Trigger, len(trs))
	for i, t := range trs {
		triggers[i] = TriggerToRepo(*t)
	}

	return triggers
}
