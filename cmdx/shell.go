package cmdx

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// SetCompletionCmd is used to generate the completion script in
// bash, zsh, fish, and powershell. It should be added on the root
// command and can be used as `completion bash` or `completion zsh`.
func SetCompletionCmd(exec string) *cobra.Command {
	var execs []interface{}
	for i := 0; i < 12; i++ {
		execs = append(execs, exec)
	}
	summary := heredoc.Docf(`To load completions:
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

		fish:

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

	return &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate shell completion scripts",
		Long:                  summary,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
}
