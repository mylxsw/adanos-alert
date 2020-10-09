package store

import (
	"encoding/json"
	"errors"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
)

type EventStore interface {
	Enqueue(msg *protocol.MessageRequest) error
	Dequeue() (*protocol.MessageRequest, error)
}

// eventStore 用于本地临时存储 message
type eventStore struct {
	db  *ledis.DB
	key []byte
}

// NewEventStore create a new eventStore
func NewEventStore(db *ledis.DB) EventStore {
	return &eventStore{db: db, key: []byte("messages")}
}

// Enqueue 事件加入队列
func (ms *eventStore) Enqueue(msg *protocol.MessageRequest) error {
	_, err := ms.db.LPush(ms.key, ms.serialize(msg))
	return err
}

// Dequeue 从队列中读取事件
func (ms *eventStore) Dequeue() (*protocol.MessageRequest, error) {
	message, err := ms.db.RPop(ms.key)
	if err != nil {
		log.Errorf("读取本地存储失败: %s", err)
		return nil, err
	}

	if message == nil {
		return nil, errors.New("读取失败")
	}

	var req protocol.MessageRequest
	ms.unserialize(message, &req)

	return &req, nil
}

func (ms *eventStore) serialize(msg interface{}) []byte {
	res, _ := json.Marshal(msg)
	return res
}

func (ms *eventStore) unserialize(data []byte, res interface{}) {
	_ = json.Unmarshal(data, &res)
}
