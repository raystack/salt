package dockertest

import "github.com/ory/dockertest/v3"

// type and function aliasing to avoid conflicting dockertest name
type Pool = dockertest.Pool
type Resource = dockertest.Resource

var NewPool = dockertest.NewPool
