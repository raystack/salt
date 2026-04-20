package cli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/raystack/salt/cli"
	"github.com/raystack/salt/cli/commander"
	"github.com/raystack/salt/cli/printer"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Sample 1: Minimal CRUD CLI ────────────────────────────────────

// Simulates a typical "resource manager" CLI like `frontier user list`.
func TestSample_ResourceManager(t *testing.T) {
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "admin"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "viewer"},
	}

	buildCLI := func() (*cobra.Command, *cli.Exporter) {
		var exporter cli.Exporter

		listCmd := &cobra.Command{
			Use:   "list",
			Short: "List all users",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if exporter != nil {
					return exporter.Write(cli.IO(cmd), users)
				}
				out := cli.Output(cmd)
				rows := [][]string{{"ID", "NAME", "EMAIL", "ROLE"}}
				for _, u := range users {
					rows = append(rows, []string{
						fmt.Sprintf("%d", u.ID), u.Name, u.Email, u.Role,
					})
				}
				out.Table(rows)
				return nil
			},
		}
		cli.AddJSONFlags(listCmd, &exporter, []string{"id", "name", "email", "role"})

		getCmd := &cobra.Command{
			Use:   "get [id]",
			Short: "Get a user by ID",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := cli.Output(cmd)
				out.Success(fmt.Sprintf("User: %s", args[0]))
				return nil
			},
		}

		deleteCmd := &cobra.Command{
			Use:   "delete [id]",
			Short: "Delete a user",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ios := cli.IO(cmd)
				if !ios.CanPrompt() {
					yes, _ := cmd.Flags().GetBool("yes")
					if !yes {
						return fmt.Errorf("--yes flag required in non-interactive mode")
					}
				}
				ios.Output().Success(fmt.Sprintf("deleted user %s", args[0]))
				return nil
			},
		}
		deleteCmd.Flags().BoolP("yes", "y", false, "skip confirmation")

		userCmd := &cobra.Command{
			Use:     "user",
			Short:   "Manage users",
			GroupID: "core",
		}
		userCmd.AddCommand(listCmd, getCmd, deleteCmd)

		rootCmd := &cobra.Command{
			Use:   "frontier",
			Short: "Identity management",
		}
		rootCmd.AddGroup(&cobra.Group{ID: "core", Title: "Core commands"})
		rootCmd.AddCommand(userCmd)

		cli.Init(rootCmd, cli.Version("1.0.0", ""))

		return rootCmd, &exporter
	}

	t.Run("table output", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "list"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "Alice")
		assert.Contains(t, stdout.String(), "Bob")
	})

	t.Run("json output with field selection", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "list", "--json", "id,name"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.NoError(t, err)
		out := stdout.String()
		assert.Contains(t, out, `"id"`)
		assert.Contains(t, out, `"name"`)
		assert.NotContains(t, out, `"email"`)
		assert.NotContains(t, out, `"role"`)
	})

	t.Run("json rejects unknown field", func(t *testing.T) {
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "list", "--json", "id,bogus"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown JSON field")
	})

	t.Run("delete requires --yes in non-interactive", func(t *testing.T) {
		ios, _, _, _ := cli.Test() // non-TTY
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "delete", "1"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--yes")
	})

	t.Run("delete succeeds with --yes", func(t *testing.T) {
		ios, _, _, stderr := cli.Test()
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "delete", "1", "--yes"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.NoError(t, err)
		assert.Contains(t, stderr.String(), "deleted user 1")
	})

	t.Run("get with args", func(t *testing.T) {
		ios, _, _, stderr := cli.Test()
		root, _ := buildCLI()
		root.SetArgs([]string{"user", "get", "42"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.NoError(t, err)
		assert.Contains(t, stderr.String(), "User: 42")
	})

	t.Run("version command works", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root, _ := buildCLI()
		root.SetArgs([]string{"version"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		err := root.Execute()
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "frontier version 1.0.0")
	})

	t.Run("help shows grouped commands", func(t *testing.T) {
		root, _ := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"--help"})
		root.Execute()
		help := buf.String()
		assert.Contains(t, help, "user")
	})

	t.Run("unknown subcommand gives suggestion", func(t *testing.T) {
		root, _ := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"usr"})
		root.Execute()
		// Should suggest "user"
	})
}

// ─── Sample 2: Output format variations ────────────────────────────

func TestSample_OutputFormats(t *testing.T) {
	data := map[string]any{
		"name":    "test-project",
		"version": "2.0.0",
		"tags":    []string{"go", "cli"},
	}

	t.Run("JSON compact (piped)", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		out := ios.Output()
		err := out.JSON(data)
		require.NoError(t, err)
		assert.NotContains(t, stdout.String(), "\n  ") // compact
		assert.Contains(t, stdout.String(), "test-project")
	})

	t.Run("JSON pretty", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		out := ios.Output()
		err := out.PrettyJSON(data)
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "  ") // indented
	})

	t.Run("YAML output", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		out := ios.Output()
		err := out.YAML(data)
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "name: test-project")
	})

	t.Run("table output", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		out := ios.Output()
		rows := [][]string{
			{"NAME", "VERSION"},
			{"alpha", "1.0"},
			{"beta", "2.0"},
		}
		out.Table(rows)
		assert.Contains(t, stdout.String(), "alpha")
		assert.Contains(t, stdout.String(), "beta")
	})

	t.Run("status messages go to stderr", func(t *testing.T) {
		ios, _, stdout, stderr := cli.Test()
		out := ios.Output()
		out.Success("done")
		out.Warning("careful")
		out.Error("oops")
		out.Info("fyi")
		out.Bold("heading")
		out.Println("data goes here")

		// Data should be on stdout only
		assert.Contains(t, stdout.String(), "data goes here")
		assert.NotContains(t, stdout.String(), "done")

		// Status should be on stderr only
		assert.Contains(t, stderr.String(), "done")
		assert.Contains(t, stderr.String(), "careful")
		assert.Contains(t, stderr.String(), "oops")
		assert.Contains(t, stderr.String(), "fyi")
		assert.Contains(t, stderr.String(), "heading")
	})

	t.Run("color functions return styled strings", func(t *testing.T) {
		// Just ensure they don't panic and return non-empty strings
		assert.NotEmpty(t, printer.Green("ok"))
		assert.NotEmpty(t, printer.Red("fail"))
		assert.NotEmpty(t, printer.Yellow("warn"))
		assert.NotEmpty(t, printer.Cyan("info"))
		assert.NotEmpty(t, printer.Grey("muted"))
		assert.NotEmpty(t, printer.Blue("link"))
		assert.NotEmpty(t, printer.Magenta("highlight"))
		assert.NotEmpty(t, printer.Italic("emphasis"))

		// Formatted variants
		assert.Contains(t, printer.Greenf("count: %d", 42), "42")
		assert.Contains(t, printer.Redf("error: %s", "bad"), "bad")
	})

	t.Run("icons return expected symbols", func(t *testing.T) {
		assert.Equal(t, "✔", printer.Icon("success"))
		assert.Equal(t, "✘", printer.Icon("failure"))
		assert.Equal(t, "ℹ", printer.Icon("info"))
		assert.Equal(t, "⚠", printer.Icon("warning"))
		assert.Equal(t, "", printer.Icon("unknown"))
	})
}

// ─── Sample 3: CLI with help topics and hooks ──────────────────────

func TestSample_TopicsAndHooks(t *testing.T) {
	buildCLI := func() *cobra.Command {
		rootCmd := &cobra.Command{
			Use:   "myapp",
			Short: "My application",
			Long:  "A sample application to test help topics and hooks.",
		}

		listCmd := &cobra.Command{
			Use:   "list",
			Short: "List items",
			RunE: func(cmd *cobra.Command, _ []string) error {
				cli.Output(cmd).Println("items listed")
				return nil
			},
		}
		rootCmd.AddCommand(listCmd)

		cli.Init(rootCmd,
			cli.Version("3.0.0", ""),
			cli.Topics(
				commander.HelpTopic{
					Name:    "auth",
					Short:   "How authentication works",
					Long:    "This app uses OAuth2 for authentication.\nSet MYAPP_TOKEN to authenticate.",
					Example: "  export MYAPP_TOKEN=abc123\n  myapp list",
				},
				commander.HelpTopic{
					Name:  "environment",
					Short: "Environment variables",
					Long:  "MYAPP_TOKEN: API token\nMYAPP_HOST: API host",
				},
			),
		)

		return rootCmd
	}

	t.Run("help topic listed", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"--help"})
		root.Execute()
		help := buf.String()
		assert.Contains(t, help, "auth")
		assert.Contains(t, help, "environment")
	})

	t.Run("help topic shows details", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"auth"})
		root.Execute()
		out := buf.String()
		assert.Contains(t, out, "OAuth2")
	})

	t.Run("completion command exists", func(t *testing.T) {
		root := buildCLI()
		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "completion" {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("reference command exists", func(t *testing.T) {
		root := buildCLI()
		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == "reference" {
				found = true
			}
		}
		assert.True(t, found)
	})
}

// ─── Sample 4: Error handling patterns ─────────────────────────────

func TestSample_ErrorHandling(t *testing.T) {
	t.Run("ErrSilent suppresses output", func(t *testing.T) {
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				cli.Output(cmd).Error("something went wrong")
				return cli.ErrSilent
			},
		}
		cli.Init(root)

		// Can't test os.Exit, but verify the error is returned
		root.SetArgs([]string{})
		err := root.Execute()
		assert.ErrorIs(t, err, cli.ErrSilent)
	})

	t.Run("ErrCancel for user cancellation", func(t *testing.T) {
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return cli.ErrCancel
			},
		}
		cli.Init(root)

		root.SetArgs([]string{})
		err := root.Execute()
		assert.ErrorIs(t, err, cli.ErrCancel)
	})

	t.Run("flag errors include usage context", func(t *testing.T) {
		root := &cobra.Command{
			Use: "app",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return nil
			},
		}
		root.Flags().Int("port", 8080, "port number")
		cli.Init(root)

		root.SetArgs([]string{"--port", "notanumber"})
		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid argument")
	})

	t.Run("missing required args", func(t *testing.T) {
		sub := &cobra.Command{
			Use:  "deploy [env]",
			Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cli.Output(cmd).Success(fmt.Sprintf("deployed to %s", args[0]))
				return nil
			},
		}
		root := &cobra.Command{Use: "app"}
		root.AddCommand(sub)
		cli.Init(root)

		root.SetArgs([]string{"deploy"})
		err := root.Execute()
		require.Error(t, err)
	})
}

// ─── Sample 5: Testing with IOStreams ──────────────────────────────

func TestSample_TestingPatterns(t *testing.T) {
	t.Run("capture table output in test", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()

		// Simulate what a command's RunE would do
		out := ios.Output()
		out.Table([][]string{
			{"ID", "NAME"},
			{"1", "Alice"},
			{"2", "Bob"},
		})

		lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
		assert.Len(t, lines, 3) // header + 2 rows
	})

	t.Run("inject IOStreams into command context", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()

		cmd := &cobra.Command{
			Use: "test",
			RunE: func(cmd *cobra.Command, _ []string) error {
				cli.Output(cmd).Println("captured")
				return nil
			},
		}
		cli.Init(cmd)

		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{})
		require.NoError(t, cmd.Execute())
		assert.Contains(t, stdout.String(), "captured")
	})

	t.Run("CanPrompt false in test by default", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		assert.False(t, ios.CanPrompt())
	})

	t.Run("simulate TTY for prompt testing", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		ios.SetStdinTTY(true)
		ios.SetStdoutTTY(true)
		assert.True(t, ios.CanPrompt())
	})

	t.Run("terminal width defaults to 80 in tests", func(t *testing.T) {
		ios, _, _, _ := cli.Test()
		assert.Equal(t, 80, ios.TerminalWidth())
	})
}

// ─── Sample 6: Embedded struct export ──────────────────────────────

func TestSample_EmbeddedStructExport(t *testing.T) {
	type Base struct {
		ID        int    `json:"id"`
		CreatedAt string `json:"created_at"`
	}
	type Project struct {
		Base
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	p := Project{
		Base:  Base{ID: 1, CreatedAt: "2024-01-01"},
		Name:  "salt",
		Owner: "raystack",
	}

	t.Run("exports top-level fields", func(t *testing.T) {
		data := cli.StructExportData(p, []string{"name", "owner"})
		assert.Equal(t, "salt", data["name"])
		assert.Equal(t, "raystack", data["owner"])
	})

	t.Run("exports embedded fields", func(t *testing.T) {
		data := cli.StructExportData(p, []string{"id", "created_at"})
		assert.Equal(t, 1, data["id"])
		assert.Equal(t, "2024-01-01", data["created_at"])
	})

	t.Run("mixes embedded and top-level", func(t *testing.T) {
		data := cli.StructExportData(p, []string{"id", "name"})
		assert.Equal(t, 1, data["id"])
		assert.Equal(t, "salt", data["name"])
	})
}

// ─── Sample 7: ConfigCommand ───────────────────────────────────────

func TestSample_ConfigCommand(t *testing.T) {
	type AppConfig struct {
		Host string `yaml:"host" default:"localhost"`
		Port int    `yaml:"port" default:"8080"`
	}

	t.Run("config command has init and list", func(t *testing.T) {
		cmd := cli.ConfigCommand("testapp", &AppConfig{})
		assert.NotNil(t, cmd)

		var names []string
		for _, sub := range cmd.Commands() {
			names = append(names, sub.Name())
		}
		assert.Contains(t, names, "init")
		assert.Contains(t, names, "list")
	})
}

// ─── Sample 8: Complex multi-group CLI ─────────────────────────────

func TestSample_MultiGroupCLI(t *testing.T) {
	buildCLI := func() *cobra.Command {
		root := &cobra.Command{
			Use:   "platform",
			Short: "Platform management CLI",
		}

		root.AddGroup(
			&cobra.Group{ID: "resources", Title: "Resource commands"},
			&cobra.Group{ID: "admin", Title: "Admin commands"},
		)

		// Resource commands
		for _, name := range []string{"project", "dataset", "job"} {
			cmd := &cobra.Command{
				Use:     name,
				Short:   fmt.Sprintf("Manage %ss", name),
				GroupID: "resources",
				RunE: func(cmd *cobra.Command, _ []string) error {
					cli.Output(cmd).Println(cmd.Name())
					return nil
				},
			}
			root.AddCommand(cmd)
		}

		// Admin commands
		for _, name := range []string{"user", "policy"} {
			cmd := &cobra.Command{
				Use:     name,
				Short:   fmt.Sprintf("Manage %ss", name),
				GroupID: "admin",
				RunE: func(cmd *cobra.Command, _ []string) error {
					cli.Output(cmd).Println(cmd.Name())
					return nil
				},
			}
			root.AddCommand(cmd)
		}

		cli.Init(root, cli.Version("2.0.0", ""))
		return root
	}

	t.Run("all groups appear in help", func(t *testing.T) {
		root := buildCLI()
		var buf strings.Builder
		root.SetOut(&buf)
		root.SetArgs([]string{"--help"})
		root.Execute()
		help := buf.String()
		assert.Contains(t, help, "Resources commands")
		assert.Contains(t, help, "Admin commands")
	})

	t.Run("commands execute correctly", func(t *testing.T) {
		ios, _, stdout, _ := cli.Test()
		root := buildCLI()
		root.SetArgs([]string{"project"})
		ctx := context.WithValue(context.Background(), cli.ContextKey(), ios)
		root.SetContext(ctx)
		require.NoError(t, root.Execute())
		assert.Contains(t, stdout.String(), "project")
	})
}

// ─── Sample 9: PreRunE hook chaining ───────────────────────────────

func TestSample_PreRunHookChaining(t *testing.T) {
	t.Run("Init preserves existing PersistentPreRunE", func(t *testing.T) {
		var hookCalled bool
		root := &cobra.Command{
			Use: "app",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				hookCalled = true
				return nil
			},
			RunE: func(cmd *cobra.Command, _ []string) error {
				return nil
			},
		}
		cli.Init(root)

		root.SetArgs([]string{})
		require.NoError(t, root.Execute())
		assert.True(t, hookCalled, "existing PersistentPreRunE should be called")
	})

	t.Run("Init preserves existing PersistentPreRun", func(t *testing.T) {
		var hookCalled bool
		root := &cobra.Command{
			Use: "app",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				hookCalled = true
			},
			RunE: func(cmd *cobra.Command, _ []string) error {
				return nil
			},
		}
		cli.Init(root)

		root.SetArgs([]string{})
		require.NoError(t, root.Execute())
		assert.True(t, hookCalled, "existing PersistentPreRun should be called")
	})

	t.Run("AddJSONFlags preserves PreRun", func(t *testing.T) {
		var hookCalled bool
		cmd := &cobra.Command{
			Use: "test",
			PreRun: func(cmd *cobra.Command, args []string) {
				hookCalled = true
			},
			RunE: func(cmd *cobra.Command, _ []string) error { return nil },
		}

		var exporter cli.Exporter
		cli.AddJSONFlags(cmd, &exporter, []string{"id"})
		cmd.SetArgs([]string{})
		require.NoError(t, cmd.Execute())
		assert.True(t, hookCalled, "PreRun should still be called after AddJSONFlags")
	})
}
