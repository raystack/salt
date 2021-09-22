package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(ctx context.Context) context.Context {
	newCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-newCtx.Done()
		stop()
	}()
	return newCtx
}
