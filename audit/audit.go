//go:generate mockery --name=repository --exported

package audit

import (
	"context"
	"errors"
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

func WithTraceIDExtractor(fn func(ctx context.Context) string) AuditOption {
	return func(s *Service) {
		s.trackIDExtractor = fn
	}
}

type Service struct {
	repository       repository
	trackIDExtractor func(ctx context.Context) string
	withMetadata     func(ctx context.Context) (context.Context, error)
}

func New(opts ...AuditOption) *Service {
	svc := &Service{}
	for _, o := range opts {
		o(svc)
	}

	return svc
}

func (s *Service) Log(ctx context.Context, action string, data interface{}) error {
	if s.withMetadata != nil {
		var err error
		ctx, err = s.withMetadata(ctx)
		if err != nil {
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

	if actor, ok := ctx.Value(actorContextKey{}).(string); ok {
		l.Actor = actor
	}

	return s.repository.Insert(ctx, l)
}
