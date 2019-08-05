package impl_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
)

type SequenceTestSuite struct {
	suite.Suite
	repo repository.SequenceRepo
}

func (s *SequenceTestSuite) SetupTest() {
	db, err := Database()
	s.NoError(err)

	s.repo = impl.NewSequenceRepo(db)
}

func (s *SequenceTestSuite) TearDownTest() {
	s.NoError(s.repo.Truncate("test"))
	s.NoError(s.repo.Truncate("test2"))
}

func (s *SequenceTestSuite) TestSequence() {
	{
		seq, err := s.repo.Next("test")
		s.NoError(err)
		s.Equal(int64(1), seq.Value)
	}

	{
		seq, err := s.repo.Next("test")
		s.NoError(err)
		s.Equal(int64(2), seq.Value)
	}

	{
		seq, err := s.repo.Next("test2")
		s.NoError(err)
		s.Equal(int64(1), seq.Value)
	}

	{
		seq, err := s.repo.Next("test2")
		s.NoError(err)
		s.Equal(int64(2), seq.Value)
	}

	s.NoError(s.repo.Truncate("not_found"))
}

func TestSequenceRepo(t *testing.T) {
	suite.Run(t, new(SequenceTestSuite))
}
