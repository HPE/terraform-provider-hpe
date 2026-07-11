// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Base Resource struct with the shared ModifyPlan method
type BaseResource struct {
	apiClient *client.OpsRampClient
}

// Configure prepares the resource with the shared API client instance.
func (r *BaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.OpsRampClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.OpsRampClient",
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
