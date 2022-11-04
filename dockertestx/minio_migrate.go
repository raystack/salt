package dockertestx

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const waitContainerTimeout = 60 * time.Second

type dockerMigrateMinioOption func(dmm *dockerMigrateMinio)

// MigrateMinioWithDockertestNetwork is an option to assign docker network
func MigrateMinioWithDockertestNetwork(network *dockertest.Network) dockerMigrateMinioOption {
	return func(dm *dockerMigrateMinio) {
		dm.network = network
	}
}

// MigrateMinioWithVersionTag is an option to assign version tag
// of a `minio/mc` image
func MigrateMinioWithVersionTag(versionTag string) dockerMigrateMinioOption {
	return func(dm *dockerMigrateMinio) {
		dm.versionTag = versionTag
	}
}

// MigrateMinioWithDockerPool is an option to assign docker pool
func MigrateMinioWithDockerPool(pool *dockertest.Pool) dockerMigrateMinioOption {
	return func(dm *dockerMigrateMinio) {
		dm.pool = pool
	}
}

type dockerMigrateMinio struct {
	network    *dockertest.Network
	pool       *dockertest.Pool
	versionTag string
}

// MigrateMinio does migration of a `bucketName` to a minio located in `minioHost`
func MigrateMinio(minioHost string, bucketName string, opts ...dockerMigrateMinioOption) error {
	var (
		err error
		dm  = &dockerMigrateMinio{}
	)

	for _, opt := range opts {
		opt(dm)
	}

	if dm.pool == nil {
		dm.pool, err = dockertest.NewPool("")
		if err != nil {
			return fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dm.versionTag == "" {
		dm.versionTag = "RELEASE.2022-08-28T20-08-11Z"
	}

	runOpts := &dockertest.RunOptions{
		Repository: "minio/mc",
		Tag:        dm.versionTag,
		Entrypoint: []string{
			"bin/sh",
			"-c",
			fmt.Sprintf(`
			/usr/bin/mc alias set myminio http://%s minio minio123;
			/usr/bin/mc rm -r --force %s;
			/usr/bin/mc mb myminio/%s;
			`, minioHost, bucketName, bucketName),
		},
	}

	if dm.network != nil {
		runOpts.NetworkID = dm.network.Network.ID
	}

	resource, err := dm.pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return err
	}

	if err := resource.Expire(120); err != nil {
		return err
	}

	waitCtx, cancel := context.WithTimeout(context.Background(), waitContainerTimeout)
	defer cancel()

	// Ensure the command completed successfully.
	status, err := dm.pool.Client.WaitContainerWithContext(resource.Container.ID, waitCtx)
	if err != nil {
		return err
	}

	if status != 0 {
		stream := new(bytes.Buffer)

		if err = dm.pool.Client.Logs(docker.LogsOptions{
			Context:      waitCtx,
			OutputStream: stream,
			ErrorStream:  stream,
			Stdout:       true,
			Stderr:       true,
			Container:    resource.Container.ID,
		}); err != nil {
			return err
		}

		return fmt.Errorf("got non-zero exit code %s", stream.String())
	}

	if err := dm.pool.Purge(resource); err != nil {
		return err
	}

	return nil
}
