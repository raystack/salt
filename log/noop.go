package log

import (
	"io"
	"io/ioutil"
)

type Noop struct{}

func (n *Noop) Info(msg string, args ...interface{})  {}
func (n *Noop) Debug(msg string, args ...interface{}) {}
func (n *Noop) Warn(msg string, args ...interface{})  {}
func (n *Noop) Error(msg string, args ...interface{}) {}
func (n *Noop) Fatal(msg string, args ...interface{}) {}

func (n *Noop) Level() string {
	return "unsupported"
}
func (n *Noop) Writer() io.Writer {
	return ioutil.Discard
}

// NewNoop returns a no operation logger, useful in tests
func NewNoop(opts ...Option) *Noop {
	return &Noop{}
}
