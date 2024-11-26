package dockertestx

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type dockerMigrateSpiceDBOption func(dmm *dockerMigrateSpiceDB)

// MigrateSpiceDBWithDockertestNetwork is an option to assign docker network
func MigrateSpiceDBWithDockertestNetwork(network *dockertest.Network) dockerMigrateSpiceDBOption {
	return func(dm *dockerMigrateSpiceDB) {
		dm.network = network
	}
}

// MigrateSpiceDBWithVersionTag is an option to assign version tag
// of a `quay.io/authzed/spicedb` image
func MigrateSpiceDBWithVersionTag(versionTag string) dockerMigrateSpiceDBOption {
	return func(dm *dockerMigrateSpiceDB) {
		dm.versionTag = versionTag
	}
}

// MigrateSpiceDBWithDockerPool is an option to assign docker pool
func MigrateSpiceDBWithDockerPool(pool *dockertest.Pool) dockerMigrateSpiceDBOption {
	return func(dm *dockerMigrateSpiceDB) {
		dm.pool = pool
	}
}

type dockerMigrateSpiceDB struct {
	network    *dockertest.Network
	pool       *dockertest.Pool
	versionTag string
}

// MigrateSpiceDB migrates spicedb with postgres backend
func MigrateSpiceDB(postgresConnectionURL string, opts ...dockerMigrateMinioOption) error {
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
		dm.versionTag = "v1.0.0"
	}

	runOpts := &dockertest.RunOptions{
		Repository: "quay.io/authzed/spicedb",
		Tag:        dm.versionTag,
		Cmd:        []string{"spicedb", "migrate", "head", "--datastore-engine", "postgres", "--datastore-conn-uri", postgresConnectionURL},
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
