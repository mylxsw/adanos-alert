package jira

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/andygrunwald/go-jira"
)

// Issue 一个 Jira Issue
type Issue struct {
	CustomFields map[string]interface{} `json:"custom_fields"`
	ProjectKey   string                 `json:"project_key"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	IssueType    string                 `json:"issue_type"`
	Priority     string                 `json:"priority"`
	Assignee     string                 `json:"assignee"`
}

// Client 用于操作 jira 的客户端对象
type Client struct {
	client *jira.Client
}

// NewClient create a new jira client
func NewClient(baseURL string, username, password string) (*Client, error) {
	httpClient := &http.Client{}
	httpClient.Transport = &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
		if username != "" && password != "" {
			req.SetBasicAuth(username, password)
		}

		return nil, nil
	}}

	jiraClient, err := jira.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{client: jiraClient}, nil
}

// IssueResp 查询到的 Issue，附加状态
type IssueResp struct {
	Issue  Issue  `json:"issue"`
	Status string `json:"status"`
}

// GetIssue 获取一个 Issue
func (client Client) GetIssue(ctx context.Context, issueID string) (IssueResp, error) {
	issue, resp, err := client.client.Issue.GetWithContext(ctx, issueID, nil)
	if err != nil {
		return IssueResp{}, fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	return IssueResp{
		Issue: Issue{
			CustomFields: issue.Fields.Unknowns,
			ProjectKey:   issue.Fields.Project.Key,
			Summary:      issue.Fields.Summary,
			Description:  issue.Fields.Description,
			IssueType:    issue.Fields.Type.ID,
			Priority:     issue.Fields.Priority.ID,
			Assignee:     issue.Fields.Assignee.Name,
		},
		Status: issue.Fields.Status.Name,
	}, nil
}

// CreateIssue create a jira issue
func (client Client) CreateIssue(ctx context.Context, issue Issue) (string, error) {
	fields := jira.IssueFields{
		Type:        jira.IssueType{ID: issue.IssueType},
		Project:     jira.Project{Key: issue.ProjectKey},
		Summary:     issue.Summary,
		Description: issue.Description,
		Unknowns:    issue.CustomFields,
	}

	if issue.Assignee != "" {
		fields.Assignee = &jira.User{Name: issue.Assignee}
	}
	if issue.Priority != "" {
		fields.Priority = &jira.Priority{ID: issue.Priority}
	}

	createdIssue, resp, err := client.client.Issue.CreateWithContext(ctx, &jira.Issue{Fields: &fields})
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	return createdIssue.ID, nil
}

// UpdateIssue 更新 Issue 的自定义字段
func (client Client) UpdateIssue(ctx context.Context, issueID string, customFields map[string]interface{}) error {
	resp, err := client.client.Issue.UpdateIssueWithContext(ctx, issueID, map[string]interface{}{"fields": customFields})
	if err != nil {
		return fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	return nil
}

// CreateComment 创建一个评论
func (client Client) CreateComment(ctx context.Context, issueID string, comment string) error {
	_, resp, err := client.client.Issue.AddCommentWithContext(ctx, issueID, &jira.Comment{Body: comment})
	if err != nil {
		return fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	return nil
}

// IssueType is a jira issue type object
type IssueType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetIssueTypes return all issue types for a project
func (client Client) GetIssueTypes(ctx context.Context, projectKey string) ([]IssueType, error) {
	metas, resp, err := client.client.Issue.GetCreateMetaWithContext(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	issueTypes := make([]IssueType, 0)
	for _, m := range metas.GetProjectWithKey(projectKey).IssueTypes {
		issueTypes = append(issueTypes, IssueType{
			ID:   m.Id,
			Name: m.Name,
		})
	}

	return issueTypes, nil
}

// IssuePriority is a jira issue priority object
type IssuePriority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetPriorities return all priorities supported by jira
func (client Client) GetPriorities(ctx context.Context) ([]IssuePriority, error) {
	priorityList, resp, err := client.client.Priority.GetListWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	priorities := make([]IssuePriority, 0)
	for _, pr := range priorityList {
		priorities = append(priorities, IssuePriority{
			ID:   pr.ID,
			Name: pr.Name,
		})
	}

	return priorities, nil
}

// CustomField 自定义字段
type CustomField struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// GetCustomFields 获取所有的自定义字段
func (client Client) GetCustomFields(ctx context.Context) ([]CustomField, error) {
	fields, resp, err := client.client.Field.GetListWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, client.extractResponse(resp))
	}

	customFields := make([]CustomField, 0)
	for _, f := range fields {
		if f.Custom {
			customFields = append(customFields, CustomField{
				ID:   f.ID,
				Name: f.Name,
				Type: f.Schema.Type,
			})
		}
	}

	return customFields, nil
}

// extractResponse 解析服务端返回的响应内容
func (client Client) extractResponse(resp *jira.Response) string {
	defer func() {
		recover()
	}()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
