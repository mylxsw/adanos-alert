package controller

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
)

type AgentController struct {
	cc infra.Resolver
}

func NewAgentController(cc infra.Resolver) web.Controller {
	return &AgentController{cc: cc}
}

func (c AgentController) Register(router web.Router) {
	router.Group("/agents", func(router web.Router) {
		router.Get("/", c.All).Name("agents:all")
		router.Delete("/{id}/", c.Remove).Name("agents:delete")
	})
}

type AgentResp struct {
	repository.Agent
	Alive bool `json:"alive"`
}

func (c AgentController) All(repo repository.AgentRepo) ([]AgentResp, error) {
	agents, err := repo.Find(bson.M{})
	if err != nil {
		return nil, err
	}

	var results []AgentResp
	if err := coll.MustNew(agents).Map(func(agent repository.Agent) AgentResp {
		return AgentResp{
			Agent: agent,
			Alive: agent.Alive(),
		}
	}).All(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func (c AgentController) Remove(req web.Request, repo repository.AgentRepo) error {
	return repo.DeleteID(repository.ID(req.PathVar("id")))
}
