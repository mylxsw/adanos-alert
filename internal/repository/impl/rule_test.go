package impl_test

import (
	"context"
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type RuleTestSuite struct {
	suite.Suite
	ruleRepo repository.RuleRepo
}

func (r *RuleTestSuite) TearDownTest() {
	r.NoError(r.ruleRepo.Delete(bson.M{}))
}

func (r *RuleTestSuite) SetupTest() {
	db, err := Database()
	r.NoError(err)

	r.ruleRepo = impl.NewRuleRepo(db)
}

func (r *RuleTestSuite) TestTags() {
	_, _ = r.ruleRepo.Add(repository.Rule{
		Name: "test1",
		Tags: []string{"mysql", "oracle", "nginx"},
	})
	_, _ = r.ruleRepo.Add(repository.Rule{
		Name: "test2",
		Tags: []string{"mysql"},
	})
	_, _ = r.ruleRepo.Add(repository.Rule{
		Name: "test3",
		Tags: nil,
	})
	_, _ = r.ruleRepo.Add(repository.Rule{
		Name: "test4",
		Tags: []string{"基础设施"},
	})
	_, _ = r.ruleRepo.Add(repository.Rule{
		Name: "test5",
		Tags: []string{"数据"},
	})

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	tags, err := r.ruleRepo.Tags(ctx)
	r.NoError(err)
	r.ElementsMatch(tags, []repository.Tag{
		{Name: "oracle", Count: 1},
		{Name: "mysql", Count: 2},
		{Name: "nginx", Count: 1},
		{Name: "数据", Count: 1},
		{Name: "基础设施", Count: 1},
	})
}

func TestRuleRepo(t *testing.T) {
	suite.Run(t, new(RuleTestSuite))
}
