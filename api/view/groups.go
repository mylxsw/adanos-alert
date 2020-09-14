package view

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/container"
)

type GroupData struct {
	Group        repository.MessageGroup
	Messages     []repository.Message
	MessageCount int64
	Next         int64
	Offset       int64
	Limit        int64
	Path         string
	HasPrev      bool
	HasNext      bool
	PrevOffset   int64
}

var defaultTemplateContent string

func init() {
	defaultTemplateContent = FSMustString(false, "/groups.html")
}

func GroupView(cc container.Container, data GroupData) (string, error) {
	templateContent, err := fileGetContent(filepath.Join(currentPath(), "groups.html"))
	if err != nil {
		templateContent = defaultTemplateContent
	}

	return template.Parse(cc, templateContent, data)
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
