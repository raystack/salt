
# cmdx

`cmdx` is a utility package designed to enhance the functionality of [Cobra](https://github.com/spf13/cobra), a popular Go library for creating command-line interfaces. It provides various helper functions and features to streamline CLI development, such as custom help topics, shell completion, command annotations, and client-specific configurations.

## Features

- **Help Topics**: Add custom help topics with descriptions and examples.
- **Shell Completions**: Generate completion scripts for Bash, Zsh, Fish, and PowerShell.
- **Command Reference**: Generate markdown documentation for all commands.
- **Client Hooks**: Apply custom configurations to commands annotated with `client:true`.


## Installation

To install the `cmdx` package, add it to your project using `go get`:

```bash
go get github.com/raystack/salt/cli/cmdx
```

## Usages

### SetHelpTopicCmd

Provides a way to define custom help topics that appear in the `help` command.

#### Example Usage
```go
topic := map[string]string{
    "short":   "Environment variables help",
    "long":    "Detailed information about environment variables used by the CLI.",
    "example": "$ mycli help env",
}

rootCmd.AddCommand(cmdx.SetHelpTopicCmd("env", topic))
```

#### Output
```bash
$ mycli help env
Detailed information about environment variables used by the CLI.

EXAMPLES

  $ mycli help env
```

---

### SetCompletionCmd

Adds a `completion` command to generate shell completion scripts for Bash, Zsh, Fish, and PowerShell.

#### Example Usage
```go
completionCmd := cmdx.SetCompletionCmd("mycli")
rootCmd.AddCommand(completionCmd)
```

#### Command Output
```bash
# Generate Bash completion script
$ mycli completion bash

# Generate Zsh completion script
$ mycli completion zsh

# Generate Fish completion script
$ mycli completion fish
```

#### Supported Shells
- **Bash**: Use `mycli completion bash`.
- **Zsh**: Use `mycli completion zsh`.
- **Fish**: Use `mycli completion fish`.
- **PowerShell**: Use `mycli completion powershell`.

---

### SetRefCmd

Adds a `reference` command to generate markdown documentation for all commands in the CLI hierarchy.

#### Example Usage
```go
refCmd := cmdx.SetRefCmd(rootCmd)
rootCmd.AddCommand(refCmd)
```

#### Command Output
```bash
$ mycli reference
# mycli reference

## `example`

A sample subcommand for the CLI.

## `another`

Another example subcommand.
```

---

### SetClientHook

Applies a custom function to commands annotated with `client:true`. Useful for client-specific configurations.

#### Example Usage
```go
cmdx.SetClientHook(rootCmd, func(cmd *cobra.Command) {
    cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
        fmt.Println("Executing client-specific setup for", cmd.Name())
    }
})
```

#### Command Example
```go
clientCmd := &cobra.Command{
    Use:   "client-action",
    Short: "A client-specific action",
    Annotations: map[string]string{
        "client": "true",
    },
}
rootCmd.AddCommand(clientCmd)
```

#### Output
```bash
$ mycli client-action
Executing client-specific setup for client-action
```

---

## Examples

Adding all features to a CLI

```go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/raystack/salt/cli/cmdx"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mycli",
		Short: "A custom CLI tool",
	}

	// Add Help Topic
	topic := map[string]string{
		"short":   "Environment variables help",
		"long":    "Details about environment variables used by the CLI.",
		"example": "$ mycli help env",
	}
	rootCmd.AddCommand(cmdx.SetHelpTopicCmd("env", topic))

	// Add Completion Command
	rootCmd.AddCommand(cmdx.SetCompletionCmd("mycli"))

	// Add Reference Command
	rootCmd.AddCommand(cmdx.SetRefCmd(rootCmd))

	// Add Client Hook
	clientCmd := &cobra.Command{
		Use:   "client-action",
		Short: "A client-specific action",
		Annotations: map[string]string{
			"client": "true",
		},
	}
	rootCmd.AddCommand(clientCmd)

	cmdx.SetClientHook(rootCmd, func(cmd *cobra.Command) {
		cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			fmt.Println("Executing client-specific setup for", cmd.Name())
		}
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
```