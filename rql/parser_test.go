package rql

import (
	"testing"
	"time"
)

type TestStruct struct {
	ID        int32     `rql:"type=number"`
	Name      string    `rql:"type=string"`
	IsActive  bool      `rql:"type=bool"`
	CreatedAt time.Time `rql:"type=datetime"`
}

func TestValidateQuery(t *testing.T) {
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
					{Name: "ID", Operator: "eq", Type: "number", Value: 123},
					{Name: "Name", Operator: "like", Type: "string", Value: "test"},
					{Name: "IsActive", Operator: "eq", Type: "bool", Value: true},
					{Name: "CreatedAt", Operator: "eq", Type: "datetime", Value: "2021-09-15T15:53:00Z"},
				},
				Sort: []Sort{
					{Key: "ID", Order: "asc"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   false,
		},
		{
			name: "Invalid filter key",
			query: Query{
				Filters: []Filter{
					{Name: "NonExistentKey", Operator: "eq", Type: "string", Value: "test"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid filter operator",
			query: Query{
				Filters: []Filter{
					{Name: "ID", Operator: "invalid", Type: "number", Value: 123},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid filter value type",
			query: Query{
				Filters: []Filter{
					{Name: "ID", Operator: "eq", Type: "number", Value: "invalid"},
				},
			},
			checkStruct: TestStruct{},
			expectErr:   true,
		},
		{
			name: "Invalid sort key",
			query: Query{
				Sort: []Sort{
					{Key: "NonExistentKey", Order: "asc"},
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
