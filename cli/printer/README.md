# Printer

The `printer` package provides utilities for terminal-based output formatting, including colorized text, progress indicators, markdown rendering, and more. It is designed for building rich and user-friendly CLI applications.

## Features

- **Text Formatting**: Bold, italic, and colorized text.
- **Progress Indicators**: Spinners and progress bars for long-running tasks.
- **Markdown Rendering**: Render Markdown with terminal-friendly styles.
- **Structured Output**: YAML and JSON rendering with support for pretty-printing.
- **Icons**: Visual indicators like success and failure icons.

## Installation

Install the package using:

```bash
go get github.com/raystack/salt/cli/printer
```

## Usage

### Text Formatting

#### Basic Colors
```go
package main

import (
	"fmt"
	"github.com/raystack/salt/cli/printer"
)

func main() {
	fmt.Println(printer.Green("Success!"))
	fmt.Println(printer.Red("Error!"))
	fmt.Println(printer.Cyanf("Hello, %s!", "World"))
}
```

Supported Colors:
- **Green**: `printer.Green`, `printer.Greenf`
- **Red**: `printer.Red`, `printer.Redf`
- **Yellow**: `printer.Yellow`, `printer.Yellowf`
- **Cyan**: `printer.Cyan`, `printer.Cyanf`
- **Grey**: `printer.Grey`, `printer.Greyf`
- **Blue**: `printer.Blue`, `printer.Bluef`
- **Magenta**: `printer.Magenta`, `printer.Magentaf`

#### Bold and Italic Text
```go
fmt.Println(printer.Bold("This is bold text."))
fmt.Println(printer.Italic("This is italic text."))
```

### Progress Indicators

#### Spinner
```go
package main

import (
	"time"
	"github.com/yourusername/printer"
)

func main() {
	indicator := printer.Spin("Processing")
	time.Sleep(2 * time.Second) // Simulate work
	indicator.Stop()
}
```

#### Progress Bar
```go
package main

import (
	"time"
	"github.com/yourusername/printer"
)

func main() {
	bar := printer.Progress(100, "Downloading")
	for i := 0; i <= 100; i++ {
		time.Sleep(50 * time.Millisecond) // Simulate work
		bar.Add(1)
	}
}
```

### Markdown Rendering

#### Render Markdown
```go
package main

import (
	"fmt"
	"github.com/yourusername/printer"
)

func main() {
	output, err := printer.Markdown("# Hello, Markdown!")
	if err != nil {
		fmt.Println("Error rendering markdown:", err)
		return
	}
	fmt.Println(output)
}
```

#### Render Markdown with Word Wrap
```go
output, err := printer.MarkdownWithWrap("# Hello, Markdown!", 80)
if err != nil {
	fmt.Println("Error rendering markdown:", err)
	return
}
fmt.Println(output)
```

### File Rendering

#### YAML
```go
package main

import (
	"github.com/yourusername/printer"
)

func main() {
	data := map[string]string{"name": "John", "age": "30"}
	err := printer.YAML(data)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

#### JSON
```go
package main

import (
	"github.com/yourusername/printer"
)

func main() {
	data := map[string]string{"name": "John", "age": "30"}
	err := printer.JSON(data)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

#### Pretty JSON
```go
package main

import (
	"github.com/yourusername/printer"
)

func main() {
	data := map[string]string{"name": "John", "age": "30"}
	err := printer.PrettyJSON(data)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### Icons

#### Visual Indicators
```go
fmt.Println(printer.FailureIcon(), printer.Red("Task failed."))
```

## Themes

The package automatically detects the terminal’s background and switches between light and dark themes. Supported colors can be customized by modifying the `Theme` struct.