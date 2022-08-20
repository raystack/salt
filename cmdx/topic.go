package cmdx

import (
	"github.com/spf13/cobra"
)

// SetHelpTopicCmd sets the help topic command.
// This should be added on the root command. e.g.
//
//	topic := map[string]string{"short": "Env","long": "Environment","example": "example",}
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
		command.Printf("\n\nEXAMPLES\n")
		command.Print(indent(command.Example, "  "))
	}
}

func helpTopicUsageFunc(command *cobra.Command) error {
	command.Printf("Usage: %s help %s", command.Root().Name(), command.Use)
	return nil
}
