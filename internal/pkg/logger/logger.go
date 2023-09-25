package logger

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger initializes a new logger.
func newLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	zapCfg := createZapConfig(cfg)

	zapCfg.Level = zap.NewAtomicLevelAt(cfg.LogLevel)
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapCfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to build zap config"))
	}

	return logger, nil
}

// createZapConfig creates a zap.Config based on the LoggerConfig.
func createZapConfig(cfg *LoggerConfig) zap.Config {
	if cfg.Production {
		return zap.NewProductionConfig()
	}
	return zap.NewDevelopmentConfig()
}

func Build() fx.Option {
	return fx.Provide(
		parseConfig,
		newLogger,
	)
}
