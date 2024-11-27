package printer

import (
	"fmt"
)

// Success prints the given message(s) in green to indicate success.
func Success(t ...string) {
	fmt.Print(Green(t...))
}

// Successln prints the given message(s) in green with a newline.
func Successln(t ...string) {
	fmt.Println(Green(t...))
}

// Successf formats and prints the success message in green.
func Successf(t string, args ...interface{}) {
	fmt.Print(Greenf(t, args...))
}

// Warning prints the given message(s) in yellow to indicate a warning.
func Warning(t ...string) {
	fmt.Print(Yellow(t...))
}

// Warningln prints the given message(s) in yellow with a newline.
func Warningln(t ...string) {
	fmt.Println(Yellow(t...))
}

// Warningf formats and prints the warning message in yellow.
func Warningf(t string, args ...interface{}) {
	fmt.Print(Yellowf(t, args...))
}

// Error prints the given message(s) in red to indicate an error.
func Error(t ...string) {
	fmt.Print(Red(t...))
}

// Errorln prints the given message(s) in red with a newline.
func Errorln(t ...string) {
	fmt.Println(Red(t...))
}

// Errorf formats and prints the error message in red.
func Errorf(t string, args ...interface{}) {
	fmt.Print(Redf(t, args...))
}

// Info prints the given message(s) in cyan to indicate informational messages.
func Info(t ...string) {
	fmt.Print(Cyan(t...))
}

// Infoln prints the given message(s) in cyan with a newline.
func Infoln(t ...string) {
	fmt.Println(Cyan(t...))
}

// Infof formats and prints the informational message in cyan.
func Infof(t string, args ...interface{}) {
	fmt.Print(Cyanf(t, args...))
}

// Bold prints the given message(s) in bold style.
func Bold(t ...string) {
	fmt.Print(bold(t...))
}

// Boldln prints the given message(s) in bold style with a newline.
func Boldln(t ...string) {
	fmt.Println(bold(t...))
}

// Boldf formats and prints the message in bold style.
func Boldf(t string, args ...interface{}) {
	fmt.Print(boldf(t, args...))
}

// Italic prints the given message(s) in italic style.
func Italic(t ...string) {
	fmt.Print(italic(t...))
}

// Italicln prints the given message(s) in italic style with a newline.
func Italicln(t ...string) {
	fmt.Println(italic(t...))
}

// Italicf formats and prints the message in italic style.
func Italicf(t string, args ...interface{}) {
	fmt.Print(italicf(t, args...))
}

// Space prints a single space to the output.
func Space() {
	fmt.Print(" ")
}
