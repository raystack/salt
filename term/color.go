package term

import (
	"fmt"

	"github.com/muesli/termenv"
)

var tp = termenv.EnvColorProfile()

type Theme struct {
	ColorSuccess termenv.Color
	ColorWarn    termenv.Color
	ColorInfo    termenv.Color
	ColorError   termenv.Color
	ColorBody    termenv.Color
	ColorPrimary termenv.Color
	ColorMagenta termenv.Color
}

var themes = map[string]Theme{
	"light": {
		ColorSuccess: tp.Color("#005F00"),
		ColorWarn:    tp.Color("#FFAF00"),
		ColorInfo:    tp.Color("#0087FF"),
		ColorError:   tp.Color("#D70000"),
		ColorBody:    tp.Color("#303030"),
		ColorPrimary: tp.Color("#000087"),
		ColorMagenta: tp.Color("#AF00FF"),
	},
	"dark": {
		ColorSuccess: tp.Color("#A8CC8C"),
		ColorWarn:    tp.Color("#DBAB79"),
		ColorInfo:    tp.Color("#66C2CD"),
		ColorError:   tp.Color("#E88388"),
		ColorBody:    tp.Color("#B9BFCA"),
		ColorPrimary: tp.Color("#71BEF2"),
		ColorMagenta: tp.Color("#D290E4"),
	},
}

type ColorScheme struct {
	theme Theme
}

func NewColorScheme() *ColorScheme {
	if !termenv.HasDarkBackground() {
		return &ColorScheme{
			theme: themes["light"],
		}
	}
	return &ColorScheme{
		theme: themes["light"],
	}
}

func (c *ColorScheme) Success(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorSuccess).String()
}

func (c *ColorScheme) Successf(t string, args ...interface{}) string {
	return c.Success(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Warn(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorWarn).String()
}

func (c *ColorScheme) Warnf(t string, args ...interface{}) string {
	return c.Warn(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Info(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorInfo).String()
}

func (c *ColorScheme) Infof(t string, args ...interface{}) string {
	return c.Info(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Error(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorWarn).String()
}

func (c *ColorScheme) Errorf(t string, args ...interface{}) string {
	return c.Error(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) SuccessIcon() string {
	return termenv.String("✓").Foreground(c.theme.ColorSuccess).String()
}

func (c *ColorScheme) WarningIcon() string {
	return termenv.String("!").Foreground(c.theme.ColorWarn).String()
}

func (c *ColorScheme) FailureIcon() string {
	return termenv.String("X").Foreground(c.theme.ColorError).String()
}
