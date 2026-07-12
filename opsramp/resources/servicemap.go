// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ServicemapResource{}
var _ resource.ResourceWithImportState = &ServicemapResource{}
var _ resource.ResourceWithModifyPlan = &ServicemapResource{}

// ServicemapResource defines the resource implementation.
type ServicemapResource struct {
	BaseResource
}

// NewServicemap creates a new instance of the resource.
func NewServicemap() resource.Resource {
	return &ServicemapResource{}
}

// Metadata returns the resource type name.
func (r *ServicemapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicemap"
}

// Recursive schema definition example
func (r *ServicemapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			// parent attribute
			"parent": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// link attribute
			"link": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// resources at root-level
			"resources": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			// search query attribute (e.g., "name CONTAINS pro")
			"search_query": schema.StringAttribute{
				Optional: true,
			},
			// alert type configuration attribute
			"alert_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("any-critical-alert", "availability", "critical-alert"),
				},
			},
			// metrics for critical-alert type
			"metrics": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},

			"threshold_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("count"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.OneOf("count", "percentage"),
				},
			},
			"threshold_limit": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				// Usa el validador de int64 del paquete de validators
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				Default: int64default.StaticInt64(1),
			},
		},
	}
}

func translatePlanToModel(plan ServicemapModel) client.CreateServicemap {
	var resources []string
	for _, i := range plan.Resources {
		resources = append(resources, i.ValueString())
	}

	var parent *client.Parent

	if plan.Parent.ValueString() != "" {
		parent = &client.Parent{
			Id: plan.Parent.ValueString(),
		}
	}

	var serviceAvailabilityMonitor *client.ServiceAvailabilityMonitor
	if !plan.AlertType.IsNull() && plan.AlertType.ValueString() != "" {
		var metrics []string
		for _, m := range plan.Metrics {
			metrics = append(metrics, m.ValueString())
		}
		if metrics == nil {
			metrics = []string{}
		}
		serviceAvailabilityMonitor = &client.ServiceAvailabilityMonitor{
			AlertType: plan.AlertType.ValueString(),
			Metrics:   metrics,
			MatchType: "ANY",
		}
	}

	var availabilityThreshold = client.AvailabilityThreshold{
		ThresholdType:  plan.ThresholdType.ValueString(),
		ThresholdLimit: int(plan.ThresholdLimit.ValueInt64()),
	}

	return client.CreateServicemap{
		Name:                       plan.Name.ValueString(),
		Type:                       plan.Type.ValueString(),
		Id:                         plan.Id.ValueString(),
		Parent:                     parent,
		Resources:                  resources,
		SearchQuery:                plan.SearchQuery.ValueString(),
		ServiceAvailabilityMonitor: serviceAvailabilityMonitor,
		AvailabilityThreshold:      availabilityThreshold,
	}
}

func (r *ServicemapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the JSON from the plan
	createServicemap := translatePlanToModel(plan)
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Create the Servicemap in the backend
	newServicemap, err := r.apiClient.CreateServicemap(tenantId, createServicemap)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	// Assign the backend response directly into the state
	plan.Id = types.StringValue(newServicemap.Id)

	// Save new state back to Terraform
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// --- Read ---
func (r *ServicemapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServicemapModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	backendServicemap, err := r.apiClient.GetServicemap(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "No Service group") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	if backendServicemap == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update modifies an existing resource.
func (r *ServicemapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Id cannot be changed; ensure Id stays the same.
	if plan.Id.ValueString() != "" && plan.Id.ValueString() != state.Id.ValueString() {
		resp.Diagnostics.AddError("Id Immutable", "Id cannot be changed after creation.")
		return
	}

	// Use the JSON from the state
	updateServicemap := translatePlanToModel(state)
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	if !plan.Name.IsNull() {
		updateServicemap.Name = plan.Name.ValueString()
	}
	if !plan.Type.IsNull() {
		updateServicemap.Type = plan.Type.ValueString()
	}
	if !reflect.DeepEqual(plan.Resources, state.Resources) {
		updateServicemap.Resources = make([]string, 0)
		for _, i := range plan.Resources {
			updateServicemap.Resources = append(updateServicemap.Resources, i.ValueString())
		}
	}
	if !plan.SearchQuery.IsNull() {
		updateServicemap.SearchQuery = plan.SearchQuery.ValueString()
	}
	if !plan.AlertType.IsNull() && plan.AlertType.ValueString() != "" {
		var metrics []string
		for _, m := range plan.Metrics {
			metrics = append(metrics, m.ValueString())
		}
		if metrics == nil {
			metrics = []string{}
		}
		updateServicemap.ServiceAvailabilityMonitor = &client.ServiceAvailabilityMonitor{
			AlertType: plan.AlertType.ValueString(),
			Metrics:   metrics,
			MatchType: "ANY",
		}
	}

	updateServicemap.AvailabilityThreshold = client.AvailabilityThreshold{
		ThresholdType:  plan.ThresholdType.ValueString(),
		ThresholdLimit: int(plan.ThresholdLimit.ValueInt64()),
	}

	updatedServicemap, err := r.apiClient.CreateServicemap(tenantId, updateServicemap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("Could not update resource: %s", err),
		)
		return
	}

	state.Id = types.StringValue(updatedServicemap.Id)
	state.Name = types.StringValue(updatedServicemap.Name)
	state.Type = types.StringValue(updatedServicemap.Type)
	state.ThresholdLimit = types.Int64Value(int64(updatedServicemap.AvailabilityThreshold.ThresholdLimit))
	state.ThresholdType = types.StringValue(updatedServicemap.AvailabilityThreshold.ThresholdType)
	state.SearchQuery = plan.SearchQuery
	state.AlertType = plan.AlertType

	// Save
	resp.State.Set(ctx, &state)
}

// Delete removes the resource from the API.
func (r *ServicemapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServicemapModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteServicemap(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing servicemap.
// Import ID format: <servicemap_id> or <client_id>:<servicemap_id> (MSP only)
func (r *ServicemapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	tenantId := r.TenantForImport(parsed)
	id := parsed.Parts[0]

	existing, err := r.apiClient.GetServicemap(tenantId, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Servicemap",
			fmt.Sprintf("Could not import servicemap with ID '%s': %s", id, err),
		)
		return
	}

	if existing == nil {
		resp.Diagnostics.AddError(
			"Error Importing Servicemap",
			fmt.Sprintf("Servicemap with ID '%s' not found", id),
		)
		return
	}

	state := ServicemapModel{
		Client:         parsed.Client,
		Id:             types.StringValue(existing.Id),
		Name:           types.StringValue(existing.Name),
		Type:           types.StringValue(existing.Type),
		ThresholdType:  types.StringValue(existing.AvailabilityThreshold.ThresholdType),
		ThresholdLimit: types.Int64Value(int64(existing.AvailabilityThreshold.ThresholdLimit)),
	}

	resp.State.Set(ctx, &state)
}

type ServicemapModel struct {
	Client         types.String   `tfsdk:"client"`
	Id             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Type           types.String   `tfsdk:"type"`
	Parent         types.String   `tfsdk:"parent"`
	Link           types.String   `tfsdk:"link"`
	Resources      []types.String `tfsdk:"resources"`
	SearchQuery    types.String   `tfsdk:"search_query"`
	AlertType      types.String   `tfsdk:"alert_type"`
	Metrics        []types.String `tfsdk:"metrics"`
	ThresholdType  types.String   `tfsdk:"threshold_type"`
	ThresholdLimit types.Int64    `tfsdk:"threshold_limit"`
}

func (r *ServicemapResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError("ServiceMaps can only be created at Client level", "Use a client-scoped provider configuration or specify the client using unique ID.")
		return
	}
}
