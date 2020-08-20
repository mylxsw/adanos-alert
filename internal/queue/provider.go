package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(NewManager)
}

func (s ServiceProvider) Boot(app infra.Glacier) {

}

func (s ServiceProvider) Daemon(ctx context.Context, app infra.Glacier) {
	app.MustResolve(func(manager Manager, conf *configs.Config) {
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
