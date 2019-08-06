package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
)

type Action interface {
	Handle(trigger repository.Trigger) error
}

var actions = make(map[string]Action)

func Factory(action string) Action {
	return actions[action]
}

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	actions["http"] = NewHttpAction()
	actions["dingding"] = NewDingdingAction()
	actions["email"] = NewEmailAction()
	actions["wechat"] = NewWechatAction()
}
