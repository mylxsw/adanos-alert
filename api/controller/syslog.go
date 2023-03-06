package controller

import (
	"fmt"
	"net/http"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
)

type SyslogController struct {
	cc infra.Resolver
}

func NewSyslogController(cc infra.Resolver) web.Controller {
	return &SyslogController{cc: cc}
}

func (u SyslogController) Register(router web.Router) {
	router.Group("/syslog/", func(router web.Router) {
		router.Get("/logs/", u.Logs).Name("syslog:logs")
	})
}

func (u SyslogController) Logs(ctx web.Context, syslogRepo repository.SyslogRepo) web.Response {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	logType := ctx.Input("type")
	if logType != "" {
		filter["type"] = repository.SyslogType(logType)
	}

	data, next, err := syslogRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("query syslog failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"logs": data,
		"next": next,
	})
}
