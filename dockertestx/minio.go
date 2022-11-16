package dockertestx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	defaultMinioRootUser     = "minio"
	defaultMinioRootPassword = "minio123"
	defaultMinioDomain       = "localhost"
)

type dockerMinioOption func(dm *dockerMinio)

// MinioWithDockertestNetwork is an option to assign docker network
func MinioWithDockertestNetwork(network *dockertest.Network) dockerMinioOption {
	return func(dm *dockerMinio) {
		dm.network = network
	}
}

// MinioWithVersionTag is an option to assign version tag
// of a `quay.io/minio/minio` image
func MinioWithVersionTag(versionTag string) dockerMinioOption {
	return func(dm *dockerMinio) {
		dm.versionTag = versionTag
	}
}

// MinioWithDockerPool is an option to assign docker pool
func MinioWithDockerPool(pool *dockertest.Pool) dockerMinioOption {
	return func(dm *dockerMinio) {
		dm.pool = pool
	}
}

type dockerMinio struct {
	network             *dockertest.Network
	pool                *dockertest.Pool
	rootUser            string
	rootPassword        string
	domain              string
	versionTag          string
	internalHost        string
	externalHost        string
	externalConsoleHost string
	dockertestResource  *dockertest.Resource
}

// CreateMinio creates a minio instance with default configurations
func CreateMinio(opts ...dockerMinioOption) (*dockerMinio, error) {
	var (
		err error
		dm  = &dockerMinio{}
	)

	for _, opt := range opts {
		opt(dm)
	}

	name := fmt.Sprintf("minio-%s", uuid.New().String())

	if dm.pool == nil {
		dm.pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dm.rootUser == "" {
		dm.rootUser = defaultMinioRootUser
	}

	if dm.rootPassword == "" {
		dm.rootPassword = defaultMinioRootPassword
	}

	if dm.domain == "" {
		dm.domain = defaultMinioDomain
	}

	if dm.versionTag == "" {
		dm.versionTag = "RELEASE.2022-09-07T22-25-02Z"
	}

	runOpts := &dockertest.RunOptions{
		Name:       name,
		Repository: "quay.io/minio/minio",
		Tag:        dm.versionTag,
		Env: []string{
			"MINIO_ROOT_USER=" + dm.rootUser,
			"MINIO_ROOT_PASSWORD=" + dm.rootPassword,
			"MINIO_DOMAIN=" + dm.domain,
		},
		Cmd:          []string{"server", "/data1", "--console-address", ":9001"},
		ExposedPorts: []string{"9000/tcp", "9001/tcp"},
	}

	if dm.network != nil {
		runOpts.NetworkID = dm.network.Network.ID
	}

	dm.dockertestResource, err = dm.pool.RunWithOptions(
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

	minioPort := dm.dockertestResource.GetPort("9000/tcp")
	minioConsolePort := dm.dockertestResource.GetPort("9001/tcp")

	dm.internalHost = fmt.Sprintf("%s:%s", name, "9000")
	dm.externalHost = fmt.Sprintf("%s:%s", "localhost", minioPort)
	dm.externalConsoleHost = fmt.Sprintf("%s:%s", "localhost", minioConsolePort)

	if err = dm.dockertestResource.Expire(120); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dm.pool.MaxWait = 60 * time.Second

	if err = dm.pool.Retry(func() error {
		httpClient := &http.Client{}
		res, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/minio/health/live", minioPort))
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("minio server return status %d", res.StatusCode)
		}

		return nil
	}); err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dm, nil
}

func (dm *dockerMinio) GetInternalHost() string {
	return dm.internalHost
}

func (dm *dockerMinio) GetExternalHost() string {
	return dm.externalHost
}

func (dm *dockerMinio) GetExternalConsoleHost() string {
	return dm.externalConsoleHost
}

// GetPool returns docker pool
func (dm *dockerMinio) GetPool() *dockertest.Pool {
	return dm.pool
}

// GetResource returns docker resource
func (dm *dockerMinio) GetResource() *dockertest.Resource {
	return dm.dockertestResource
}
