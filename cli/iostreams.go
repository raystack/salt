package cli

import (
	"bytes"
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/raystack/salt/cli/printer"
	"github.com/raystack/salt/cli/prompt"
	"github.com/raystack/salt/cli/terminal"
	"golang.org/x/term"
)

// IOStreams provides centralized access to standard I/O streams and
// terminal capabilities for CLI commands.
//
// Use [System] for production and [Test] for tests. Commands access
// it via [IO]:
//
//	ios := cli.IO(cmd)
//	if !ios.CanPrompt() {
//	    return fmt.Errorf("--yes required in non-interactive mode")
//	}
type IOStreams struct {
	In     io.ReadCloser // standard input
	Out    io.Writer     // standard output (may become pager pipe)
	ErrOut io.Writer     // standard error

	inTTY  bool
	outTTY bool
	errTTY bool

	colorEnabled bool
	neverPrompt  bool

	pager        *terminal.Pager
	pagerStarted bool
	origOut      io.Writer

	// lazily created
	output   *printer.Output
	prompter prompt.Prompter
}

// System creates IOStreams wired to the real terminal.
func System() *IOStreams {
	outTTY := isTTY(os.Stdout)
	return &IOStreams{
		In:           os.Stdin,
		Out:          os.Stdout,
		ErrOut:       os.Stderr,
		inTTY:        isTTY(os.Stdin),
		outTTY:       outTTY,
		errTTY:       isTTY(os.Stderr),
		colorEnabled: outTTY && !termenv.EnvNoColor(),
	}
}

// Test creates IOStreams backed by buffers for deterministic testing.
// All TTY flags default to false and color is disabled.
func Test() (ios *IOStreams, stdin *bytes.Buffer, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	stdin = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	ios = &IOStreams{
		In:     io.NopCloser(stdin),
		Out:    stdout,
		ErrOut: stderr,
	}
	return
}

// IsStdinTTY reports whether standard input is a terminal.
func (s *IOStreams) IsStdinTTY() bool { return s.inTTY }

// IsStdoutTTY reports whether standard output is a terminal.
func (s *IOStreams) IsStdoutTTY() bool { return s.outTTY }

// IsStderrTTY reports whether standard error is a terminal.
func (s *IOStreams) IsStderrTTY() bool { return s.errTTY }

// SetStdinTTY overrides the stdin TTY flag (useful in tests).
func (s *IOStreams) SetStdinTTY(v bool) { s.inTTY = v }

// SetStdoutTTY overrides the stdout TTY flag (useful in tests).
func (s *IOStreams) SetStdoutTTY(v bool) { s.outTTY = v; s.output = nil }

// SetStderrTTY overrides the stderr TTY flag (useful in tests).
func (s *IOStreams) SetStderrTTY(v bool) { s.errTTY = v }

// SetColorEnabled overrides color detection (useful in tests).
func (s *IOStreams) SetColorEnabled(v bool) { s.colorEnabled = v }

// SetNeverPrompt disables interactive prompting regardless of TTY state.
func (s *IOStreams) SetNeverPrompt(v bool) { s.neverPrompt = v }

// ColorEnabled reports whether color output is active.
func (s *IOStreams) ColorEnabled() bool { return s.colorEnabled }

// CanPrompt reports whether interactive prompting is possible.
// Returns false if prompting is disabled, or stdin/stdout are not terminals.
func (s *IOStreams) CanPrompt() bool {
	return !s.neverPrompt && s.inTTY && s.outTTY
}

// TerminalWidth returns the terminal width in columns.
// Returns 80 if the width cannot be determined.
func (s *IOStreams) TerminalWidth() int {
	if f, ok := s.Out.(*os.File); ok {
		if w, _, err := term.GetSize(int(f.Fd())); err == nil && w > 0 {
			return w
		}
	}
	return 80
}

// StartPager starts a pager process and redirects Out through it.
// Does nothing if stdout is not a TTY or no pager command is configured.
func (s *IOStreams) StartPager() error {
	if !s.outTTY {
		return nil
	}
	p := terminal.NewPager()
	p.Out = s.Out
	p.ErrOut = s.ErrOut
	if err := p.Start(); err != nil {
		return err
	}
	s.origOut = s.Out
	s.Out = p.Out
	s.pager = p
	s.pagerStarted = true
	s.output = nil // invalidate cached Output
	return nil
}

// StopPager stops the pager process and restores the original Out.
func (s *IOStreams) StopPager() {
	if s.pager != nil && s.pagerStarted {
		s.pager.Stop()
		s.pagerStarted = false
		if s.origOut != nil {
			s.Out = s.origOut
			s.origOut = nil
			s.output = nil // invalidate cached Output
		}
	}
}

// Output returns the formatting layer, creating it lazily.
func (s *IOStreams) Output() *printer.Output {
	if s.output == nil {
		s.output = printer.NewOutputFrom(s.Out, s.ErrOut, s.outTTY)
	}
	return s.output
}

// Prompter returns the prompt layer, creating it lazily.
func (s *IOStreams) Prompter() prompt.Prompter {
	if s.prompter == nil {
		s.prompter = prompt.New()
	}
	return s.prompter
}

func isTTY(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
