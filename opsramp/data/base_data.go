// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/opsramp/utils/clientfactory"
)

// Base Data struct
type BaseData struct {
	clientFactory *clientfactory.ClientFactory
	apiClient     *client.OpsRampClient
}

// Configure prepares the data by resolving the API client from the factory.
// Data Configure is called on each Terraform operation, so transient failures
// are retried automatically on the next plan/apply.
func (d *BaseData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	f, ok := req.ProviderData.(*clientfactory.ClientFactory)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.ClientFactory",
		)

		return
	}

	d.clientFactory = f

	c, err := f.Client()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpsRamp API Client",
			err.Error(),
		)

		return
	}

	d.apiClient = c
}
