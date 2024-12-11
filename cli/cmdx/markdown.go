package cmdx

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// AddMarkdownCommand integrates a hidden `markdown` command into the root command.
// This command generates a Markdown documentation tree for all commands in the hierarchy.
func (m *Commander) AddMarkdownCommand(outputPath string) {
	markdownCmd := &cobra.Command{
		Use:    "markdown",
		Short:  "Generate Markdown documentation for all commands",
		Hidden: true,
		Annotations: map[string]string{
			"group": "help",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.generateMarkdownTree(outputPath, m.RootCmd)
		},
	}

	m.RootCmd.AddCommand(markdownCmd)
}

// generateMarkdownTree generates a Markdown documentation tree for the given command hierarchy.
//
// Parameters:
//   - rootOutputPath: The root directory where the Markdown files will be generated.
//   - cmd: The root Cobra command whose hierarchy will be documented.
//
// Returns:
//   - An error if any part of the process (file creation, directory creation) fails.
func (m *Commander) generateMarkdownTree(rootOutputPath string, cmd *cobra.Command) error {
	dirFilePath := filepath.Join(rootOutputPath, cmd.Name())

	// Handle subcommands by creating a directory and iterating through subcommands.
	if len(cmd.Commands()) > 0 {
		if err := ensureDir(dirFilePath); err != nil {
			return fmt.Errorf("failed to create directory for command %q: %w", cmd.Name(), err)
		}
		for _, subCmd := range cmd.Commands() {
			if err := m.generateMarkdownTree(dirFilePath, subCmd); err != nil {
				return err
			}
		}
	} else {
		// Generate a Markdown file for leaf commands.
		outFilePath := filepath.Join(rootOutputPath, cmd.Name()+".md")

		f, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Generate Markdown with a custom link handler.
		return doc.GenMarkdownCustom(cmd, f, func(s string) string {
			return filepath.Join(dirFilePath, s)
		})
	}

	return nil
}

// ensureDir ensures that the given directory exists, creating it if necessary.
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
