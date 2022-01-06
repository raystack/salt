package cmdx

import (
	"github.com/spf13/cobra"
)

func SetHelpTopic(title string, topic map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     title,
		Short:   topic["short"],
		Long:    topic["long"],
		Example: topic["example"],
		Hidden:  true,
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
