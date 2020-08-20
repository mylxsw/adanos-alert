package store

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(NewMessageStore)
}

func (s ServiceProvider) Boot(app infra.Glacier) {

}