package db_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/odpf/salt/db"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const (
	dialect  = "postgres"
	user     = "postgres"
	password = "pass"
	database = "postgres"
	host     = "localhost"
	port     = "5432"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + database,
		},
	}

	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not expire resource: %s", err.Error())
	}

	pool.MaxWait = 60 * time.Second

	var client *db.Client

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
