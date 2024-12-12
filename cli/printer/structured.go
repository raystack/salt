package printer

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAML prints the given data in YAML format.
func YAML(data interface{}) error {
	return File(data, "yaml")
}

// JSON prints the given data in JSON format.
func JSON(data interface{}) error {
	return File(data, "json")
}

// PrettyJSON prints the given data in pretty-printed JSON format.
func PrettyJSON(data interface{}) error {
	return File(data, "prettyjson")
}

// File marshals and prints the given data in the specified format.
func File(data interface{}, format string) (err error) {
	var output []byte
	switch format {
	case "yaml":
		output, err = yaml.Marshal(data)
	case "json":
		output, err = json.Marshal(data)
	case "prettyjson":
		output, err = json.MarshalIndent(data, "", "\t")
	default:
		return fmt.Errorf("unknown format: %v", format)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
