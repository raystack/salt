package cli

import (
	"fmt"

	"github.com/raystack/salt/config"
	"github.com/spf13/cobra"
)

// ConfigCommand returns a "config" command with "init" and "list" subcommands
// for managing client-side CLI configuration.
//
// The appName is used to determine the config file location
// (~/.config/raystack/<appName>.yml). The defaultCfg is a pointer to a struct
// with default values used when initializing a new config file.
//
// Usage:
//
//	rootCmd.AddCommand(cli.ConfigCommand("frontier", &Config{}))
func ConfigCommand(appName string, defaultCfg interface{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config <command>",
		Short:   "Manage client configuration",
		Example: fmt.Sprintf("  $ %s config init\n  $ %s config list", appName, appName),
	}

	cmd.AddCommand(configInitCmd(appName, defaultCfg))
	cmd.AddCommand(configListCmd(appName))

	return cmd
}

func configInitCmd(appName string, defaultCfg interface{}) *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new configuration file",
		Example: fmt.Sprintf("  $ %s config init", appName),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			loader := config.NewLoader(config.WithAppConfig(appName))
			if err := loader.Init(defaultCfg); err != nil {
				return err
			}
			Output(cmd).Success("config initialized")
			return nil
		},
	}
}

func configListCmd(appName string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List current configuration",
		Example: fmt.Sprintf("  $ %s config list", appName),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			loader := config.NewLoader(config.WithAppConfig(appName))
			data, err := loader.View()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			Output(cmd).Println(data)
			return nil
		},
	}
}
