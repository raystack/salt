package cli_test

import (
	"context"
	"testing"

	"github.com/raystack/salt/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResource struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Secret string `json:"-"`
}

type customResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	tags []string
}

func (r *customResource) ExportData(fields []string) map[string]any {
	data := cli.StructExportData(r, fields)
	for _, f := range fields {
		if f == "tags" {
			data["tags"] = r.tags
		}
	}
	return data
}

func TestStructExportData(t *testing.T) {
	r := testResource{ID: 1, Name: "alice", Status: "active", Secret: "s3cret"}

	t.Run("extracts requested fields by json tag", func(t *testing.T) {
		data := cli.StructExportData(r, []string{"id", "name"})
		assert.Equal(t, map[string]any{"id": 1, "name": "alice"}, data)
	})

	t.Run("handles all fields", func(t *testing.T) {
		data := cli.StructExportData(r, []string{"id", "name", "status"})
		assert.Equal(t, 3, len(data))
	})

	t.Run("skips unknown fields", func(t *testing.T) {
		data := cli.StructExportData(r, []string{"id", "nonexistent"})
		assert.Equal(t, map[string]any{"id": 1}, data)
	})

	t.Run("works with pointer", func(t *testing.T) {
		data := cli.StructExportData(&r, []string{"name"})
		assert.Equal(t, map[string]any{"name": "alice"}, data)
	})
}

func TestExporter_Write(t *testing.T) {
	resources := []testResource{
		{ID: 1, Name: "alice", Status: "active"},
		{ID: 2, Name: "bob", Status: "inactive"},
	}

	t.Run("compact JSON when not TTY", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()

		cmd := &cobra.Command{Use: "test"}
		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id", "name", "status"})
		cmd.SetArgs([]string{"--json", "id,name"})
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return exporter.Write(ios, resources)
		}

		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())

		assert.Contains(t, stdout.String(), `"id":1`)
		assert.Contains(t, stdout.String(), `"name":"alice"`)
		assert.NotContains(t, stdout.String(), `"status"`)
		// Compact: no leading spaces.
		assert.NotContains(t, stdout.String(), "  ")
	})

	t.Run("indented JSON on TTY", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.SetStdoutTTY(true)

		cmd := &cobra.Command{Use: "test"}
		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id", "name", "status"})
		cmd.SetArgs([]string{"--json", "id,name"})
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return exporter.Write(ios, resources)
		}

		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())

		assert.Contains(t, stdout.String(), "  ")
	})
}

func TestExporter_Exportable(t *testing.T) {
	resources := []*customResource{
		{ID: 1, Name: "proj", tags: []string{"go", "cli"}},
	}

	ios, _, stdout, _ := cli.Test()

	cmd := &cobra.Command{Use: "test"}
	var exporter cli.Exporter
	cli.AddJSONFlags(cmd, &exporter, []string{"id", "name", "tags"})
	cmd.SetArgs([]string{"--json", "name,tags"})
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return exporter.Write(ios, resources)
	}

	ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
	cmd.SetContext(ctx)
	require.NoError(t, cmd.Execute())

	assert.Contains(t, stdout.String(), `"tags":["go","cli"]`)
	assert.Contains(t, stdout.String(), `"name":"proj"`)
}

func TestAddJSONFlags(t *testing.T) {
	t.Run("exporter is nil when --json not used", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id", "name"})
		cmd.SetArgs([]string{})
		require.NoError(t, cmd.Execute())
		assert.Nil(t, exporter)
	})

	t.Run("rejects unknown fields", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id", "name"})
		cmd.SetArgs([]string{"--json", "id,bogus"})
		err := cmd.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unknown JSON field: "bogus"`)
		assert.Contains(t, err.Error(), "id")
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("returns fields from exporter", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id", "name", "status"})
		cmd.SetArgs([]string{"--json", "id,status"})
		require.NoError(t, cmd.Execute())
		require.NotNil(t, exporter)
		assert.Equal(t, []string{"id", "status"}, exporter.Fields())
	})
}
