package rpc

import (
	"context"
	"encoding/json"

	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/glacier/infra"
)

// EventService is a service server for message processing
type EventService struct {
	cc         infra.Resolver
	msgService service.EventService `autowire:"@"`
}

// NewEventService create a new message service
func NewEventService(cc infra.Resolver) *EventService {
	ms := &EventService{cc: cc}
	cc.Must(cc.AutoWire(ms))
	return ms
}

// Push add a new message
func (ms *EventService) Push(ctx context.Context, request *protocol.MessageRequest) (*protocol.IDResponse, error) {
	var commonMessage extension.CommonEvent
	if err := json.Unmarshal([]byte(request.Data), &commonMessage); err != nil {
		return nil, err
	}

	id, err := ms.msgService.Add(ctx, commonMessage)
	if err != nil {
		return nil, err
	}

	return &protocol.IDResponse{Id: id.Hex()}, nil
}
