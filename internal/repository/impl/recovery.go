package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type RecoveryRepo struct {
	col *mongo.Collection
}

func NewRecoveryRepo(db *mongo.Database) repository.RecoveryRepo {
	return &RecoveryRepo{col: db.Collection("recovery")}
}

func (r RecoveryRepo) Register(ctx context.Context, recoveryAt time.Time, recoveryID string, refID primitive.ObjectID) error {
	rec, err := r.get(ctx, recoveryID)
	if err != nil {
		if err == repository.ErrNotFound {
			_, err := r.col.InsertOne(ctx, repository.Recovery{
				RecoveryID: recoveryID,
				RecoveryAt: recoveryAt,
				RefIDs:     misc.IfElse(refID != primitive.NilObjectID, []primitive.ObjectID{refID}, []primitive.ObjectID{}).([]primitive.ObjectID),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	if refID != primitive.NilObjectID {
		rec.RefIDs = append(rec.RefIDs, refID)
	}

	rec.UpdatedAt = time.Now()
	rec.RecoveryAt = recoveryAt

	_, err = r.col.ReplaceOne(ctx, bson.M{"recovery_id": recoveryID}, rec)
	return err
}

func (r RecoveryRepo) get(ctx context.Context, recoveryID string) (rec repository.Recovery, err error) {
	err = r.col.FindOne(ctx, bson.M{"recovery_id": recoveryID}).Decode(&rec)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (r RecoveryRepo) RecoverableEvents(ctx context.Context, deadline time.Time) ([]repository.Recovery, error) {
	results := make([]repository.Recovery, 0)
	cursor, err := r.col.Find(ctx, bson.M{"recovery_at": bson.M{"$lt": deadline}})
	if err != nil {
		return results, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(ctx) {
		var rec repository.Recovery
		if err = cursor.Decode(&rec); err != nil {
			log.Errorf("decode recovery message from mongodb failed: %v", err)
			continue
		}

		results = append(results, rec)
	}

	return results, nil
}

func (r RecoveryRepo) Delete(ctx context.Context, recoveryID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"recovery_id": recoveryID})
	return err
}
