package commander

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/raystack/salt/cli/printer"

	"github.com/spf13/cobra"
)

// addReferenceCommand adds a `reference` command to the CLI.
// The `reference` command generates markdown documentation for all commands
// in the CLI command tree.
func (m *Manager) addReferenceCommand() {
	var isPlain bool
	refCmd := &cobra.Command{
		Use:   "reference",
		Short: "Comprehensive reference of all commands",
		Long:  m.generateReferenceMarkdown(),
		Run:   m.runReferenceCommand(&isPlain),
		Annotations: map[string]string{
			"group": "help",
		},
	}
	refCmd.SetHelpFunc(m.runReferenceCommand(&isPlain))
	refCmd.Flags().BoolVarP(&isPlain, "plain", "p", true, "output in plain markdown (without ANSI color)")

	m.RootCmd.AddCommand(refCmd)
}

// runReferenceCommand handles the output generation for the `reference` command.
// It renders the documentation either as plain markdown or with ANSI color.
func (m *Manager) runReferenceCommand(isPlain *bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var (
			output string
			err    error
		)

		if *isPlain {
			output = cmd.Long
		} else {
			output, err = printer.Markdown(cmd.Long)
			if err != nil {
				fmt.Println("Error generating markdown:", err)
				return
			}
		}

		fmt.Print(output)
	}
}

// generateReferenceMarkdown generates a complete markdown representation
// of the command tree for the `reference` command.
func (m *Manager) generateReferenceMarkdown() string {
	buf := bytes.NewBufferString(fmt.Sprintf("# %s reference\n\n", m.RootCmd.Name()))
	for _, c := range m.RootCmd.Commands() {
		if c.Hidden {
			continue
		}
		m.generateCommandReference(buf, c, 2)
	}
	return buf.String()
}

// generateCommandReference recursively generates markdown for a given command
// and its subcommands.
func (m *Manager) generateCommandReference(w io.Writer, cmd *cobra.Command, depth int) {
	// Name + Description
	fmt.Fprintf(w, "%s `%s`\n\n", strings.Repeat("#", depth), cmd.UseLine())
	fmt.Fprintf(w, "%s\n\n", cmd.Short)

	// Flags
	if flagUsages := cmd.Flags().FlagUsages(); flagUsages != "" {
		fmt.Fprintf(w, "```\n%s```\n\n", dedent(flagUsages))
	}

	// Subcommands
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		m.generateCommandReference(w, c, depth+1)
	}
}
