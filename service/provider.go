package service

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct{}

func (p ServiceProvider) Register(app container.Container) {
	app.MustSingleton(NewEventService)
	app.MustSingleton(NewEventGroupService)
}

func (p ServiceProvider) Boot(app infra.Glacier) {
}
