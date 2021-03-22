package service

import (
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (p Provider) Register(app infra.Binder) {
	app.MustSingleton(NewEventService)
	app.MustSingleton(NewEventGroupService)
}

func (p Provider) Boot(app infra.Resolver) {
}
