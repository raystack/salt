package db_test

import (
	"embed"
	"fmt"
	"log"
	"testing"

	"github.com/raystack/salt/db"
	"github.com/stretchr/testify/assert"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

func TestRunMigrations(t *testing.T) {
	if _, err := client.Exec(dropTableQuery); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}

	dsn := fmt.Sprintf(dsn, user, password, port, database)
	var (
		pgConfig = db.Config{
			Driver: "postgres",
			URL:    dsn,
		}
	)

	err := db.RunMigrations(pgConfig, migrationFs, "migrations")
	assert.NoError(t, err)

	// User table should exist
	var tableExist bool
	result := client.QueryRow(checkTableQuery)
	result.Scan(&tableExist)
	assert.Equal(t, true, tableExist)
}

func TestRunRollback(t *testing.T) {
	if _, err := client.Exec(dropTableQuery); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}

	dsn := fmt.Sprintf(dsn, user, password, port, database)
	var (
		pgConfig = db.Config{
			Driver: "postgres",
			URL:    dsn,
		}
	)

	err := db.RunRollback(pgConfig, migrationFs, "migrations")
	assert.NoError(t, err)

	// User table should not exist
	var tableExist bool
	result := client.QueryRow(checkTableQuery)
	result.Scan(&tableExist)
	assert.Equal(t, false, tableExist)
}
