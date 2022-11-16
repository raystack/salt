package dockertestx

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type dockerCortexOption func(dc *dockerCortex)

// CortexWithDockertestNetwork is an option to assign docker network
func CortexWithDockertestNetwork(network *dockertest.Network) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.network = network
	}
}

// CortexWithDockertestNetwork is an option to assign version tag
// of a `quay.io/cortexproject/cortex` image
func CortexWithVersionTag(versionTag string) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.versionTag = versionTag
	}
}

// CortexWithDockerPool is an option to assign docker pool
func CortexWithDockerPool(pool *dockertest.Pool) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.pool = pool
	}
}

// CortexWithModule is an option to assign cortex module name
// e.g. all, alertmanager, querier, etc
func CortexWithModule(moduleName string) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.moduleName = moduleName
	}
}

// CortexWithAlertmanagerURL is an option to assign alertmanager url
func CortexWithAlertmanagerURL(amURL string) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.alertManagerURL = amURL
	}
}

// CortexWithS3Endpoint is an option to assign external s3/minio storage
func CortexWithS3Endpoint(s3URL string) dockerCortexOption {
	return func(dc *dockerCortex) {
		dc.s3URL = s3URL
	}
}

type dockerCortex struct {
	network            *dockertest.Network
	pool               *dockertest.Pool
	moduleName         string
	alertManagerURL    string
	s3URL              string
	internalHost       string
	externalHost       string
	versionTag         string
	dockertestResource *dockertest.Resource
}

// CreateCortex is a function to create a dockerized single-process cortex with
// s3/minio as the backend storage
func CreateCortex(opts ...dockerCortexOption) (*dockerCortex, error) {
	var (
		err error
		dc  = &dockerCortex{}
	)

	for _, opt := range opts {
		opt(dc)
	}

	name := fmt.Sprintf("cortex-%s", uuid.New().String())

	if dc.pool == nil {
		dc.pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dc.versionTag == "" {
		dc.versionTag = "master-63703f5"
	}

	if dc.moduleName == "" {
		dc.moduleName = "all"
	}

	runOpts := &dockertest.RunOptions{
		Name:       name,
		Repository: "quay.io/cortexproject/cortex",
		Tag:        dc.versionTag,
		Env: []string{
			"minio_host=siren_nginx_1",
		},
		Cmd: []string{
			fmt.Sprintf("-target=%s", dc.moduleName),
			"-config.file=/etc/single-process-config.yaml",
			fmt.Sprintf("-ruler.storage.s3.endpoint=%s", dc.s3URL),
			fmt.Sprintf("-ruler.alertmanager-url=%s", dc.alertManagerURL),
			fmt.Sprintf("-alertmanager.storage.s3.endpoint=%s", dc.s3URL),
		},
		ExposedPorts: []string{"9009/tcp"},
		ExtraHosts: []string{
			"cortex.siren_nginx_1:127.0.0.1",
		},
	}

	if dc.network != nil {
		runOpts.NetworkID = dc.network.Network.ID
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var (
		rulesFolder        = fmt.Sprintf("%s/tmp/dockertest-configs/cortex/rules", pwd)
		alertManagerFolder = fmt.Sprintf("%s/tmp/dockertest-configs/cortex/alertmanager", pwd)
	)

	foldersPath := []string{rulesFolder, alertManagerFolder}
	for _, fp := range foldersPath {
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if err := os.MkdirAll(fp, 0777); err != nil {
				return nil, err
			}
		}
	}

	_, thisFileName, _, ok := runtime.Caller(0)
	if !ok {
		return nil, err
	}
	thisFileFolder := path.Dir(thisFileName)

	dc.dockertestResource, err = dc.pool.RunWithOptions(
		runOpts,
		func(config *docker.HostConfig) {
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			config.Mounts = []docker.HostMount{
				{
					Target: "/etc/single-process-config.yaml",
					Source: fmt.Sprintf("%s/configs/cortex/single_process_cortex.yaml", thisFileFolder),
					Type:   "bind",
				},
				{
					Target: "/tmp/cortex/rules",
					Source: rulesFolder,
					Type:   "bind",
				},
				{
					Target: "/tmp/cortex/alertmanager",
					Source: alertManagerFolder,
					Type:   "bind",
				},
			}
		},
	)
	if err != nil {
		return nil, err
	}

	externalPort := dc.dockertestResource.GetPort("9009/tcp")
	dc.internalHost = fmt.Sprintf("%s:9009", name)
	dc.externalHost = fmt.Sprintf("localhost:%s", externalPort)

	if err = dc.dockertestResource.Expire(120); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dc.pool.MaxWait = 60 * time.Second

	if err = dc.pool.Retry(func() error {
		httpClient := &http.Client{}
		res, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/config", externalPort))
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("cortex server return status %d", res.StatusCode)
		}

		return nil
	}); err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dc, nil
}

// GetInternalHost returns internal hostname and port
// e.g. internal-xxxxxx:8080
func (dc *dockerCortex) GetInternalHost() string {
	return dc.internalHost
}

// GetExternalHost returns localhost and port
// e.g. localhost:51113
func (dc *dockerCortex) GetExternalHost() string {
	return dc.externalHost
}

// GetPool returns docker pool
func (dc *dockerCortex) GetPool() *dockertest.Pool {
	return dc.pool
}

// GetResource returns docker resource
func (dc *dockerCortex) GetResource() *dockertest.Resource {
	return dc.dockertestResource
}
