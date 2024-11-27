package cmdx

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/raystack/salt/cli/printer"

	"github.com/spf13/cobra"
)

// SetRefCmd adds a `reference` command to the root command to generate
// comprehensive reference documentation for the command tree.
//
// The `reference` command outputs the documentation in markdown format
// and supports a `--plain` flag to control whether ANSI colors are used.
func SetRefCmd(root *cobra.Command) *cobra.Command {
	var isPlain bool
	cmd := &cobra.Command{
		Use:   "reference",
		Short: "Comprehensive reference of all commands",
		Long:  referenceLong(root),
		Run:   referenceHelpFn(&isPlain),
		Annotations: map[string]string{
			"group": "help",
		},
	}
	cmd.SetHelpFunc(referenceHelpFn(&isPlain))
	cmd.Flags().BoolVarP(&isPlain, "plain", "p", true, "output in plain markdown (without ansi color)")
	return cmd
}

// referenceHelpFn generates the output for the `reference` command.
// It renders the documentation either as plain markdown or with ANSI color.
func referenceHelpFn(isPlain *bool) func(*cobra.Command, []string) {
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

// referenceLong generates the complete reference documentation
// for the command tree in markdown format.
func referenceLong(cmd *cobra.Command) string {
	buf := bytes.NewBufferString(fmt.Sprintf("# %s reference\n\n", cmd.Name()))
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		generateCommandReference(buf, c, 2)
	}
	return buf.String()
}

func generateCommandReference(w io.Writer, cmd *cobra.Command, depth int) {
	// Name + Description
	fmt.Fprintf(w, "%s `%s`\n\n", strings.Repeat("#", depth), cmd.UseLine())
	fmt.Fprintf(w, "%s\n\n", cmd.Short)

	if flagUsages := cmd.Flags().FlagUsages(); flagUsages != "" {
		fmt.Fprintf(w, "```\n%s````\n\n", dedent(flagUsages))
	}

	// Subcommands
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		generateCommandReference(w, c, depth+1)
	}
}
