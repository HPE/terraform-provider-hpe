package tests

import (
	"testing"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
)

func TestIsEmpty(t *testing.T) {
	// Custom struct for testing
	type MyStruct struct {
		Name  string
		Age   int
		Valid bool
	}

	// Another custom struct with a pointer field
	type MyStructWithPtr struct {
		Name *string
		Age  int
	}

	// Struct with a slice field
	type StructWithSlice struct {
		ID      int
		Items   []string
		Data    []byte
		Numbers []int
	}

	// Struct with a map field
	type StructWithMap struct {
		ID        string
		Settings  map[string]string
		Counters  map[string]int
		UserProps map[int]interface{}
	}

	// Struct with both slice and map fields
	type StructWithCollections struct {
		Version int
		Tags    []string
		Config  map[string]interface{}
		Enabled bool
	}

	// Test cases
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		// Primitive types
		{"nil", nil, true},
		{"int zero", 0, true},
		{"int non-zero", 1, false},
		{"string empty", "", true},
		{"string non-empty", "hello", false},
		{"bool false", false, true},
		{"bool true", true, false},
		{"float64 zero", 0.0, true},
		{"float64 non-zero", 1.23, false},

		// Slices
		{"slice nil", ([]int)(nil), true},
		{"slice empty", []int{}, true},
		{"slice non-empty", []int{1, 2}, false},

		// Maps
		{"map nil", (map[string]int)(nil), true},
		{"map empty", map[string]int{}, true},
		{"map non-empty", map[string]int{"a": 1}, false},

		// Other collection/reference types
		{"interface nil", (interface{})(nil), true},

		// Pointers
		{"int pointer nil", (*int)(nil), true},
		{"string pointer nil", (*string)(nil), true},
		{"struct pointer nil", (*MyStruct)(nil), true},
		{"struct pointer to zero", &MyStruct{}, true},
		{"struct pointer to non-zero", &MyStruct{Name: "Alice"}, false},

		// Structs (basic)
		{"struct zero", MyStruct{}, true},
		{"struct with string field", MyStruct{Name: "Bob"}, false},
		{"struct with int field", MyStruct{Age: 30}, false},
		{"struct with bool field", MyStruct{Valid: true}, false},
		{"struct with all fields zero", MyStruct{Name: "", Age: 0, Valid: false}, true},

		// Structs with pointer fields
		{"struct with nil pointer field", MyStructWithPtr{Name: nil, Age: 0}, true},
		{"struct with non-nil pointer field (non-empty string)", func() MyStructWithPtr {
			s := "data"
			return MyStructWithPtr{Name: &s, Age: 0}
		}(), false},
		{"struct with int field in ptr struct", MyStructWithPtr{Age: 25}, false},

		// Structs with Slice Fields
		{"StructWithSlice zero", StructWithSlice{}, true},
		// {"StructWithSlice with empty slice", StructWithSlice{Items: []string{}}, true},
		{"StructWithSlice with non-empty slice (string)", StructWithSlice{Items: []string{"a"}}, false},
		{"StructWithSlice with non-empty slice (byte)", StructWithSlice{Data: []byte{1}}, false},
		{"StructWithSlice with non-empty slice (int)", StructWithSlice{Numbers: []int{5}}, false},
		{"StructWithSlice with int field non-zero", StructWithSlice{ID: 1}, false},
		{"StructWithSlice with int field and empty slice", StructWithSlice{ID: 1, Items: []string{}}, false},

		// Structs with Map Fields
		{"StructWithMap zero", StructWithMap{}, true},
		// {"StructWithMap with empty map", StructWithMap{Settings: map[string]string{}}, true},
		{"StructWithMap with non-empty map (string)", StructWithMap{Settings: map[string]string{"k": "v"}}, false},
		{"StructWithMap with non-empty map (int)", StructWithMap{Counters: map[string]int{"c": 1}}, false},
		{"StructWithMap with non-empty map (interface)", StructWithMap{UserProps: map[int]any{1: "value"}}, false},
		{"StructWithMap with string field non-empty", StructWithMap{ID: "test"}, false},
		{"StructWithMap with string field and empty map", StructWithMap{ID: "test", Settings: map[string]string{}}, false},

		// Structs with both Slice and Map Fields
		{"StructWithCollections zero", StructWithCollections{}, true},
		// {"StructWithCollections with empty slice", StructWithCollections{Tags: []string{}}, true},
		// {"StructWithCollections with empty map", StructWithCollections{Config: map[string]interface{}{}}, true},
		// {"StructWithCollections with empty slice and map", StructWithCollections{Tags: []string{}, Config: map[string]interface{}{}}, true},
		{"StructWithCollections with non-empty slice", StructWithCollections{Tags: []string{"tag1"}}, false},
		{"StructWithCollections with non-empty map", StructWithCollections{Config: map[string]any{"key": 1}}, false},
		{"StructWithCollections with non-zero primitive field", StructWithCollections{Version: 1}, false},
		{"StructWithCollections with non-zero bool field", StructWithCollections{Enabled: true}, false},
		{"StructWithCollections with non-empty slice and non-zero primitive", StructWithCollections{Tags: []string{"tag1"}, Version: 1}, false},

		// Time.Time
		{"time.Time zero", time.Time{}, true},
		{"time.Time non-zero", time.Now(), false},

		// Arrays
		{"array int zero", [3]int{}, true},
		{"array int non-zero", [3]int{1, 0, 0}, false},
		{"array string zero", [2]string{}, true},
		{"array string non-zero", [2]string{"", "a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sdk.IsEmpty(tt.input)
			if got != tt.expected {
				t.Errorf("sdk.IsEmpty(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
