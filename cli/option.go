package cli

import (
	"github.com/raystack/salt/cli/commander"
	"github.com/spf13/cobra"
)

// Option configures a CLI application.
type Option func(*CLI)

// Name sets the CLI application name (required).
func Name(name string) Option {
	return func(c *CLI) { c.name = name }
}

// Description sets the CLI application description.
func Description(desc string) Option {
	return func(c *CLI) { c.description = desc }
}

// Version sets the version string and GitHub repo for update checking.
// The repo should be in "owner/repo" format (e.g. "raystack/frontier").
func Version(ver, repo string) Option {
	return func(c *CLI) {
		c.version = ver
		c.repo = repo
	}
}

// Commands adds subcommands to the CLI.
func Commands(cmds ...*cobra.Command) Option {
	return func(c *CLI) { c.commands = append(c.commands, cmds...) }
}

// Topics adds help topics to the CLI.
func Topics(topics ...commander.HelpTopic) Option {
	return func(c *CLI) { c.topics = append(c.topics, topics...) }
}

// Hooks adds hook behaviors applied to commands.
func Hooks(hooks ...commander.HookBehavior) Option {
	return func(c *CLI) { c.hooks = append(c.hooks, hooks...) }
}
