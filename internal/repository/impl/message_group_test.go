package impl_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageGroupRepoTestSuit struct {
	suite.Suite
	repo repository.MessageGroupRepo
}

func (m *MessageGroupRepoTestSuit) TearDownTest() {
	m.NoError(m.repo.Delete(bson.M{}))
}

func (m *MessageGroupRepoTestSuit) SetupTest() {
	db, err := Database()
	m.NoError(err)

	m.repo = impl.NewMessageGroupRepo(db)
}

func (m *MessageGroupRepoTestSuit) TestMessageGroup() {
	grp := repository.MessageGroup{
		Status: repository.MessageGroupStatusCollecting,
	}

	id, err := m.repo.Add(grp)
	m.NoError(err)
	m.NotEmpty(id.String())

	grp2, err := m.repo.Get(id)
	m.NoError(err)
	m.Equal(grp.Status, grp2.Status)
	m.NotEmpty(grp2.CreatedAt)
	m.NotEmpty(grp2.UpdatedAt)

	grp.Status = repository.MessageGroupStatusPending
	id2, err := m.repo.Add(grp)
	m.NoError(err)
	m.NotEmpty(id2.String())
	m.NotEqual(id.String(), id2.String())

	for i := 0; i < 10; i++ {
		grp.Status = repository.MessageGroupStatusOK
		id3, err := m.repo.Add(grp)
		m.NoError(err)
		m.NotEmpty(id3.String())
		m.NotEqual(id.String(), id3.String())
	}

	m.NoError(m.repo.DeleteID(id2))

	count, err := m.repo.Count(bson.M{})
	m.NoError(err)
	m.EqualValues(11, count)

	grps, err := m.repo.Find(bson.M{"status": repository.MessageGroupStatusOK})
	m.NoError(err)
	m.EqualValues(10, len(grps))

	grps, next, err := m.repo.Paginate(bson.M{"status": repository.MessageGroupStatusOK}, 0, 10)
	m.NoError(err)
	m.EqualValues(10, len(grps))
	m.EqualValues(10, next)

	grps, next, err = m.repo.Paginate(bson.M{"status": repository.MessageGroupStatusOK}, next, 10)
	m.NoError(err)
	m.Empty(grps)
	m.EqualValues(0, next)

	m.NoError(m.repo.Traverse(bson.M{"status": repository.MessageGroupStatusOK}, func(grp repository.MessageGroup) error {
		grp.Status = repository.MessageGroupStatusFailed
		return m.repo.Update(grp.ID, grp)
	}))

	count, err = m.repo.Count(bson.M{"status": repository.MessageGroupStatusFailed})
	m.NoError(err)
	m.EqualValues(10, count)
}

func TestMessageGroupRepo(t *testing.T) {
	suite.Run(t, new(MessageGroupRepoTestSuit))
}
