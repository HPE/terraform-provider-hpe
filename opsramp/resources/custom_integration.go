// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &CustomIntegrationResource{}
var _ resource.ResourceWithImportState = &CustomIntegrationResource{}
var _ resource.ResourceWithModifyPlan = &CustomIntegrationResource{}

// CustomIntegrationResource defines the resource implementation.
type CustomIntegrationResource struct {
	BaseResource
}

// CustomIntegrationModel maps Terraform schema attributes to the provider model.
type CustomIntegrationModel struct {
	Id           types.String `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	Client       types.String `tfsdk:"client"`
	ClientId     types.String `tfsdk:"api_client_id"`
	ClientSecret types.String `tfsdk:"api_client_secret"`
	RoleName     types.String `tfsdk:"role_name"`
}

// NewCustomIntegration creates a new instance of the resource.
func NewCustomIntegration() resource.Resource {
	return &CustomIntegrationResource{}
}

// Metadata returns the resource type name.
func (r *CustomIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_integration"
}

// Schema defines the schema for the resource.
func (r *CustomIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Custom Integration resource. This resource creates an API OAuth2 token for integrations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the custom integration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the custom integration.",
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant) where the integration should be created. If not provided, the integration is created at the partner (provider tenant) level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_client_id": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The OAuth2 client ID generated for this integration. Use this for API authentication.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_client_secret": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The OAuth2 client secret generated for this integration. Use this for API authentication.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the role to assign to the integration. This determines the permissions for the API token.",
			},
		},
	}
}

// getTenantId determines which tenant ID to use based on the optional client parameter
func (r *CustomIntegrationResource) getTenantId(clientId types.String) string {
	// If client is specified, use the client ID as the tenant for client-level integration
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}
	// Otherwise, use the provider's tenant ID for partner-level integration
	return r.apiClient.TenantId
}

// Create handles the creation of the resource.
func (r *CustomIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the tenant ID (either client or provider tenant)
	tenantId := r.getTenantId(plan.Client)

	// Look up role ID by name
	role, err := r.apiClient.FindRoleByName(tenantId, plan.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Role Lookup Error",
			fmt.Sprintf("Could not find role with name '%s': %s", plan.RoleName.ValueString(), err.Error()))
		return
	}

	// Build the create request using role ID
	createIntegration := client.CreateCustomIntegration{
		DisplayName: plan.DisplayName.ValueString(),
		Category:    "Custom",
		InboundConfig: client.InboundConfig{
			Authentication: client.AuthenticationConfig{
				AuthType: "OAUTH2",
				Role: client.RoleClientRef{
					UniqueId: role.UniqueId,
				},
			},
		},
	}

	// Create the custom integration
	created, err := r.apiClient.CreateCustomIntegration(tenantId, createIntegration)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	// Update state with response values
	plan.Id = types.StringValue(created.Id)
	plan.DisplayName = types.StringValue(created.DisplayName)

	// Store OAuth2 credentials from inboundConfig.authentication.apiKeyPairs
	// These are only returned on creation and must be stored immediately
	apiKey, apiSecret := created.GetAPICredentials()
	if apiKey != "" {
		plan.ClientId = types.StringValue(apiKey)
		plan.ClientSecret = types.StringValue(apiSecret)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *CustomIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the tenant ID
	tenantId := r.getTenantId(state.Client)

	// Get the integration from the API
	existing, err := r.apiClient.GetCustomIntegration(tenantId, state.Id.ValueString())
	if err != nil {
		// Handle cases where integration no longer exists:
		// - 404 Not Found
		// - 500 with "No installed integration found" (API returns this when uninstalled)
		errStr := err.Error()
		if strings.Contains(errStr, "404") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "No installed integration found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Update state with current values (keep credentials from state as they're not returned on GET)
	state.DisplayName = types.StringValue(existing.DisplayName)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *CustomIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CustomIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CustomIntegrationModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the tenant ID
	tenantId := r.getTenantId(plan.Client)

	// Look up role ID by name
	role, err := r.apiClient.FindRoleByName(tenantId, plan.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Role Lookup Error",
			fmt.Sprintf("Could not find role with name '%s': %s", plan.RoleName.ValueString(), err.Error()))
		return
	}

	// Build the update request using role ID
	updateIntegration := client.CreateCustomIntegration{
		DisplayName: plan.DisplayName.ValueString(),
		Category:    "Custom",
		InboundConfig: client.InboundConfig{
			Authentication: client.AuthenticationConfig{
				AuthType: "OAUTH2",
				Role: client.RoleClientRef{
					UniqueId: role.UniqueId,
				},
			},
		},
	}

	// Update the custom integration
	updated, err := r.apiClient.UpdateCustomIntegration(tenantId, state.Id.ValueString(), updateIntegration)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	// Update state with response values
	plan.Id = state.Id
	plan.DisplayName = types.StringValue(updated.DisplayName)
	// Keep the credentials from state as they're not returned on update
	plan.ClientId = state.ClientId
	plan.ClientSecret = state.ClientSecret

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *CustomIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the tenant ID
	tenantId := r.getTenantId(state.Client)

	// Delete the custom integration
	err := r.apiClient.DeleteCustomIntegration(tenantId, state.Id.ValueString())
	if err != nil {
		// If already deleted, ignore the error
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())
			return
		}
	}
}

// ImportState handles resource import.
// Import ID format: <integration_id> or <client_id>:<integration_id> (MSP only)
func (r *CustomIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	var state CustomIntegrationModel
	state.Client = parsed.Client
	state.Id = types.StringValue(parsed.Parts[0])

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
