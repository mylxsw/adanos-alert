package repository

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventGroupRepo struct {
	Groups []repository.EventGroup
}

func (m *EventGroupRepo) StatByDatetimeCount(ctx context.Context, filter bson.M, startTime, endTime time.Time, hour int64) ([]repository.EventGroupByDatetimeCount, error) {
	panic("implement me")
}

func (m *EventGroupRepo) StatByAggCountInPeriod(ctx context.Context, ruleID primitive.ObjectID, startTime, endTime time.Time, hour int64) ([]repository.EventGroupAggByDatetimeCount, error) {
	panic("implement me")
}

func (m *EventGroupRepo) StatByAggCount(ctx context.Context, ruleID primitive.ObjectID, startTime, endTime time.Time) ([]repository.EventGroupAggCount, error) {
	panic("implement me")
}

func (m *EventGroupRepo) StatByRuleCount(ctx context.Context, startTime, endTime time.Time) ([]repository.EventGroupByRuleCount, error) {
	panic("implement me")
}

func (m *EventGroupRepo) StatByUserCount(ctx context.Context, startTime, endTime time.Time) ([]repository.EventGroupByUserCount, error) {
	panic("implement me")
}


func (m *EventGroupRepo) LastGroup(filter bson.M) (grp repository.EventGroup, err error) {
	panic("implement me")
}

func NewMessageGroupRepo() repository.EventGroupRepo {
	return &EventGroupRepo{Groups: make([]repository.EventGroup, 0)}
}

func (m *EventGroupRepo) Add(grp repository.EventGroup) (id primitive.ObjectID, err error) {
	panic("implement me")
}

func (m *EventGroupRepo) Get(id primitive.ObjectID) (grp repository.EventGroup, err error) {
	panic("implement me")
}

func (m *EventGroupRepo) Find(filter bson.M) (grps []repository.EventGroup, err error) {
	panic("implement me")
}

func (m *EventGroupRepo) Paginate(filter bson.M, offset, limit int64) (grps []repository.EventGroup, next int64, err error) {
	panic("implement me")
}

func (m *EventGroupRepo) Delete(filter bson.M) error {
	m.Groups = m.filter(filter)
	return nil
}

func (m *EventGroupRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m *EventGroupRepo) Traverse(filter bson.M, cb func(grp repository.EventGroup) error) error {
	for _, grp := range m.filter(filter) {
		if err := cb(grp); err != nil {
			return err
		}
	}

	return nil
}

func (m *EventGroupRepo) UpdateID(id primitive.ObjectID, grp repository.EventGroup) error {
	for i, g := range m.Groups {
		if g.ID == id {
			m.Groups[i] = grp
			break
		}
	}

	return nil
}

func (m *EventGroupRepo) Count(filter bson.M) (int64, error) {
	return int64(len(m.filter(filter))), nil
}

func (m *EventGroupRepo) CollectingGroup(rule repository.EventGroupRule) (group repository.EventGroup, err error) {
	groups := m.filter(bson.M{"rule._id": rule.ID, "status": repository.EventGroupStatusCollecting})
	if len(groups) == 0 {
		group = repository.EventGroup{
			ID:        primitive.NewObjectID(),
			Rule:      rule,
			Status:    repository.EventGroupStatusCollecting,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		m.Groups = append(m.Groups, group)
		return
	}

	return groups[0], nil
}

func (m *EventGroupRepo) filter(filter bson.M) (groups []repository.EventGroup) {
	err := coll.MustNew(m.Groups).Filter(func(grp repository.EventGroup) bool {
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
