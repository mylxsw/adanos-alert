package view

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
)

type GroupData struct {
	Group       repository.EventGroup
	Events      []repository.Event
	EventsCount int64
	Next        int64
	Offset      int64
	Limit       int64
	Path        string
	HasPrev     bool
	HasNext     bool
	PrevOffset  int64
}

var defaultTemplateContent string

func init() {
	defaultTemplateContent = FSMustString(false, "/groups.html")
}

// GroupView 分组视图展示
func GroupView(cc infra.Resolver, data GroupData) (string, error) {
	templateContent, err := fileGetContent(filepath.Join(currentPath(), "groups.html"))
	if err != nil {
		templateContent = defaultTemplateContent
	}

	return template.Parse(cc, strings.ReplaceAll(defaultLayout, "{{--BODY--}}", templateContent), data)
}

// ReportView 报表视图展示
func ReportView(cc infra.Resolver, templateContent string, data GroupData) (string, error) {
	if templateContent == "" {
		return GroupView(cc, data)
	}

	parsed, err := template.Parse(cc, strings.ReplaceAll(defaultLayout, "{{--BODY--}}", templateContent), data)
	if err != nil {
		log.WithFields(log.Fields{
			"template": templateContent,
			"data":     data,
		}).Errorf("parse report template failed: %v", err)
		return GroupView(cc, data)
	}

	return parsed, nil
}

// fileGetContent 读取文件内容
func fileGetContent(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// currentPath 获取当前工作目录
func currentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	return dir
}
