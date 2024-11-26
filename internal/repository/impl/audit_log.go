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

// SyslogRepo 系统日志仓库
type SyslogRepo struct {
	col *mongo.Collection
}

// NewSyslogRepo 创建一个 系统日志仓库
func NewSyslogRepo(db *mongo.Database) repository.SyslogRepo {
	return &SyslogRepo{
		col: db.Collection("syslog"),
	}
}

// Add 添加日志
func (alr *SyslogRepo) Add(al repository.Syslog) (id primitive.ObjectID, err error) {
	al.CreatedAt = time.Now()

	rs, err := alr.col.InsertOne(context.TODO(), al)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

// Get 查询单个日志
func (alr *SyslogRepo) Get(id primitive.ObjectID) (al repository.Syslog, err error) {
	err = alr.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&al)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

// Paginate 分页查询
func (alr *SyslogRepo) Paginate(filter bson.M, offset, limit int64) (als []repository.Syslog, next int64, err error) {
	als = make([]repository.Syslog, 0)
	cur, err := alr.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return
	}
	defer func() {
		_ = cur.Close(context.TODO())
	}()

	for cur.Next(context.TODO()) {
		var al repository.Syslog
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
func (alr *SyslogRepo) Delete(filter bson.M) error {
	_, err := alr.col.DeleteMany(context.TODO(), filter)
	return err
}

// DeleteID 删除单条
func (alr *SyslogRepo) DeleteID(id primitive.ObjectID) error {
	return alr.Delete(bson.M{"_id": id})
}
