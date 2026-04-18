// Package cli provides CLI enhancements for raystack applications.
//
// Usage:
//
//	rootCmd := &cobra.Command{Use: "frontier", Short: "identity management"}
//	rootCmd.AddCommand(serverCmd, userCmd)
//
//	cli.Init(rootCmd,
//	    cli.Version("0.1.0", "raystack/frontier"),
//	    cli.Topics(authTopic, envTopic),
//	)
//
//	rootCmd.Execute()
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

// Init enhances a cobra root command with standard CLI features:
// help, completion, reference docs, output/prompter context, and
// optionally a version command with update checking.
//
// The developer owns the root command — Init only adds features to it.
func Init(rootCmd *cobra.Command, opts ...Option) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Inject shared output and prompter into command context.
	existing := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(cmd.Context(), contextKey{}, &cliContext{
			output:   printer.NewOutput(os.Stdout),
			prompter: prompt.New(),
		})
		cmd.SetContext(ctx)
		if existing != nil {
			existing(cmd, args)
		}
	}

	// Wire commander features.
	var managerOpts []func(*commander.Manager)
	if len(cfg.topics) > 0 {
		managerOpts = append(managerOpts, commander.WithTopics(cfg.topics))
	}
	if len(cfg.hooks) > 0 {
		managerOpts = append(managerOpts, commander.WithHooks(cfg.hooks))
	}
	mgr := commander.New(rootCmd, managerOpts...)
	mgr.Init()

	// Add version command if configured.
	if cfg.version != "" {
		rootCmd.AddCommand(versionCmd(rootCmd.Name(), cfg.version, cfg.repo))
	}
}

func versionCmd(name, ver, repo string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, _ []string) {
			out := Output(cmd)
			out.Println(fmt.Sprintf("%s version %s", name, ver))
			if repo != "" {
				if msg := version.CheckForUpdate(ver, repo); msg != "" {
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
