package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	// Add add a new message to repository
	Add(ctx context.Context, msg misc.CommonMessage) (primitive.ObjectID, error)
}

type messageService struct {
	cc      container.Container
	kvRepo  repository.KVRepo      `autowire:"@"`
	msgRepo repository.MessageRepo `autowire:"@"`
}

func NewMessageService(cc container.Container) MessageService {
	ms := &messageService{cc: cc}
	cc.Must(cc.AutoWire(ms))
	return ms
}

func (m *messageService) Add(ctx context.Context, msg misc.CommonMessage) (primitive.ObjectID, error) {
	controlMessage := msg.GetControlMessage()

	var msgID primitive.ObjectID
	var err error

	defer func() {
		// 自动恢复
		recoveryAfter := controlMessage.GetRecoveryAfter()
		if controlMessage.ID != "" && recoveryAfter > 0 {
			recoveryRepo := m.cc.MustGet(new(repository.RecoveryRepo)).(repository.RecoveryRepo)
			if err := recoveryRepo.Register(ctx, time.Now().Add(recoveryAfter), controlMessage.ID, msgID); err != nil {
				log.Errorf("register recovery message(id=%s) failed: %v", msgID.Hex(), err)
			}
		}
	}()

	if controlMessage.ID != "" {
		key := fmt.Sprintf("msgctl:inhibit:%s", controlMessage.ID)
		// 消息抑制
		inhibitInterval := controlMessage.GetInhibitInterval()
		if inhibitInterval > 0 {
			if _, err := m.kvRepo.Get(key); err != nil {
				if err := m.kvRepo.SetWithTTL(key, time.Now().String(), inhibitInterval); err != nil {
					log.Errorf("set inhibit interval for %s failed: %v", key, err)
				}
			} else {
				// 消息被抑制，直接丢弃
				log.WithFields(log.Fields{
					"key": key,
					"ctl": msg.GetControlMessage(),
					"msg": msg.GetRepoMessage(),
				}).Debugf("message is discard because it's been inhibited")
				return primitive.NilObjectID, nil
			}
		}
	}

	// 保存消息
	msgID, err = m.msgRepo.AddWithContext(ctx, msg.GetRepoMessage())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return msgID, nil
}
