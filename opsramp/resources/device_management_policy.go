// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                   = &DeviceManagementPolicyResource{}
	_ resource.ResourceWithModifyPlan     = &DeviceManagementPolicyResource{}
	_ resource.ResourceWithImportState    = &DeviceManagementPolicyResource{}
	_ resource.ResourceWithValidateConfig = &DeviceManagementPolicyResource{}
)

// DeviceManagementPolicyResource defines the resource implementation.
type DeviceManagementPolicyResource struct {
	BaseResource
}

// DeviceManagementPolicyModel maps Terraform schema attributes to the provider model.
type DeviceManagementPolicyModel struct {
	Client              types.String                        `tfsdk:"client"`
	Id                  types.Int64                         `tfsdk:"id"`
	UID                 types.String                        `tfsdk:"uid"`
	Name                types.String                        `tfsdk:"name"`
	Actions             []DeviceManagementPolicyActionModel `tfsdk:"actions"`
	ExecutionOrder      types.Int64                         `tfsdk:"execution_order"`
	ClientMatchingType  types.String                        `tfsdk:"client_matching_type"`
	Clients             types.List                          `tfsdk:"clients"`
	SearchQuery         types.String                        `tfsdk:"search_query"`
	ParentPolicy        types.Int64                         `tfsdk:"parent_policy"`
	SystemQuery         types.String                        `tfsdk:"system_query"`
	CompleteSearchQuery types.String                        `tfsdk:"complete_search_query"`
	ParentPolicyUID     types.String                        `tfsdk:"parent_policy_uid"`
	Customized          types.Bool                          `tfsdk:"customized"`
	GlobalPolicy        types.Bool                          `tfsdk:"global_policy"`
}

// DeviceManagementPolicyActionModel represents a single policy action.
type DeviceManagementPolicyActionModel struct {
	Action                types.String                            `tfsdk:"action"`
	Items                 []DeviceManagementPolicyActionItemModel `tfsdk:"items"`
	ForceAssignOrUnassign types.Bool                              `tfsdk:"force_assign_or_unassign"`
}

// DeviceManagementPolicyActionItemModel represents a single item assigned by an action.
type DeviceManagementPolicyActionItemModel struct {
	Item       types.String `tfsdk:"item"`
	ItemType   types.String `tfsdk:"item_type"`
	ItemUID    types.String `tfsdk:"item_uid"`
	Parameters types.Map    `tfsdk:"parameters"`
}

// NewDeviceManagementPolicy creates a new instance of the resource.
func NewDeviceManagementPolicy() resource.Resource {
	return &DeviceManagementPolicyResource{}
}

// ValidateConfig validates that item_type is set on each item when action is "Assign Jobs".
func (r *DeviceManagementPolicyResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config DeviceManagementPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, action := range config.Actions {
		if action.Action.ValueString() != "Assign Jobs" {
			continue
		}
		for j, item := range action.Items {
			if item.ItemType.IsNull() || item.ItemType.ValueString() == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("actions").AtListIndex(i).AtName("items").AtListIndex(j).AtName("item_type"),
					"Missing required attribute",
					"item_type is required for each item when action is 'Assign Jobs'.",
				)
			}
		}
	}
}

// Metadata returns the resource type name.
func (r *DeviceManagementPolicyResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_device_management_policy"
}

// Schema defines the schema for the resource.
func (r *DeviceManagementPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Device Management Policy. Filters resources with rules and assigns " +
			"monitoring templates, availability rules, credential sets, custom attributes, jobs, and knowledge articles " +
			"to the matching resources. The policy is automatically run (saved and applied) on create and update.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this policy should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric identifier of the device management policy.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"uid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier (e.g. `POLICY-xxxxxxxx`) of the device management policy, used for import.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the device management policy.",
			},
			"actions": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Actions applied to resources matched by `rules`.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The action to perform. Supported values: `ASSIGN MONITORINGTEMPLATE`, " +
								"`UNASSIGN MONITORINGTEMPLATE`, `ASSIGN KNOWLEDGEARTICLE`, `Assign Availability Rule`, " +
								"`ASSIGN CUSTOMATTRIBUTES`, `Assign Jobs`, `ASSIGN CREDENTIALSET`.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"ASSIGN MONITORINGTEMPLATE",
									"UNASSIGN MONITORINGTEMPLATE",
									"ASSIGN KNOWLEDGEARTICLE",
									"Assign Availability Rule",
									"ASSIGN CUSTOMATTRIBUTES",
									"Assign Jobs",
									"ASSIGN CREDENTIALSET",
								),
							},
						},
						"force_assign_or_unassign": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to overwrite existing resource-level assignments for this action.",
						},
						"items": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The items assigned by this action.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"item": schema.StringAttribute{
										Optional: true,
										Computed: true,
										MarkdownDescription: "The identifier of the item to assign (e.g. template ID, credential set ID, `AND`/`ANY` for " +
											"availability rules). One of `item` or `item_uid` is required.",
										Validators: []validator.String{
											stringvalidator.AtLeastOneOf(path.Expressions{
												path.MatchRelative(),
												path.MatchRelative().AtParent().AtName("item_uid"),
											}...),
										},
									},
									"item_type": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "The type of the item, required for `Assign Jobs` (e.g. `assetinfo`).",
									},
									"item_uid": schema.StringAttribute{
										Optional: true,
										Computed: true,
										MarkdownDescription: "The unique identifier of the item, required for `Assign Jobs`. One of `item` or " +
											"`item_uid` is required.",
									},
									"parameters": schema.MapAttribute{
										Optional:            true,
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Configuration parameters to apply for `ASSIGN MONITORINGTEMPLATE` actions.",
									},
								},
							},
						},
					},
				},
			},
			"execution_order": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The execution order of the policy.",
			},
			"client_matching_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Applicable at partner-level only. How `clients` is applied. Valid values: `ALL`, " +
					"`INCLUDE`, `EXCLUDE`. Defaults to `ALL`. If `INCLUDE`/`EXCLUDE`, `clients` is required.",
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "INCLUDE", "EXCLUDE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clients": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The client UUIDs the policy applies to when `client_matching_type` is `INCLUDE` or `EXCLUDE`.",
			},
			"search_query": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Search query string used to filter resources, equivalent to the rules built in the UI.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_policy": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The numeric `id` of the parent policy to clone this policy from.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"parent_policy_uid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The `uid` of the parent policy, when this policy was cloned. Read-only.",
			},
			"system_query": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The system-generated query string applied in addition to `search_query`.",
			},
			"complete_search_query": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The combined user-defined and system-generated search query.",
			},
			"customized": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this policy has been customized (cloned) from a parent policy.",
			},
			"global_policy": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this is a global policy rather than a tenant-specific policy.",
			},
		},
	}
}

func buildDeviceManagementPolicyRequest(
	ctx context.Context,
	plan DeviceManagementPolicyModel,
) (client.DeviceManagementPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	policy := client.DeviceManagementPolicy{
		Name:               plan.Name.ValueString(),
		ResourceType:       "",
		MatchType:          "ALL",
		ClientMatchingType: plan.ClientMatchingType.ValueString(),
		SearchQuery:        plan.SearchQuery.ValueString(),
	}

	if !plan.Id.IsNull() {
		id := plan.Id.ValueInt64()
		policy.Id = int(id)
	}

	if !plan.ParentPolicy.IsNull() && !plan.ParentPolicy.IsUnknown() {
		policy.ParentPolicy = int(plan.ParentPolicy.ValueInt64())
	}

	if !plan.Clients.IsNull() && !plan.Clients.IsUnknown() {
		var clients []string
		diags.Append(plan.Clients.ElementsAs(ctx, &clients, false)...)
		policy.Clients = clients
	}

	policy.Rules = make([]client.DeviceManagementPolicyRule, 0)

	for _, action := range plan.Actions {
		a := client.DeviceManagementPolicyAction{
			Action: action.Action.ValueString(),
		}
		if !action.ForceAssignOrUnassign.IsNull() {
			force := action.ForceAssignOrUnassign.ValueBool()
			a.ForceAssignOrUnassign = &force
		}
		for _, item := range action.Items {
			i := client.DeviceManagementPolicyActionItem{
				Item:     item.Item.ValueString(),
				ItemType: item.ItemType.ValueString(),
				ItemUID:  item.ItemUID.ValueString(),
			}
			if !item.Parameters.IsNull() && !item.Parameters.IsUnknown() {
				params := map[string]string{}
				diags.Append(item.Parameters.ElementsAs(ctx, &params, false)...)
				i.Parameters = params
			}
			a.Items = append(a.Items, i)
		}
		policy.Actions = append(policy.Actions, a)
	}

	return policy, diags
}

func mapDeviceManagementPolicyToState(
	ctx context.Context,
	resp *client.DeviceManagementPolicy,
	state *DeviceManagementPolicyModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Id = types.Int64Value(int64(resp.Id))
	state.UID = types.StringValue(resp.UID)
	state.Name = types.StringValue(resp.Name)
	state.ExecutionOrder = types.Int64Value(int64(resp.ExecutionOrder))
	state.ClientMatchingType = types.StringValue(resp.ClientMatchingType)
	state.SearchQuery = types.StringValue(resp.SearchQuery)
	state.SystemQuery = types.StringValue(resp.SystemQuery)
	state.CompleteSearchQuery = types.StringValue(resp.CompleteSearchQuery)
	state.ParentPolicyUID = types.StringValue(resp.ParentPolicyUID)
	state.Customized = types.BoolValue(resp.Customized)
	state.GlobalPolicy = types.BoolValue(resp.GlobalPolicy)

	if resp.ParentPolicy != 0 {
		state.ParentPolicy = types.Int64Value(int64(resp.ParentPolicy))
	}

	if len(resp.Clients) > 0 {
		clients, d := types.ListValueFrom(ctx, types.StringType, resp.Clients)
		diags.Append(d...)
		state.Clients = clients
	} else if state.Clients.IsNull() || state.Clients.IsUnknown() {
		state.Clients = types.ListNull(types.StringType)
	}

	// The API resolves/echoes item identifiers and parameters differently than
	// submitted (e.g. "Assign Jobs" swaps item/item_uid, and monitoring template
	// actions can gain server-injected default parameters). Capture the
	// configured values before they are overwritten below so they can be
	// preserved instead of the response's, avoiding "inconsistent result after
	// apply" errors.
	priorActions := state.Actions

	state.Actions = make([]DeviceManagementPolicyActionModel, len(resp.Actions))
	for i, action := range resp.Actions {
		am := DeviceManagementPolicyActionModel{
			Action: types.StringValue(action.Action),
		}
		if action.ForceAssignOrUnassign != nil {
			am.ForceAssignOrUnassign = types.BoolValue(*action.ForceAssignOrUnassign)
		}

		var priorItems []DeviceManagementPolicyActionItemModel
		if i < len(priorActions) {
			priorItems = priorActions[i].Items
		}

		am.Items = make([]DeviceManagementPolicyActionItemModel, len(action.Items))
		for j, item := range action.Items {
			im := DeviceManagementPolicyActionItemModel{
				Item:     types.StringValue(item.Item),
				ItemType: types.StringValue(item.ItemType),
				ItemUID:  types.StringValue(item.ItemUID),
			}

			var priorItem *DeviceManagementPolicyActionItemModel
			if j < len(priorItems) {
				priorItem = &priorItems[j]
			}

			// Only preserve the configured value when it was already known; an
			// unknown prior value (e.g. derived from another resource) must be
			// resolved from the API response, which defaults to "" when the API
			// does not echo it back.
			if priorItem != nil {
				if !priorItem.Item.IsNull() && !priorItem.Item.IsUnknown() {
					im.Item = priorItem.Item
				}
				if !priorItem.ItemType.IsNull() && !priorItem.ItemType.IsUnknown() {
					im.ItemType = priorItem.ItemType
				}
				if !priorItem.ItemUID.IsNull() && !priorItem.ItemUID.IsUnknown() {
					im.ItemUID = priorItem.ItemUID
				}
			}

			if priorItem != nil && !priorItem.Parameters.IsNull() && !priorItem.Parameters.IsUnknown() {
				im.Parameters = priorItem.Parameters
			} else if len(item.Parameters) > 0 {
				params, d := types.MapValueFrom(ctx, types.StringType, item.Parameters)
				diags.Append(d...)
				im.Parameters = params
			} else {
				im.Parameters = types.MapNull(types.StringType)
			}

			am.Items[j] = im
		}
		state.Actions[i] = am
	}

	return diags
}

// Create handles the creation of the resource.
func (r *DeviceManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceManagementPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy, diags := buildDeviceManagementPolicyRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.apiClient.CreateDeviceManagementPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())

		return
	}

	// Persist policy before trying to run
	resp.Diagnostics.Append(mapDeviceManagementPolicyToState(ctx, created, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if err := r.apiClient.RunDeviceManagementPolicy(tenantId, created.UID); err != nil {
		resp.Diagnostics.AddWarning("Run Policy warning", fmt.Sprintf(
			"Policy %s (%q) was created but could not be run: %s", created.Name, created.UID, err.Error(),
		))

		return
	}
}

// Read handles reading the resource state.
func (r *DeviceManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceManagementPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetDeviceManagementPolicy(tenantId, state.UID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	resp.Diagnostics.Append(mapDeviceManagementPolicyToState(ctx, existing, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *DeviceManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DeviceManagementPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	plan.Id = state.Id

	policy, diags := buildDeviceManagementPolicyRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.apiClient.UpdateDeviceManagementPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())

		return
	}

	if err := r.apiClient.RunDeviceManagementPolicy(tenantId, updated.UID); err != nil {
		resp.Diagnostics.AddError("Run Policy Error", fmt.Sprintf(
			"Policy %q was updated but could not be run: %s", updated.UID, err.Error(),
		))

		return
	}

	plan.Client = state.Client
	resp.Diagnostics.Append(mapDeviceManagementPolicyToState(ctx, updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *DeviceManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceManagementPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteDeviceManagementPolicy(tenantId, state.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports an existing device management policy by its uid, optionally
// prefixed with a client UUID (MSP scope only): "<client_id>:<uid>".
func (r *DeviceManagementPolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	tenantId := r.TenantForImport(parsed)
	uid := parsed.Parts[0]

	existing, err := r.apiClient.GetDeviceManagementPolicy(tenantId, uid)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Device Management Policy",
			fmt.Sprintf("Could not import device management policy with uid '%s': %s", uid, err))

		return
	}

	var state DeviceManagementPolicyModel
	state.Client = parsed.Client
	resp.Diagnostics.Append(mapDeviceManagementPolicyToState(ctx, existing, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
