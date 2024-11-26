# dockertestx

This package is an abstraction of several dockerized data storages using `ory/dockertest` to bootstrap a specific dockerized instance.

Example postgres

```go
// create postgres instance
pgDocker, err := dockertest.CreatePostgres(
    dockertest.PostgresWithDetail(
        pgUser, pgPass, pgDBName,
    ),
)

// get connection string
connString := pgDocker.GetExternalConnString()

// purge docker
if err := pgDocker.GetPool().Purge(pgDocker.GetResouce()); err != nil {
    return fmt.Errorf("could not purge resource: %w", err)
}
```

Example spice db

- bootsrap spice db with postgres and wire them internally via network bridge

```go
// create custom pool
pool, err := dockertest.NewPool("")
if err != nil {
    return nil, err
}

// create a bridge network for testing
network, err = pool.Client.CreateNetwork(docker.CreateNetworkOptions{
    Name: fmt.Sprintf("bridge-%s", uuid.New().String()),
})
if err != nil {
    return nil, err
}


// create postgres instance
pgDocker, err := dockertest.CreatePostgres(
    dockertest.PostgresWithDockerPool(pool),
    dockertest.PostgresWithDockertestNetwork(network),
    dockertest.PostgresWithDetail(
        pgUser, pgPass, pgDBName,
    ),
)

// get connection string
connString := pgDocker.GetInternalConnString()

// create spice db instance
spiceDocker, err := dockertest.CreateSpiceDB(connString,
    dockertest.SpiceDBWithDockerPool(pool),
    dockertest.SpiceDBWithDockertestNetwork(network),
)

if err := dockertest.MigrateSpiceDB(connString,
    dockertest.MigrateSpiceDBWithDockerPool(pool),
    dockertest.MigrateSpiceDBWithDockertestNetwork(network),
); err != nil {
    return err
}

// purge docker resources
if err := pool.Purge(spiceDocker.GetResouce()); err != nil {
    return fmt.Errorf("could not purge resource: %w", err)
}
if err := pool.Purge(pgDocker.GetResouce()); err != nil {
    return fmt.Errorf("could not purge resource: %w", err)
}
```
