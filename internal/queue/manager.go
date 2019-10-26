package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
)

type Manager interface {
	Enqueue(item repository.QueueItem) (string, error)
	StartWorker(ctx context.Context, workerID string)
	Pause(pause bool)
	RegisterHandler(name string, handler Handler)
}

type Handler func(item repository.QueueItem) error

type queueManager struct {
	lock     sync.RWMutex
	cc       *container.Container
	repo     repository.QueueRepo
	handlers map[string]Handler

	maxRetryTimes int
	paused        bool
}

// NewManager create a QueueManager
func NewManager(cc *container.Container) Manager {
	manager := queueManager{
		cc:       cc,
		paused:   true,
		handlers: make(map[string]Handler),
	}
	cc.MustResolve(func(repo repository.QueueRepo, conf *configs.Config) {
		manager.repo = repo
		manager.maxRetryTimes = conf.QueueJobMaxRetryTimes
	})
	return &manager
}

// RegisterHandler register a handler for job processing
func (manager *queueManager) RegisterHandler(name string, handler Handler) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.handlers[name] = handler
}

// Pause control whether the queue is working or paused
func (manager *queueManager) Pause(pause bool) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.paused = pause
}

// Enqueue add an item to queue
func (manager *queueManager) Enqueue(item repository.QueueItem) (string, error) {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	if _, ok := manager.handlers[item.Name]; !ok {
		return "", errors.New("not support such queueItem")
	}

	id, err := manager.repo.Enqueue(item)
	if err != nil {
		return "", fmt.Errorf("enqueu failed: %w", err)
	}

	return id.Hex(), nil
}

// StartWorker start a worker
func (manager *queueManager) StartWorker(ctx context.Context, workID string) {
	log.Debugf("queue worker [%s] started", workID)
	defer log.Debugf("queue worker [%s] stopped", workID)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			manager.run(ctx)
		}
	}
}

// isPaused whether the queue is paused
func (manager *queueManager) isPaused() bool {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.paused
}

func (manager *queueManager) run(ctx context.Context) {
	if manager.isPaused() {
		return
	}

	for item, err := manager.repo.Dequeue(); err == nil; {
		manager.handle(ctx, item)
		if manager.isPaused() {
			return
		}
	}
}

func (manager *queueManager) handle(ctx context.Context, item repository.QueueItem) {
	manager.lock.RLock()
	handler, ok := manager.handlers[item.Name]
	manager.lock.RUnlock()
	if !ok {
		item.Status = repository.QueueItemStatusCanceled
		if err := manager.repo.UpdateID(item.ID, item); err != nil {
			log.WithFields(log.Fields{
				"err":     err.Error(),
				"item_id": item.ID,
			}).Errorf("can not update queue item: %v", err)
		}

		return
	}

	// execute queue job handler
	if err := eliminatePanic(handler)(item); err != nil {
		// if job failed, check execute times, if requeue times > max requeueTimes, set job as failed
		// otherwise requeue it and try again latter
		if item.RequeueTimes > manager.maxRetryTimes {
			item.Status = repository.QueueItemStatusFailed
			if err := manager.repo.UpdateID(item.ID, item); err != nil {
				log.WithFields(log.Fields{
					"err":  err.Error(),
					"item": item,
				}).Errorf("can not update queue item: %v", err)
			}

			return
		}

		// try again latter
		item.NextExecuteAt = time.Now().Add(time.Duration((item.RequeueTimes+1)*30) * time.Second)
		if _, err := manager.repo.Enqueue(item); err != nil {
			log.WithFields(log.Fields{
				"err":  err.Error(),
				"item": item,
			}).Errorf("can not requeue item: %v", err)
		}

		return
	}

	// handler finished successful, update job status to succeed
	item.Status = repository.QueueItemStatusSucceed
	if err := manager.repo.UpdateID(item.ID, item); err != nil {
		log.WithFields(log.Fields{
			"err":  err.Error(),
			"item": item,
		}).Errorf("can not update queue item: %v", err)
	}
}

func eliminatePanic(cb Handler) Handler {
	return func(item repository.QueueItem) (err error) {
		defer func() {
			if err2 := recover(); err2 != nil {
				err = fmt.Errorf("handler panic with: %v", err2)
			}
		}()

		return cb(item)
	}
}
