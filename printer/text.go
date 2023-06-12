package printer

import (
	"fmt"

	"github.com/raystack/salt/term"
)

func Success(t ...string) {
	fmt.Print(term.Green(t...))
}

func Successln(t ...string) {
	fmt.Println(term.Green(t...))
}

func Successf(t string, args ...interface{}) {
	fmt.Print(term.Greenf(t, args...))
}

func Warning(t ...string) {
	fmt.Print(term.Yellow(t...))
}

func Warningln(t ...string) {
	fmt.Println(term.Yellow(t...))
}

func Warningf(t string, args ...interface{}) {
	fmt.Print(term.Yellowf(t, args...))
}

func Error(t ...string) {
	fmt.Print(term.Red(t...))
}

func Errorln(t ...string) {
	fmt.Println(term.Red(t...))
}

func Errorf(t string, args ...interface{}) {
	fmt.Print(term.Redf(t, args...))
}

func Info(t ...string) {
	fmt.Print(term.Cyan(t...))
}

func Infoln(t ...string) {
	fmt.Println(term.Cyan(t...))
}

func Infof(t string, args ...interface{}) {
	fmt.Print(term.Cyanf(t, args...))
}

func Bold(t ...string) {
	fmt.Print(term.Bold(t...))
}

func Boldln(t ...string) {
	fmt.Println(term.Bold(t...))
}

func Boldf(t string, args ...interface{}) {
	fmt.Print(term.Boldf(t, args...))
}

func Italic(t ...string) {
	fmt.Print(term.Italic(t...))
}

func Italicln(t ...string) {
	fmt.Println(term.Italic(t...))
}

func Italicf(t string, args ...interface{}) {
	fmt.Print(term.Italicf(t, args...))
}

func Text(t ...string) {
	fmt.Print(term.Grey(t...))
}

func Textln(t ...string) {
	fmt.Println(term.Grey(t...))
}

func Textf(t string, args ...interface{}) {
	fmt.Print(term.Greyf(t, args...))
}

func SuccessIcon() {
	fmt.Print(term.Green("✓"))
}

func WarningIcon() {
	fmt.Print(term.Yellow("!"))
}

func ErrorIcon() {
	fmt.Print(term.Red("✗"))
}

func InfoIcon() {
	fmt.Print(term.Cyan("⛭"))
}

func Space() {
	fmt.Print(" ")
}
