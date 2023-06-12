package printer

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/raystack/salt/term"
)

type Indicator struct {
	spinner *spinner.Spinner
}

func (s *Indicator) Stop() {
	if s.spinner == nil {
		return
	}
	s.spinner.Stop()
}

func Spin(label string) *Indicator {
	set := spinner.CharSets[11]
	if !term.IsTTY() {
		return &Indicator{}
	}
	s := spinner.New(set, 120*time.Millisecond, spinner.WithColor("fgCyan"))
	if label != "" {
		s.Prefix = label + " "
	}

	s.Start()

	return &Indicator{s}
}
