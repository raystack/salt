# config

This package is an wrapper over viper wwith opinionated defaults that allows loading config from a yaml file and environment variables.

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/odpf/salt/config"
)

type Config struct {
	Port     int          `mapstructure:"port" default:"8080"`
	NewRelic NestedConfig `mapstructure:"newrelic"`
	LogLevel string       `mapstructure:"log_level" default:"info"`
}

type NestedConfig struct {
	Enabled bool   `mapstructure:"enabled" default:"false"`
	AppName string `mapstructure:"appname" default:"app"`
	License string `mapstructure:"license"`
}

func main() {
	var c Config
	config.LoadConfig(&c) // pass pointer to the struct into which you want to load config
	s, _ := json.MarshalIndent(c, "", "  ") // spaces: 2 | tabs: 1 😛
	fmt.Println(string(s))
}
```

In the above program a YAML file or environment variables can be used to configure.

```yaml
port: 9000
newrelic:
    enabled: true
    appname: config-test-yaml
    license: ____LICENSE_STRING_OF_40_CHARACTERS_____
log_level: debug
```

or

```sh
export PORT=9001
export NEWRELIC_ENABLED=true
export NEWRELIC_APPNAME=config-test-env
export NEWRELIC_LICENSE=____LICENSE_STRING_OF_40_CHARACTERS_____
export LOG_LEVEL=debug
```

or a mix of both. 

**Configs set in environment will override the ones set as default and in yaml file.**

## TODO
 - allow overiding viper configs and injecting viper instance
 - function to print/return config keys in yaml path and env format as helper with defaults
 - add support for flags
