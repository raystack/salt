package printer

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAML prints the given data in YAML format.
//
// Parameters:
//   - data: The data to be marshaled into YAML and printed.
//
// Returns:
//   - An error if the data cannot be marshaled into YAML.
//
// Example Usage:
//
//	config := map[string]string{"key": "value"}
//	err := printer.YAML(config)
func YAML(data interface{}) error {
	return File(data, "yaml")
}

// JSON prints the given data in JSON format.
//
// Parameters:
//   - data: The data to be marshaled into JSON and printed.
//
// Returns:
//   - An error if the data cannot be marshaled into JSON.
//
// Example Usage:
//
//	config := map[string]string{"key": "value"}
//	err := printer.JSON(config)
func JSON(data interface{}) error {
	return File(data, "json")
}

// PrettyJSON prints the given data in pretty-printed JSON format.
//
// Parameters:
//   - data: The data to be marshaled into indented JSON and printed.
//
// Returns:
//   - An error if the data cannot be marshaled into JSON.
//
// Example Usage:
//
//	config := map[string]string{"key": "value"}
//	err := printer.PrettyJSON(config)
func PrettyJSON(data interface{}) error {
	return File(data, "prettyjson")
}

// File marshals and prints the given data in the specified format.
//
// Supported formats:
//   - "yaml": Prints the data as YAML.
//   - "json": Prints the data as compact JSON.
//   - "prettyjson": Prints the data as pretty-printed JSON.
//
// Parameters:
//   - data: The data to be marshaled and printed.
//   - format: The desired output format ("yaml", "json", or "prettyjson").
//
// Returns:
//   - An error if the data cannot be marshaled into the specified format or if the format is unsupported.
//
// Example Usage:
//
//	config := map[string]string{"key": "value"}
//	err := printer.File(config, "yaml")
func File(data interface{}, format string) (err error) {
	var output []byte
	switch format {
	case "yaml":
		// Marshal the data into YAML format.
		output, err = yaml.Marshal(data)
	case "json":
		// Marshal the data into compact JSON format.
		output, err = json.Marshal(data)
	case "prettyjson":
		// Marshal the data into pretty-printed JSON format.
		output, err = json.MarshalIndent(data, "", "\t")
	default:
		// Return an error for unsupported formats.
		return fmt.Errorf("unknown format: %v", format)
	}

	if err != nil {
		return err
	}

	// Print the marshaled data to stdout.
	fmt.Println(string(output))
	return nil
}
