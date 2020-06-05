package job

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/network"
	"github.com/pkg/errors"
)

func heartbeatJob(cc container.Container, hs protocol.HeartbeatClient) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	pingReq := protocol.PingRequest{
		AgentTs:       time.Now().Unix(),
		Agent:         agentName(),
		ClientVersion: cc.MustGet(glacier.VersionKey).(string),
	}

	pong, err := hs.Ping(ctx, &pingReq)
	if err != nil {
		log.Warning("心跳上报失败: %v", err)
		return errors.Wrap(err, "心跳上报失败")
	}

	log.Debugf("心跳上报成功，服务端版本: %s, 服务端时间戳: %v", pong.ServerVersion, pong.ServerTs)
	return nil
}

func agentName() string {
	ips, err := network.GetLanIPs()
	if err != nil || len(ips) == 0 {
		return "unknown"
	}

	return ips[0]
}
