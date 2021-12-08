package log

import (
	"io"

	"go.uber.org/zap"
)

type Zap struct {
	log  *zap.SugaredLogger
	conf zap.Config
}

func (z Zap) Debug(msg string, args ...interface{}) {
	z.log.With(args...).Debug(msg)
}

func (z Zap) Info(msg string, args ...interface{}) {
	z.log.With(args...).Info(msg)
}

func (z Zap) Warn(msg string, args ...interface{}) {
	z.log.With(args...).Warn(msg, args)
}

func (z Zap) Error(msg string, args ...interface{}) {
	z.log.With(args...).Error(msg, args)
}

func (z Zap) Fatal(msg string, args ...interface{}) {
	z.log.With(args...).Fatal(msg, args)
}

func (z Zap) Level() string {
	return z.conf.Level.String()
}

func (z Zap) Writer() io.Writer {
	panic("not supported")
}

func ZapWithConfig(conf zap.Config, opts ...zap.Option) Option {
	return func(z interface{}) {
		z.(*Zap).conf = conf
		prodLogger, err := z.(*Zap).conf.Build(opts...)
		if err != nil {
			panic(err)
		}
		z.(*Zap).log = prodLogger.Sugar()
	}
}

func ZapWithNoop() Option {
	return func(z interface{}) {
		z.(*Zap).log = zap.NewNop().Sugar()
		z.(*Zap).conf = zap.Config{}
	}
}

// NewZap returns a zap logger instance with info level as default log level
func NewZap(opts ...Option) *Zap {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level.SetLevel(zap.InfoLevel)
	logger, err := defaultConfig.Build()
	if err != nil {
		panic(err)
	}

	zapper := &Zap{
		log:  logger.Sugar(),
		conf: defaultConfig,
	}
	for _, opt := range opts {
		opt(zapper)
	}
	return zapper
}
