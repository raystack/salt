package cmdx

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/salt/config"
	"gopkg.in/yaml.v3"
)

const (
	ODPF_CONFIG_DIR = "ODPF_CONFIG_DIR"
	XDG_CONFIG_HOME = "XDG_CONFIG_HOME"
	APP_DATA        = "AppData"
	LOCAL_APP_DATA  = "LocalAppData"
)

// SetConfig allows to set a client config file.
// It is used to load and save a config file
// for command line clients.
func SetConfig(app string) *Config {
	return &Config{
		filename: configFile(app),
	}
}

type Config struct {
	filename string
}

func (c *Config) File() string {
	return c.filename
}

func (c *Config) Defaults(cfg interface{}) {
	defaults.SetDefaults(cfg)
}

func (c *Config) Init(cfg interface{}) error {
	defaults.SetDefaults(cfg)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if fileExist(c.filename) {
		return errors.New("config file already exists")
	}

	if _, err := os.Stat(c.filename); os.IsNotExist(err) {
		os.MkdirAll(configDir("odpf"), 0700)
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

func (c *Config) Load(cfg interface{}) error {
	loader := config.NewLoader(config.WithFile(c.filename))

	if err := loader.Load(cfg); err != nil {
		return err
	}
	return nil
}

func configFile(app string) string {
	file := app + ".yml"
	return filepath.Join(configDir("odpf"), file)
}

func configDir(root string) string {
	var path string
	if a := os.Getenv(ODPF_CONFIG_DIR); a != "" {
		path = a
	} else if b := os.Getenv(XDG_CONFIG_HOME); b != "" {
		path = filepath.Join(b, root)
	} else if c := os.Getenv(APP_DATA); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, root)
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", root)
	}

	if !dirExists(path) {
		_ = os.MkdirAll(filepath.Dir(path), 0755)
	}

	return path
}

func dirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
