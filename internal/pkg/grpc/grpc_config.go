package grpc

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"

	"github.com/marlonlorram/go-microservices/internal/pkg/config"
)

// GrpcConfig represents environment variable configuration parameters for grpc.
type GrpcConfig struct {
	Host string `envconfig:"GRPC_HOST" required:"true" desc:"The hostname or IP address of the gRPC server"`
	Port string `envconfig:"GRPC_PORT" required:"true" desc:"The port number on which the gRPC server is listening"`
}

// parseConfig parses the environment variables and returns a configured instance of GrpcConfig.
// Returns an error if the configuration could not be loaded.
func parseConfig() (*GrpcConfig, error) {
	cfg, err := config.Parse[GrpcConfig]()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse grpc configuration"))
	}

	return cfg, nil
}
