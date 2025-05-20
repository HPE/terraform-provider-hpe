// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/auth"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/httptrace"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
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

	if cf.model.Insecure.ValueBool() {
		options = append(options, WithInsecureTLS())
	}

	f := func(ctx context.Context) (*sdk.APIClient, error) {
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
	newClient  func(context.Context) (*sdk.APIClient, error)
}

func (c ClientFactory) NewClient(ctx context.Context) (*sdk.APIClient, error) {
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
	token string,
	opts ...ClientOption,
) *sdk.APIClient {
	var options clientOpts

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
			Timeout:   15 * time.Second,
		}
	}

	return c
}
