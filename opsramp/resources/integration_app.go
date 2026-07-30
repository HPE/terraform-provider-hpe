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
var (
	_ resource.Resource                = &IntegrationAppResource{}
	_ resource.ResourceWithImportState = &IntegrationAppResource{}
)

// IntegrationAppResource defines the resource implementation for SDK APP integrations (v3 API).
type IntegrationAppResource struct {
	BaseResource
}

// IntegrationAppModel maps Terraform schema attributes to the provider model.
type IntegrationAppModel struct {
	Id                           types.String `tfsdk:"id"`
	Application                  types.String `tfsdk:"application"`
	DisplayName                  types.String `tfsdk:"display_name"`
	Version                      types.String `tfsdk:"version"`
	Client                       types.String `tfsdk:"client"`
	Status                       types.String `tfsdk:"status"`
	State                        types.String `tfsdk:"state"`
	BypassResourceReconciliation types.Bool   `tfsdk:"bypass_resource_reconciliation"`
}

// NewIntegrationApp creates a new instance of the resource.
func NewIntegrationApp() resource.Resource {
	return &IntegrationAppResource{}
}

// Metadata returns the resource type name.
func (r *IntegrationAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_app"
}

// Schema defines the schema for the resource.
func (r *IntegrationAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp SDK APP Integration (v3 API). Use this resource for integrations with category `SDK APP` (e.g. hpe-alletra-6000, snmp, kubernetes).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the installed app (e.g. INTG-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The application identifier (e.g. `hpe-alletra-6000`, `snmp`, `kubernetes`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name assigned by OpsRamp after installation.",
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The version of the app to install (e.g. `1.0.0`, `2.0.0`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant) where the app should be installed. If not provided, uses the provider tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the installed app (e.g. Install Request Sent).",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The state of the installed app (e.g. Installed).",
			},
			"bypass_resource_reconciliation": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable bypass resource reconciliation during discovery for this installation.",
			},
		},
	}
}

// getTenantId determines which tenant ID to use based on the optional client parameter
func (r *IntegrationAppResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}

	return r.apiClient.TenantId
}

// Create handles the creation of the resource.
func (r *IntegrationAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntegrationAppModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	v3Req := client.InstallIntegrationV3Request{
		App:                       plan.Application.ValueString(),
		Version:                   plan.Version.ValueString(),
		MultiAppsDiscoveryEnabled: plan.BypassResourceReconciliation.ValueBool(),
	}

	installed, err := r.apiClient.InstallIntegrationV3(tenantId, v3Req)
	if err != nil {
		resp.Diagnostics.AddError("Integration App Install Error",
			fmt.Sprintf("Could not install app '%s': %s", plan.Application.ValueString(), err.Error()))

		return
	}

	plan.Id = types.StringValue(installed.ID)
	plan.DisplayName = types.StringValue(installed.DisplayName)
	plan.Status = types.StringValue(installed.Status)
	plan.State = types.StringValue(installed.State)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *IntegrationAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IntegrationAppModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetIntegrationV3(tenantId, state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	state.DisplayName = types.StringValue(existing.DisplayName)
	state.Status = types.StringValue(existing.Status)
	state.State = types.StringValue(existing.State)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
// SDK APP integrations are immutable (all mutable fields use RequiresReplace),
// so this should not normally be called.
func (r *IntegrationAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported",
		"opsramp_integration_app does not support in-place updates. All changes require replacement.")
}

// Delete handles deleting the resource.
func (r *IntegrationAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IntegrationAppModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteIntegrationV3(tenantId, state.Id.ValueString(), "Terraform - Resource destroyed")
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())

			return
		}
	}
}

// ImportState handles resource import.
// Import ID format: <integration_id> or <client_id>:<integration_id> (MSP only)
func (r *IntegrationAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	var state IntegrationAppModel
	state.Client = parsed.Client
	state.Id = types.StringValue(parsed.Parts[0])

	// These will be populated on the next Read
	state.Application = types.StringUnknown()
	state.DisplayName = types.StringUnknown()
	state.Version = types.StringUnknown()
	// state.ProfileId = types.StringUnknown()
	state.Status = types.StringUnknown()
	state.State = types.StringUnknown()

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
