package cmdx_test

import (
	"testing"

	"github.com/odpf/salt/cmdx"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Host string `default:"localhost"`
}

func TestInit(t *testing.T) {
	t.Run("should return config filename", func(t *testing.T) {
		c := cmdx.SetConfig("stencil")

		assert.Contains(t, c.File(), "stencil.yml")
	})

	t.Run("should return default config", func(t *testing.T) {
		cliconfig := &TestConfig{}

		c := cmdx.SetConfig("stencil")
		c.Defaults(cliconfig)

		assert.Equal(t, "localhost", cliconfig.Host)
	})
}
