//go:generate mockery --name=repository --exported

package audit

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	TimeNow = time.Now

	ErrInvalidMetadata = errors.New("failed to cast existing metadata to map[string]interface{} type")
)

type actorContextKey struct{}
type metadataContextKey struct{}

func WithActor(ctx context.Context, actor string) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func WithMetadata(ctx context.Context, md map[string]interface{}) (context.Context, error) {
	existingMetadata := ctx.Value(metadataContextKey{})
	if existingMetadata == nil {
		return context.WithValue(ctx, metadataContextKey{}, md), nil
	}

	// append new metadata
	mapMd, ok := existingMetadata.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidMetadata
	}
	for k, v := range md {
		mapMd[k] = v
	}

	return context.WithValue(ctx, metadataContextKey{}, mapMd), nil
}

type repository interface {
	Init(context.Context) error
	Insert(context.Context, *Log) error
}

type AuditOption func(*Service)

func WithRepository(r repository) AuditOption {
	return func(s *Service) {
		s.repository = r
	}
}

func WithMetadataExtractor(fn func(context.Context) map[string]interface{}) AuditOption {
	return func(s *Service) {
		s.withMetadata = func(ctx context.Context) (context.Context, error) {
			md := fn(ctx)
			return WithMetadata(ctx, md)
		}
	}
}

func WithActorExtractor(fn func(context.Context) (string, error)) AuditOption {
	return func(s *Service) {
		s.actorExtractor = fn
	}
}

func defaultActorExtractor(ctx context.Context) (string, error) {
	if actor, ok := ctx.Value(actorContextKey{}).(string); ok {
		return actor, nil
	}
	return "", nil
}

type Service struct {
	repository     repository
	actorExtractor func(context.Context) (string, error)
	withMetadata   func(context.Context) (context.Context, error)
}

func New(opts ...AuditOption) *Service {
	svc := &Service{
		actorExtractor: defaultActorExtractor,
	}
	for _, o := range opts {
		o(svc)
	}

	return svc
}

func (s *Service) Log(ctx context.Context, action string, data interface{}) error {
	if s.withMetadata != nil {
		var err error
		if ctx, err = s.withMetadata(ctx); err != nil {
			return err
		}
	}

	l := &Log{
		Timestamp: TimeNow(),
		Action:    action,
		Data:      data,
	}

	if md, ok := ctx.Value(metadataContextKey{}).(map[string]interface{}); ok {
		l.Metadata = md
	}

	if s.actorExtractor != nil {
		actor, err := s.actorExtractor(ctx)
		if err != nil {
			return fmt.Errorf("extracting actor: %w", err)
		}
		l.Actor = actor
	}

	return s.repository.Insert(ctx, l)
}
