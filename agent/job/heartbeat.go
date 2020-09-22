package job

import (
	"context"
	"time"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/agent/config"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func heartbeatJob(cc container.Container, db *ledis.DB, conf *config.Config, hs protocol.HeartbeatClient) error {
	agentID, _ := db.Get([]byte("agent-id"))
	pingReq := protocol.PingRequest{
		AgentTs:       time.Now().Unix(),
		AgentIP:       misc.ServerIP(),
		AgentID:       string(agentID),
		ClientVersion: cc.MustGet(infra.VersionKey).(string),
		Agent: &protocol.AgentInfo{
			Listen:        conf.Listen,
			LogPath:       conf.LogPath,
			Host:          buildAgentInfoHost(),
			Load:          buildAgentInfoLoad(),
			MemorySwap:    buildAgentInfoMemorySwap(),
			MemoryVirtual: buildAgentInfoMemoryVirtual(),
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	pong, err := hs.Ping(ctx, &pingReq)
	if err != nil {
		log.Warningf("心跳上报失败: %v", err)
		return nil
	}

	log.Debugf("心跳上报成功，服务端版本: %s, 服务端时间戳: %v", pong.ServerVersion, pong.ServerTs)
	return nil
}

func buildAgentInfoHost() *protocol.AgentInfoHost {
	info, err := host.Info()
	if err != nil {
		log.Errorf("read agent host info failed: %v", err)
		return nil
	}

	return &protocol.AgentInfoHost{
		Hostname:        info.Hostname,
		Uptime:          int64(info.Uptime),
		BootTime:        int64(info.BootTime),
		Procs:           int64(info.Procs),
		Os:              info.OS,
		Platform:        info.Platform,
		PlatformFamily:  info.PlatformFamily,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}
}

func buildAgentInfoLoad() *protocol.AgentInfoLoad {
	la, err := load.Avg()
	if err != nil {
		log.Errorf("read system avg load failed: %v", err)
		return nil
	}

	return &protocol.AgentInfoLoad{
		Load1:  la.Load1,
		Load5:  la.Load5,
		Load15: la.Load15,
	}
}

func buildAgentInfoMemorySwap() *protocol.AgentInfoMemorySwap {
	sw, err := mem.SwapMemory()
	if err != nil {
		log.Errorf("read swap memory failed: %v", err)
		return nil
	}

	return &protocol.AgentInfoMemorySwap{
		Total:       int64(sw.Total),
		Used:        int64(sw.Used),
		Free:        int64(sw.Free),
		UsedPercent: sw.UsedPercent,
	}
}

func buildAgentInfoMemoryVirtual() *protocol.AgentInfoMemoryVirtual {
	vw, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("read virtual memory failed: %v", err)
		return nil
	}

	return &protocol.AgentInfoMemoryVirtual{
		Total:       int64(vw.Total),
		Available:   int64(vw.Available),
		Used:        int64(vw.Used),
		UsedPercent: vw.UsedPercent,
		Free:        int64(vw.Free),
		Buffers:     int64(vw.Buffers),
		Cached:      int64(vw.Cached),
	}
}
