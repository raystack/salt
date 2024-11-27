# Prompt Package

The `prompt` package simplifies interactive CLI input using the `survey` library. It provides a consistent interface for user prompts such as single and multi-selection, text input, and confirmation.

## Features

- **Select**: Prompt users to select one option from a list.
- **MultiSelect**: Prompt users to select multiple options.
- **Input**: Prompt users to provide text input.
- **Confirm**: Prompt users for a yes/no confirmation.

## Installation

Add the package to your Go project:

```bash
go get github.com/raystack/salt/cli/prompt
```

## Usage

Here’s an example of how to use the `prompt` package:

```go
package main

import (
    "fmt"
    "github.com/raystack/salt/cli/prompt"
)

func main() {
    p := prompt.New()

    // Single selection
    index, err := p.Select("Choose an option:", "Option 1", []string{"Option 1", "Option 2", "Option 3"})
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Selected option index:", index)

    // Multi-selection
    indices, err := p.MultiSelect("Choose multiple options:", nil, []string{"Option A", "Option B", "Option C"})
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Selected option indices:", indices)

    // Text input
    input, err := p.Input("Enter your name:", "John Doe")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Input:", input)

    // Confirmation
    confirm, err := p.Confirm("Do you want to proceed?", true)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Confirmation:", confirm)
}
```