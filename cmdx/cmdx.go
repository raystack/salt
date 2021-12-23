package cmdx

import (
	"github.com/odpf/salt/cmdx/cfg"
	"github.com/odpf/salt/cmdx/help"
	"github.com/spf13/cobra"
)

// SetHelp sets a custom help and usage function.
// It allows to group commands in different sections
// based on cobra commands annotations.
func SetHelp(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		help.RootHelpFunc(cmd, args)
	})
	cmd.SetUsageFunc(help.RootUsageFunc)
	cmd.SetFlagErrorFunc(help.RootFlagErrorFunc)
}

// SetConfig allows to set a client config file.
// It is used to load and save a config file
// for command line clients.
func SetConfig(app string, configMap interface{}) *cfg.Config {
	return cfg.NewCfg(app, configMap)
}
