package repository

import (
	"errors"

	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
)

var ErrNotFound = errors.New("not found")

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(impl.NewSequenceRepo)
	app.MustSingleton(impl.NewKVRepo)
	app.MustSingleton(impl.NewMessageRepo)
	app.MustSingleton(impl.NewMessageGroupRepo)
	app.MustSingleton(impl.NewUserRepo)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {}
