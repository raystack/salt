package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignals wraps context so that it is marked done
// when one of SIGINT or SIGTERM is received
func HandleSignals(ctx context.Context) context.Context {
	newCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-newCtx.Done()
		stop()
	}()
	return newCtx
}
