package impl_test

import (
	"context"
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageGroupTestSuite struct {
	suite.Suite
	repo    repository.EventGroupRepo
	seqRepo repository.SequenceRepo
}

func (m *MessageGroupTestSuite) TearDownTest() {
	m.NoError(m.repo.Delete(bson.M{}))
	m.NoError(m.seqRepo.Truncate("group_seq"))
}

func (m *MessageGroupTestSuite) SetupTest() {
	db, err := Database()
	m.NoError(err)

	m.seqRepo = impl.NewSequenceRepo(db)
	m.repo = impl.NewEventGroupRepo(db, m.seqRepo)
}

func (m *MessageGroupTestSuite) TestMessageGroup() {
	grp := repository.EventGroup{
		Status: repository.EventGroupStatusCollecting,
	}

	id, err := m.repo.Add(grp)
	m.NoError(err)
	m.NotEmpty(id.String())

	grp2, err := m.repo.Get(id)
	m.NoError(err)
	m.Equal(grp.Status, grp2.Status)
	m.NotEmpty(grp2.CreatedAt)
	m.NotEmpty(grp2.UpdatedAt)

	_, err = m.repo.Get(primitive.NewObjectID())
	m.Error(err)
	m.Equal(repository.ErrNotFound, err)

	grp.Status = repository.EventGroupStatusPending
	id2, err := m.repo.Add(grp)
	m.NoError(err)
	m.NotEmpty(id2.String())
	m.NotEqual(id.String(), id2.String())

	for i := 0; i < 10; i++ {
		grp.Status = repository.EventGroupStatusOK
		id3, err := m.repo.Add(grp)
		m.NoError(err)
		m.NotEmpty(id3.String())
		m.NotEqual(id.String(), id3.String())
	}

	m.NoError(m.repo.DeleteID(id2))

	count, err := m.repo.Count(bson.M{})
	m.NoError(err)
	m.EqualValues(11, count)

	grps, err := m.repo.Find(bson.M{"status": repository.EventGroupStatusOK})
	m.NoError(err)
	m.EqualValues(10, len(grps))

	grps, next, err := m.repo.Paginate(bson.M{"status": repository.EventGroupStatusOK}, 0, 10, false)
	m.NoError(err)
	m.EqualValues(10, len(grps))
	m.EqualValues(10, next)

	grps, next, err = m.repo.Paginate(bson.M{"status": repository.EventGroupStatusOK}, next, 10, false)
	m.NoError(err)
	m.Empty(grps)
	m.EqualValues(0, next)

	m.NoError(m.repo.Traverse(bson.M{"status": repository.EventGroupStatusOK}, func(grp repository.EventGroup) error {
		grp.Status = repository.EventGroupStatusFailed
		return m.repo.UpdateID(grp.ID, grp)
	}))

	count, err = m.repo.Count(bson.M{"status": repository.EventGroupStatusFailed})
	m.NoError(err)
	m.EqualValues(10, count)

	// Test collecting group
	rule := repository.Rule{
		ID:          primitive.NewObjectID(),
		Name:        "test",
		Description: "test rule",
		Interval:    30,
		Rule:        `"php" in Tags`,
		Status:      repository.RuleStatusEnabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	groupRule := rule.ToGroupRule("", repository.EventTypePlain)
	m.Equal(rule.ID, groupRule.ID)
	m.Equal(rule.Name, groupRule.Name)
	m.Equal(rule.Rule, groupRule.Rule)

	collectingGroup, err := m.repo.CollectingGroup(groupRule)
	m.NoError(err)
	m.Equal(repository.EventGroupStatusCollecting, collectingGroup.Status)
	m.Equal(groupRule.Rule, collectingGroup.Rule.Rule)
	m.NotEmpty(collectingGroup.CreatedAt)

	collectingGroup2, err := m.repo.CollectingGroup(groupRule)
	m.NoError(err)
	m.Equal(collectingGroup.ID, collectingGroup2.ID)
	m.EqualValues(collectingGroup.CreatedAt.Unix(), collectingGroup2.CreatedAt.Unix())

	ruleCount, err := m.repo.StatByRuleCount(context.TODO(), time.Now().Add(-365*24*time.Hour), time.Now())
	m.NoError(err)
	m.NotEmpty(ruleCount)

	_, err = m.repo.StatByUserCount(context.TODO(), time.Now().Add(-365*24*time.Hour), time.Now())
	m.NoError(err)

	res, err := m.repo.StatByDatetimeCount(context.TODO(), nil, time.Now().Add(-365*24*time.Hour), time.Now(), 1)
	m.NoError(err)
	m.NotEmpty(res)
}

func TestMessageGroupRepo(t *testing.T) {
	suite.Run(t, new(MessageGroupTestSuite))
}
