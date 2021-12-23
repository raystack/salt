package cfg

import (
	"errors"
	"io/ioutil"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/salt/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	filename  string
	configMap interface{}
}

func (c *Config) Init() error {
	defaults.SetDefaults(c.configMap)

	data, err := yaml.Marshal(c.configMap)
	if err != nil {
		return err
	}

	if fileExist(c.filename) {
		return errors.New("config file already exists")
	}

	if err := ioutil.WriteFile(c.filename, data, 0655); err != nil {
		return err
	}
	return nil
}

func (c *Config) Read() (string, error) {
	cfg, err := ioutil.ReadFile(c.filename)
	return string(cfg), err
}

func (c *Config) File() string {
	return c.filename
}

func (c *Config) Load() (interface{}, error) {
	loader := config.NewLoader(config.WithFile(c.filename))

	if err := loader.Load(c.configMap); err != nil {
		return c.configMap, err
	}
	return c.configMap, nil
}

func NewCfg(app string, configMap interface{}) *Config {
	return &Config{
		filename:  ConfigFile(app),
		configMap: configMap,
	}
}
