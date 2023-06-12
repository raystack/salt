# config

This package is an wrapper over viper with opinionated defaults that allows loading config from a yaml file and environment variables.

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/raystack/salt/config"
)

type Config struct {
	Port     int            `mapstructure:"port" default:"8080"`
	DB       DBConfig       `mapstructure:"db"`
	NewRelic NewRelicConfig `mapstructure:"new_relic"`
	LogLevel string         `mapstructure:"log_level" default:"info"`
}

type DBConfig struct {
	Port int    `mapstructure:"port" default:"5432"`
	Host string `mapstructure:"host" default:"localhost"`
}

type NewRelicConfig struct {
	Enabled bool   `mapstructure:"enabled" default:"false"`
	AppName string `mapstructure:"app_name" default:"test-app"`
	License string `mapstructure:"license"`
}

func main() {
	var c Config
	l := config.NewLoader(
		// config.WithViper(viper.New()), // default
		// config.WithName("config"), // default
		// config.WithType("yaml"), // default
		// config.WithEnvKeyReplacer(".", "_"), // default
		config.WithPath("$HOME/.test"),
		config.WithEnvPrefix("CONFIG"),
	)

	if err := l.Load(&c); err != nil { // pass pointer to the struct into which you want to load config
		panic(err)
	}
	s, _ := json.MarshalIndent(c, "", "  ") // spaces: 2 | tabs: 1 😛
	fmt.Println(string(s))
}
```

In the above program a YAML file or environment variables can be used to configure.

```yaml
port: 9000
db:
  port: 5432
  host: db-host-yaml
new_relic:
  enabled: true
  app_name: config-test-yaml
  license: ____LICENSE_STRING_OF_40_CHARACTERS_____
log_level: debug
```

or

```sh
export CONFIG_PORT=9001
export CONFIG_DB_PORT=5432
export CONFIG_DB_HOST=db-host-env
export CONFIG_NEW_RELIC_ENABLED=true
export CONFIG_NEW_RELIC_APP_NAME=config-test-env
export CONFIG_NEW_RELIC_LICENSE=____LICENSE_STRING_OF_40_CHARACTERS_____
export CONFIG_LOG_LEVEL=debug
```

or a mix of both.

**Configs set in environment will override the ones set as default and in yaml file.**

## TODO

- function to print/return config keys in yaml path and env format with defaults as helper
- add support for flags
