package jsondiff

import (
	"cmp"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type DiffEntry struct {
	FieldName  string  `json:"field_name"`
	ChangeType string  `json:"change_type"`
	FromValue  *string `json:"from_value,omitempty"`
	ToValue    *string `json:"to_value,omitempty"`
	FullPath   string  `json:"full_path"`
	ValueType  string  `json:"value_type"`
}

type JSONDiffer struct{}

func NewJSONDiffer() *JSONDiffer {
	return &JSONDiffer{}
}

func (jd *JSONDiffer) Compare(json1, json2 string) ([]DiffEntry, error) {
	var obj1, obj2 any

	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return nil, fmt.Errorf("error parsing first JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return nil, fmt.Errorf("error parsing second JSON: %w", err)
	}

	var diffs []DiffEntry
	jd.compareObjects(obj1, obj2, "", &diffs)

	slices.SortFunc(diffs, func(a, b DiffEntry) int {
		return cmp.Compare(a.FullPath, b.FullPath)
	})

	return diffs, nil
}

func (jd *JSONDiffer) compareObjects(obj1, obj2 any, path string, diffs *[]DiffEntry) {
	if reflect.DeepEqual(obj1, obj2) {
		return
	}

	// Handle null value changes when we have a path context
	if path != "" {
		if obj1 == nil && obj2 != nil {
			*diffs = append(*diffs, jd.createDiffEntry(path, "modified", obj1, obj2))
			return
		}
		if obj1 != nil && obj2 == nil {
			*diffs = append(*diffs, jd.createDiffEntry(path, "modified", obj1, obj2))
			return
		}
	} else {
		// Top-level comparison logic for when fields might not exist
		if obj1 == nil && obj2 != nil {
			*diffs = append(*diffs, jd.createDiffEntry(path, "added", nil, obj2))
			return
		}
		if obj1 != nil && obj2 == nil {
			*diffs = append(*diffs, jd.createDiffEntry(path, "removed", obj1, nil))
			return
		}
	}

	type1 := reflect.TypeOf(obj1)
	type2 := reflect.TypeOf(obj2)

	if type1 != type2 {
		*diffs = append(*diffs, jd.createDiffEntry(path, "modified", obj1, obj2))
		return
	}

	switch v1 := obj1.(type) {
	case map[string]any:
		v2 := obj2.(map[string]any)
		jd.compareObjects_internal(v1, v2, path, diffs)

	case []any:
		v2 := obj2.([]any)
		jd.compareArrays(v1, v2, path, diffs)

	default:
		if !reflect.DeepEqual(obj1, obj2) {
			*diffs = append(*diffs, jd.createDiffEntry(path, "modified", obj1, obj2))
		}
	}
}

func (jd *JSONDiffer) compareObjects_internal(obj1, obj2 map[string]any, path string, diffs *[]DiffEntry) {
	allKeys := make(map[string]bool)
	for key := range obj1 {
		allKeys[key] = true
	}
	for key := range obj2 {
		allKeys[key] = true
	}

	for key := range allKeys {
		newPath := jd.buildPath(path, key)

		val1, exists1 := obj1[key]
		val2, exists2 := obj2[key]

		if !exists1 && exists2 {
			*diffs = append(*diffs, jd.createDiffEntry(newPath, "added", nil, val2))
		} else if exists1 && !exists2 {
			*diffs = append(*diffs, jd.createDiffEntry(newPath, "removed", val1, nil))
		} else if exists1 && exists2 {
			jd.compareObjects(val1, val2, newPath, diffs)
		}
	}
}

func (jd *JSONDiffer) compareArrays(arr1, arr2 []any, path string, diffs *[]DiffEntry) {
	if !reflect.DeepEqual(arr1, arr2) {
		*diffs = append(*diffs, jd.createDiffEntry(path, "modified", arr1, arr2))
	}
}

func (jd *JSONDiffer) buildPath(parentPath, key string) string {
	if parentPath == "" {
		return "/" + key
	}
	return parentPath + "/" + key
}

func (jd *JSONDiffer) createDiffEntry(path, changeType string, oldValue, newValue any) DiffEntry {
	entry := DiffEntry{
		FullPath:   path,
		ChangeType: changeType,
		FieldName:  jd.extractFieldName(path),
	}

	switch changeType {
	case "added":
		toVal := jd.formatValue(newValue)
		entry.ToValue = &toVal
		entry.ValueType = jd.getValueType(newValue)

	case "removed":
		fromVal := jd.formatValue(oldValue)
		entry.FromValue = &fromVal
		entry.ValueType = jd.getValueType(oldValue)

	case "modified":
		fromVal := jd.formatValue(oldValue)
		toVal := jd.formatValue(newValue)
		entry.FromValue = &fromVal
		entry.ToValue = &toVal
		// For reconstruction, we need the old value's type, not the new value's type
		entry.ValueType = jd.getValueType(oldValue)
	}

	return entry
}

func (jd *JSONDiffer) extractFieldName(path string) string {
	if path == "" {
		return "root"
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return "root"
	}

	return parts[len(parts)-1]
}

func (jd *JSONDiffer) getValueType(val any) string {
	if val == nil {
		return "null"
	}

	switch val.(type) {
	case string:
		return "string"
	case float64, int, int64:
		return "number"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return "unknown"
	}
}

func (jd *JSONDiffer) formatValue(val any) string {
	if val == nil {
		return "null"
	}

	switch v := val.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'g', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case map[string]any, []any:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	default:
		return fmt.Sprintf("%v", v)
	}
}
