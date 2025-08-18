package jsondiff

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JSONReconstructor struct{}

func NewJSONReconstructor() *JSONReconstructor {
	return &JSONReconstructor{}
}

func (jr *JSONReconstructor) ReverseDiff(currentJSON string, diffs []DiffEntry) (string, error) {
	var current interface{}
	if err := json.Unmarshal([]byte(currentJSON), &current); err != nil {
		return "", fmt.Errorf("error parsing current JSON: %w", err)
	}

	reconstructed := jr.deepCopy(current)

	for i := len(diffs) - 1; i >= 0; i-- {
		diff := diffs[i]
		if err := jr.applyReverseDiff(reconstructed, diff); err != nil {
			return "", fmt.Errorf("error applying reverse diff at %s: %w", diff.FullPath, err)
		}
	}

	result, err := json.MarshalIndent(reconstructed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling result: %w", err)
	}

	return string(result), nil
}

func (jr *JSONReconstructor) applyReverseDiff(obj interface{}, diff DiffEntry) error {
	pathParts := jr.parsePath(diff.FullPath)

	switch diff.ChangeType {
	case "added":
		return jr.removeAtPath(obj, pathParts)
	case "removed":
		if diff.FromValue == nil {
			return fmt.Errorf("missing from_value for removed field")
		}
		value, err := jr.parseValue(*diff.FromValue, diff.ValueType)
		if err != nil {
			return err
		}
		return jr.setAtPath(obj, pathParts, value)
	case "modified":
		if diff.FromValue == nil {
			return fmt.Errorf("missing from_value for modified field")
		}
		value, err := jr.parseValue(*diff.FromValue, diff.ValueType)
		if err != nil {
			return err
		}
		return jr.setAtPath(obj, pathParts, value)
	}

	return nil
}

func (jr *JSONReconstructor) parsePath(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}

func (jr *JSONReconstructor) setAtPath(obj interface{}, pathParts []string, value interface{}) error {
	if len(pathParts) == 0 {
		return fmt.Errorf("empty path")
	}

	if len(pathParts) == 1 {
		switch o := obj.(type) {
		case map[string]interface{}:
			o[pathParts[0]] = value
		default:
			return fmt.Errorf("cannot set value on non-object")
		}
		return nil
	}

	switch o := obj.(type) {
	case map[string]interface{}:
		key := pathParts[0]
		if o[key] == nil {
			o[key] = make(map[string]interface{})
		}
		return jr.setAtPath(o[key], pathParts[1:], value)
	default:
		return fmt.Errorf("cannot navigate through non-object")
	}
}

func (jr *JSONReconstructor) removeAtPath(obj interface{}, pathParts []string) error {
	if len(pathParts) == 0 {
		return fmt.Errorf("empty path")
	}

	if len(pathParts) == 1 {
		switch o := obj.(type) {
		case map[string]interface{}:
			delete(o, pathParts[0])
		default:
			return fmt.Errorf("cannot remove from non-object")
		}
		return nil
	}

	switch o := obj.(type) {
	case map[string]interface{}:
		key := pathParts[0]
		if o[key] == nil {
			return nil
		}
		return jr.removeAtPath(o[key], pathParts[1:])
	default:
		return fmt.Errorf("cannot navigate through non-object")
	}
}

func (jr *JSONReconstructor) parseValue(valueStr, valueType string) (interface{}, error) {
	// Handle null values first
	if valueStr == "null" || valueType == "null" {
		return nil, nil
	}

	switch valueType {
	case "string":
		return valueStr, nil
	case "number":
		if intVal, err := strconv.Atoi(valueStr); err == nil {
			return float64(intVal), nil
		}
		if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return floatVal, nil
		}
		return nil, fmt.Errorf("invalid number: %s", valueStr)
	case "boolean":
		switch valueStr {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
		return strconv.ParseBool(valueStr)
	case "array", "object":
		var result interface{}
		if err := json.Unmarshal([]byte(valueStr), &result); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return result, nil
	case "null":
		return nil, nil
	default:
		return valueStr, nil
	}
}

func (jr *JSONReconstructor) deepCopy(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case map[string]interface{}:
		copy := make(map[string]interface{})
		for key, value := range v {
			copy[key] = jr.deepCopy(value)
		}
		return copy
	case []interface{}:
		copy := make([]interface{}, len(v))
		for i, value := range v {
			copy[i] = jr.deepCopy(value)
		}
		return copy
	default:
		return v
	}
}
