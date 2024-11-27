package cmdx

import "github.com/spf13/cobra"

// SetClientHook recursively applies a custom function to all commands
// with the annotation `client:true` in the given Cobra command tree.
//
// This is particularly useful for applying client-specific configurations
// to commands annotated as "client".
//
// Parameters:
//   - rootCmd: The root Cobra command to start traversing from.
//   - applyFunc: A function that applies the desired configuration
//     to commands with the `client:true` annotation.
//
// Example Usage:
//
//	cmdx.SetClientHook(rootCmd, func(cmd *cobra.Command) {
//	    cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
//	        fmt.Println("Client-specific setup")
//	    }
//	})
func SetClientHook(rootCmd *cobra.Command, applyFunc func(cmd *cobra.Command)) {
	for _, subCmd := range rootCmd.Commands() {
		if subCmd.Annotations != nil && subCmd.Annotations["client"] == "true" {
			applyFunc(subCmd)
		}
		SetClientHook(subCmd, applyFunc)
	}
}
