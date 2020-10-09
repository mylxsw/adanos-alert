package matcher_test

import (
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type triggerMatcherTestCase struct {
	Cond    string
	Matched bool
}

func TestTriggerMatcher(t *testing.T) {

	currentTs, _ := time.Parse(time.RFC3339, "2019-08-06T20:44:26+08:00")
	grp := repository.EventGroup{
		ID:           primitive.NewObjectID(),
		MessageCount: 10,
		CreatedAt:    currentTs,
	}
	triggerCtx := matcher.NewTriggerContext(container.New(), repository.Trigger{}, grp, func() []repository.Event {
		return []repository.Event{
			{
				Content: "Hello, world",
				Meta:    repository.EventMeta{"environment": "prod", "server": "192.168.1.2"},
				Tags:    []string{"php", "nodejs"},
				Origin:  "filebeat",
			},
			{
				Content: "Are you ready?",
				Meta:    repository.EventMeta{"environment": "prod", "server": "192.168.1.3"},
				Tags:    []string{"java", "nodejs"},
				Origin:  "elasticsearch",
			},
			{
				Content: "Nice day!",
				Meta:    repository.EventMeta{"environment": "prod", "server": "192.168.1.3"},
				Tags:    []string{"java"},
				Origin:  "elasticsearch",
			},
		}
	})

	var testcases = []triggerMatcherTestCase{
		{Cond: "Group.EventsCount > 9", Matched: true},
		{Cond: "Group.EventsCount > 10", Matched: false},
		{Cond: "Now().Sub(Group.CreatedAt).Minutes() > 10", Matched: true},
		{Cond: "ParseTime(\"2006-01-02T15:04:05Z07:00\", \"2019-08-06T20:00:00+08:00\").Before(Group.CreatedAt)", Matched: true},
		{Cond: "Group.CreatedAt.Hour() in 20..21", Matched: true},
		{Cond: "Group.CreatedAt.Hour() not in 9..18", Matched: true},
		{Cond: "len(filter(Events(), {#.Content matches 'ready'})) > 0", Matched: true},
		{Cond: `any(Events(), {"php" in #.Tags}) and none(Events(), {"swift" in #.Tags})`, Matched: true},
	}

	for _, ts := range testcases {
		mt, err := matcher.NewTriggerMatcher(repository.Trigger{PreCondition: ts.Cond})
		assert.NoError(t, err)

		matched, err := mt.Match(triggerCtx)
		assert.NoError(t, err)
		assert.Equal(t, ts.Matched, matched)

		assert.Equal(t, ts.Cond, mt.Trigger().PreCondition)
	}

	_, err := matcher.NewTriggerMatcher(repository.Trigger{PreCondition: "xxxxx"})
	assert.Error(t, err)
}
