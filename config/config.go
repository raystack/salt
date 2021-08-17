package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Loader struct {
	v *viper.Viper
}

type LoaderOption func(*Loader)

func WithViper(in *viper.Viper) LoaderOption {
	return func(l *Loader) {
		l.v = in
	}
}

func WithConfigName(in string) LoaderOption {
	return func(l *Loader) {
		l.v.SetConfigName(in)
	}
}

func AddConfigPath(in string) LoaderOption {
	return func(l *Loader) {
		l.v.AddConfigPath(in)
	}
}

func WithConfigType(in string) LoaderOption {
	return func(l *Loader) {
		l.v.SetConfigType(in)
	}
}

func WithEnvKeyReplacer(old string, new string) LoaderOption {
	return func(l *Loader) {
		l.v.SetEnvKeyReplacer(strings.NewReplacer(old, new))
	}
}

func NewLoader(options ...LoaderOption) *Loader {
	loader := &Loader{
		v: getViperWithDefaults(),
	}

	for _, option := range options {
		option(loader)
	}
	return loader
}

// Load loads configuration into the given mapstructure (https://github.com/mitchellh/mapstructure)
// from a config.yaml file and overrides with any values set in env variables
func (l *Loader) Load(config interface{}) {
	verifyParamIsPtrToStructElsePanic(config)

	l.v.AutomaticEnv()

	err := l.v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panicf("unable to read configs using viper: %v\n", err)
		}
	}

	err, configKeys := getFlattenedStructKeys(config)
	if err != nil {
		panicf("unable to get all config keys from struct: %v\n", err)
	}

	// Bind each conf fields to environment vars
	for key := range configKeys {
		err := l.v.BindEnv(configKeys[key])
		if err != nil {
			panicf("unable to bind env keys: %v\n", err)
		}
	}

	// set defaults using the default struct tag
	defaults.SetDefaults(config)

	err = l.v.Unmarshal(config)
	if err != nil {
		panicf("unable to load config to struct: %v\n", err)
	}
}

func verifyParamIsPtrToStructElsePanic(param interface{}) {
	value := reflect.ValueOf(param)
	if value.Kind() != reflect.Ptr {
		panicf("Require Ptr to a Struct for Load. Got %v\n", value.Kind())
	} else {
		value = reflect.Indirect(value)
		if value.Kind() != reflect.Struct {
			panicf("Require Ptr to a Struct for Load. Got Ptr to %v\n", value.Kind())
		}
	}
}

func getViperWithDefaults() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}

func getFlattenedStructKeys(config interface{}) (error, []string) {
	var structMap map[string]interface{}
	err := mapstructure.Decode(config, &structMap)
	if err != nil {
		return err, nil
	}

	flat, err := flatten.Flatten(structMap, "", flatten.DotStyle)
	if err != nil {
		return err, nil
	}

	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}

	return nil, keys
}

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
