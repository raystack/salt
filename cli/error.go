package cli

import "errors"

// ErrSilent indicates the command already printed its error.
// Execute will exit 1 without printing anything.
var ErrSilent = errors.New("silent error")

// ErrCancel indicates the user cancelled the operation (e.g. ctrl-c).
// Execute will exit 0 without printing anything.
var ErrCancel = errors.New("cancelled")

// flagError wraps an error caused by bad flags or arguments.
// Execute prints the error and shows the failing command's usage.
type flagError struct {
	err error
}

func (e *flagError) Error() string { return e.err.Error() }
func (e *flagError) Unwrap() error { return e.err }
