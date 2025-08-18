package jsondiff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wI2L/jsondiff"
)

// WI2LDiffer uses wI2L/jsondiff library as an alternative approach
// This is kept as a standalone implementation for comparison with the custom JSONDiffer
type WI2LDiffer struct{}

// NewWI2LDiffer creates a new wI2L-based differ
func NewWI2LDiffer() *WI2LDiffer {
	return &WI2LDiffer{}
}

// GetName returns the differ name
func (w *WI2LDiffer) GetName() string {
	return "wI2L/jsondiff"
}

// Compare computes diff using wI2L/jsondiff with array-friendly configurations
func (w *WI2LDiffer) Compare(json1, json2 string) ([]DiffEntry, error) {
	// Use minimal configuration to avoid over-consolidation
	opts := []jsondiff.Option{
		// Use LCS for better array handling
		jsondiff.LCS(),
	}

	patch, err := jsondiff.CompareJSON([]byte(json1), []byte(json2), opts...)
	if err != nil {
		return nil, fmt.Errorf("wI2L jsondiff comparison failed: %w", err)
	}

	// Convert patch operations to our DiffEntry format
	diffs := w.convertPatchToDiffEntries(patch)

	// Apply post-processing to achieve arrays-as-whole behavior
	return w.postProcessArrays(diffs, json1, json2), nil
}

// convertPatchToDiffEntries converts operations to DiffEntry format
func (w *WI2LDiffer) convertPatchToDiffEntries(patch jsondiff.Patch) []DiffEntry {
	var diffs []DiffEntry

	for _, op := range patch {
		switch op.Type {
		case "add":
			diffs = append(diffs, DiffEntry{
				FullPath:   op.Path,
				FieldName:  w.extractFieldNameFromPath(op.Path),
				ChangeType: "added",
				ToValue:    w.formatPatchValue(op.Value),
				ValueType:  w.inferValueType(op.Value),
			})
		case "remove":
			diffs = append(diffs, DiffEntry{
				FullPath:   op.Path,
				FieldName:  w.extractFieldNameFromPath(op.Path),
				ChangeType: "removed",
				FromValue:  w.formatPatchValue(op.OldValue),
				ValueType:  w.inferValueType(op.OldValue),
			})
		case "replace":
			diffs = append(diffs, DiffEntry{
				FullPath:   op.Path,
				FieldName:  w.extractFieldNameFromPath(op.Path),
				ChangeType: "modified",
				FromValue:  w.formatPatchValue(op.OldValue),
				ToValue:    w.formatPatchValue(op.Value),
				// For reconstruction compatibility, use old value type
				ValueType:  w.inferValueType(op.OldValue),
			})
		case "move":
			// Handle move operations - could be array reordering
			diffs = append(diffs, DiffEntry{
				FullPath:   op.Path,
				FieldName:  w.extractFieldNameFromPath(op.Path),
				ChangeType: "modified",
				FromValue:  w.formatPatchValue(op.OldValue),
				ToValue:    w.formatPatchValue(op.Value),
				ValueType:  w.inferValueType(op.OldValue),
			})
		case "copy":
			// Handle copy operations
			diffs = append(diffs, DiffEntry{
				FullPath:   op.Path,
				FieldName:  w.extractFieldNameFromPath(op.Path),
				ChangeType: "added",
				ToValue:    w.formatPatchValue(op.Value),
				ValueType:  w.inferValueType(op.Value),
			})
		}
	}

	return diffs
}

// postProcessArrays consolidates all array element operations into array-as-whole operations
func (w *WI2LDiffer) postProcessArrays(diffs []DiffEntry, json1, json2 string) []DiffEntry {
	// Always consolidate any array element changes to whole array operations
	arrayChanges := make(map[string][]DiffEntry)
	var result []DiffEntry
	
	// Parse original JSONs once
	var obj1, obj2 interface{}
	json.Unmarshal([]byte(json1), &obj1)
	json.Unmarshal([]byte(json2), &obj2)
	
	for _, diff := range diffs {
		if w.isArrayElementPath(diff.FullPath) {
			arrayPath := w.getArrayPathFromElement(diff.FullPath)
			arrayChanges[arrayPath] = append(arrayChanges[arrayPath], diff)
		} else {
			result = append(result, diff)
		}
	}
	
	// Consolidate ALL array changes to whole array operations
	for arrayPath, changes := range arrayChanges {
		if len(changes) > 0 {
			// Get the complete arrays from both JSONs
			oldArray := w.getValueAtPath(obj1, arrayPath)
			newArray := w.getValueAtPath(obj2, arrayPath)
			
			// Determine change type based on array existence
			var changeType string
			var fromValue, toValue *string
			
			if oldArray == nil && newArray != nil {
				changeType = "added"
				toValue = w.formatPatchValue(newArray)
			} else if oldArray != nil && newArray == nil {
				changeType = "removed"
				fromValue = w.formatPatchValue(oldArray)
			} else {
				changeType = "modified"
				fromValue = w.formatPatchValue(oldArray)
				toValue = w.formatPatchValue(newArray)
			}
			
			result = append(result, DiffEntry{
				FullPath:   arrayPath,
				FieldName:  w.extractFieldNameFromPath(arrayPath),
				ChangeType: changeType,
				FromValue:  fromValue,
				ToValue:    toValue,
				ValueType:  "array",
			})
		}
	}
	
	return result
}

// isArrayElementPath checks if a path refers to an array element
func (w *WI2LDiffer) isArrayElementPath(path string) bool {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for _, part := range parts {
		if w.isArrayIndex(part) {
			return true
		}
	}
	return false
}

// getArrayPathFromElement extracts the array path from an element path
func (w *WI2LDiffer) getArrayPathFromElement(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for i, part := range parts {
		if w.isArrayIndex(part) {
			if i == 0 {
				return ""
			}
			return "/" + strings.Join(parts[:i], "/")
		}
	}
	return path
}

// isArrayIndex checks if a string represents an array index
func (w *WI2LDiffer) isArrayIndex(s string) bool {
	if s == "-" {
		return true
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// getValueAtPath retrieves value at JSON path
func (w *WI2LDiffer) getValueAtPath(obj interface{}, path string) interface{} {
	if path == "" || path == "/" {
		return obj
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	current := obj

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		default:
			return nil
		}
	}

	return current
}

// extractFieldNameFromPath extracts field name from path
func (w *WI2LDiffer) extractFieldNameFromPath(path string) string {
	if path == "" || path == "/" {
		return "root"
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return "root"
	}

	return parts[len(parts)-1]
}

// formatPatchValue formats a value to string
func (w *WI2LDiffer) formatPatchValue(value interface{}) *string {
	if value == nil {
		str := "null"
		return &str
	}

	switch v := value.(type) {
	case string:
		return &v
	case bool:
		str := fmt.Sprintf("%t", v)
		return &str
	case float64:
		str := fmt.Sprintf("%g", v)
		return &str
	default:
		bytes, _ := json.Marshal(v)
		str := string(bytes)
		return &str
	}
}

// inferValueType determines the JSON type
func (w *WI2LDiffer) inferValueType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch value.(type) {
	case string:
		return "string"
	case bool:
		return "boolean"
	case float64, int:
		return "number"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}