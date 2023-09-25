package grpc

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/fx"
)

// mount configures the lifecycle of the gRPC service to ensure that the connection is
// properly closed when the application is terminated. This function is invoked during
// the application's initialization.
func mount(grpcClient GrpcClient, grpcServer GrpcServer, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := grpcServer.start(ctx); err != nil {
				return fault.Wrap(err, fmsg.With("failed to start gRPC server"))
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			grpcServer.shutdown()

			if err := grpcClient.close(); err != nil {
				return fault.Wrap(err, fmsg.With("failed to properly shutdown gRPC client"))
			}

			return nil
		},
	})
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			parseConfig,
			newServer,
			newClient,
		),
		fx.Invoke(mount),
	)
}
