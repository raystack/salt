# JSON Diff

A Go library for calculating differences between JSON documents and reconstructing original JSON from diffs. This implementation treats array changes as complete units rather than element-by-element differences.

## Features

- **JSON Comparison**: Generate detailed diffs between two JSON documents
- **Reconstruction**: Reverse diffs to reconstruct the original JSON
- **Array Handling**: Arrays are compared as complete units for cleaner diffs
- **Type Safety**: Preserves data types during comparison and reconstruction
- **Zero Dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/raystack/salt/data/jsondiff
```

## Usage

### Custom JSONDiffer

```go
package main

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
    "github.com/raystack/salt/data/jsondiff"
)

func main() {
    originalJSON := `{
        "name": "John",
        "matrix4D": [[[[1,2,3],[4,5,6],[4,5,6]]]],
        "old_param": "value",
        "fruits": {
            "apple": 1,
            "pears": 2,
            "arr": [1,2,3]
        }
    }`

    currentJSON := `{
        "name": "John Doe",
        "matrix4D": [[[[1,2,3],[4,5,6,7]]]],
        "age": 30,
        "fruits": {
            "apple": 1,
            "beans": 2,
            "arr": [1,3,4]
        },
        "new_object": {
            "a":1,
            "b":[1,2]
        }
    }`

    // Generate diff using custom differ
    differ := jsondiff.NewJSONDiffer()
    diffs, err := differ.Compare(originalJSON, currentJSON)
    if err != nil {
        fmt.Printf("Error generating diff: %v\n", err)
        return
    }

    fmt.Printf("Generated %d diff entries:\n", len(diffs))
    diffJSON, err := json.MarshalIndent(diffs, "", "  ")
    if err != nil {
        fmt.Printf("Error formatting diff: %v\n", err)
        return
    }
    fmt.Println(string(diffJSON))

    // Test reconstruction
    reconstructor := jsondiff.NewJSONReconstructor()
    reconstructed, err := reconstructor.ReverseDiff(currentJSON, diffs)
    if err != nil {
        fmt.Printf("Error reconstructing: %v\n", err)
        return
    }

    // Verify reconstruction accuracy
    var originalObj, reconstructedObj any
    json.Unmarshal([]byte(originalJSON), &originalObj)
    json.Unmarshal([]byte(reconstructed), &reconstructedObj)

    if reflect.DeepEqual(originalObj, reconstructedObj) {
        fmt.Println("✓ Reconstruction successful")
    } else {
        fmt.Println("✗ Reconstruction failed")
    }
}
```

### wI2L JSONDiffer

```go
package main

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
    "github.com/raystack/salt/data/jsondiff"
)

func main() {
    originalJSON := `{
        "name": "John",
        "matrix4D": [[[[1,2,3],[4,5,6],[4,5,6]]]],
        "old_param": "value",
        "fruits": {
            "apple": 1,
            "pears": 2,
            "arr": [1,2,3]
        }
    }`

    currentJSON := `{
        "name": "John Doe",
        "matrix4D": [[[[1,2,3],[4,5,6,7]]]],
        "age": 30,
        "fruits": {
            "apple": 1,
            "beans": 2,
            "arr": [1,3,4]
        },
        "new_object": {
            "a":1,
            "b":[1,2]
        }
    }`

    // Generate diff using wI2L differ
    differ := jsondiff.NewWI2LDiffer()
    diffs, err := differ.Compare(originalJSON, currentJSON)
    if err != nil {
        fmt.Printf("Error generating diff: %v\n", err)
        return
    }

    fmt.Printf("Generated %d diff entries:\n", len(diffs))
    diffJSON, err := json.MarshalIndent(diffs, "", "  ")
    if err != nil {
        fmt.Printf("Error formatting diff: %v\n", err)
        return
    }
    fmt.Println(string(diffJSON))

    // Test reconstruction
    reconstructor := jsondiff.NewJSONReconstructor()
    reconstructed, err := reconstructor.ReverseDiff(currentJSON, diffs)
    if err != nil {
        fmt.Printf("Error reconstructing: %v\n", err)
        return
    }

    // Verify reconstruction accuracy
    var originalObj, reconstructedObj any
    json.Unmarshal([]byte(originalJSON), &originalObj)
    json.Unmarshal([]byte(reconstructed), &reconstructedObj)

    if reflect.DeepEqual(originalObj, reconstructedObj) {
        fmt.Println("✓ Reconstruction successful")
    } else {
        fmt.Println("✗ Reconstruction failed")
    }
}
```

## Data Structures

### DiffEntry

Each difference is represented as a `DiffEntry`:

```go
type DiffEntry struct {
    FieldName  string  `json:"field_name"`   // Name of the changed field
    ChangeType string  `json:"change_type"`  // "added", "removed", or "modified"
    FromValue  *string `json:"from_value,omitempty"`  // Original value (for removed/modified)
    ToValue    *string `json:"to_value,omitempty"`    // New value (for added/modified)
    FullPath   string  `json:"full_path"`    // Full JSON path (e.g., "/fruits/apple")
    ValueType  string  `json:"value_type"`   // "string", "number", "boolean", "array", "object", "null"
}
```

### Change Types

- **"added"**: Field was added in the new JSON
- **"removed"**: Field was deleted from the original JSON  
- **"modified"**: Field value was changed

## Array Handling

Unlike element-by-element array diffs, this library treats entire arrays as single units. When an array changes, the entire array is marked as "modified" with the complete old and new values.

**Example:**
```json
// Original: [1,2,3]
// Current:  [1,3,4]
// Result:   One "modified" entry with full array values
```

This approach provides cleaner, more predictable diffs for complex nested structures.

## API Reference

### JSONDiffer

```go
func NewJSONDiffer() *JSONDiffer
func (jd *JSONDiffer) Compare(json1, json2 string) ([]DiffEntry, error)
```

### WI2LDiffer

```go
func NewWI2LDiffer() *WI2LDiffer
func (w *WI2LDiffer) Compare(json1, json2 string) ([]DiffEntry, error)
```

### JSONReconstructor

```go
func NewJSONReconstructor() *JSONReconstructor
func (jr *JSONReconstructor) ReverseDiff(currentJSON string, diffs []DiffEntry) (string, error)
```

## Error Handling

The library returns detailed error messages for:
- Invalid JSON syntax
- Malformed diff entries
- Path resolution issues during reconstruction
