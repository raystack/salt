package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/jeremywohl/flatten"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// bindFlags dynamically binds flags to configuration fields based on `cmdx` tags.
func bindFlags(v *viper.Viper, flagSet *pflag.FlagSet, structType reflect.Type, parentKey string) error {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("cmdx")
		if tag == "" {
			continue
		}

		if parentKey != "" {
			tag = parentKey + "." + tag
		}

		if field.Type.Kind() == reflect.Struct {
			// Recurse into nested structs
			if err := bindFlags(v, flagSet, field.Type, tag); err != nil {
				return err
			}
		} else {
			flag := flagSet.Lookup(tag)
			if flag == nil {
				return fmt.Errorf("missing flag for tag: %s", tag)
			}
			if err := v.BindPFlag(tag, flag); err != nil {
				return fmt.Errorf("failed to bind flag for tag: %s, error: %w", tag, err)
			}
		}
	}
	return nil
}

// validateStructPtr ensures the provided value is a pointer to a struct.
func validateStructPtr(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("load requires a pointer to a struct")
	}
	return nil
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

// Utilities for app-specific configuration paths
func getConfigFilePath(app string) (string, error) {
	dirPath := getConfigDir("raystack")
	if err := ensureDir(dirPath); err != nil {
		return "", err
	}
	return filepath.Join(dirPath, app+".yml"), nil
}

func getConfigDir(root string) string {
	switch {
	case envSet("RAYSTACK_CONFIG_DIR"):
		return filepath.Join(os.Getenv("RAYSTACK_CONFIG_DIR"), root)
	case envSet("XDG_CONFIG_HOME"):
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), root)
	case runtime.GOOS == "windows" && envSet("APPDATA"):
		return filepath.Join(os.Getenv("APPDATA"), root)
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", root)
	}
}

func ensureDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func envSet(key string) bool {
	return os.Getenv(key) != ""
}
