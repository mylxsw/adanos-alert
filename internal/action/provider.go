package action

import (
	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/pkg/errors"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewManager)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.MustResolve(func(manager Manager, queueManager queue.Manager) {
		manager.Register("http", NewHttpAction(manager))
		manager.Register("dingding", NewDingdingAction(manager))
		manager.Register("email", NewEmailAction(manager))
		manager.Register("wechat", NewWechatAction(manager))

		queueManager.RegisterHandler("action", func(item repository.QueueJob) error {
			var payload Payload
			if err := payload.Decode([]byte(item.Payload)); err != nil {
				log.WithFields(log.Fields{
					"item": item,
					"err":  err.Error(),
				}).Errorf("can not decode payload: %s", err)
				return errors.Wrap(err, "can not decode payload")
			}

			return manager.Run(payload.Action).Handle(payload.Rule, payload.Trigger, payload.Group)
		})
	})
}
