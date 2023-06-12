package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx/types"
	"github.com/raystack/salt/audit"
	"github.com/raystack/salt/audit/repositories"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite

	repository *repositories.PostgresRepository
}

func TestPostgresRepository(t *testing.T) {
	suite.Run(t, new(PostgresRepositoryTestSuite))
}

func (s *PostgresRepositoryTestSuite) SetupSuite() {
	var err error
	repository, pool, dockerResource, err := newTestRepository(log.NewLogrus())
	if err != nil {
		s.T().Fatal(err)
	}
	s.repository = repository

	s.T().Cleanup(func() {
		if err := s.repository.DB().Close(); err != nil {
			s.T().Fatal(err)
		}
		if err := purgeTestDocker(pool, dockerResource); err != nil {
			s.T().Fatal(err)
		}
	})
}

func (s *PostgresRepositoryTestSuite) TestInsert() {
	s.Run("should insert record to db", func() {
		l := &audit.Log{
			Timestamp: time.Now(),
			Action:    "test-action",
			Actor:     "user@example.com",
			Data: types.NullJSONText{
				JSONText: []byte(`{"test": "data"}`),
				Valid:    true,
			},
			Metadata: types.NullJSONText{
				JSONText: []byte(`{"test": "metadata"}`),
				Valid:    true,
			},
		}

		err := s.repository.Insert(context.Background(), l)
		s.Require().NoError(err)

		rows, err := s.repository.DB().Query("SELECT * FROM audit_logs")
		var actualResult repositories.AuditModel
		for rows.Next() {
			err := rows.Scan(&actualResult.Timestamp, &actualResult.Action, &actualResult.Actor, &actualResult.Data, &actualResult.Metadata)
			s.Require().NoError(err)
		}

		s.NoError(err)
		s.NotNil(actualResult)
		if diff := cmp.Diff(l.Timestamp, actualResult.Timestamp, cmpopts.EquateApproxTime(time.Microsecond)); diff != "" {
			s.T().Errorf("result not match, diff: %v", diff)
		}
		s.Equal(l.Action, actualResult.Action)
		s.Equal(l.Actor, actualResult.Actor)
		s.Equal(l.Data, actualResult.Data)
		s.Equal(l.Metadata, actualResult.Metadata)
	})

	s.Run("should return error if data marshalling returns error", func() {
		l := &audit.Log{
			Data: make(chan int),
		}

		err := s.repository.Insert(context.Background(), l)
		s.EqualError(err, "marshalling data: json: unsupported type: chan int")
	})

	s.Run("should return error if metadata marshalling returns error", func() {
		l := &audit.Log{
			Metadata: map[string]interface{}{
				"foo": make(chan int),
			},
		}

		err := s.repository.Insert(context.Background(), l)
		s.EqualError(err, "marshalling metadata: json: unsupported type: chan int")
	})
}
