package cmdx

import (
	"github.com/spf13/cobra"
)

// SetHelpTopicCmd creates a custom help topic command.
//
// This function allows you to define a help topic that provides detailed information
// about a specific subject. The generated command should be added to the root command.
//
// Parameters:
//   - title: The name of the help topic (e.g., "env").
//   - topic: A map containing the following keys:
//   - "short": A brief description of the topic.
//   - "long": A detailed explanation of the topic.
//   - "example": An example usage of the topic.
//
// Returns:
//   - A pointer to the configured help topic `cobra.Command`.
//
// Example:
//
//	topic := map[string]string{
//	    "short": "Environment variables help",
//	    "long": "Details about environment variables used by the CLI.",
//	    "example": "$ mycli help env",
//	}
//	rootCmd.AddCommand(cmdx.SetHelpTopicCmd("env", topic))
func SetHelpTopicCmd(title string, topic map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     title,
		Short:   topic["short"],
		Long:    topic["long"],
		Example: topic["example"],
		Hidden:  false,
		Annotations: map[string]string{
			"group": "help",
		},
	}

	cmd.SetHelpFunc(helpTopicHelpFunc)
	cmd.SetUsageFunc(helpTopicUsageFunc)

	return cmd
}

func helpTopicHelpFunc(command *cobra.Command, args []string) {
	command.Print(command.Long)
	if command.Example != "" {
		command.Printf("\nEXAMPLES\n")
		command.Print(indent(command.Example, "  "))
	}
}

func helpTopicUsageFunc(command *cobra.Command) error {
	command.Printf("Usage: %s help %s\n", command.Root().Name(), command.Use)
	return nil
}
