package job

import (
	"fmt"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/cron"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewAggregationJob)
	app.MustSingleton(NewTrigger)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.Cron(func(cr cron.Manager, cc *container.Container) error {
		return cc.Resolve(func(conf *configs.Config, aggregationJob *AggregationJob, alertJob *TriggerJob) {
			_ = cr.Add(AggregationJobName, fmt.Sprintf("@every %s", conf.AggregationPeriod), aggregationJob.Handle)
			_ = cr.Add(TriggerJobName, fmt.Sprintf("@every %s", conf.ActionTriggerPeriod), alertJob.Handle)
		})
	})
}
