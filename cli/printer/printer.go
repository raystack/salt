// Package printer provides terminal output utilities for CLI applications.
//
// Create an Output for your command and use it for all text, structured data,
// and progress indicators:
//
//	out := printer.NewOutput(os.Stdout)
//	out.Success("deployed to prod")
//	out.Table(rows)
//	out.JSON(data)
package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

// Output handles all terminal output for a CLI command.
//
// Data output (Table, JSON, YAML, Println) goes to the primary writer (stdout).
// Status output (Spin, Warning, Error, Info, Success) goes to the error writer (stderr).
// This separation ensures spinners and status messages don't corrupt
// piped data output (e.g. myapp list --json | jq).
type Output struct {
	w     io.Writer
	errW  io.Writer
	theme Theme
	tty   bool
}

// NewOutput creates a new Output that writes data to w and status to stderr.
// It auto-detects TTY and color support from the writer.
func NewOutput(w io.Writer) *Output {
	tty := false
	if f, ok := w.(*os.File); ok {
		tty = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return &Output{w: w, errW: os.Stderr, theme: newTheme(), tty: tty}
}

// --- Text output ---

// Success prints a green success message to stderr.
func (o *Output) Success(msg string) {
	fmt.Fprintln(o.errW, o.color(o.theme.Green, msg))
}

// Warning prints a yellow warning message to stderr.
func (o *Output) Warning(msg string) {
	fmt.Fprintln(o.errW, o.color(o.theme.Yellow, msg))
}

// Error prints a red error message to stderr.
func (o *Output) Error(msg string) {
	fmt.Fprintln(o.errW, o.color(o.theme.Red, msg))
}

// Info prints a cyan informational message to stderr.
func (o *Output) Info(msg string) {
	fmt.Fprintln(o.errW, o.color(o.theme.Cyan, msg))
}

// Bold prints a bold message.
func (o *Output) Bold(msg string) {
	fmt.Fprintln(o.w, termenv.String(msg).Bold().String())
}

// Print prints a plain message.
func (o *Output) Print(msg string) {
	fmt.Fprint(o.w, msg)
}

// Println prints a plain message with a newline.
func (o *Output) Println(msg string) {
	fmt.Fprintln(o.w, msg)
}

// --- Structured output ---

// JSON writes data as compact JSON.
func (o *Output) JSON(data interface{}) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintln(o.w, string(out))
	return nil
}

// PrettyJSON writes data as indented JSON.
func (o *Output) PrettyJSON(data interface{}) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(o.w, string(out))
	return nil
}

// YAML writes data as YAML.
func (o *Output) YAML(data interface{}) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprint(o.w, string(out))
	return nil
}

// Table writes rows as a tab-aligned table when the output is a TTY.
// When piped (non-TTY), it writes tab-separated values for easy
// processing with tools like awk, cut, or jq.
func (o *Output) Table(rows [][]string) {
	if o.tty {
		tw := tabwriter.NewWriter(o.w, 0, 0, 2, ' ', 0)
		for _, row := range rows {
			fmt.Fprintln(tw, strings.Join(row, "\t"))
		}
		tw.Flush()
	} else {
		for _, row := range rows {
			fmt.Fprintln(o.w, strings.Join(row, "\t"))
		}
	}
}

// --- Markdown ---

// Markdown renders and prints markdown text with terminal styling.
func (o *Output) Markdown(text string) error {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	tr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(0),
		glamour.WithStylesFromJSONBytes([]byte(`{"document":{"margin":0},"code_block":{"margin":0}}`)),
	)
	if err != nil {
		return err
	}
	rendered, err := tr.Render(text)
	if err != nil {
		return err
	}
	fmt.Fprint(o.w, rendered)
	return nil
}

// MarkdownWithWrap renders markdown with a specified word wrap width.
func (o *Output) MarkdownWithWrap(text string, wrap int) error {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	tr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(wrap),
		glamour.WithStylesFromJSONBytes([]byte(`{"document":{"margin":0},"code_block":{"margin":0}}`)),
	)
	if err != nil {
		return err
	}
	rendered, err := tr.Render(text)
	if err != nil {
		return err
	}
	fmt.Fprint(o.w, rendered)
	return nil
}

// --- Progress indicators ---

// Indicator wraps a terminal spinner.
type Indicator struct {
	spinner *spinner.Spinner
}

// Stop halts the spinner.
func (i *Indicator) Stop() {
	if i.spinner != nil {
		i.spinner.Stop()
	}
}

// Spin creates and starts a spinner. Returns a no-op indicator if not a TTY.
func (o *Output) Spin(label string) *Indicator {
	if !o.tty {
		return &Indicator{}
	}
	s := spinner.New(spinner.CharSets[11], 120*time.Millisecond, spinner.WithColor("fgCyan"))
	if label != "" {
		s.Prefix = label + " "
	}
	s.Writer = o.errW
	s.Start()
	return &Indicator{s}
}

// Progress creates a progress bar on stderr.
func (o *Output) Progress(max int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionSetWriter(o.errW),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
	)
}

// --- Formatting helpers (return styled strings for composition) ---

// Green returns text styled in green.
func Green(t string) string { return colorize(newTheme().Green, t) }

// Yellow returns text styled in yellow.
func Yellow(t string) string { return colorize(newTheme().Yellow, t) }

// Cyan returns text styled in cyan.
func Cyan(t string) string { return colorize(newTheme().Cyan, t) }

// Red returns text styled in red.
func Red(t string) string { return colorize(newTheme().Red, t) }

// Grey returns text styled in grey.
func Grey(t string) string { return colorize(newTheme().Grey, t) }

// Blue returns text styled in blue.
func Blue(t string) string { return colorize(newTheme().Blue, t) }

// Magenta returns text styled in magenta.
func Magenta(t string) string { return colorize(newTheme().Magenta, t) }

// --- Formatted color helpers ---

// Greenf returns formatted text styled in green.
func Greenf(format string, a ...interface{}) string { return Green(fmt.Sprintf(format, a...)) }

// Yellowf returns formatted text styled in yellow.
func Yellowf(format string, a ...interface{}) string { return Yellow(fmt.Sprintf(format, a...)) }

// Cyanf returns formatted text styled in cyan.
func Cyanf(format string, a ...interface{}) string { return Cyan(fmt.Sprintf(format, a...)) }

// Redf returns formatted text styled in red.
func Redf(format string, a ...interface{}) string { return Red(fmt.Sprintf(format, a...)) }

// Greyf returns formatted text styled in grey.
func Greyf(format string, a ...interface{}) string { return Grey(fmt.Sprintf(format, a...)) }

// Bluef returns formatted text styled in blue.
func Bluef(format string, a ...interface{}) string { return Blue(fmt.Sprintf(format, a...)) }

// Magentaf returns formatted text styled in magenta.
func Magentaf(format string, a ...interface{}) string { return Magenta(fmt.Sprintf(format, a...)) }

// Italic returns text styled in italic.
func Italic(t string) string { return termenv.String(t).Italic().String() }

// Icon returns a symbol for the given name: "success"→✔, "failure"→✘, "info"→ℹ, "warning"→⚠.
func Icon(name string) string {
	icons := map[string]string{"failure": "✘", "success": "✔", "info": "ℹ", "warning": "⚠"}
	return icons[name]
}

// --- Theme ---

// Theme defines terminal colors.
type Theme struct {
	Green   termenv.Color
	Yellow  termenv.Color
	Cyan    termenv.Color
	Red     termenv.Color
	Grey    termenv.Color
	Blue    termenv.Color
	Magenta termenv.Color
}

func newTheme() Theme {
	tp := termenv.EnvColorProfile()
	if !termenv.HasDarkBackground() {
		return Theme{
			Green: tp.Color("#005F00"), Yellow: tp.Color("#FFAF00"),
			Cyan: tp.Color("#0087FF"), Red: tp.Color("#D70000"),
			Grey: tp.Color("#303030"), Blue: tp.Color("#000087"),
			Magenta: tp.Color("#AF00FF"),
		}
	}
	return Theme{
		Green: tp.Color("#A8CC8C"), Yellow: tp.Color("#DBAB79"),
		Cyan: tp.Color("#66C2CD"), Red: tp.Color("#E88388"),
		Grey: tp.Color("#B9BFCA"), Blue: tp.Color("#71BEF2"),
		Magenta: tp.Color("#D290E4"),
	}
}

func (o *Output) color(c termenv.Color, t string) string {
	return colorize(c, t)
}

func colorize(c termenv.Color, t string) string {
	return termenv.String(t).Foreground(c).String()
}
