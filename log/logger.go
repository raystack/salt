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
//
//	timeTaken := time.Duration(time.Second * 1)
//	l.Debug("processed request", "time taken", timeTaken)
//
// here key should always be a `string` and value could be of any type as
// long as it is printable.
//
//	l.Info("processed request", "time taken", timeTaken, "started at", startedAt)
type Logger interface {

	// Debug level message with alternating key/value pairs
	// key should be string, value could be anything printable
	Debug(msg string, args ...interface{})

	// Info level message with alternating key/value pairs
	// key should be string, value could be anything printable
	Info(msg string, args ...interface{})

	// Warn level message with alternating key/value pairs
	// key should be string, value could be anything printable
	Warn(msg string, args ...interface{})

	// Error level message with alternating key/value pairs
	// key should be string, value could be anything printable
	Error(msg string, args ...interface{})

	// Fatal level message with alternating key/value pairs
	// key should be string, value could be anything printable
	Fatal(msg string, args ...interface{})

	// Level returns priority level for which this logger will filter logs
	Level() string

	// Writer used to print logs
	Writer() io.Writer
}
