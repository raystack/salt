package cli

import "github.com/raystack/salt/cli/commander"

type config struct {
	version string
	repo    string
	topics  []commander.HelpTopic
	hooks   []commander.HookBehavior
}

// Option configures cli.Init.
type Option func(*config)

// Version enables a version command with update checking.
// The repo should be in "owner/repo" format (e.g. "raystack/frontier").
func Version(ver, repo string) Option {
	return func(c *config) {
		c.version = ver
		c.repo = repo
	}
}

// Topics adds help topics to the CLI.
func Topics(topics ...commander.HelpTopic) Option {
	return func(c *config) { c.topics = append(c.topics, topics...) }
}

// Hooks adds hook behaviors applied to commands.
func Hooks(hooks ...commander.HookBehavior) Option {
	return func(c *config) { c.hooks = append(c.hooks, hooks...) }
}
