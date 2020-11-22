package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventGroupService interface {
	// CutGroup 缩减分组中 event 的数量，只保留  keepCount 条（relation_ids 不为空的 events 不能删除）
	CutGroup(ctx context.Context, groupID primitive.ObjectID, keepCount int) error
}
