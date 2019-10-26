package api

import (
	"errors"

	"github.com/mylxsw/adanos-alert/api/controller"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
)

func routers(cc *container.Container) func(router *hades.Router, mw hades.RequestMiddleware) {
	conf := cc.MustGet(&configs.Config{}).(*configs.Config)
	return func(router *hades.Router, mw hades.RequestMiddleware) {
		mws := make([]hades.HandlerDecorator, 0)
		mws = append(mws, mw.AccessLog(), mw.CORS("*"))
		if conf.APIToken != "" {
			authMiddleware := mw.AuthHandler(func(typ string, credential string) error {
				if typ != "Bearer" {
					return errors.New("invalid auth type, only support Bearer")
				}

				if credential != conf.APIToken {
					return errors.New("token not match")
				}

				return nil
			})

			mws = append(mws, authMiddleware)
		}

		router.WithMiddleware(mws...).Controllers(
			"/api",
			controller.NewWelcomeController(cc),
			controller.NewMessageController(cc),
		)
	}
}