package impl

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentRepo struct {
	col *mongo.Collection
}

func NewAgentRepo(db *mongo.Database) repository.AgentRepo {
	col := db.Collection("agent")
	return &AgentRepo{col: col}
}

func (a AgentRepo) Update(agent repository.Agent) (primitive.ObjectID, error) {
	panic("implement me")
}

func (a AgentRepo) Find(filter bson.M) (rules []repository.Agent, err error) {
	panic("implement me")
}

func (a AgentRepo) Delete(filter bson.M) error {
	panic("implement me")
}

func (a AgentRepo) DeleteID(id primitive.ObjectID) error {
	panic("implement me")
}
