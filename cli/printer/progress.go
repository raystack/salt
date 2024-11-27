package printer

import (
	"github.com/schollz/progressbar/v3"
)

// Progress creates and returns a progress bar for tracking the progress of an operation.
//
// The progress bar supports color, shows a description, and displays the current progress count.
//
// Parameters:
//   - max: The maximum value of the progress bar, indicating 100% completion.
//   - description: A brief description of the task associated with the progress bar.
//
// Returns:
//   - A pointer to a `progressbar.ProgressBar` instance for managing the progress.
//
// Example Usage:
//
//	bar := printer.Progress(100, "Downloading files")
//	for i := 0; i < 100; i++ {
//	    bar.Add(1) // Increment progress by 1.
//	}
func Progress(max int, description string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
	)
	return bar
}
