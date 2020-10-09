package job

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/agent/config"
	"github.com/mylxsw/adanos-alert/agent/store"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
)

func eventSyncJob(eventStore store.EventStore, conf *config.Config, msgRPCServer protocol.MessageClient) error {
	for {
		message, err := eventStore.Dequeue()
		if err != nil || message == nil {
			break
		}

		if err := sendToServer(message, msgRPCServer, conf); err != nil {
			log.Warningf("事件同步失败，重新加入队列: %s", err)
			if err := eventStore.Enqueue(message); err != nil {
				log.Warningf("事件重新写入队列失败: %s, 事件内容：%s", err, message.Data)
			}

			// 如果写入出错，则休息1s，防止重试速度过快
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func sendToServer(msg *protocol.MessageRequest, msgRPCServer protocol.MessageClient, conf *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := msgRPCServer.Push(ctx, msg)
	if err != nil {
		return fmt.Errorf("RPC请求失败: %s", err)
	}

	log.WithFields(log.Fields{
		"id":   resp.Id,
		"body": msg.Data,
	}).Debugf("事件同步成功")

	return nil
}
