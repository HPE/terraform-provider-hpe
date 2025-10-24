// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/auth"
	sdklegacy "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

// A constructor for a new "legacy" Morpheus SDK client with custom
// RoundTrippers for per-request token and creds authentication.
// We use this for resources inherited from the legacy Morpheus provider,
// and to cover gaps in generated SDK implementations.
func NewLegacyClient(
	_ context.Context,
	url string,
	username string,
	password string,
	token string,
	opts ...sdklegacy.ClientOption,
) *sdklegacy.Client {
	c := sdklegacy.NewClient(url, opts...)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.IsInsecure(),
		},
		Proxy: http.ProxyFromEnvironment,
	}

	var authRoundTripper http.RoundTripper
	if token != "" {
		authRoundTripper = auth.NewTokenRoundTripper(
			context.Background(),
			transport,
			token,
		)
	} else {
		authRoundTripper = auth.NewCredsRoundTripper(
			context.Background(),
			transport,
			url,
			username,
			password,
		)
	}

	c.HTTPClient = &http.Client{
		Transport: authRoundTripper,
		Timeout:   15 * time.Second, // increased timeout for complex API endpoints
	}

	return c
}
