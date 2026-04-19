package cli_test

import (
	"github.com/raystack/salt/cli"
	"github.com/raystack/salt/cli/commander"
	"github.com/spf13/cobra"
)

func ExampleInit() {
	rootCmd := &cobra.Command{
		Use:   "frontier",
		Short: "identity management",
	}
	rootCmd.PersistentFlags().String("host", "", "API server host")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cli.Output(cmd)
			out.Table([][]string{
				{"ID", "NAME"},
				{"1", "Alice"},
			})
			return nil
		},
	}
	rootCmd.AddCommand(listCmd)

	cli.Init(rootCmd,
		cli.Version("0.1.0", "raystack/frontier"),
	)

	cli.Execute(rootCmd)
}

func ExampleExecute() {
	rootCmd := &cobra.Command{
		Use: "myapp",
		RunE: func(cmd *cobra.Command, _ []string) error {
			p := cli.Prompter(cmd)
			ok, _ := p.Confirm("Continue?", true)
			if !ok {
				return cli.ErrCancel // exit 0, no output
			}
			return nil
		},
	}

	cli.Init(rootCmd)
	cli.Execute(rootCmd)
}

func ExampleInit_withTopics() {
	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "my application",
	}

	cli.Init(rootCmd,
		cli.Topics(
			commander.HelpTopic{
				Name:  "auth",
				Short: "How authentication works",
				Long:  "Detailed explanation of authentication...",
			},
		),
	)

	cli.Execute(rootCmd)
}
