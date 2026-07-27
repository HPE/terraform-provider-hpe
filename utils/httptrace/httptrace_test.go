// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package httptrace_test

import (
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/utils/httptrace"
)

// secret values that must never survive Redact.
const (
	secretPassword = "sup3r-s3cret-pw"
	secretToken    = "abcdef0123456789"
)

// These tests mutate the environment, so they must not run in parallel.

func TestRedact(t *testing.T) {
	// Redaction is the default; pin it so an ambient export cannot skew this.
	t.Setenv(httptrace.EnvShowSecrets, "")

	tests := map[string]struct {
		dump        string
		mustNotHave []string
		mustHave    []string
	}{
		"oauth password grant body": {
			dump: "POST /oauth/token HTTP/1.1\r\n" +
				"Host: morpheus.example.com\r\n" +
				"Content-Type: application/x-www-form-urlencoded\r\n" +
				"\r\n" +
				"grant_type=password&username=svc@example.com&password=" +
				secretPassword,
			mustNotHave: []string{secretPassword},
			// The non-secret parts of the trace stay intact for debugging.
			mustHave: []string{"POST /oauth/token", "username=svc@example.com"},
		},
		"authorization header": {
			dump: "GET /api/networks/routers/28 HTTP/1.1\r\n" +
				"Authorization: Bearer " + secretToken + "\r\n" +
				"Connection: keep-alive\r\n",
			mustNotHave: []string{secretToken},
			mustHave: []string{
				"GET /api/networks/routers/28",
				"Connection: keep-alive",
			},
		},
		"token response body": {
			dump: "HTTP/1.1 200 OK\r\n" +
				"Content-Type: application/json\r\n" +
				"\r\n" +
				`{"access_token":"` + secretToken +
				`","refresh_token":"` + secretToken + `","expires_in":3600}`,
			mustNotHave: []string{secretToken},
			mustHave:    []string{`"expires_in":3600`},
		},
		"json password field": {
			dump:        `{"username":"svc@example.com","password":"` + secretPassword + `"}`,
			mustNotHave: []string{secretPassword},
			mustHave:    []string{`"username":"svc@example.com"`},
		},
		"set-cookie header": {
			dump: "HTTP/1.1 200 OK\r\n" +
				"Set-Cookie: JSESSIONID=" + secretToken + "; Path=/\r\n",
			mustNotHave: []string{secretToken},
		},
		"nothing sensitive is left untouched": {
			dump:     "GET /api/whoami HTTP/1.1\r\nAccept: application/json\r\n",
			mustHave: []string{"GET /api/whoami", "Accept: application/json"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := httptrace.Redact(tc.dump)

			for _, secret := range tc.mustNotHave {
				if strings.Contains(got, secret) {
					t.Errorf("secret leaked into trace:\n%s", got)
				}
			}

			for _, want := range tc.mustHave {
				if !strings.Contains(got, want) {
					t.Errorf("expected %q to be preserved, got:\n%s", want, got)
				}
			}
		})
	}
}

func TestShowSecrets(t *testing.T) {
	tests := map[string]struct {
		value string
		set   bool
		want  bool
	}{
		"unset":            {set: false, want: false},
		"empty":            {value: "", set: true, want: false},
		"zero":             {value: "0", set: true, want: false},
		"false":            {value: "false", set: true, want: false},
		"FALSE mixed case": {value: "FaLsE", set: true, want: false},
		"no":               {value: "no", set: true, want: false},
		"off":              {value: "off", set: true, want: false},
		"whitespace only":  {value: "  ", set: true, want: false},
		"one":              {value: "1", set: true, want: true},
		"true":             {value: "true", set: true, want: true},
		"TRUE mixed case":  {value: "TrUe", set: true, want: true},
		"yes":              {value: "yes", set: true, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.set {
				t.Setenv(httptrace.EnvShowSecrets, tc.value)
			} else {
				// t.Setenv restores the prior value on cleanup, so setting a
				// falsey value is equivalent to "not requested" here.
				t.Setenv(httptrace.EnvShowSecrets, "")
			}

			if got := httptrace.ShowSecrets(); got != tc.want {
				t.Errorf("ShowSecrets() = %v, want %v (value %q)",
					got, tc.want, tc.value)
			}
		})
	}
}

// TestRedactDisabled proves the opt-out returns the dump untouched.
func TestRedactDisabled(t *testing.T) {
	t.Setenv(httptrace.EnvShowSecrets, "1")

	dump := "GET /api/whoami HTTP/1.1\r\n" +
		"Authorization: Bearer " + secretToken + "\r\n"

	if got := httptrace.Redact(dump); got != dump {
		t.Errorf("expected dump to be returned verbatim, got:\n%s", got)
	}
}

func TestIsEnabled(t *testing.T) {
	t.Setenv(httptrace.EnvEnabled, "true")

	if !httptrace.IsEnabled() {
		t.Error("IsEnabled() = false, want true when the variable is set")
	}
}
