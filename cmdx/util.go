package cmdx

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IsCmdErr returns true if erorr is cobra command error.
// This is useful for distinguishing between a human error
// and a program error and displaying the correct message.
func IsCmdErr(err error) bool {
	errstr := err.Error()

	strs := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
	}

	for _, str := range strs {
		if strings.Contains(errstr, str) {
			return true
		}
	}
	return false
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

var lineRE = regexp.MustCompile(`(?m)^`)

func indent(s, indent string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}

func dirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func toTitle(text string) string {
	heading := cases.Title(language.Und).String(text)
	return heading
}

func bold(text string) termenv.Style {
	h := termenv.String(text).Bold()
	return h
}
