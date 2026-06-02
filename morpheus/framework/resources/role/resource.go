// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package role

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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
func populateGetRoleAsStatePermissions(
	ctx context.Context,
	r *sdk.GetRole200Response,
) (PermissionsValue, diag.Diagnostics) {
	var features []FeaturePermissionsValue
	for _, v := range r.FeaturePermissions {
		features = append(features, FeaturePermissionsValue{
			Code:        convert.StrToType(v.Code),
			Access:      convert.StrToType(v.Access),
			Id:          convert.Int64ToType(v.Id),
			Name:        convert.StrToType(v.Name),
			SubCategory: convert.StrToType(v.SubCategory),
			state:       attr.ValueStateKnown,
		})
	}

	var blueprints []BlueprintPermissionsValue
	for _, v := range r.AppTemplatePermissions {
		blueprints = append(blueprints, BlueprintPermissionsValue{
			Name:   convert.StrToType(v.Name),
			Id:     convert.Int64ToType(v.Id),
			Access: convert.StrToType(v.Access),
			state:  attr.ValueStateKnown,
		})
	}

	var catalogItemTypes []CatalogItemTypePermissionsValue
	for _, v := range r.CatalogItemTypePermissions {
		catalogItemTypes = append(catalogItemTypes, CatalogItemTypePermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Access: convert.StrToType(v.Access),
			Name:   convert.StrToType(v.Name),
			state:  attr.ValueStateKnown,
		})
	}

	var clouds []CloudPermissionsValue
	for _, v := range r.Zones {
		clouds = append(clouds, CloudPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Access: convert.StrToType(v.Access),
			Name:   convert.StrToType(v.Name),
			state:  attr.ValueStateKnown,
		})
	}

	var groups []GroupPermissionsValue
	for _, v := range r.Sites {
		groups = append(groups, GroupPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Access: convert.StrToType(v.Access),
			Name:   convert.StrToType(v.Name),
			state:  attr.ValueStateKnown,
		})
	}

	var instanceTypes []InstanceTypePermissionsValue
	for _, v := range r.InstanceTypePermissions {
		instanceTypes = append(instanceTypes, InstanceTypePermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
			state:  attr.ValueStateKnown,
		})
	}

	var personas []PersonaPermissionsValue
	for _, v := range r.PersonaPermissions {
		personas = append(personas, PersonaPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
			Code:   convert.StrToType(v.Code),
			state:  attr.ValueStateKnown,
		})
	}

	var reportTypes []ReportTypePermissionsValue
	for _, v := range r.ReportTypePermissions {
		reportTypes = append(reportTypes, ReportTypePermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
			Code:   convert.StrToType(v.Code),
			state:  attr.ValueStateKnown,
		})
	}

	var tasks []TaskPermissionsValue
	for _, v := range r.TaskPermissions {
		tasks = append(tasks, TaskPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
			Code:   types.StringPointerValue(v.Code.Get()),
			state:  attr.ValueStateKnown,
		})
	}

	var vdiPools []VdiPoolPermissionsValue
	for _, v := range r.VdiPoolPermissions {
		vdiPools = append(vdiPools, VdiPoolPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
			state:  attr.ValueStateKnown,
		})
	}

	var workflows []WorkflowPermissionsValue
	for _, v := range r.TaskSetPermissions {
		workflows = append(workflows, WorkflowPermissionsValue{
			Id:     convert.Int64ToType(v.Id),
			Name:   convert.StrToType(v.Name),
			Access: convert.StrToType(v.Access),
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

	catalogItemTypesSet, diags := types.SetValueFrom(
		ctx,
		CatalogItemTypePermissionsValue{}.Type(ctx),
		catalogItemTypes,
	)
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

	instanceTypesSet, diags := types.SetValueFrom(
		ctx,
		InstanceTypePermissionsValue{}.Type(ctx),
		instanceTypes,
	)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	personasSet, diags := types.SetValueFrom(ctx, PersonaPermissionsValue{}.Type(ctx), personas)
	if diags.HasError() {
		return PermissionsValue{}, diags
	}

	reportTypesSet, diags := types.SetValueFrom(
		ctx,
		ReportTypePermissionsValue{}.Type(ctx),
		reportTypes,
	)
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

// Helper function to break out the logic of setting permissions in update.
// It also handles the resetting of permissions that are not in the plan.
// We need to use the values from API state obtained from a prior GET to reset
// those fine-grained permissions that are not in our plan.
func setPermissionsInUpdate(
	ctx context.Context,
	apiState *RoleModel,
	plan *RoleModel,
	updateRole *sdk.UpdateRoleRequestRole,
) diag.Diagnostics {
	var diags diag.Diagnostics

	// `resetAllAccess` is used to make the provider permissions settings
	// behave as overrides.
	// In the Morpheus API, the permissions get reset first, then the new
	// access levels applied.

	// Currently, this field is bugged and doesn't
	// affect the fine-grained access levels of non-feature permissions.
	// So later in this function, we handle the reset logic for those
	// permissions manually.

	updateRole.ResetAllAccess = sdk.PtrBool(true)

	if !plan.Permissions.DefaultBlueprintAccess.IsUnknown() {
		updateRole.GlobalAppTemplateAccess = sdk.PtrString(plan.Permissions.DefaultBlueprintAccess.ValueString())
	}

	if !plan.Permissions.DefaultCatalogItemTypeAccess.IsUnknown() {
		updateRole.GlobalCatalogItemTypeAccess = sdk.PtrString(plan.Permissions.DefaultCatalogItemTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultCloudAccess.IsUnknown() {
		updateRole.GlobalZoneAccess = sdk.PtrString(plan.Permissions.DefaultCloudAccess.ValueString())
	}

	if !plan.Permissions.DefaultGroupAccess.IsUnknown() {
		updateRole.GlobalSiteAccess = sdk.PtrString(plan.Permissions.DefaultGroupAccess.ValueString())
	}

	if !plan.Permissions.DefaultInstanceTypeAccess.IsUnknown() {
		updateRole.GlobalInstanceTypeAccess = sdk.PtrString(plan.Permissions.DefaultInstanceTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultPersonaAccess.IsUnknown() {
		updateRole.GlobalPersonaAccess = sdk.PtrString(plan.Permissions.DefaultPersonaAccess.ValueString())
	}

	if !plan.Permissions.DefaultReportTypeAccess.IsUnknown() {
		updateRole.GlobalReportTypeAccess = sdk.PtrString(plan.Permissions.DefaultReportTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultTaskAccess.IsUnknown() {
		updateRole.GlobalTaskAccess = sdk.PtrString(plan.Permissions.DefaultTaskAccess.ValueString())
	}

	if !plan.Permissions.DefaultVdiPoolAccess.IsUnknown() {
		updateRole.GlobalVdiPoolAccess = sdk.PtrString(plan.Permissions.DefaultVdiPoolAccess.ValueString())
	}

	if !plan.Permissions.DefaultWorkflowAccess.IsUnknown() {
		updateRole.GlobalTaskSetAccess = sdk.PtrString(plan.Permissions.DefaultWorkflowAccess.ValueString())
	}

	if !plan.Permissions.FeaturePermissions.IsUnknown() &&
		!plan.Permissions.FeaturePermissions.IsNull() {
		var planFeaturePermissions []FeaturePermissionsValue
		diags := plan.Permissions.FeaturePermissions.ElementsAs(ctx, &planFeaturePermissions, false)
		if diags.HasError() {
			return diags
		}

		var updateRoleFeaturePermissions []sdk.UpdateRoleRequestRoleFeaturePermissionsInner
		for _, v := range planFeaturePermissions {
			updateRoleFeaturePermissions = append(
				updateRoleFeaturePermissions,
				sdk.UpdateRoleRequestRoleFeaturePermissionsInner{
					Access: v.Access.ValueString(),
					Code:   v.Code.ValueString(),
				},
			)
		}

		updateRole.FeaturePermissions = updateRoleFeaturePermissions
	}

	if !plan.Permissions.BlueprintPermissions.IsUnknown() {

		var updateRoleBlueprintPermissions []sdk.UpdateRoleRequestRoleAppTemplatePermissionsInner

		var apiStateBlueprintPermissions []BlueprintPermissionsValue
		diags = apiState.Permissions.BlueprintPermissions.ElementsAs(
			ctx,
			&apiStateBlueprintPermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.BlueprintPermissions.IsNull() {

			var planBlueprintPermissions []BlueprintPermissionsValue
			diags = plan.Permissions.BlueprintPermissions.ElementsAs(ctx, &planBlueprintPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateBlueprintPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planBlueprintPermissions, func(vv BlueprintPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleBlueprintPermissions = append(
						updateRoleBlueprintPermissions,
						sdk.UpdateRoleRequestRoleAppTemplatePermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planBlueprintPermissions {
				updateRoleBlueprintPermissions = append(
					updateRoleBlueprintPermissions,
					sdk.UpdateRoleRequestRoleAppTemplatePermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.AppTemplatePermissions = updateRoleBlueprintPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateBlueprintPermissions {
				updateRoleBlueprintPermissions = append(
					updateRoleBlueprintPermissions,
					sdk.UpdateRoleRequestRoleAppTemplatePermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.AppTemplatePermissions = updateRoleBlueprintPermissions
		}
	}

	if !plan.Permissions.CatalogItemTypePermissions.IsUnknown() {

		var updateRoleCatalogItemTypePermissions []sdk.UpdateRoleRequestRoleCatalogItemTypePermissionsInner

		var apiStateCatalogItemTypePermissions []CatalogItemTypePermissionsValue
		diags = apiState.Permissions.CatalogItemTypePermissions.ElementsAs(
			ctx,
			&apiStateCatalogItemTypePermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.CatalogItemTypePermissions.IsNull() {

			var planCatalogItemTypePermissions []CatalogItemTypePermissionsValue
			diags = plan.Permissions.CatalogItemTypePermissions.ElementsAs(
				ctx,
				&planCatalogItemTypePermissions,
				false,
			)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateCatalogItemTypePermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(
					planCatalogItemTypePermissions,
					func(vv CatalogItemTypePermissionsValue) bool {
						return vv.Id.Equal(v.Id)
					},
				) {
					updateRoleCatalogItemTypePermissions = append(
						updateRoleCatalogItemTypePermissions,
						sdk.UpdateRoleRequestRoleCatalogItemTypePermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planCatalogItemTypePermissions {
				updateRoleCatalogItemTypePermissions = append(
					updateRoleCatalogItemTypePermissions,
					sdk.UpdateRoleRequestRoleCatalogItemTypePermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.CatalogItemTypePermissions = updateRoleCatalogItemTypePermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateCatalogItemTypePermissions {
				updateRoleCatalogItemTypePermissions = append(
					updateRoleCatalogItemTypePermissions,
					sdk.UpdateRoleRequestRoleCatalogItemTypePermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.CatalogItemTypePermissions = updateRoleCatalogItemTypePermissions
		}
	}

	if !plan.Permissions.CloudPermissions.IsUnknown() {

		var updateRoleCloudPermissions []sdk.UpdateRoleRequestRoleZonesInner

		var apiStateCloudPermissions []CloudPermissionsValue
		diags = apiState.Permissions.CloudPermissions.ElementsAs(ctx, &apiStateCloudPermissions, false)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.CloudPermissions.IsNull() {

			var planCloudPermissions []CloudPermissionsValue
			diags := plan.Permissions.CloudPermissions.ElementsAs(ctx, &planCloudPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateCloudPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planCloudPermissions, func(vv CloudPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleCloudPermissions = append(
						updateRoleCloudPermissions,
						sdk.UpdateRoleRequestRoleZonesInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planCloudPermissions {
				updateRoleCloudPermissions = append(
					updateRoleCloudPermissions,
					sdk.UpdateRoleRequestRoleZonesInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.Zones = updateRoleCloudPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateCloudPermissions {
				updateRoleCloudPermissions = append(
					updateRoleCloudPermissions,
					sdk.UpdateRoleRequestRoleZonesInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.Zones = updateRoleCloudPermissions
		}
	}

	if !plan.Permissions.GroupPermissions.IsUnknown() {

		var updateRoleGroupPermissions []sdk.UpdateRoleRequestRoleSitesInner

		var apiStateGroupPermissions []GroupPermissionsValue
		diags = apiState.Permissions.GroupPermissions.ElementsAs(ctx, &apiStateGroupPermissions, false)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.GroupPermissions.IsNull() {

			var planGroupPermissions []GroupPermissionsValue
			diags := plan.Permissions.GroupPermissions.ElementsAs(ctx, &planGroupPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateGroupPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planGroupPermissions, func(vv GroupPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleGroupPermissions = append(
						updateRoleGroupPermissions,
						sdk.UpdateRoleRequestRoleSitesInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planGroupPermissions {
				updateRoleGroupPermissions = append(
					updateRoleGroupPermissions,
					sdk.UpdateRoleRequestRoleSitesInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.Sites = updateRoleGroupPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateGroupPermissions {
				updateRoleGroupPermissions = append(
					updateRoleGroupPermissions,
					sdk.UpdateRoleRequestRoleSitesInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.Sites = updateRoleGroupPermissions
		}
	}

	if !plan.Permissions.InstanceTypePermissions.IsUnknown() {

		var updateRoleInstanceTypePermissions []sdk.UpdateRoleRequestRoleInstanceTypePermissionsInner

		var apiStateInstanceTypePermissions []InstanceTypePermissionsValue
		diags = apiState.Permissions.InstanceTypePermissions.ElementsAs(
			ctx,
			&apiStateInstanceTypePermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.InstanceTypePermissions.IsNull() {

			var planInstanceTypePermissions []InstanceTypePermissionsValue
			diags := plan.Permissions.InstanceTypePermissions.ElementsAs(
				ctx,
				&planInstanceTypePermissions,
				false,
			)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateInstanceTypePermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(
					planInstanceTypePermissions,
					func(vv InstanceTypePermissionsValue) bool {
						return vv.Id.Equal(v.Id)
					},
				) {
					updateRoleInstanceTypePermissions = append(
						updateRoleInstanceTypePermissions,
						sdk.UpdateRoleRequestRoleInstanceTypePermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planInstanceTypePermissions {
				updateRoleInstanceTypePermissions = append(
					updateRoleInstanceTypePermissions,
					sdk.UpdateRoleRequestRoleInstanceTypePermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.InstanceTypePermissions = updateRoleInstanceTypePermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateInstanceTypePermissions {
				updateRoleInstanceTypePermissions = append(
					updateRoleInstanceTypePermissions,
					sdk.UpdateRoleRequestRoleInstanceTypePermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.InstanceTypePermissions = updateRoleInstanceTypePermissions
		}
	}

	if !plan.Permissions.PersonaPermissions.IsUnknown() {

		var updateRolePersonaPermissions []sdk.UpdateRoleRequestRolePersonaPermissionsInner

		var apiStatePersonaPermissions []PersonaPermissionsValue
		diags = apiState.Permissions.PersonaPermissions.ElementsAs(
			ctx,
			&apiStatePersonaPermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.PersonaPermissions.IsNull() {

			var planPersonaPermissions []PersonaPermissionsValue
			diags := plan.Permissions.PersonaPermissions.ElementsAs(ctx, &planPersonaPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStatePersonaPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planPersonaPermissions, func(vv PersonaPermissionsValue) bool {
					return vv.Code.Equal(v.Code)
				}) {
					updateRolePersonaPermissions = append(
						updateRolePersonaPermissions,
						sdk.UpdateRoleRequestRolePersonaPermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Code:   v.Code.ValueString(),
						},
					)
				}
			}

			for _, v := range planPersonaPermissions {
				updateRolePersonaPermissions = append(
					updateRolePersonaPermissions,
					sdk.UpdateRoleRequestRolePersonaPermissionsInner{
						Access: v.Access.ValueString(),
						Code:   v.Code.ValueString(),
					},
				)
			}

			updateRole.PersonaPermissions = updateRolePersonaPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStatePersonaPermissions {
				updateRolePersonaPermissions = append(
					updateRolePersonaPermissions,
					sdk.UpdateRoleRequestRolePersonaPermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Code:   v.Code.ValueString(),
					},
				)
			}

			updateRole.PersonaPermissions = updateRolePersonaPermissions
		}
	}

	if !plan.Permissions.ReportTypePermissions.IsUnknown() {

		var updateRoleReportTypePermissions []sdk.UpdateRoleRequestRoleReportTypePermissionsInner

		var apiStateReportTypePermissions []ReportTypePermissionsValue
		diags = apiState.Permissions.ReportTypePermissions.ElementsAs(
			ctx,
			&apiStateReportTypePermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.ReportTypePermissions.IsNull() {

			var planReportTypePermissions []ReportTypePermissionsValue
			diags := plan.Permissions.ReportTypePermissions.ElementsAs(
				ctx,
				&planReportTypePermissions,
				false,
			)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateReportTypePermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planReportTypePermissions, func(vv ReportTypePermissionsValue) bool {
					return vv.Code.Equal(v.Code)
				}) {
					updateRoleReportTypePermissions = append(
						updateRoleReportTypePermissions,
						sdk.UpdateRoleRequestRoleReportTypePermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Code:   v.Code.ValueString(),
						},
					)
				}
			}

			for _, v := range planReportTypePermissions {
				updateRoleReportTypePermissions = append(
					updateRoleReportTypePermissions,
					sdk.UpdateRoleRequestRoleReportTypePermissionsInner{
						Access: v.Access.ValueString(),
						Code:   v.Code.ValueString(),
					},
				)
			}

			updateRole.ReportTypePermissions = updateRoleReportTypePermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateReportTypePermissions {
				updateRoleReportTypePermissions = append(
					updateRoleReportTypePermissions,
					sdk.UpdateRoleRequestRoleReportTypePermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Code:   v.Code.ValueString(),
					},
				)
			}

			updateRole.ReportTypePermissions = updateRoleReportTypePermissions
		}
	}

	if !plan.Permissions.TaskPermissions.IsUnknown() {

		var updateRoleTaskPermissions []sdk.UpdateRoleRequestRoleTaskPermissionsInner

		var apiStateTaskPermissions []TaskPermissionsValue
		diags = apiState.Permissions.TaskPermissions.ElementsAs(ctx, &apiStateTaskPermissions, false)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.TaskPermissions.IsNull() {

			var planTaskPermissions []TaskPermissionsValue
			diags := plan.Permissions.TaskPermissions.ElementsAs(ctx, &planTaskPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateTaskPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planTaskPermissions, func(vv TaskPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleTaskPermissions = append(
						updateRoleTaskPermissions,
						sdk.UpdateRoleRequestRoleTaskPermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planTaskPermissions {
				updateRoleTaskPermissions = append(
					updateRoleTaskPermissions,
					sdk.UpdateRoleRequestRoleTaskPermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.TaskPermissions = updateRoleTaskPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateTaskPermissions {
				updateRoleTaskPermissions = append(
					updateRoleTaskPermissions,
					sdk.UpdateRoleRequestRoleTaskPermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.TaskPermissions = updateRoleTaskPermissions
		}
	}

	if !plan.Permissions.VdiPoolPermissions.IsUnknown() {

		var updateRoleVdiPoolPermissions []sdk.UpdateRoleRequestRoleVdiPoolPermissionsInner

		var apiStateVdiPoolPermissions []VdiPoolPermissionsValue
		diags = apiState.Permissions.VdiPoolPermissions.ElementsAs(
			ctx,
			&apiStateVdiPoolPermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.VdiPoolPermissions.IsNull() {

			var planVdiPoolPermissions []VdiPoolPermissionsValue
			diags := plan.Permissions.VdiPoolPermissions.ElementsAs(ctx, &planVdiPoolPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateVdiPoolPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planVdiPoolPermissions, func(vv VdiPoolPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleVdiPoolPermissions = append(
						updateRoleVdiPoolPermissions,
						sdk.UpdateRoleRequestRoleVdiPoolPermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planVdiPoolPermissions {
				updateRoleVdiPoolPermissions = append(
					updateRoleVdiPoolPermissions,
					sdk.UpdateRoleRequestRoleVdiPoolPermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.VdiPoolPermissions = updateRoleVdiPoolPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateVdiPoolPermissions {
				updateRoleVdiPoolPermissions = append(
					updateRoleVdiPoolPermissions,
					sdk.UpdateRoleRequestRoleVdiPoolPermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.VdiPoolPermissions = updateRoleVdiPoolPermissions
		}
	}

	if !plan.Permissions.WorkflowPermissions.IsUnknown() {

		var updateRoleWorkflowPermissions []sdk.UpdateRoleRequestRoleTaskSetPermissionsInner

		var apiStateWorkflowPermissions []WorkflowPermissionsValue
		diags = apiState.Permissions.WorkflowPermissions.ElementsAs(
			ctx,
			&apiStateWorkflowPermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		if !plan.Permissions.WorkflowPermissions.IsNull() {

			var planWorkflowPermissions []WorkflowPermissionsValue
			diags := plan.Permissions.WorkflowPermissions.ElementsAs(ctx, &planWorkflowPermissions, false)
			if diags.HasError() {
				return diags
			}

			for _, v := range apiStateWorkflowPermissions {
				// If the permission setting exists in API state, but
				// NOT in the plan, then reset it to "default".
				if !slices.ContainsFunc(planWorkflowPermissions, func(vv WorkflowPermissionsValue) bool {
					return vv.Id.Equal(v.Id)
				}) {
					updateRoleWorkflowPermissions = append(
						updateRoleWorkflowPermissions,
						sdk.UpdateRoleRequestRoleTaskSetPermissionsInner{
							Access: DefaultPermissionAccessLevel,
							Id:     v.Id.ValueInt64(),
						},
					)
				}
			}

			for _, v := range planWorkflowPermissions {
				updateRoleWorkflowPermissions = append(
					updateRoleWorkflowPermissions,
					sdk.UpdateRoleRequestRoleTaskSetPermissionsInner{
						Access: v.Access.ValueString(),
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.TaskSetPermissions = updateRoleWorkflowPermissions

		} else {
			// For when we remove permissions from config.
			// Resets everything obtained from the GET to their default values.
			for _, v := range apiStateWorkflowPermissions {
				updateRoleWorkflowPermissions = append(
					updateRoleWorkflowPermissions,
					sdk.UpdateRoleRequestRoleTaskSetPermissionsInner{
						Access: DefaultPermissionAccessLevel,
						Id:     v.Id.ValueInt64(),
					},
				)
			}

			updateRole.TaskSetPermissions = updateRoleWorkflowPermissions
		}
	}

	return diags
}

// Helper function to break out the logic of setting permissions in create.
func setPermissionsInCreate(
	ctx context.Context,
	plan *RoleModel,
	addRole *sdk.AddRolesRequestRole,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if !plan.Permissions.DefaultBlueprintAccess.IsUnknown() {
		addRole.GlobalAppTemplateAccess = sdk.PtrString(plan.Permissions.DefaultBlueprintAccess.ValueString())
	}

	if !plan.Permissions.DefaultCatalogItemTypeAccess.IsUnknown() {
		addRole.GlobalCatalogItemTypeAccess = sdk.PtrString(plan.Permissions.DefaultCatalogItemTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultCloudAccess.IsUnknown() {
		addRole.GlobalZoneAccess = sdk.PtrString(plan.Permissions.DefaultCloudAccess.ValueString())
	}

	if !plan.Permissions.DefaultGroupAccess.IsUnknown() {
		addRole.GlobalSiteAccess = sdk.PtrString(plan.Permissions.DefaultGroupAccess.ValueString())
	}

	if !plan.Permissions.DefaultInstanceTypeAccess.IsUnknown() {
		addRole.GlobalInstanceTypeAccess = sdk.PtrString(plan.Permissions.DefaultInstanceTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultPersonaAccess.IsUnknown() {
		addRole.GlobalPersonaAccess = sdk.PtrString(plan.Permissions.DefaultPersonaAccess.ValueString())
	}

	if !plan.Permissions.DefaultReportTypeAccess.IsUnknown() {
		addRole.GlobalReportTypeAccess = sdk.PtrString(plan.Permissions.DefaultReportTypeAccess.ValueString())
	}

	if !plan.Permissions.DefaultTaskAccess.IsUnknown() {
		addRole.GlobalTaskAccess = sdk.PtrString(plan.Permissions.DefaultTaskAccess.ValueString())
	}

	if !plan.Permissions.DefaultVdiPoolAccess.IsUnknown() {
		addRole.GlobalVdiPoolAccess = sdk.PtrString(plan.Permissions.DefaultVdiPoolAccess.ValueString())
	}

	if !plan.Permissions.DefaultWorkflowAccess.IsUnknown() {
		addRole.GlobalTaskSetAccess = sdk.PtrString(plan.Permissions.DefaultWorkflowAccess.ValueString())
	}

	if !plan.Permissions.FeaturePermissions.IsUnknown() {
		var featurePermissions []FeaturePermissionsValue
		diags := plan.Permissions.FeaturePermissions.ElementsAs(ctx, &featurePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleFeaturePermissions []sdk.AddRolesRequestRoleFeaturePermissionsInner
		for _, v := range featurePermissions {
			addRoleFeaturePermissions = append(
				addRoleFeaturePermissions,
				sdk.AddRolesRequestRoleFeaturePermissionsInner{
					Access: v.Access.ValueString(),
					Code:   v.Code.ValueString(),
				},
			)
		}

		addRole.FeaturePermissions = addRoleFeaturePermissions
	}

	if !plan.Permissions.BlueprintPermissions.IsUnknown() {
		var blueprintPermissions []BlueprintPermissionsValue
		diags = plan.Permissions.BlueprintPermissions.ElementsAs(ctx, &blueprintPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleBlueprintPermissions []sdk.AddRolesRequestRoleAppTemplatePermissionsInner
		for _, v := range blueprintPermissions {
			addRoleBlueprintPermissions = append(
				addRoleBlueprintPermissions,
				sdk.AddRolesRequestRoleAppTemplatePermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.AppTemplatePermissions = addRoleBlueprintPermissions
	}

	if !plan.Permissions.CatalogItemTypePermissions.IsUnknown() {
		var catalogItemTypePermissions []CatalogItemTypePermissionsValue
		diags = plan.Permissions.CatalogItemTypePermissions.ElementsAs(
			ctx,
			&catalogItemTypePermissions,
			false,
		)
		if diags.HasError() {
			return diags
		}

		var addRoleCatalogItemTypePermissions []sdk.AddRolesRequestRoleCatalogItemTypePermissionsInner
		for _, v := range catalogItemTypePermissions {
			addRoleCatalogItemTypePermissions = append(
				addRoleCatalogItemTypePermissions,
				sdk.AddRolesRequestRoleCatalogItemTypePermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.CatalogItemTypePermissions = addRoleCatalogItemTypePermissions
	}

	if !plan.Permissions.CloudPermissions.IsUnknown() {
		var cloudPermissions []CloudPermissionsValue
		diags := plan.Permissions.CloudPermissions.ElementsAs(ctx, &cloudPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleCloudPermissions []sdk.AddRolesRequestRoleZonesInner
		for _, v := range cloudPermissions {
			addRoleCloudPermissions = append(
				addRoleCloudPermissions,
				sdk.AddRolesRequestRoleZonesInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.Zones = addRoleCloudPermissions
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

		addRole.Sites = addRoleGroupPermissions
	}

	if !plan.Permissions.InstanceTypePermissions.IsUnknown() {
		var instanceTypePermissions []InstanceTypePermissionsValue
		diags := plan.Permissions.InstanceTypePermissions.ElementsAs(ctx, &instanceTypePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleInstanceTypePermissions []sdk.AddRolesRequestRoleInstanceTypePermissionsInner
		for _, v := range instanceTypePermissions {
			addRoleInstanceTypePermissions = append(
				addRoleInstanceTypePermissions,
				sdk.AddRolesRequestRoleInstanceTypePermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.InstanceTypePermissions = addRoleInstanceTypePermissions
	}

	if !plan.Permissions.PersonaPermissions.IsUnknown() {
		var personaPermissions []PersonaPermissionsValue
		diags := plan.Permissions.PersonaPermissions.ElementsAs(ctx, &personaPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRolePersonaPermissions []sdk.AddRolesRequestRolePersonaPermissionsInner
		for _, v := range personaPermissions {
			addRolePersonaPermissions = append(
				addRolePersonaPermissions,
				sdk.AddRolesRequestRolePersonaPermissionsInner{
					Access: v.Access.ValueString(),
					Code:   v.Code.ValueString(),
				},
			)
		}

		addRole.PersonaPermissions = addRolePersonaPermissions
	}

	if !plan.Permissions.ReportTypePermissions.IsUnknown() {
		var reportTypePermissions []ReportTypePermissionsValue
		diags := plan.Permissions.ReportTypePermissions.ElementsAs(ctx, &reportTypePermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleReportTypePermissions []sdk.AddRolesRequestRoleReportTypePermissionsInner
		for _, v := range reportTypePermissions {
			addRoleReportTypePermissions = append(
				addRoleReportTypePermissions,
				sdk.AddRolesRequestRoleReportTypePermissionsInner{
					Access: v.Access.ValueString(),
					Code:   v.Code.ValueString(),
				},
			)
		}

		addRole.ReportTypePermissions = addRoleReportTypePermissions
	}

	if !plan.Permissions.TaskPermissions.IsUnknown() {

		var taskPermissions []TaskPermissionsValue
		diags := plan.Permissions.TaskPermissions.ElementsAs(ctx, &taskPermissions, false)
		if diags.HasError() {
			return diags
		}

		var addRoleTaskPermissions []sdk.AddRolesRequestRoleTaskPermissionsInner
		for _, v := range taskPermissions {
			addRoleTaskPermissions = append(
				addRoleTaskPermissions,
				sdk.AddRolesRequestRoleTaskPermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.TaskPermissions = addRoleTaskPermissions
	}

	if !plan.Permissions.VdiPoolPermissions.IsUnknown() {
		var addRoleVdiPoolPermissions []sdk.AddRolesRequestRoleVdiPoolPermissionsInner
		var vdiPoolPermissions []VdiPoolPermissionsValue
		diags := plan.Permissions.VdiPoolPermissions.ElementsAs(ctx, &vdiPoolPermissions, false)
		if diags.HasError() {
			return diags
		}

		for _, v := range vdiPoolPermissions {
			addRoleVdiPoolPermissions = append(
				addRoleVdiPoolPermissions,
				sdk.AddRolesRequestRoleVdiPoolPermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.VdiPoolPermissions = addRoleVdiPoolPermissions
	}

	if !plan.Permissions.WorkflowPermissions.IsUnknown() {
		var addRoleWorkflowPermissions []sdk.AddRolesRequestRoleTaskSetPermissionsInner
		var workflowPermissions []WorkflowPermissionsValue
		diags := plan.Permissions.WorkflowPermissions.ElementsAs(ctx, &workflowPermissions, false)
		if diags.HasError() {
			return diags
		}

		for _, v := range workflowPermissions {
			addRoleWorkflowPermissions = append(
				addRoleWorkflowPermissions,
				sdk.AddRolesRequestRoleTaskSetPermissionsInner{
					Access: v.Access.ValueString(),
					Id:     v.Id.ValueInt64(),
				},
			)
		}

		addRole.TaskSetPermissions = addRoleWorkflowPermissions
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
			fmt.Sprintf("role %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	permissions, diags := populateGetRoleAsStatePermissions(ctx, r)
	if diags.HasError() {
		return state, diags
	}

	state.Id = convert.Int64ToType(r.Role.Id)
	state.Name = convert.StrToType(r.Role.Name)
	state.DefaultPersonaCode = types.StringNull()
	if r.Role.DefaultPersona != nil {
		state.DefaultPersonaCode = convert.StrToType(r.Role.DefaultPersona.Code)
	}
	state.Description = convert.StrToType(r.Role.Description.Get())
	state.LandingUrl = convert.StrToType(r.Role.LandingUrl.Get())
	state.Multitenant = convert.BoolToType(r.Role.Multitenant)
	state.MultitenantLocked = convert.BoolToType(r.Role.MultitenantLocked)
	state.RoleType = convert.StrToType(r.Role.RoleType)
	state.Permissions = permissions

	// Convert the `account` role type from API to `tenant` Tor Terraform
	if state.RoleType.ValueString() == RoleTypeAccountAPI {
		state.RoleType = types.StringValue(RoleTypeTenant)
	}

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

	addRole := &sdk.AddRolesRequestRole{}
	name := plan.Name.ValueString()

	// required
	addRole.Authority = name

	// optional
	if !plan.DefaultPersonaCode.IsUnknown() && !plan.DefaultPersonaCode.IsNull() {
		addRole.DefaultPersona.Set(sdk.PtrString(plan.DefaultPersonaCode.ValueString()))
	}

	if !plan.Description.IsUnknown() {
		addRole.Description.Set(sdk.PtrString(plan.Description.ValueString()))
	}

	if !plan.LandingUrl.IsUnknown() {
		addRole.LandingUrl.Set(sdk.PtrString(plan.LandingUrl.ValueString()))
	}

	// optional_computed
	if !plan.Multitenant.IsUnknown() {
		// default: false
		addRole.Multitenant = sdk.PtrBool(plan.Multitenant.ValueBool())
	}

	if !plan.MultitenantLocked.IsUnknown() {
		// default: false
		addRole.MultitenantLocked = sdk.PtrBool(plan.MultitenantLocked.ValueBool())
	}

	if !plan.RoleType.IsUnknown() {
		// default: user
		if plan.RoleType.ValueString() == RoleTypeUser {
			addRole.RoleType = sdk.PtrString(plan.RoleType.ValueString())
		}

		if plan.RoleType.ValueString() == RoleTypeTenant {
			addRole.RoleType = sdk.PtrString(RoleTypeAccountAPI)
		}
	}

	// Only add to create request if user has set permissions explicitly.
	if !plan.Permissions.IsUnknown() && !plan.Permissions.IsNull() {
		diags := setPermissionsInCreate(ctx, &plan, addRole)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}
	}

	addRoleReq := &sdk.AddRolesRequest{Role: *addRole}

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
			"role "+name+" POST failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if role.Role == nil || role.Role.Id == nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": id is nil",
		)

		return
	}

	id := *role.Role.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "role",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	apiState, diags := getRoleAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"create role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)
		taintResourceState(id)

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

	// for the case of omitting permissions field
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
		if !plan.Permissions.FeaturePermissions.IsNull() &&
			!plan.Permissions.FeaturePermissions.IsUnknown() {

			var planFeaturePermissions []FeaturePermissionsValue
			diags := plan.Permissions.FeaturePermissions.ElementsAs(
				ctx,
				&planFeaturePermissions,
				false,
			)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				taintResourceState(id)

				return
			}

			var apiStateFeaturePermissions []FeaturePermissionsValue
			diags = apiState.Permissions.FeaturePermissions.ElementsAs(
				ctx,
				&apiStateFeaturePermissions,
				false,
			)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				taintResourceState(id)

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
					taintResourceState(id)

					return
				}
			}

			featuresSetWithComputed, diags := types.SetValueFrom(
				ctx,
				FeaturePermissionsValue{}.Type(ctx),
				planFeaturePermissions,
			)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				taintResourceState(id)

				return
			}

			apiState.Permissions.FeaturePermissions = featuresSetWithComputed
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set role state",
			fmt.Sprintf("Role %d was created but state could not be saved", id),
		)
		taintResourceState(id)

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

	importing := state.Name.IsNull()
	if importing {
		// Convert the legacy `custom` defaultAccessLevel value to `none` on import
		// As per API behaviour, `custom` gets treated as `none` since
		// the change to the newer permissions model (the one we use in the provider).
		// By doing this, we can correctly handle the importing of legacy Roles using `custom`
		// which were previously unmanaged by Terraform.
		const customVal = "custom"
		noneVal := "none"
		if apiState.Permissions.DefaultBlueprintAccess.ValueString() == customVal {
			apiState.Permissions.DefaultBlueprintAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultCatalogItemTypeAccess.ValueString() == customVal {
			apiState.Permissions.DefaultCatalogItemTypeAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultCloudAccess.ValueString() == customVal {
			apiState.Permissions.DefaultCloudAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultGroupAccess.ValueString() == customVal {
			apiState.Permissions.DefaultGroupAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultInstanceTypeAccess.ValueString() == customVal {
			apiState.Permissions.DefaultInstanceTypeAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultPersonaAccess.ValueString() == customVal {
			apiState.Permissions.DefaultPersonaAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultReportTypeAccess.ValueString() == customVal {
			apiState.Permissions.DefaultReportTypeAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultTaskAccess.ValueString() == customVal {
			apiState.Permissions.DefaultTaskAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultVdiPoolAccess.ValueString() == customVal {
			apiState.Permissions.DefaultVdiPoolAccess = convert.StrToType(&noneVal)
		}

		if apiState.Permissions.DefaultWorkflowAccess.ValueString() == customVal {
			apiState.Permissions.DefaultWorkflowAccess = convert.StrToType(&noneVal)
		}
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

	// for the case of omitting permissions field
	if state.Permissions.IsNull() {
		apiState.Permissions = NewPermissionsValueNull()
	}

	if !state.Permissions.IsNull() && !state.Permissions.IsUnknown() {
		// We extract all feature permissions from API state into a []FeaturePermissionsValue.
		// Then we extract the feature permissions from Terraform state to a []FeaturePermissionsValue.
		// Then we check if the feature permissions in Terraform state are a subset of those in API state.
		// If they are a subset, we set the values from the GET in state.

		// We need to do this because the API returns ALL feature permissions in a GET,
		// not just the ones whose default values were overridden by the user.

		if !state.Permissions.FeaturePermissions.IsNull() &&
			!state.Permissions.FeaturePermissions.IsUnknown() {

			var apiStateFeaturePermissions []FeaturePermissionsValue
			diags := apiState.Permissions.FeaturePermissions.ElementsAs(
				ctx,
				&apiStateFeaturePermissions,
				false,
			)
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
					// We set access to detect drift in API and state
					stateFeaturePermissions[k].Access = apiStateFeaturePermissions[n].Access
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
			featuresSetWithComputed, diags := types.SetValueFrom(
				ctx,
				FeaturePermissionsValue{}.Type(ctx),
				stateFeaturePermissions,
			)
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
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state RoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRole := &sdk.UpdateRoleRequestRole{}
	id := plan.Id.ValueInt64()

	// required - authority (name)
	updateRole.Authority = sdk.PtrString(plan.Name.ValueString())

	// optional fields
	if plan.DefaultPersonaCode.IsNull() {
		updateRole.DefaultPersona.Set(nil)
	} else {
		updateRole.DefaultPersona.Set(sdk.PtrString(plan.DefaultPersonaCode.ValueString()))
	}

	if plan.Description.IsNull() {
		updateRole.Description.Set(nil)
	} else {
		updateRole.Description.Set(sdk.PtrString(plan.Description.ValueString()))
	}

	if plan.LandingUrl.IsNull() {
		updateRole.LandingUrl.Set(nil)
	} else {
		updateRole.LandingUrl.Set(sdk.PtrString(plan.LandingUrl.ValueString()))
	}

	if !plan.Multitenant.IsNull() {
		updateRole.Multitenant = sdk.PtrBool(plan.Multitenant.ValueBool())
	}

	if !plan.MultitenantLocked.IsNull() {
		updateRole.MultitenantLocked = sdk.PtrBool(plan.MultitenantLocked.ValueBool())
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d: failed to create client: ", id)+err.Error(),
		)

		return
	}

	// The update section:
	// 1. Perform a GET so we know which non-feature permissions to reset.
	// 2. Perform a PUT to both reset the existing permissions levels and
	// apply the permissions levels from the Terraform plan in the same PUT.

	// Doing the steps in that order will ensure that our Terraform config
	// will act as an override for defaults.

	// 1. Perform a GET so we know which non-feature permissions to reset
	getRole, diags := getRoleAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	// Set permissions regardless of whether the
	// permissions block is Null or Unknown.
	// This allows us to reset permissions levels even when the
	// permissions block has been removed from the config.
	diags = setPermissionsInUpdate(ctx, &getRole, &plan, updateRole)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. Perform a PUT to apply the permissions levels from the Terraform plan.
	updateRoleReq := &sdk.UpdateRoleRequest{Role: *updateRole}

	role, hresp, err := client.RolesAPI.UpdateRole(ctx, id).
		UpdateRoleRequest(*updateRoleReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d PUT failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if role.Role == nil || role.Role.Id == nil {
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d: id is nil", id),
		)

		return
	}

	newID := *role.Role.Id
	if newID != id {
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d: id mismatch %d != %d", id, id, newID),
		)

		return
	}

	apiState, diags := getRoleAsState(ctx, newID, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	// Handle optional behavior for default access levels (similar to Create/Read)
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

	// Handle computed feature permissions (similar to Create/Read)
	if plan.Permissions.FeaturePermissions.IsNull() {
		apiState.Permissions.FeaturePermissions = types.SetNull(FeaturePermissionsValue{}.Type(ctx))
	}

	// Handle permissions omission case
	if plan.Permissions.IsNull() {
		apiState.Permissions = NewPermissionsValueNull()
	}

	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		// We extract all feature permissions from API state into a []FeaturePermissionsValue.
		// Then we extract the feature permissions from Terraform state to a []FeaturePermissionsValue.
		// Then we check if the feature permissions in Terraform state are a subset of those in API state.
		// If they are a subset, we set the values from the GET in state.

		// We need to do this because the API returns ALL feature permissions in a GET,
		// not just the ones whose default values were overridden by the user.

		if !plan.Permissions.FeaturePermissions.IsNull() &&
			!plan.Permissions.FeaturePermissions.IsUnknown() {

			var apiStateFeaturePermissions []FeaturePermissionsValue
			diags := apiState.Permissions.FeaturePermissions.ElementsAs(
				ctx,
				&apiStateFeaturePermissions,
				false,
			)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			var planFeaturePermissions []FeaturePermissionsValue
			diags = plan.Permissions.FeaturePermissions.ElementsAs(ctx, &planFeaturePermissions, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			for k, v := range planFeaturePermissions {
				// If apiStateFeaturePermissions contains v with the conditions in the closure...
				if n := slices.IndexFunc(apiStateFeaturePermissions, func(vv FeaturePermissionsValue) bool {
					// We should only compare on code as it acts as an ID and the other fields are computed.
					return vv.Code.Equal(v.Code)
				}); n > -1 {
					// If there's a match, update the permissions to store to state with the computed values.
					// We set access to detect drift in API and plan
					planFeaturePermissions[k].Access = apiStateFeaturePermissions[n].Access
					planFeaturePermissions[k].Id = apiStateFeaturePermissions[n].Id
					planFeaturePermissions[k].Name = apiStateFeaturePermissions[n].Name
					planFeaturePermissions[k].SubCategory = apiStateFeaturePermissions[n].SubCategory
					// We don't need to set planFeaturePermissions[k].state;
					// its value is already attr.ValueStateKnown.

				} else {
					resp.Diagnostics.AddError(
						"update role resource",
						fmt.Sprintf("role %d: permission with code %s not found", id, v.Code.String()),
					)

					return
				}
			}

			// If we get to here, the permissions in plan + state are a subset of those in API state.
			featuresSetWithComputed, diags := types.SetValueFrom(
				ctx,
				FeaturePermissionsValue{}.Type(ctx),
				planFeaturePermissions,
			)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)

				return
			}

			apiState.Permissions.FeaturePermissions = featuresSetWithComputed
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
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
			fmt.Sprintf("role %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
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
	emptyPermissions, diags := NewPermissionsValue(
		PermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
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
			"catalog_item_type_permissions": types.SetUnknown(
				CatalogItemTypePermissionsValue{}.Type(ctx),
			),
			"cloud_permissions":         types.SetUnknown(CloudPermissionsValue{}.Type(ctx)),
			"group_permissions":         types.SetUnknown(GroupPermissionsValue{}.Type(ctx)),
			"instance_type_permissions": types.SetUnknown(InstanceTypePermissionsValue{}.Type(ctx)),
			"persona_permissions":       types.SetUnknown(PersonaPermissionsValue{}.Type(ctx)),
			"report_type_permissions":   types.SetUnknown(ReportTypePermissionsValue{}.Type(ctx)),
			"task_permissions":          types.SetUnknown(TaskPermissionsValue{}.Type(ctx)),
			"vdi_pool_permissions":      types.SetUnknown(VdiPoolPermissionsValue{}.Type(ctx)),
			"workflow_permissions":      types.SetUnknown(WorkflowPermissionsValue{}.Type(ctx)),
		},
	)
	if diags.HasError() {
		return
	}
	emptyPermissions.state = attr.ValueStateKnown

	diags = resp.State.SetAttribute(ctx, path.Root("permissions"), emptyPermissions)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(diags...)
}

// This method is called by Terraform's ValidateResourceConfig RPC.
// We use this to perform the validation of attributes specific to user and tenant roles.
// We need to use the ValidateConfig method as schema validators
// do not have access to config values other than the attribute they're defined for.
// Only user roles can set group permissions.
// Only tenant roles can set cloud permissions.
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
			`cloud_permissions not available for role_type "`+RoleTypeUser+`". `+
				`Set role_type to "`+RoleTypeTenant+`" to set cloud_permissions.`,
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
			`default_cloud_access not available for role_type "`+RoleTypeUser+`". `+
				`Set role_type to "`+RoleTypeTenant+`" to set default_cloud_access.`,
		)

		return
	}

	// if roleType is "tenant" and group_permissions has been set...
	if roleType == RoleTypeTenant &&
		!config.Permissions.GroupPermissions.IsNull() &&
		!config.Permissions.GroupPermissions.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.group_permissions"),
			"Conflicting attributes in configuration",
			`group_permissions not available for role_type "`+RoleTypeTenant+`". `+
				`Set role_type to "`+RoleTypeUser+`" to set group_permissions.`,
		)

		return
	}

	// if roleType is "tenant" and default_group_access has been set...
	if roleType == RoleTypeTenant &&
		!config.Permissions.DefaultGroupAccess.IsNull() &&
		!config.Permissions.DefaultGroupAccess.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions.default_group_access"),
			"Conflicting attributes in configuration",
			`default_group_access not available for role_type "`+RoleTypeTenant+`". `+
				`Set role_type to "`+RoleTypeUser+`" to set default_group_access.`,
		)

		return
	}

	// if roleType is "tenant" and multitenant has been set...
	if roleType == RoleTypeTenant &&
		!config.Multitenant.IsNull() &&
		!config.Multitenant.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("multitenant"),
			"Conflicting attributes in configuration",
			`multitenant not available for role_type "`+RoleTypeTenant+`". `+
				`Set role_type to "`+RoleTypeUser+`" to set multitenant.`,
		)

		return
	}

	// if roleType is "tenant" and multitenant_locked has been set...
	if roleType == RoleTypeTenant &&
		!config.MultitenantLocked.IsNull() &&
		!config.MultitenantLocked.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("multitenant_locked"),
			"Conflicting attributes in configuration",
			`multitenant_locked not available for role_type "`+RoleTypeTenant+`". `+
				`Set role_type to "`+RoleTypeUser+`" to set multitenant_locked.`,
		)

		return
	}
}
