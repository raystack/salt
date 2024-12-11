package cmdx

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AddHelpTopics adds all configured help topics to the CLI.
//
// Help topics provide detailed information about specific subjects,
// such as environment variables or configuration.
func (m *Commander) AddHelpTopics() {
	for _, topic := range m.Topics {
		m.addHelpTopicCommand(topic)
	}
}

// addHelpTopicCommand adds a single help topic command to the CLI.
func (m *Commander) addHelpTopicCommand(topic HelpTopic) {
	helpCmd := &cobra.Command{
		Use:     topic.Name,
		Short:   topic.Short,
		Long:    topic.Long,
		Example: topic.Example,
		Hidden:  false,
		Annotations: map[string]string{
			"group": "help",
		},
	}

	helpCmd.SetHelpFunc(helpTopicHelpFunc)
	helpCmd.SetUsageFunc(helpTopicUsageFunc)

	m.RootCmd.AddCommand(helpCmd)
}

// helpTopicHelpFunc customizes the help message for a help topic command.
func helpTopicHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), cmd.Long)
	if cmd.Example != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "\nEXAMPLES")
		fmt.Fprintln(cmd.OutOrStdout(), indent(cmd.Example, "  "))
	}
}

// helpTopicUsageFunc customizes the usage message for a help topic command.
func helpTopicUsageFunc(cmd *cobra.Command) error {
	fmt.Fprintf(cmd.OutOrStdout(), "Usage: %s help %s\n", cmd.Root().Name(), cmd.Use)
	return nil
}
