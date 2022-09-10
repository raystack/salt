package term

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

type Pager struct {
	Out    io.Writer
	ErrOut io.Writer

	pagerCommand string
	pagerProcess *os.Process
}

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

func (s *Pager) Set(cmd string) {
	s.pagerCommand = cmd
}

func (s *Pager) Get() string {
	return s.pagerCommand
}

func (s *Pager) Start() error {
	if s.pagerCommand == "" || s.pagerCommand == "cat" {
		return nil
	}

	pagerArgs, err := shlex.Split(s.pagerCommand)
	if err != nil {
		return err
	}

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

	pagerExe, err := safeexec.LookPath(pagerArgs[0])
	if err != nil {
		return err
	}
	pagerCmd := exec.Command(pagerExe, pagerArgs[1:]...)
	pagerCmd.Env = pagerEnv
	pagerCmd.Stdout = s.Out
	pagerCmd.Stderr = s.ErrOut
	pagedOut, err := pagerCmd.StdinPipe()
	if err != nil {
		return err
	}
	s.Out = &pagerWriter{pagedOut}
	err = pagerCmd.Start()
	if err != nil {
		return err
	}
	s.pagerProcess = pagerCmd.Process
	return nil
}

func (s *Pager) Stop() {
	if s.pagerProcess == nil {
		return
	}

	_ = s.Out.(io.WriteCloser).Close()
	_, _ = s.pagerProcess.Wait()
	s.pagerProcess = nil
}

// pagerWriter implements a WriteCloser that wraps all EPIPE errors in an ErrClosedPagerPipe type.
type pagerWriter struct {
	io.WriteCloser
}

func (w *pagerWriter) Write(d []byte) (int, error) {
	n, err := w.WriteCloser.Write(d)
	if err != nil && (errors.Is(err, io.ErrClosedPipe) || isEpipeError(err)) {
		return n, &ErrClosedPagerPipe{err}
	}
	return n, err
}

func isEpipeError(err error) bool {
	return errors.Is(err, syscall.EPIPE)
}

// ErrClosedPagerPipe is the error returned when writing to a pager that has been closed.
type ErrClosedPagerPipe struct {
	error
}
