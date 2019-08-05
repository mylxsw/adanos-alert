package job

import (
	"github.com/mylxsw/go-toolkit/container"
)

type AlertJob struct {
	app *container.Container
}

func NewAlertJob(app *container.Container) *AlertJob {
	return &AlertJob{app: app}
}

func (a AlertJob) Handle() {
	a.app.MustResolve(func() {

	})
}
