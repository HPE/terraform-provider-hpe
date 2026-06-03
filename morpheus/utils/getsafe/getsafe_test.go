package getsafe_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
	"github.com/stretchr/testify/assert"
)

// Mirrors SDK struct patterns: pointer fields with nested structs.
type simpleStruct struct {
	name *string
}

type complexStruct struct {
	name   *string
	nested *simpleStruct
}

func TestGetSafe(t *testing.T) {
	t.Parallel()

	stringVal := "test"
	structVal := simpleStruct{name: &stringVal}
	complexVal := complexStruct{name: &stringVal, nested: &structVal}
	mapVal := map[string]any{"key": "value"}

	testGetSafe(t, "Primitive", &stringVal, stringVal)
	testGetSafe(t, "Struct", &structVal, structVal)
	testGetSafe(t, "NestedStructPtr", complexVal.nested, *complexVal.nested)
	testGetSafe(t, "NestedFieldPtr", complexVal.nested.name, *complexVal.nested.name)
	testGetSafe(t, "Map", &mapVal, mapVal)

	// Nil: returns zero value — typed nil vars mimic SDK nil pointer fields
	var (
		nilString *string
		nilStruct *simpleStruct
		nilSlice  *[]simpleStruct
		nilArray  *[2]simpleStruct
		nilMap    *map[string]any
	)

	testGetSafe(t, "PrimitiveNil", nilString, "")
	testGetSafe(t, "StructNil", nilStruct, simpleStruct{})
	testGetSafe(t, "SliceNil", nilSlice, nil)
	testGetSafe(t, "ArrayNil", nilArray, [2]simpleStruct{})
	testGetSafe(t, "MapNil", nilMap, nil)
}

func TestGetSafeOk(t *testing.T) {
	t.Parallel()

	stringVal := "test"
	structVal := simpleStruct{name: &stringVal}
	complexVal := complexStruct{name: &stringVal, nested: &structVal}
	mapVal := map[string]any{"key": "value"}

	// Non-nil: returns the same pointer and ok=true
	testGetSafeOk(t, "Primitive", &stringVal)
	testGetSafeOk(t, "Struct", &structVal)
	testGetSafeOk(t, "NestedStructPtr", complexVal.nested)
	testGetSafeOk(t, "NestedFieldPtr", complexVal.nested.name)
	testGetSafeOk(t, "Map", &mapVal)

	// Nil: returns nil and ok=false — typed nil vars mimic SDK nil pointer fields
	var (
		nilString *string
		nilStruct *simpleStruct
		nilSlice  *[]simpleStruct
		nilArray  *[2]simpleStruct
		nilMap    *map[string]any
	)

	testGetSafeOkNil(t, "PrimitiveNil", nilString)
	testGetSafeOkNil(t, "StructNil", nilStruct)
	testGetSafeOkNil(t, "SliceNil", nilSlice)
	testGetSafeOkNil(t, "ArrayNil", nilArray)
	testGetSafeOkNil(t, "MapNil", nilMap)
}

// --- helpers ---

func testGetSafe[T any](t *testing.T, name string, ptr *T, expected T) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		got := getsafe.GetSafe(ptr)
		assert.Equal(t, expected, got)
	})
}

func testGetSafeOk[T any](t *testing.T, name string, ptr *T) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		got, ok := getsafe.GetSafeOk(ptr)
		assert.True(t, ok)
		assert.Same(t, ptr, got)
		assert.Equal(t, ptr, got)
	})
}

func testGetSafeOkNil[T any](t *testing.T, name string, ptr *T) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		got, ok := getsafe.GetSafeOk(ptr)
		assert.Nil(t, got)
		assert.False(t, ok)
	})
}
