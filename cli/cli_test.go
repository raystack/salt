package cli_test

import (
	"bytes"
	"testing"

	"github.com/raystack/salt/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRoot() *cobra.Command {
	return &cobra.Command{Use: "testcli", Short: "test app"}
}

func TestInit(t *testing.T) {
	t.Run("adds completion command", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root)

		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "completion" {
				found = true
			}
		}
		assert.True(t, found, "completion command should be added")
	})

	t.Run("adds reference command", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root)

		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "reference" {
				found = true
			}
		}
		assert.True(t, found, "reference command should be added")
	})

	t.Run("adds version command when configured", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root, cli.Version("1.0.0", "raystack/test"))

		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "version" {
				found = true
			}
		}
		assert.True(t, found, "version command should be added")
	})

	t.Run("no version command without option", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root)

		for _, cmd := range root.Commands() {
			assert.NotEqual(t, "version", cmd.Name(), "version command should not be added without option")
		}
	})

	t.Run("version command prints version", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root, cli.Version("2.5.0", ""))

		var buf bytes.Buffer
		root.SetOut(&buf)
		root.SetArgs([]string{"version"})
		err := root.Execute()
		require.NoError(t, err)
	})

	t.Run("silences cobra error and usage output", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root)

		assert.True(t, root.SilenceErrors, "SilenceErrors should be true")
		assert.True(t, root.SilenceUsage, "SilenceUsage should be true")
	})

	t.Run("wraps flag errors for Execute", func(t *testing.T) {
		root := newTestRoot()
		root.Flags().Int("port", 8080, "server port")
		root.RunE = func(cmd *cobra.Command, args []string) error { return nil }
		cli.Init(root)

		// Unknown flag returns an error (wrapped internally as flagError).
		root.SetArgs([]string{"--unknown"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown flag")

		// Invalid flag value also returns an error.
		root.SetArgs([]string{"--port", "abc"})
		err = root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid argument")
	})
}

func TestOutput(t *testing.T) {
	t.Run("returns output from context after Init", func(t *testing.T) {
		root := newTestRoot()
		cli.Init(root)

		var out *bytes.Buffer
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, _ []string) {
				o := cli.Output(cmd)
				assert.NotNil(t, o)
				out = &bytes.Buffer{}
			},
		}
		root.AddCommand(child)
		root.SetArgs([]string{"child"})
		root.Execute()
		assert.NotNil(t, out)
	})

	t.Run("returns fallback when no context", func(t *testing.T) {
		cmd := &cobra.Command{Use: "bare"}
		out := cli.Output(cmd)
		assert.NotNil(t, out, "should return fallback output")
	})
}

func TestPrompter(t *testing.T) {
	t.Run("returns fallback when no context", func(t *testing.T) {
		cmd := &cobra.Command{Use: "bare"}
		p := cli.Prompter(cmd)
		assert.NotNil(t, p, "should return fallback prompter")
	})
}
