package cmdx

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// AddCompletionCommand adds a `completion` command to the CLI.
//
// The completion command generates shell completion scripts for Bash, Zsh,
// Fish, and PowerShell.
//
// Example:
//
//	manager := cmdx.NewManager(rootCmd)
//	manager.AddCompletionCommand()
//
// Usage:
//
//	$ mycli completion bash
//	$ mycli completion zsh
func (m *Manager) AddCompletionCommand() {
	summary := m.generateCompletionSummary(m.RootCmd.Use)

	completionCmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate shell completion scripts",
		Long:                  summary,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run:                   m.runCompletionCommand,
	}

	m.RootCmd.AddCommand(completionCmd)
}

// runCompletionCommand executes the appropriate shell completion generation logic.
func (m *Manager) runCompletionCommand(cmd *cobra.Command, args []string) {
	switch args[0] {
	case "bash":
		cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	}
}

// generateCompletionSummary creates the long description for the `completion` command.
func (m *Manager) generateCompletionSummary(exec string) string {
	var execs []interface{}
	for i := 0; i < 12; i++ {
		execs = append(execs, exec)
	}
	return heredoc.Docf(`To load completions:
		`+"```"+`
		Bash:

		  $ source <(%s completion bash)

		  # To load completions for each session, execute once:
		  # Linux:
		  $ %s completion bash > /etc/bash_completion.d/%s
		  # macOS:
		  $ %s completion bash > /usr/local/etc/bash_completion.d/%s

		Zsh:

		  # If shell completion is not already enabled in your environment,
		  # you will need to enable it.  You can execute the following once:

		  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

		  # To load completions for each session, execute once:
		  $ %s completion zsh > "${fpath[1]}/_yourprogram"

		  # You will need to start a new shell for this setup to take effect.

		Fish:

		  $ %s completion fish | source

		  # To load completions for each session, execute once:
		  $ %s completion fish > ~/.config/fish/completions/%s.fish

		PowerShell:

		  PS> %s completion powershell | Out-String | Invoke-Expression

		  # To load completions for every new session, run:
		  PS> %s completion powershell > %s.ps1
		  # and source this file from your PowerShell profile.
		`+"```"+`
	`, execs...)
}
