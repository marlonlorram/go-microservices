package grpc

import (
	"context"
	"net"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"

	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConnectionIdleTime = 5 * time.Minute
	gRPCTimeout           = 15 * time.Second
	maxConnectionAge      = 5 * time.Minute
	gRPCLifeTime          = 10 * time.Minute
)

// GrpcServer represents a gRPC server interface.
type GrpcServer interface {
	start(context.Context) error
	shutdown()
}

// grpcServer implements the GrpcServer interface.
type grpcServer struct {
	server *googleGrpc.Server
	config *GrpcConfig
	logger *zap.Logger
}

// newServer initializes a new gRPC server.
func newServer(l *zap.Logger, cfg *GrpcConfig) GrpcServer {
	s := &grpcServer{
		server: setupGRPCServer(),
		config: cfg,
		logger: l,
	}

	return s
}

// setupGRPCServer sets up and returns a new gRPC server with the specified options.
func setupGRPCServer() *googleGrpc.Server {
	unaryInterceptors := []googleGrpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
	}
	streamInterceptors := []googleGrpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
	}

	opts := []googleGrpc.ServerOption{
		googleGrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdleTime,
			Timeout:           gRPCTimeout,
			MaxConnectionAge:  maxConnectionAge,
			Time:              gRPCLifeTime,
		}),
		googleGrpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(streamInterceptors...)),
		googleGrpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(unaryInterceptors...)),
	}

	return googleGrpc.NewServer(opts...)
}

// start starts the gRPC server.
func (s *grpcServer) start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return fault.Wrap(err,
			fmsg.With("failed to listen"),
			fmsg.WithDesc("port", s.config.Port))
	}

	go func() {
		if err := s.server.Serve(listener); err != nil {
			s.logger.Fatal("failed to start grpc server", zap.Error(err))
		}
	}()

	s.logger.Info("grpc server started", zap.String("port", s.config.Port))
	return nil
}

// shutdown stops the gRPC server gracefully.
func (s *grpcServer) shutdown() {
	s.server.GracefulStop()
	s.logger.Info("grpc server stopped gracefully.")
}
