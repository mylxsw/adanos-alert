package repository

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrAlreadyLocked is returned by a locking operation when a resource
	// is already locked.
	ErrAlreadyLocked = errors.New("unable to acquire lock (resource is already locked)")

	// ErrLockNotFound is returned when a lock cannot be found.
	ErrLockNotFound = errors.New("unable to find lock")
)

type Lock struct {
	LockID    primitive.ObjectID `bson:"lock_id" json:"lock_id"`
	Resource  string             `bson:"resource" json:"resource"`
	Acquired  bool               `bson:"acquired" json:"acquired"`
	Owner     string             `bson:"owner" json:"owner"`
	TTL       uint               `bson:"ttl" json:"ttl"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	RenewedAt time.Time          `bson:"renewed_at" json:"renewed_at"`
	ExpiredAt time.Time          `bson:"expired_at" json:"expired_at"`
}

type LockRepo interface {
	Lock(resource string, owner string, ttl uint) (*Lock, error)
	Renew(lockID primitive.ObjectID, ttl uint) (*Lock, error)
	UnLock(lockID primitive.ObjectID) error
	Remove(resource string) error
}
