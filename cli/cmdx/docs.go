package cmdx

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// GenerateMarkdownTree generate cobra cmd commands tree as markdown file
// rootOutputPath determines the folder where the markdown files are written
func GenerateMarkdownTree(rootOutputPath string, cmd *cobra.Command) error {
	dirFilePath := filepath.Join(rootOutputPath, cmd.Name())
	if len(cmd.Commands()) != 0 {
		if _, err := os.Stat(dirFilePath); os.IsNotExist(err) {
			if err := os.Mkdir(dirFilePath, os.ModePerm); err != nil {
				return err
			}
		}
		for _, subCmd := range cmd.Commands() {
			GenerateMarkdownTree(dirFilePath, subCmd)
		}
	} else {
		outFilePath := filepath.Join(rootOutputPath, cmd.Name())
		outFilePath = outFilePath + ".md"

		f, err := os.Create(outFilePath)
		if err != nil {
			return err
		}

		return doc.GenMarkdownCustom(cmd, f, func(s string) string {
			return filepath.Join(dirFilePath, s)
		})
	}

	return nil
}
