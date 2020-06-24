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

func messageSyncJob(messageStore store.MessageStore, conf *config.Config, msgRPCServer protocol.MessageClient) error {
	for {
		message, err := messageStore.Dequeue()
		if err != nil || message == nil {
			break
		}

		if err := sendToServer(message, msgRPCServer, conf); err != nil {
			log.Warningf("消息同步失败，重新加入队列: %s", err)
			if err := messageStore.Enqueue(message); err != nil {
				log.Warningf("消息重新写入队列失败: %s, 消息内容：%s", err, message.Data)
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
	}).Debugf("消息同步成功")

	return nil
}
