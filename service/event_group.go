package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventGroupService 用于对 eventGroup 操作的 service
type EventGroupService interface {
	// CutGroup 缩减分组中 event 的数量，只保留  keepCount 条（relation_ids 不为空的 events 不能删除）
	CutGroup(ctx context.Context, groupID primitive.ObjectID, keepCount int64) (int64, error)
	// EventShouldRealtime 检查事件是否应该即时发送，如果是，则执行 shouldFn 函数
	// 每次调用该函数，无论是否事件应该即时发送，都会更新时间静默周期，重新计算剩余时间
	EventShouldRealtime(ctx context.Context, evtKey string, silentPeriod time.Duration, shouldFn func())
}

type eventGroupService struct {
	cc           container.Container
	evtRepo      repository.EventRepo      `autowire:"@"`
	evtGroupRepo repository.EventGroupRepo `autowire:"@"`
	kvRepo       repository.KVRepo         `autowire:"@"`
}

// NewEventGroupService create a new event group service
func NewEventGroupService(cc container.Container) EventGroupService {
	eg := &eventGroupService{cc: cc}
	cc.Must(cc.AutoWire(eg))
	return eg
}

// CutGroup 实现 EventGroupService 接口
func (eg *eventGroupService) CutGroup(ctx context.Context, groupID primitive.ObjectID, keepCount int64) (int64, error) {
	if groupID.IsZero() {
		return 0, nil
	}

	allEventCount, err := eg.evtRepo.Count(bson.M{"group_ids": groupID})
	if err != nil {
		return 0, err
	}

	if allEventCount <= keepCount {
		return 0, nil
	}

	keepEventIDs, err := eg.evtRepo.FindIDs(ctx, bson.M{"group_ids": groupID}, keepCount)
	if err != nil {
		return 0, err
	}

	return allEventCount - keepCount, eg.evtRepo.Delete(bson.M{"group_ids": groupID, "_id": bson.M{"$nin": keepEventIDs}})
}

// EventShouldRealtime 返回事件是否应该实时发送
func (eg *eventGroupService) EventShouldRealtime(ctx context.Context, evtKey string, silentPeriod time.Duration, shouldFn func()) {
	realtimeKey := fmt.Sprintf("realtime_silent:%s", evtKey)
	_, err := eg.kvRepo.Get(realtimeKey)
	if err == nil || err == repository.ErrNotFound {
		if err == repository.ErrNotFound {
			shouldFn()
		}

		if err := eg.kvRepo.SetWithTTL(realtimeKey, time.Now().Format(time.RFC3339), silentPeriod); err != nil {
			log.WithFields(log.Fields{"key": evtKey}).Errorf("update event realtime silent failed: %v", err)
		}

		return
	}

	log.WithFields(log.Fields{"key": evtKey}).Errorf("query event realtime key failed: %v", err)
}
