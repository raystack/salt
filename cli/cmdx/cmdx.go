package cmdx

import (
	"strings"

	"github.com/spf13/cobra"
)

// Commander manages and configures features for a CLI tool.
type Commander struct {
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

// NewCommander creates a new CLI Commander using the provided root command and optional configurations.
//
// Parameters:
// - rootCmd: The root Cobra command for the CLI.
// - options: Functional options for configuring the Commander.
//
// Example:
//
//	rootCmd := &cobra.Command{Use: "mycli"}
//	manager := cmdx.NewCommander(rootCmd, cmdx.WithTopics(...), cmdx.WithHooks(...))
func NewCommander(rootCmd *cobra.Command, options ...func(*Commander)) *Commander {
	// Create Commander with defaults
	manager := &Commander{
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

// Init sets up the CLI features based on the Commander's configuration.
//
// It enables or disables features like custom help, reference documentation,
// shell completion, help topics, and client hooks based on the Commander's settings.
func (m *Commander) Init() {
	if m.Help {
		m.SetCustomHelp()
	}
	if m.Reference {
		m.AddReferenceCommand()
	}
	if m.Completion {
		m.AddCompletionCommand()
	}
	if m.Docs {
		m.AddMarkdownCommand("./docs")
	}
	if len(m.Topics) > 0 {
		m.AddHelpTopics()
	}

	if len(m.Hooks) > 0 {
		m.AddClientHooks()
	}
}

// WithTopics sets the help topics for the Commander.
func WithTopics(topics []HelpTopic) func(*Commander) {
	return func(m *Commander) {
		m.Topics = topics
	}
}

// WithHooks sets the hook behaviors for the Commander.
func WithHooks(hooks []HookBehavior) func(*Commander) {
	return func(m *Commander) {
		m.Hooks = hooks
	}
}

// IsCLIErr checks if the given error is related to a Cobra command error.
//
// This is useful for distinguishing between user errors (e.g., incorrect commands or flags)
// and program errors, allowing the application to display appropriate messages.
func IsCLIErr(err error) bool {
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
