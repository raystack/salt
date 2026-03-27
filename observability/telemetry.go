package observability

import (
	"context"
	"time"

	"github.com/raystack/salt/observability/logger"
)

const gracePeriod = 5 * time.Second

type Config struct {
	AppVersion    string
	AppName       string              `yaml:"app_name" mapstructure:"app_name" default:"service"`
	OpenTelemetry OpenTelemetryConfig `yaml:"open_telemetry" mapstructure:"open_telemetry"`
}

func Init(ctx context.Context, cfg Config, logger logger.Logger) (cleanUp func(), err error) {
	shutdown, err := initOTLP(ctx, cfg, logger)
	if err != nil {
		return noOp, err
	}
	return shutdown, nil
}
