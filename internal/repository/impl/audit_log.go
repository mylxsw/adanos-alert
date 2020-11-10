package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditLogRepo 审计日志仓库
type AuditLogRepo struct {
	col *mongo.Collection
}

// NewAuditLogRepo 创建一个 审计日志仓库
func NewAuditLogRepo(db *mongo.Database) repository.AuditLogRepo {
	return &AuditLogRepo{
		col: db.Collection("audit_log"),
	}
}

// Add 添加日志
func (alr *AuditLogRepo) Add(al repository.AuditLog) (id primitive.ObjectID, err error) {
	al.CreatedAt = time.Now()

	rs, err := alr.col.InsertOne(context.TODO(), al)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

// Get 查询单个日志
func (alr *AuditLogRepo) Get(id primitive.ObjectID) (al repository.AuditLog, err error) {
	err = alr.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&al)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

// Paginate 分页查询
func (alr *AuditLogRepo) Paginate(filter bson.M, offset, limit int64) (als []repository.AuditLog, next int64, err error) {
	als = make([]repository.AuditLog, 0)
	cur, err := alr.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var al repository.AuditLog
		if err = cur.Decode(&al); err != nil {
			return
		}

		als = append(als, al)
	}

	if int64(len(als)) == limit {
		next = offset + limit
	}

	return
}

// Delete 删除
func (alr *AuditLogRepo) Delete(filter bson.M) error {
	_, err := alr.col.DeleteMany(context.TODO(), filter)
	return err
}

// DeleteID 删除单条
func (alr *AuditLogRepo) DeleteID(id primitive.ObjectID) error {
	return alr.Delete(bson.M{"_id": id})
}
