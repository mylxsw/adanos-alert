package job

import (
	"context"
	"time"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

func heartbeatJob(cc container.Container, db *ledis.DB, hs protocol.HeartbeatClient) error {
	agentID, _ := db.Get([]byte("agent-id"))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	pingReq := protocol.PingRequest{
		AgentTs:       time.Now().Unix(),
		AgentIP:       misc.ServerIP(),
		AgentID:       string(agentID),
		ClientVersion: cc.MustGet(infra.VersionKey).(string),
	}

	pong, err := hs.Ping(ctx, &pingReq)
	if err != nil {
		log.Warningf("心跳上报失败: %v", err)
		return nil
	}

	log.Debugf("心跳上报成功，服务端版本: %s, 服务端时间戳: %v", pong.ServerVersion, pong.ServerTs)
	return nil
}
