package rpc

import (
	"context"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
)

// MessageService is a service server for message processing
type MessageService struct {
	cc      container.Container
	msgRepo repository.MessageRepo
}

// NewMessageService create a new message service
func NewMessageService(cc container.Container) *MessageService {
	ms := MessageService{cc: cc}
	cc.MustResolve(func(mr repository.MessageRepo) {
		ms.msgRepo = mr
	})
	return &ms
}

// Push add a new message
func (ms *MessageService) Push(ctx context.Context, request *protocol.MessageRequest) (*protocol.IDResponse, error) {
	var meta repository.MessageMeta
	_ = coll.Map(request.Meta, &meta, func(val string) interface{} { return val })

	id, err := ms.msgRepo.AddWithContext(ctx, repository.Message{
		Content: request.Content,
		Meta:    meta,
		Tags:    request.Tags,
		Origin:  request.Origin,
	})
	if err != nil {
		return nil, err
	}

	return &protocol.IDResponse{Id: id.Hex()}, nil
}
