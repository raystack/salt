package printer

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/raystack/salt/cli/terminator"
)

// Indicator represents a terminal spinner used for indicating progress or ongoing operations.
type Indicator struct {
	spinner *spinner.Spinner // The spinner instance.
}

// Stop halts the spinner animation.
//
// This method ensures the spinner is stopped gracefully. If the spinner is nil (e.g., when the
// terminal does not support TTY), the method does nothing.
//
// Example Usage:
//
//	indicator := printer.Spin("Loading")
//	// Perform some operation...
//	indicator.Stop()
func (s *Indicator) Stop() {
	if s.spinner == nil {
		return
	}
	s.spinner.Stop()
}

// Spin creates and starts a terminal spinner to indicate an ongoing operation.
//
// The spinner uses a predefined character set and updates at a fixed interval. It automatically
// disables itself if the terminal does not support TTY.
//
// Parameters:
//   - label: A string to prefix the spinner (e.g., "Loading").
//
// Returns:
//   - An *Indicator instance that manages the spinner lifecycle.
//
// Example Usage:
//
//	indicator := printer.Spin("Processing data")
//	// Perform some long-running operation...
//	indicator.Stop()
func Spin(label string) *Indicator {
	// Predefined spinner character set (dots style).
	set := spinner.CharSets[11]

	// Check if the terminal supports TTY; if not, return a no-op Indicator.
	if !terminator.IsTTY() {
		return &Indicator{}
	}

	// Create a new spinner instance with a 120ms update interval and cyan color.
	s := spinner.New(set, 120*time.Millisecond, spinner.WithColor("fgCyan"))

	// Add a label prefix if provided.
	if label != "" {
		s.Prefix = label + " "
	}

	// Start the spinner animation.
	s.Start()

	// Return the Indicator wrapping the spinner instance.
	return &Indicator{s}
}
