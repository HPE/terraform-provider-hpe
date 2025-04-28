// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/client"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/httptrace"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
)

func New(m model.SubModel) *ClientFactory {
	t := http.DefaultTransport

	if httptrace.Enabled() {
		t = httptrace.New(t)
	}

	return &ClientFactory{
		transport: t,
		model:     m,
	}
}

type ClientFactory struct {
	transport http.RoundTripper
	model     model.SubModel
}

func (c ClientFactory) New(ctx context.Context) (*client.Client, error) {
	morpheus := client.New(ctx, client.Config{
		URL:         c.model.URL.ValueString(),
		InsecureTLS: true,
	})

	if !c.model.AccessToken.IsNull() {
		if err := morpheus.SetAccessToken(ctx, c.model.AccessToken.ValueString()); err != nil {
			return nil, fmt.Errorf("morpheus access token authentication failed: %w", err)
		}
	}

	if !c.model.Username.IsNull() {
		if err := morpheus.SetCredentials(ctx,
			c.model.Username.ValueString(),
			c.model.Password.ValueString(),
		); err != nil {
			return nil, fmt.Errorf("morpheus credential authentication failed: %w", err)
		}
	}

	return morpheus, nil
}
