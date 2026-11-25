// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
	"github.com/HPE/terraform-provider-hpe/utils/httptrace"
)

// factory options
type FactoryOption func(*ClientFactory)

func WithFactoryHTTPClient(c *http.Client) FactoryOption {
	return func(cf *ClientFactory) {
		cf.httpclient = c
	}
}

func New(m model.MorpheusProviderModel, opts ...FactoryOption) *ClientFactory {
	var options []ClientOption

	cf := &ClientFactory{
		model: m,
	}

	for _, opt := range opts {
		opt(cf)
	}

	if cf.httpclient != nil {
		// Custom http client
		options = append(options, WithHTTPClient(cf.httpclient))
	}

	if cf.model.Insecure.ValueBool() {
		options = append(options, WithInsecureTLS())
	}

	f := func(ctx context.Context) (*sdk.APIClient, error) {
		client := NewAPIClient(
			ctx,
			cf.model.URL.ValueString(),
			cf.model.Username.ValueString(),
			cf.model.Password.ValueString(),
			cf.model.TenantSubdomain.ValueString(),
			cf.model.AccessToken.ValueString(),
			options...,
		)

		return client, nil
	}

	cf.newClient = f

	return cf
}

type ClientFactory struct {
	httpclient *http.Client
	model      model.MorpheusProviderModel
	newClient  func(context.Context) (*sdk.APIClient, error)
}

func (c ClientFactory) NewClient(ctx context.Context) (*sdk.APIClient, error) {
	// Surface error here to avoid panic on missing Morpheus config.
	// Checking for nil prevents a panic when using the provider adapter architecture.
	if c.newClient == nil {
		msg := `
morpheus client not configured - possible missing morpheus provider block.

provider "hpe" {
  morpheus { <- missing or duplicate?
    url = "https://example.com"
  }
}`

		return nil, errors.New(msg)
	}

	return c.newClient(ctx)
}

// Unauthenticated is an HTTP client for the few appliance endpoints that take
// no credentials (the first-time setup API), built from the same url and
// insecure settings as the authenticated SDK client.
type Unauthenticated struct {
	// BaseURL is the provider's Morpheus url without a trailing slash.
	BaseURL string
	// HTTP carries the provider's TLS settings and tracing but no
	// authentication; it has no timeout, so callers bound requests by context.
	HTTP *http.Client
}

// NewUnauthenticated returns an Unauthenticated client for the configured
// appliance. It fails the same way NewClient does when the provider block is
// missing.
func (c ClientFactory) NewUnauthenticated() (*Unauthenticated, error) {
	if c.newClient == nil {
		return nil, errors.New("morpheus client not configured - possible missing morpheus provider block")
	}

	u := &Unauthenticated{BaseURL: auth.NormalizeBaseURL(c.model.URL.ValueString())}

	u.HTTP = &http.Client{Transport: newTransport(c.model.Insecure.ValueBool())}
	if c.httpclient != nil {
		client := *c.httpclient
		u.HTTP = &client
	}
	u.HTTP.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return u, nil
}

// newTransport builds the provider's HTTP transport: TLS verification as
// configured, proxy from the environment, tracing when enabled, and the
// in-process observer hook tests use.
func newTransport(insecure bool) http.RoundTripper {
	var transport http.RoundTripper

	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure, //nolint: gosec
		},
		Proxy: http.ProxyFromEnvironment,
	}

	if httptrace.IsEnabled() {
		transport = httptrace.New(transport)
	}

	// Pass-through unless an in-process test installed an observer.
	return httptrace.ObserveRoundTripper{Transport: transport}
}

type clientOpts struct {
	httpclient *http.Client
	insecure   bool
}

// client options
type ClientOption func(*clientOpts)

func WithHTTPClient(h *http.Client) ClientOption {
	return func(o *clientOpts) {
		o.httpclient = h
	}
}

func WithInsecureTLS() ClientOption {
	return func(o *clientOpts) {
		o.insecure = true
	}
}

func NewAPIClient(
	_ context.Context,
	url,
	username string,
	password string,
	tenantSubdomain string,
	token string,
	opts ...ClientOption,
) *sdk.APIClient {
	var options clientOpts

	url = auth.NormalizeBaseURL(url)
	morpheusCfg := sdk.NewConfiguration()
	morpheusCfg.Servers[0].URL = url

	c := sdk.NewAPIClient(morpheusCfg)

	for _, opt := range opts {
		opt(&options)
	}

	if options.httpclient == nil {
		transport := newTransport(options.insecure)

		var authRoundTripper http.RoundTripper
		if token != "" {
			authRoundTripper = auth.NewTokenRoundTripper(
				context.Background(),
				transport,
				token,
			)
		} else {
			if tenantSubdomain != "" {
				username = fmt.Sprintf(`%s\\%s`, tenantSubdomain, username)
			}
			authRoundTripper = auth.NewCredsRoundTripper(
				context.Background(),
				transport,
				url,
				username,
				password,
			)
		}
		c.GetConfig().HTTPClient = &http.Client{
			Transport: authRoundTripper,
			Timeout:   15 * time.Minute,
		}
	}

	return c
}
