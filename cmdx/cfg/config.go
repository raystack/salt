package cfg

import (
	"errors"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/salt/config"
	"gopkg.in/yaml.v3"
)

// Config represents a config file manager.
type Config struct {
	filename  string
	configMap interface{}
}

// Init initializes a client config file from config map.
func (c *Config) Init() error {
	defaults.SetDefaults(c.configMap)
	data, err := yaml.Marshal(c.configMap)
	if err != nil {
		return err
	}
	if fileExist(c.filename) {
		return errors.New("config file already exists")
	}
	return writeFile(c.filename, data)
}

// Read reads all configurations for client config file.
func (c *Config) Read() (string, error) {
	cfg, err := readFile(c.filename)
	return string(cfg), err
}

// File returns the config file name.
func (c *Config) File() string {
	return c.filename
}

// Load loads a config file into a config map.
func (c *Config) Load() (interface{}, error) {
	loader := config.NewLoader(config.WithFile(c.filename))

	if err := loader.Load(c.configMap); err != nil {
		return c.configMap, err
	}
	return c.configMap, nil
}

// NewCfg creates a new config file manager.
func NewCfg(app string, configMap interface{}) *Config {
	return &Config{
		filename:  ConfigFile(app),
		configMap: configMap,
	}
}
