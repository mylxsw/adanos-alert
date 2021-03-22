package store

import (
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (s Provider) Register(app infra.Binder) {
	app.MustSingleton(NewEventStore)
}

func (s Provider) Boot(app infra.Resolver) {}
