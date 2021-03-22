package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (s Provider) Register(app infra.Binder) {
	app.MustSingleton(NewManager)
}

func (s Provider) Boot(app infra.Resolver) {}

func (s Provider) Daemon(ctx context.Context, app infra.Resolver) {
	app.MustResolve(func(manager Manager, conf *configs.Config) {
		if conf.NoJobMode {
			return
		}

		var wg sync.WaitGroup
		wg.Add(conf.QueueWorkerNum)

		for i := 0; i < conf.QueueWorkerNum; i++ {
			go func(i int) {
				defer wg.Done()
				manager.StartWorker(ctx, fmt.Sprintf("worker #%d", i))
			}(i)
		}

		manager.Pause(false)

		wg.Wait()
	})
}
