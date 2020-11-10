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

type EventRelationNoteRepo struct {
	col *mongo.Collection
}

func NewEventRelationNoteRepo(db *mongo.Database) repository.EventRelationNoteRepo {
	col := db.Collection("event_relation_note")
	if _, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"relation_id": 1},
		Options: options.Index().SetUnique(false),
	}); err != nil {
		log.Errorf("can not create index for event_relation_note.relation_id: %v", err)
	}

	return &EventRelationNoteRepo{col: col}
}

func (m *EventRelationNoteRepo) AddNote(ctx context.Context, note repository.EventRelationNote) (repository.ID, error) {
	note.ID = primitive.NewObjectID()
	note.CreatedAt = time.Now()
	note.DeletedAt = time.Time{}

	rs, err := m.col.InsertOne(ctx, note)
	if err != nil {
		return "", err
	}

	return repository.ID(rs.InsertedID.(primitive.ObjectID).Hex()), nil
}

func (m *EventRelationNoteRepo) PaginateNotes(ctx context.Context, relID primitive.ObjectID, filter bson.M, offset, limit int64) (notes []repository.EventRelationNote, next int64, err error) {
	notes = make([]repository.EventRelationNote, 0)

	filter["relation_id"] = relID
	cur, err := m.col.Find(ctx, filter, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}).SetSkip(offset))
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var note repository.EventRelationNote
		if err = cur.Decode(&note); err != nil {
			return
		}

		notes = append(notes, note)
	}

	if int64(len(notes)) == limit {
		next = offset + limit
	}

	return notes, next, err
}

func (m *EventRelationNoteRepo) DeleteNote(ctx context.Context, relID primitive.ObjectID, filter bson.M) error {
	filter["relation_id"] = relID
	_, err := m.col.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}
