// (C) Copyright 2024-2026 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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
