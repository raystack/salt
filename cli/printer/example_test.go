package printer_test

import (
	"os"

	"github.com/raystack/salt/cli/printer"
)

func ExampleNewOutput() {
	out := printer.NewOutput(os.Stdout)

	out.Success("deployed to prod")
	out.Warning("check logs for warnings")
	out.Error("connection failed")
	out.Info("3 items found")
	out.Bold("important message")
}

func ExampleOutput_Table() {
	out := printer.NewOutput(os.Stdout)

	rows := [][]string{
		{"ID", "NAME", "STATUS"},
		{"1", "Alice", "active"},
		{"2", "Bob", "inactive"},
	}
	out.Table(rows)
}

func ExampleOutput_JSON() {
	out := printer.NewOutput(os.Stdout)

	data := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}
	out.JSON(data)
}

func ExampleOutput_Spin() {
	out := printer.NewOutput(os.Stdout)

	spinner := out.Spin("loading...")
	// ... do work ...
	spinner.Stop()
}

func Example_colorFormatting() {
	// Color functions return styled strings for composition.
	status := printer.Green("passing") + " — " + printer.Red("2 failing")
	_ = status

	// Formatted variants work like fmt.Sprintf.
	count := printer.Greenf("found %d items", 42)
	_ = count

	// Icons for status indicators.
	ok := printer.Icon("success") + " all tests passed"
	_ = ok
}
