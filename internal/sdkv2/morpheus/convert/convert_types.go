package convert

import (
	"strconv"
)

// This file contains some helper methods for things like
// converting values to one type or another.

func StringToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)

	return v
}

func IntToString(n int) string {
	return strconv.FormatInt(int64(n), 10)
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

// Bool returns a pointer To the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}

	return false
}

// BoolSlice converts a slice of bool values inTo a slice of
// bool pointers
func BoolSlice(src []bool) []*bool {
	dst := make([]*bool, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

// BoolValueSlice converts a slice of bool pointers inTo a slice of
// bool values
func BoolValueSlice(src []*bool) []bool {
	dst := make([]bool, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}
