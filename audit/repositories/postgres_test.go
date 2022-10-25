package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx/types"
	"github.com/odpf/salt/audit"
	"github.com/odpf/salt/audit/repositories"
	"github.com/odpf/salt/log"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite

	repository     *repositories.PostgresRepository
	pool           *dockertest.Pool
	dockerResource *dockertest.Resource
}

func TestPostgresRepository(t *testing.T) {
	suite.Run(t, new(PostgresRepositoryTestSuite))
}

func (s *PostgresRepositoryTestSuite) SetupSuite() {
	var err error
	s.repository, s.pool, s.dockerResource, err = newTestRepository(log.NewLogrus())
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PostgresRepositoryTestSuite) TearDownSuite() {
	if err := s.repository.DB().Close(); err != nil {
		s.T().Fatal(err)
	}
	if err := purgeTestDocker(s.pool, s.dockerResource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PostgresRepositoryTestSuite) TestInit() {
	s.Run("should migrate audit log model", func() {
		expoectedColumns := []string{
			"timestamp", "action", "actor", "data", "metadata",
		}
		rows, err := s.repository.DB().Query("SELECT * FROM audit_logs")
		s.NoError(err)

		defer rows.Close()
		actualRows, err := rows.Columns()
		s.NoError(err)

		s.Equal(expoectedColumns, actualRows)
	})
}

func (s *PostgresRepositoryTestSuite) TestInsert() {
	s.Run("should insert record to db", func() {
		l := &audit.Log{
			Timestamp: time.Now(),
			Action:    "test-action",
			Actor:     "user@example.com",
			Data:      types.JSONText(`{"test": "data"}`),
			Metadata:  types.JSONText(`{"test": "metadata"}`),
		}

		err := s.repository.Insert(context.Background(), l)
		s.Require().NoError(err)

		rows, err := s.repository.DB().Query("SELECT * FROM audit_logs")
		var actualResult repositories.AuditPostgresModel
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

	s.Run("should return error if data marshaling returns error", func() {
		l := &audit.Log{
			Data: make(chan int),
		}

		err := s.repository.Insert(context.Background(), l)
		s.EqualError(err, "marshaling data: json: unsupported type: chan int")
	})

	s.Run("should return error if metadata marshaling returns error", func() {
		l := &audit.Log{
			Metadata: map[string]interface{}{
				"foo": make(chan int),
			},
		}

		err := s.repository.Insert(context.Background(), l)
		s.EqualError(err, "marshaling metadata: json: unsupported type: chan int")
	})
}
