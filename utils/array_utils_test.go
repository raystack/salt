package utils

import (
	"testing"
)

func TestStringFoundInArray(t *testing.T) {
	tests := []struct {
		name     string
		item     string
		arr      []string
		expected bool
	}{
		{name: "String found in array", item: "apple", arr: []string{"apple", "banana", "cherry"}, expected: true},
		{name: "String not found in array", item: "orange", arr: []string{"apple", "banana", "cherry"}, expected: false},
		{name: "Empty array", item: "apple", arr: []string{}, expected: false},
		{name: "Empty string in array", item: "", arr: []string{"apple", "banana", "cherry"}, expected: false},
		{name: "Empty string found in array", item: "", arr: []string{"", "banana", "cherry"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringFoundInArray(tt.item, tt.arr)
			if result != tt.expected {
				t.Errorf("StringFoundInArray(%v, %v) = %v; want %v", tt.item, tt.arr, result, tt.expected)
			}
		})
	}
}
