package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateController struct {
	cc *container.Container
}

func NewTemplateController(cc *container.Container) web.Controller {
	return &TemplateController{cc: cc}
}

func (t *TemplateController) Register(router *web.Router) {
	router.Group("/templates/", func(router *web.Router) {
		router.Get("/", t.Templates).Name("template:all")
		router.Post("/", t.Add).Name("template:add")
		router.Get("/{id}/", t.Get).Name("template:one")
		router.Post("/{id}/", t.Update).Name("template:update")
		router.Delete("/{id}/", t.Delete).Name("template:delete")
	})
}

func (t *TemplateController) Get(ctx web.Context, repo repository.TemplateRepo) (*repository.Template, error) {
	templateID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	template, err := repo.Get(templateID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(err, http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &template, nil
}

func (t *TemplateController) Templates(ctx web.Context, repo repository.TemplateRepo) ([]repository.Template, error) {
	filter := bson.M{}

	templateType := ctx.Input("type")
	if templateType != "" {
		filter["type"] = templateType
	}

	return repo.Find(filter)
}

type TemplateForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Type        string `json:"type"`
}

func (t TemplateForm) Validate() error {
	if t.Name == "" {
		return errors.New("invalid argument: name is required")
	}

	if !govalidator.IsIn(
		t.Type,
		string(repository.TemplateTypeMatchRule),
		string(repository.TemplateTypeTriggerRule),
		string(repository.TemplateTypeTemplate),
	) {
		return errors.New("invalid argument: type")
	}

	return nil
}

func (t *TemplateController) Add(ctx web.Context, repo repository.TemplateRepo) (*repository.Template, error) {
	var templateForm TemplateForm
	if err := ctx.Unmarshal(&templateForm); err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	ctx.Validate(templateForm, true)

	id, err := repo.Add(repository.Template{
		Name:        templateForm.Name,
		Description: templateForm.Description,
		Content:     templateForm.Content,
		Type:        repository.TemplateType(templateForm.Type),
	})
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	temp, err := repo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &temp, nil
}

func (t *TemplateController) Update(ctx web.Context, repo repository.TemplateRepo) (*repository.Template, error) {
	templateID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	var templateForm TemplateForm
	if err := ctx.Unmarshal(&templateForm); err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	ctx.Validate(templateForm, true)

	template, err := repo.Get(templateID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(err, http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	if template.Predefined {
		return nil, web.WrapJSONError(errors.New("predefined template is readonly"), http.StatusUnprocessableEntity)
	}

	template.Name = templateForm.Name
	template.Description = templateForm.Description
	template.Content = templateForm.Content
	template.Type = repository.TemplateType(templateForm.Type)

	if err := repo.Update(templateID, template); err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &template, nil
}

func (t *TemplateController) Delete(ctx web.Context, repo repository.TemplateRepo) error {
	templateID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	template, err := repo.Get(templateID)
	if err != nil {
		if err == repository.ErrNotFound {
			return web.WrapJSONError(err, http.StatusNotFound)
		}

		return web.WrapJSONError(err, http.StatusInternalServerError)
	}

	if template.Predefined {
		return web.WrapJSONError(errors.New("predefined template is readonly"), http.StatusUnprocessableEntity)
	}

	if err := repo.DeleteID(templateID); err != nil {
		return web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return nil
}
