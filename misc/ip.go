package misc

import "github.com/mylxsw/go-toolkit/network"

func ServerIP() string {
	ips, err := network.GetLanIPs()
	if err != nil || len(ips) == 0 {
		return "unknown"
	}

	return ips[0]
}
