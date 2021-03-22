package controller

import (
	"fmt"
	"net/http"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
)

type AuditController struct {
	cc infra.Resolver
}

func NewAuditController(cc infra.Resolver) web.Controller {
	return &AuditController{cc: cc}
}

func (u AuditController) Register(router web.Router) {
	router.Group("/audit/", func(router web.Router) {
		router.Get("/logs/", u.Logs).Name("audit:logs")
	})
}

func (u AuditController) Logs(ctx web.Context, auditRepo repository.AuditLogRepo) web.Response {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	logType := ctx.Input("type")
	if logType != "" {
		filter["type"] = repository.AuditLogType(logType)
	}

	data, next, err := auditRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("query audit logs failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"logs": data,
		"next": next,
	})
}
