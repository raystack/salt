package cli

import "github.com/raystack/salt/cli/commander"

type options struct {
	version string
	repo    string
	topics  []commander.HelpTopic
	hooks   []commander.HookBehavior
}

// Option configures cli.Init.
type Option func(*options)

// Version enables a version command with update checking.
// The repo should be in "owner/repo" format (e.g. "raystack/frontier").
func Version(ver, repo string) Option {
	return func(c *options) {
		c.version = ver
		c.repo = repo
	}
}

// Topics adds help topics to the CLI.
func Topics(topics ...commander.HelpTopic) Option {
	return func(c *options) { c.topics = append(c.topics, topics...) }
}

// Hooks adds hook behaviors applied to commands.
func Hooks(hooks ...commander.HookBehavior) Option {
	return func(c *options) { c.hooks = append(c.hooks, hooks...) }
}
