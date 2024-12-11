// Package cmdx extends the capabilities of the Cobra library to build advanced CLI tools.
// It provides features such as custom help, shell completion, reference documentation generation,
// help topics, and client-specific hooks.
//
// # Features
//
//  1. **Custom Help**:
//     Enhance the default help output with a structured and detailed format.
//
//  2. **Reference Command**:
//     Generate markdown documentation for the entire CLI command tree.
//
//  3. **Shell Completion**:
//     Generate shell completion scripts for Bash, Zsh, Fish, and PowerShell.
//
//  4. **Help Topics**:
//     Add custom help topics to provide detailed information about specific subjects.
//
//  5. **Client Hooks**:
//     Apply custom logic to commands annotated with `client:true`.
//
// # Example
//
// The following example demonstrates how to use the cmdx package:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/spf13/cobra"
//	    "github.com/your-username/cmdx"
//	)
//
//	func main() {
//	    rootCmd := &cobra.Command{
//	        Use:   "mycli",
//	        Short: "A sample CLI tool",
//	    }
//
//	    // Define help topics
//	    helpTopics := []cmdx.HelpTopic{
//	        {
//	            Name:    "env",
//	            Short:   "Environment variables help",
//	            Long:    "Details about environment variables used by the CLI.",
//	            Example: "$ mycli help env",
//	        },
//	    }
//
//	    // Define hooks
//	    hooks := []cmdx.HookBehavior{
//	        {
//	            Name: "setup",
//	            Behavior: func(cmd *cobra.Command) {
//	                cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
//	                    fmt.Println("Setting up for", cmd.Name())
//	                }
//	            },
//	        },
//	    }
//
//	    // Create the Commander with configurations
//	    manager := cmdx.NewCommander(
//	        rootCmd,
//	        cmdx.WithTopics(helpTopics),
//	        cmdx.WithHooks(hooks),
//	        cmdx.EnableConfig(),
//	        cmdx.EnableDocs(),
//	    )
//
//	    // Initialize the manager
//	    if err := manager.Initialize(); err != nil {
//	        fmt.Println("Error initializing CLI:", err)
//	        return
//	    }
//
//	    // Execute the CLI
//	    if err := rootCmd.Execute(); err != nil {
//	        fmt.Println("Command execution failed:", err)
//	    }
//	}
package cmdx
