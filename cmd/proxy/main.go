package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/connector"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/asteria/log"
	"github.com/urfave/cli"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := &cli.App{
		Name:    "adanos-proxy",
		Usage:   "adanos-proxy 是一个简单的命令行工具，你可以通过管道的方式，把消息直接发送给 adanos-alert 用于告警通知",
		Version: fmt.Sprintf("%s (%s)", Version, GitCommit[:8]),
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "adanos-server",
				Usage: "adanos-alert server addr",
			},
			&cli.StringFlag{
				Name:  "adanos-token",
				Value: "",
				Usage: "adanos-alert server token",
			},
			&cli.StringSliceFlag{
				Name:  "tag",
				Usage: "specify tags for alert message",
			},
			&cli.StringSliceFlag{
				Name:  "meta",
				Usage: "specify meta for alert message, meta must be key=value",
			},
			&cli.StringFlag{
				Name:  "origin",
				Usage: "specify origin for alert message",
			},
			&cli.StringFlag{
				Name:  "id",
				Usage: "指定当前 message 的 ID，用于消息抑制和自动恢复",
			},
			&cli.StringFlag{
				Name:  "inhibit-interval",
				Usage: "消息抑制周期，比如 1m 表示 1 分钟内只有一条相同 id 的消息会被接收",
			},
			&cli.StringFlag{
				Name:  "recovery-after",
				Usage: "自动恢复周期，比如 1m 表示 1 分钟内如果没有新的相同 id 的消息到达，自动创建一条已恢复的消息",
			},
			&cli.IntFlag{
				Name:  "max-lines",
				Value: 1000,
			},
			&cli.BoolFlag{
				Name:  "multi-line",
				Usage: "是否为多行日志，添加该选项后，所有输入将作为一个告警事件发送",
			},
		},
		Action: func(c *cli.Context) error {
			stdinLines := readStdin(c.Int("max-lines"))
			if stdinLines == "" {
				return nil
			}

			adanosServers := c.StringSlice("adanos-server")
			if len(adanosServers) == 0 {
				adanosServers = append(adanosServers, "http://localhost:19999")
			}

			ctl := connector.EventControl{
				ID:              c.String("id"),
				InhibitInterval: c.String("inhibit-interval"),
				RecoveryAfter:   c.String("recovery-after"),
			}

			multiLine := c.Bool("multi-line")
			if multiLine {
				evt := connector.NewEvent(stdinLines).
					WithTags(c.StringSlice("tag")...).
					WithOrigin(c.String("origin")).
					WithMetas(createMessageMeta(c.StringSlice("meta"))).
					WithCtl(ctl)

				ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
				return connector.NewConnector(c.String("adanos-token"), adanosServers...).Send(ctx, evt)
			}

			for _, line := range strings.Split(stdinLines, "\n") {
				if strings.TrimSpace(line) == "" {
					continue
				}

				evt := connector.NewEvent(line).
					WithTags(c.StringSlice("tag")...).
					WithOrigin(c.String("origin")).
					WithMetas(createMessageMeta(c.StringSlice("meta"))).
					WithCtl(ctl)

				ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
				if err := connector.NewConnector(c.String("adanos-token"), adanosServers...).Send(ctx, evt); err != nil {
					log.WithFields(log.Fields{
						"event": line,
					}).Errorf("send event to adanos-alert server failed: %v", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorf("error: %v", err)
		panic(err)
	}
}

func createMessageMeta(meta []string) repository.EventMeta {
	metas := make(repository.EventMeta)
	if meta != nil && len(meta) > 0 {
		for _, m := range meta {
			segs := strings.SplitN(m, "=", 2)
			metas[segs[0]] = misc.IfElse(len(segs) == 2, segs[1], "")
		}
	}
	return metas
}

func readStdin(maxLines int) string {
	result := ""

	reader := bufio.NewReader(os.Stdin)
	lineNo := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		lineNo++
		if lineNo > maxLines {
			break
		}

		result += line
	}

	return strings.TrimSpace(strings.TrimSuffix(result, "\n"))
}
