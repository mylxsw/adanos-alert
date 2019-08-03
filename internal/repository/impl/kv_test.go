package impl_test

import (
	"context"
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Database() (*mongo.Database, error) {
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	db := conn.Database("adanos_test")
	return db, nil
}

func TestKvRepo(t *testing.T) {
	db, err := Database()
	assert.NoError(t, err)

	kvRepo := impl.NewKVRepo(db)
	assert.NoError(t, kvRepo.Set("environment", "dev"))
	assert.NoError(t, kvRepo.Set("timeout", 3000))
	assert.NoError(t, kvRepo.SetWithTTL("expired", "yes", 1*time.Nanosecond))
	{
		pair, err := kvRepo.Get("environment")
		assert.NoError(t, err)
		assert.Equal(t, "dev", pair.Value)
		assert.False(t, pair.WithTTL)
	}

	{
		assert.NoError(t, kvRepo.SetWithTTL("lock", "locked", time.Nanosecond))
		time.Sleep(2 * time.Nanosecond)
		pair, err := kvRepo.Get("lock")
		assert.Equal(t, repository.ErrNotFound, err)
		assert.Equal(t, "locked", pair.Value)
		assert.True(t, pair.WithTTL)

		pair, err = kvRepo.Get("not_found")
		assert.Equal(t, repository.ErrNotFound, err)
	}

	{
		pairs, err := kvRepo.All(bson.M{})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(pairs))
	}

	{
		removeCount, err := kvRepo.Remove("environment")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), removeCount)
	}

	{
		removeCount, err := kvRepo.Remove("timeout")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), removeCount)
	}

	{
		removeCount, err := kvRepo.Remove("not_found")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), removeCount)
	}
}
