package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recovery struct {
	RecoveryID string               `json:"recovery_id" bson:"recovery_id"`
	RecoveryAt time.Time            `json:"recovery_at" bson:"recovery_at"`
	RefIDs     []primitive.ObjectID `json:"ref_ids" bson:"ref_ids"`
	CreatedAt  time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at" json:"updated_at"`
}

type RecoveryRepo interface {
	Register(ctx context.Context, recoveryAt time.Time, recoveryID string, refID primitive.ObjectID) (err error)
	RecoverableMessages(ctx context.Context, deadline time.Time) ([]Recovery, error)
	Delete(ctx context.Context, recoveryID string) error
}
