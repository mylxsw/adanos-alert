package rpc

import (
	"context"
	"encoding/json"

	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/container"
)

// MessageService is a service server for message processing
type MessageService struct {
	cc         container.Container
	msgService service.MessageService `autowire:"@"`
}

// NewMessageService create a new message service
func NewMessageService(cc container.Container) *MessageService {
	ms := &MessageService{cc: cc}
	cc.Must(cc.AutoWire(ms))
	return ms
}

// Push add a new message
func (ms *MessageService) Push(ctx context.Context, request *protocol.MessageRequest) (*protocol.IDResponse, error) {
	var commonMessage misc.CommonMessage
	if err := json.Unmarshal([]byte(request.Data), &commonMessage); err != nil {
		return nil, err
	}

	id, err := ms.msgService.Add(ctx, commonMessage)
	if err != nil {
		return nil, err
	}

	return &protocol.IDResponse{Id: id.Hex()}, nil
}
