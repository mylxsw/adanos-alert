package job_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository"
	mockAction "github.com/mylxsw/adanos-alert/test/mock/action"
	mockRepo "github.com/mylxsw/adanos-alert/test/mock/repository"
	"github.com/mylxsw/container"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerTestSuite struct {
	suite.Suite
	app *container.Container
}

func (t *TriggerTestSuite) SetupTest() {
	cc := container.New()
	cc.MustSingleton(mockRepo.NewMessageGroupRepo)
	cc.MustSingleton(mockRepo.NewRuleRepo)
	cc.MustSingleton(mockRepo.NewMessageRepo)

	// create action mock
	action.Register("http", mockAction.NewHttpAction())

	// prepare test data
	cc.MustResolve(t.createRules)
	cc.MustResolve(t.createMessageAndGroups)

	t.app = cc
}

func (t *TriggerTestSuite) TearDownTest() {
	t.app.MustResolve(func(groupRepo repository.MessageGroupRepo, messageRepo repository.MessageRepo, ruleRepo repository.RuleRepo) {
		t.NoError(groupRepo.Delete(bson.M{}))
		t.NoError(messageRepo.Delete(bson.M{}))
		t.NoError(ruleRepo.Delete(bson.M{}))
	})
}

func (t *TriggerTestSuite) createRules(ruleRepo repository.RuleRepo) {
	// add a test rule
	rul := t.createRule()
	_, err := ruleRepo.Add(rul)
	t.NoError(err)
}

func (t *TriggerTestSuite) createMessageAndGroups(groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo) {
	// TODO create test data
}

func (t *TriggerTestSuite) TestTriggerJobHandle() {
	t.app.MustResolve(func(groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo, ruleRepo repository.RuleRepo) {
		// mockGroupRepo := groupRepo.(*mockRepo.MessageGroupRepo)
		// mockMessageRepo := msgRepo.(*mockRepo.MessageRepo)
		// mockRuleRepo := ruleRepo.(*mockRepo.RuleRepo)

		job.NewTrigger(t.app).Handle()

		mockHttpAction := action.Factory("http").(*mockAction.HttpAction)
		fmt.Println(mockHttpAction.Histories)
		// TODO assert

	})
}

func (t *TriggerTestSuite) createRule() repository.Rule {
	return repository.Rule{
		Name:            "",
		Interval:        0,
		Threshold:       0,
		Rule:            `Upper(Meta["level"]) == "ERROR"`,
		SummaryTemplate: "",
		Triggers: []repository.Trigger{
			{
				PreCondition: "Group.MessageCount > 2",
				Action:       "http",
			},
			{
				PreCondition: `any(Messages(), {"php" in #.Tags})`,
				Action:       "http",
			},
			{
				PreCondition: `any(Messages(), {"swift" in #.Tags})`,
				Action:       "http",
			},
		},
		Status: repository.RuleStatusEnabled,
	}
}

func TestTriggerJob_Handle(t *testing.T) {
	suite.Run(t, new(TriggerTestSuite))
}
