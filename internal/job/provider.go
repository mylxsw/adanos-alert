package job

import (
	"time"

	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/period_job"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewAggregationJob)
	app.MustSingleton(NewTrigger)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.PeriodJob(func(pj *period_job.Manager, cc *container.Container) {
		cc.MustResolve(func(aggregationJob *AggregationJob, alertJob *TriggerJob) {
			pj.Run("aggregation", aggregationJob, 30*time.Second)
			pj.Run("trigger", alertJob, 15*time.Second)
		})
	})
}
