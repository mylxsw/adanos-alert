package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AgentRepo struct {
	col *mongo.Collection
}

func NewAgentRepo(db *mongo.Database) repository.AgentRepo {
	col := db.Collection("agent")
	return &AgentRepo{col: col}
}

func (a AgentRepo) Update(agent repository.Agent) (repository.ID, error) {
	if agent.AgentID == "" {
		return "", errors.New("agent_id is required")
	}

	agents, err := a.Find(bson.M{"agent_id": agent.AgentID})
	if err != nil {
		return "", err
	}

	if len(agents) == 0 {
		agent.CreatedAt = time.Now()
		agent.UpdatedAt = agent.CreatedAt

		rs, err := a.col.InsertOne(context.TODO(), agent)
		if err != nil {
			return "", err
		}

		return repository.ID(rs.InsertedID.(primitive.ObjectID).Hex()), nil
	}
	agent.CreatedAt = agents[0].CreatedAt
	agent.UpdatedAt = time.Now()

	_, err = a.col.ReplaceOne(context.TODO(), bson.M{"_id": agents[0].ID}, agent)
	return repository.ID(agents[0].ID.Hex()), err
}

func (a AgentRepo) Get(id repository.ID) (agent repository.Agent, err error) {
	mID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return agent, fmt.Errorf("invalid id: %w", err)
	}

	err = a.col.FindOne(context.TODO(), bson.M{"_id": mID}).Decode(&agent)
	return
}

func (a AgentRepo) Find(filter bson.M) (agents []repository.Agent, err error) {
	agents = make([]repository.Agent, 0)
	cur, err := a.col.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return
	}
	defer func() {
		_ = cur.Close(context.TODO())
	}()

	for cur.Next(context.TODO()) {
		var agent repository.Agent
		if err = cur.Decode(&agent); err != nil {
			return
		}

		agents = append(agents, agent)
	}

	return
}

func (a AgentRepo) Delete(filter bson.M) error {
	_, err := a.col.DeleteMany(context.TODO(), filter)
	return err
}

func (a AgentRepo) DeleteID(id repository.ID) error {
	mID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	return a.Delete(bson.M{"_id": mID})
}
