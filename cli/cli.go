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
//	cli.Execute(rootCmd)
package cli

import (
	"context"
	"errors"
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
	cfg := &options{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Set error prefix for consistent error messages.
	rootCmd.SetErrPrefix(rootCmd.Name() + ":")

	// Silence cobra's default error and usage printing.
	// Errors are handled by Execute; usage is shown only for flag errors.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

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

	// Wrap flag parsing errors so Execute can show contextual usage.
	// Must be set after mgr.Init() which also configures a flag error func.
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return &flagError{err: err}
	})

	// Add version command if configured.
	if cfg.version != "" {
		rootCmd.AddCommand(versionCmd(rootCmd.Name(), cfg.version, cfg.repo))
	}
}

// Execute runs the root command and handles errors with appropriate
// exit codes and output. It uses ExecuteC to obtain the failing command
// so flag errors can show contextual usage.
//
// This function never returns on error — it calls os.Exit.
func Execute(rootCmd *cobra.Command) {
	cmd, err := rootCmd.ExecuteC()
	if err == nil {
		return
	}

	var flagErr *flagError
	switch {
	case errors.Is(err, ErrCancel):
		os.Exit(0)
	case errors.Is(err, ErrSilent):
		os.Exit(1)
	case errors.As(err, &flagErr):
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, cmd.UsageString())
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
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
	if ctx := cmd.Context(); ctx != nil {
		if cc, ok := ctx.Value(contextKey{}).(*cliContext); ok {
			return cc.output
		}
	}
	return printer.NewOutput(os.Stdout)
}

// Prompter extracts the shared prompter from a command's context.
func Prompter(cmd *cobra.Command) prompt.Prompter {
	if ctx := cmd.Context(); ctx != nil {
		if cc, ok := ctx.Value(contextKey{}).(*cliContext); ok {
			return cc.prompter
		}
	}
	return prompt.New()
}
