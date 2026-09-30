// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package errfmt_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func resp(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}
}

// decode decodes body into a generated delete-response model the way the SDK
// does, so the test exercises the real types (Success pointer, remaining
// fields in AdditionalProperties).
func decode(t *testing.T, body string) *sdk.DeleteUser200Response {
	t.Helper()

	var r sdk.DeleteUser200Response
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatal(err)
	}

	return &r
}

func TestUnitCheckDeleted(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		result  any
		err     error
		resp    *http.Response
		wantErr string // "" for nil
		refused bool
	}{
		{"200 success:true", decode(t, `{"success":true}`), nil, resp(200, ""), "", false},
		{"200 no success field", decode(t, `{}`), nil, resp(200, ""), "", false},
		{
			"404 is already gone", nil, errors.New("404 Not Found"),
			resp(404, `{"success":false,"msg":"not found"}`), "", false,
		},
		{"transport error", nil, errors.New("dial tcp: connection refused"), nil, "connection refused", false},
		{"500", nil, errors.New("500"), resp(500, `{"msg":"boom"}`), "500", false},
		{
			"200 success:false with msg", decode(t, `{"success":false,"msg":"Network is in use"}`), nil, resp(200, ""),
			"Morpheus refused the delete: Network is in use", true,
		},
		{
			"200 success:false with errors",
			decode(t, `{"success":false,"errors":{"hosts":"still attached","zone":"locked"}}`), nil, resp(200, ""),
			"Morpheus refused the delete: hosts: still attached; zone: locked", true,
		},
		{
			"200 success:false, msg and errors", decode(t, `{"success":false,"msg":"Refused","errors":{"a":"b"}}`), nil, resp(200, ""),
			"Morpheus refused the delete: Refused; a: b", true,
		},
		{
			"200 success:false, no reason", decode(t, `{"success":false}`), nil, resp(200, ""),
			"Morpheus refused the delete: no reason given", true,
		},
		{"nil result on 200", nil, nil, resp(200, ""), "", false},
		{"result of an unrelated type", struct{ X int }{1}, nil, resp(200, ""), "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := errfmt.CheckDeleted(tc.result, tc.err, tc.resp)

			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("expected an error containing %q", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
			}

			if errfmt.IsDeleteRefused(err) != tc.refused {
				t.Fatalf("IsDeleteRefused = %v, want %v (err: %v)", !tc.refused, tc.refused, err)
			}
		})
	}
}

// Every generated response model shape the Deletes return must satisfy the
// helper: a pointer to a struct with Success *bool and AdditionalProperties.
func TestUnitCheckDeletedAcceptsGeneratedModels(t *testing.T) {
	t.Parallel()

	for name, r := range map[string]any{
		"DeleteUser200Response":       &sdk.DeleteUser200Response{Success: sdk.PtrBool(false), AdditionalProperties: map[string]any{"msg": "x"}},
		"DeleteSubnet200Response":     &sdk.DeleteSubnet200Response{Success: sdk.PtrBool(false)},
		"DeleteOptionList200Response": &sdk.DeleteOptionList200Response{Success: sdk.PtrBool(false)},
	} {
		if err := errfmt.CheckDeleted(r, nil, resp(200, "")); !errfmt.IsDeleteRefused(err) {
			t.Errorf("%s: success:false was not reported as refused (err=%v)", name, err)
		}
	}
}

func TestUnitRetryRefusedDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	refused := errfmt.DeleteRefusedError{Msg: "busy"}

	// Refused twice, then succeeds: three calls, nil.
	calls := 0
	err := errfmt.RetryRefusedDelete(ctx, "thing", 5, time.Millisecond, func() error {
		calls++
		if calls < 3 {
			return refused
		}

		return nil
	})

	if err != nil || calls != 3 {
		t.Fatalf("expected success after 3 calls, got err=%v calls=%d", err, calls)
	}

	// Always refused: attempts exhausted, the refusal is returned.
	calls = 0
	err = errfmt.RetryRefusedDelete(ctx, "thing", 3, time.Millisecond, func() error {
		calls++

		return refused
	})

	if !errfmt.IsDeleteRefused(err) || calls != 3 {
		t.Fatalf("expected the refusal after 3 calls, got err=%v calls=%d", err, calls)
	}

	// A non-refusal error is not retried.
	calls = 0
	other := errors.New("500")
	err = errfmt.RetryRefusedDelete(ctx, "thing", 3, time.Millisecond, func() error {
		calls++

		return other
	})

	if !errors.Is(err, other) || calls != 1 {
		t.Fatalf("expected the error at once, got err=%v calls=%d", err, calls)
	}

	// attempts < 1 still runs once.
	calls = 0
	_ = errfmt.RetryRefusedDelete(ctx, "thing", 0, time.Millisecond, func() error {
		calls++

		return refused
	})

	if calls != 1 {
		t.Fatalf("attempts=0 must run once, ran %d", calls)
	}

	// A cancelled context stops the wait and reports both.
	cctx, cancel := context.WithCancel(ctx)
	cancel()

	err = errfmt.RetryRefusedDelete(cctx, "thing", 3, time.Hour, func() error { return refused })
	if !errfmt.IsDeleteRefused(err) || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected refusal wrapped with the context error, got %v", err)
	}
}
