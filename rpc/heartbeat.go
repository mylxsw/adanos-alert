package rpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
)

// HeartbeatService is a service server for heartbeat
type HeartbeatService struct {
	cc infra.Resolver
}

func NewHeartbeatService(cc infra.Resolver) *HeartbeatService {
	return &HeartbeatService{cc: cc}
}

func (h *HeartbeatService) Ping(ctx context.Context, request *protocol.PingRequest) (*protocol.PongResponse, error) {
	if log.DebugEnabled() {
		log.Debugf("agent heartbeat received, id=%s, ip=%s, version=%s, ts=%v", request.AgentID, request.AgentIP, request.ClientVersion, request.AgentTs)
	}

	h.cc.MustResolve(func(agent repository.AgentRepo) {
		lastStat, _ := json.Marshal(request.GetAgent())
		if _, err := agent.Update(repository.Agent{
			IP:          request.GetAgentIP(),
			AgentID:     request.GetAgentID(),
			Version:     request.GetClientVersion(),
			LastAliveAt: time.Now(),
			LastStat:    string(lastStat),
		}); err != nil {
			log.WithFields(log.Fields{
				"req": request,
			}).Errorf("agent status update failed: %v", err)
		}
	})
	return &protocol.PongResponse{
		ServerTs:      time.Now().Unix(),
		ServerVersion: h.cc.MustGet(infra.VersionKey).(string),
	}, nil
}
