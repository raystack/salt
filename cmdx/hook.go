package cmdx

import "github.com/spf13/cobra"

// SetClientHook applies custom cobra config specific
// for client or cmd that have annotation `client:true`
func SetClientHook(rootCmd *cobra.Command, applyFunc func(cmd *cobra.Command)) {
	for _, c := range rootCmd.Commands() {
		if c.Annotations != nil && c.Annotations["client"] == "true" {
			applyFunc(c)
		}
		SetClientHook(rootCmd, applyFunc)
	}
}
