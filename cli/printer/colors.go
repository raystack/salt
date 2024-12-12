package printer

import (
	"fmt"
	"github.com/muesli/termenv"
)

var tp = termenv.EnvColorProfile()

// Theme defines a collection of colors for terminal outputs.
type Theme struct {
	Green   termenv.Color
	Yellow  termenv.Color
	Cyan    termenv.Color
	Red     termenv.Color
	Grey    termenv.Color
	Blue    termenv.Color
	Magenta termenv.Color
}

var themes = map[string]Theme{
	"light": {
		Green:   tp.Color("#005F00"),
		Yellow:  tp.Color("#FFAF00"),
		Cyan:    tp.Color("#0087FF"),
		Red:     tp.Color("#D70000"),
		Grey:    tp.Color("#303030"),
		Blue:    tp.Color("#000087"),
		Magenta: tp.Color("#AF00FF"),
	},
	"dark": {
		Green:   tp.Color("#A8CC8C"),
		Yellow:  tp.Color("#DBAB79"),
		Cyan:    tp.Color("#66C2CD"),
		Red:     tp.Color("#E88388"),
		Grey:    tp.Color("#B9BFCA"),
		Blue:    tp.Color("#71BEF2"),
		Magenta: tp.Color("#D290E4"),
	},
}

// NewTheme initializes a Theme based on the terminal background (light or dark).
func NewTheme() Theme {
	if !termenv.HasDarkBackground() {
		return themes["light"]
	}
	return themes["dark"]
}

var theme = NewTheme()

// formatColorize applies the given color to the formatted text.
func formatColorize(color termenv.Color, t string, args ...interface{}) string {
	return colorize(color, fmt.Sprintf(t, args...))
}

func Green(t ...string) string {
	return colorize(theme.Green, t...)
}

func Greenf(t string, args ...interface{}) string {
	return formatColorize(theme.Green, t, args...)
}

func Yellow(t ...string) string {
	return colorize(theme.Yellow, t...)
}

func Yellowf(t string, args ...interface{}) string {
	return formatColorize(theme.Yellow, t, args...)
}

func Cyan(t ...string) string {
	return colorize(theme.Cyan, t...)
}

func Cyanf(t string, args ...interface{}) string {
	return formatColorize(theme.Cyan, t, args...)
}

func Red(t ...string) string {
	return colorize(theme.Red, t...)
}

func Redf(t string, args ...interface{}) string {
	return formatColorize(theme.Red, t, args...)
}

func Grey(t ...string) string {
	return colorize(theme.Grey, t...)
}

func Greyf(t string, args ...interface{}) string {
	return formatColorize(theme.Grey, t, args...)
}

func Blue(t ...string) string {
	return colorize(theme.Blue, t...)
}

func Bluef(t string, args ...interface{}) string {
	return formatColorize(theme.Blue, t, args...)
}

func Magenta(t ...string) string {
	return colorize(theme.Magenta, t...)
}

func Magentaf(t string, args ...interface{}) string {
	return formatColorize(theme.Magenta, t, args...)
}

func Icon(name string) string {
	icons := map[string]string{"failure": "✘", "success": "✔", "info": "ℹ", "warning": "⚠"}
	if icon, exists := icons[name]; exists {
		return icon
	}
	return ""
}

// colorize applies the given color to the text.
func colorize(color termenv.Color, t ...string) string {
	return termenv.String(t...).Foreground(color).String()
}
