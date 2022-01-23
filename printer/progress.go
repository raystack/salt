package printer

import (
	"github.com/schollz/progressbar/v3"
)

func Progress(max int, description string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
	)
	return bar
}
