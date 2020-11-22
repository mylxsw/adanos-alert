package service

import (
	"context"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventGroupService 用于对 evengGroup 操作的 service
type EventGroupService interface {
	// CutGroup 缩减分组中 event 的数量，只保留  keepCount 条（relation_ids 不为空的 events 不能删除）
	CutGroup(ctx context.Context, groupID primitive.ObjectID, keepCount int64) (int64, error)
}

type eventGroupService struct {
	cc           container.Container
	evtRepo      repository.EventRepo      `autowire:"@"`
	evtGroupRepo repository.EventGroupRepo `autowire:"@"`
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
