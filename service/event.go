package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventService interface {
	// Add add a new event to repository
	Add(ctx context.Context, msg extension.CommonEvent) (primitive.ObjectID, error)
}

type eventService struct {
	cc      container.Container
	kvRepo  repository.KVRepo    `autowire:"@"`
	msgRepo repository.EventRepo `autowire:"@"`
}

func NewEventService(cc container.Container) EventService {
	ms := &eventService{cc: cc}
	cc.Must(cc.AutoWire(ms))
	return ms
}

func (m *eventService) Add(ctx context.Context, msg extension.CommonEvent) (primitive.ObjectID, error) {
	controlMessage := msg.GetControl()

	var msgID primitive.ObjectID
	var err error

	defer func() {
		// 自动恢复
		recoveryAfter := controlMessage.GetRecoveryAfter()
		if controlMessage.ID != "" && recoveryAfter > 0 {
			recoveryRepo := m.cc.MustGet(new(repository.RecoveryRepo)).(repository.RecoveryRepo)
			if err := recoveryRepo.Register(ctx, time.Now().Add(recoveryAfter), controlMessage.ID, msgID); err != nil {
				log.Errorf("register recovery event(id=%s) failed: %v", msgID.Hex(), err)
			}
		}
	}()

	if controlMessage.ID != "" {
		key := fmt.Sprintf("msgctl:inhibit:%s", controlMessage.ID)
		// 事件抑制
		inhibitInterval := controlMessage.GetInhibitInterval()
		if inhibitInterval > 0 {
			if _, err := m.kvRepo.Get(key); err != nil {
				if err := m.kvRepo.SetWithTTL(key, time.Now().String(), inhibitInterval); err != nil {
					log.Errorf("set inhibit interval for %s failed: %v", key, err)
				}
			} else {
				// 事件被抑制，直接丢弃
				log.WithFields(log.Fields{
					"key": key,
					"ctl": msg.GetControl(),
					"msg": msg.CreateRepoEvent(),
				}).Debugf("event is discard because it's been inhibited")
				return primitive.NilObjectID, nil
			}
		}
	}

	// 保存事件
	msgID, err = m.msgRepo.AddWithContext(ctx, msg.CreateRepoEvent())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return msgID, nil
}
