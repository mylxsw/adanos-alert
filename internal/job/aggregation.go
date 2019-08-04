package job

import (
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-toolkit/container"
)

type AggregationJob struct {
	app *container.Container
}

func NewAggregationJob(app *container.Container) *AggregationJob {
	return &AggregationJob{app: app}
}

func (a *AggregationJob) Handle() {
	log.Debug("Hello, world")
}
