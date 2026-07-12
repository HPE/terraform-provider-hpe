// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

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

// ImportID holds the parsed result of a standardized import ID.
type ImportID struct {
	// Client is the resolved client value: non-empty when MSP + client_id provided, empty otherwise.
	Client types.String
	// Parts contains the resource-specific ID segments (everything after the optional client prefix).
	Parts []string
}

// ParseImportID parses an import ID string with a standardized convention:
//
//	Single resource ID:                          <resource_id>
//	With client override (MSP only):             <client_id>:<resource_id>
//	Multi-segment (e.g. parent:child):           <parent_id>:<child_id>
//	Multi-segment with client (MSP only):        <client_id>:<parent_id>:<child_id>
//
// The expectedParts parameter specifies how many resource-specific segments are expected
// (not counting the optional client prefix). For example:
//   - A simple resource (just an ID): expectedParts = 1
//   - integration_event (integration_id:event_id): expectedParts = 2
//   - integration_config (integration_id:config_id): expectedParts = 2
//
// Scope rules:
//   - If the provider is CLIENT-scoped, the client prefix is never expected. All parts are resource segments.
//   - If the provider is MSP-scoped and len(parts) == expectedParts+1, the first part is treated as client_id.
//   - If the provider is MSP-scoped and len(parts) == expectedParts, no client override (resource at MSP level).
func (r *BaseResource) ParseImportID(importID string, expectedParts int) (*ImportID, error) {
	parts := strings.Split(importID, ":")

	result := &ImportID{
		Client: types.StringNull(),
	}

	if r.apiClient.Scope == "CLIENT" {
		// CLIENT scope: all parts are resource segments, ignore any client prefix interpretation
		if len(parts) != expectedParts {
			return nil, importIDError(expectedParts, false)
		}
		result.Parts = parts
		return result, nil
	}

	// MSP scope: optionally has a client prefix
	switch len(parts) {
	case expectedParts:
		// No client prefix — resource at MSP level
		result.Parts = parts
	case expectedParts + 1:
		// First segment is client_id
		result.Client = types.StringValue(parts[0])
		result.Parts = parts[1:]
	default:
		return nil, importIDError(expectedParts, true)
	}

	return result, nil
}

// TenantForImport returns the tenant ID to use for API calls during import.
// If the import specified a client override, returns that; otherwise returns the provider's tenant.
func (r *BaseResource) TenantForImport(parsed *ImportID) string {
	if !parsed.Client.IsNull() && parsed.Client.ValueString() != "" {
		return parsed.Client.ValueString()
	}
	return r.apiClient.TenantId
}

func importIDError(expectedParts int, isMSP bool) error {
	if expectedParts == 1 {
		if isMSP {
			return fmt.Errorf("expected format: <resource_id> or <client_id>:<resource_id>")
		}
		return fmt.Errorf("expected format: <resource_id>")
	}
	// Build a generic multi-segment description
	segments := make([]string, expectedParts)
	for i := range segments {
		segments[i] = "<id>"
	}
	base := strings.Join(segments, ":")
	if isMSP {
		return fmt.Errorf("expected format: %s or <client_id>:%s", base, base)
	}
	return fmt.Errorf("expected format: %s", base)
}
