package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/messager/jira"
	"github.com/mylxsw/asteria/log"
)

type JiraAction struct {
	manager Manager
}

func (act JiraAction) Validate(meta string, userRefs []string) error {
	var jiraMeta JiraMeta
	if err := json.Unmarshal([]byte(meta), &jiraMeta); err != nil {
		return err
	}

	if len(userRefs) > 2 {
		return errors.New("invalid users, only support one user")
	}

	return nil
}

func NewJiraAction(manager Manager) *JiraAction {
	return &JiraAction{manager: manager}
}

type JiraMeta struct {
	Issue jira.Issue `json:"issue"`
}

func (act JiraAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	var meta JiraMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse jira meta failed: %v", err)
	}

	return act.manager.Resolve(func(conf *configs.Config, evtRepo repository.EventRepo, userRepo repository.UserRepo) error {
		jiraClient, err := jira.NewClient(conf.Jira.BaseURL, conf.Jira.Username, conf.Jira.Password)
		if err != nil {
			return fmt.Errorf("create jira client failed: %w", err)
		}

		payload, description := createPayloadAndSummary(act.manager, "jira", conf, evtRepo, rule, trigger, grp)
		if meta.Issue.Description != "" {
			description = parseTemplate(act.manager, meta.Issue.Description, payload)
		}

		summary := rule.Name
		if meta.Issue.Summary != "" {
			summary = parseTemplate(act.manager, meta.Issue.Summary, payload)
		}

		customFields := make(map[string]interface{})
		for k, v := range meta.Issue.CustomFields {
			customFields[k] = parseTemplate(act.manager, fmt.Sprintf("%v", v), payload)
		}

		issue := jira.Issue{
			CustomFields: customFields,
			ProjectKey:   meta.Issue.ProjectKey,
			Summary:      summary,
			Description:  description,
			IssueType:    meta.Issue.IssueType,
			Priority:     meta.Issue.Priority,
		}

		if len(trigger.UserRefs) > 0 && !trigger.UserRefs[0].IsZero() {
			user, err := userRepo.Get(trigger.UserRefs[0])
			if err != nil {
				log.WithFields(log.Fields{
					"user_id":    trigger.UserRefs[0].Hex(),
					"trigger_id": trigger.ID.Hex(),
					"rule_id":    rule.ID.Hex(),
				}).Errorf("no such user")
			} else {
				jiraUser := user.Metas.Get("jira")
				if jiraUser != "" {
					issue.Assignee = jiraUser
				}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		issueID, err := jiraClient.CreateIssue(ctx, issue)
		if err != nil {
			log.WithFields(log.Fields{
				"title":       rule.Name,
				"description": description,
				"err":         err,
				"meta":        meta,
			}).Errorf("send message to jira failed: %v", err)
			return err
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"title":       rule.Name,
				"description": description,
				"meta":        meta,
				"issue_id":    issueID,
			}).Debug("send message to jira succeed")
		}

		return nil
	})
}
