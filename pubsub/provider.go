package pubsub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/color"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(app container.Container) {
}

func (s ServiceProvider) Boot(app infra.Glacier) {
	app.MustResolve(func(em event.Manager, auditRepo repository.AuditLogRepo) {
		// 用户变更事件监听
		em.Listen(func(ev UserChangedEvent) {
			auditRepo.Add(repository.AuditLog{
				Type: repository.AuditLogTypeAction,
				Body: fmt.Sprintf("[%s] User %s: %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.User)),
			})
		})

		// 规则变更事件监听
		em.Listen(func(ev RuleChangedEvent) {
			auditRepo.Add(repository.AuditLog{
				Type: repository.AuditLogTypeAction,
				Body: fmt.Sprintf("[%s] Rule %s: %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.Rule)),
			})
		})

		// 钉钉机器人变更事件监听
		em.Listen(func(ev DingdingRobotEvent) {
			auditRepo.Add(repository.AuditLog{
				Type: repository.AuditLogTypeAction,
				Body: fmt.Sprintf("[%s] DingdingRobot %s: %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.DingDingRobot)),
			})
		})

		// 系统启停事件监听
		em.Listen(func(ev SystemUpDownEvent) {
			status := "down"
			if ev.Up {
				status = "up"
			}

			auditRepo.Add(repository.AuditLog{
				Type: repository.AuditLogTypeSystem,
				Body: fmt.Sprintf("[%s] System is changed to %s", ev.CreatedAt.Format(time.RFC3339), status),
			})
		})
	})
}

func serialize(data interface{}) string {
	res, _ := json.Marshal(data)
	return color.TextWrap(color.LightGrey, string(res))
}