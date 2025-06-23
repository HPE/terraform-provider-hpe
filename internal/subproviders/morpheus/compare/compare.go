package compare

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	ErrorRecoveredFromPanic = "IsSubset recovered from panic: "
	ErrorNotSubset          = "sub is not a subset of super"
)

// We originally needed this for comparing permission sets between plan and API.
// Permission sets can be either set by the user, or computed by the API if the user
// leaves the field blank.

// In addition to being set by the user, the API will compute additional values
// that the user may not have explicitly set. This complicates comparison using
// solely unmarshaling into one of the generated API structs as we can't range over
// and recursively walk the struct without using reflection.
func IsSubset(sub, super any) (eq bool, err error) {
	// catch-all for panic edge cases
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s%v", ErrorRecoveredFromPanic, r)
			eq = false
		}
	}()

	// this will return the root error back up the call stack
	// this is especially required if we recover from a panic
	if err != nil {
		return false, err
	}

	vSub := reflect.ValueOf(sub)
	vSuper := reflect.ValueOf(super)

	// not equal if types are not the same
	if vSub.Type() != vSuper.Type() {
		return false, errors.New(ErrorNotSubset)
	}

	if vSub.Kind() == reflect.Pointer {
		vSub = vSub.Elem()
	}

	if vSuper.Kind() == reflect.Pointer {
		vSuper = vSuper.Elem()
	}

	// check for different types of fields
	switch vSub.Kind() {
	case reflect.Struct:
		for i := range vSub.NumField() {
			fieldSub := vSub.Field(i)
			fieldSuper := vSuper.Field(i)
			fieldType := vSub.Type().Field(i)

			// Skip unexported/private struct fields; their PkgPath will be non-empty
			// trying to access unexported struct fields will result in a panic
			// https://pkg.go.dev/reflect#Value.Interface
			if !fieldType.IsExported() {
				continue
			}

			// If the field's value is equal to its zero value, skip it
			if reflect.DeepEqual(fieldSub.Interface(), reflect.Zero(fieldSub.Type()).Interface()) {
				continue
			}

			if _, err := IsSubset(fieldSub.Interface(), fieldSuper.Interface()); err != nil {
				// propagate the root error
				return false, err
			}

		}

		return true, nil

	case reflect.Map:
		if vSuper.Kind() != reflect.Map {
			return false, nil
		}

		for _, key := range vSub.MapKeys() {
			valSub := vSub.MapIndex(key)
			valSuper := vSuper.MapIndex(key)

			if !valSuper.IsValid() {
				return false, errors.New(ErrorNotSubset)
			}

			if _, err := IsSubset(valSub.Interface(), valSuper.Interface()); err != nil {
				// propagate the root error
				return false, err
			}
		}

		return true, nil

	case reflect.Slice, reflect.Array:
		// Order-insensitive slice comparison: every element in sub must exist in super
		used := make([]bool, vSuper.Len())
		for i := range vSub.Len() {
			// has vSub.Index(i) been found in vSuper?
			found := false
			for j := range vSuper.Len() {
				if used[j] {
					continue
				}
				if _, err := IsSubset(vSub.Index(i).Interface(), vSuper.Index(j).Interface()); err == nil {
					used[j] = true
					found = true

					break
				}
			}

			if !found {
				return false, errors.New(ErrorNotSubset)
			}
		}

		return true, nil

	// for all other types
	default:
		eq := reflect.DeepEqual(sub, super)
		if eq {
			return true, nil
		}

		return false, errors.New(ErrorNotSubset)
	}
}
