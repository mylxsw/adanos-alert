package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	cc *container.Container
}

func NewUserController(cc *container.Container) web.Controller {
	return &UserController{cc: cc}
}

func (u UserController) Register(router *web.Router) {
	router.Group("/users/", func(router *web.Router) {
		router.Get("/", u.Users).Name("users:all")
		router.Post("/", u.Add).Name("users:add")
		router.Post("/{id}/", u.Update).Name("users:update")
		router.Get("/{id}/", u.User).Name("users:one")
		router.Delete("/{id}/", u.Delete).Name("users:delete")
	})

	router.Group("/users-helper/", func(router *web.Router) {
		router.Get("/names/", u.UserNames).Name("users-helper:names")
	})
}

type UserForm struct {
	Name   string                `json:"name"`
	Metas  []repository.UserMeta `json:"metas"`
	Status string                `json:"status"`
}

func (userForm UserForm) GetMetas() []repository.UserMeta {
	results := make([]repository.UserMeta, 0)
	for _, v := range userForm.Metas {
		if strings.TrimSpace(v.Key) != "" {
			results = append(results, v)
		}
	}

	return results
}

func (userForm UserForm) Validate() error {
	if userForm.Name == "" {
		return errors.New("invalid argument: name is required")
	}

	if userForm.Status != "" && !govalidator.IsIn(
		userForm.Status,
		string(repository.UserStatusDisabled),
		string(repository.UserStatusEnabled),
	) {
		return errors.New("invalid argument: status must be enabled/disabled")
	}

	return nil
}

type UserNameResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserNames return all user names only
func (u UserController) UserNames(ctx web.Context, userRepo repository.UserRepo) ([]UserNameResp, error) {
	users, err := userRepo.Find(bson.M{})
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	resps := make([]UserNameResp, 0)
	for _, u := range users {
		resps = append(resps, UserNameResp{
			ID:   u.ID.Hex(),
			Name: u.Name,
		})
	}

	return resps, nil
}

func (u UserController) Add(ctx web.Context, userRepo repository.UserRepo) (*repository.User, error) {
	var userForm UserForm
	if err := ctx.Unmarshal(&userForm); err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	ctx.Validate(userForm, true)

	id, err := userRepo.Add(repository.User{
		Name:   userForm.Name,
		Metas:  userForm.GetMetas(),
		Status: repository.UserStatus(userForm.Status),
	})
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	user, err := userRepo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (u UserController) Update(ctx web.Context, userRepo repository.UserRepo) (*repository.User, error) {
	userID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	var userForm UserForm
	if err := ctx.Unmarshal(&userForm); err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	ctx.Validate(userForm, true)

	user, err := userRepo.Get(userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(err, http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	user.Name = userForm.Name
	user.Metas = userForm.GetMetas()
	user.Status = repository.UserStatus(userForm.Status)

	if err := userRepo.Update(userID, user); err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (u UserController) Delete(ctx web.Context, userRepo repository.UserRepo) error {
	userID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return userRepo.DeleteID(userID)
}

func (u UserController) User(ctx web.Context, userRepo repository.UserRepo) (*repository.User, error) {
	userID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	user, err := userRepo.Get(userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(errors.New("no such user"), http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (u UserController) Users(ctx web.Context, userRepo repository.UserRepo) web.Response {
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

	return ctx.JSON(web.M{
		"users": users,
		"next":  next,
	})
}
