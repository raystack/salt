package mux

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/oklog/run"
)

const (
	defaultGracePeriod = 10 * time.Second
)

// Serve starts TCP listeners and serves the registered protocol servers of the
// given serveTarget(s) and blocks until the servers exit. Context can be
// cancelled to perform graceful shutdown.
func Serve(ctx context.Context, opts ...Option) error {
	mux := muxServer{gracePeriod: defaultGracePeriod}
	for _, opt := range opts {
		if err := opt(&mux); err != nil {
			return err
		}
	}

	if len(mux.targets) == 0 {
		return errors.New("mux serve: at least one serve target must be set")
	}

	return mux.Serve(ctx)
}

type muxServer struct {
	targets     []serveTarget
	gracePeriod time.Duration
}

func (mux *muxServer) Serve(ctx context.Context) error {
	var g run.Group
	for _, t := range mux.targets {
		l, err := net.Listen("tcp", t.Address())
		if err != nil {
			return fmt.Errorf("mux serve: %w", err)
		}

		t := t // redeclare to avoid referring to updated value inside closures.
		g.Add(func() error {
			err := t.Serve(l)
			if err != nil {
				log.Print("[ERROR] Serve:", err)
			}
			return err
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), mux.gracePeriod)
			defer cancel()

			if err := t.Shutdown(ctx); err != nil {
				log.Print("[ERROR] Shutdown server gracefully:", err)
			}
		})
	}

	g.Add(func() error {
		<-ctx.Done()
		return ctx.Err()
	}, func(error) {
	})

	return g.Run()
}
