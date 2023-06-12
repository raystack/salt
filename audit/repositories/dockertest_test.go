package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/raystack/salt/audit/repositories"
	"github.com/raystack/salt/log"
)

func newTestRepository(logger log.Logger) (*repositories.PostgresRepository, *dockertest.Pool, *dockertest.Resource, error) {
	host := "localhost"
	port := "5433"
	user := "test_user"
	password := "test_pass"
	dbName := "test_db"
	sslMode := "disable"

	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + user,
			"POSTGRES_DB=" + dbName,
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not create dockertest pool: %w", err)
	}

	resource, err := pool.RunWithOptions(opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	port = resource.GetPort("5432/tcp")

	// attach terminal logger to container if exists
	// for debugging purpose
	if logger.Level() == "debug" {
		logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
			Container:    resource.Container.ID,
			OutputStream: logger.Writer(),
			ErrorStream:  logger.Writer(),
			Stderr:       true,
			Stdout:       true,
			Stream:       true,
		})
		if err != nil {
			logger.Fatal("could not connect to postgres container log output", "error", err)
		}
		defer func() {
			if err = logWaiter.Close(); err != nil {
				logger.Fatal("could not close container log", "error", err)
			}

			if err = logWaiter.Wait(); err != nil {
				logger.Fatal("could not wait for container log to close", "error", err)
			}
		}()
	}

	// Tell docker to hard kill the container in 120 seconds
	if err := resource.Expire(120); err != nil {
		return nil, nil, nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 60 * time.Second

	var repo *repositories.PostgresRepository
	time.Sleep(5 * time.Second)
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode))
		if err != nil {
			return err
		}
		repo = repositories.NewPostgresRepository(db)

		return db.Ping()
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	if err := setup(repo); err != nil {
		logger.Fatal("failed to setup and migrate DB", "error", err)
	}
	return repo, pool, resource, nil
}

func setup(repo *repositories.PostgresRepository) error {
	var queries = []string{
		"DROP SCHEMA public CASCADE",
		"CREATE SCHEMA public",
	}
	for _, query := range queries {
		repo.DB().Exec(query)
	}

	if err := repo.Init(context.Background()); err != nil {
		return err
	}

	return nil
}

func purgeTestDocker(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %w", err)
	}
	return nil
}
