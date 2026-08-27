// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package httptrace

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ http.RoundTripper = TraceRoundTripper{}

const (
	// EnvEnabled turns API request/response tracing on when set.
	EnvEnabled = "MORPHEUS_API_HTTPTRACE"

	// EnvShowSecrets disables redaction of credentials in those traces.
	// #nosec G101 -- environment variable name, not an embedded credential.
	EnvShowSecrets = "MORPHEUS_API_HTTPTRACE_SHOW_SECRETS"
)

// redactedValue replaces any secret found in a traced request or response.
const redactedValue = "[REDACTED]"

var (
	// Credential-bearing headers. Matched per line against a dumped
	// request/response, which uses CRLF separators.
	secretHeaderRE = regexp.MustCompile(
		`(?im)^(Authorization|Proxy-Authorization|Cookie|Set-Cookie):[^\r\n]*`,
	)

	// Secrets in form-encoded bodies, e.g. the OAuth password grant
	// "grant_type=password&username=u&password=p" sent to /oauth/token.
	secretFormRE = regexp.MustCompile(
		`(?i)\b(password|client_secret|access_token|refresh_token)=[^&\r\n]*`,
	)

	// Secrets in JSON bodies, e.g. the token response
	// {"access_token":"...","refresh_token":"..."}.
	secretJSONRE = regexp.MustCompile(
		`(?i)"(password|clientSecret|client_secret|accessToken|access_token` +
			`|refreshToken|refresh_token)"(\s*:\s*)"[^"]*"`,
	)

	// showSecretsWarning ensures the clear-text warning is emitted once per
	// process rather than on every request.
	showSecretsWarning sync.Once
)

// IsEnabled reports whether API tracing is turned on.
func IsEnabled() bool {
	_, enabled := os.LookupEnv(EnvEnabled)

	return enabled
}

// ShowSecrets reports whether credential redaction has been switched off.
//
// Redaction is on by default so that traces captured by CI can be published
// safely. Set MORPHEUS_API_HTTPTRACE_SHOW_SECRETS to an affirmative value
// (e.g. "1" or "true") to emit raw traces when debugging an authentication
// problem locally.
//
// Unlike IsEnabled, which only tests for presence, an empty or explicitly
// falsey value keeps redaction on. This is a security control, so an
// accidentally blank export -- easily produced when a value is threaded
// through CI tooling -- must not silently expose credentials.
func ShowSecrets() bool {
	value, ok := os.LookupEnv(EnvShowSecrets)
	if !ok {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "0", "false", "no", "off":
		return false
	default:
		return true
	}
}

// Redact removes credentials from a dumped HTTP request or response so that
// traces can be published safely.
//
// Tracing is typically enabled on CI runners whose `go test -json` output is
// published as a build artifact. Because the dumps include bodies, an
// appliance password (sent to /oauth/token when the provider authenticates
// with username/password) and the resulting bearer tokens would otherwise be
// written to that artifact in clear text.
//
// Redaction is skipped entirely when ShowSecrets reports true.
func Redact(dump string) string {
	if ShowSecrets() {
		return dump
	}

	dump = secretHeaderRE.ReplaceAllString(dump, "$1: "+redactedValue)
	dump = secretFormRE.ReplaceAllString(dump, "$1="+redactedValue)
	dump = secretJSONRE.ReplaceAllString(dump, `"$1"$2"`+redactedValue+`"`)

	return dump
}

func New(
	transport http.RoundTripper,
) http.RoundTripper {
	return TraceRoundTripper{
		Transport: transport,
	}
}

type TraceRoundTripper struct {
	Transport http.RoundTripper
}

func (t TraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if ShowSecrets() {
		showSecretsWarning.Do(func() {
			tflog.Warn(req.Context(),
				EnvShowSecrets+" is set: API traces include credentials in "+
					"clear text. Do not publish this log.")
		})
	}

	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		tflog.Info(req.Context(),
			"\n\n->TX\n\n"+
				Redact(string(reqBytes))+
				"--\n")
	} else {
		msg := fmt.Sprintf("Error tracing request: %v", err)
		tflog.Error(req.Context(), msg)
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	respBytes, err := httputil.DumpResponse(resp, true)
	if err == nil {
		tflog.Info(req.Context(),
			"\n\n<-RX\n\n"+
				Redact(string(respBytes))+
				"\n--\n")
	} else {
		msg := fmt.Sprintf("Error tracing response: %v", err)
		tflog.Error(req.Context(), msg)
	}

	return resp, nil
}
