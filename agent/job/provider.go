package job

import (
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/cron"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(protocol.NewMessageClient)
	app.MustSingleton(protocol.NewHeartbeatClient)
}

func (s ServiceProvider) Boot(app glacier.Glacier) {
	app.Cron(func(cr cron.Manager, cc container.Container) error {
		cc.Must(cr.Add("sync-message", "@every 5s", messageSyncJob))
		cc.Must(cr.Add("heartbeat", "@every 10s", heartbeatJob))

		return nil
	})
}
