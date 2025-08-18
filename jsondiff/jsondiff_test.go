package jsondiff

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONDifferComprehensive(t *testing.T) {
	testCases := []struct {
		name  string
		json1 string
		json2 string
	}{
		// SIMPLE TEST CASES
		{
			name:  "Simple String Change",
			json1: `{"name": "John"}`,
			json2: `{"name": "Jane"}`,
		},
		{
			name:  "Simple Number Change",
			json1: `{"age": 25}`,
			json2: `{"age": 30}`,
		},
		{
			name:  "Simple Boolean Change",
			json1: `{"active": true}`,
			json2: `{"active": false}`,
		},
		{
			name:  "Simple Addition",
			json1: `{"name": "John"}`,
			json2: `{"name": "John", "age": 25}`,
		},
		{
			name:  "Simple Deletion",
			json1: `{"name": "John", "age": 25}`,
			json2: `{"name": "John"}`,
		},
		{
			name:  "Simple Array Change",
			json1: `{"items": [1, 2, 3]}`,
			json2: `{"items": [1, 2, 4]}`,
		},

		// BASIC TEST CASES
		{
			name: "Basic Property Changes",
			json1: `{
				"name": "John",
				"age": 25,
				"active": true
			}`,
			json2: `{
				"name": "John Doe",
				"age": 30,
				"active": false,
				"email": "john@example.com"
			}`,
		},
		{
			name: "Array Changes (Arrays-as-Whole)",
			json1: `{
				"fruits": ["apple", "banana"],
				"numbers": [1, 2, 3]
			}`,
			json2: `{
				"fruits": ["apple", "orange", "grape"],
				"numbers": [1, 3, 4, 5]
			}`,
		},

		// COMPLEX TEST CASES
		{
			name: "Nested Object Changes",
			json1: `{
				"user": {
					"name": "John",
					"address": {
						"city": "NYC",
						"zip": "10001"
					}
				}
			}`,
			json2: `{
				"user": {
					"name": "Jane",
					"address": {
						"city": "NYC",
						"zip": "10002",
						"country": "USA"
					},
					"phone": "555-1234"
				}
			}`,
		},
		{
			name: "Complex Multi-dimensional Arrays",
			json1: `{
				"matrix": [[[1,2],[3,4]], [[5,6]]],
				"data": {"arr": [1,2,3]}
			}`,
			json2: `{
				"matrix": [[[1,2,3],[3,4]], [[5,6,7]]],
				"data": {"arr": [4,5,6]}
			}`,
		},
		{
			name: "Deeply Nested Structure",
			json1: `{
				"level1": {
					"level2": {
						"level3": {
							"level4": {
								"value": "deep",
								"array": [1, 2, 3]
							}
						}
					}
				}
			}`,
			json2: `{
				"level1": {
					"level2": {
						"level3": {
							"level4": {
								"value": "deeper",
								"array": [1, 2, 3, 4],
								"new": "added"
							}
						}
					}
				}
			}`,
		},
		{
			name: "Mixed Data Types",
			json1: `{
				"string": "text",
				"number": 42,
				"boolean": true,
				"null_value": null,
				"array": [1, "two", true, null],
				"object": {"nested": "value"}
			}`,
			json2: `{
				"string": "updated",
				"number": 100,
				"boolean": false,
				"null_value": "not null",
				"array": [1, "three", false, {"obj": "in array"}],
				"object": {"nested": "updated", "added": "field"},
				"new_field": "added"
			}`,
		},

		// EDGE CASES
		{
			name: "Empty to Non-empty",
			json1: `{
				"items": [],
				"metadata": {}
			}`,
			json2: `{
				"items": [1, 2, 3],
				"metadata": {"created": "2023-01-01"}
			}`,
		},
		{
			name: "Non-empty to Empty",
			json1: `{
				"items": [1, 2, 3],
				"metadata": {"created": "2023-01-01"}
			}`,
			json2: `{
				"items": [],
				"metadata": {}
			}`,
		},
		{
			name: "Null Value Handling",
			json1: `{
				"field1": null,
				"field2": "value",
				"field3": null
			}`,
			json2: `{
				"field1": "not null",
				"field2": null,
				"field4": null
			}`,
		},
		{
			name: "Type Changes",
			json1: `{
				"value1": "string",
				"value2": 42,
				"value3": true,
				"value4": [1, 2, 3]
			}`,
			json2: `{
				"value1": 123,
				"value2": "number to string",
				"value3": {"bool": "to object"},
				"value4": "array to string"
			}`,
		},
		{
			name:  "Empty JSON Objects",
			json1: `{}`,
			json2: `{"new": "field"}`,
		},
		{
			name:  "Single Field JSON",
			json1: `{"only": "field"}`,
			json2: `{}`,
		},
		{
			name: "Large Array Changes",
			json1: `{
				"large_array": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
				"metadata": {"size": 10}
			}`,
			json2: `{
				"large_array": [1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21],
				"metadata": {"size": 11}
			}`,
		},
		{
			name: "Array with Mixed Types",
			json1: `{
				"mixed": [1, "two", true, null, [1,2], {"nested": "obj"}]
			}`,
			json2: `{
				"mixed": ["one", 2, false, "not null", [3,4,5], {"nested": "updated"}]
			}`,
		},
		{
			name: "Duplicate Keys Edge Case",
			json1: `{
				"key": "value1",
				"nested": {
					"key": "nested_value1"
				}
			}`,
			json2: `{
				"key": "value2",
				"nested": {
					"key": "nested_value2",
					"another": "added"
				}
			}`,
		},
		{
			name: "Number Precision",
			json1: `{
				"integer": 42,
				"float": 3.14,
				"large": 1234567890,
				"small": 0.001
			}`,
			json2: `{
				"integer": 43,
				"float": 3.141592,
				"large": 1234567891,
				"small": 0.002
			}`,
		},
		{
			name: "Unicode and Special Characters",
			json1: `{
				"unicode": "Hello 世界",
				"special": "Special chars: !@#$%^&*()",
				"emoji": "😀🎉🚀"
			}`,
			json2: `{
				"unicode": "Hi 世界",
				"special": "More special: <>?[]{}",
				"emoji": "😎🎊🛸",
				"new_unicode": "新しい"
			}`,
		},
		{
			name: "Very Nested Arrays",
			json1: `{
				"matrix4D": [[[[1,2,3],[4,5,6]]], [[[7,8,9]]]]
			}`,
			json2: `{
				"matrix4D": [[[[1,2,3,10],[4,5,6,11]]], [[[7,8,9,12],[13,14,15]]]]
			}`,
		},

		// OBJECT DRILLING VS ARRAY-AS-WHOLE TEST CASES
		{
			name: "Object Drilling - Nested Object Changes",
			json1: `{
				"user": {
					"name": "John",
					"age": 25,
					"tags": ["user", "admin"]
				}
			}`,
			json2: `{
				"user": {
					"name": "Jane",
					"age": 30,
					"tags": ["user", "member", "active"]
				}
			}`,
		},
		{
			name: "Mixed Object and Array Changes",
			json1: `{
				"config": {
					"debug": true,
					"ports": [8080, 3000],
					"database": {
						"host": "localhost",
						"tables": ["users", "logs"]
					}
				}
			}`,
			json2: `{
				"config": {
					"debug": false,
					"ports": [8080, 3000, 9000],
					"database": {
						"host": "remote.db",
						"port": 5432,
						"tables": ["users", "logs", "cache"]
					}
				}
			}`,
		},
		{
			name: "Object with Array Fields",
			json1: `{
				"fruits": {
					"apple": 1,
					"pears": 2,
					"arr": [1, 2, 3]
				},
				"vegetables": {
					"carrot": 5,
					"items": ["onion", "potato"]
				}
			}`,
			json2: `{
				"fruits": {
					"apple": 1,
					"beans": 2,
					"arr": [1, 3, 4]
				},
				"vegetables": {
					"carrot": 5,
					"items": ["onion", "tomato", "pepper"]
				}
			}`,
		},
		{
			name: "Deeply Nested Object vs Array",
			json1: `{
				"level1": {
					"level2": {
						"properties": {
							"name": "test",
							"value": 100
						},
						"items": [1, 2, 3]
					}
				}
			}`,
			json2: `{
				"level1": {
					"level2": {
						"properties": {
							"name": "updated",
							"value": 200,
							"active": true
						},
						"items": [1, 2, 4, 5]
					}
				}
			}`,
		},
		{
			name: "Object Field Addition and Removal",
			json1: `{
				"settings": {
					"theme": "dark",
					"notifications": true,
					"features": ["feature1", "feature2"]
				},
				"metadata": {
					"version": "1.0",
					"author": "user"
				}
			}`,
			json2: `{
				"settings": {
					"theme": "light",
					"privacy": false,
					"features": ["feature1", "feature3", "feature4"]
				},
				"metadata": {
					"version": "1.1",
					"author": "user",
					"created": "2023-01-01"
				}
			}`,
		},
		{
			name: "Array in Array vs Object in Object",
			json1: `{
				"matrix": [
					[1, 2, 3],
					[4, 5, 6]
				],
				"config": {
					"database": {
						"host": "localhost",
						"port": 3306
					},
					"cache": {
						"type": "redis",
						"ttl": 3600
					}
				}
			}`,
			json2: `{
				"matrix": [
					[1, 2, 3, 7],
					[4, 5, 8]
				],
				"config": {
					"database": {
						"host": "remote",
						"port": 5432,
						"ssl": true
					},
					"cache": {
						"type": "memory",
						"size": "1GB"
					}
				}
			}`,
		},
		{
			name: "Objects with Null and Array Fields",
			json1: `{
				"user": {
					"name": "John",
					"email": null,
					"roles": ["user"]
				},
				"preferences": {
					"theme": null,
					"languages": ["en", "es"]
				}
			}`,
			json2: `{
				"user": {
					"name": "Jane",
					"email": "jane@example.com",
					"roles": ["admin", "user"]
				},
				"preferences": {
					"theme": "dark",
					"languages": ["en", "fr", "de"],
					"timezone": "UTC"
				}
			}`,
		},
		{
			name: "Mixed Data Types in Objects",
			json1: `{
				"data": {
					"string_field": "text",
					"number_field": 42,
					"array_field": [1, 2, 3],
					"object_field": {
						"nested": "value"
					},
					"null_field": null
				}
			}`,
			json2: `{
				"data": {
					"string_field": "updated",
					"number_field": 100,
					"array_field": [4, 5, 6, 7],
					"object_field": {
						"nested": "updated",
						"added": "field"
					},
					"null_field": "not null",
					"new_field": true
				}
			}`,
		},
		{
			name: "Complex Nested Structure",
			json1: `{
				"application": {
					"name": "MyApp",
					"version": "1.0.0",
					"modules": [
						{
							"name": "auth",
							"routes": ["/login", "/logout"]
						}
					],
					"config": {
						"database": {
							"connections": [
								{"host": "db1", "port": 5432}
							]
						}
					}
				}
			}`,
			json2: `{
				"application": {
					"name": "MyApp",
					"version": "2.0.0",
					"modules": [
						{
							"name": "auth",
							"routes": ["/login", "/logout", "/register"]
						},
						{
							"name": "api",
							"routes": ["/api/users"]
						}
					],
					"config": {
						"database": {
							"connections": [
								{"host": "db1", "port": 5432},
								{"host": "db2", "port": 5433}
							],
							"pool_size": 10
						}
					}
				}
			}`,
		},
	}

	differ := NewJSONDiffer()
	reconstructor := NewJSONReconstructor()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test diff generation
			diffs, err := differ.Compare(tc.json1, tc.json2)
			if err != nil {
				t.Fatalf("Diff generation failed: %v", err)
			}

			// Test reconstruction
			reconstructed, err := reconstructor.ReverseDiff(tc.json2, diffs)
			if err != nil {
				t.Fatalf("Reconstruction failed: %v", err)
			}

			// Verify reconstruction
			var originalObj, reconstructedObj interface{}
			if err := json.Unmarshal([]byte(tc.json1), &originalObj); err != nil {
				t.Fatalf("Failed to unmarshal original JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(reconstructed), &reconstructedObj); err != nil {
				t.Fatalf("Failed to unmarshal reconstructed JSON: %v", err)
			}

			if !reflect.DeepEqual(originalObj, reconstructedObj) {
				t.Errorf("Reconstruction failed.\nExpected: %s\nGot: %s", tc.json1, reconstructed)
			}
		})
	}
}

func TestJSONDifferBasicFunctionality(t *testing.T) {
	differ := NewJSONDiffer()

	t.Run("No differences", func(t *testing.T) {
		json1 := `{"name": "John", "age": 30}`
		json2 := `{"name": "John", "age": 30}`

		diffs, err := differ.Compare(json1, json2)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(diffs) != 0 {
			t.Errorf("Expected 0 diffs, got %d", len(diffs))
		}
	})

	t.Run("Simple modification", func(t *testing.T) {
		json1 := `{"name": "John"}`
		json2 := `{"name": "Jane"}`

		diffs, err := differ.Compare(json1, json2)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(diffs) != 1 {
			t.Fatalf("Expected 1 diff, got %d", len(diffs))
		}

		diff := diffs[0]
		if diff.ChangeType != "modified" {
			t.Errorf("Expected change type 'modified', got '%s'", diff.ChangeType)
		}
		if diff.FullPath != "/name" {
			t.Errorf("Expected path '/name', got '%s'", diff.FullPath)
		}
		if diff.FromValue == nil || *diff.FromValue != "John" {
			t.Errorf("Expected from value 'John', got %v", diff.FromValue)
		}
		if diff.ToValue == nil || *diff.ToValue != "Jane" {
			t.Errorf("Expected to value 'Jane', got %v", diff.ToValue)
		}
	})

	t.Run("Field addition", func(t *testing.T) {
		json1 := `{"name": "John"}`
		json2 := `{"name": "John", "age": 30}`

		diffs, err := differ.Compare(json1, json2)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(diffs) != 1 {
			t.Fatalf("Expected 1 diff, got %d", len(diffs))
		}

		diff := diffs[0]
		if diff.ChangeType != "added" {
			t.Errorf("Expected change type 'added', got '%s'", diff.ChangeType)
		}
		if diff.FullPath != "/age" {
			t.Errorf("Expected path '/age', got '%s'", diff.FullPath)
		}
	})

	t.Run("Field removal", func(t *testing.T) {
		json1 := `{"name": "John", "age": 30}`
		json2 := `{"name": "John"}`

		diffs, err := differ.Compare(json1, json2)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(diffs) != 1 {
			t.Fatalf("Expected 1 diff, got %d", len(diffs))
		}

		diff := diffs[0]
		if diff.ChangeType != "removed" {
			t.Errorf("Expected change type 'removed', got '%s'", diff.ChangeType)
		}
		if diff.FullPath != "/age" {
			t.Errorf("Expected path '/age', got '%s'", diff.FullPath)
		}
	})
}

func TestJSONDifferInvalidJSON(t *testing.T) {
	differ := NewJSONDiffer()

	t.Run("Invalid first JSON", func(t *testing.T) {
		json1 := `{"name": "John",}`
		json2 := `{"name": "Jane"}`

		_, err := differ.Compare(json1, json2)
		if err == nil {
			t.Fatal("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("Invalid second JSON", func(t *testing.T) {
		json1 := `{"name": "John"}`
		json2 := `{"name": "Jane",}`

		_, err := differ.Compare(json1, json2)
		if err == nil {
			t.Fatal("Expected error for invalid JSON, got nil")
		}
	})
}
