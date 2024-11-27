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

// This ensures the rendered markdown does not wrap lines, useful for wide terminals.
func withoutWrap() glamour.TermRendererOption {
	return glamour.WithWordWrap(0)
}

// render applies the given rendering options to the provided markdown text.
func render(text string, opts RenderOpts) (string, error) {
	// Ensure input text uses consistent line endings.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// Create a new terminal renderer with the provided options.
	tr, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return "", err
	}

	// Render the markdown text and return the result.
	return tr.Render(text)
}

// Markdown renders the given markdown text with default options.
//
// This includes automatic styling, emoji rendering, no indentation, and no word wrapping.
//
// Parameters:
//   - text: The markdown text to render.
//
// Returns:
//   - The rendered markdown string.
//   - An error if rendering fails.
//
// Example Usage:
//
//	output, err := printer.Markdown("# Hello, Markdown!")
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
//
// Parameters:
//   - text: The markdown text to render.
//   - wrap: The desired word wrapping width (e.g., 80 for 80 characters).
//
// Returns:
//   - The rendered markdown string.
//   - An error if rendering fails.
//
// Example Usage:
//
//	output, err := printer.MarkdownWithWrap("# Hello, Markdown!", 80)
func MarkdownWithWrap(text string, wrap int) (string, error) {
	opts := RenderOpts{
		glamour.WithAutoStyle(),    // Automatically determine styling based on terminal settings.
		glamour.WithEmoji(),        // Enable emoji rendering.
		glamour.WithWordWrap(wrap), // Enable word wrapping with the specified width.
		withoutIndentation(),       // Disable indentation for a compact view.
	}

	return render(text, opts)
}
