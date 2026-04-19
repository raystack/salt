package cli_test

import (
	"context"
	"testing"

	"github.com/raystack/salt/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {
	ios := cli.System()
	assert.NotNil(t, ios.In)
	assert.NotNil(t, ios.Out)
	assert.NotNil(t, ios.ErrOut)
}

func TestTest(t *testing.T) {
	ios, stdin, stdout, stderr := cli.Test()

	assert.False(t, ios.IsStdinTTY(), "test stdin should not be TTY")
	assert.False(t, ios.IsStdoutTTY(), "test stdout should not be TTY")
	assert.False(t, ios.IsStderrTTY(), "test stderr should not be TTY")
	assert.False(t, ios.ColorEnabled(), "test should have color disabled")

	// Verify buffers are wired correctly.
	_, err := ios.Out.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, "hello", stdout.String())

	_, err = ios.ErrOut.Write([]byte("warn"))
	require.NoError(t, err)
	assert.Equal(t, "warn", stderr.String())

	assert.NotNil(t, stdin, "stdin buffer should be returned")
}

func TestIOStreams_CanPrompt(t *testing.T) {
	t.Run("false when not TTY", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		assert.False(t, ios.CanPrompt())
	})

	t.Run("true when both stdin and stdout are TTY", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetStdinTTY(true)
		ios.SetStdoutTTY(true)
		assert.True(t, ios.CanPrompt())
	})

	t.Run("false when NeverPrompt is set", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetStdinTTY(true)
		ios.SetStdoutTTY(true)
		ios.SetNeverPrompt(true)
		assert.False(t, ios.CanPrompt())
	})

	t.Run("false when only stdin is TTY", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetStdinTTY(true)
		assert.False(t, ios.CanPrompt())
	})
}

func TestIOStreams_TTYOverrides(t *testing.T) {
	ios, _, _, _ := cli.Test()

	ios.SetStdinTTY(true)
	assert.True(t, ios.IsStdinTTY())

	ios.SetStdoutTTY(true)
	assert.True(t, ios.IsStdoutTTY())

	ios.SetStderrTTY(true)
	assert.True(t, ios.IsStderrTTY())
}

func TestIOStreams_ColorEnabled(t *testing.T) {
	ios, _, _, _ := cli.Test()
	assert.False(t, ios.ColorEnabled())

	ios.SetColorEnabled(true)
	assert.True(t, ios.ColorEnabled())
}

func TestIOStreams_TerminalWidth(t *testing.T) {
	ios, _, _, _ := cli.Test()
	// Non-file writer returns default 80.
	assert.Equal(t, 80, ios.TerminalWidth())
}

func TestIOStreams_Output(t *testing.T) {
	ios, _, stdout, _ := cli.Test()

	out := ios.Output()
	assert.NotNil(t, out)

	// Same instance on repeated calls.
	assert.Same(t, out, ios.Output())

	// Writes go to the stdout buffer.
	out.Println("test output")
	assert.Contains(t, stdout.String(), "test output")
}

func TestIOStreams_OutputResetsOnTTYChange(t *testing.T) {
	ios, _, _, _ := cli.Test()
	out1 := ios.Output()

	ios.SetStdoutTTY(true)
	out2 := ios.Output()
	assert.NotSame(t, out1, out2, "Output should be recreated after TTY change")
}

func TestIOStreams_Prompter(t *testing.T) {
	ios, _, _, _ := cli.Test()
	p := ios.Prompter()
	assert.NotNil(t, p)
	// Same instance on repeated calls.
	assert.Same(t, p, ios.Prompter())
}

func TestIO(t *testing.T) {
	t.Run("extracts IOStreams from context", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)

		got := cli.IO(cmd)
		assert.Same(t, ios, got)
	})

	t.Run("returns fallback when no context", func(t *testing.T) {
		cmd := &cobra.Command{Use: "bare"}
		got := cli.IO(cmd)
		assert.NotNil(t, got)
	})
}
