package printer

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// Table renders a terminal-friendly table to the provided writer.
//
// Create a table with customized formatting and styles,
// suitable for displaying data in CLI applications.
//
// Parameters:
//   - target: The `io.Writer` where the table will be written (e.g., os.Stdout).
//   - rows: A 2D slice of strings representing the table rows and columns.
//     Each inner slice represents a single row, with its elements as column values.
//
// Example Usage:
//
//	rows := [][]string{
//	    {"ID", "Name", "Age"},
//	    {"1", "Alice", "30"},
//	    {"2", "Bob", "25"},
//	}
//	printer.Table(os.Stdout, rows)
//
// Behavior:
//   - Disables text wrapping for better terminal rendering.
//   - Aligns headers and rows to the left.
//   - Removes borders and separators for a clean look.
//   - Formats the table using tab padding for better alignment in terminals.
func Table(target io.Writer, rows [][]string) {
	table := tablewriter.NewWriter(target)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows)
	table.Render()
}
