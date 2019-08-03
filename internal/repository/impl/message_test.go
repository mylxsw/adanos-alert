package impl_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageTestSuit struct {
	suite.Suite
	repo repository.MessageRepo
}

func (s *MessageTestSuit) SetupTest() {
	db, err := Database()
	s.NoError(err)

	s.repo = impl.NewMessageRepo(db)
}

func (s *MessageTestSuit) TearDownTest() {
	s.NoError(s.repo.Delete(bson.M{}))
}

func (s *MessageTestSuit) TestMessageCURD() {
	msg := repository.Message{
		Content: "message content",
		Tags: []repository.MessageTag{
			{Key: "level", Value: "error"},
			{Key: "biz_code", Value: "laravel"},
		},
		Origin: "elasticsearch",
	}

	id, err := s.repo.Add(msg)
	s.NoError(err)
	s.NotEmpty(id.String())

	m, err := s.repo.Get(id)
	s.NoError(err)
	s.Equal(msg.Content, m.Content)
	s.NotEmpty(m.CreatedAt)
	s.Equal(2, len(m.Tags))

	for i := 0; i < 100; i++ {
		msg.Content = fmt.Sprintf("new message content %d", i)
		msg.Tags = append(msg.Tags, repository.MessageTag{Key: "filename", Value: "/var/log/message"})

		id, err := s.repo.Add(msg)
		s.NoError(err)
		s.NotEmpty(id)
	}

	count, err := s.repo.Count(bson.M{})
	s.NoError(err)
	s.EqualValues(101, count)

	groupId := primitive.NewObjectID()
	s.NoError(s.repo.Traverse(bson.M{"tags.key": "filename", "tags.value": "/var/log/message"}, func(msg repository.Message) error {
		msg.GroupID = groupId
		return s.repo.UpdateID(msg.ID, msg)
	}))

	msgs, err := s.repo.Find(bson.M{"group_id": primitive.NilObjectID})
	s.NoError(err)
	s.EqualValues(1, len(msgs))

	msgs, next, err := s.repo.Paginate(bson.M{}, 0, 10)
	s.NoError(err)
	s.EqualValues(10, len(msgs))
	s.EqualValues(10, next)

	msgs, next, err = s.repo.Paginate(bson.M{}, next, 90)
	s.NoError(err)
	s.EqualValues(90, len(msgs))
	s.EqualValues(100, next)

	msgs, next, err = s.repo.Paginate(bson.M{}, next, 10)
	s.NoError(err)
	s.EqualValues(1, len(msgs))
	s.EqualValues(0, next)
}

func TestMessageRepo(t *testing.T) {
	suite.Run(t, new(MessageTestSuit))
}
