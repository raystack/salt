package telemetry

import (
	"context"
	"log/slog"
	"time"
)

const gracePeriod = 5 * time.Second

// Config holds the telemetry configuration.
type Config struct {
	AppVersion    string
	AppName       string              `yaml:"app_name" mapstructure:"app_name" default:"service"`
	OpenTelemetry OpenTelemetryConfig `yaml:"open_telemetry" mapstructure:"open_telemetry"`
}

// Init initializes OpenTelemetry and returns a cleanup function.
func Init(ctx context.Context, cfg Config, logger *slog.Logger) (cleanUp func(), err error) {
	shutdown, err := initOTLP(ctx, cfg, logger)
	if err != nil {
		return noOp, err
	}
	return shutdown, nil
}
