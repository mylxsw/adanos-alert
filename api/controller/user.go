package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	cc *container.Container
}

func NewUserController(cc *container.Container) hades.Controller {
	return &UserController{cc: cc}
}

func (u UserController) Register(router *hades.Router) {
	router.Group("/users/", func(router *hades.Router) {
		router.Get("/", u.Users).Name("users:all")
		router.Post("/", u.Add).Name("users:add")
		router.Get("/{id}/", u.User).Name("users:one")
		router.Delete("/{id}/", u.Delete).Name("users:delete")
	})
}

type UserForm struct {
	Name   string                `json:"name"`
	Metas  []repository.UserMeta `json:"metas"`
	Status string                `json:"status"`
}

func (userForm UserForm) Validate() error {
	if userForm.Name == "" {
		return errors.New("invalid argument [name]")
	}

	if userForm.Status != "" && !govalidator.IsIn(
		userForm.Status,
		string(repository.UserStatusDisabled),
		string(repository.UserStatusEnabled),
	) {
		return errors.New("invalid argument [status]")
	}

	return nil
}

func (u UserController) Add(ctx hades.Context, userRepo repository.UserRepo) (*repository.User, error) {
	var userForm UserForm
	if err := ctx.Unmarshal(&userForm); err != nil {
		return nil, hades.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	ctx.Validate(userForm, true)

	id, err := userRepo.Add(repository.User{
		Name:   userForm.Name,
		Metas:  userForm.Metas,
		Status: repository.UserStatus(userForm.Status),
	})
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	user, err := userRepo.Get(id)
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (u UserController) Delete(ctx hades.Context, userRepo repository.UserRepo) error {
	userID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return hades.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return userRepo.DeleteID(userID)
}

func (u UserController) User(ctx hades.Context, userRepo repository.UserRepo) (*repository.User, error) {
	userID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, hades.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	user, err := userRepo.Get(userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, hades.WrapJSONError(errors.New("no such user"), http.StatusNotFound)
		}

		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (u UserController) Users(ctx hades.Context, userRepo repository.UserRepo) hades.Response {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	name := ctx.Input("name")
	if name != "" {
		filter["name"] = name
	}

	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	meta := ctx.Input("meta")
	if meta != "" {
		filter["metas.value"] = meta
	}

	users, next, err := userRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("query failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.JSON(hades.M{
		"users": users,
		"next":  next,
	})
}
