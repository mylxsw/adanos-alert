package rpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/graceful"
	"google.golang.org/grpc"
	"net"
)

type Provider struct{}

func (p Provider) Register(app container.Container) {
	app.MustSingleton(func(conf *configs.Config) *grpc.Server {
		auth := authFunc(app, conf)
		return grpc.NewServer(
			grpc_middleware.WithStreamServerChain(
				grpc_auth.StreamServerInterceptor(auth),
				grpc_recovery.StreamServerInterceptor(),
			),
			grpc_middleware.WithUnaryServerChain(
				grpc_auth.UnaryServerInterceptor(auth),
				grpc_recovery.UnaryServerInterceptor(),
			),
		)
	})

}

func (p Provider) Boot(app glacier.Glacier) {
	app.MustResolve(func(serv *grpc.Server) {
		protocol.RegisterMessageServer(serv, NewMessageService(app.Container()))
	})
}

func (p Provider) Daemon(_ context.Context, app glacier.Glacier) {
	app.MustResolve(func(serv *grpc.Server, conf *configs.Config, gf graceful.Graceful) {
		listener, err := net.Listen("tcp", conf.GRPCListen)
		if err != nil {
			panic(fmt.Sprintf("can not create listener for grpc: %v", err))
		}

		gf.AddShutdownHandler(func() {
			serv.GracefulStop()
			log.Debug("grpc server has been stopped")
		})

		log.Debugf("grpc server started, listening on %s", conf.GRPCListen)
		if err := serv.Serve(listener); err != nil {
			log.Errorf("GRPC Server has been stopped: %v", err)
		}
	})
}
