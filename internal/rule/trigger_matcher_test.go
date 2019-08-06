package rule_test

import (
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/rule"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type triggerMatcherTestCase struct {
	Cond    string
	Matched bool
}

func TestTriggerMatcher(t *testing.T) {

	currentTs, _ := time.Parse(time.RFC3339, "2019-08-06T20:44:26+08:00")
	grp := repository.MessageGroup{
		ID:           primitive.NewObjectID(),
		MessageCount: 10,
		CreatedAt:    currentTs,
	}
	triggerCtx := rule.NewTriggerContext(grp, func() []repository.Message {
		return []repository.Message{
			{
				Content: "Hello, world",
				Meta:    repository.MessageMeta{"environment": "prod", "server": "192.168.1.2"},
				Tags:    []string{"php", "nodejs"},
				Origin:  "filebeat",
			},
			{
				Content: "Are you ready?",
				Meta:    repository.MessageMeta{"environment": "prod", "server": "192.168.1.3"},
				Tags:    []string{"java", "nodejs"},
				Origin:  "elasticsearch",
			},
			{
				Content: "Nice day!",
				Meta:    repository.MessageMeta{"environment": "prod", "server": "192.168.1.3"},
				Tags:    []string{"java"},
				Origin:  "elasticsearch",
			},
		}
	})

	var testcases = []triggerMatcherTestCase{
		{Cond: "Group.MessageCount > 9", Matched: true},
		{Cond: "Group.MessageCount > 10", Matched: false},
		{Cond: "Now().Sub(Group.CreatedAt).Minutes() > 10", Matched: true},
		{Cond: "ParseTime(\"2006-01-02T15:04:05Z07:00\", \"2019-08-06T20:00:00+08:00\").Before(Group.CreatedAt)", Matched: true},
		{Cond: "Group.CreatedAt.Hour() in 20..21", Matched: true},
		{Cond: "Group.CreatedAt.Hour() not in 9..18", Matched: true},
		{Cond: "len(filter(Messages(), {#.Content matches 'ready'})) > 0", Matched: true},
		{Cond: `any(Messages(), {"php" in #.Tags}) and none(Messages(), {"swift" in #.Tags})`, Matched: true},
	}

	for _, ts := range testcases {
		matcher, err := rule.NewTriggerMatcher(repository.Trigger{PreCondition: ts.Cond,})
		assert.NoError(t, err)

		matched, err := matcher.Match(triggerCtx)
		assert.NoError(t, err)
		assert.Equal(t, ts.Matched, matched)

		assert.Equal(t, ts.Cond, matcher.Trigger().PreCondition)
	}

	_, err := rule.NewTriggerMatcher(repository.Trigger{PreCondition: "xxxxx"})
	assert.Error(t, err)
}
