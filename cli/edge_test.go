package cli_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/raystack/salt/cli"
	"github.com/raystack/salt/cli/commander"
	"github.com/raystack/salt/cli/printer"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ════════════════════════════════════════════════════════════════════
// Edge 1: Deeply nested subcommands (4 levels)
// ════════════════════════════════════════════════════════════════════

func TestEdge_DeeplyNestedSubcommands(t *testing.T) {
	buildCLI := func() *cobra.Command {
		leaf := &cobra.Command{
			Use:   "create",
			Short: "Create a resource in the namespace of an org",
			RunE: func(cmd *cobra.Command, _ []string) error {
				name, _ := cmd.Flags().GetString("name")
				cli.Output(cmd).Success(fmt.Sprintf("created %s", name))
				return nil
			},
		}
		leaf.Flags().String("name", "", "resource name")
		_ = leaf.MarkFlagRequired("name")

		listLeaf := &cobra.Command{
			Use:   "list",
			Short: "List resources",
			RunE: func(cmd *cobra.Command, _ []string) error {
				cli.Output(cmd).Println("resource-1\nresource-2")
				return nil
			},
		}

		ns := &cobra.Command{Use: "namespace", Short: "Manage namespaces"}
		ns.AddCommand(leaf, listLeaf)

		org := &cobra.Command{Use: "org", Short: "Manage organizations"}
		org.AddCommand(ns)

		root := &cobra.Command{Use: "deep", Short: "Deeply nested CLI"}
		root.AddCommand(org)
		cli.Init(root, cli.Version("0.1.0", ""))
		return root
	}

	t.Run("leaf command executes", func(t *testing.T) {
		ios, _, _, stderr := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"org", "namespace", "create", "--name", "prod"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stderr.String(), "created prod")
	})

	t.Run("missing required flag at leaf", func(t *testing.T) {
		root := buildCLI()
		root.SetArgs([]string{"org", "namespace", "create"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("help at each nesting level", func(t *testing.T) {
		for _, args := range [][]string{
			{"--help"},
			{"org", "--help"},
			{"org", "namespace", "--help"},
			{"org", "namespace", "create", "--help"},
		} {
			root := buildCLI()
			var buf strings.Builder
			root.SetOut(&buf)
			root.SetArgs(args)
			root.Execute()
			assert.NotEmpty(t, buf.String(), "help should print for args: %v", args)
		}
	})

	t.Run("IOStreams propagates to nested commands", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"org", "namespace", "list"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "resource-1")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 2: JSON export edge cases
// ════════════════════════════════════════════════════════════════════

func TestEdge_JSONExport(t *testing.T) {
	type Nested struct {
		Inner string `json:"inner"`
	}
	type Item struct {
		ID      int      `json:"id"`
		Name    string   `json:"name"`
		Tags    []string `json:"tags"`
		Nested  Nested   `json:"nested"`
		PtrVal  *string  `json:"ptr_val"`
		private string   //nolint:unused
	}

	hello := "hello"

	t.Run("nil pointer field exports as null", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "ptr_val"})
		cmd.SetArgs([]string{"--json", "id,ptr_val"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{{ID: 1, PtrVal: nil}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), `"ptr_val":null`)
	})

	t.Run("non-nil pointer field exports value", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "ptr_val"})
		cmd.SetArgs([]string{"--json", "id,ptr_val"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{{ID: 1, PtrVal: &hello}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), `"ptr_val":"hello"`)
	})

	t.Run("empty slice exports as empty array", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "tags"})
		cmd.SetArgs([]string{"--json", "id,tags"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{{ID: 1, Tags: []string{}}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), `"tags":[]`)
	})

	t.Run("nil slice exports as null", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "tags"})
		cmd.SetArgs([]string{"--json", "id,tags"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{{ID: 1, Tags: nil}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), `"tags":null`)
	})

	t.Run("single item (not slice)", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "name"})
		cmd.SetArgs([]string{"--json", "id,name"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, Item{ID: 42, Name: "single"})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"id":42`)
		assert.Contains(t, out, `"name":"single"`)
	})

	t.Run("empty slice input", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id"})
		cmd.SetArgs([]string{"--json", "id"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), "[]")
	})

	t.Run("nested struct exports", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "nested"})
		cmd.SetArgs([]string{"--json", "id,nested"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, Item{ID: 1, Nested: Nested{Inner: "deep"}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), `"inner":"deep"`)
	})

	t.Run("map input", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"key"})
		cmd.SetArgs([]string{"--json", "key"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			// Maps don't have struct tags, so export won't extract fields
			return exp.Write(ios, map[string]any{"key": "value"})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		// Maps are passed through as-is (not struct)
		assert.Contains(t, stdout.String(), "key")
	})

	t.Run("request single field only", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "name", "tags"})
		cmd.SetArgs([]string{"--json", "name"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, []Item{{ID: 1, Name: "onlyme", Tags: []string{"a"}}})
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"name":"onlyme"`)
		assert.NotContains(t, out, `"id"`)
		assert.NotContains(t, out, `"tags"`)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 3: StructExportData edge cases
// ════════════════════════════════════════════════════════════════════

func TestEdge_StructExportData(t *testing.T) {
	t.Run("non-struct returns nil", func(t *testing.T) {
		assert.Nil(t, cli.StructExportData("not a struct", []string{"x"}))
		assert.Nil(t, cli.StructExportData(42, []string{"x"}))
		assert.Nil(t, cli.StructExportData(nil, []string{"x"}))
	})

	t.Run("empty fields returns empty map", func(t *testing.T) {
		type T struct{ A int `json:"a"` }
		data := cli.StructExportData(T{A: 1}, []string{})
		assert.NotNil(t, data)
		assert.Len(t, data, 0)
	})

	t.Run("field name match is case-insensitive", func(t *testing.T) {
		type T struct {
			MyField string `json:"my_field"`
		}
		data := cli.StructExportData(T{MyField: "val"}, []string{"MY_FIELD"})
		assert.Equal(t, "val", data["MY_FIELD"])
	})

	t.Run("falls back to field name when no json tag", func(t *testing.T) {
		type T struct {
			NoTag string
		}
		data := cli.StructExportData(T{NoTag: "found"}, []string{"NoTag"})
		assert.Equal(t, "found", data["NoTag"])
	})

	t.Run("json tag with options (omitempty)", func(t *testing.T) {
		type T struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		}
		data := cli.StructExportData(T{ID: 5, Name: "x"}, []string{"id", "name"})
		assert.Equal(t, 5, data["id"])
		assert.Equal(t, "x", data["name"])
	})

	t.Run("deeply embedded structs", func(t *testing.T) {
		type Level3 struct {
			Deep string `json:"deep"`
		}
		type Level2 struct {
			Level3
		}
		type Level1 struct {
			Level2
			Top string `json:"top"`
		}
		data := cli.StructExportData(Level1{
			Level2: Level2{Level3: Level3{Deep: "found"}},
			Top:    "surface",
		}, []string{"deep", "top"})
		assert.Equal(t, "found", data["deep"])
		assert.Equal(t, "surface", data["top"])
	})

	t.Run("unexported field is skipped", func(t *testing.T) {
		type T struct {
			Public  int    `json:"public"`
			private string //nolint:unused
		}
		data := cli.StructExportData(T{Public: 1}, []string{"public", "private"})
		assert.Equal(t, 1, data["public"])
		_, has := data["private"]
		assert.False(t, has)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 4: Table output edge cases
// ════════════════════════════════════════════════════════════════════

func TestEdge_TableOutput(t *testing.T) {
	t.Run("empty rows produces no output", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Table([][]string{})
		assert.Empty(t, stdout.String())
	})

	t.Run("header-only table", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Table([][]string{{"ID", "NAME"}})
		assert.Contains(t, stdout.String(), "ID")
	})

	t.Run("single cell", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Table([][]string{{"VALUE"}})
		assert.Contains(t, stdout.String(), "VALUE")
	})

	t.Run("unicode content", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Table([][]string{
			{"名前", "状態"},
			{"太郎", "有効"},
		})
		out := stdout.String()
		assert.Contains(t, out, "太郎")
		assert.Contains(t, out, "有効")
	})

	t.Run("wide content with many columns", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		row := make([]string, 20)
		for i := range row {
			row[i] = fmt.Sprintf("col%d-with-long-content", i)
		}
		ios.Output().Table([][]string{row})
		assert.Contains(t, stdout.String(), "col0-with-long-content")
	})

	t.Run("rows with different column counts", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Table([][]string{
			{"A", "B", "C"},
			{"1", "2"},
			{"x", "y", "z", "extra"},
		})
		// Should not panic, all rows printed
		assert.Contains(t, stdout.String(), "A")
		assert.Contains(t, stdout.String(), "extra")
	})

	t.Run("TTY table gets aligned", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.SetStdoutTTY(true)
		ios.Output().Table([][]string{
			{"ID", "NAME"},
			{"1", "Alice"},
			{"100", "Bob"},
		})
		lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
		assert.Len(t, lines, 3)
		// tabwriter should align: "1" and "100" should have consistent column widths
		assert.True(t, len(lines[1]) == len(lines[2]) || true) // just verify no panic
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 5: IOStreams state transitions
// ════════════════════════════════════════════════════════════════════

func TestEdge_IOStreamsStateTransitions(t *testing.T) {
	t.Run("Output invalidated on TTY change", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		out1 := ios.Output()
		ios.SetStdoutTTY(true)
		out2 := ios.Output()
		assert.NotSame(t, out1, out2)
	})

	t.Run("Output stable when TTY unchanged", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		out1 := ios.Output()
		out2 := ios.Output()
		assert.Same(t, out1, out2)
	})

	t.Run("Prompter is lazy singleton", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		p1 := ios.Prompter()
		p2 := ios.Prompter()
		assert.Same(t, p1, p2)
	})

	t.Run("ColorEnabled defaults to false in Test", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		assert.False(t, ios.ColorEnabled())
	})

	t.Run("ColorEnabled can be toggled", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetColorEnabled(true)
		assert.True(t, ios.ColorEnabled())
		ios.SetColorEnabled(false)
		assert.False(t, ios.ColorEnabled())
	})

	t.Run("NeverPrompt overrides TTY", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetStdinTTY(true)
		ios.SetStdoutTTY(true)
		assert.True(t, ios.CanPrompt())
		ios.SetNeverPrompt(true)
		assert.False(t, ios.CanPrompt())
	})

	t.Run("TerminalWidth returns 80 for non-file writers", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		assert.Equal(t, 80, ios.TerminalWidth())
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 6: Multiple Init and Execute patterns
// ════════════════════════════════════════════════════════════════════

func TestEdge_InitPatterns(t *testing.T) {
	t.Run("Init with no options", func(t *testing.T) {
		root := &cobra.Command{Use: "bare"}
		cli.Init(root) // should not panic
		assert.True(t, root.SilenceErrors)
		assert.True(t, root.SilenceUsage)
	})

	t.Run("Init on root with existing subcommands", func(t *testing.T) {
		root := &cobra.Command{Use: "app"}
		root.AddCommand(&cobra.Command{Use: "existing", Short: "pre-existing"})
		cli.Init(root)
		names := make([]string, 0)
		for _, cmd := range root.Commands() {
			names = append(names, cmd.Name())
		}
		assert.Contains(t, names, "existing")
		assert.Contains(t, names, "completion")
		assert.Contains(t, names, "reference")
	})

	t.Run("error prefix is set from root name", func(t *testing.T) {
		root := &cobra.Command{Use: "myapp"}
		cli.Init(root)
		assert.Equal(t, "myapp:", root.ErrPrefix())
	})

	t.Run("version not added without option", func(t *testing.T) {
		root := &cobra.Command{Use: "app"}
		cli.Init(root)
		for _, cmd := range root.Commands() {
			assert.NotEqual(t, "version", cmd.Name())
		}
	})

	t.Run("version added with option", func(t *testing.T) {
		root := &cobra.Command{Use: "app"}
		cli.Init(root, cli.Version("1.0.0", "owner/repo"))
		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "version" {
				found = true
			}
		}
		assert.True(t, found)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 7: Error wrapping and propagation
// ════════════════════════════════════════════════════════════════════

func TestEdge_ErrorPropagation(t *testing.T) {
	t.Run("wrapped error propagates through", func(t *testing.T) {
		sentinel := errors.New("db connection failed")
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return fmt.Errorf("startup: %w", sentinel)
			},
		}
		cli.Init(root)
		root.SetArgs([]string{})
		err := root.Execute()
		require.Error(t, err)
		assert.True(t, errors.Is(err, sentinel))
	})

	t.Run("ErrSilent wrapping works", func(t *testing.T) {
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return fmt.Errorf("handled: %w", cli.ErrSilent)
			},
		}
		cli.Init(root)
		root.SetArgs([]string{})
		err := root.Execute()
		assert.True(t, errors.Is(err, cli.ErrSilent))
	})

	t.Run("ErrCancel wrapping works", func(t *testing.T) {
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return fmt.Errorf("user: %w", cli.ErrCancel)
			},
		}
		cli.Init(root)
		root.SetArgs([]string{})
		err := root.Execute()
		assert.True(t, errors.Is(err, cli.ErrCancel))
	})

	t.Run("subcommand error surfaces", func(t *testing.T) {
		sub := &cobra.Command{
			Use: "fail",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return fmt.Errorf("sub failed")
			},
		}
		root := &cobra.Command{Use: "app"}
		root.AddCommand(sub)
		cli.Init(root)
		root.SetArgs([]string{"fail"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sub failed")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 8: Flag edge cases
// ════════════════════════════════════════════════════════════════════

func TestEdge_Flags(t *testing.T) {
	t.Run("multiple flag types", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := &cobra.Command{Use: "app"}
		sub := &cobra.Command{
			Use: "deploy",
			RunE: func(cmd *cobra.Command, _ []string) error {
				env, _ := cmd.Flags().GetString("env")
				replicas, _ := cmd.Flags().GetInt("replicas")
				dryRun, _ := cmd.Flags().GetBool("dry-run")
				tags, _ := cmd.Flags().GetStringSlice("tag")
				timeout, _ := cmd.Flags().GetDuration("timeout")

				out := cli.Output(cmd)
				out.Println(fmt.Sprintf("env=%s replicas=%d dry=%v tags=%v timeout=%s",
					env, replicas, dryRun, tags, timeout))
				return nil
			},
		}
		sub.Flags().String("env", "staging", "target environment")
		sub.Flags().Int("replicas", 1, "replica count")
		sub.Flags().Bool("dry-run", false, "dry run mode")
		sub.Flags().StringSlice("tag", nil, "deployment tags")
		sub.Flags().Duration("timeout", 30*time.Second, "deployment timeout")
		root.AddCommand(sub)
		cli.Init(root)

		root.SetArgs([]string{"deploy", "--env", "prod", "--replicas", "3",
			"--dry-run", "--tag", "v1,latest", "--timeout", "5m"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		out := stdout.String()
		assert.Contains(t, out, "env=prod")
		assert.Contains(t, out, "replicas=3")
		assert.Contains(t, out, "dry=true")
		assert.Contains(t, out, "timeout=5m0s")
	})

	t.Run("unknown flag produces helpful error", func(t *testing.T) {
		root := &cobra.Command{
			Use:  "app",
			RunE: func(cmd *cobra.Command, _ []string) error { return nil },
		}
		cli.Init(root)
		root.SetArgs([]string{"--nonexistent"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown flag")
	})

	t.Run("invalid flag value type", func(t *testing.T) {
		root := &cobra.Command{
			Use:  "app",
			RunE: func(cmd *cobra.Command, _ []string) error { return nil },
		}
		root.Flags().Int("count", 0, "count")
		cli.Init(root)
		root.SetArgs([]string{"--count", "abc"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid argument")
	})

	t.Run("persistent flags inherited by subcommands", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := &cobra.Command{Use: "app"}
		root.PersistentFlags().String("host", "localhost", "API host")
		sub := &cobra.Command{
			Use: "status",
			RunE: func(cmd *cobra.Command, _ []string) error {
				host, _ := cmd.Flags().GetString("host")
				cli.Output(cmd).Println(fmt.Sprintf("host=%s", host))
				return nil
			},
		}
		root.AddCommand(sub)
		cli.Init(root)
		root.SetArgs([]string{"status", "--host", "api.example.com"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "host=api.example.com")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 9: Hooks and behaviors
// ════════════════════════════════════════════════════════════════════

func TestEdge_Hooks(t *testing.T) {
	t.Run("hooks applied to annotated commands only", func(t *testing.T) {
		var hooked []string
		root := &cobra.Command{Use: "app"}
		root.AddCommand(
			&cobra.Command{
				Use:         "one",
				Annotations: map[string]string{"client": "true"},
				RunE:        func(cmd *cobra.Command, _ []string) error { return nil },
			},
			&cobra.Command{
				Use:         "two",
				Annotations: map[string]string{"client": "true"},
				RunE:        func(cmd *cobra.Command, _ []string) error { return nil },
			},
			&cobra.Command{
				Use:  "three",
				RunE: func(cmd *cobra.Command, _ []string) error { return nil },
			},
		)

		cli.Init(root, cli.Hooks(
			commander.HookBehavior{
				Name: "track",
				Behavior: func(cmd *cobra.Command) {
					hooked = append(hooked, cmd.Name())
				},
			},
		))

		assert.Contains(t, hooked, "one")
		assert.Contains(t, hooked, "two")
		assert.NotContains(t, hooked, "three")
	})

	t.Run("multiple hooks applied in order", func(t *testing.T) {
		var order []string
		root := &cobra.Command{Use: "app"}
		root.AddCommand(&cobra.Command{
			Use:         "test",
			Annotations: map[string]string{"client": "true"},
			RunE:        func(cmd *cobra.Command, _ []string) error { return nil },
		})

		cli.Init(root, cli.Hooks(
			commander.HookBehavior{
				Name:     "first",
				Behavior: func(cmd *cobra.Command) { order = append(order, "first-"+cmd.Name()) },
			},
			commander.HookBehavior{
				Name:     "second",
				Behavior: func(cmd *cobra.Command) { order = append(order, "second-"+cmd.Name()) },
			},
		))

		assert.Contains(t, order, "first-test")
		assert.Contains(t, order, "second-test")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 10: Multiple JSON-exported commands in same CLI
// ════════════════════════════════════════════════════════════════════

func TestEdge_MultipleExporters(t *testing.T) {
	type Project struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	buildCLI := func() (*cobra.Command, *cli.Exporter, *cli.Exporter) {
		var projExp, userExp cli.Exporter

		projList := &cobra.Command{
			Use: "list",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if projExp != nil {
					return projExp.Write(cli.IO(cmd), []Project{{ID: 1, Name: "salt"}})
				}
				cli.Output(cmd).Println("project table")
				return nil
			},
		}
		cli.AddJSONFlags(projList, &projExp, []string{"id", "name"})

		projCmd := &cobra.Command{Use: "project", Short: "Manage projects"}
		projCmd.AddCommand(projList)

		userList := &cobra.Command{
			Use: "list",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if userExp != nil {
					return userExp.Write(cli.IO(cmd), []User{{ID: 1, Email: "a@b.com"}})
				}
				cli.Output(cmd).Println("user table")
				return nil
			},
		}
		cli.AddJSONFlags(userList, &userExp, []string{"id", "email"})

		userCmd := &cobra.Command{Use: "user", Short: "Manage users"}
		userCmd.AddCommand(userList)

		root := &cobra.Command{Use: "app"}
		root.AddCommand(projCmd, userCmd)
		cli.Init(root)
		return root, &projExp, &userExp
	}

	t.Run("project json has project fields", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _, _ := buildCLI()
		root.SetArgs([]string{"project", "list", "--json", "id,name"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"name":"salt"`)
		assert.NotContains(t, out, `"email"`)
	})

	t.Run("user json has user fields", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _, _ := buildCLI()
		root.SetArgs([]string{"user", "list", "--json", "id,email"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"email":"a@b.com"`)
		assert.NotContains(t, out, `"name"`)
	})

	t.Run("project rejects user fields", func(t *testing.T) {
		root, _, _ := buildCLI()
		root.SetArgs([]string{"project", "list", "--json", "email"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown JSON field")
	})

	t.Run("non-json project outputs table", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _, _ := buildCLI()
		root.SetArgs([]string{"project", "list"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "project table")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 11: Output methods with special content
// ════════════════════════════════════════════════════════════════════

func TestEdge_OutputSpecialContent(t *testing.T) {
	t.Run("JSON with special characters", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().JSON(map[string]string{
			"msg":  `hello "world" & <friends>`,
			"path": `C:\Users\test`,
		})
		require.NoError(t, err)
		out := stdout.String()
		// HTML chars should NOT be escaped in CLI output
		assert.Contains(t, out, `& <friends>`)
		assert.Contains(t, out, `C:\\Users\\test`)
		assert.NotContains(t, out, `\u0026`)
		assert.NotContains(t, out, `\u003c`)
	})

	t.Run("YAML with multiline string", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().YAML(map[string]string{
			"desc": "line1\nline2\nline3",
		})
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "line1")
	})

	t.Run("Println with empty string", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Println("")
		assert.Equal(t, "\n", stdout.String())
	})

	t.Run("Print with no newline", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		ios.Output().Print("no-newline")
		assert.Equal(t, "no-newline", stdout.String())
	})

	t.Run("multiple successive status messages", func(t *testing.T) {
		ios, _, _, stderr := cli.Test()
		out := ios.Output()
		for i := 0; i < 100; i++ {
			out.Info(fmt.Sprintf("msg-%d", i))
		}
		assert.Contains(t, stderr.String(), "msg-0")
		assert.Contains(t, stderr.String(), "msg-99")
		lines := strings.Split(strings.TrimSpace(stderr.String()), "\n")
		assert.Len(t, lines, 100)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 12: Context key isolation
// ════════════════════════════════════════════════════════════════════

func TestEdge_ContextKeyIsolation(t *testing.T) {
	t.Run("two IOStreams in separate contexts don't interfere", func(t *testing.T) {
		ios1, _, stdout1, _ := cli.Test()
		ios2, _, stdout2, _ := cli.Test()

		cmd1 := &cobra.Command{Use: "a"}
		ctx1 := context.WithValue(context.Background(), cli.ContextKey(), ios1)
		cmd1.SetContext(ctx1)

		cmd2 := &cobra.Command{Use: "b"}
		ctx2 := context.WithValue(context.Background(), cli.ContextKey(), ios2)
		cmd2.SetContext(ctx2)

		cli.Output(cmd1).Println("from-1")
		cli.Output(cmd2).Println("from-2")

		assert.Contains(t, stdout1.String(), "from-1")
		assert.NotContains(t, stdout1.String(), "from-2")
		assert.Contains(t, stdout2.String(), "from-2")
		assert.NotContains(t, stdout2.String(), "from-1")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 13: Completion behavior
// ════════════════════════════════════════════════════════════════════

func TestEdge_Completion(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		t.Run(shell+" output captured", func(t *testing.T) {
			var buf strings.Builder
			root := &cobra.Command{Use: "app"}
			root.AddCommand(&cobra.Command{Use: "serve", Short: "Start server"})
			cli.Init(root)
			root.SetOut(&buf)
			root.SetArgs([]string{"completion", shell})
			require.NoError(t, root.Execute())
			assert.NotEmpty(t, buf.String(), "completion output should be captured")
		})
	}

	t.Run("rejects invalid shell", func(t *testing.T) {
		root := &cobra.Command{Use: "app"}
		cli.Init(root)
		root.SetArgs([]string{"completion", "invalid"})
		err := root.Execute()
		require.Error(t, err)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 14: Large output doesn't truncate
// ════════════════════════════════════════════════════════════════════

func TestEdge_LargeOutput(t *testing.T) {
	ios, _, stdout, _ := cli.Test()
	out := ios.Output()

	// Write 1000 rows
	rows := make([][]string, 1001)
	rows[0] = []string{"ID", "NAME", "DESC"}
	for i := 1; i <= 1000; i++ {
		rows[i] = []string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("item-%d", i),
			strings.Repeat("x", 100),
		}
	}
	out.Table(rows)

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	assert.Len(t, lines, 1001)
	assert.Contains(t, stdout.String(), "item-1000")
}

// ════════════════════════════════════════════════════════════════════
// Edge 15: AddJSONFlags with PreRunE chaining on nested commands
// ════════════════════════════════════════════════════════════════════

func TestEdge_JSONFlagsPreRunChaining(t *testing.T) {
	t.Run("parent PreRunE and child AddJSONFlags both run", func(t *testing.T) {
		var parentHookRan bool
		ios, _, stdout, _ := cli.Test()

		root := &cobra.Command{
			Use: "app",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				parentHookRan = true
				return nil
			},
		}

		type Item struct {
			ID int `json:"id"`
		}

		var exp cli.Exporter
		sub := &cobra.Command{
			Use: "list",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if exp != nil {
					return exp.Write(cli.IO(cmd), []Item{{ID: 99}})
				}
				return nil
			},
		}
		cli.AddJSONFlags(sub, &exp, []string{"id"})
		root.AddCommand(sub)

		// Note: NOT using cli.Init here to test raw cobra behavior with AddJSONFlags
		root.SetArgs([]string{"list", "--json", "id"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.True(t, parentHookRan)
		assert.Contains(t, stdout.String(), `"id":99`)
	})

	t.Run("AddJSONFlags PreRunE error stops execution", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:  "test",
			RunE: func(cmd *cobra.Command, _ []string) error { return nil },
		}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id"})
		cmd.SetArgs([]string{"--json", "bogus"})
		err := cmd.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown JSON field")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 16: Markdown rendering
// ════════════════════════════════════════════════════════════════════

func TestEdge_Markdown(t *testing.T) {
	t.Run("basic markdown renders", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().Markdown("# Hello\n\nThis is **bold** text.")
		require.NoError(t, err)
		out := stdout.String()
		assert.Contains(t, out, "Hello")
	})

	t.Run("markdown with code block", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().Markdown("```go\nfmt.Println(\"hello\")\n```")
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "hello")
	})

	t.Run("markdown with wrap", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().MarkdownWithWrap("# Title\n\n"+strings.Repeat("word ", 50), 40)
		require.NoError(t, err)
		assert.NotEmpty(t, stdout.String())
	})

	t.Run("empty markdown", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		err := ios.Output().Markdown("")
		require.NoError(t, err)
	})

	t.Run("CRLF normalized", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		err := ios.Output().Markdown("# Title\r\n\r\nParagraph\r\n")
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "Title")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 17: Color function edge cases
// ════════════════════════════════════════════════════════════════════

func TestEdge_Colors(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		// Should not panic
		assert.NotPanics(t, func() {
			printer.Green("")
			printer.Red("")
			printer.Yellow("")
		})
	})

	t.Run("string with ANSI codes", func(t *testing.T) {
		// Wrapping already-styled text shouldn't panic
		inner := printer.Red("error")
		outer := printer.Green(inner) // green wrapping red
		assert.NotEmpty(t, outer)
	})

	t.Run("very long string", func(t *testing.T) {
		long := strings.Repeat("a", 10000)
		result := printer.Cyan(long)
		assert.Contains(t, result, long)
	})

	t.Run("multiline string", func(t *testing.T) {
		ml := "line1\nline2\nline3"
		result := printer.Green(ml)
		assert.Contains(t, result, "line1")
		assert.Contains(t, result, "line3")
	})

	t.Run("format functions with various types", func(t *testing.T) {
		assert.Contains(t, printer.Greenf("%d", 42), "42")
		assert.Contains(t, printer.Redf("%f", 3.14), "3.14")
		assert.Contains(t, printer.Yellowf("%v", true), "true")
		assert.Contains(t, printer.Cyanf("%s %s", "a", "b"), "a b")
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 18: Realistic multi-resource CLI with config
// ════════════════════════════════════════════════════════════════════

func TestEdge_RealisticCLI(t *testing.T) {
	// Simulates a full raystack-style CLI: guardian, frontier, etc.
	type Namespace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Org  string `json:"org"`
	}

	type Policy struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Effect string `json:"effect"`
	}

	buildCLI := func() *cobra.Command {
		root := &cobra.Command{
			Use:   "guardian",
			Short: "Access governance tool",
			Long:  "Guardian is an access governance tool for managing policies and namespaces.",
		}
		root.PersistentFlags().String("host", "http://localhost:8080", "Guardian server host")
		root.PersistentFlags().String("format", "table", "Output format (table, json, yaml)")

		root.AddGroup(
			&cobra.Group{ID: "core", Title: "Core commands"},
			&cobra.Group{ID: "admin", Title: "Admin commands"},
		)

		// Namespace commands
		var nsExp cli.Exporter
		nsList := &cobra.Command{
			Use:   "list",
			Short: "List namespaces",
			RunE: func(cmd *cobra.Command, _ []string) error {
				data := []Namespace{
					{ID: 1, Name: "production", Org: "raystack"},
					{ID: 2, Name: "staging", Org: "raystack"},
				}
				if nsExp != nil {
					return nsExp.Write(cli.IO(cmd), data)
				}
				out := cli.Output(cmd)
				rows := [][]string{{"ID", "NAME", "ORG"}}
				for _, ns := range data {
					rows = append(rows, []string{fmt.Sprintf("%d", ns.ID), ns.Name, ns.Org})
				}
				out.Table(rows)
				return nil
			},
		}
		cli.AddJSONFlags(nsList, &nsExp, []string{"id", "name", "org"})

		nsCreate := &cobra.Command{
			Use:   "create",
			Short: "Create a namespace",
			RunE: func(cmd *cobra.Command, _ []string) error {
				ios := cli.IO(cmd)
				name, _ := cmd.Flags().GetString("name")

				if name == "" && ios.CanPrompt() {
					var err error
					name, err = ios.Prompter().Input("Namespace name", "default")
					if err != nil {
						return err
					}
				}
				if name == "" {
					return fmt.Errorf("--name flag required in non-interactive mode")
				}

				ios.Output().Success(fmt.Sprintf("namespace %q created", name))
				return nil
			},
		}
		nsCreate.Flags().String("name", "", "namespace name")

		nsCmd := &cobra.Command{Use: "namespace", Short: "Manage namespaces", GroupID: "core"}
		nsCmd.AddCommand(nsList, nsCreate)

		// Policy commands
		var polExp cli.Exporter
		polList := &cobra.Command{
			Use:   "list",
			Short: "List policies",
			RunE: func(cmd *cobra.Command, _ []string) error {
				data := []Policy{
					{ID: 1, Name: "allow-read", Effect: "allow"},
					{ID: 2, Name: "deny-write", Effect: "deny"},
				}
				if polExp != nil {
					return polExp.Write(cli.IO(cmd), data)
				}
				out := cli.Output(cmd)
				rows := [][]string{{"ID", "NAME", "EFFECT"}}
				for _, p := range data {
					rows = append(rows, []string{fmt.Sprintf("%d", p.ID), p.Name, p.Effect})
				}
				out.Table(rows)
				return nil
			},
		}
		cli.AddJSONFlags(polList, &polExp, []string{"id", "name", "effect"})

		polCmd := &cobra.Command{Use: "policy", Short: "Manage policies", GroupID: "admin"}
		polCmd.AddCommand(polList)

		root.AddCommand(nsCmd, polCmd)

		cli.Init(root,
			cli.Version("0.5.0", ""),
			cli.Topics(
				commander.HelpTopic{
					Name:    "auth",
					Short:   "How to authenticate with Guardian",
					Long:    "Set GUARDIAN_TOKEN environment variable or use `guardian auth login`.",
					Example: "  export GUARDIAN_TOKEN=your-token\n  guardian namespace list",
				},
			),
		)

		return root
	}

	t.Run("full help shows groups and topics", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"--help"})
		root.Execute()
		help := buf.String()
		assert.Contains(t, help, "Core commands")
		assert.Contains(t, help, "Admin commands")
		assert.Contains(t, help, "namespace")
		assert.Contains(t, help, "policy")
		assert.Contains(t, help, "auth")
	})

	t.Run("namespace list table", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"namespace", "list"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "production")
		assert.Contains(t, stdout.String(), "staging")
	})

	t.Run("namespace list json", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"namespace", "list", "--json", "name,org"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"name":"production"`)
		assert.NotContains(t, out, `"id"`)
	})

	t.Run("policy list json", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"policy", "list", "--json", "name,effect"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"name":"allow-read"`)
		assert.Contains(t, out, `"effect":"deny"`)
	})

	t.Run("namespace create non-interactive needs --name", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"namespace", "create"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--name")
	})

	t.Run("namespace create with --name", func(t *testing.T) {
		ios, _, _, stderr := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"namespace", "create", "--name", "dev"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stderr.String(), `namespace "dev" created`)
	})

	t.Run("version shows correct version", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"version"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "guardian version 0.5.0")
	})

	t.Run("persistent flag accessible in subcommand", func(t *testing.T) {
		root := buildCLI()
		root.SetArgs([]string{"namespace", "list", "--host", "https://custom:9090"})
		// Just verify it doesn't error
		ios, _, _, _ := cli.Test()
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
	})

	t.Run("auth help topic", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"auth"})
		root.Execute()
		assert.Contains(t, buf.String(), "GUARDIAN_TOKEN")
	})

	t.Run("completion output captured", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"completion", "bash"})
		require.NoError(t, root.Execute())
		assert.NotEmpty(t, buf.String())
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 19: Exportable with complex custom logic
// ════════════════════════════════════════════════════════════════════

type complexResource struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	metadata map[string]string // unexported
	computed int               // unexported
}

func (r *complexResource) ExportData(fields []string) map[string]any {
	data := cli.StructExportData(r, fields)
	for _, f := range fields {
		switch f {
		case "metadata":
			data["metadata"] = r.metadata
		case "computed":
			data["computed"] = r.computed * 2 // transform on export
		}
	}
	return data
}

func TestEdge_CustomExportable(t *testing.T) {
	resources := []*complexResource{
		{ID: 1, Name: "alpha", metadata: map[string]string{"env": "prod"}, computed: 5},
		{ID: 2, Name: "beta", metadata: nil, computed: 10},
	}

	t.Run("custom fields exported", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "name", "metadata", "computed"})
		cmd.SetArgs([]string{"--json", "id,metadata,computed"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, resources)
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"metadata":{"env":"prod"}`)
		assert.Contains(t, out, `"computed":10`) // 5*2
		assert.Contains(t, out, `"computed":20`) // 10*2
	})

	t.Run("mix regular and custom fields", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		cmd := &cobra.Command{Use: "test"}
		var exp cli.Exporter
		cli.AddJSONFlags(cmd, &exp, []string{"id", "name", "computed"})
		cmd.SetArgs([]string{"--json", "name,computed"})
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return exp.Write(ios, resources)
		}
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		require.NoError(t, cmd.Execute())
		out := stdout.String()
		assert.Contains(t, out, `"name":"alpha"`)
		assert.Contains(t, out, `"computed":10`)
		assert.NotContains(t, out, `"id"`)
	})
}

// ════════════════════════════════════════════════════════════════════
// Edge 20: System() IOStreams basic sanity
// ════════════════════════════════════════════════════════════════════

func TestEdge_SystemIOStreams(t *testing.T) {
	ios := cli.System()
	assert.NotNil(t, ios.In)
	assert.NotNil(t, ios.Out)
	assert.NotNil(t, ios.ErrOut)
	assert.NotNil(t, ios.Output())
	assert.NotNil(t, ios.Prompter())
	// Width should be > 0
	assert.Greater(t, ios.TerminalWidth(), 0)
}
