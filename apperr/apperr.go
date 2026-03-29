package apperr

import (
	"fmt"
	"runtime"
)

// Component defines where the error occurred
type Component string

const (
	DataLayer Component = "DataLayer"
	LogicLayer Component = "LogicLayer"
	APILayer Component = "APILayer"
)

// AppError is our custom error type
type AppError struct {
	Component     Component // Where it failed
	PublicMessage string    // Safe to show to the end-user
	OriginalErr   error     // The actual error that triggered this
	File          string    // Traceback: File name
	Line          int       // Traceback: Line number
}

// New creates a new AppError and captures the caller's traceback automatically
func New(comp Component, publicMsg string, err error) *AppError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	return &AppError{
		Component:     comp,
		PublicMessage: publicMsg,
		OriginalErr:   err,
		File:          file,
		Line:          line,
	}
}

func (e *AppError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("[%s] %s:%d - %s: %v", e.Component, e.File, e.Line, e.PublicMessage, e.OriginalErr)
	}
	return fmt.Sprintf("[%s] %s:%d - %s", e.Component, e.File, e.Line, e.PublicMessage)
}

func (e *AppError) Unwrap() error {
	return e.OriginalErr
}

// ClientError provides the safe message for the API/CLI response.
func (e *AppError) ClientError() string {
	return e.PublicMessage
}