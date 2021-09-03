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
	ColorNeutral termenv.Color
	ColorBase    termenv.Color
	ColorPrimary termenv.Color
}

var themes = map[string]Theme{
	"light": {
		ColorSuccess: tp.Color("#005F00"),
		ColorWarn:    tp.Color("#FFAF00"),
		ColorInfo:    tp.Color("#0087FF"),
		ColorError:   tp.Color("#D70000"),
		ColorNeutral: tp.Color("#303030"),
		ColorBase:    tp.Color("#000087"),
		ColorPrimary: tp.Color("#AF00FF"),
	},
	"dark": {
		ColorSuccess: tp.Color("#A8CC8C"),
		ColorWarn:    tp.Color("#DBAB79"),
		ColorInfo:    tp.Color("#66C2CD"),
		ColorError:   tp.Color("#E88388"),
		ColorNeutral: tp.Color("#B9BFCA"),
		ColorBase:    tp.Color("#71BEF2"),
		ColorPrimary: tp.Color("#D290E4"),
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
		theme: themes["dark"],
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
	return termenv.String(t).Foreground(c.theme.ColorError).String()
}

func (c *ColorScheme) Errorf(t string, args ...interface{}) string {
	return c.Error(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Neutral(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorNeutral).String()
}

func (c *ColorScheme) Neutralf(t string, args ...interface{}) string {
	return c.Neutral(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Base(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorBase).String()
}

func (c *ColorScheme) Basef(t string, args ...interface{}) string {
	return c.Base(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Primary(t string) string {
	return termenv.String(t).Foreground(c.theme.ColorPrimary).String()
}

func (c *ColorScheme) Primaryf(t string, args ...interface{}) string {
	return c.Primary(fmt.Sprintf(t, args...))
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
