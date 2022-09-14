package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logrus struct {
	log *logrus.Logger
}

func (l Logrus) getFields(args ...interface{}) map[string]interface{} {
	fieldMap := map[string]interface{}{}
	if len(args) > 1 && len(args)%2 == 0 {
		for i := 1; i < len(args); i += 2 {
			fieldMap[args[i-1].(string)] = args[i]
		}
	}
	return fieldMap
}

func (l *Logrus) Info(msg string, args ...interface{}) {
	l.log.WithFields(l.getFields(args...)).Info(msg)
}

func (l *Logrus) Debug(msg string, args ...interface{}) {
	l.log.WithFields(l.getFields(args...)).Debug(msg)
}

func (l *Logrus) Warn(msg string, args ...interface{}) {
	l.log.WithFields(l.getFields(args...)).Warn(msg)
}

func (l *Logrus) Error(msg string, args ...interface{}) {
	l.log.WithFields(l.getFields(args...)).Error(msg)
}

func (l *Logrus) Fatal(msg string, args ...interface{}) {
	l.log.WithFields(l.getFields(args...)).Fatal(msg)
}

func (l *Logrus) Level() string {
	return l.log.Level.String()
}

func (l *Logrus) Writer() io.Writer {
	return l.log.Writer()
}

func (l *Logrus) Entry(args ...interface{}) *logrus.Entry {
	return l.log.WithFields(l.getFields(args...))
}

func LogrusWithLevel(level string) Option {
	return func(logger interface{}) {
		logLevel, err := logrus.ParseLevel(level)
		if err != nil {
			panic(err)
		}
		logger.(*Logrus).log.SetLevel(logLevel)
	}
}

func LogrusWithWriter(writer io.Writer) Option {
	return func(logger interface{}) {
		logger.(*Logrus).log.SetOutput(writer)
	}
}

// LogrusWithFormatter can be used to change default formatting
// by implementing logrus.Formatter
// For example:
//
//	type PlainFormatter struct{}
//	func (p *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
//	    return []byte(entry.Message), nil
//	}
//	l := log.NewLogrus(log.LogrusWithFormatter(&PlainFormatter{}))
func LogrusWithFormatter(f logrus.Formatter) Option {
	return func(logger interface{}) {
		logger.(*Logrus).log.SetFormatter(f)
	}
}

// NewLogrus returns a logrus logger instance with info level as default log level
func NewLogrus(opts ...Option) *Logrus {
	logger := &Logrus{
		log: logrus.New(),
	}
	logger.log.Level = logrus.InfoLevel
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}
