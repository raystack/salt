package commander

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// addCompletionCommand adds a `completion` command to the CLI.
// The `completion` command generates shell completion scripts
// for Bash, Zsh, Fish, and PowerShell.
// Usage:
//
//	$ mycli completion bash
//	$ mycli completion zsh
func (m *Manager) addCompletionCommand() {
	summary := m.generateCompletionSummary(m.RootCmd.Use)

	completionCmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate shell completion scripts",
		Long:                  summary,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(out)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(out)
			}
			return nil
		},
	}

	m.RootCmd.AddCommand(completionCmd)
}

// generateCompletionSummary creates the long description for the `completion` command.
func (m *Manager) generateCompletionSummary(exec string) string {
	var execs []any
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
