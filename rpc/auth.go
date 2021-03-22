package rpc

import (
	"context"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/glacier/infra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// authFunc Rpc权限校验
func authFunc(cc infra.Binder, conf *configs.Config) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return ctx, status.Errorf(codes.Unauthenticated, "auth failed: empty token")
		}

		token := meta.Get("token")
		if len(token) != 1 {
			return ctx, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		if conf.GRPCToken != token[0] {
			return ctx, status.Errorf(codes.Unauthenticated, "auth failed: token not match")
		}

		return ctx, nil
	}
}
