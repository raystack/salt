package log

import (
	"io"
)

// Option modifies the logger behavior
type Option func(interface{})

// Logger is a convenient interface to use provided loggers
// either use it as it is or implement your own interface where
// the logging implementations are used
// Each log method must take first string as message and then one or
// more key,value arguments.
// For example:
//     timeTaken := time.Duration(time.Second * 1)
//     l.Debug("processed request", "time taken", timeTaken)
// here key should always be a `string` and value could be of any type as
// long as it is printable.
//     l.Info("processed request", "time taken", timeTaken, "started at", startedAt)
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	Level() string
	Writer() io.Writer
}
