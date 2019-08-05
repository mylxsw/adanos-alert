package repository

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/go-toolkit/collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageRepo struct {
	Messages []repository.Message
}

func NewMessageRepo() repository.MessageRepo {
	return &MessageRepo{Messages: make([]repository.Message, 0),}
}

func (m *MessageRepo) Add(msg repository.Message) (id primitive.ObjectID, err error) {
	msg.ID = primitive.NewObjectID()
	msg.CreatedAt = time.Now()

	m.Messages = append(m.Messages, msg)
	return msg.ID, nil
}

func (m *MessageRepo) Get(id primitive.ObjectID) (msg repository.Message, err error) {
	for _, msg := range m.Messages {
		if msg.ID == id {
			return msg, nil
		}
	}

	return msg, repository.ErrNotFound
}

func (m *MessageRepo) Find(filter bson.M) (messages []repository.Message, err error) {
	panic("implement me")
}

func (m *MessageRepo) Paginate(filter bson.M, offset, limit int64) (messages []repository.Message, next int64, err error) {
	panic("implement me")
}

func (m *MessageRepo) Delete(filter bson.M) error {
	m.Messages = m.filter(filter)
	return nil
}

func (m *MessageRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m *MessageRepo) Traverse(filter bson.M, cb func(msg repository.Message) error) error {
	for _, msg := range m.filter(filter) {
		if err := cb(msg); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageRepo) UpdateID(id primitive.ObjectID, update repository.Message) error {
	for i, msg := range m.Messages {
		if msg.ID == id {
			m.Messages[i] = update
			break
		}
	}

	return nil
}

func (m *MessageRepo) Count(filter bson.M) (int64, error) {
	return int64(len(m.filter(filter))), nil
}

func (m *MessageRepo) filter(filter bson.M) (messages []repository.Message) {
	err := collection.MustNew(m.Messages).Filter(func(msg repository.Message) bool {
		if status, ok := filter["status"]; ok && msg.Status != status {
			return false
		}

		if id, ok := filter["_id"]; ok && id != msg.ID {
			return false
		}

		return true
	}).All(&messages)

	if err != nil {
		panic(err)
	}

	return
}
