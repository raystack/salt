// Package cli provides a CLI bootstrap for raystack applications.
//
// Usage:
//
//	cli.Execute(
//	    cli.Name("frontier"),
//	    cli.Description("identity management"),
//	    cli.Version("0.1.0", "raystack/frontier"),
//	    cli.Commands(serverCmd, configCmd),
//	)
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/raystack/salt/cli/commander"
	"github.com/raystack/salt/cli/printer"
	"github.com/raystack/salt/cli/prompt"
	"github.com/raystack/salt/cli/version"
	"github.com/spf13/cobra"
)

type contextKey struct{}

type cliContext struct {
	output   *printer.Output
	prompter prompt.Prompter
}

// CLI holds the configured CLI application.
type CLI struct {
	name        string
	description string
	version     string
	repo        string
	commands    []*cobra.Command
	topics      []commander.HelpTopic
	hooks       []commander.HookBehavior
}

// Execute creates and runs a CLI application with sensible defaults.
// Help, completion, and reference commands are enabled automatically.
func Execute(opts ...Option) error {
	c, err := New(opts...)
	if err != nil {
		return err
	}
	return c.execute()
}

// New creates a CLI without executing, for advanced wiring.
func New(opts ...Option) (*CLI, error) {
	c := &CLI{}
	for _, opt := range opts {
		opt(c)
	}
	if c.name == "" {
		return nil, fmt.Errorf("cli: Name is required")
	}
	return c, nil
}

func (c *CLI) execute() error {
	rootCmd := &cobra.Command{
		Use:   c.name,
		Short: c.description,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			ctx := context.WithValue(cmd.Context(), contextKey{}, &cliContext{
				output:   printer.NewOutput(os.Stdout),
				prompter: prompt.New(),
			})
			cmd.SetContext(ctx)
		},
		SilenceUsage: true,
	}

	// Wire commander features.
	var managerOpts []func(*commander.Manager)
	if len(c.topics) > 0 {
		managerOpts = append(managerOpts, commander.WithTopics(c.topics))
	}
	if len(c.hooks) > 0 {
		managerOpts = append(managerOpts, commander.WithHooks(c.hooks))
	}
	mgr := commander.New(rootCmd, managerOpts...)
	mgr.Init()

	// Add version command if configured.
	if c.version != "" {
		rootCmd.AddCommand(c.versionCmd())
	}

	// Add user commands.
	rootCmd.AddCommand(c.commands...)

	return rootCmd.Execute()
}

func (c *CLI) versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, _ []string) {
			out := Output(cmd)
			out.Println(fmt.Sprintf("%s version %s", c.name, c.version))
			if c.repo != "" {
				if msg := version.CheckForUpdate(c.version, c.repo); msg != "" {
					out.Warning(msg)
				}
			}
		},
	}
}

// Output extracts the shared printer from a command's context.
func Output(cmd *cobra.Command) *printer.Output {
	if ctx, ok := cmd.Context().Value(contextKey{}).(*cliContext); ok {
		return ctx.output
	}
	return printer.NewOutput(os.Stdout)
}

// Prompter extracts the shared prompter from a command's context.
func Prompter(cmd *cobra.Command) prompt.Prompter {
	if ctx, ok := cmd.Context().Value(contextKey{}).(*cliContext); ok {
		return ctx.prompter
	}
	return prompt.New()
}
