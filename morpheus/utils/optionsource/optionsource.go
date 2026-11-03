// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package optionsource reads Morpheus option sources — the enumerations behind
// the UI's dropdowns, served at GET /api/options/{source} and
// GET /api/options/{serviceType}/{source} — for the hpe_morpheus_option_source_value
// and hpe_morpheus_option_source_values data sources.
//
// The generated SDK cannot be used for this: its method for the endpoint takes
// no query parameters, and the parameter space is open (each source reads what
// it likes). So the request is built here and sent through the SDK client's own
// HTTP client, which carries the provider's authentication, TLS settings and
// timeouts.
//
// Two server behaviours shape the design and are documented on the types below:
// a parameter sent in the wrong position is silently ignored rather than
// rejected, so every parameter is sent in both positions by default; and the
// response shape is whatever the source returns, so the result model tolerates
// arrays and objects, header rows, and values keyed "id" instead of "value".
package optionsource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Style controls where each parameter is placed in the request.
//
// Morpheus reads option-source parameters from two places: the top-level query
// string, and a config[<key>] query parameter. Which one a given source reads is
// not discoverable, and a parameter in the wrong place is not an error — the
// source silently runs unfiltered. StyleAuto therefore sends every parameter in
// both places, which every source tested tolerates. The other two styles exist
// for a source that does not.
type Style string

const (
	// StyleAuto sends each parameter both top-level and as config[key].
	StyleAuto Style = "auto"
	// StyleTopLevel sends each parameter top-level only.
	StyleTopLevel Style = "top_level"
	// StyleConfig sends each parameter as config[key] only.
	StyleConfig Style = "config"
)

// Styles lists the accepted parameter_style values, for schema validation.
func Styles() []string {
	return []string{string(StyleAuto), string(StyleTopLevel), string(StyleConfig)}
}

// ParseStyle converts the configured string to a Style; "" means StyleAuto.
func ParseStyle(s string) (Style, error) {
	switch Style(s) {
	case "", StyleAuto:
		return StyleAuto, nil
	case StyleTopLevel, StyleConfig:
		return Style(s), nil
	default:
		return "", fmt.Errorf("unknown parameter_style %q; expected one of %s", s, strings.Join(Styles(), ", "))
	}
}

// Request describes one option-source call.
type Request struct {
	// Source is the option-source name, e.g. "zonePools".
	Source string
	// ServiceType, when set, selects the two-segment path
	// /api/options/{serviceType}/{source}, which is how a source defined on a
	// cloud-specific service (e.g. "scvmm") is reached.
	ServiceType string
	// Parameters are sent according to Style. Keys are the names the source
	// reads (zoneId, layoutId, ...).
	Parameters map[string]string
	// Style is where the parameters are placed; zero value means StyleAuto.
	Style Style
	// Collection selects a list from an object-shaped response; see parseData.
	Collection string
}

// Path returns the request path for the source.
func (r Request) Path() string {
	if r.ServiceType != "" {
		return "/api/options/" + url.PathEscape(r.ServiceType) + "/" + url.PathEscape(r.Source)
	}

	return "/api/options/" + url.PathEscape(r.Source)
}

// Query returns the encoded query string for the parameters, per Style. Keys are
// sorted so the result is stable for messages and tests. The config[key] form is
// percent-encoded — the server rejects raw brackets.
func (r Request) Query() string {
	style := r.Style
	if style == "" {
		style = StyleAuto
	}

	keys := make([]string, 0, len(r.Parameters))
	for k := range r.Parameters {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	values := url.Values{}

	for _, k := range keys {
		v := r.Parameters[k]

		if style == StyleAuto || style == StyleTopLevel {
			values.Set(k, v)
		}

		if style == StyleAuto || style == StyleConfig {
			values.Set("config["+k+"]", v)
		}
	}

	return values.Encode()
}

// Value is one option as the data sources expose it.
type Value struct {
	// Name is the display label.
	Name string
	// Value is the value to use in configuration — the source's "value", or its
	// "id" when there is no "value" (the zoneNetworkOptions shape).
	Value string
	// Attributes carries every other key of the option, stringified.
	Attributes map[string]string
}

// Result is a successful call.
type Result struct {
	// Values are the options, in the order returned, with header rows removed.
	Values []Value
	// Path and Query are what was sent, for diagnostics.
	Path  string
	Query string
}

// ServerError is a call the server refused, or answered with success:false.
type ServerError struct {
	Source      string
	ServiceType string
	Path        string
	Query       string
	Status      int
	// Message is the server's "msg" when the body was the usual envelope, else
	// a trimmed excerpt of the body.
	Message string
}

// gasket is the text Morpheus returns for an unhandled failure inside the
// source. It is the same for an unknown source name and for a source given
// missing or malformed parameters, so the caller cannot tell them apart from
// the response alone; Advice adds what the catalogue knows.
const gasket = "threw a gasket"

// Gasket reports whether the failure is the generic unhandled one.
func (e *ServerError) Gasket() bool {
	return strings.Contains(e.Message, gasket)
}

// Describe is the failure without the request line, so a caller can insert
// what it knows between the two; Error is Describe followed by RequestLine.
func (e *ServerError) Describe() string {
	var b strings.Builder

	fmt.Fprintf(&b, "option source %q returned HTTP %d", e.Source, e.Status)

	if e.Message != "" {
		fmt.Fprintf(&b, " (%s)", e.Message)
	}

	b.WriteString(".")

	if e.Gasket() {
		b.WriteString(" Morpheus returns this for an unknown source name and for a source " +
			"whose parameters are missing or malformed alike.")
	}

	return b.String()
}

// RequestLine names the request that failed.
func (e *ServerError) RequestLine() string {
	if e.Query != "" {
		return fmt.Sprintf("Request: GET %s?%s", e.Path, e.Query)
	}

	return "Request: GET " + e.Path
}

func (e *ServerError) Error() string {
	return e.Describe() + " " + e.RequestLine()
}

// CollectionError is an object-shaped response read without a usable
// Collection: none was given, or the one given is not a list in the response.
type CollectionError struct {
	Source     string
	Collection string
	Available  []string
}

func (e *CollectionError) Error() string {
	if e.Collection == "" {
		return fmt.Sprintf("option source %q returns an object, not a list; set collection to one of: %s",
			e.Source, strings.Join(e.Available, ", "))
	}

	return fmt.Sprintf("option source %q has no list named %q; the lists available are: %s",
		e.Source, e.Collection, strings.Join(e.Available, ", "))
}

// ShapeError is a Collection given for a response that is already a list.
type ShapeError struct {
	Source     string
	Collection string
}

func (e *ShapeError) Error() string {
	return fmt.Sprintf("option source %q returns a list, so collection %q does not apply; remove it",
		e.Source, e.Collection)
}

// maxBody bounds the response read. The largest source observed returns a few
// hundred rows; this is generous.
const maxBody = 32 << 20

// Fetch calls the option source and returns its values.
//
// The SDK client supplies the appliance URL and an HTTP client whose transport
// adds the bearer token, so the call is authenticated like every other.
func Fetch(ctx context.Context, client *sdk.APIClient, req Request) (*Result, error) {
	if client == nil {
		return nil, errors.New("option source: nil client")
	}

	if strings.TrimSpace(req.Source) == "" {
		return nil, errors.New("option source: source is required")
	}

	cfg := client.GetConfig()
	if cfg == nil || len(cfg.Servers) == 0 || cfg.HTTPClient == nil {
		return nil, errors.New("option source: client has no server URL or HTTP client")
	}

	path := req.Path()
	query := req.Query()

	full := strings.TrimRight(cfg.Servers[0].URL, "/") + path
	if query != "" {
		full += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, full, nil)
	if err != nil {
		return nil, fmt.Errorf("option source %q: building request: %w", req.Source, err)
	}

	httpReq.Header.Set("Accept", "application/json")

	for k, v := range cfg.DefaultHeader {
		httpReq.Header.Set(k, v)
	}

	resp, err := cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("option source %q: GET %s: %w", req.Source, path, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return nil, fmt.Errorf("option source %q: reading response: %w", req.Source, err)
	}

	env := envelope{}
	envelopeOK := json.Unmarshal(body, &env) == nil

	if resp.StatusCode != http.StatusOK || (envelopeOK && env.Success != nil && !*env.Success) {
		msg := env.Msg
		if msg == "" {
			msg = excerpt(body)
		}

		return nil, &ServerError{
			Source: req.Source, ServiceType: req.ServiceType,
			Path: path, Query: query, Status: resp.StatusCode, Message: msg,
		}
	}

	if !envelopeOK {
		return nil, fmt.Errorf("option source %q: response is not JSON: %s", req.Source, excerpt(body))
	}

	values, err := parseData(req.Source, env.Data, req.Collection)
	if err != nil {
		return nil, err
	}

	return &Result{Values: values, Path: path, Query: query}, nil
}

// envelope is the response wrapper every option source shares.
type envelope struct {
	Success *bool           `json:"success"`
	Msg     string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

// excerpt trims a body for a message.
func excerpt(body []byte) string {
	s := strings.TrimSpace(string(body))
	if len(s) > 200 {
		s = s[:200] + "…"
	}

	return s
}
