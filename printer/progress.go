package printer

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/odpf/salt/term"
)

type Spinner struct {
	indicator *spinner.Spinner
}

func (s *Spinner) Stop() {
	if s.indicator == nil {
		return
	}
	s.indicator.Stop()
}

func Progress(label string) *Spinner {
	set := spinner.CharSets[11]
	if !term.IsTTY() {
		return &Spinner{}
	}
	s := spinner.New(set, 120*time.Millisecond, spinner.WithColor("fgCyan"))
	if label != "" {
		s.Prefix = label + " "
	}

	s.Start()

	return &Spinner{s}
}
