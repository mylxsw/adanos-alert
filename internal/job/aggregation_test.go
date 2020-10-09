package job_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository"
	mockRepo "github.com/mylxsw/adanos-alert/test/mock/repository"
	"github.com/mylxsw/container"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type AggregationTestSuite struct {
	suite.Suite
	app container.Container
}

func (a *AggregationTestSuite) SetupTest() {
	cc := container.New()
	cc.MustSingleton(mockRepo.NewMessageRepo)
	cc.MustSingleton(mockRepo.NewMessageGroupRepo)
	cc.MustSingleton(mockRepo.NewRuleRepo)

	a.app = cc
}

func (a *AggregationTestSuite) TearDownTest() {
	a.app.MustResolve(func(msgRepo repository.EventRepo, msgGroupRepo repository.EventGroupRepo, ruleRepo repository.RuleRepo) {
		a.NoError(msgRepo.Delete(bson.M{}))
		a.NoError(ruleRepo.Delete(bson.M{}))
		a.NoError(msgGroupRepo.Delete(bson.M{}))
	})
}

func (a *AggregationTestSuite) TestAggregationJob() {
	a.app.MustResolve(func(msgRepo repository.EventRepo, msgGroupRepo repository.EventGroupRepo, ruleRepo repository.RuleRepo) {
		mockMsgRepo := msgRepo.(*mockRepo.MessageRepo)
		mockMsgGroupRepo := msgGroupRepo.(*mockRepo.EventGroupRepo)
		mockRuleRepo := ruleRepo.(*mockRepo.RuleRepo)

		// add a rule for test
		ruleID, err := mockRuleRepo.Add(repository.Rule{
			Name:     "test",
			Rule:     `"php" in Tags`,
			Interval: 30,
			Status:   repository.RuleStatusEnabled,
		})
		a.NoError(err)

		// add some messages 
		for i := 0; i < 10; i++ {
			_, err = mockMsgRepo.Add(repository.Event{
				Content: fmt.Sprintf("Hello, world #%d", i),
				Meta:    repository.EventMeta{"environment": "dev", "seq": strconv.Itoa(i)},
				Tags:    []string{"php", "nodejs", fmt.Sprintf("tag_%d", i/2)},
				Origin:  "filebeat",
				Status:  repository.EventStatusPending,
			})
			a.NoError(err)
		}

		// add some not matched messages
		for i := 0; i < 5; i++ {
			_, err = mockMsgRepo.Add(repository.Event{
				Content: fmt.Sprintf("Hello, world #%d", i),
				Meta:    repository.EventMeta{"environment": "dev", "seq": strconv.Itoa(i + 10)},
				Tags:    []string{"java", "closure", fmt.Sprintf("tag_%d", i/2)},
				Origin:  "logstash",
				Status:  repository.EventStatusPending,
			})
			a.NoError(err)
		}

		// execute job
		job.NewAggregationJob(a.app).Handle()

		// check message group
		a.EqualValues(1, len(mockMsgGroupRepo.Groups))
		a.EqualValues(ruleID, mockMsgGroupRepo.Groups[0].Rule.ID)
		a.Equal(repository.EventGroupStatusCollecting, mockMsgGroupRepo.Groups[0].Status)

		// check message
		groupedMsgCount, err := mockMsgRepo.Count(bson.M{"status": repository.EventStatusGrouped})
		a.NoError(err)
		a.EqualValues(10, groupedMsgCount)

		canceledMsgCount, err := mockMsgRepo.Count(bson.M{"status": repository.EventStatusCanceled})
		a.NoError(err)
		a.EqualValues(5, canceledMsgCount)

		// message grouping
		// change create_at -10s
		mockMsgGroupRepo.Groups[0].CreatedAt = mockMsgGroupRepo.Groups[0].CreatedAt.Add(-10 * time.Second)
		job.NewAggregationJob(a.app).Handle()
		a.Equal(repository.EventGroupStatusCollecting, mockMsgGroupRepo.Groups[0].Status)

		// change created_at -30s, reach grouping condition
		mockMsgGroupRepo.Groups[0].CreatedAt = mockMsgGroupRepo.Groups[0].CreatedAt.Add(-20 * time.Second)
		job.NewAggregationJob(a.app).Handle()
		a.Equal(repository.EventGroupStatusPending, mockMsgGroupRepo.Groups[0].Status)
	})
}

func TestAggregationJob_Handle(t *testing.T) {
	suite.Run(t, new(AggregationTestSuite))
}
