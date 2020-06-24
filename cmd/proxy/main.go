package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/misc"
	"github.com/mylxsw/asteria/log"
	"github.com/pkg/errors"
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
			&cli.IntFlag{
				Name:  "max-lines",
				Value: 100,
			},
		},
		Action: func(c *cli.Context) error {
			message := readStdin(c.Int("max-lines"))
			if message == "" {
				return nil
			}

			adanosServers := c.StringSlice("adanos-server")
			if len(adanosServers) == 0 {
				adanosServers = append(adanosServers, "http://localhost:19999")
			}

			return sendMessage(
				message,
				createMessageMeta(c.StringSlice("meta")),
				c.StringSlice("tag"),
				c.String("origin"),
				adanosServers,
				c.String("adanos-token"),
			)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorf("error: %v", err)
		panic(err)
	}
}

func createMessageMeta(meta []string) repository.MessageMeta {
	metas := make(repository.MessageMeta)
	if meta != nil && len(meta) > 0 {
		for _, m := range meta {
			segs := strings.SplitN(m, "=", 2)
			if len(segs) == 2 {
				metas[segs[0]] = segs[1]
			} else {
				metas[segs[0]] = ""
			}
		}
	}
	return metas
}

func sendMessage(message string, metas repository.MessageMeta, tags []string, origin string, adanosServers []string, adanosToken string) error {
	commonMessage := misc.CommonMessage{
		Content: message,
		Meta:    metas,
		Tags:    tags,
		Origin:  origin,
	}
	data, _ := json.Marshal(commonMessage)

	var err error
	for _, s := range adanosServers {
		if err = sendMessageToServer(commonMessage, data, s, adanosToken); err == nil {
			break
		}

		log.Warningf("send to server %s failed: %v", s, err)
	}

	return err
}

func sendMessageToServer(commonMessage misc.CommonMessage, data []byte, adanosServer, adanosToken string) error {
	reqURL := fmt.Sprintf("%s/api/messages/", strings.TrimRight(adanosServer, "/"))

	log.WithFields(log.Fields{
		"message": commonMessage,
	}).Debugf("request: %v", reqURL)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}

	if adanosToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adanosToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body failed")
	}

	log.Debugf("response: %v", string(respBody))
	return nil
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
