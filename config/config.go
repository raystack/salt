package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigFileNotFoundError indicates that the configuration file was not found.
// In this case, Viper will attempt to load configurations from environment variables or defaults.
type ConfigFileNotFoundError struct {
	Err error
}

// Error returns the error message for ConfigFileNotFoundError.
func (e ConfigFileNotFoundError) Error() string {
	return fmt.Sprintf("config file not found, falling back to environment and defaults: %v", e.Err)
}

// Unwrap provides compatibility for error unwrapping.
func (e *ConfigFileNotFoundError) Unwrap() error {
	return e.Err
}

// Loader is responsible for managing configuration loading and decoding.
type Loader struct {
	viperInstance *viper.Viper
	decoderOpts   []viper.DecoderConfigOption
}

// LoaderOption defines a functional option for configuring a Loader instance.
type LoaderOption func(*Loader)

// WithViper allows using a custom Viper instance.
func WithViper(v *viper.Viper) LoaderOption {
	return func(l *Loader) {
		l.viperInstance = v
	}
}

// WithFile specifies an explicit configuration file path.
func WithFile(file string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.SetConfigFile(file)
	}
}

// WithName sets the base name of the configuration file (excluding extension).
func WithName(name string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.SetConfigName(name)
	}
}

// WithPath adds a directory to search for the configuration file.
// Can be called multiple times to add multiple paths.
func WithPath(path string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.AddConfigPath(path)
	}
}

// WithType specifies the configuration file format (e.g., "yaml", "json").
func WithType(fileType string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.SetConfigType(fileType)
	}
}

// WithBindPFlags binds command-line flags to the configuration based on struct tags (`cmdx`).
func WithBindPFlags(flagSet *pflag.FlagSet, config interface{}) LoaderOption {
	return func(l *Loader) {
		structType := reflect.TypeOf(config).Elem()
		for i := 0; i < structType.NumField(); i++ {
			if tag := structType.Field(i).Tag.Get("cmdx"); tag != "" {
				l.viperInstance.BindPFlag(tag, flagSet.Lookup(tag))
			}
		}
	}
}

// WithEnvPrefix sets a prefix for environment variable keys.
func WithEnvPrefix(prefix string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.SetEnvPrefix(prefix)
	}
}

// WithEnvKeyReplacer customizes key transformation for environment variables.
func WithEnvKeyReplacer(old, new string) LoaderOption {
	return func(l *Loader) {
		l.viperInstance.SetEnvKeyReplacer(strings.NewReplacer(old, new))
	}
}

// WithDecoderConfigOption sets custom decoding options for the configuration loader.
func WithDecoderConfigOption(opts ...viper.DecoderConfigOption) LoaderOption {
	return func(l *Loader) {
		l.decoderOpts = append(l.decoderOpts, opts...)
	}
}

// NewLoader initializes a Loader instance with the specified options.
func NewLoader(options ...LoaderOption) *Loader {
	loader := &Loader{
		viperInstance: defaultViperInstance(),
		decoderOpts: []viper.DecoderConfigOption{
			viper.DecodeHook(
				mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					StringToJsonFunc(),
				),
			),
		},
	}
	for _, opt := range options {
		opt(loader)
	}
	return loader
}

// Load populates the provided config struct with values from the configuration sources.
func (l *Loader) Load(config interface{}) error {
	if err := validateStructPtr(config); err != nil {
		return err
	}

	l.viperInstance.AutomaticEnv()

	if err := l.viperInstance.ReadInConfig(); err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) || errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return ConfigFileNotFoundError{Err: err}
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	keys, err := extractFlattenedKeys(config)
	if err != nil {
		return fmt.Errorf("failed to extract config keys: %w", err)
	}

	for _, key := range keys {
		l.viperInstance.BindEnv(key)
	}

	defaults.SetDefaults(config)

	if err := l.viperInstance.Unmarshal(config, l.decoderOpts...); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// validateStructPtr ensures the provided value is a pointer to a struct.
func validateStructPtr(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("Load requires a pointer to a struct")
	}
	return nil
}

// defaultViperInstance initializes a Viper instance with default settings.
func defaultViperInstance() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}

// extractFlattenedKeys retrieves all keys from the struct in a flattened format.
func extractFlattenedKeys(config interface{}) ([]string, error) {
	var structMap map[string]interface{}
	if err := mapstructure.Decode(config, &structMap); err != nil {
		return nil, err
	}
	flatMap, err := flatten.Flatten(structMap, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(flatMap))
	for k := range flatMap {
		keys = append(keys, k)
	}
	return keys, nil
}

// StringToJsonFunc is a decode hook for parsing JSON strings into maps.
func StringToJsonFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t.Kind() == reflect.Map {
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(data.(string)), &result); err == nil {
				return result, nil
			}
		}
		return data, nil
	}
}
