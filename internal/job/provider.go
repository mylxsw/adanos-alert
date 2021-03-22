package job

import (
	"fmt"
	"os"
	"sync"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mylxsw/adanos-alert/configs"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	hostname, _ := os.Hostname()

	return []infra.Provider{
		scheduler.Provider(s.jobs, scheduler.SetDistributeLockManagerOption(func(cc infra.Resolver) scheduler.DistributeLockManager {
			return NewDistributeLockManager(
				cc.MustGet((*repository.LockRepo)(nil)).(repository.LockRepo),
				fmt.Sprintf("%s(%s)", hostname, configs.Get(cc).Listen),
			)
		})),
	}
}

func (s Provider) jobs(cc infra.Resolver, creator scheduler.JobCreator) {
	cc.MustResolve(func(conf *configs.Config, aggregationJob *AggregationJob, alertJob *TriggerJob, recoveryJob *RecoveryJob, lockRepo repository.LockRepo) {
		if conf.NoJobMode {
			return
		}

		_ = creator.Add(AggregationJobName, fmt.Sprintf("@every %s", conf.AggregationPeriod), aggregationJob.Handle)
		_ = creator.Add(TriggerJobName, fmt.Sprintf("@every %s", conf.ActionTriggerPeriod), alertJob.Handle)
		_ = creator.Add(RecoveryJobName, fmt.Sprintf("@every %s", conf.AggregationPeriod), recoveryJob.Handle)
	})
}

func (s Provider) Register(app infra.Binder) {
	app.MustSingleton(NewAggregationJob)
	app.MustSingleton(NewTrigger)
	app.MustSingleton(NewRecoveryJob)
}

func (s Provider) Boot(app infra.Resolver) {}

type DistributeLockManager struct {
	syncLock sync.RWMutex
	lockRepo repository.LockRepo
	lockID   primitive.ObjectID
	locked   bool
	owner    string
}

func NewDistributeLockManager(lockRepo repository.LockRepo, owner string) scheduler.DistributeLockManager {
	return &DistributeLockManager{lockRepo: lockRepo, locked: false, owner: owner}
}

var lockResource = "crontab-lock"

func (d *DistributeLockManager) TryLock() error {
	d.syncLock.Lock()
	defer d.syncLock.Unlock()

	if d.locked {
		if _, err := d.lockRepo.Renew(d.lockID, 90); err != nil {
			if err == repository.ErrLockNotFound {
				if err := d.lock(); err != nil {
					return err
				}
			} else {
				return errors.Wrap(err, "renew lock failed")
			}
		}
	} else {
		if err := d.lock(); err != nil {
			return err
		}
	}

	return nil
}

func (d *DistributeLockManager) lock() error {
	lock, err := d.lockRepo.Lock(lockResource, d.owner, 90)
	if err != nil {
		if err == repository.ErrAlreadyLocked {
			d.lockID = primitive.NilObjectID
			d.locked = false
			return nil
		}

		return errors.Wrap(err, "acquire lock failed")
	}

	d.lockID = lock.LockID
	d.locked = true

	if log.DebugEnabled() {
		log.Debugf("got distribute lock, owner=%s", d.owner)
	}

	return nil
}

func (d *DistributeLockManager) TryUnLock() error {
	d.syncLock.Lock()
	defer d.syncLock.Unlock()

	if !d.locked {
		return nil
	}

	if err := d.lockRepo.UnLock(d.lockID); err != nil {
		return err
	}

	d.locked = false
	d.lockID = primitive.NilObjectID

	if log.DebugEnabled() {
		log.Debugf("distribute lock has been released")
	}

	return nil
}

func (d *DistributeLockManager) HasLock() bool {
	d.syncLock.RLock()
	defer d.syncLock.RUnlock()

	return d.locked
}
