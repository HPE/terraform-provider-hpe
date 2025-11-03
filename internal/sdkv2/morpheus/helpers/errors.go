package helpers

import (
	"fmt"
)

// TypeAssertFail builds an error for reporting a failed type assertion.
// It is useful for working with sdkv2 as it relies heavily on type
// assertions to pull information from state and plan.
// k is the key or name of the object being type asserted, v is its value.
func TypeAssertFailError(k string, v any) error {
	return fmt.Errorf("%s: Type assertion failed for value: %v (type: %T)", k, v, v)
}

// EmptySliceAccess builds an error for reporting an attempt to access
// an element from an empty slice.
// It is useful for working with the legacy provider code as it often attempts
// to access the 0th index of a slice without checking its length first.
// k is the key or name of the slice being accessed.
func EmptySliceError(k string) error {
	return fmt.Errorf("%s: Slice is empty", k)
}

// NilPointer builds an error for reporting a nil pointer dereference.
// It is useful for working with the legacy provider code as it often
// dereferences pointers without checking for nil first.
// k is the key or name of the pointer being dereferenced.
func NilPointerError(k string) error {
	return fmt.Errorf("%s: Pointer is nil", k)
}
