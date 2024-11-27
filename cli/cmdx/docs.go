package cmdx

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// GenerateMarkdownTree generates a Markdown documentation tree for all commands
// in the given Cobra command hierarchy.
//
// Parameters:
//   - rootOutputPath: The root directory where the Markdown files will be generated.
//   - cmd: The root Cobra command whose hierarchy will be documented.
//
// Returns:
//   - An error if any part of the process (file creation, directory creation) fails.
//
// Example Usage:
//
//	rootCmd := &cobra.Command{Use: "mycli"}
//	cmdx.GenerateMarkdownTree("./docs", rootCmd)
func GenerateMarkdownTree(rootOutputPath string, cmd *cobra.Command) error {
	dirFilePath := filepath.Join(rootOutputPath, cmd.Name())
	if len(cmd.Commands()) > 0 {
		if err := ensureDir(dirFilePath); err != nil {
			return fmt.Errorf("failed to create directory for command %q: %w", cmd.Name(), err)
		}
		for _, subCmd := range cmd.Commands() {
			if err := GenerateMarkdownTree(dirFilePath, subCmd); err != nil {
				return err
			}
		}
	} else {
		outFilePath := filepath.Join(rootOutputPath, cmd.Name())
		outFilePath = outFilePath + ".md"

		f, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer f.Close()

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
