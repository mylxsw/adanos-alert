package rpc

import (
	"context"
	"encoding/json"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
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
	var commonMessage misc.CommonMessage
	if err := json.Unmarshal([]byte(request.Data), &commonMessage); err != nil {
		return nil, err
	}

	id, err := ms.msgRepo.AddWithContext(ctx, commonMessage.GetRepoMessage())
	if err != nil {
		return nil, err
	}

	return &protocol.IDResponse{Id: id.Hex()}, nil
}
