package template

import (
	"strconv"
	"strings"

	"github.com/mylxsw/asteria/log"
)

type OpenFalconIM struct {
	Priority    int
	Status      string
	Endpoint    string
	Body        string
	CurrentStep int
	FormatTime  string
}

// ParseOpenFalconImMessage 解析 Open-Falcon IM 消息
// https://github.com/open-falcon/falcon-plus/blob/2648553f82dd3986a91239d590461c0d795f63a4/modules/alarm/cron/builder.go#L43:6
func ParseOpenFalconImMessage(msg string) OpenFalconIM {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("parse open falcon message failed")
		}
	}()

	res := OpenFalconIM{}
	// [P3][PROBLEM][192.168.200.4][][ all(#1) agent.alive  1==1][O1 2019-07-08 23:35:00]
	segs := strings.Split(strings.TrimRight(strings.TrimLeft(msg, "["), "]"), "][")
	if len(segs) > 0 {
		res.Priority, _ = strconv.Atoi(strings.TrimLeft(segs[0], "P"))
	}

	if len(segs) > 1 {
		res.Status = segs[1]
	}

	if len(segs) > 2 {
		res.Endpoint = segs[2]
	}

	if len(segs) > 4 {
		res.Body = segs[4]
	}

	if len(segs) > 5 {
		ss := strings.SplitN(segs[5], " ", 2)
		res.CurrentStep, _ = strconv.Atoi(ss[0][1:])
		res.FormatTime = ss[1]
	}

	return res
}