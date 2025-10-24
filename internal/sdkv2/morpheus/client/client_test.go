// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sdklegacy "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/stretchr/testify/assert"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/client"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that the legacy client we construct for the provider
// can correctly use custom transports
func TestLegacyClientCustomTransport(t *testing.T) {
	defer testhelpers.RecordResult(t)

	mux := http.NewServeMux()

	// NOTE: The first write to a response implicitly sets status to 200 OK
	// https://pkg.go.dev/net/http#ResponseWriter
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		// simulate granting token (for creds)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"abc"}`)
	})

	mux.HandleFunc("/api/whoami", func(w http.ResponseWriter, r *http.Request) {
		// simulate checking token
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		// if we don't have valid JSON the SDK will fail to parse response
		fmt.Fprint(w, `{}`)
	})

	server := httptest.NewServer(mux)

	defer server.Close()

	testCases := []struct {
		name      string
		url       string
		username  string
		password  string
		token     string
		expectErr bool
	}{
		{
			name:      "credentials",
			url:       server.URL,
			username:  "test",
			password:  "test123",
			expectErr: false,
		},
		{
			name:      "token",
			url:       server.URL,
			token:     "abc",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			morpheusClient := client.NewLegacyClient(
				context.Background(),
				tc.url,
				tc.username,
				tc.password,
				tc.token,
				sdklegacy.SkipLogin(),
			)

			resp, err := morpheusClient.Whoami()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
