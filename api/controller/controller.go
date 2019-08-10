package controller

import (
	"github.com/mylxsw/hades"
)

type Controller interface {
	Register(router *hades.Router)
}
