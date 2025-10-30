package helpers

import (
	"fmt"
)

// TypeAssertFail builds an error for reporting a failed type assertion.
// It is useful for working with sdkv2 as it relies heavily on type
// assertions to pull information from state and plan.
// k is the key or name of the object being type asserted, v is its value.
func TypeAssertFail(k string, v any) error {
	return fmt.Errorf("%s: Type assertion failed for value: %v (type: %T)", k, v, v)
}
