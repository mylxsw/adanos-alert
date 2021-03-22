package job

import (
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		scheduler.Provider(func(cc infra.Resolver, creator scheduler.JobCreator) {
			cc.Must(creator.Add("sync-events", "@every 5s", eventSyncJob))
			cc.Must(creator.Add("heartbeat", "@every 10s", heartbeatJob))
		}),
	}
}

func (s Provider) Register(app infra.Binder) {
	app.MustSingleton(protocol.NewMessageClient)
	app.MustSingleton(protocol.NewHeartbeatClient)
}

func (s Provider) Boot(app infra.Resolver) {}
