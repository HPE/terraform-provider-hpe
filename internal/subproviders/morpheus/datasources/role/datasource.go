// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package role

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/role/consts"
	providererrors "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

const summary = "read role data source"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_role"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = RoleDataSourceSchema(ctx)
}

// This function breaks out the logic of reading permissions from API response to store to state.
func populateRoleAsStatePermissions(ctx context.Context, r *sdk.GetRole200Response) (PermissionsValue, diag.Diagnostics) {

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

func roleAsState(
	ctx context.Context,
	role *sdk.GetRole200Response,
) (RoleModel, diag.Diagnostics) {
	var state RoleModel
	var diags diag.Diagnostics

	permissions, diags := populateRoleAsStatePermissions(ctx, role)
	if diags.HasError() {

		return state, diags
	}

	state.Id = convert.Int64ToType(role.Role.Id)
	state.Name = convert.StrToType(role.Role.Name)
	state.Description = convert.StrToType(role.Role.Description.Get())
	state.LandingUrl = convert.StrToType(role.Role.LandingUrl.Get())
	state.Multitenant = convert.BoolToType(role.Role.Multitenant)
	state.MultitenantLocked = convert.BoolToType(role.Role.MultitenantLocked)
	state.RoleType = convert.StrToType(role.Role.RoleType)
	state.Permissions = permissions

	return state, diags
}

func getRoleByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetRole200Response, error) {
	r, hresp, err := apiClient.RolesAPI.GetRole(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for role %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	return r, nil
}

func getRoleByName(
	ctx context.Context,
	data RoleModel,
	apiClient *sdk.APIClient,
) (*sdk.GetRole200Response, error) {
	name := data.Name.ValueString()

	rs, hresp, err := apiClient.RolesAPI.ListRoles(ctx).Authority(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for role %s", name)
	}

	var roles []sdk.ListRoles200ResponseAllOfRolesInner

	for _, c := range rs.Roles {
		if c.GetName() == name {
			roles = append(roles, c)
		}
	}

	if len(roles) == 1 {
		return getRoleByID(ctx, roles[0].GetId(), apiClient)
	} else if len(roles) > 1 {
		return nil, errors.New(consts.ErrorMultipleRoles)
	}

	return nil, errors.New(consts.ErrorNoRoleFound)
}

func getRole(
	ctx context.Context,
	data RoleModel,
	apiClient *sdk.APIClient,
) (*sdk.GetRole200Response, error) {
	if !data.Id.IsNull() {
		return getRoleByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getRoleByName(ctx, data, apiClient)
	}

	return nil, errors.New(consts.ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data RoleModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	role, err := getRole(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	apiState, diags := roleAsState(ctx, role)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	// Perform additional validation of default group/cloud access based on the role_type.
	// Morpheus API does not perform validation like this, but the Morpheus UI does.

	// Only account roles should be able to set default cloud access
	if apiState.RoleType.ValueString() == consts.RoleTypeUser {
		apiState.Permissions.DefaultCloudAccess = types.StringNull()
	}

	// Only user roles should be able to set default group access
	if apiState.RoleType.ValueString() == consts.RoleTypeAccount {
		apiState.Permissions.DefaultGroupAccess = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
}
