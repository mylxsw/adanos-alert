package impl_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLockRepo_Lock(t *testing.T) {
	db, err := Database()
	assert.NoError(t, err)

	lockRepo := impl.NewLockRepo(db)

	resource := "test-lock"

	// check renew unknown lock
	_, err = lockRepo.Renew(primitive.NewObjectID(), 33)
	assert.Equal(t, repository.ErrLockNotFound, err)

	// check unlock unknown lock
	err = lockRepo.UnLock(primitive.NewObjectID())
	assert.Equal(t, repository.ErrLockNotFound, err)

	// aquire a lock
	lock, err := lockRepo.Lock(resource, "mylxsw", 300)
	assert.NoError(t, err)

	defer func(lockID primitive.ObjectID) {
		// release a lock
		err := lockRepo.UnLock(lockID)
		assert.NoError(t, err)

		err = lockRepo.Remove(resource)
		assert.NoError(t, err)
	}(lock.LockID)

	// renew a lock
	_, err = lockRepo.Renew(lock.LockID, 200)
	assert.NoError(t, err)

	// acquire a lock has been acquired by others
	_, err = lockRepo.Lock(resource, "zhangsan", 300)
	assert.Equal(t, repository.ErrAlreadyLocked, err)
}
