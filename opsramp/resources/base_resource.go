// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/opsramp/utils/clientfactory"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Base Resource struct with the shared ModifyPlan method
type BaseResource struct {
	clientFactory *clientfactory.ClientFactory
	apiClient     *client.OpsRampClient
}

// Configure prepares the resource by resolving the API client from the factory.
// Resource Configure is called on each Terraform operation, so transient failures
// are retried automatically on the next plan/apply.
func (r *BaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clientFactory = f

	c, err := f.Client()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpsRamp API Client",
			"An unexpected error occurred when creating the OpsRamp API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OpsRamp Client Error: "+err.Error(),
		)
		return
	}

	r.apiClient = c
}

func (r *BaseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Shared plan modification logic goes here
	if r.apiClient == nil {
		return
	}

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var clientVal types.String
	diag := req.Plan.GetAttribute(ctx, path.Root("client"), &clientVal)
	if diag.HasError() {
		// The resource schema does not have a "client" attribute; skip scope validation.
		return
	}

	// Validate client attribute is not used when provider is in Client scope
	if r.apiClient.Scope == "CLIENT" && !clientVal.IsNull() && clientVal.ValueString() != "" {
		resp.Diagnostics.AddError(
			"Client Attribute Not Allowed",
			"The 'client' attribute cannot be used when the provider is configured for a Client (non-MSP) tenant. "+
				"All resources will use the provider's configured tenant.",
		)
	}
}
