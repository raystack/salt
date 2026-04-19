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
	rootCmd.PersistentFlags().StringP("host", "h", "", "API server host")

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

	if err := rootCmd.Execute(); err != nil {
		cli.HandleError(err)
	}
}

func ExampleHandleError() {
	rootCmd := &cobra.Command{
		Use: "myapp",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Return ErrSilent if you already printed the error
			out := cli.Output(cmd)
			out.Error("connection failed: timeout")
			return cli.ErrSilent

			// Return ErrCancel if user cancelled
			// return cli.ErrCancel

			// Return FlagError for bad input
			// return cli.NewFlagError(fmt.Errorf("--port must be positive"))
		},
	}

	cli.Init(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		cli.HandleError(err)
	}
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

	rootCmd.Execute()
}
