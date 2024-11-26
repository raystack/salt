package telemetry

import (
	"context"
	"time"

	"github.com/raystack/salt/log"
)

const gracePeriod = 5 * time.Second

type Config struct {
	AppVersion    string
	AppName       string              `yaml:"app_name" mapstructure:"app_name" default:"service"`
	OpenTelemetry OpenTelemetryConfig `yaml:"open_telemetry" mapstructure:"open_telemetry"`
}

func Init(ctx context.Context, cfg Config, logger log.Logger) (cleanUp func(), err error) {
	shutdown, err := initOTLP(ctx, cfg, logger)
	if err != nil {
		return noOp, err
	}
	return shutdown, nil
}
