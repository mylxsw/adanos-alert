package template

import (
	"fmt"
	"testing"
)

func TestParseOpenFalconImMessage(t *testing.T) {
	im := ParseOpenFalconImMessage("[P3][PROBLEM][192.168.200.4][][ all(#1) agent.alive  1==1][O1 2019-07-08 23:35:00]")
	res := fmt.Sprintf("priority=%d, status=%s, body=%s, endpoint=%s, step=%d, t=%s", im.Priority, im.Status, im.Body, im.Endpoint, im.CurrentStep, im.FormatTime)
	if res != "priority=3, status=PROBLEM, body= all(#1) agent.alive  1==1, endpoint=192.168.200.4, step=1, t=2019-07-08 23:35:00" {
		t.Errorf("test failed")
	}
}
