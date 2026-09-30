// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package errfmt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Morpheus renders a delete it accepted as a request but refused as an
// operation as HTTP 200 with a body of {"success": false, "msg": ...} -- the
// shared delete response view, used by every API controller. The status code
// alone therefore says nothing about whether the object is gone; the body has
// to be read. A Delete that judges by status code removes the resource from
// state while it remains on the appliance, and the next plan tries to create
// it again (MORPH-17987).

// DeleteRefusedError is a delete Morpheus refused: the request was accepted
// (HTTP 200) and the body said success:false. Msg is Morpheus's own message,
// with any per-field errors appended, or "no reason given".
type DeleteRefusedError struct {
	Msg string
}

func (e DeleteRefusedError) Error() string {
	return "Morpheus refused the delete: " + e.Msg
}

// IsDeleteRefused reports whether err is (or wraps) a DeleteRefusedError.
func IsDeleteRefused(err error) bool {
	var refused DeleteRefusedError

	return errors.As(err, &refused)
}

// structOf dereferences result to the struct it points at; invalid when it is
// nil or not a struct.
func structOf(result any) reflect.Value {
	if result == nil {
		return reflect.Value{}
	}

	v := reflect.ValueOf(result)
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	return v
}

// additionalProperties returns the generated model's AdditionalProperties
// field -- where the unmodelled "msg" and "errors" of a response land.
func additionalProperties(result any) map[string]any {
	v := structOf(result)
	if !v.IsValid() {
		return nil
	}

	f := v.FieldByName("AdditionalProperties")
	if !f.IsValid() {
		return nil
	}

	m, _ := f.Interface().(map[string]any)

	return m
}

// CheckDeleted judges a DELETE by its status *and* its body:
//
//   - a 404 is not an error: the object is already gone, which is what a
//     delete wants;
//   - a transport error or a non-2xx status is the error CheckResponse gives;
//   - a 2xx whose body says success:false is a DeleteRefusedError carrying
//     Morpheus's message;
//   - anything else is success.
//
// result is the decoded response the SDK returned (it may be nil, e.g. on a
// transport error); its dynamic type is one of the generated *200Response
// models.
func CheckDeleted(result any, err error, httpResp *http.Response) error {
	if IsNotFound(httpResp) {
		return nil
	}

	if err := CheckResponse(err, httpResp); err != nil {
		return err
	}

	if refused, msg := refusedBody(result); refused {
		return DeleteRefusedError{Msg: msg}
	}

	return nil
}

// refusedBody reports whether the decoded body says success:false, and the
// message it carries. The generated models have no getters, so the uniform
// `Success *bool` field is read by reflection; a result without one (or a nil
// result) cannot be a refusal.
func refusedBody(result any) (bool, string) {
	v := structOf(result)
	if !v.IsValid() {
		return false, ""
	}

	f := v.FieldByName("Success")
	if !f.IsValid() {
		return false, ""
	}

	switch f.Kind() {
	case reflect.Pointer:
		if f.IsNil() || f.Elem().Kind() != reflect.Bool || f.Elem().Bool() {
			return false, ""
		}
	case reflect.Bool:
		if f.Bool() {
			return false, ""
		}
	default:
		return false, ""
	}

	return true, RefusalMessage(additionalProperties(result))
}

// RefusalMessage renders the "msg" and "errors" of a Morpheus response body
// as one line: the message, then each error as "field: text" in field order;
// "no reason given" when the body has neither.
func RefusalMessage(body map[string]any) string {
	var parts []string

	if m, ok := body["msg"].(string); ok && m != "" {
		parts = append(parts, m)
	}

	if errs, ok := body["errors"].(map[string]any); ok {
		keys := make([]string, 0, len(errs))
		for k := range errs {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s: %v", k, errs[k]))
		}
	}

	if len(parts) == 0 {
		return "no reason given"
	}

	return strings.Join(parts, "; ")
}

// RetryRefusedDelete runs del, and when Morpheus refuses the delete tries
// again up to attempts times, pausing between attempts -- for resources whose
// delete does host-side work that can fail transiently. Any other error, and
// success, return at once; the last refusal is returned when the attempts run
// out. attempts < 1 is treated as 1.
func RetryRefusedDelete(ctx context.Context, what string, attempts int, pause time.Duration, del func() error) error {
	if attempts < 1 {
		attempts = 1
	}

	for attempt := 1; ; attempt++ {
		err := del()
		if err == nil || !IsDeleteRefused(err) || attempt >= attempts {
			return err
		}

		tflog.Warn(ctx, what+": Morpheus refused the delete; retrying",
			map[string]any{"attempt": attempt, "of": attempts, "error": err.Error()})

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w (while waiting to retry: %v)", err, ctx.Err())
		case <-time.After(pause):
		}
	}
}
