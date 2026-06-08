package getsafe

// Get returns the zero value if v is nil, otherwise returns *v.
func Get[T any](v *T) T {
	if v == nil {
		var zero T

		return zero
	}

	return *v
}

// GetOk returns (nil, false) if v is nil, otherwise (v, true).
func GetOk[T any](v *T) (*T, bool) {
	if v == nil {
		return nil, false
	}

	return v, true
}
