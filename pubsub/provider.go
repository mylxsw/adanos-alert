package pubsub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/asteria/color"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		event.Provider(func(cc infra.Resolver, em event.Listener) {
			cc.MustResolve(func(syslogRepo repository.SyslogRepo) {
				// 用户变更事件监听
				em.Listen(func(ev UserChangedEvent) {
					_, _ = syslogRepo.Add(repository.Syslog{
						Type: repository.SyslogTypeAction,
						Body: fmt.Sprintf("[%s] User %s %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.User)),
					})
				})

				// 规则变更事件监听
				em.Listen(func(ev RuleChangedEvent) {
					_, _ = syslogRepo.Add(repository.Syslog{
						Type: repository.SyslogTypeAction,
						Body: fmt.Sprintf("[%s] Rule %s %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.Rule)),
					})
				})

				// 钉钉机器人变更事件监听
				em.Listen(func(ev DingdingRobotEvent) {
					_, _ = syslogRepo.Add(repository.Syslog{
						Type: repository.SyslogTypeAction,
						Body: fmt.Sprintf("[%s] DingdingRobot %s %s", ev.CreatedAt.Format(time.RFC3339), ev.Type, serialize(ev.DingDingRobot)),
					})
				})

				// 系统启停事件监听
				em.Listen(func(ev SystemUpDownEvent) {
					_, _ = syslogRepo.Add(repository.Syslog{
						Type: repository.SyslogTypeSystem,
						Body: fmt.Sprintf("[%s] System is changed to %s", ev.CreatedAt.Format(time.RFC3339), misc.IfElse(ev.Up, "up", "down")),
					})
				})

				// 事件组事件清理
				em.Listen(func(ev EventGroupReduceEvent) {
					_, _ = syslogRepo.Add(repository.Syslog{
						Type: repository.SyslogTypeAction,
						Body: fmt.Sprintf("[%s] EventGroup's (%s) event count reduced to %d, deleted count=%d", ev.CreatedAt.Format(time.RFC3339), ev.GroupID.Hex(), ev.KeepCount, ev.DeleteCount),
					})
				})
			})
		}),
	}
}

func (s Provider) Register(infra.Binder) {}
func (s Provider) Boot(infra.Resolver)   {}

func serialize(data interface{}) string {
	res, _ := json.Marshal(data)
	return color.TextWrap(color.LightGrey, string(res))
}
