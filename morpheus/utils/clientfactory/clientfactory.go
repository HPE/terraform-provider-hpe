// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/httptrace"
)

// tenancyCache lazily records whether the configured caller is the master
// tenant. It is held by pointer on ClientFactory so that by-value copies of the
// factory (see morpheus/configure) share the same cache. Only successful
// determinations are cached; errors are not, so a transient whoami failure does
// not poison later checks.
type tenancyCache struct {
	mu       sync.Mutex
	resolved bool
	isMaster bool
}

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
		model:   m,
		tenancy: &tenancyCache{},
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
	tenancy    *tenancyCache
}

// CallerIsMaster reports whether the authenticated tenant is the master tenant,
// via whoami. whoami introspects the current user (not an arbitrary account), so
// it is reachable by any authenticated caller including a subtenant.
//
// The whoami response omits isMasterAccount unless it is true, so a nil value
// means non-master, not "unknown".
//
// Successful determinations are cached on the pointer-held tenancy cache (shared
// across by-value copies of the factory). A nil cache (zero-value factory) falls
// back to an uncached lookup. Errors are never cached.
func (c ClientFactory) CallerIsMaster(ctx context.Context) (bool, error) {
	if c.tenancy != nil {
		c.tenancy.mu.Lock()
		if c.tenancy.resolved {
			isMaster := c.tenancy.isMaster
			c.tenancy.mu.Unlock()

			return isMaster, nil
		}
		c.tenancy.mu.Unlock()
	}

	client, err := c.NewClient(ctx)
	if err != nil {
		return false, err
	}

	who, hresp, err := client.AuthenticationAPI.Whoami(ctx).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		return false, err
	}
	if who == nil {
		return false, fmt.Errorf("whoami returned an empty response")
	}

	// isMasterAccount is only serialized when true; absence means non-master.
	isMaster := who.IsMasterAccount != nil && *who.IsMasterAccount

	if c.tenancy != nil {
		c.tenancy.mu.Lock()
		c.tenancy.resolved = true
		c.tenancy.isMaster = isMaster
		c.tenancy.mu.Unlock()
	}

	return isMaster, nil
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
		var transport http.RoundTripper

		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.insecure, //nolint: gosec
			},
			Proxy: http.ProxyFromEnvironment,
		}

		if httptrace.IsEnabled() {
			transport = httptrace.New(transport)
		}

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
