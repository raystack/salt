package config

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// Commands returns a list of Cobra commands for managing the configuration.
func Commands(app string, cfgTemplate interface{}) (*cobra.Command, error) {
	cfg, err := New(app)
	if err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration",
		Annotations: map[string]string{
			"group": "core",
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "init",
			Short: "Initialize configuration with default values",
			Annotations: map[string]string{
				"group": "core",
			},
			Run: func(cmd *cobra.Command, args []string) {
				if err := cfg.Init(cfgTemplate); err != nil {
					log.Fatalf("Error initializing config: %v", err)
				}
				fmt.Println("Configuration initialized successfully.")
			},
		},
		&cobra.Command{
			Use:   "view",
			Short: "View the current configuration",
			Annotations: map[string]string{
				"group": "core",
			},
			Run: func(cmd *cobra.Command, args []string) {
				content, err := cfg.Read()
				if err != nil {
					log.Fatalf("Error reading config: %v", err)
				}
				fmt.Println("Current Configuration:")
				fmt.Println(content)
			},
		},
	)

	return cmd, nil
}
