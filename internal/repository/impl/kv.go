package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KVRepo struct {
	col *mongo.Collection
}

func NewKVRepo(db *mongo.Database) repository.KVRepo {
	return &KVRepo{col: db.Collection("kv")}
}

func (repo KVRepo) Set(key string, value interface{}) error {
	return repo.SetWithTTL(key, value, 0)
}

func (repo KVRepo) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	withTTL := false
	if ttl > 0 {
		withTTL = true
	}

	kv := repository.KV{
		Key:       key,
		Value:     value,
		WithTTL:   withTTL,
		ExpiredAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	_, _ = repo.col.DeleteOne(context.TODO(), bson.M{"key": key})
	_, err := repo.col.InsertOne(context.TODO(), kv)
	return err
}

func (repo KVRepo) Get(key string) (pair repository.KV, err error) {
	var kv repository.KV

	if err := repo.col.FindOne(context.TODO(), bson.M{"key": key}).Decode(&kv); err != nil {
		if err == mongo.ErrNoDocuments {
			return pair, repository.ErrNotFound
		}
		return pair, err
	}

	if kv.WithTTL && kv.ExpiredAt.Before(time.Now()) {
		_, _ = repo.Remove(key)
		return kv, repository.ErrNotFound
	}

	return kv, nil
}

func (repo KVRepo) Remove(key string) (removeCount int64, err error) {
	rs, err := repo.col.DeleteOne(context.TODO(), bson.M{"key": key})
	if err != nil {
		return 0, err
	}

	return rs.DeletedCount, nil
}

func (repo KVRepo) All(filter bson.M) (pairs []repository.KV, err error) {
	err = repo.GC()
	if err != nil {
		return
	}

	cur, err := repo.col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var kv repository.KV
		if err := cur.Decode(&kv); err != nil {
			return nil, err
		}

		pairs = append(pairs, kv)
	}

	return pairs, nil
}

func (repo KVRepo) GC() error {
	_, err := repo.col.DeleteMany(context.TODO(), bson.M{"with_ttl": true, "expired_at": bson.M{"$lt": time.Now()}})
	return err
}
