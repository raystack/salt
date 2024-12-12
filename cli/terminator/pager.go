package terminator

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/cli/safeexec"
	"github.com/google/shlex"
)

// Pager manages a pager process for displaying output in a paginated format.
//
// It supports configuring the pager command, starting the pager process,
// and ensuring proper cleanup when the pager is no longer needed.
type Pager struct {
	Out          io.Writer   // The writer to send output to the pager.
	ErrOut       io.Writer   // The writer to send error output to.
	pagerCommand string      // The command to run the pager (e.g., "less", "more").
	pagerProcess *os.Process // The running pager process, if any.
}

// NewPager creates a new Pager instance with default settings.
//
// If the "PAGER" environment variable is not set, the default command is "more".
func NewPager() *Pager {
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "more"
	}

	return &Pager{
		pagerCommand: pagerCmd,
		Out:          os.Stdout,
		ErrOut:       os.Stderr,
	}
}

// Set updates the pager command used to display output.
//
// Parameters:
//   - cmd: The pager command (e.g., "less", "more").
func (p *Pager) Set(cmd string) {
	p.pagerCommand = cmd
}

// Get returns the current pager command.
//
// Returns:
//   - The pager command as a string.
func (p *Pager) Get() string {
	return p.pagerCommand
}

// Start begins the pager process to display output.
//
// If the pager command is "cat" or empty, it does nothing.
// The function also sets environment variables to optimize the behavior of
// certain pagers, like "less" and "lv".
//
// Returns:
//   - An error if the pager command fails to start or if arguments cannot be parsed.
func (p *Pager) Start() error {
	if p.pagerCommand == "" || p.pagerCommand == "cat" {
		return nil
	}

	pagerArgs, err := shlex.Split(p.pagerCommand)
	if err != nil {
		return err
	}

	// Prepare the environment variables for the pager process.
	pagerEnv := os.Environ()
	for i := len(pagerEnv) - 1; i >= 0; i-- {
		if strings.HasPrefix(pagerEnv[i], "PAGER=") {
			pagerEnv = append(pagerEnv[0:i], pagerEnv[i+1:]...)
		}
	}
	if _, ok := os.LookupEnv("LESS"); !ok {
		pagerEnv = append(pagerEnv, "LESS=FRX")
	}
	if _, ok := os.LookupEnv("LV"); !ok {
		pagerEnv = append(pagerEnv, "LV=-c")
	}

	// Locate the pager executable using safeexec for added security.
	pagerExe, err := safeexec.LookPath(pagerArgs[0])
	if err != nil {
		return err
	}

	pagerCmd := exec.Command(pagerExe, pagerArgs[1:]...)
	pagerCmd.Env = pagerEnv
	pagerCmd.Stdout = p.Out
	pagerCmd.Stderr = p.ErrOut
	pagedOut, err := pagerCmd.StdinPipe()
	if err != nil {
		return err
	}
	p.Out = &pagerWriter{pagedOut}

	// Start the pager process.
	err = pagerCmd.Start()
	if err != nil {
		return err
	}
	p.pagerProcess = pagerCmd.Process
	return nil
}

// Stop terminates the running pager process and cleans up resources.
func (p *Pager) Stop() {
	if p.pagerProcess == nil {
		return
	}

	// Close the output writer and wait for the process to exit.
	_ = p.Out.(io.WriteCloser).Close()
	_, _ = p.pagerProcess.Wait()
	p.pagerProcess = nil
}

// pagerWriter is a custom writer that wraps WriteCloser and handles EPIPE errors.
//
// If a write fails due to a closed pipe, it returns an ErrClosedPagerPipe error.
type pagerWriter struct {
	io.WriteCloser
}

// Write writes data to the underlying WriteCloser and handles EPIPE errors.
//
// Parameters:
//   - d: The data to write.
//
// Returns:
//   - The number of bytes written and an error if the write fails.
func (w *pagerWriter) Write(d []byte) (int, error) {
	n, err := w.WriteCloser.Write(d)
	if err != nil && (errors.Is(err, io.ErrClosedPipe) || isEpipeError(err)) {
		return n, &ErrClosedPagerPipe{err}
	}
	return n, err
}

// isEpipeError checks if an error is a broken pipe (EPIPE) error.
//
// Returns:
//   - A boolean indicating whether the error is an EPIPE error.
func isEpipeError(err error) bool {
	return errors.Is(err, syscall.EPIPE)
}

// ErrClosedPagerPipe is an error type returned when writing to a closed pager pipe.
type ErrClosedPagerPipe struct {
	error
}
