// (C) Copyright 2024-2026 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	consts "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/common"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/models"
)

const (
	testServiceInstanceID    = "18ba6409-ac59-4eac-9414-0147e72d615e"
	testAccessToken          = "2b9fba7f-7c14-4773-a970-a9ad393811ac" //nolint:gosec // test fixture
	testRefreshToken         = "7806acfb-f847-48b1-a6d5-6119dccb3ffe" //nolint:gosec // test fixture
	testMorpheusURL          = "https://1234-mp.private.greenlake.hpe-gl-intg.com"
	testIAMToken             = "iam-token-fixture" //nolint:gosec // test fixture
	testAccessTokenExpires   = 1758034360176
	testAccessTokenExpiresIn = 3600
)

const cmpDetailsBody = `{
	"ServiceInstanceID": "` + testServiceInstanceID + `",
	"TenantID": "1234",
	"TenantName": "tenant",
	"LocationName": "BLR",
	"URL": "` + testMorpheusURL + `/",
	"TokenDetails": {
		"access_token": "` + testAccessToken + `",
		"expires": 1758034360176,
		"refresh_token": "` + testRefreshToken + `",
		"expires_in": 3600
	}
}`

// newTestBrokerClient wires an APIClient at the supplied host with the
// location/space query params and a bearer token, mirroring how the
// pce_identity token exchange configures the broker client.
func newTestBrokerClient(host string) *APIClient {
	cfg := NewConfiguration()
	cfg.Host = host
	cfg.DefaultQueryParams["location"] = "BLR"
	cfg.DefaultQueryParams["space"] = "default"

	c := NewAPIClient(cfg)
	c.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, _ interface{}) {
		*ctx = context.WithValue(*ctx, ContextAccessToken, testIAMToken)
	})

	return c
}

func TestGetMorpheusDetails_Success(t *testing.T) {
	t.Parallel()

	var (
		gotPath   string
		gotQuery  url.Values
		gotAuth   string
		gotAccept string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")

		w.Header().Set("Content-Type", consts.ContentType)
		_, _ = w.Write([]byte(cmpDetailsBody))
	}))
	defer srv.Close()

	got, err := newTestBrokerClient(srv.URL).GetCMPDetails(context.Background())
	if err != nil {
		t.Fatalf("GetCMPDetails() unexpected error: %v", err)
	}

	// The broker path must not carry the vmaas-cmp API base path.
	if want := "/" + consts.CMPDetails; gotPath != want {
		t.Errorf("request path = %q, want %q", gotPath, want)
	}

	if v := gotQuery.Get("location"); v != "BLR" {
		t.Errorf("location query = %q, want %q", v, "BLR")
	}

	if v := gotQuery.Get("space"); v != "default" {
		t.Errorf("space query = %q, want %q", v, "default")
	}

	if want := "Bearer " + testIAMToken; gotAuth != want {
		t.Errorf("Authorization header = %q, want %q", gotAuth, want)
	}

	if gotAccept != consts.ContentType {
		t.Errorf("Accept header = %q, want %q", gotAccept, consts.ContentType)
	}

	// URL is returned with any trailing slash trimmed.
	want := models.TFMorpheusDetails{
		ID:          testServiceInstanceID,
		AccessToken: testAccessToken,
		ValidTill:   testAccessTokenExpiresIn,
		URL:         testMorpheusURL,
	}

	if got != want {
		t.Errorf("GetCMPDetails() = %+v, want %+v", got, want)
	}
}

func TestGetMorpheusDetails_ErrorStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", consts.ContentType)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer srv.Close()

	got, err := newTestBrokerClient(srv.URL).GetCMPDetails(context.Background())
	if err == nil {
		t.Fatal("GetCMPDetails() expected an error for a 401 response, got nil")
	}

	if got != (models.TFMorpheusDetails{}) {
		t.Errorf("GetCMPDetails() = %+v, want zero value on error", got)
	}
}

func TestGetMorpheusDetails_MalformedBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", consts.ContentType)
		_, _ = w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	if _, err := newTestBrokerClient(srv.URL).GetCMPDetails(context.Background()); err == nil {
		t.Fatal("GetCMPDetails() expected an error for a malformed body, got nil")
	}
}

// trackingBody records whether it was closed, so a test can prove the client
// releases the response body.
type trackingBody struct {
	io.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true

	return nil
}

// roundTripperFunc adapts a function to http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// newBodyTrackingClient returns a broker client whose transport serves the
// given status and body, along with the body it will serve so that the test can
// check whether it was closed.
func newBodyTrackingClient(status int, payload string) (*APIClient, *trackingBody) {
	body := &trackingBody{Reader: strings.NewReader(payload)}

	cfg := NewConfiguration()
	cfg.Host = "https://broker.example.invalid"
	cfg.HTTPClient = &http.Client{
		Transport: roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Body:       body,
				Header:     make(http.Header),
			}, nil
		}),
	}

	c := NewAPIClient(cfg)
	c.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, _ interface{}) {
		*ctx = context.WithValue(*ctx, ContextAccessToken, testIAMToken)
	})

	return c, body
}

// An error response must not leak its body. ParseError reads the body without
// closing it, so the close has to be deferred before the status is checked.
func TestGetCMPDetailsClosesBodyOnErrorResponse(t *testing.T) {
	t.Parallel()

	client, body := newBodyTrackingClient(
		http.StatusUnauthorized,
		`{"message":"unauthorized"}`,
	)

	if _, err := client.GetCMPDetails(context.Background()); err == nil {
		t.Fatal("GetCMPDetails() expected an error for a 401 response, got nil")
	}

	if !body.closed {
		t.Error("response body was not closed on an error response")
	}
}

// The success path must release the body too.
func TestGetCMPDetailsClosesBodyOnSuccess(t *testing.T) {
	t.Parallel()

	client, body := newBodyTrackingClient(http.StatusOK, cmpDetailsBody)

	if _, err := client.GetCMPDetails(context.Background()); err != nil {
		t.Fatalf("GetCMPDetails() unexpected error: %v", err)
	}

	if !body.closed {
		t.Error("response body was not closed on a successful response")
	}
}

// An unconfigured api must report an error rather than panicking: a panic in a
// provider SDK takes down the whole plugin process.
func TestAPIDoRejectsUnconfiguredAPI(t *testing.T) {
	t.Parallel()

	client, _ := newBodyTrackingClient(http.StatusOK, cmpDetailsBody)

	testcases := map[string]*api{
		"no path": {
			method:     http.MethodGet,
			client:     client,
			jsonParser: func([]byte) error { return nil },
		},
		"no method": {
			path:       consts.CMPDetails,
			client:     client,
			jsonParser: func([]byte) error { return nil },
		},
		"no client": {
			path:       consts.CMPDetails,
			method:     http.MethodGet,
			jsonParser: func([]byte) error { return nil },
		},
		"no json parser": {
			path:   consts.CMPDetails,
			method: http.MethodGet,
			client: client,
		},
	}

	for name, testcase := range testcases {
		a := testcase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// The nil client case would panic on a nil dereference if the
			// configuration were not checked before the version lookup.
			err := a.do(context.Background(), nil, nil)
			if err == nil {
				t.Fatal("do() returned no error for an unconfigured api")
			}

			if got, want := err.Error(), "api not properly configured"; got != want {
				t.Errorf("do() error = %q, want %q", got, want)
			}
		})
	}
}

// A version that cannot be parsed must also be reported rather than panicking.
func TestAPIDoRejectsUnparseableVersion(t *testing.T) {
	t.Parallel()

	client, _ := newBodyTrackingClient(http.StatusOK, cmpDetailsBody)

	a := &api{
		path:              consts.CMPDetails,
		method:            http.MethodGet,
		client:            client,
		compatibleVersion: "not-a-version",
		jsonParser:        func([]byte) error { return nil },
	}

	err := a.do(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("do() returned no error for an unparseable compatible version")
	}

	if !strings.Contains(err.Error(), "failed to parse the compatible version") {
		t.Errorf("do() error = %q, want it to mention the compatible version", err.Error())
	}
}
