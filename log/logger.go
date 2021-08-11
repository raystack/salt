package log

import (
	"io"
)

// Option modifies the logger behavior
type Option func(interface{})

// Logger is a convenient interface to use provided loggers
// either use it as it is or implement your own interface where
// the logging implementations are used
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	Level() string
	Writer() io.Writer
}
