// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"net/http"
	"time"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/auth"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/httptrace"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HewlettPackard/hpe-morpheus-client/client"
)

// factory options
type FactoryOption func(*ClientFactory)

func WithFactoryHTTPClient(c *http.Client) FactoryOption {
	return func(cf *ClientFactory) {
		cf.httpclient = c
	}
}

func New(m model.SubModel, opts ...FactoryOption) *ClientFactory {
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

	f := func(ctx context.Context) (*client.APIClient, error) {
		client := NewAPIClient(
			ctx,
			cf.model.URL.ValueString(),
			cf.model.Username.ValueString(),
			cf.model.Password.ValueString(),
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
	model      model.SubModel
	newClient  func(context.Context) (*client.APIClient, error)
}

func (c ClientFactory) NewClient(ctx context.Context) (*client.APIClient, error) {
	return c.newClient(ctx)
}

// client options
type ClientOption func(*client.APIClient)

func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *client.APIClient) {
		c.GetConfig().HTTPClient = h
	}
}

func NewAPIClient(
	_ context.Context,
	url,
	username string,
	password string,
	token string,
	opts ...ClientOption,
) *client.APIClient {
	morpheusCfg := client.NewConfiguration()
	morpheusCfg.Servers[0].URL = url

	c := client.NewAPIClient(morpheusCfg)

	if c.GetConfig().HTTPClient == http.DefaultClient {
		var authTransport http.RoundTripper
		if token != "" {
			authTransport = auth.NewTokenRoundTripper(
				context.Background(),
				token,
			)
		} else {
			authTransport = auth.NewCredsRoundTripper(
				context.Background(),
				url,
				username,
				password,
			)
		}
		c.GetConfig().HTTPClient = &http.Client{
			Transport: authTransport,
			Timeout:   15 * time.Second,
		}
	}

	if httptrace.Enabled() {
		c.GetConfig().HTTPClient.Transport = httptrace.New(
			c.GetConfig().HTTPClient.Transport,
		)
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
