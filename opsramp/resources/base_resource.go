// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/opsramp/utils/clientfactory"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
		// No client prefix - resource at MSP level
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

// ValidateOrganizationMatching enforces the shared organization_matching_type / included_clients
// rules for policies that support MSP-only organization matching (alert correlation, alert
// prediction, and first response policies):
//   - At MSP level (no client override), organization_matching_type is required.
//   - When organization_matching_type is "INCLUDE", included_clients is required and non-empty.
//   - When organization_matching_type is "ALL", included_clients must not be configured.
//   - At CLIENT level, neither attribute may be configured; they are never sent to the API.
//
// included_clients entries (names or UUIDs) are passed through unchanged: the API accepts and
// echoes back whatever identifier form the user configures, so the plan must not rewrite them.
func (r *BaseResource) ValidateOrganizationMatching(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.apiClient == nil || req.Plan.Raw.IsNull() {
		return
	}

	orgPath := path.Root("organization_matching_type")
	clientsPath := path.Root("included_clients")

	var clientVal types.String
	if diags := req.Plan.GetAttribute(ctx, path.Root("client"), &clientVal); diags.HasError() {
		return
	}

	// An unknown "client" (e.g. referencing a resource created in the same apply) can't yet be
	// classified as MSP or CLIENT scope; defer validation until it resolves to a known value.
	if clientVal.IsUnknown() {
		return
	}

	var orgConfig types.String
	req.Config.GetAttribute(ctx, orgPath, &orgConfig)

	var clientsConfig types.List
	req.Config.GetAttribute(ctx, clientsPath, &clientsConfig)

	isMSP := r.apiClient.Scope == "MSP" && (clientVal.IsNull() || clientVal.ValueString() == "")

	if !isMSP {
		if !orgConfig.IsNull() && orgConfig.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(orgPath, "Not Applicable at Client Level",
				"'organization_matching_type' only applies when the operation is performed at the MSP level "+
					"and must not be set when operating on a specific client.")
		}

		if !clientsConfig.IsNull() {
			resp.Diagnostics.AddAttributeError(clientsPath, "Not Applicable at Client Level",
				"'included_clients' only applies when the operation is performed at the MSP level "+
					"and must not be set when operating on a specific client.")
		}

		if !resp.Diagnostics.HasError() {
			resp.Plan.SetAttribute(ctx, clientsPath, types.ListNull(types.StringType))
		}

		return
	}

	var orgPlan types.String
	req.Plan.GetAttribute(ctx, orgPath, &orgPlan)

	// An unknown "organization_matching_type" can't yet be checked for presence; defer
	// validation until it resolves to a known value.
	if orgPlan.IsUnknown() {
		return
	}

	if orgPlan.IsNull() || orgPlan.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(orgPath, "Missing Required Attribute",
			"'organization_matching_type' is required when the operation is performed at the MSP level.")

		return
	}

	if orgPlan.ValueString() != "INCLUDE" {
		if !clientsConfig.IsNull() {
			resp.Diagnostics.AddAttributeError(clientsPath, "Not Applicable",
				"'included_clients' only applies when 'organization_matching_type' is 'INCLUDE'.")

			return
		}

		resp.Plan.SetAttribute(ctx, clientsPath, types.ListNull(types.StringType))

		return
	}

	var clientsPlan types.List
	req.Plan.GetAttribute(ctx, clientsPath, &clientsPlan)

	// A wholly unknown "included_clients" (e.g. built from a resource created in the same
	// apply) can't yet be checked for emptiness; defer validation until it resolves.
	if clientsPlan.IsUnknown() {
		return
	}

	if clientsPlan.IsNull() || len(clientsPlan.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(clientsPath, "Missing Required Attribute",
			"'included_clients' is required and must include at least one client when "+
				"'organization_matching_type' is 'INCLUDE'.")

		return
	}
}

// ValidateOrganizationMatchingSet is the Set-based counterpart of ValidateOrganizationMatching,
// used by policies (such as alert escalation) whose included_clients attribute is a Set rather
// than a List.
func (r *BaseResource) ValidateOrganizationMatchingSet(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.apiClient == nil || req.Plan.Raw.IsNull() {
		return
	}

	orgPath := path.Root("organization_matching_type")
	clientsPath := path.Root("included_clients")

	var clientVal types.String
	if diags := req.Plan.GetAttribute(ctx, path.Root("client"), &clientVal); diags.HasError() {
		return
	}

	// An unknown "client" (e.g. referencing a resource created in the same apply) can't yet be
	// classified as MSP or CLIENT scope; defer validation until it resolves to a known value.
	if clientVal.IsUnknown() {
		return
	}

	var orgConfig types.String
	req.Config.GetAttribute(ctx, orgPath, &orgConfig)

	var clientsConfig types.Set
	req.Config.GetAttribute(ctx, clientsPath, &clientsConfig)

	isMSP := r.apiClient.Scope == "MSP" && (clientVal.IsNull() || clientVal.ValueString() == "")

	if !isMSP {
		if !orgConfig.IsNull() && orgConfig.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(orgPath, "Not Applicable at Client Level",
				"'organization_matching_type' only applies when the operation is performed at the MSP level "+
					"and must not be set when operating on a specific client.")
		}

		if !clientsConfig.IsNull() {
			resp.Diagnostics.AddAttributeError(clientsPath, "Not Applicable at Client Level",
				"'included_clients' only applies when the operation is performed at the MSP level "+
					"and must not be set when operating on a specific client.")
		}

		if !resp.Diagnostics.HasError() {
			resp.Plan.SetAttribute(ctx, clientsPath, types.SetNull(types.StringType))
		}

		return
	}

	var orgPlan types.String
	req.Plan.GetAttribute(ctx, orgPath, &orgPlan)

	// An unknown "organization_matching_type" can't yet be checked for presence; defer
	// validation until it resolves to a known value.
	if orgPlan.IsUnknown() {
		return
	}

	if orgPlan.IsNull() || orgPlan.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(orgPath, "Missing Required Attribute",
			"'organization_matching_type' is required when the operation is performed at the MSP level.")

		return
	}

	if orgPlan.ValueString() != "INCLUDE" {
		if !clientsConfig.IsNull() {
			resp.Diagnostics.AddAttributeError(clientsPath, "Not Applicable",
				"'included_clients' only applies when 'organization_matching_type' is 'INCLUDE'.")

			return
		}

		resp.Plan.SetAttribute(ctx, clientsPath, types.SetNull(types.StringType))

		return
	}

	var clientsPlan types.Set
	req.Plan.GetAttribute(ctx, clientsPath, &clientsPlan)

	// A wholly unknown "included_clients" (e.g. built from a resource created in the same
	// apply) can't yet be checked for emptiness; defer validation until it resolves.
	if clientsPlan.IsUnknown() {
		return
	}

	if clientsPlan.IsNull() || len(clientsPlan.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(clientsPath, "Missing Required Attribute",
			"'included_clients' is required and must include at least one client when "+
				"'organization_matching_type' is 'INCLUDE'.")

		return
	}
}

// applyIncludedClientsTargeting resolves organization_matching_type/included_clients (a List of
// plain client identifiers) from an alert policy response, shared by the alert correlation,
// prediction, and first-response policies. organization_matching_type is Computed, so it must
// always end up with a known value: at CLIENT scope (isMSP false) it's always null (the
// attribute doesn't apply there); at MSP scope it resolves to "INCLUDE" when the response has
// included clients, otherwise to "clientsIncluded" if the API supplied it, or "ALL" as the
// default when neither is present (some API responses omit "clientsIncluded" entirely for
// ALL-scoped policies). On Read/Import (fromRead), the fields are always refreshed from the
// response so remote targeting changes (e.g. a different client, or a switch to ALL) are
// reflected. On Create/Update (write responses), a field is only resolved when Terraform
// doesn't already have a known (user-configured) value, so the plan's representation is
// preserved untouched. Either way, an empty response list must map to a null list, not an empty
// one, or Terraform reports a "produced an unexpected new value" plan inconsistency for
// ALL-scoped policies.
func applyIncludedClientsTargeting(
	includedClients []string,
	clientsIncluded string,
	isMSP bool,
	fromRead bool,
	matchingType *types.String,
	stateIncludedClients *types.List,
) {
	if !isMSP {
		*matchingType = types.StringNull()
		*stateIncludedClients = types.ListNull(types.StringType)

		return
	}

	if len(includedClients) > 0 {
		if fromRead || matchingType.IsNull() || matchingType.IsUnknown() {
			*matchingType = types.StringValue("INCLUDE")
		}
		if fromRead || stateIncludedClients.IsNull() || stateIncludedClients.IsUnknown() {
			elems := make([]attr.Value, len(includedClients))
			for i, c := range includedClients {
				elems[i] = types.StringValue(c)
			}
			list, _ := types.ListValue(types.StringType, elems)
			*stateIncludedClients = list
		}

		return
	}

	resolvedMatchingType := "ALL"
	if clientsIncluded != "" {
		resolvedMatchingType = clientsIncluded
	}
	if fromRead || matchingType.IsNull() || matchingType.IsUnknown() {
		*matchingType = types.StringValue(resolvedMatchingType)
	}
	if fromRead || stateIncludedClients.IsNull() || stateIncludedClients.IsUnknown() {
		// Also covers the Go zero-value of types.List (e.g. a freshly declared model during
		// import), which reports IsNull() true but carries no element type until set here.
		*stateIncludedClients = types.ListNull(types.StringType)
	}
}
