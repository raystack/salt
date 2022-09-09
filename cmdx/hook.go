package cmdx

import "github.com/spf13/cobra"

// SetClientHook applies custom cobra config specific
// for client or cmd that have annotation `client:true`
func SetClientHook(rootCmd *cobra.Command, applyFunc func(cmd *cobra.Command)) {
	for _, subCmd := range rootCmd.Commands() {
		if subCmd.Annotations != nil && subCmd.Annotations["client"] == "true" {
			applyFunc(subCmd)
		}
		SetClientHook(subCmd, applyFunc)
	}
}
