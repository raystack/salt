package printer

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// RenderOpts is a type alias for a slice of glamour.TermRendererOption,
// representing the rendering options for the markdown renderer.
type RenderOpts []glamour.TermRendererOption

// This ensures the rendered markdown has no extra indentation or margins, providing a compact view.
func withoutIndentation() glamour.TermRendererOption {
	overrides := []byte(`
	  {
			"document": {
				"margin": 0
			},
			"code_block": {
				"margin": 0
			}
	  }`)

	return glamour.WithStylesFromJSONBytes(overrides)
}

// withoutWrap ensures the rendered markdown does not wrap lines, useful for wide terminals.
func withoutWrap() glamour.TermRendererOption {
	return glamour.WithWordWrap(0)
}

// render applies the given rendering options to the provided markdown text.
func render(text string, opts RenderOpts) (string, error) {
	// Ensure input text uses consistent line endings.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	tr, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return "", err
	}

	return tr.Render(text)
}

// Markdown renders the given markdown text with default options.
func Markdown(text string) (string, error) {
	opts := RenderOpts{
		glamour.WithAutoStyle(), // Automatically determine styling based on terminal settings.
		glamour.WithEmoji(),     // Enable emoji rendering.
		withoutIndentation(),    // Disable indentation for a compact view.
		withoutWrap(),           // Disable word wrapping.
	}

	return render(text, opts)
}

// MarkdownWithWrap renders the given markdown text with a specified word wrapping width.
func MarkdownWithWrap(text string, wrap int) (string, error) {
	opts := RenderOpts{
		glamour.WithAutoStyle(),    // Automatically determine styling based on terminal settings.
		glamour.WithEmoji(),        // Enable emoji rendering.
		glamour.WithWordWrap(wrap), // Enable word wrapping with the specified width.
		withoutIndentation(),       // Disable indentation for a compact view.
	}

	return render(text, opts)
}
