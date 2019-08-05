package job_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository"
	mockRepo "github.com/mylxsw/adanos-alert/test/mock/repository"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type AggregationTestSuite struct {
	suite.Suite
	app *container.Container
}

func (a *AggregationTestSuite) SetupTest() {
	cc := container.New()
	cc.MustSingleton(mockRepo.NewMessageRepo)
	cc.MustSingleton(mockRepo.NewMessageGroupRepo)
	cc.MustSingleton(mockRepo.NewRuleRepo)

	a.app = cc
}

func (a *AggregationTestSuite) TearDownTest() {
	a.app.MustResolve(func(msgRepo repository.MessageRepo, msgGroupRepo repository.MessageGroupRepo, ruleRepo repository.RuleRepo) {
		a.NoError(msgRepo.Delete(bson.M{}))
		a.NoError(ruleRepo.Delete(bson.M{}))
		a.NoError(msgGroupRepo.Delete(bson.M{}))
	})
}

func (a *AggregationTestSuite) TestAggregationJob() {
	a.app.MustResolve(func(msgRepo repository.MessageRepo, msgGroupRepo repository.MessageGroupRepo, ruleRepo repository.RuleRepo) {
		mockMsgRepo := msgRepo.(*mockRepo.MessageRepo)
		mockMsgGroupRepo := msgGroupRepo.(*mockRepo.MessageGroupRepo)
		mockRuleRepo := ruleRepo.(*mockRepo.RuleRepo)

		// add a rule for test
		ruleID, err := mockRuleRepo.Add(repository.Rule{
			Name:   "test",
			Rule:   `"php" in Tags`,
			Status: repository.RuleStatusEnabled,
		})
		a.NoError(err)

		// add some messages 
		for i := 0; i < 10; i++ {
			_, err = mockMsgRepo.Add(repository.Message{
				Content: fmt.Sprintf("Hello, world #%d", i),
				Meta:    repository.MessageMeta{"environment": "dev", "seq": strconv.Itoa(i)},
				Tags:    []string{"php", "nodejs", fmt.Sprintf("tag_%d", i/2)},
				Origin:  "filebeat",
				Status:  repository.MessageStatusPending,
			})
			a.NoError(err)
		}

		// add some not matched messages
		for i := 0; i < 5; i++ {
			_, err = mockMsgRepo.Add(repository.Message{
				Content: fmt.Sprintf("Hello, world #%d", i),
				Meta:    repository.MessageMeta{"environment": "dev", "seq": strconv.Itoa(i + 10)},
				Tags:    []string{"java", "closure", fmt.Sprintf("tag_%d", i/2)},
				Origin:  "logstash",
				Status:  repository.MessageStatusPending,
			})
			a.NoError(err)
		}

		// execute job
		job.NewAggregationJob(a.app).Handle()

		// check message group
		a.EqualValues(1, len(mockMsgGroupRepo.Groups))
		a.EqualValues(ruleID, mockMsgGroupRepo.Groups[0].Rule.ID)
		a.Equal(repository.MessageGroupStatusCollecting, mockMsgGroupRepo.Groups[0].Status)

		// check message
		groupedMsgCount, err := mockMsgRepo.Count(bson.M{"status": repository.MessageStatusGrouped})
		a.NoError(err)
		a.EqualValues(10, groupedMsgCount)

		canceledMsgCount, err := mockMsgRepo.Count(bson.M{"status": repository.MessageStatusCanceled})
		a.NoError(err)
		a.EqualValues(5, canceledMsgCount)
	})
}

func TestAggregationJob_Handle(t *testing.T) {
	suite.Run(t, new(AggregationTestSuite))
}
