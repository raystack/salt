package printer

import (
	"fmt"
)

func Success(t ...string) {
	fmt.Print(Green(t...))
}

func Successln(t ...string) {
	fmt.Println(Green(t...))
}

func Successf(t string, args ...interface{}) {
	fmt.Print(Greenf(t, args...))
}

func Warning(t ...string) {
	fmt.Print(Yellow(t...))
}

func Warningln(t ...string) {
	fmt.Println(Yellow(t...))
}

func Warningf(t string, args ...interface{}) {
	fmt.Print(Yellowf(t, args...))
}

func Error(t ...string) {
	fmt.Print(Red(t...))
}

func Errorln(t ...string) {
	fmt.Println(Red(t...))
}

func Errorf(t string, args ...interface{}) {
	fmt.Print(Redf(t, args...))
}

func Info(t ...string) {
	fmt.Print(Cyan(t...))
}

func Infoln(t ...string) {
	fmt.Println(Cyan(t...))
}

func Infof(t string, args ...interface{}) {
	fmt.Print(Cyanf(t, args...))
}

func Bold(t ...string) {
	fmt.Print(bold(t...))
}

func Boldln(t ...string) {
	fmt.Println(bold(t...))
}

func Boldf(t string, args ...interface{}) {
	fmt.Print(boldf(t, args...))
}

func Italic(t ...string) {
	fmt.Print(italic(t...))
}

func Italicln(t ...string) {
	fmt.Println(italic(t...))
}

func Italicf(t string, args ...interface{}) {
	fmt.Print(italicf(t, args...))
}

func Text(t ...string) {
	fmt.Print(Grey(t...))
}

func Textln(t ...string) {
	fmt.Println(Grey(t...))
}

func Textf(t string, args ...interface{}) {
	fmt.Print(Greyf(t, args...))
}

func SuccessIcon() {
	fmt.Print(Green("✓"))
}

func WarningIcon() {
	fmt.Print(Yellow("!"))
}

func ErrorIcon() {
	fmt.Print(Red("✗"))
}

func InfoIcon() {
	fmt.Print(Cyan("⛭"))
}

func Space() {
	fmt.Print(" ")
}
