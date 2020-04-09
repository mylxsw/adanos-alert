package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LockRepo struct {
	col *mongo.Collection
}

// Lock acquire a lock
// One of three things will happen when we run this change (upsert).
// 1) If the resource exists and has any locks on it (shared or
//    exclusive), we will get a duplicate key error.
// 2) If the resource exists but doesn't have any locks on it, we will
//    update it to obtain an exclusive lock.
// 3) If the resource doesn't exist yet, it will be inserted which will
//    give us an exclusive lock on it.
func (l *LockRepo) Lock(resource string, owner string, ttl uint) (*repository.Lock, error) {
	now := time.Now()
	rs := l.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{"resource": resource, "$or": bson.A{
			bson.M{"acquired": false},
			bson.M{"expired_at": bson.M{"$lte": now}},
		}},
		bson.M{
			"$set": bson.M{
				"lock_id":    primitive.NewObjectID(),
				"resource":   resource,
				"acquired":   true,
				"ttl":        ttl,
				"owner":      owner,
				"expired_at": now.Add(time.Duration(ttl) * time.Second),
				"created_at": now,
				"renewed_at": now,
			},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	var lock repository.Lock
	if err := rs.Decode(&lock); err != nil {
		if commandErr, ok := err.(mongo.CommandError); ok {
			if commandErr.Name == "DuplicateKey" {
				return nil, repository.ErrAlreadyLocked
			}
		}
		return nil, err
	}

	return &lock, nil
}

func (l *LockRepo) UnLock(lockID primitive.ObjectID) error {
	rs := l.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{
			"lock_id": lockID,
		},
		bson.M{
			"$set": bson.M{
				"acquired":   false,
				"ttl":        0,
				"owner":      "",
				"expired_at": time.Now().Add(-time.Hour),
			},
		},
		options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After),
	)

	var lock repository.Lock
	if err := rs.Decode(&lock); err != nil {
		if err == mongo.ErrNoDocuments {
			return repository.ErrLockNotFound
		}

		return err
	}

	return nil
}

func (l *LockRepo) Renew(lockID primitive.ObjectID, ttl uint) (*repository.Lock, error) {
	now := time.Now()
	rs := l.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{
			"lock_id": lockID,
		},
		bson.M{
			"$set": bson.M{
				"acquired":   true,
				"ttl":        ttl,
				"expired_at": now.Add(time.Duration(ttl) * time.Second),
				"renewed_at": now,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After),
	)

	var lock repository.Lock
	if err := rs.Decode(&lock); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrLockNotFound
		}

		return nil, err
	}

	return &lock, nil
}

func (l *LockRepo) Remove(resource string) error {
	_, err := l.col.DeleteOne(
		context.TODO(),
		bson.M{"resource": resource},
	)

	return err
}

func NewLockRepo(db *mongo.Database) repository.LockRepo {
	col := db.Collection("lock")
	name, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"resource": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Errorf("create unique index for lock collection failed: %v", err)
	} else {
		log.Debugf("ensure unique index (%s) for lock collection", name)
	}

	return &LockRepo{col: col}
}
