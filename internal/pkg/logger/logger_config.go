package logger

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap/zapcore"

	"github.com/marlonlorram/go-microservices/internal/pkg/config"
)

// LoggerConfig represents environment variable configuration parameters for logging.
type LoggerConfig struct {
	Production bool          `envconfig:"PRODUCTION" default:"false"`
	LogLevel   zapcore.Level `envconfig:"LOG_LEVEL"  default:"info"`
}

// parseConfig parses the environment variables and returns a configured instance of LoggerConfig.
// Returns an error if the configuration could not be loaded.
func parseConfig() (*LoggerConfig, error) {
	cfg, err := config.Parse[LoggerConfig]()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse logger configuration"))
	}

	return cfg, nil
}
