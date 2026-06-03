package getsafe

// GetSafe returns the zero value if v is nil, otherwise returns *v.
func GetSafe[T any](v *T) T {
	if v == nil {
		var zero T

		return zero
	}

	return *v
}

// GetSafeOk returns (nil, false) if v is nil, otherwise (v, true).
func GetSafeOk[T any](v *T) (*T, bool) {
	if v == nil {
		return nil, false
	}

	return v, true
}
