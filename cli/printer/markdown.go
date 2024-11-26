package printer

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

type RenderOpts []glamour.TermRendererOption

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

func withoutWrap() glamour.TermRendererOption {
	return glamour.WithWordWrap(0)
}

func render(text string, opts RenderOpts) (string, error) {
	// Glamour rendering preserves carriage return characters in code blocks, but
	// we need to ensure that no such characters are present in the output.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	tr, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return "", err
	}

	return tr.Render(text)
}

func Markdown(text string) (string, error) {
	opts := RenderOpts{
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		withoutIndentation(),
		withoutWrap(),
	}

	return render(text, opts)
}

func MarkdownWithWrap(text string, wrap int) (string, error) {
	opts := RenderOpts{
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(wrap),
		withoutIndentation(),
	}

	return render(text, opts)
}
