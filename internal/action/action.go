package action

import (
	"sync"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
)

type Action interface {
	Handle(trigger repository.Trigger) error
}

var actions = make(map[string]Action)
var lock sync.RWMutex

// Factory get a action by name
func Factory(action string) Action {
	lock.RLock()
	defer lock.RUnlock()

	return actions[action]
}

// Register register a new action
func Register(name string, action Action) {
	lock.Lock()
	defer lock.Unlock()

	actions[name] = action
}

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	Register("http", NewHttpAction())
	Register("dingding", NewDingdingAction())
	Register("email", NewEmailAction())
	Register("wechat", NewWechatAction())
}
