package cli_test

import (
	"github.com/raystack/salt/cli"
	"github.com/raystack/salt/cli/commander"
	"github.com/spf13/cobra"
)

func ExampleExecute() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cli.Output(cmd)
			out.Table([][]string{
				{"ID", "NAME"},
				{"1", "Alice"},
				{"2", "Bob"},
			})
			return nil
		},
	}

	cli.Execute(
		cli.Name("myapp"),
		cli.Description("my application"),
		cli.Version("0.1.0", "raystack/myapp"),
		cli.Commands(listCmd),
	)
}

func ExampleExecute_withTopics() {
	cli.Execute(
		cli.Name("myapp"),
		cli.Description("my application"),
		cli.Commands(),
		cli.Topics(
			commander.HelpTopic{
				Name:  "auth",
				Short: "How authentication works",
				Long:  "Detailed explanation of authentication...",
			},
		),
	)
}
