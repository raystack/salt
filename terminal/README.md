# Terminal

The `terminal` package provides a collection of utilities to manage terminal interactions, including pager handling, TTY detection, and environment configuration for better command-line application support.

## Features

- **Pager Management**: Easily manage pagers like `less` or `more` to display output in a paginated format.
- **TTY Detection**: Check if the program is running in a terminal environment.
- **Color Management**: Determine if color output is disabled based on environment variables.
- **CI Detection**: Identify if the program is running in a Continuous Integration (CI) environment.
- **Homebrew Utilities**: Check for Homebrew installation and verify binary paths.
- **Browser Launching**: Open URLs in the default web browser, with cross-platform support.
- 
## Installation

To include this package in your Go project, use:

```bash
go get github.com/raystack/salt
```

## Usage

### 1. Creating and Using a Pager

The `Pager` struct manages a pager process for displaying output in a paginated format.

```go
package main

import (
    "fmt"
    "github.com/raystack/salt/terminal"
)

func main() {
    // Create a new Pager instance
    pager := terminal.NewPager()

    // Optionally, set a custom pager command
    pager.Set("less -R")

    // Start the pager
    err := pager.Start()
    if err != nil {
        fmt.Println("Error starting pager:", err)
        return
    }
    defer pager.Stop()

    // Output text to the pager
    fmt.Fprintln(pager.Out, "This is a sample text output to the pager.")
}
```

### 2. Checking if the Terminal is a TTY

Use `IsTTY` to check if the output is a TTY (teletypewriter).

```go
if terminal.IsTTY() {
    fmt.Println("Running in a terminal!")
} else {
    fmt.Println("Not running in a terminal.")
}
```

### 3. Checking if Color Output is Disabled

Use `IsColorDisabled` to determine if color output should be suppressed.

```go
if terminal.IsColorDisabled() {
    fmt.Println("Color output is disabled.")
} else {
    fmt.Println("Color output is enabled.")
}
```

### 4. Checking if Running in a CI Environment

Use `IsCI` to check if the program is running in a CI environment.

```go
if terminal.IsCI() {
    fmt.Println("Running in a Continuous Integration environment.")
} else {
    fmt.Println("Not running in a CI environment.")
}
```


### 4. Checking if Running in a CI Environment

Use `IsCI` to check if the program is running in a CI environment.

```go
if termutil.IsCI() {
    fmt.Println("Running in a Continuous Integration environment.")
} else {
    fmt.Println("Not running in a CI environment.")
}
```

### 5. Checking for Homebrew Installation

Use `HasHomebrew` to check if Homebrew is installed on the system.

```go
if termuinal.HasHomebrew() {
    fmt.Println("Homebrew is installed!")
} else {
    fmt.Println("Homebrew is not installed.")
}
```

### 6. Checking if a Binary is Under Homebrew Path

Use `IsUnderHomebrew` to determine if a binary is managed by Homebrew.

```go
binaryPath := "/usr/local/bin/somebinary"
if terminal.IsUnderHomebrew(binaryPath) {
    fmt.Println("The binary is under the Homebrew path.")
} else {
    fmt.Println("The binary is not under the Homebrew path.")
}
```

### 7. Opening a URL in the Default Web Browser

Use `OpenBrowser` to launch the default web browser with a specified URL.

```go
goos := "darwin" // Use runtime.GOOS to get the current OS in a real scenario
url := "https://www.example.com"
cmd := terminal.OpenBrowser(goos, url)
if err := cmd.Run(); err != nil {
    fmt.Println("Failed to open browser:", err)
}
```