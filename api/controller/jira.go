package controller

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/pkg/messager/jira"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/go-utils/str"
)

type JiraController struct {
	cc container.Container
}

func NewJiraController(cc container.Container) web.Controller {
	return &JiraController{cc: cc}
}

func (j JiraController) Register(router *web.Router) {
	router.Group("/jira/issue", func(router *web.Router) {
		router.Get("/options/", j.IssueOptions).Name("jira.issue:options")
		router.Get("/types/", j.IssueTypes).Name("jira:issue:types")
		router.Get("/priorities/", j.Priorities).Name("jira:issue:priorities")
		router.Get("/custom-fields/", j.CustomFields).Name("jira:issue:custom-fields")
	})
}

func (j JiraController) IssueOptions(webCtx web.Context, conf *configs.Config) web.Response {
	res := web.M{"priorities": nil, "issue_types": nil}
	var lock sync.Mutex

	jiraClient, err := jira.NewClient(conf.Jira.BaseURL, conf.Jira.Username, conf.Jira.Password)
	if err != nil {
		log.Errorf("create jira client failed: %v", err)
		return webCtx.JSON(res)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
		defer cancel()

		priorities, err := jiraClient.GetPriorities(ctx)
		if err != nil {
			log.Errorf("get jira priorities failed: %v", err)
		} else {
			lock.Lock()
			res["priorities"] = priorities
			lock.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		projectKey := strings.TrimSpace(webCtx.Input("project_key"))
		if projectKey != "" {
			ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
			defer cancel()

			issueTypes, err := jiraClient.GetIssueTypes(ctx, projectKey)
			if err != nil {
				log.Errorf("get jira issue types failed: %v", err)
			} else {
				lock.Lock()
				res["issue_types"] = issueTypes
				lock.Unlock()
			}
		}
	}()

	wg.Wait()

	return webCtx.JSON(res)
}

func (j JiraController) Priorities(webCtx web.Context, conf *configs.Config) web.Response {
	jiraClient, err := jira.NewClient(conf.Jira.BaseURL, conf.Jira.Username, conf.Jira.Password)
	if err != nil {
		log.Errorf("create jira client failed: %v", err)
		return webCtx.JSON(web.M{"priorities": nil})
	}

	ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
	defer cancel()

	priorities, err := jiraClient.GetPriorities(ctx)
	if err != nil {
		log.Errorf("get jira priorities failed: %v", err)
		return webCtx.JSON(web.M{"priorities": nil})
	}

	return webCtx.JSON(web.M{
		"priorities": priorities,
	})
}

func (j JiraController) IssueTypes(webCtx web.Context, conf *configs.Config) web.Response {
	projectKey := strings.TrimSpace(webCtx.Input("project_key"))
	if conf.Jira.BaseURL == "" || projectKey == "" {
		return webCtx.JSON(web.M{"issue_types": nil})
	}

	jiraClient, err := jira.NewClient(conf.Jira.BaseURL, conf.Jira.Username, conf.Jira.Password)
	if err != nil {
		log.Errorf("create jira client failed: %v", err)
		return webCtx.JSON(web.M{"issue_types": nil})
	}

	ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
	defer cancel()

	issueTypes, err := jiraClient.GetIssueTypes(ctx, projectKey)
	if err != nil {
		log.Errorf("get jira issue types failed: %v", err)
		return webCtx.JSON(web.M{"issue_types": nil})
	}

	return webCtx.JSON(web.M{
		"issue_types": issueTypes,
	})
}

func (j JiraController) CustomFields(webCtx web.Context, conf *configs.Config) web.Response {
	if conf.Jira.BaseURL == "" {
		return webCtx.JSON(web.M{"fields": nil})
	}

	jiraClient, err := jira.NewClient(conf.Jira.BaseURL, conf.Jira.Username, conf.Jira.Password)
	if err != nil {
		log.Errorf("create jira client failed: %v", err)
		return webCtx.JSON(web.M{"issue_types": nil})
	}

	ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
	defer cancel()

	fields, err := jiraClient.GetCustomFields(ctx)
	if err != nil {
		log.Errorf("get jira issue custom fields failed: %v", err)
		return webCtx.JSON(web.M{"fields": nil})
	}

	_ = coll.Filter(fields, &fields, func(cf jira.CustomField) bool {
		return str.In(cf.Type, []string{"string", "number", "datetime"})
	})

	return webCtx.JSON(web.M{"fields": fields})
}
