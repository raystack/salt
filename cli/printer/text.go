package printer

import (
	"fmt"
	"github.com/muesli/termenv"
)

// Success prints the given message(s) in green to indicate success.
func Success(t ...string) {
	printWithColor(Green, t...)
}

// Successln prints the given message(s) in green with a newline.
func Successln(t ...string) {
	printWithColorln(Green, t...)
}

// Successf formats and prints the success message in green.
func Successf(t string, args ...interface{}) {
	printWithColorf(Greenf, t, args...)
}

// Warning prints the given message(s) in yellow to indicate a warning.
func Warning(t ...string) {
	printWithColor(Yellow, t...)
}

// Warningln prints the given message(s) in yellow with a newline.
func Warningln(t ...string) {
	printWithColorln(Yellow, t...)
}

// Warningf formats and prints the warning message in yellow.
func Warningf(t string, args ...interface{}) {
	printWithColorf(Yellowf, t, args...)
}

// Error prints the given message(s) in red to indicate an error.
func Error(t ...string) {
	printWithColor(Red, t...)
}

// Errorln prints the given message(s) in red with a newline.
func Errorln(t ...string) {
	printWithColorln(Red, t...)
}

// Errorf formats and prints the error message in red.
func Errorf(t string, args ...interface{}) {
	printWithColorf(Redf, t, args...)
}

// Info prints the given message(s) in cyan to indicate informational messages.
func Info(t ...string) {
	printWithColor(Cyan, t...)
}

// Infoln prints the given message(s) in cyan with a newline.
func Infoln(t ...string) {
	printWithColorln(Cyan, t...)
}

// Infof formats and prints the informational message in cyan.
func Infof(t string, args ...interface{}) {
	printWithColorf(Cyanf, t, args...)
}

// Bold prints the given message(s) in bold style.
func Bold(t ...string) string {
	return termenv.String(t...).Bold().String()
}

// Boldln prints the given message(s) in bold style with a newline.
func Boldln(t ...string) {
	fmt.Println(Bold(t...))
}

// Boldf formats and prints the message in bold style.
func Boldf(t string, args ...interface{}) string {
	return Bold(fmt.Sprintf(t, args...))
}

// Italic prints the given message(s) in italic style.
func Italic(t ...string) string {
	return termenv.String(t...).Italic().String()
}

// Italicln prints the given message(s) in italic style with a newline.
func Italicln(t ...string) {
	fmt.Println(Italic(t...))
}

// Italicf formats and prints the message in italic style.
func Italicf(t string, args ...interface{}) string {
	return Italic(fmt.Sprintf(t, args...))
}

// Space prints a single space to the output.
func Space() {
	fmt.Print(" ")
}

// printWithColor prints the given message(s) with the specified color function.
func printWithColor(colorFunc func(...string) string, t ...string) {
	fmt.Print(colorFunc(t...))
}

// printWithColorln prints the given message(s) with the specified color function and a newline.
func printWithColorln(colorFunc func(...string) string, t ...string) {
	fmt.Println(colorFunc(t...))
}

// printWithColorf formats and prints the message with the specified color function.
func printWithColorf(colorFunc func(string, ...interface{}) string, t string, args ...interface{}) {
	fmt.Print(colorFunc(t, args...))
}
