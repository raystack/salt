package config_test

import (
	"fmt"
	"log"

	"github.com/raystack/salt/config"
)

func ExampleNewLoader() {
	type Config struct {
		Server struct {
			Port int    `mapstructure:"port" default:"8080"`
			Host string `mapstructure:"host" default:"localhost"`
		} `mapstructure:"server"`
		LogLevel string `mapstructure:"log_level" default:"info"`
	}

	var cfg Config
	loader := config.NewLoader(
		config.WithFile("./config.yaml"),
		config.WithEnvPrefix("MYAPP"),
	)

	if err := loader.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("server: %s:%d, log: %s\n", cfg.Server.Host, cfg.Server.Port, cfg.LogLevel)
}
