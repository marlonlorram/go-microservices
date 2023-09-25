package grpc

import (
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient represents a gRPC client interface.
type GrpcClient interface {
	Connection() *grpc.ClientConn
	close() error
}

// grpcClient implements the GrpcClient interface.
type grpcClient struct {
	conn *grpc.ClientConn
}

// newClient initializes a new gRPC client.
func newClient(l *zap.Logger, cfg *GrpcConfig) (GrpcClient, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fault.Wrap(err,
			fmsg.With("failed to dial gRPC server"),
			fmsg.WithDesc("endpoint", endpoint))
	}

	l.Info("gRPC client connected", zap.String("endpoint", endpoint))
	return &grpcClient{conn}, nil
}

// Connection returns the gRPC client connection.
func (g *grpcClient) Connection() *grpc.ClientConn {
	return g.conn
}

// close closes the gRPC client connection.
func (g *grpcClient) close() error {
	return g.conn.Close()
}
