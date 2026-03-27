package rql

import (
	"strings"
	"testing"
	"time"
)

func TestValidateQuery(t *testing.T) {
	type TestStruct struct {
		ID        int32     `rql:"name=id,type=number"`
		Name      string    `rql:"name=name,type=string"`
		IsActive  bool      `rql:"name=is_active,type=bool"`
		CreatedAt time.Time `rql:"name=created_at,type=datetime"`
	}

	tests := []struct {
		name        string
		query       Query
		checkStruct TestStruct
		expectErr   bool
	}{
		{
			name: "Valid filters and sort",
			query: Query{
				Filters: []Filter{
					{Name: "ID", Operator: "eq", Value: 123},
					{Name: "Name", Operator: "like", Value: "test"},
					{Name: "is_active", Operator: "eq", Value: true},
					{Name: "created_at", Operator: "eq", Value: "2021-09-15T15:53:00Z"},
				},
				Sort: []Sort{
					{Name: "ID", Order: "asc"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   false,
		},
		{
			name: "Invalid filter key",
			query: Query{
				Filters: []Filter{
					{Name: "NonExistentKey", Operator: "eq", Value: "test"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid filter operator",
			query: Query{
				Filters: []Filter{
					{Name: "ID", Operator: "invalid", Value: 123},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid filter value type",
			query: Query{
				Filters: []Filter{
					{Name: "ID", Operator: "eq", Value: "invalid"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid sort key",
			query: Query{
				Sort: []Sort{
					{Name: "NonExistentKey", Order: "asc"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(&tt.query, tt.checkStruct)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateQuery() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestGetDataTypeOfField(t *testing.T) {
	type TestStruct struct {
		StringField   string    `rql:"name=string_field,type=string"`
		NumberField   int       `rql:"name=number_field,type=number"`
		BoolField     bool      `rql:"name=bool_field,type=bool"`
		DateTimeField time.Time `rql:"name=datetime_field,type=datetime"`
		InvalidField  string    `rql:"name=invalid_field,type=invalid"`
		NoTypeField   string    `rql:"name=no_type_field"` // No type specified
		NoTagField    string    // No tag at all
	}

	tests := []struct {
		name          string
		fieldName     string
		expectedType  string
		expectedError bool
		errorContains string
	}{
		{
			name:          "String field by struct name",
			fieldName:     "StringField",
			expectedType:  "string",
			expectedError: false,
		},
		{
			name:          "String field by tag name",
			fieldName:     "string_field",
			expectedType:  "string",
			expectedError: false,
		},
		{
			name:          "Number field by struct name",
			fieldName:     "NumberField",
			expectedType:  "number",
			expectedError: false,
		},
		{
			name:          "Number field by tag name",
			fieldName:     "number_field",
			expectedType:  "number",
			expectedError: false,
		},
		{
			name:          "Bool field by struct name",
			fieldName:     "BoolField",
			expectedType:  "bool",
			expectedError: false,
		},
		{
			name:          "DateTime field by struct name",
			fieldName:     "DateTimeField",
			expectedType:  "datetime",
			expectedError: false,
		},
		{
			name:          "Invalid field name",
			fieldName:     "NonExistentField",
			expectedType:  "",
			expectedError: true,
			errorContains: "is not a valid field",
		},
		{
			name:          "No type specified in tag",
			fieldName:     "NoTypeField",
			expectedType:  "string", // Should default to string
			expectedError: false,
		},
		{
			name:          "No tag field",
			fieldName:     "NoTagField",
			expectedType:  "string", // Should default to string
			expectedError: false,
		},
	}

	testStruct := TestStruct{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataType, err := GetDataTypeOfField(tt.fieldName, testStruct)

			// Check error cases
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if dataType != tt.expectedType {
				t.Errorf("Expected type '%s', got '%s'", tt.expectedType, dataType)
			}
		})
	}
}
