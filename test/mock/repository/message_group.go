package repository

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/go-toolkit/collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageGroupRepo struct {
	Groups []repository.MessageGroup
}

func NewMessageGroupRepo() repository.MessageGroupRepo {
	return &MessageGroupRepo{Groups: make([]repository.MessageGroup, 0),}
}

func (m *MessageGroupRepo) Add(grp repository.MessageGroup) (id primitive.ObjectID, err error) {
	panic("implement me")
}

func (m *MessageGroupRepo) Get(id primitive.ObjectID) (grp repository.MessageGroup, err error) {
	panic("implement me")
}

func (m *MessageGroupRepo) Find(filter bson.M) (grps []repository.MessageGroup, err error) {
	panic("implement me")
}

func (m *MessageGroupRepo) Paginate(filter bson.M, offset, limit int64) (grps []repository.MessageGroup, next int64, err error) {
	panic("implement me")
}

func (m *MessageGroupRepo) Delete(filter bson.M) error {
	m.Groups = m.filter(filter)
	return nil
}

func (m *MessageGroupRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m *MessageGroupRepo) Traverse(filter bson.M, cb func(grp repository.MessageGroup) error) error {
	panic("implement me")
}

func (m *MessageGroupRepo) Update(id primitive.ObjectID, grp repository.MessageGroup) error {
	panic("implement me")
}

func (m *MessageGroupRepo) Count(filter bson.M) (int64, error) {
	return int64(len(m.filter(filter))), nil
}

func (m *MessageGroupRepo) CollectingGroup(rule repository.MessageGroupRule) (group repository.MessageGroup, err error) {
	groups := m.filter(bson.M{"rule._id": rule.ID, "status": repository.MessageGroupStatusCollecting})
	if len(groups) == 0 {
		group = repository.MessageGroup{
			ID:        primitive.NewObjectID(),
			Rule:      rule,
			Status:    repository.MessageGroupStatusCollecting,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		m.Groups = append(m.Groups, group)
		return
	}

	return groups[0], nil
}

func (m *MessageGroupRepo) filter(filter bson.M) (groups []repository.MessageGroup) {
	err := collection.MustNew(m.Groups).Filter(func(grp repository.MessageGroup) bool {
		if status, ok := filter["status"]; ok && grp.Status != status {
			return false
		}

		if ruleId, ok := filter["rule._id"]; ok && grp.Rule.ID != ruleId {
			return false
		}

		if id, ok := filter["_id"]; ok && id != grp.ID {
			return false
		}

		return true
	}).All(&groups)

	if err != nil {
		panic(err)
	}

	return
}
