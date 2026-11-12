// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package errfmt_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v5"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
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

// clientErr wraps a cause the way http.Client.Do does: in a *url.Error.
func clientErr(cause error) error {
	return &url.Error{Op: "Get", URL: "https://m/api/clusters/3", Err: cause}
}

// timeoutErr stands in for the unexported net.Error the client returns when
// Client.Timeout is exceeded (Timeout() true, and errors.Is DeadlineExceeded).
type timeoutErr struct{}

func (timeoutErr) Error() string        { return "Get: context deadline exceeded (Client.Timeout exceeded)" }
func (timeoutErr) Timeout() bool        { return true }
func (timeoutErr) Temporary() bool      { return true }
func (timeoutErr) Is(target error) bool { return target == context.DeadlineExceeded }

// refusedLogin / incompleteLogin are what the round-tripper returns when the
// re-login after an expired token is answered with a 401, or never completes.
func refusedLogin() error {
	return &auth.AuthenticationError{
		Err:     fmt.Errorf("could not authenticate: %w", errors.New("401 Unauthorized")),
		Refused: true,
	}
}

func incompleteLogin(cause error) error {
	return &auth.AuthenticationError{Err: fmt.Errorf("could not authenticate: %w", clientErr(cause))}
}

// tokenEndpointDown is the re-login answered 503 by the token endpoint.
func tokenEndpointDown() error {
	return &auth.AuthenticationError{
		Err:        fmt.Errorf("could not authenticate: %w", errors.New("503 Service Unavailable")),
		StatusCode: 503,
	}
}

// TestUnitPollDeleted: the operation of the wait-for-delete loops. A
// successful GET is StillExistsError (poll again), only a 404 ends the wait, a
// connection failure or a 5xx is TransientError (poll again), and every other
// failure -- a 4xx, a refused login, a certificate error, a done context, an
// error the client never produces for a network blip -- is permanent. The
// no-response rows use the shapes http.Client.Do returns (a *url.Error around
// the cause). Returning nil on a 200 ended the cluster's wait on its first
// poll and let a replacement race the deletion job (MORPH-18379).
func TestUnitPollDeleted(t *testing.T) {
	t.Parallel()

	const (
		done = iota
		still
		transient
		permanent
	)

	doneCtx, cancel := context.WithCancel(context.Background())
	cancel()

	refused := clientErr(&net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")})
	reset := clientErr(&net.OpError{Op: "read", Err: errors.New("connection reset by peer")})
	badCert := &tls.CertificateVerificationError{Err: x509.UnknownAuthorityError{}}
	plainPort := tls.RecordHeaderError{Msg: "first record does not look like a TLS handshake"}
	blipLogin := incompleteLogin(&net.OpError{Op: "dial", Err: errors.New("refused")})
	badURL := &url.Error{Op: "parse", URL: "::bad", Err: errors.New("missing protocol scheme")}

	cases := []struct {
		name string
		ctx  context.Context
		err  error
		resp *http.Response
		want int
	}{
		{"404 ends the wait", nil, errors.New("404 Not Found"), resp(404, `{"success":false}`), done},
		{"200 still exists", nil, nil, resp(200, `{"cluster":{"status":"deprovisioning"}}`), still},
		{"connection refused is transient", nil, refused, nil, transient},
		{"connection reset is transient", nil, reset, nil, transient},
		{"DNS failure is transient", nil, clientErr(&net.DNSError{Err: "no such host", Name: "m"}), nil, transient},
		{"client timeout is transient", nil, clientErr(timeoutErr{}), nil, transient},
		{"dial timeout is transient", nil, clientErr(&net.OpError{Op: "dial", Err: timeoutErr{}}), nil, transient},
		{"server closed before answering is transient", nil, clientErr(io.EOF), nil, transient},
		{"truncated response is transient", nil, clientErr(io.ErrUnexpectedEOF), nil, transient},
		{"500 is transient", nil, errors.New("500"), resp(500, `{"msg":"boom"}`), transient},
		{"502 is transient", nil, errors.New("502"), resp(502, `<html>bad gateway</html>`), transient},
		{"503 is transient", nil, errors.New("503"), resp(503, ``), transient},
		{"incomplete re-login with a blip is transient", nil, blipLogin, nil, transient},
		{"re-login answered 503 is transient", nil, tokenEndpointDown(), nil, transient},
		{"wrapped re-login answered 503 is transient", nil, fmt.Errorf("get: %w", tokenEndpointDown()), nil, transient},
		{"400 is permanent", nil, errors.New("400"), resp(400, `{"msg":"bad"}`), permanent},
		{"401 is permanent", nil, errors.New("401"), resp(401, `{"error":"invalid_token"}`), permanent},
		{"403 is permanent", nil, errors.New("403"), resp(403, `{"msg":"forbidden"}`), permanent},
		{"refused re-login is permanent", nil, refusedLogin(), nil, permanent},
		{"wrapped refused re-login is permanent", nil, fmt.Errorf("get: %w", refusedLogin()), nil, permanent},
		{"incomplete re-login with a certificate error is permanent", nil, incompleteLogin(badCert), nil, permanent},
		{"untrusted certificate is permanent", nil, clientErr(badCert), nil, permanent},
		{"hostname mismatch is permanent", nil, clientErr(x509.HostnameError{Host: "m"}), nil, permanent},
		{"TLS to a plain port is permanent", nil, clientErr(plainPort), nil, permanent},
		{"done context makes a refused connection permanent", doneCtx, refused, nil, permanent},
		{"done context makes a client timeout permanent", doneCtx, clientErr(timeoutErr{}), nil, permanent},
		{"round-tripper error is permanent", nil, clientErr(errors.New("round trip failed")), nil, permanent},
		{"unparsable URL is permanent", nil, badURL, nil, permanent},
		{"bare error with no response is permanent", nil, errors.New("request could not be built"), nil, permanent},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := tc.ctx
			if ctx == nil {
				ctx = context.Background()
			}

			err := errfmt.PollDeleted(ctx, "cluster 3", tc.err, tc.resp)

			var (
				stillErr     errfmt.StillExistsError
				transientErr errfmt.TransientError
				permErr      *backoff.PermanentError
			)

			got := done
			switch {
			case err == nil:
			case errors.As(err, &permErr):
				got = permanent
			case errors.As(err, &transientErr):
				got = transient
			case errors.As(err, &stillErr):
				got = still
			default:
				t.Fatalf("PollDeleted() = %T %v: not one of the documented results", err, err)
			}

			if got != tc.want {
				t.Fatalf("PollDeleted() = %v (%d), want %d", err, got, tc.want)
			}

			if err != nil && !strings.Contains(err.Error(), "cluster 3") {
				t.Errorf("error %q does not name the resource", err)
			}

			if got == transient && strings.Contains(err.Error(), "retry") {
				t.Errorf("a transient error is also the final diagnostic when the wait runs out; it must not promise a retry: %q", err)
			}

			if got == transient {
				// Unwrap reaches the (secret-scrubbed) failure CheckResponse built.
				if inner := errors.Unwrap(err); inner == nil || !strings.Contains(err.Error(), inner.Error()) {
					t.Errorf("TransientError does not unwrap to its cause: %v", err)
				}
			}
		})
	}
}

// TestUnitDeletePoller: a poller tolerates MaxConsecutiveTransient-1
// transient failures in a row, gives up on the next, and a poll that reaches
// the appliance resets the count -- so a blip is ridden out but a persistent
// 500 (MORPH-18060) is reported after five intervals, not after the delete
// timeout.
func TestUnitDeletePoller(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := &errfmt.DeletePoller{What: "image 9"}

	var perm *backoff.PermanentError

	for i := 1; i < errfmt.MaxConsecutiveTransient; i++ {
		err := p.Poll(ctx, errors.New("500"), resp(500, `{"msg":"boom"}`))

		var transient errfmt.TransientError
		if !errors.As(err, &transient) {
			t.Fatalf("failure %d: got %T %v, want TransientError", i, err, err)
		}
	}

	// A poll that reaches the appliance resets the count.
	var still errfmt.StillExistsError
	if err := p.Poll(ctx, nil, resp(200, `{"image":{}}`)); !errors.As(err, &still) {
		t.Fatalf("200 after failures: got %T %v, want StillExistsError", err, err)
	}

	for i := 1; i < errfmt.MaxConsecutiveTransient; i++ {
		if err := p.Poll(ctx, errors.New("500"), resp(500, `{"msg":"boom"}`)); errors.As(err, &perm) {
			t.Fatalf("failure %d after a reset: gave up early: %v", i, err)
		}
	}

	err := p.Poll(ctx, errors.New("500"), resp(500, `{"msg":"boom"}`))
	if !errors.As(err, &perm) {
		t.Fatalf("failure %d: got %T %v, want permanent", errfmt.MaxConsecutiveTransient, err, err)
	}

	for _, want := range []string{"image 9", "times in a row", "boom"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("final error %q should contain %q", err, want)
		}
	}

	// Non-transient results pass through unchanged.
	if err := p.Poll(ctx, errors.New("404"), resp(404, ``)); err != nil {
		t.Errorf("404: got %v, want nil", err)
	}

	if err := p.Poll(ctx, errors.New("403"), resp(403, `{}`)); !errors.As(err, &perm) {
		t.Errorf("403: got %T, want permanent", err)
	}
}
