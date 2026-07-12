// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces.
var _ resource.Resource = &DeviceGroupResource{}
var _ resource.ResourceWithModifyPlan = &DeviceGroupResource{}

// DeviceGroupModel maps Terraform schema attributes to the provider model.
type DeviceGroupModel struct {
	Client      types.String `tfsdk:"client"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	EntityType  types.String `tfsdk:"entity_type"`
	ParentId    types.String `tfsdk:"parent_id"`
	SearchQuery types.String `tfsdk:"search_query"`
	Resources   types.Set    `tfsdk:"resources"`
}

// DeviceGroupResource defines the resource implementation.
type DeviceGroupResource struct {
	BaseResource
}

// NewDeviceGroup creates a new instance of the resource.
func NewDeviceGroup() resource.Resource {
	return &DeviceGroupResource{}
}

// Metadata returns the resource type name.
func (r *DeviceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

// Recursive schema definition example
func (r *DeviceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this device group should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Unique identifier for the group.",
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"parent_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"entity_type": schema.StringAttribute{
				Computed: true,
			},
			"search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The search query for selecting resources.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resources": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of resource IDs to assign to this device group.",
			},
		},
	}

}

func translatePlanToDeviceGroupModel(plan DeviceGroupModel) client.DeviceGroupAPI {
	var parent *client.Parent

	if plan.ParentId.ValueString() != "" {
		parent = &client.Parent{
			Id: plan.ParentId.ValueString(),
		}
	}

	var filterCriteria *client.FilterCriteria
	if !plan.SearchQuery.IsNull() && !plan.SearchQuery.IsUnknown() && plan.SearchQuery.ValueString() != "" {
		filterCriteria = &client.FilterCriteria{
			SearchQuery: plan.SearchQuery.ValueString(),
		}
	}

	deviceGroup := client.DeviceGroupAPI{
		Id:             plan.Id.ValueString(),
		Name:           plan.Name.ValueString(),
		Parent:         parent,
		EntityType:     "DEVICE_GROUP",
		FilterCriteria: filterCriteria,
	}

	return deviceGroup
}

func mapDeviceGroupResponseToModel(api *client.DeviceGroupAPI, model *DeviceGroupModel) {
	model.Id = types.StringValue(api.Id)
	model.Name = types.StringValue(api.Name)
	model.EntityType = types.StringValue(api.EntityType)

	if api.FilterCriteria != nil && api.FilterCriteria.SearchQuery != "" {
		model.SearchQuery = types.StringValue(api.FilterCriteria.SearchQuery)
	} else {
		model.SearchQuery = types.StringValue("")
	}

	if api.Parent != nil {
		model.ParentId = types.StringValue(api.Parent.Id)
	}
}

func deviceGroupResourcesSet(ctx context.Context, ids []string) (types.Set, diag.Diagnostics) {
	if ids == nil {
		ids = []string{}
	}

	return types.SetValueFrom(ctx, types.StringType, ids)
}

func (r *DeviceGroupResource) readDeviceGroupResources(ctx context.Context, tenantId string, deviceGroupId string) (types.Set, diag.Diagnostics) {
	resources, err := r.apiClient.GetDeviceGroupChilds(tenantId, deviceGroupId)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Read Error", err.Error())
		return types.SetUnknown(types.StringType), diags
	}

	ids := make([]string, 0, len(resources))
	for _, resource := range resources {
		// Empty type for child groups
		if resource.ResourceType != "" {
			ids = append(ids, resource.Uuid)
		}
	}

	return deviceGroupResourcesSet(ctx, ids)
}

func createDeviceGroupWithRetry(apiClient *client.OpsRampClient, tenantId string, req client.DeviceGroupAPI, hasParent bool) (*client.DeviceGroupAPI, error) {
	const maxAttempts = 4
	const retryDelay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		created, err := apiClient.CreateDeviceGroup(tenantId, req)
		if err == nil {
			return created, nil
		}

		lastErr = err
		// Retry only for transient server errors while attaching to a parent.
		if !hasParent || !strings.Contains(err.Error(), "status: 500") || attempt == maxAttempts {
			break
		}

		time.Sleep(retryDelay)
	}

	return nil, lastErr
}

func (r *DeviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceGroup := translatePlanToDeviceGroupModel(plan)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Create the device group. Parent-scoped creation may intermittently return 500,
	// so retry a few times for that specific transient error.
	newDeviceGroup, err := createDeviceGroupWithRetry(r.apiClient, tenantId, deviceGroup, !plan.ParentId.IsNull() && plan.ParentId.ValueString() != "")
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	mapDeviceGroupResponseToModel(newDeviceGroup, &plan)

	// Unset resources is equivalent to []: always reconcile.

	resource_ids := []string{}
	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() && len(plan.Resources.Elements()) > 0 {
		resource_ids = setToStringSlice(plan.Resources)

		// Add child resources if specified - don't fail the entire creation if this part fails, but do report it in diagnostics.
		if err := r.apiClient.AddDeviceGroupChilds(tenantId, newDeviceGroup.Id, resource_ids); err != nil {
			resp.Diagnostics.AddError("Error adding resources to device group", err.Error())
		}
	}

	plan.Resources, diags = r.readDeviceGroupResources(ctx, tenantId, newDeviceGroup.Id)
	resp.Diagnostics.Append(diags...)

	// Save new state back to Terraform
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// --- Read ---
func (r *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	backendDeviceGroup, err := r.apiClient.GetDeviceGroup(tenantId, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "No Device group exists") {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	if backendDeviceGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resources, diags := r.readDeviceGroupResources(ctx, tenantId, backendDeviceGroup.Id)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState := DeviceGroupModel{Client: state.Client}
	mapDeviceGroupResponseToModel(backendDeviceGroup, &newState)
	newState.Resources = resources

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

// Update modifies an existing resource.
func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DeviceGroupModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	deviceGroup := translatePlanToDeviceGroupModel(plan)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Update the device group
	newDeviceGroup, err := r.apiClient.CreateDeviceGroup(tenantId, deviceGroup)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	mapDeviceGroupResponseToModel(newDeviceGroup, &state)
	state.Client = plan.Client

	// Unset resources is equivalent to []: always reconcile.
	oldIds := setToStringSlice(state.Resources)
	newIds := setToStringSlice(plan.Resources)
	toAdd := stringSetDiff(newIds, oldIds)
	toRemove := stringSetDiff(oldIds, newIds)

	if len(toAdd) > 0 {
		if err := r.apiClient.AddDeviceGroupChilds(tenantId, newDeviceGroup.Id, toAdd); err != nil {
			resp.Diagnostics.AddError("Error adding resources to device group", err.Error())
			return
		}
	}

	if len(toRemove) > 0 {
		if err := r.apiClient.RemoveDeviceGroupChilds(tenantId, newDeviceGroup.Id, toRemove); err != nil {
			resp.Diagnostics.AddError("Error removing resources from device group", err.Error())
			return
		}
	}

	var diags diag.Diagnostics
	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() {
		state.Resources = plan.Resources
	} else {
		state.Resources, diags = deviceGroupResourcesSet(ctx, newIds)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save
	resp.State.Set(ctx, &state)
}

// Delete removes the resource from the API.
func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceGroupModel
	req.State.Get(ctx, &state)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteDeviceGroup(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DeviceGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.BaseResource.ModifyPlan(ctx, req, resp)
	if resp.Diagnostics.HasError() || req.Plan.Raw.IsNull() {
		return
	}

	var plan DeviceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Removing the resources attribute from config is equivalent to resources = [].
	if plan.Resources.IsUnknown() {
		plan.Resources, _ = deviceGroupResourcesSet(ctx, nil)
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}
