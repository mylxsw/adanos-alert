package rpc

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
)

// HeartbeatService is a service server for heartbeat
type HeartbeatService struct {
	cc container.Container
}

func NewHeartbeatService(cc container.Container) *HeartbeatService {
	return &HeartbeatService{cc: cc}
}

func (h *HeartbeatService) Ping(ctx context.Context, request *protocol.PingRequest) (*protocol.PongResponse, error) {
	log.Debugf("agent heartbeat received, ip=%s, version=%s, ts=%v", request.Agent, request.ClientVersion, request.AgentTs)
	return &protocol.PongResponse{
		ServerTs:      time.Now().Unix(),
		ServerVersion: h.cc.MustGet(glacier.VersionKey).(string),
	}, nil
}
