package dockertestx

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/raystack/salt/log"
)

const (
	defaultPGUname  = "test_user"
	defaultPGPasswd = "test_pass"
	defaultDBname   = "test_db"
)

type dockerPostgresOption func(dpg *dockerPostgres)

func PostgresWithLogger(logger log.Logger) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.logger = logger
	}
}

// PostgresWithDockertestNetwork is an option to assign docker network
func PostgresWithDockertestNetwork(network *dockertest.Network) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.network = network
	}
}

// PostgresWithDockertestResourceExpiry is an option to assign docker resource expiry time
func PostgresWithDockertestResourceExpiry(expiryInSeconds uint) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.expiryInSeconds = expiryInSeconds
	}
}

// PostgresWithDetail is an option to assign custom details
// like username, password, and database name
func PostgresWithDetail(
	username string,
	password string,
	dbName string,
) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.username = username
		dpg.password = password
		dpg.dbName = dbName
	}
}

// PostgresWithVersionTag is an option to assign version tag
// of a `postgres` image
func PostgresWithVersionTag(versionTag string) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.versionTag = versionTag
	}
}

// PostgresWithDockerPool is an option to assign docker pool
func PostgresWithDockerPool(pool *dockertest.Pool) dockerPostgresOption {
	return func(dpg *dockerPostgres) {
		dpg.pool = pool
	}
}

type dockerPostgres struct {
	logger             log.Logger
	network            *dockertest.Network
	pool               *dockertest.Pool
	username           string
	password           string
	dbName             string
	versionTag         string
	connStringInternal string
	connStringExternal string
	expiryInSeconds    uint
	dockertestResource *dockertest.Resource
}

// CreatePostgres creates a postgres instance with default configurations
func CreatePostgres(opts ...dockerPostgresOption) (*dockerPostgres, error) {
	var (
		err error
		dpg = &dockerPostgres{}
	)

	for _, opt := range opts {
		opt(dpg)
	}

	name := fmt.Sprintf("postgres-%s", uuid.New().String())

	if dpg.pool == nil {
		dpg.pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dpg.username == "" {
		dpg.username = defaultPGUname
	}

	if dpg.password == "" {
		dpg.password = defaultPGPasswd
	}

	if dpg.dbName == "" {
		dpg.dbName = defaultDBname
	}

	if dpg.versionTag == "" {
		dpg.versionTag = "12"
	}

	if dpg.expiryInSeconds == 0 {
		dpg.expiryInSeconds = 120
	}

	runOpts := &dockertest.RunOptions{
		Name:       name,
		Repository: "postgres",
		Tag:        dpg.versionTag,
		Env: []string{
			"POSTGRES_PASSWORD=" + dpg.password,
			"POSTGRES_USER=" + dpg.username,
			"POSTGRES_DB=" + dpg.dbName,
		},
		ExposedPorts: []string{"5432/tcp"},
	}

	if dpg.network != nil {
		runOpts.NetworkID = dpg.network.Network.ID
	}

	dpg.dockertestResource, err = dpg.pool.RunWithOptions(
		runOpts,
		func(config *docker.HostConfig) {
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		return nil, err
	}

	pgPort := dpg.dockertestResource.GetPort("5432/tcp")
	dpg.connStringInternal = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dpg.username, dpg.password, name, "5432", dpg.dbName)
	dpg.connStringExternal = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dpg.username, dpg.password, "localhost", pgPort, dpg.dbName)

	if err = dpg.dockertestResource.Expire(dpg.expiryInSeconds); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dpg.pool.MaxWait = 60 * time.Second

	if err = dpg.pool.Retry(func() error {
		if _, err := sqlx.Connect("postgres", dpg.connStringExternal); err != nil {
			return err
		}
		return nil
	}); err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dpg, nil
}

// GetInternalConnString returns internal connection string of a postgres instance
func (dpg *dockerPostgres) GetInternalConnString() string {
	return dpg.connStringInternal
}

// GetExternalConnString returns external connection string of a postgres instance
func (dpg *dockerPostgres) GetExternalConnString() string {
	return dpg.connStringExternal
}

// GetPool returns docker pool
func (dpg *dockerPostgres) GetPool() *dockertest.Pool {
	return dpg.pool
}

// GetResource returns docker resource
func (dpg *dockerPostgres) GetResource() *dockertest.Resource {
	return dpg.dockertestResource
}
