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
	repo    repository.MessageRepo
	seqRepo repository.SequenceRepo
}

func (s *MessageTestSuit) SetupTest() {
	db, err := Database()
	s.NoError(err)

	s.seqRepo = impl.NewSequenceRepo(db)
	s.repo = impl.NewMessageRepo(db, s.seqRepo)
}

func (s *MessageTestSuit) TearDownTest() {
	s.NoError(s.repo.Delete(bson.M{}))
	s.NoError(s.seqRepo.Truncate("message_seq"))
}

func (s *MessageTestSuit) TestMessageCURD() {
	msg := repository.Message{
		Content: "message content",
		Meta: repository.MessageMeta{
			"level":       "error",
			"environment": "dev",
		},
		Tags:    []string{"test", "test2"},
		Origin:  "elasticsearch",
		GroupID: make([]primitive.ObjectID, 0),
	}

	id, err := s.repo.Add(msg)
	s.NoError(err)
	s.NotEmpty(id.String())

	m, err := s.repo.Get(id)
	s.NoError(err)
	s.Equal(msg.Content, m.Content)
	s.NotEmpty(m.CreatedAt)
	s.Equal(2, len(m.Meta))

	for i := 0; i < 100; i++ {
		msg.Content = fmt.Sprintf("new message content %d", i)
		msg.Meta["filename"] = "/var/log/message"

		id, err := s.repo.Add(msg)
		s.NoError(err)
		s.NotEmpty(id)
	}

	count, err := s.repo.Count(bson.M{})
	s.NoError(err)
	s.EqualValues(101, count)

	groupId := primitive.NewObjectID()
	s.NoError(s.repo.Traverse(bson.M{"meta.filename": "/var/log/message"}, func(msg repository.Message) error {
		msg.GroupID = []primitive.ObjectID{groupId}
		return s.repo.UpdateID(msg.ID, msg)
	}))

	msgs, err := s.repo.Find(bson.D{{"$or", bson.A{bson.M{"group_ids": nil}, bson.M{"group_ids": bson.M{"$size": 0}}}}})
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
