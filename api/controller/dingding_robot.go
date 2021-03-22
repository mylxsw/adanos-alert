package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DingdingRobotController struct {
	cc infra.Resolver
}

func NewDingdingRobotController(cc infra.Resolver) web.Controller {
	return &DingdingRobotController{cc: cc}
}

func (u DingdingRobotController) Register(router web.Router) {
	router.Group("/dingding-robots/", func(router web.Router) {
		router.Get("/", u.DingdingRobots).Name("robots:all")
		router.Post("/", u.Add).Name("robots:add")
		router.Post("/{id}/", u.Update).Name("robots:update")
		router.Get("/{id}/", u.DingdingRobot).Name("robots:one")
		router.Delete("/{id}/", u.Delete).Name("robots:delete")
	})

	router.Group("/dingding-robots-helper/", func(router web.Router) {
		router.Get("/names/", u.DingdingRobotNames).Name("robots-helper:names")
	})
}

type DingdingRobotForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Token       string `json:"token"`
	Secret      string `json:"secret"`
}

func (robotForm *DingdingRobotForm) Validate(req web.Request) error {
	if robotForm.Name == "" {
		return errors.New("invalid argument: name is required")
	}

	if len(robotForm.Token) < 32 {
		return errors.New("invalid argument: token is invalid")
	}

	return nil
}

type DingdingRobotNameResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DingdingRobotNames return all robot names only
func (u DingdingRobotController) DingdingRobotNames(ctx web.Context, robotRepo repository.DingdingRobotRepo) ([]DingdingRobotNameResp, error) {
	robots, err := robotRepo.Find(bson.M{})
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	resps := make([]DingdingRobotNameResp, 0)
	for _, u := range robots {
		resps = append(resps, DingdingRobotNameResp{
			ID:          u.ID.Hex(),
			Name:        u.Name,
			Description: u.Description,
		})
	}

	return resps, nil
}

func (u DingdingRobotController) Add(ctx web.Context, em event.Manager, robotRepo repository.DingdingRobotRepo) (*repository.DingdingRobot, error) {
	var robotForm *DingdingRobotForm
	if err := ctx.Unmarshal(&robotForm); err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	ctx.Validate(robotForm, true)

	robot := repository.DingdingRobot{
		Name:        robotForm.Name,
		Description: robotForm.Description,
		Token:       robotForm.Token,
		Secret:      robotForm.Secret,
	}

	id, err := robotRepo.Add(robot)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	em.Publish(pubsub.DingdingRobotEvent{
		DingDingRobot: robot,
		Type:          pubsub.EventTypeAdd,
		CreatedAt:     time.Now(),
	})

	robot, err = robotRepo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &robot, nil
}

func (u DingdingRobotController) Update(ctx web.Context, em event.Manager, robotRepo repository.DingdingRobotRepo) (*repository.DingdingRobot, error) {
	robotID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	var robotForm *DingdingRobotForm
	if err := ctx.Unmarshal(&robotForm); err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	ctx.Validate(robotForm, true)

	robot, err := robotRepo.Get(robotID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(err, http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	robot.Name = robotForm.Name
	robot.Description = robotForm.Description
	robot.Token = robotForm.Token
	robot.Secret = robotForm.Secret

	if err := robotRepo.Update(robotID, robot); err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	em.Publish(pubsub.DingdingRobotEvent{
		DingDingRobot: robot,
		Type:          pubsub.EventTypeUpdate,
		CreatedAt:     time.Now(),
	})

	return &robot, nil
}

func (u DingdingRobotController) Delete(ctx web.Context, em event.Manager, robotRepo repository.DingdingRobotRepo) error {
	robotID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	robot, err := robotRepo.Get(robotID)
	if err != nil {
		return err
	}

	em.Publish(pubsub.DingdingRobotEvent{
		DingDingRobot: robot,
		Type:          pubsub.EventTypeDelete,
		CreatedAt:     time.Now(),
	})

	return robotRepo.DeleteID(robotID)
}

func (u DingdingRobotController) DingdingRobot(ctx web.Context, robotRepo repository.DingdingRobotRepo) (*repository.DingdingRobot, error) {
	robotID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	robot, err := robotRepo.Get(robotID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(errors.New("no such robot"), http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &robot, nil
}

func (u DingdingRobotController) DingdingRobots(ctx web.Context, robotRepo repository.DingdingRobotRepo) web.Response {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	name := ctx.Input("name")
	if name != "" {
		filter["name"] = bson.M{"$regex": name}
	}

	robots, next, err := robotRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("query failed: %v", err), http.StatusInternalServerError)
	}

	for i, robot := range robots {
		robots[i].Secret = "secreted"
		robots[i].Token = template.StringMask(robot.Token, 12)
	}

	return ctx.JSON(web.M{
		"robots": robots,
		"next":   next,
		"search": web.M{
			"name": name,
		},
	})
}
