package term

import (
	"fmt"

	"github.com/muesli/termenv"
)

var tp = termenv.EnvColorProfile()

// Theme represents a color theme.
type Theme struct {
	ColorGreen   termenv.Color
	ColorYellow  termenv.Color
	ColorCyan    termenv.Color
	ColorRed     termenv.Color
	ColorGrey    termenv.Color
	ColorBlue    termenv.Color
	ColorMagenta termenv.Color
}

var themes = map[string]Theme{
	"light": {
		ColorGreen:   tp.Color("#005F00"),
		ColorYellow:  tp.Color("#FFAF00"),
		ColorCyan:    tp.Color("#0087FF"),
		ColorRed:     tp.Color("#D70000"),
		ColorGrey:    tp.Color("#303030"),
		ColorBlue:    tp.Color("#000087"),
		ColorMagenta: tp.Color("#AF00FF"),
	},
	"dark": {
		ColorGreen:   tp.Color("#A8CC8C"),
		ColorYellow:  tp.Color("#DBAB79"),
		ColorCyan:    tp.Color("#66C2CD"),
		ColorRed:     tp.Color("#E88388"),
		ColorGrey:    tp.Color("#B9BFCA"),
		ColorBlue:    tp.Color("#71BEF2"),
		ColorMagenta: tp.Color("#D290E4"),
	},
}

// ColorScheme is a color scheme.
type ColorScheme struct {
	theme Theme
}

// NewColorScheme returns a new ColorScheme with the given theme.
func NewColorScheme() *ColorScheme {
	if !termenv.HasDarkBackground() {
		return &ColorScheme{
			theme: themes["light"],
		}
	}
	return &ColorScheme{
		theme: themes["dark"],
	}
}

func (c *ColorScheme) Bold(t string) string {
	return termenv.String(t).Bold().String()
}

func (c *ColorScheme) Boldf(t string, args ...interface{}) string {
	return c.Bold(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Italic(t string) string {
	return termenv.String(t).Italic().String()
}

func (c *ColorScheme) Italicf(t string, args ...interface{}) string {
	return c.Italic(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Green(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorGreen).String()
}

func (c *ColorScheme) Greenf(t string, args ...interface{}) string {
	return c.Green(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Yellow(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorYellow).String()
}

func (c *ColorScheme) Yellowf(t string, args ...interface{}) string {
	return c.Yellow(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Cyan(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorCyan).String()
}

func (c *ColorScheme) Cyanf(t string, args ...interface{}) string {
	return c.Cyan(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Red(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorRed).String()
}

func (c *ColorScheme) Redf(t string, args ...interface{}) string {
	return c.Red(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Grey(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorGrey).String()
}

func (c *ColorScheme) Greyf(t string, args ...interface{}) string {
	return c.Grey(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Blue(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorBlue).String()
}

func (c *ColorScheme) Bluef(t string, args ...interface{}) string {
	return c.Blue(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Magenta(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorMagenta).String()
}

func (c *ColorScheme) Magentaf(t string, args ...interface{}) string {
	return c.Magenta(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) SuccessIcon() string {
	return termenv.String("✓").Foreground(c.theme.ColorGreen).String()
}

func (c *ColorScheme) WarningIcon() string {
	return termenv.String("!").Foreground(c.theme.ColorYellow).String()
}

func (c *ColorScheme) FailureIcon() string {
	return termenv.String("✘").Foreground(c.theme.ColorRed).String()
}
