// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package role

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_role"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = RoleResourceSchema(ctx)
}

// This function breaks out the logic of reading permissions from API response to store to state.
func populateGetRoleAsStatePermissions(ctx context.Context, r *sdk.GetRole200Response) (PermissionsValue, diag.Diagnostics) {

	var features []FeaturePermissionsValue
	for _, v := range r.FeaturePermissions {
		features = append(features, FeaturePermissionsValue{
			Code:        types.StringValue(v.GetCode()),
			Access:      types.StringValue(v.GetAccess()),
			Id:          types.Int64Value(v.GetId()),
			Name:        types.StringValue(v.GetName()),
			SubCategory: types.StringValue(v.GetSubCategory()),
			state:       attr.ValueStateKnown,
		})
	}

	var blueprints []BlueprintPermissionsValue
	for _, v := range r.AppTemplatePermissions {
		blueprints = append(blueprints, BlueprintPermissionsValue{
			Name:   types.StringValue(v.GetName()),
			Id:     types.Int64Value(v.GetId()),
			Access: types.StringValue(v.GetAccess()),
			state:  attr.ValueStateKnown,
		})
	}

	var catalogItemTypes []CatalogItemTypePermissionsValue
	for _, v := range r.CatalogItemTypePermissions {
		catalogItemTypes = append(catalogItemTypes, CatalogItemTypePermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Access: types.StringValue(v.GetAccess()),
			Name:   types.StringValue(v.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var clouds []CloudPermissionsValue
	for _, v := range r.Zones {
		clouds = append(clouds, CloudPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Access: types.StringValue(v.GetAccess()),
			Name:   types.StringValue(v.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var groups []GroupPermissionsValue
	for _, v := range r.Sites {
		groups = append(groups, GroupPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Access: types.StringValue(v.GetAccess()),
			Name:   types.StringValue(v.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var instanceTypes []InstanceTypePermissionsValue
	for _, v := range r.InstanceTypePermissions {
		instanceTypes = append(instanceTypes, InstanceTypePermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			state:  attr.ValueStateKnown,
		})
	}

	var personas []PersonaPermissionsValue
	for _, v := range r.PersonaPermissions {
		personas = append(personas, PersonaPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			Code:   types.StringValue(v.GetCode()),
			state:  attr.ValueStateKnown,
		})
	}

	var reportTypes []ReportTypePermissionsValue
	for _, v := range r.ReportTypePermissions {
		reportTypes = append(reportTypes, ReportTypePermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			Code:   types.StringValue(v.GetCode()),
			state:  attr.ValueStateKnown,
		})
	}

	var tasks []TaskPermissionsValue
	for _, v := range r.TaskPermissions {
		tasks = append(tasks, TaskPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			Code:   types.StringPointerValue(v.Code.Get()),
			state:  attr.ValueStateKnown,
		})
	}

	var vdiPools []VdiPoolPermissionsValue
	for _, v := range r.VdiPoolPermissions {
		vdiPools = append(vdiPools, VdiPoolPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			state:  attr.ValueStateKnown,
		})
	}

	var workflows []WorkflowPermissionsValue
	for _, v := range r.TaskSetPermissions {
		workflows = append(workflows, WorkflowPermissionsValue{
			Id:     types.Int64Value(v.GetId()),
			Name:   types.StringValue(v.GetName()),
			Access: types.StringValue(v.GetAccess()),
			state:  attr.ValueStateKnown,
		})
	}

	featuresSet, diags := types.SetValueFrom(ctx, FeaturePermissionsValue{}.Type(ctx), features)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	blueprintsSet, diags := types.SetValueFrom(ctx, BlueprintPermissionsValue{}.Type(ctx), blueprints)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	catalogItemTypesSet, diags := types.SetValueFrom(ctx, CatalogItemTypePermissionsValue{}.Type(ctx), catalogItemTypes)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	cloudsSet, diags := types.SetValueFrom(ctx, CloudPermissionsValue{}.Type(ctx), clouds)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	groupsSet, diags := types.SetValueFrom(ctx, GroupPermissionsValue{}.Type(ctx), groups)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	instanceTypesSet, diags := types.SetValueFrom(ctx, InstanceTypePermissionsValue{}.Type(ctx), instanceTypes)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	personasSet, diags := types.SetValueFrom(ctx, PersonaPermissionsValue{}.Type(ctx), personas)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	reportTypesSet, diags := types.SetValueFrom(ctx, ReportTypePermissionsValue{}.Type(ctx), reportTypes)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	tasksSet, diags := types.SetValueFrom(ctx, TaskPermissionsValue{}.Type(ctx), tasks)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	vdiPoolsSet, diags := types.SetValueFrom(ctx, VdiPoolPermissionsValue{}.Type(ctx), vdiPools)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	workflowsSet, diags := types.SetValueFrom(ctx, WorkflowPermissionsValue{}.Type(ctx), workflows)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	return NewPermissionsValue(PermissionsValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"default_blueprint_access":         convert.StrToType(r.GlobalAppTemplateAccess),
		"default_catalog_item_type_access": convert.StrToType(r.GlobalCatalogItemTypeAccess),
		"default_cloud_access":             convert.StrToType(r.GlobalZoneAccess),
		"default_group_access":             convert.StrToType(r.GlobalSiteAccess),
		"default_instance_type_access":     convert.StrToType(r.GlobalInstanceTypeAccess),
		"default_persona_access":           convert.StrToType(r.GlobalPersonaAccess),
		"default_report_type_access":       convert.StrToType(r.GlobalReportTypeAccess),
		"default_task_access":              convert.StrToType(r.GlobalTaskAccess),
		"default_vdi_pool_access":          convert.StrToType(r.GlobalVdiPoolAccess),
		"default_workflow_access":          convert.StrToType(r.GlobalTaskSetAccess),
		"feature_permissions":              featuresSet,
		"blueprint_permissions":            blueprintsSet,
		"catalog_item_type_permissions":    catalogItemTypesSet,
		"cloud_permissions":                cloudsSet,
		"group_permissions":                groupsSet,
		"instance_type_permissions":        instanceTypesSet,
		"persona_permissions":              personasSet,
		"report_type_permissions":          reportTypesSet,
		"task_permissions":                 tasksSet,
		"vdi_pool_permissions":             vdiPoolsSet,
		"workflow_permissions":             workflowsSet,
	})
}

// Helper function to break out the logic of setting permissions in create.
func setPermissionsInCreate(
	ctx context.Context,
	plan *RoleModel,
	addRole *sdk.AddRolesRequestRole,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if !plan.Permissions.DefaultBlueprintAccess.IsUnknown() {
		addRole.SetGlobalAppTemplateAccess(plan.Permissions.DefaultBlueprintAccess.ValueString())
	}

	if !plan.Permissions.DefaultCatalogItemTypeAccess.IsUnknown() {
		addRole.SetGlobalCatalogItemTypeAccess(plan.Permissions.DefaultCatalogItemTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultCloudAccess.IsUnknown() {
		addRole.SetGlobalZoneAccess(plan.Permissions.DefaultCloudAccess.ValueString())
	}

	if !plan.Permissions.DefaultGroupAccess.IsUnknown() {
		addRole.SetGlobalSiteAccess(plan.Permissions.DefaultGroupAccess.ValueString())
	}

	if !plan.Permissions.DefaultInstanceTypeAccess.IsUnknown() {
		addRole.SetGlobalInstanceTypeAccess(plan.Permissions.DefaultInstanceTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultPersonaAccess.IsUnknown() {
		addRole.SetGlobalPersonaAccess(plan.Permissions.DefaultPersonaAccess.ValueString())
	}

	if !plan.Permissions.DefaultReportTypeAccess.IsUnknown() {
		addRole.SetGlobalReportTypeAccess(plan.Permissions.DefaultReportTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultTaskAccess.IsUnknown() {
		addRole.SetGlobalTaskAccess(plan.Permissions.DefaultTaskAccess.ValueString())
	}

	if !plan.Permissions.DefaultVdiPoolAccess.IsUnknown() {
		addRole.SetGlobalVdiPoolAccess(plan.Permissions.DefaultVdiPoolAccess.ValueString())
	}

	if !plan.Permissions.DefaultWorkflowAccess.IsUnknown() {
		addRole.SetGlobalTaskSetAccess(plan.Permissions.DefaultWorkflowAccess.ValueString())
	}

	if !plan.Permissions.FeaturePermissions.IsUnknown() {
		var featurePermissions []FeaturePermissionsValue
		diags := plan.Permissions.FeaturePermissions.ElementsAs(ctx, &featurePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleFeaturePermissions []sdk.AddRolesRequestRoleFeaturePermissionsInner
		for _, v := range featurePermissions {
			addRoleFeaturePermissions = append(addRoleFeaturePermissions, sdk.AddRolesRequestRoleFeaturePermissionsInner{
				Access: v.Access.ValueString(),
				Code:   v.Code.ValueString(),
			})
		}

		addRole.SetFeaturePermissions(addRoleFeaturePermissions)
	}

	if !plan.Permissions.BlueprintPermissions.IsUnknown() {
		var blueprintPermissions []BlueprintPermissionsValue
		diags = plan.Permissions.BlueprintPermissions.ElementsAs(ctx, &blueprintPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleBlueprintPermissions []sdk.AddRolesRequestRoleAppTemplatePermissionsInner
		for _, v := range blueprintPermissions {
			addRoleBlueprintPermissions = append(addRoleBlueprintPermissions, sdk.AddRolesRequestRoleAppTemplatePermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetAppTemplatePermissions(addRoleBlueprintPermissions)
	}

	if !plan.Permissions.CatalogItemTypePermissions.IsUnknown() {
		var catalogItemTypePermissions []CatalogItemTypePermissionsValue
		diags = plan.Permissions.CatalogItemTypePermissions.ElementsAs(ctx, &catalogItemTypePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleCatalogItemTypePermissions []sdk.AddRolesRequestRoleCatalogItemTypePermissionsInner
		for _, v := range catalogItemTypePermissions {
			addRoleCatalogItemTypePermissions = append(addRoleCatalogItemTypePermissions, sdk.AddRolesRequestRoleCatalogItemTypePermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetCatalogItemTypePermissions(addRoleCatalogItemTypePermissions)
	}

	if !plan.Permissions.CloudPermissions.IsUnknown() {
		var cloudPermissions []CloudPermissionsValue
		diags := plan.Permissions.CloudPermissions.ElementsAs(ctx, &cloudPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleCloudPermissions []sdk.AddRolesRequestRoleZonesInner
		for _, v := range cloudPermissions {
			addRoleCloudPermissions = append(addRoleCloudPermissions, sdk.AddRolesRequestRoleZonesInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetZones(addRoleCloudPermissions)
	}

	if !plan.Permissions.GroupPermissions.IsUnknown() {
		var groupPermissions []GroupPermissionsValue
		diags := plan.Permissions.GroupPermissions.ElementsAs(ctx, &groupPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleGroupPermissions []sdk.AddRolesRequestRoleSitesInner
		for _, v := range groupPermissions {
			addRoleGroupPermissions = append(addRoleGroupPermissions, sdk.AddRolesRequestRoleSitesInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetSites(addRoleGroupPermissions)
	}

	if !plan.Permissions.InstanceTypePermissions.IsUnknown() {
		var instanceTypePermissions []InstanceTypePermissionsValue
		diags := plan.Permissions.InstanceTypePermissions.ElementsAs(ctx, &instanceTypePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleInstanceTypePermissions []sdk.AddRolesRequestRoleInstanceTypePermissionsInner
		for _, v := range instanceTypePermissions {
			addRoleInstanceTypePermissions = append(addRoleInstanceTypePermissions, sdk.AddRolesRequestRoleInstanceTypePermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetInstanceTypePermissions(addRoleInstanceTypePermissions)
	}

	if !plan.Permissions.PersonaPermissions.IsUnknown() {
		var personaPermissions []PersonaPermissionsValue
		diags := plan.Permissions.PersonaPermissions.ElementsAs(ctx, &personaPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRolePersonaPermissions []sdk.AddRolesRequestRolePersonaPermissionsInner
		for _, v := range personaPermissions {
			addRolePersonaPermissions = append(addRolePersonaPermissions, sdk.AddRolesRequestRolePersonaPermissionsInner{
				Access: v.Access.ValueString(),
				Code:   v.Code.ValueString(),
			})
		}

		addRole.SetPersonaPermissions(addRolePersonaPermissions)
	}

	if !plan.Permissions.ReportTypePermissions.IsUnknown() {
		var reportTypePermissions []ReportTypePermissionsValue
		diags := plan.Permissions.ReportTypePermissions.ElementsAs(ctx, &reportTypePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleReportTypePermissions []sdk.AddRolesRequestRoleReportTypePermissionsInner
		for _, v := range reportTypePermissions {
			addRoleReportTypePermissions = append(addRoleReportTypePermissions, sdk.AddRolesRequestRoleReportTypePermissionsInner{
				Access: v.Access.ValueString(),
				Code:   v.Code.ValueString(),
			})
		}

		addRole.SetReportTypePermissions(addRoleReportTypePermissions)
	}

	if !plan.Permissions.TaskPermissions.IsUnknown() {

		var taskPermissions []TaskPermissionsValue
		diags := plan.Permissions.TaskPermissions.ElementsAs(ctx, &taskPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleTaskPermissions []sdk.AddRolesRequestRoleTaskPermissionsInner
		for _, v := range taskPermissions {
			addRoleTaskPermissions = append(addRoleTaskPermissions, sdk.AddRolesRequestRoleTaskPermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetTaskPermissions(addRoleTaskPermissions)
	}

	if !plan.Permissions.VdiPoolPermissions.IsUnknown() {
		var addRoleVdiPoolPermissions []sdk.AddRolesRequestRoleVdiPoolPermissionsInner
		var vdiPoolPermissions []VdiPoolPermissionsValue
		diags := plan.Permissions.VdiPoolPermissions.ElementsAs(ctx, &vdiPoolPermissions, false)
		if diags.HasError() {
			return diags
		}

		for _, v := range vdiPoolPermissions {
			addRoleVdiPoolPermissions = append(addRoleVdiPoolPermissions, sdk.AddRolesRequestRoleVdiPoolPermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetVdiPoolPermissions(addRoleVdiPoolPermissions)
	}

	if !plan.Permissions.WorkflowPermissions.IsUnknown() {
		var addRoleWorkflowPermissions []sdk.AddRolesRequestRoleTaskSetPermissionsInner
		var workflowPermissions []WorkflowPermissionsValue
		diags := plan.Permissions.WorkflowPermissions.ElementsAs(ctx, &workflowPermissions, false)
		if diags.HasError() {
			return diags
		}

		for _, v := range workflowPermissions {
			addRoleWorkflowPermissions = append(addRoleWorkflowPermissions, sdk.AddRolesRequestRoleTaskSetPermissionsInner{
				Access: v.Access.ValueString(),
				Id:     v.Id.ValueInt64(),
			})
		}

		addRole.SetTaskSetPermissions(addRoleWorkflowPermissions)
	}

	return diags
}

// populate role resource model with current API values
func getRoleAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (RoleModel, diag.Diagnostics) {
	var state RoleModel
	var diags diag.Diagnostics

	r, hresp, err := client.RolesAPI.GetRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate role resource",
			fmt.Sprintf("role %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	permissions, diags := populateGetRoleAsStatePermissions(ctx, r)
	if diags.HasError() {

		return state, diags
	}

	state.Id = convert.Int64ToType(r.Role.Id)
	state.Name = convert.StrToType(r.Role.Name)
	state.Description = convert.StrToType(r.Role.Description.Get())
	state.LandingUrl = convert.StrToType(r.Role.LandingUrl.Get())
	state.Multitenant = convert.BoolToType(r.Role.Multitenant)
	state.MultitenantLocked = convert.BoolToType(r.Role.MultitenantLocked)
	state.RoleType = convert.StrToType(r.Role.RoleType)
	state.Permissions = permissions

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan RoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRole := sdk.NewAddRolesRequestRoleWithDefaults()
	name := plan.Name.ValueString()

	// required
	addRole.SetAuthority(name)

	// optional
	if !plan.Description.IsUnknown() {
		addRole.SetDescription(plan.Description.ValueString())
	}

	if !plan.LandingUrl.IsUnknown() {
		addRole.SetLandingUrl(plan.LandingUrl.ValueString())
	}

	// optional_computed
	if !plan.Multitenant.IsUnknown() {
		// default: false
		addRole.SetMultitenant(plan.Multitenant.ValueBool())
	}

	if !plan.MultitenantLocked.IsUnknown() {
		// default: false
		addRole.SetMultitenantLocked(plan.MultitenantLocked.ValueBool())
	}

	if !plan.RoleType.IsUnknown() {
		// default: user
		addRole.SetRoleType(plan.RoleType.ValueString())
	}

	// Only add to create request if user has set permissions explicitly.
	if !plan.Permissions.IsUnknown() && !plan.Permissions.IsNull() {
		diags := setPermissionsInCreate(ctx, &plan, addRole)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}
	}

	addRoleReq := sdk.NewAddRolesRequest(*addRole)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	role, hresp, err := client.RolesAPI.AddRoles(ctx).
		AddRolesRequest(*addRoleReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if role.GetRole().Id == nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": id is nil",
		)

		return
	}

	id := *role.GetRole().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiState, diags := getRoleAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"create role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	// for optional behaviour on the default access levels
	if plan.Permissions.DefaultBlueprintAccess.IsNull() {
		apiState.Permissions.DefaultBlueprintAccess = types.StringNull()
	}

	if plan.Permissions.DefaultCatalogItemTypeAccess.IsNull() {
		apiState.Permissions.DefaultCatalogItemTypeAccess = types.StringNull()
	}

	if plan.Permissions.DefaultCloudAccess.IsNull() {
		apiState.Permissions.DefaultCloudAccess = types.StringNull()
	}

	if plan.Permissions.DefaultGroupAccess.IsNull() {
		apiState.Permissions.DefaultGroupAccess = types.StringNull()
	}

	if plan.Permissions.DefaultInstanceTypeAccess.IsNull() {
		apiState.Permissions.DefaultInstanceTypeAccess = types.StringNull()
	}

	if plan.Permissions.DefaultPersonaAccess.IsNull() {
		apiState.Permissions.DefaultPersonaAccess = types.StringNull()
	}

	if plan.Permissions.DefaultReportTypeAccess.IsNull() {
		apiState.Permissions.DefaultReportTypeAccess = types.StringNull()
	}

	if plan.Permissions.DefaultTaskAccess.IsNull() {
		apiState.Permissions.DefaultTaskAccess = types.StringNull()
	}

	if plan.Permissions.DefaultVdiPoolAccess.IsNull() {
		apiState.Permissions.DefaultVdiPoolAccess = types.StringNull()
	}

	if plan.Permissions.DefaultWorkflowAccess.IsNull() {
		apiState.Permissions.DefaultWorkflowAccess = types.StringNull()
	}

	// for the case of ommitting permissions field
	if plan.Permissions.IsNull() {
		apiState.Permissions = NewPermissionsValueNull()
	}

	if plan.Permissions.FeaturePermissions.IsNull() {
		apiState.Permissions.FeaturePermissions = types.SetNull(FeaturePermissionsValue{}.Type(ctx))
	}

	// If the user provided a config with feature permissions as part of the create,
	// then set the feature permissions to what was in the plan (optional).
	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {

		// Only feature permissions requires this more complicated create logic.
		// This is because if the user sets feature permissions, we can only store to state
		// the set of feature permissions that were set by the user.
		if !plan.Permissions.FeaturePermissions.IsNull() && !plan.Permissions.FeaturePermissions.IsUnknown() {

			var planFeaturePermissions []FeaturePermissionsValue
			diags := plan.Permissions.FeaturePermissions.ElementsAs(ctx, &planFeaturePermissions, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			var apiStateFeaturePermissions []FeaturePermissionsValue
			diags = apiState.Permissions.FeaturePermissions.ElementsAs(ctx, &apiStateFeaturePermissions, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			for k, v := range planFeaturePermissions {
				if n := slices.IndexFunc(apiStateFeaturePermissions, func(vv FeaturePermissionsValue) bool {
					// We don't know the values of the Id, Name, and SubCategory fields at create time,
					// so we use Code to find those values for v (codes are unique).
					return vv.Code.Equal(v.Code)
				}); n > -1 {
					// If there's a match, update the permissions to store to state with the computed values.
					planFeaturePermissions[k].Id = apiStateFeaturePermissions[n].Id
					planFeaturePermissions[k].Name = apiStateFeaturePermissions[n].Name
					planFeaturePermissions[k].SubCategory = apiStateFeaturePermissions[n].SubCategory
					// We don't need to set planFeaturePermissions[k].state,
					// its value is already attr.ValueStateKnown.
				} else {
					// the case where the permission is not found - error
					resp.Diagnostics.AddError(
						"create role resource",
						fmt.Sprintf("role %d: permission with code %s not found", id, v.Code.String()),
					)

					return
				}
			}

			featuresSetWithComputed, diags := types.SetValueFrom(ctx, FeaturePermissionsValue{}.Type(ctx), planFeaturePermissions)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			apiState.Permissions.FeaturePermissions = featuresSetWithComputed
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state RoleModel

	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read role resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	apiState, diags := getRoleAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"read role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	// for optional behaviour on the default access levels
	if state.Permissions.DefaultBlueprintAccess.IsNull() {
		apiState.Permissions.DefaultBlueprintAccess = types.StringNull()
	}

	if state.Permissions.DefaultCatalogItemTypeAccess.IsNull() {
		apiState.Permissions.DefaultCatalogItemTypeAccess = types.StringNull()
	}

	if state.Permissions.DefaultCloudAccess.IsNull() {
		apiState.Permissions.DefaultCloudAccess = types.StringNull()
	}

	if state.Permissions.DefaultGroupAccess.IsNull() {
		apiState.Permissions.DefaultGroupAccess = types.StringNull()
	}

	if state.Permissions.DefaultInstanceTypeAccess.IsNull() {
		apiState.Permissions.DefaultInstanceTypeAccess = types.StringNull()
	}

	if state.Permissions.DefaultPersonaAccess.IsNull() {
		apiState.Permissions.DefaultPersonaAccess = types.StringNull()
	}

	if state.Permissions.DefaultReportTypeAccess.IsNull() {
		apiState.Permissions.DefaultReportTypeAccess = types.StringNull()
	}

	if state.Permissions.DefaultTaskAccess.IsNull() {
		apiState.Permissions.DefaultTaskAccess = types.StringNull()
	}

	if state.Permissions.DefaultVdiPoolAccess.IsNull() {
		apiState.Permissions.DefaultVdiPoolAccess = types.StringNull()
	}

	if state.Permissions.DefaultWorkflowAccess.IsNull() {
		apiState.Permissions.DefaultWorkflowAccess = types.StringNull()
	}

	if state.Permissions.FeaturePermissions.IsNull() {
		apiState.Permissions.FeaturePermissions = types.SetNull(FeaturePermissionsValue{}.Type(ctx))
	}

	// for the case of ommitting permissions field
	if state.Permissions.IsNull() {
		apiState.Permissions = NewPermissionsValueNull()
	}

	if !state.Permissions.IsNull() && !state.Permissions.IsUnknown() {

		// We extract all feature permissions from API state into a []FeaturePermissionsValue.
		// Then we extract the feature permissions from Terraform state to a []FeaturePermissionsValue.
		// Then we check if the feature permissions in Terraform state are a subset of those in API state.
		// If they are a subset, we use the permissions in state in the Read.
		// We need to do this because the API returns ALL feature permissions in a GET,
		// not just the ones that were overridden by the user.

		if !state.Permissions.FeaturePermissions.IsNull() && !state.Permissions.FeaturePermissions.IsUnknown() {

			var apiStateFeaturePermissions []FeaturePermissionsValue
			diags := apiState.Permissions.FeaturePermissions.ElementsAs(ctx, &apiStateFeaturePermissions, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			var stateFeaturePermissions []FeaturePermissionsValue
			diags = state.Permissions.FeaturePermissions.ElementsAs(ctx, &stateFeaturePermissions, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			for k, v := range stateFeaturePermissions {
				// If apiStateFeaturePermissions contains v with the conditions in the closure...
				if n := slices.IndexFunc(apiStateFeaturePermissions, func(vv FeaturePermissionsValue) bool {
					// We should only compare on code and access, as the other fields are computed.
					// If we compare on the other fields when we have a tainted state with computed values missing,
					// then we'll incorrectly error that the state is not a subset

					// For the case of a tainted state, so we can still find the permissions
					// and get an accurate view of the plan.
					if v.Name.IsUnknown() && v.Id.IsUnknown() && v.SubCategory.IsUnknown() {
						return vv.Code.Equal(v.Code)
					}

					// all other times, when computed state values are OK
					return vv.Id.Equal(v.Id) &&
						vv.Code.Equal(v.Code)
				}); n > -1 {
					// If there's a match, update the permissions to store to state with the computed values.
					stateFeaturePermissions[k].Id = apiStateFeaturePermissions[n].Id
					stateFeaturePermissions[k].Name = apiStateFeaturePermissions[n].Name
					stateFeaturePermissions[k].SubCategory = apiStateFeaturePermissions[n].SubCategory
					// We don't need to set planFeaturePermissions[k].state,
					// its value is already attr.ValueStateKnown.

				} else {
					resp.Diagnostics.AddError(
						"read role resource",
						fmt.Sprintf("role %d: permission with code %s not found", id, v.Code.String()),
					)

					return
				}
			}

			// If we get to here, the permissions in state are a subset of those in API state.
			featuresSetWithComputed, diags := types.SetValueFrom(ctx, FeaturePermissionsValue{}.Type(ctx), stateFeaturePermissions)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			apiState.Permissions.FeaturePermissions = featuresSetWithComputed
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update role resource",
		"update of 'role' resources has not been implemented",
	)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data RoleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	client, _ := r.NewClient(ctx)
	_, hresp, err := client.RolesAPI.DeleteRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete role resource",
			fmt.Sprintf("role %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import role resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)
	if diags.HasError() {
		return
	}

	// We need to set permissions to be empty so that Read will correctly populate it with API values.
	// For import, we're effectively ignoring the IsNull() checks that we've put in place to
	// support the optional typing of the various permissions fields.
	// By doing this, import will populate permissions with all values read from the API,
	// while maintaining the optional behaviour on Create.
	emptyPermissions, diags := NewPermissionsValue(PermissionsValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"default_blueprint_access":         types.StringUnknown(),
		"default_catalog_item_type_access": types.StringUnknown(),
		"default_cloud_access":             types.StringUnknown(),
		"default_group_access":             types.StringUnknown(),
		"default_instance_type_access":     types.StringUnknown(),
		"default_persona_access":           types.StringUnknown(),
		"default_report_type_access":       types.StringUnknown(),
		"default_task_access":              types.StringUnknown(),
		"default_vdi_pool_access":          types.StringUnknown(),
		"default_workflow_access":          types.StringUnknown(),
		"feature_permissions":              types.SetUnknown(FeaturePermissionsValue{}.Type(ctx)),
		"blueprint_permissions":            types.SetUnknown(BlueprintPermissionsValue{}.Type(ctx)),
		"catalog_item_type_permissions":    types.SetUnknown(CatalogItemTypePermissionsValue{}.Type(ctx)),
		"cloud_permissions":                types.SetUnknown(CloudPermissionsValue{}.Type(ctx)),
		"group_permissions":                types.SetUnknown(GroupPermissionsValue{}.Type(ctx)),
		"instance_type_permissions":        types.SetUnknown(InstanceTypePermissionsValue{}.Type(ctx)),
		"persona_permissions":              types.SetUnknown(PersonaPermissionsValue{}.Type(ctx)),
		"report_type_permissions":          types.SetUnknown(ReportTypePermissionsValue{}.Type(ctx)),
		"task_permissions":                 types.SetUnknown(TaskPermissionsValue{}.Type(ctx)),
		"vdi_pool_permissions":             types.SetUnknown(VdiPoolPermissionsValue{}.Type(ctx)),
		"workflow_permissions":             types.SetUnknown(WorkflowPermissionsValue{}.Type(ctx)),
	})
	emptyPermissions.state = attr.ValueStateKnown

	diags = resp.State.SetAttribute(ctx, path.Root("permissions"), emptyPermissions)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(diags...)
}

// This method is called by Terraform's ValidateResourceConfig RPC.
// We use this to perform the validation of permissions specific to user and account roles.
// We need to use the ValidateConfig method as schema validators
// do not have access to config values other than the attribute they're defined for.
// Only user roles can set group permissions.
// Only account roles can set cloud permissions.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config RoleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleType := config.RoleType.ValueString()

	// The ValidateConfigRequest has no knowledge of the plan,
	// so we have to simulate the default value of "user" here.
	if roleType == "" {
		roleType = RoleTypeUser
	}

	// if roleType is "user" and cloud_permissions has been set...
	if roleType == RoleTypeUser &&
		!config.Permissions.CloudPermissions.IsNull() &&
		!config.Permissions.CloudPermissions.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.cloud_permissions"),
			"Conflicting attributes in configuration",
			`cloud_permissions not available for role_type "user". `+
				`Set role_type to "account" to set cloud_permissions.`,
		)

		return
	}

	// if roleType is "user" and default_cloud_access has been set...
	if roleType == RoleTypeUser &&
		!config.Permissions.DefaultCloudAccess.IsNull() &&
		!config.Permissions.DefaultCloudAccess.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.default_cloud_access"),
			"Conflicting attributes in configuration",
			`default_cloud_access not available for role_type "user". `+
				`Set role_type to "account" to set default_cloud_access.`,
		)

		return
	}

	// if roleType is "account" and group_permissions has been set...
	if roleType == RoleTypeAccount &&
		!config.Permissions.GroupPermissions.IsNull() &&
		!config.Permissions.GroupPermissions.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.group_permissions"),
			"Conflicting attributes in configuration",
			`group_permissions not available for role_type "account". `+
				`Set role_type to "user" to set group_permissions.`,
		)

		return
	}

	// if roleType is "account" and default_group_access has been set...
	if roleType == RoleTypeAccount &&
		!config.Permissions.DefaultGroupAccess.IsNull() &&
		!config.Permissions.DefaultGroupAccess.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.default_group_access"),
			"Conflicting attributes in configuration",
			`default_group_access not available for role_type "account". `+
				`Set role_type to "user" to set default_group_access.`,
		)

		return
	}
}
