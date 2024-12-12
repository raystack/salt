package commander

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Section Titles for Help Output
const (
	usage     = "Usage"
	corecmd   = "Core commands"
	othercmd  = "Other commands"
	helpcmd   = "Help topics"
	flags     = "Flags"
	iflags    = "Inherited flags"
	arguments = "Arguments"
	examples  = "Examples"
	envs      = "Environment variables"
	learn     = "Learn more"
	feedback  = "Feedback"
)

// setCustomHelp configures a custom help function for the CLI.
// The custom help function organizes commands into sections and provides
// detailed error messages for incorrect flag usage.
func (m *Manager) setCustomHelp() {
	m.RootCmd.PersistentFlags().Bool("help", false, "Show help for command")

	m.RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		displayHelp(cmd, args)
	})
	m.RootCmd.SetUsageFunc(generateUsage)
	m.RootCmd.SetFlagErrorFunc(handleFlagError)
}

// generateUsage customizes the usage function for a command.
func generateUsage(cmd *cobra.Command) error {
	cmd.Printf("Usage:  %s\n", cmd.UseLine())

	subcommands := cmd.Commands()
	if len(subcommands) > 0 {
		cmd.Print("\nAvailable commands:\n")
		for _, subCmd := range subcommands {
			if !subCmd.Hidden {
				cmd.Printf("  %s\n", subCmd.Name())
			}
		}
	}

	flagUsages := cmd.LocalFlags().FlagUsages()
	if flagUsages != "" {
		cmd.Println("\nFlags:")
		cmd.Print(indent(dedent(flagUsages), "  "))
	}
	return nil
}

// handleFlagError processes flag-related errors, including the special case of help flags.
func handleFlagError(cmd *cobra.Command, err error) error {
	if errors.Is(err, pflag.ErrHelp) {
		return err
	}
	return err
}

// displayHelp generates a custom help message for a Cobra command.
func displayHelp(cmd *cobra.Command, args []string) {
	if isRootCommand(cmd.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		showSuggestions(cmd, args[1])
		return
	}

	helpEntries := buildHelpEntries(cmd)
	printHelpEntries(cmd, helpEntries)
}

// buildHelpEntries constructs a structured help message for a command.
func buildHelpEntries(cmd *cobra.Command) []helpEntry {
	var coreCommands, helpCommands, otherCommands []string
	groupCommands := map[string][]string{}

	for _, c := range cmd.Commands() {
		if c.Short == "" || c.Hidden {
			continue
		}

		entry := fmt.Sprintf("%s%s", rpad(c.Name(), c.NamePadding()+3), c.Short)
		if group, ok := c.Annotations["group"]; ok {
			switch group {
			case "core":
				coreCommands = append(coreCommands, entry)
			case "help":
				helpCommands = append(helpCommands, entry)
			default:
				groupCommands[group] = append(groupCommands[group], entry)
			}
		} else {
			otherCommands = append(otherCommands, entry)
		}
	}

	// Treat all commands as core if no groups are specified
	if len(coreCommands) == 0 && len(groupCommands) == 0 {
		coreCommands = otherCommands
		otherCommands = []string{}
	}

	helpEntries := []helpEntry{}
	if text := cmd.Long; text != "" {
		helpEntries = append(helpEntries, helpEntry{"", text})
	}

	helpEntries = append(helpEntries, helpEntry{usage, cmd.UseLine()})
	if len(coreCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{corecmd, strings.Join(coreCommands, "\n")})
	}
	for group, cmds := range groupCommands {
		helpEntries = append(helpEntries, helpEntry{fmt.Sprintf("%s commands", toTitle(group)), strings.Join(cmds, "\n")})
	}
	if len(otherCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{othercmd, strings.Join(otherCommands, "\n")})
	}
	if len(helpCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{helpcmd, strings.Join(helpCommands, "\n")})
	}
	if flagUsages := cmd.LocalFlags().FlagUsages(); flagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{flags, dedent(flagUsages)})
	}
	if inheritedFlagUsages := cmd.InheritedFlags().FlagUsages(); inheritedFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{iflags, dedent(inheritedFlagUsages)})
	}
	if argsAnnotation, ok := cmd.Annotations["help:arguments"]; ok {
		helpEntries = append(helpEntries, helpEntry{arguments, argsAnnotation})
	}
	if cmd.Example != "" {
		helpEntries = append(helpEntries, helpEntry{examples, cmd.Example})
	}
	if argsAnnotation, ok := cmd.Annotations["help:environment"]; ok {
		helpEntries = append(helpEntries, helpEntry{envs, argsAnnotation})
	}
	if argsAnnotation, ok := cmd.Annotations["help:learn"]; ok {
		helpEntries = append(helpEntries, helpEntry{learn, argsAnnotation})
	}
	if argsAnnotation, ok := cmd.Annotations["help:feedback"]; ok {
		helpEntries = append(helpEntries, helpEntry{feedback, argsAnnotation})
	}
	return helpEntries
}

// printHelpEntries displays help entries to the command's output.
func printHelpEntries(cmd *cobra.Command, entries []helpEntry) {
	out := cmd.OutOrStdout()
	for _, entry := range entries {
		if entry.Title != "" {
			fmt.Fprintln(out, bold(entry.Title))
			fmt.Fprintln(out, indent(strings.Trim(entry.Body, "\r\n"), "  "))
		} else {
			fmt.Fprintln(out, entry.Body)
		}
		fmt.Fprintln(out)
	}
}

// showSuggestions displays suggestions for mistyped subcommands.
func showSuggestions(cmd *cobra.Command, arg string) {
	cmd.Printf("unknown command %q for %q\n", arg, cmd.CommandPath())

	var suggestions []string
	if arg == "help" {
		suggestions = []string{"--help"}
	} else {
		if cmd.SuggestionsMinimumDistance <= 0 {
			cmd.SuggestionsMinimumDistance = 2
		}
		suggestions = cmd.SuggestionsFor(arg)
	}

	if len(suggestions) > 0 {
		cmd.Println("\nDid you mean this?")
		for _, suggestion := range suggestions {
			cmd.Printf("  %s\n", suggestion)
		}
	}

	cmd.Println()
	_ = generateUsage(cmd)
}

// isRootCommand checks if the given command is the root command.
func isRootCommand(cmd *cobra.Command) bool {
	return cmd != nil && !cmd.HasParent()
}

// Utility types and functions
type helpEntry struct {
	Title string
	Body  string
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

var lineRE = regexp.MustCompile(`(?m)^`)

func indent(s, indent string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}

func toTitle(text string) string {
	heading := cases.Title(language.Und).String(text)
	return heading
}

func bold(text string) termenv.Style {
	h := termenv.String(text).Bold()
	return h
}
