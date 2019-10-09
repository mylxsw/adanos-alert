package job

import (
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/period_job"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewAggregationJob)
	app.MustSingleton(NewTrigger)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.PeriodJob(func(pj *period_job.Manager, cc *container.Container) {
		cc.MustResolve(func(conf *configs.Config, aggregationJob *AggregationJob, alertJob *TriggerJob) {
			pj.Run("aggregation", aggregationJob, conf.AggregationPeriod)
			pj.Run("trigger", alertJob, conf.ActionTriggerPeriod)
		})
	})
}
