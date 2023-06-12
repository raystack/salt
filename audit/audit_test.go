package audit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raystack/salt/audit"
	"github.com/raystack/salt/audit/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuditTestSuite struct {
	suite.Suite

	now time.Time

	mockRepository *mocks.Repository
	service        *audit.Service
}

func (s *AuditTestSuite) setupTest() {
	s.mockRepository = new(mocks.Repository)
	s.service = audit.New(
		audit.WithMetadataExtractor(func(context.Context) map[string]interface{} {
			return map[string]interface{}{
				"trace_id":    "test-trace-id",
				"app_name":    "guardian_test",
				"app_version": 1,
			}
		}),
		audit.WithRepository(s.mockRepository),
	)

	s.now = time.Now()
	audit.TimeNow = func() time.Time {
		return s.now
	}
}

func TestAudit(t *testing.T) {
	suite.Run(t, new(AuditTestSuite))
}

func (s *AuditTestSuite) TestLog() {
	s.Run("should insert to repository", func() {
		s.setupTest()

		s.mockRepository.On("Insert", mock.Anything, &audit.Log{
			Timestamp: s.now,
			Action:    "action",
			Actor:     "user@example.com",
			Data:      map[string]interface{}{"foo": "bar"},
			Metadata: map[string]interface{}{
				"trace_id":    "test-trace-id",
				"app_name":    "guardian_test",
				"app_version": 1,
			},
		}).Return(nil)

		ctx := context.Background()
		ctx = audit.WithActor(ctx, "user@example.com")
		err := s.service.Log(ctx, "action", map[string]interface{}{"foo": "bar"})
		s.NoError(err)
	})

	s.Run("actor extractor", func() {
		s.Run("should use actor extractor if option given", func() {
			expectedActor := "test-actor"
			s.service = audit.New(
				audit.WithActorExtractor(func(ctx context.Context) (string, error) {
					return expectedActor, nil
				}),
				audit.WithRepository(s.mockRepository),
			)

			s.mockRepository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				log := args.Get(1).(*audit.Log)
				s.Equal(expectedActor, log.Actor)
			}).Return(nil).Once()

			err := s.service.Log(context.Background(), "", nil)
			s.NoError(err)
		})

		s.Run("should return error if extractor returns error", func() {
			expectedError := errors.New("test error")
			s.service = audit.New(
				audit.WithActorExtractor(func(ctx context.Context) (string, error) {
					return "", expectedError
				}),
			)

			err := s.service.Log(context.Background(), "", nil)
			s.ErrorIs(err, expectedError)
		})
	})

	s.Run("metadata", func() {
		s.Run("should pass empty trace id if extractor not found", func() {
			s.service = audit.New(
				audit.WithMetadataExtractor(func(ctx context.Context) map[string]interface{} {
					return map[string]interface{}{
						"app_name":    "guardian_test",
						"app_version": 1,
					}
				}),
				audit.WithRepository(s.mockRepository),
			)

			s.mockRepository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				l := args.Get(1).(*audit.Log)
				s.IsType(map[string]interface{}{}, l.Metadata)

				md := l.Metadata.(map[string]interface{})
				s.Empty(md["trace_id"])
				s.NotEmpty(md["app_name"])
				s.NotEmpty(md["app_version"])
			}).Return(nil).Once()

			err := s.service.Log(context.Background(), "", nil)
			s.NoError(err)
		})

		s.Run("should append new metadata to existing one", func() {
			s.service = audit.New(
				audit.WithMetadataExtractor(func(ctx context.Context) map[string]interface{} {
					return map[string]interface{}{
						"existing": "foobar",
					}
				}),
				audit.WithRepository(s.mockRepository),
			)

			expectedMetadata := map[string]interface{}{
				"existing": "foobar",
				"new":      "foobar",
			}
			s.mockRepository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				log := args.Get(1).(*audit.Log)
				s.Equal(expectedMetadata, log.Metadata)
			}).Return(nil).Once()

			ctx, err := audit.WithMetadata(context.Background(), map[string]interface{}{
				"new": "foobar",
			})
			s.Require().NoError(err)

			err = s.service.Log(ctx, "", nil)
			s.NoError(err)
		})
	})

	s.Run("should return error if repository.Insert fails", func() {
		s.setupTest()

		expectedError := errors.New("test error")
		s.mockRepository.On("Insert", mock.Anything, mock.Anything).Return(expectedError)

		err := s.service.Log(context.Background(), "", nil)
		s.ErrorIs(err, expectedError)
	})
}
