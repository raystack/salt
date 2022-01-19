package term

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	indicator *spinner.Spinner
}

func Spin(label string) *Spinner {
	set := spinner.CharSets[14]
	if !IsTTY() {
		return &Spinner{}
	}
	s := spinner.New(set, 100*time.Millisecond)
	s.Prefix = label + " "
	s.Start()

	return &Spinner{s}
}

func (s *Spinner) Stop() {
	if s.indicator == nil {
		return
	}
	s.indicator.Stop()
}
