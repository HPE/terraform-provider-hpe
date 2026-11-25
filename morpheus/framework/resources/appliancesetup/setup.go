// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// The first-time setup wizard is a front end for two unauthenticated
// endpoints. GET /api/setup/check answers {"success":true} once the application
// layer is up and adds "setupNeeded":true while the appliance has no tenant;
// POST /api/setup/init creates the master tenant and its first user and is
// refused with HTTP 400 "Appliance has already been setup" afterwards.

const (
	checkPath = "/api/setup/check"
	initPath  = "/api/setup/init"

	// defaultPollInterval is how often create asks a starting appliance
	// whether its API is up yet. A first boot takes many minutes, so a few
	// seconds between polls costs nothing.
	defaultPollInterval = 15 * time.Second

	// maxBodyBytes bounds how much of an answer is read: the setup answers
	// are a few hundred bytes, and a loading page from the front-end proxy a
	// few kilobytes.
	maxBodyBytes = 1 << 20
)

// checkResult is the appliance's answer to GET /api/setup/check.
type checkResult struct {
	Success      bool   `json:"success"`
	SetupNeeded  bool   `json:"setupNeeded"`
	ApplianceURL string `json:"applianceUrl"`
	BuildVersion string `json:"buildVersion"`
}

// initPayload is the body of POST /api/setup/init. The field names are the
// appliance's. licenseKey is omitted when empty: an empty key fails decoding.
type initPayload struct {
	ApplianceURL  string `json:"applianceUrl"`
	ApplianceName string `json:"applianceName,omitempty"`
	AccountName   string `json:"accountName"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	HubMode       string `json:"hubMode"`
	LicenseKey    string `json:"licenseKey,omitempty"`
}

// initResult is the appliance's answer to POST /api/setup/init. On success it
// carries the check() fields as well.
type initResult struct {
	Success      bool              `json:"success"`
	Msg          string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	BuildVersion string            `json:"buildVersion"`
}

// setupClient talks to the setup endpoints of one appliance.
type setupClient struct {
	baseURL      string
	http         *http.Client
	pollInterval time.Duration
}

func newSetupClient(u *clientfactory.Unauthenticated) *setupClient {
	return &setupClient{baseURL: u.BaseURL, http: u.HTTP, pollInterval: defaultPollInterval}
}

// errNotReady is returned by check while the appliance is still starting:
// the connection is refused or dropped, the front-end proxy answers with a
// loading page, or the answer is not the setup API's JSON.
var errNotReady = errors.New("appliance API not ready")

// errAlreadySetUp is returned by init when another party initialised the
// appliance between check and init.
var errAlreadySetUp = errors.New("appliance has already been set up")

// check asks the appliance whether its API is up and whether setup is needed.
// A connection failure or an answer that is not the setup API's JSON is
// reported as errNotReady; the caller decides whether to wait.
func (c *setupClient) check(ctx context.Context) (checkResult, error) {
	var result checkResult

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+checkPath, nil)
	if err != nil {
		return result, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		return result, fmt.Errorf("%w: %s", errNotReady, redactedErr(err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return result, fmt.Errorf("%w: reading the answer: %s", errNotReady, redactedErr(err))
	}

	if resp.StatusCode != http.StatusOK || json.Unmarshal(body, &result) != nil || !result.Success {
		return result, fmt.Errorf("%w: HTTP %d: %s", errNotReady, resp.StatusCode, excerpt(body))
	}

	return result, nil
}

// awaitAPI polls check until the appliance answers, the context expires, or a
// non-readiness error occurs. It returns the first good answer.
func (c *setupClient) awaitAPI(ctx context.Context) (checkResult, error) {
	last := "none"

	for {
		result, err := c.check(ctx)
		if err == nil {
			return result, nil
		}

		if ctx.Err() != nil {
			return result, fmt.Errorf("timed out waiting for the appliance API to answer (last answer: %s)", last)
		}

		if !errors.Is(err, errNotReady) {
			return result, err
		}

		last = strings.TrimPrefix(err.Error(), errNotReady.Error()+": ")
		tflog.Info(ctx, "waiting for the appliance API", map[string]any{"last_answer": last})

		select {
		case <-ctx.Done():
			return result, fmt.Errorf("timed out waiting for the appliance API to answer (last answer: %s)", last)
		case <-time.After(c.pollInterval):
		}
	}
}

// init performs the first-time setup. The appliance's refusal of an already
// initialised appliance is errAlreadySetUp; any other refusal is an error
// carrying the appliance's message and per-field errors.
func (c *setupClient) init(ctx context.Context, payload initPayload) (initResult, error) {
	var result initResult

	// The payload carries the administrator's password: that is what the
	// setup endpoint takes, and the body goes nowhere but to the appliance.
	//nolint:gosec
	body, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+initPath, bytes.NewReader(body))
	if err != nil {
		return result, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return result, fmt.Errorf("posting the setup request: %s", redactedErr(err))
	}
	defer resp.Body.Close()

	answer, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return result, fmt.Errorf("reading the setup answer: %s", redactedErr(err))
	}

	if err := json.Unmarshal(answer, &result); err != nil {
		return result, fmt.Errorf("HTTP %d: %s", resp.StatusCode, excerpt(answer))
	}

	if result.Success {
		return result, nil
	}

	if strings.Contains(strings.ToLower(result.Msg), "already been setup") {
		return result, errAlreadySetUp
	}

	return result, fmt.Errorf("the appliance refused the setup (HTTP %d): %s", resp.StatusCode, result.message())
}

// message renders a refusal: the appliance's message followed by its
// per-field errors in a stable order.
func (r initResult) message() string {
	msg := r.Msg
	if msg == "" {
		msg = "no message"
	}

	if len(r.Errors) == 0 {
		return msg
	}

	fields := make([]string, 0, len(r.Errors))
	for field := range r.Errors {
		fields = append(fields, field)
	}

	sort.Strings(fields)

	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, field+": "+r.Errors[field])
	}

	return msg + " (" + strings.Join(parts, "; ") + ")"
}

// excerpt renders the start of a non-JSON answer (a proxy's loading page,
// say) on one line for a diagnostic.
func excerpt(body []byte) string {
	const limit = 120

	text := strings.Join(strings.Fields(string(body)), " ")
	if text == "" {
		return "empty body"
	}

	if len(text) > limit {
		return text[:limit] + "..."
	}

	return text
}

// redactedErr renders a transport error without the request URL's query or
// user info, which net/url errors otherwise include.
func redactedErr(err error) string {
	var urlErr interface{ Unwrap() error }
	if errors.As(err, &urlErr) && urlErr.Unwrap() != nil {
		return urlErr.Unwrap().Error()
	}

	return err.Error()
}
