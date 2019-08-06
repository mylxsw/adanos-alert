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

type RuleRepoTestSuite struct {
	suite.Suite
	repo repository.RuleRepo
}

func (r *RuleRepoTestSuite) TearDownTest() {
	r.NoError(r.repo.Delete(bson.M{}))
}

func (r *RuleRepoTestSuite) SetupTest() {
	db, err := Database()
	r.NoError(err)

	r.repo = impl.NewRuleRepo(db)
}

func (r *RuleRepoTestSuite) TestRuleRepo() {
	var lastId primitive.ObjectID
	for i := 0; i < 10; i++ {
		id, err := r.repo.Add(repository.Rule{
			Name:        fmt.Sprintf("Test #%d", i),
			Description: fmt.Sprintf("This is a test rule #%d", i),
			Interval:    60,
			Threshold:   0,
			Priority:    int64(i * 10),
			Status:      repository.RuleStatusEnabled,
			Triggers: []repository.Trigger{
				{
					PreCondition: "a == b",
					Action:       "email",
					Meta:         "",
				},
			},
		})
		r.NoError(err)
		r.NotEmpty(id.String())

		lastId = id
	}

	rule, err := r.repo.Get(lastId)
	r.NoError(err)
	r.Equal("Test #9", rule.Name)
	r.NotEmpty(rule.Triggers[0].ID.String())

	rule.Threshold = 100
	rule.Triggers = append(rule.Triggers, repository.Trigger{
		PreCondition: "c == d",
		Action:       "dingding",
	})
	r.NoError(r.repo.UpdateID(lastId, rule))

	rule, err = r.repo.Get(lastId)
	r.NoError(err)
	r.EqualValues(100, rule.Threshold)
	r.EqualValues(2, len(rule.Triggers))
	r.NotEmpty(rule.Triggers[1].ID.String())

	rules, err := r.repo.Find(bson.M{})
	r.NoError(err)
	r.EqualValues(10, len(rules))
	r.Equal("Test #9", rules[0].Name)
	r.Equal("Test #2", rules[7].Name)
	r.Equal("Test #0", rules[9].Name)

	var matchedCount = 0
	r.NoError(r.repo.Traverse(bson.M{"threshold": bson.M{"$lt": 10}}, func(rule repository.Rule) error {
		matchedCount++
		r.Equal(fmt.Sprintf("Test #%d", 9-matchedCount), rule.Name)

		return nil
	}))
	r.EqualValues(9, matchedCount)

	// Delete two rules
	r.NoError(r.repo.Delete(bson.M{"name": "Test #5"}))
	r.NoError(r.repo.Delete(bson.M{"name": "Test #3"}))
	// delete a not exist rule
	r.NoError(r.repo.Delete(bson.M{"name": "Test #15"}))
	// delete by id
	r.NoError(r.repo.DeleteID(lastId))

	resCount, err := r.repo.Count(bson.M{})
	r.NoError(err)
	r.EqualValues(7, resCount)
}

func TestRuleRepo(t *testing.T) {
	suite.Run(t, new(RuleRepoTestSuite))
}
