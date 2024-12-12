package commander

import (
	"strings"

	"github.com/spf13/cobra"
)

// Manager manages and configures features for a CLI tool.
type Manager struct {
	RootCmd    *cobra.Command
	Help       bool           // Enable custom help.
	Reference  bool           // Enable reference command.
	Completion bool           // Enable shell completion.
	Config     bool           // Enable configuration management.
	Docs       bool           // Enable markdown documentation
	Hooks      []HookBehavior // Hook behaviors to apply to commands
	Topics     []HelpTopic    // Help topics with their details.
}

// HelpTopic defines a single help topic with its details.
type HelpTopic struct {
	Name    string
	Short   string
	Long    string
	Example string
}

// HookBehavior defines a specific behavior applied to commands.
type HookBehavior struct {
	Name     string                   // Name of the hook (e.g., "setup", "auth").
	Behavior func(cmd *cobra.Command) // Function to apply to commands.
}

// New creates a new CLI Manager using the provided root command and optional configurations.
//
// Parameters:
// - rootCmd: The root Cobra command for the CLI.
// - options: Functional options for configuring the Manager.
//
// Example:
//
//	rootCmd := &cobra.Command{Use: "mycli"}
//	manager := cmdx.NewCommander(rootCmd, cmdx.WithTopics(...), cmdx.WithHooks(...))
func New(rootCmd *cobra.Command, options ...func(*Manager)) *Manager {
	// Create Manager with defaults
	manager := &Manager{
		RootCmd:    rootCmd,
		Help:       true,  // Default enabled
		Reference:  true,  // Default enabled
		Completion: true,  // Default enabled
		Docs:       false, // Default disabled
		Topics:     []HelpTopic{},
		Hooks:      []HookBehavior{},
	}

	// Apply functional options
	for _, opt := range options {
		opt(manager)
	}

	return manager
}

// Init sets up the CLI features based on the Manager's configuration.
// It enables or disables features like custom help, reference documentation,
// shell completion, help topics, and client hooks based on the Manager's settings.
func (m *Manager) Init() {
	if m.Help {
		m.setCustomHelp()
	}
	if m.Reference {
		m.addReferenceCommand()
	}
	if m.Completion {
		m.addCompletionCommand()
	}
	if m.Docs {
		m.addMarkdownCommand("./docs")
	}
	if len(m.Topics) > 0 {
		m.addHelpTopics()
	}

	if len(m.Hooks) > 0 {
		m.addClientHooks()
	}
}

// WithTopics sets the help topics for the Manager.
func WithTopics(topics []HelpTopic) func(*Manager) {
	return func(m *Manager) {
		m.Topics = topics
	}
}

// WithHooks sets the hook behaviors for the Manager.
func WithHooks(hooks []HookBehavior) func(*Manager) {
	return func(m *Manager) {
		m.Hooks = hooks
	}
}

// IsCommandErr checks if the given error is related to a Cobra command error.
// This is useful for distinguishing between user errors (e.g., incorrect commands or flags)
// and program errors, allowing the application to display appropriate messages.
func IsCommandErr(err error) bool {
	if err == nil {
		return false
	}

	// Known Cobra command error keywords
	cmdErrorKeywords := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
	}

	errMessage := err.Error()
	for _, keyword := range cmdErrorKeywords {
		if strings.Contains(errMessage, keyword) {
			return true
		}
	}
	return false
}
