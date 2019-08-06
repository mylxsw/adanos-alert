package controller

import (
	"github.com/mylxsw/go-toolkit/web"
)

type Controller interface {
	Register(router *web.Router)
}
