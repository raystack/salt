package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/raystack/salt/db"
	"github.com/stretchr/testify/assert"
)

const (
	dialect  = "postgres"
	user     = "root"
	password = "pass"
	database = "postgres"
	host     = "localhost"
	port     = "5432"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

var (
	createTableQuery = "CREATE TABLE IF NOT EXISTS users (id VARCHAR(36) PRIMARY KEY, name VARCHAR(50))"
	dropTableQuery   = "DROP TABLE IF EXISTS users"
	checkTableQuery  = "SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users');"
)

var client *db.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + database,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	fmt.Println(resource.GetPort("5432/tcp"))

	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not expire resource: %s", err.Error())
	}

	pool.MaxWait = 60 * time.Second

	dsn := fmt.Sprintf(dsn, user, password, port, database)
	var (
		pgConfig = db.Config{
			Driver: "postgres",
			URL:    dsn,
		}
	)

	if err = pool.Retry(func() error {
		client, err = db.New(pgConfig)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		client.Close()
	}()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestWithTxn(t *testing.T) {
	if _, err := client.Exec(dropTableQuery); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}
	err := client.WithTxn(context.Background(), sql.TxOptions{}, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(createTableQuery); err != nil {
			return err
		}
		if _, err := tx.Exec(dropTableQuery); err != nil {
			return err
		}

		return nil
	})
	assert.NoError(t, err)

	// Table should be dropped
	var tableExist bool
	result := client.QueryRow(checkTableQuery)
	result.Scan(&tableExist)
	assert.Equal(t, false, tableExist)
}

func TestWithTxnCommit(t *testing.T) {
	if _, err := client.Exec(dropTableQuery); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}
	query2 := "SELECT 1"

	err := client.WithTxn(context.Background(), sql.TxOptions{}, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(createTableQuery); err != nil {
			return err
		}
		if _, err := tx.Exec(query2); err != nil {
			return err
		}

		return nil
	})
	// WithTx should not return an error
	assert.NoError(t, err)

	// User table should exist
	var tableExist bool
	result := client.QueryRow(checkTableQuery)
	result.Scan(&tableExist)
	assert.Equal(t, true, tableExist)
}

func TestWithTxnRollback(t *testing.T) {
	if _, err := client.Exec(dropTableQuery); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}
	query2 := "WRONG QUERY"

	err := client.WithTxn(context.Background(), sql.TxOptions{}, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(createTableQuery); err != nil {
			return err
		}
		if _, err := tx.Exec(query2); err != nil {
			return err
		}

		return nil
	})
	// WithTx should return an error
	assert.Error(t, err)

	// Table should not be created
	var tableExist bool
	result := client.QueryRow(checkTableQuery)
	result.Scan(&tableExist)
	assert.Equal(t, false, tableExist)
}
