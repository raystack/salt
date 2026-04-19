package cli

import (
	"errors"
	"fmt"
	"os"
)

// ErrSilent indicates the command already printed its error.
// The error handler should exit 1 without printing anything.
var ErrSilent = errors.New("silent error")

// ErrCancel indicates the user cancelled the operation (e.g. ctrl-c).
// The error handler should exit 0.
var ErrCancel = errors.New("cancelled")

// FlagError wraps an error caused by bad flags or arguments.
// The error handler should print the error and show usage.
type FlagError struct {
	Err error
}

func (e *FlagError) Error() string { return e.Err.Error() }
func (e *FlagError) Unwrap() error { return e.Err }

// NewFlagError creates a FlagError.
func NewFlagError(err error) *FlagError {
	return &FlagError{Err: err}
}

// HandleError handles a command error by type.
// FlagError prints the error (usage is shown by cobra).
// SilentError exits without printing.
// CancelError exits with code 0.
// Other errors print "Error: <message>".
func HandleError(err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, ErrCancel):
		os.Exit(0)
	case errors.Is(err, ErrSilent):
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
