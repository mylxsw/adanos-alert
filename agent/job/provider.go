package job

import (
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/cron"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(protocol.NewMessageClient)
	app.MustSingleton(protocol.NewHeartbeatClient)
}

func (s ServiceProvider) Boot(app infra.Glacier) {
	app.Cron(func(cr cron.Manager, cc container.Container) error {
		cc.Must(cr.Add("sync-events", "@every 5s", eventSyncJob))
		cc.Must(cr.Add("heartbeat", "@every 10s", heartbeatJob))

		return nil
	})
}
