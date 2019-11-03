package impl

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/cron"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewSequenceRepo)
	app.MustSingleton(NewKVRepo)
	app.MustSingleton(NewMessageRepo)
	app.MustSingleton(NewMessageGroupRepo)
	app.MustSingleton(NewUserRepo)
	app.MustSingleton(NewRuleRepo)
	app.MustSingleton(NewQueueRepo)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.Cron(func(cr cron.Manager, cc *container.Container) error {
		return cc.Resolve(func(repo repository.KVRepo) {
			_ = cr.Add("kv_repository_gc", "@hourly", func() {
				if err := repo.GC(); err != nil {
					log.Errorf("kv repo gc failed: %v", err)
				}
			})
		})
	})
}
