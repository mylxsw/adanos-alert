package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Agent struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	IP          string    `bson:"ip" json:"ip"`
	AgentID     string    `bson:"agent_id" json:"agent_id"`
	Version     string    `bson:"version" json:"version"`
	LastAliveAt time.Time `bson:"last_alive_at" json:"last_alive_at"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (agent Agent) Alive() bool {
	return time.Now().Unix()-agent.LastAliveAt.Unix() < 60
}

type AgentRepo interface {
	Update(agent Agent) (primitive.ObjectID, error)
	Get(id primitive.ObjectID) (agent Agent, err error)
	Find(filter bson.M) (agents []Agent, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
}
