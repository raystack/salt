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

func bold(t ...string) string {
	return termenv.String(t...).Bold().String()
}

func boldf(t string, args ...interface{}) string {
	return bold(fmt.Sprintf(t, args...))
}

func italic(t ...string) string {
	return termenv.String(t...).Italic().String()
}

func italicf(t string, args ...interface{}) string {
	return italic(fmt.Sprintf(t, args...))
}

func Green(t ...string) string {
	return termenv.String(t...).Foreground(theme.Green).String()
}

func Greenf(t string, args ...interface{}) string {
	return Green(fmt.Sprintf(t, args...))
}

func Yellow(t ...string) string {
	return termenv.String(t...).Foreground(theme.Yellow).String()
}

func Yellowf(t string, args ...interface{}) string {
	return Yellow(fmt.Sprintf(t, args...))
}

func Cyan(t ...string) string {
	return termenv.String(t...).Foreground(theme.Cyan).String()
}

func Cyanf(t string, args ...interface{}) string {
	return Cyan(fmt.Sprintf(t, args...))
}

func Red(t ...string) string {
	return termenv.String(t...).Foreground(theme.Red).String()
}

func Redf(t string, args ...interface{}) string {
	return Red(fmt.Sprintf(t, args...))
}

func Grey(t ...string) string {
	return termenv.String(t...).Foreground(theme.Grey).String()
}

func Greyf(t string, args ...interface{}) string {
	return Grey(fmt.Sprintf(t, args...))
}

func Blue(t ...string) string {
	return termenv.String(t...).Foreground(theme.Blue).String()
}

func Bluef(t string, args ...interface{}) string {
	return Blue(fmt.Sprintf(t, args...))
}

func Magenta(t ...string) string {
	return termenv.String(t...).Foreground(theme.Magenta).String()
}

func Magentaf(t string, args ...interface{}) string {
	return Magenta(fmt.Sprintf(t, args...))
}

func FailureIcon() string {
	return termenv.String("✘").String()
}
