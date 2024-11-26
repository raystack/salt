package dockertestx

import (
	"context"
	"fmt"
	"time"

	authzedpb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultPreSharedKey = "default-preshared-key"
	defaultLogLevel     = "debug"
)

type dockerSpiceDBOption func(dsp *dockerSpiceDB)

func SpiceDBWithLogLevel(logLevel string) dockerSpiceDBOption {
	return func(dsp *dockerSpiceDB) {
		dsp.logLevel = logLevel
	}
}

// SpiceDBWithDockertestNetwork is an option to assign docker network
func SpiceDBWithDockertestNetwork(network *dockertest.Network) dockerSpiceDBOption {
	return func(dsp *dockerSpiceDB) {
		dsp.network = network
	}
}

// SpiceDBWithVersionTag is an option to assign version tag
// of a `quay.io/authzed/spicedb` image
func SpiceDBWithVersionTag(versionTag string) dockerSpiceDBOption {
	return func(dsp *dockerSpiceDB) {
		dsp.versionTag = versionTag
	}
}

// SpiceDBWithDockerPool is an option to assign docker pool
func SpiceDBWithDockerPool(pool *dockertest.Pool) dockerSpiceDBOption {
	return func(dsp *dockerSpiceDB) {
		dsp.pool = pool
	}
}

// SpiceDBWithPreSharedKey is an option to assign pre-shared-key
func SpiceDBWithPreSharedKey(preSharedKey string) dockerSpiceDBOption {
	return func(dsp *dockerSpiceDB) {
		dsp.preSharedKey = preSharedKey
	}
}

type dockerSpiceDB struct {
	network            *dockertest.Network
	pool               *dockertest.Pool
	preSharedKey       string
	versionTag         string
	logLevel           string
	externalPort       string
	dockertestResource *dockertest.Resource
}

// CreateSpiceDB creates a spicedb instance with postgres backend and default configurations
func CreateSpiceDB(postgresConnectionURL string, opts ...dockerSpiceDBOption) (*dockerSpiceDB, error) {
	var (
		err error
		dsp = &dockerSpiceDB{}
	)

	for _, opt := range opts {
		opt(dsp)
	}

	name := fmt.Sprintf("spicedb-%s", uuid.New().String())

	if dsp.pool == nil {
		dsp.pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dsp.preSharedKey == "" {
		dsp.preSharedKey = defaultPreSharedKey
	}

	if dsp.logLevel == "" {
		dsp.logLevel = defaultLogLevel
	}

	if dsp.versionTag == "" {
		dsp.versionTag = "v1.0.0"
	}

	runOpts := &dockertest.RunOptions{
		Name:         name,
		Repository:   "quay.io/authzed/spicedb",
		Tag:          dsp.versionTag,
		Cmd:          []string{"spicedb", "serve", "--log-level", dsp.logLevel, "--grpc-preshared-key", dsp.preSharedKey, "--grpc-no-tls", "--datastore-engine", "postgres", "--datastore-conn-uri", postgresConnectionURL},
		ExposedPorts: []string{"50051/tcp"},
	}

	if dsp.network != nil {
		runOpts.NetworkID = dsp.network.Network.ID
	}

	dsp.dockertestResource, err = dsp.pool.RunWithOptions(
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

	dsp.externalPort = dsp.dockertestResource.GetPort("50051/tcp")

	if err = dsp.dockertestResource.Expire(120); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dsp.pool.MaxWait = 60 * time.Second

	if err = dsp.pool.Retry(func() error {
		client, err := authzed.NewClient(
			fmt.Sprintf("localhost:%s", dsp.externalPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpcutil.WithInsecureBearerToken(dsp.preSharedKey),
		)
		if err != nil {
			return err
		}
		_, err = client.ReadSchema(context.Background(), &authzedpb.ReadSchemaRequest{})
		grpCStatus := status.Convert(err)
		if grpCStatus.Code() == codes.Unavailable {
			return err
		}
		return nil
	}); err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dsp, nil
}

// GetExternalPort returns exposed port of the spicedb instance
func (dsp *dockerSpiceDB) GetExternalPort() string {
	return dsp.externalPort
}

// GetPreSharedKey returns pre-shared-key used in the spicedb instance
func (dsp *dockerSpiceDB) GetPreSharedKey() string {
	return dsp.preSharedKey
}

// GetPool returns docker pool
func (dsp *dockerSpiceDB) GetPool() *dockertest.Pool {
	return dsp.pool
}

// GetResource returns docker resource
func (dsp *dockerSpiceDB) GetResource() *dockertest.Resource {
	return dsp.dockertestResource
}
