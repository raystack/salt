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

var cs = NewColorScheme()

func Bold(t ...string) string {
	return termenv.String(t...).Bold().String()
}

func Boldf(t string, args ...interface{}) string {
	return Bold(fmt.Sprintf(t, args...))
}

func Italic(t ...string) string {
	return termenv.String(t...).Italic().String()
}

func Italicf(t string, args ...interface{}) string {
	return Italic(fmt.Sprintf(t, args...))
}

func Green(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorGreen).String()
}

func Greenf(t string, args ...interface{}) string {
	return Green(fmt.Sprintf(t, args...))
}

func Yellow(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorYellow).String()
}

func Yellowf(t string, args ...interface{}) string {
	return Yellow(fmt.Sprintf(t, args...))
}

func Cyan(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorCyan).String()
}

func Cyanf(t string, args ...interface{}) string {
	return Cyan(fmt.Sprintf(t, args...))
}

func Red(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorRed).String()
}

func Redf(t string, args ...interface{}) string {
	return Red(fmt.Sprintf(t, args...))
}

func Grey(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorGrey).String()
}

func Greyf(t string, args ...interface{}) string {
	return Grey(fmt.Sprintf(t, args...))
}

func Blue(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorBlue).String()
}

func Bluef(t string, args ...interface{}) string {
	return Blue(fmt.Sprintf(t, args...))
}

func Magenta(t ...string) string {
	return termenv.String(t...).Foreground(cs.theme.ColorMagenta).String()
}

func Magentaf(t string, args ...interface{}) string {
	return Magenta(fmt.Sprintf(t, args...))
}

func SuccessIcon() string {
	return termenv.String("✓").String()
}

func WarningIcon() string {
	return termenv.String("!").String()
}

func FailureIcon() string {
	return termenv.String("✘").String()
}
