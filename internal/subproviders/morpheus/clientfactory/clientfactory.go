// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
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

func (c ClientFactory) New(_ context.Context) client.Client {
	httpClient := &http.Client{
		Transport: c.transport,
	}

	return client.Client{
		HTTPClient: httpClient,
		URL:        c.model.URL.ValueString(),
	}
}
