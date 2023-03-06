package queue

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
)

// Manager 队列管理接口
type Manager interface {
	Enqueue(item repository.QueueJob) (string, error)
	StartWorker(ctx context.Context, workerID string)
	Pause(pause bool)
	Paused() bool
	Info() Info
	RegisterHandler(name string, handler Handler)
}

// Handler 队列消息处理器
type Handler func(item repository.QueueJob) error

// Info 队列状态信息
type Info struct {
	StartAt        time.Time `json:"start_at"`
	WorkerNum      int       `json:"worker_num"`
	ProcessedCount int64     `json:"processed_count"`
	FailedCount    int64     `json:"failed_count"`
}

type queueManager struct {
	lock     sync.RWMutex
	cc       infra.Resolver
	repo     repository.QueueRepo
	handlers map[string]Handler

	info Info

	maxRetryTimes int
	paused        bool
}

// NewManager create a QueueManager
func NewManager(cc infra.Resolver) Manager {
	manager := queueManager{
		cc:       cc,
		paused:   true,
		handlers: make(map[string]Handler),
		info: Info{
			StartAt:        time.Now(),
			WorkerNum:      0,
			ProcessedCount: 0,
		},
	}
	cc.MustResolve(func(repo repository.QueueRepo, conf *configs.Config) {
		manager.repo = repo
		manager.maxRetryTimes = conf.QueueJobMaxRetryTimes
	})
	return &manager
}

// Info return queue runtime info
func (manager *queueManager) Info() Info {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.info
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
func (manager *queueManager) Enqueue(item repository.QueueJob) (string, error) {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	if _, ok := manager.handlers[item.Name]; !ok {
		return "", errors.New("not support such queueItem")
	}

	id, err := manager.repo.Enqueue(context.TODO(), item)
	if err != nil {
		return "", fmt.Errorf("enqueu failed: %w", err)
	}

	return id.Hex(), nil
}

// StartWorker start a worker
func (manager *queueManager) StartWorker(ctx context.Context, workID string) {

	manager.lock.Lock()
	manager.info.WorkerNum++
	manager.lock.Unlock()

	if log.DebugEnabled() {
		log.Debugf("queue worker [%s] started", workID)
	}

	defer func() {
		manager.lock.Lock()
		manager.info.WorkerNum--
		manager.lock.Unlock()

		if log.DebugEnabled() {
			log.Debugf("queue worker [%s] stopped", workID)
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
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

// Paused whether the queue is paused
func (manager *queueManager) Paused() bool {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.paused
}

func (manager *queueManager) run(ctx context.Context) {
	if manager.Paused() {
		return
	}

	var item repository.QueueJob
	var err error
	item, err = manager.repo.Dequeue(ctx)

	for err == nil {
		manager.lock.Lock()
		manager.info.ProcessedCount++
		manager.lock.Unlock()

		manager.handle(ctx, item)
		if manager.Paused() {
			return
		}

		item, err = manager.repo.Dequeue(ctx)
	}
}

func (manager *queueManager) handle(ctx context.Context, item repository.QueueJob) {
	manager.lock.RLock()
	handler, ok := manager.handlers[item.Name]
	manager.lock.RUnlock()
	if !ok {
		item.Status = repository.QueueItemStatusCanceled
		if err := manager.repo.UpdateID(ctx, item.ID, item); err != nil {
			log.WithFields(log.Fields{
				"err":     err.Error(),
				"item_id": item.ID,
			}).Errorf("can not update queue item: %v", err)
		}

		return
	}

	// execute queue job handler
	if err := eliminatePanic(handler)(item); err != nil {
		manager.lock.Lock()
		manager.info.FailedCount++
		manager.lock.Unlock()

		item.LastError = err.Error()

		// if job failed, check execute times, if requeue times > max requeueTimes, set job as failed
		// otherwise requeue it and try again latter
		if item.RequeueTimes > manager.maxRetryTimes {
			item.Status = repository.QueueItemStatusFailed

			log.WithFields(log.Fields{
				"item":     item,
				"last_err": err,
			}).Errorf("job execute failed after max retry times")

			if err := manager.repo.UpdateID(ctx, item.ID, item); err != nil {
				log.WithFields(log.Fields{
					"err":  err.Error(),
					"item": item,
				}).Errorf("can not update queue item: %v", err)
			}

			return
		}

		log.WithFields(log.Fields{
			"err":  err.Error(),
			"item": item,
		}).Warningf("job execute failed, we will retry later")

		// try again latter
		item.NextExecuteAt = time.Now().Add(time.Duration((item.RequeueTimes+1)*30+rand.Intn(10)) * time.Second)
		if _, err := manager.repo.Enqueue(ctx, item); err != nil {
			log.WithFields(log.Fields{
				"err":  err.Error(),
				"item": item,
			}).Errorf("can not requeue item: %v", err)
		}

		return
	}

	// handler finished successful, update job status to succeed
	item.Status = repository.QueueItemStatusSucceed
	if err := manager.repo.UpdateID(ctx, item.ID, item); err != nil {
		log.WithFields(log.Fields{
			"err":  err.Error(),
			"item": item,
		}).Errorf("can not update queue item: %v", err)
	}
}

func eliminatePanic(cb Handler) Handler {
	return func(item repository.QueueJob) (err error) {
		defer func() {
			if err2 := recover(); err2 != nil {
				err = fmt.Errorf("handler panic with: %v", err2)
			}
		}()

		return cb(item)
	}
}
