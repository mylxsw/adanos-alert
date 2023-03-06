package job

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RecoveryJobName = "recovery"

type RecoveryJob struct {
	app       infra.Resolver
	executing chan interface{} // 标识当前Job是否在执行中
}

func NewRecoveryJob(app infra.Resolver) *RecoveryJob {
	return &RecoveryJob{app: app, executing: make(chan interface{}, 1)}
}

func (a *RecoveryJob) Handle() {
	select {
	case a.executing <- struct{}{}:
		defer func() { <-a.executing }()

		a.app.MustResolve(func(recoveryRepo repository.RecoveryRepo, eventRepo repository.EventRepo) {
			events, err := recoveryRepo.RecoverableEvents(context.TODO(), time.Now())
			if err != nil {
				log.Errorf("query recoverable events from mongodb failed: %v", err)
				return
			}

			for _, m := range events {
				(func(m repository.Recovery) {
					defer func() {
						if err := recover(); err != nil {
							log.With(m).Errorf("add recovery event failed: %v", err)
						} else {
							if err := recoveryRepo.Delete(context.TODO(), m.RecoveryID); err != nil {
								log.With(m).Errorf("remove recovery event from mongodb failed: %v", err)
							}
						}
					}()
					if len(m.RefIDs) == 0 {
						return
					}

					msgSample, err := eventRepo.Get(m.RefIDs[len(m.RefIDs)-1])
					if err != nil {
						log.With(m).Errorf("get recovery event sample failed: %v", err)
					}

					msgSample.Type = repository.EventTypeRecovery
					msgSample.ID = primitive.NilObjectID
					msgSample.GroupID = nil
					msgSample.CreatedAt = time.Now()
					msgSample.Status = ""
					msgSample.Meta["recovery-refs"] = m.RefIDs
					msgSample.Tags = append(misc.IfElse(
						msgSample.Tags == nil,
						make([]string, 0),
						msgSample.Tags,
					).([]string), "adanos-recovery")

					if _, err := eventRepo.AddWithContext(context.TODO(), msgSample); err != nil {
						log.With(m).Errorf("add recovery event failed: %v", err)
					}
				})(m)
			}
		})

	default:
		log.Warningf("the last recovery job is not finished yet, skip for this time")
	}
}
