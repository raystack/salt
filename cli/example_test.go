package cli_test

import (
	"fmt"

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

func ExampleIO() {
	deleteCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			ios := cli.IO(cmd)

			// Guard interactive prompts in non-TTY environments.
			if !ios.CanPrompt() {
				return fmt.Errorf("--yes flag required in non-interactive mode")
			}

			ok, _ := ios.Prompter().Confirm("Delete resource?", false)
			if !ok {
				return cli.ErrCancel
			}

			ios.Output().Success("deleted")
			return nil
		},
	}

	rootCmd := &cobra.Command{Use: "myapp"}
	rootCmd.AddCommand(deleteCmd)
	cli.Init(rootCmd)
	cli.Execute(rootCmd)
}

func ExampleTest() {
	// Use cli.Test() in unit tests to capture output.
	ios, _, stdout, _ := cli.Test()
	ios.SetStdoutTTY(true) // simulate a terminal

	out := ios.Output()
	out.Println("hello from test")

	fmt.Print(stdout.String())
	// Output: hello from test
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
