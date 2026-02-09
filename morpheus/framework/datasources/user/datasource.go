// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/user/consts"
	internalErrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read user data source"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_user"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = UserDataSourceSchema(ctx)
}

func getUserByUsername(
	ctx context.Context,
	data *UserModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	name := data.Username.ValueString()
	us, hresp, err := apiClient.UsersAPI.ListUsers(ctx).Username(name).Execute()
	if us == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for user with username %s: %s",
			name, internalErrors.ErrMsg(err, hresp)))

		return diags
	}

	users := us.GetUsers()

	// Additional filtering to ensure exact username match (API might return partial matches)
	var filteredUsers []sdk.ListUsers200ResponseAllOfUsersInner
	for _, u := range users {
		if u.GetUsername() == data.Username.ValueString() {
			filteredUsers = append(filteredUsers, u)
		}
	}
	users = filteredUsers

	if len(users) > 1 {
		diags.AddError(summary, consts.ErrorMultipleUsers)

		return diags
	} else if len(users) == 0 {
		diags.AddError(summary, consts.ErrorNoUserFound)

		return diags
	}

	user := users[0]

	return getUserByID(ctx, user.GetId(), data, apiClient)
}

func getUserByID(
	ctx context.Context,
	id int64,
	data *UserModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Create user request with access included
	getUserRequest := apiClient.UsersAPI.GetUser(ctx, id)
	getUserRequest = getUserRequest.IncludeAccess(true)

	u, hresp, err := getUserRequest.Execute()
	if u == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for user with id %d: %s",
			id, internalErrors.ErrMsg(err, hresp)))

		return diags
	}
	user, ok := u.GetUserOk()
	if !ok {
		diags.AddError(summary, consts.ErrorNoUserFound)

		return diags
	}

	roleValues := []attr.Value{}
	if roles, ok := user.GetRolesOk(); ok {
		for _, role := range roles {
			roleValue, roleDiags := NewRolesValue(
				RolesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"authority":   convert.StrToType(role.Authority),
					"description": convert.StrToType(role.Description.Get()),
					"id":          convert.Int64ToType(role.Id),
					"name":        convert.StrToType(role.Name),
				},
			)
			if roleDiags.HasError() {
				diags.Append(roleDiags...)

				return diags
			}
			roleValues = append(roleValues, roleValue)
		}
	}

	var roleSet types.Set
	if len(roleValues) == 0 {
		roleSet = types.SetValueMust(RolesType{
			ObjectType: types.ObjectType{
				AttrTypes: RolesValue{}.AttributeTypes(ctx),
			},
		}, []attr.Value{})
	} else {
		var roleDiags diag.Diagnostics
		roleSet, roleDiags = types.SetValueFrom(ctx, RolesType{
			ObjectType: types.ObjectType{
				AttrTypes: RolesValue{}.AttributeTypes(ctx),
			},
		}, roleValues)
		if roleDiags.HasError() {
			diags.Append(roleDiags...)

			return diags
		}
	}

	// Handle Access complex object
	access, accessDiags := getAccessAsState(ctx, user)
	if accessDiags.HasError() {
		diags.Append(accessDiags...)

		return diags
	}
	data.Access = access

	if tenant, tenantOk := user.GetAccountOk(); tenantOk {
		tenantValue, tenantDiags := NewTenantValue(
			TenantValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(tenant.Id),
				"name": convert.StrToType(tenant.Name),
			},
		)
		if tenantDiags.HasError() {
			diags.Append(tenantDiags...)

			return diags
		}
		data.Tenant = tenantValue
	} else {
		data.Tenant = NewTenantValueNull()
	}

	if defaultPersona, defaultPersonaOk := user.GetDefaultPersonaOk(); defaultPersonaOk {
		defaultPersonaValue, defaultPersonaDiags := NewDefaultPersonaValue(
			DefaultPersonaValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(defaultPersona.Id),
				"code": convert.StrToType(defaultPersona.Code),
				"name": convert.StrToType(defaultPersona.Name),
			},
		)
		if defaultPersonaDiags.HasError() {
			diags.Append(defaultPersonaDiags...)

			return diags
		}
		data.DefaultPersona = defaultPersonaValue
	} else {
		data.DefaultPersona = NewDefaultPersonaValueNull()
	}

	data.AccountExpired = convert.BoolToType(user.AccountExpired)
	data.AccountLocked = convert.BoolToType(user.AccountLocked)
	data.DisplayName = convert.StrToType(user.DisplayName)
	data.Email = convert.StrToType(user.Email)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.FirstName = convert.StrToType(user.FirstName)
	data.Id = convert.Int64ToType(user.Id)
	data.IsUsing2fa = convert.BoolToType(user.IsUsing2FA)
	data.LastName = convert.StrToType(user.LastName)
	data.LinuxKeyPairId = convert.Int64ToType(user.LinuxKeyPairId.Get())
	data.LinuxUsername = convert.StrToType(user.LinuxUsername.Get())
	data.PasswordExpired = convert.BoolToType(user.PasswordExpired)
	data.ReceiveNotifications = convert.BoolToType(user.ReceiveNotifications)
	data.Roles = roleSet
	data.Username = convert.StrToType(user.Username)
	data.WindowsUsername = convert.StrToType(user.WindowsUsername.Get())

	return diags
}

func getAccessAsState(
	ctx context.Context,
	user *sdk.GetUser200ResponseUser,
) (AccessValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	access, accessOk := user.GetAccessOk()
	if !accessOk {
		return NewAccessValueNull(), diags
	}

	// Collect element structs first
	var blueprints []BlueprintsValue
	for _, in := range access.GetAppTemplates() {
		blueprints = append(blueprints, BlueprintsValue{
			Access: types.StringValue(in.GetAccess()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var catalogItemTypes []CatalogItemTypesValue
	for _, in := range access.GetCatalogItemTypes() {
		catalogItemTypes = append(catalogItemTypes, CatalogItemTypesValue{
			Access: types.StringValue(in.GetAccess()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var features []FeaturesValue
	for _, in := range access.GetFeatures() {
		features = append(features, FeaturesValue{
			Access:      types.StringValue(in.GetAccess()),
			Code:        types.StringValue(in.GetCode()),
			Name:        types.StringValue(in.GetName()),
			SubCategory: types.StringValue(in.GetSubCategory()),
			state:       attr.ValueStateKnown,
		})
	}

	var instanceTypes []InstanceTypesValue
	for _, in := range access.GetInstanceTypes() {
		instanceTypes = append(instanceTypes, InstanceTypesValue{
			Access: types.StringValue(in.GetAccess()),
			Code:   types.StringValue(in.GetCode()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var personas []PersonasValue
	for _, in := range access.GetPersonas() {
		personas = append(personas, PersonasValue{
			Access: types.StringValue(in.GetAccess()),
			Code:   types.StringValue(in.GetCode()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var reportTypes []ReportTypesValue
	for _, in := range access.GetReportTypes() {
		reportTypes = append(reportTypes, ReportTypesValue{
			Access: types.StringValue(in.GetAccess()),
			Code:   types.StringValue(in.GetCode()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var groups []GroupsValue
	for _, in := range access.GetSites() {
		groups = append(groups, GroupsValue{
			Access: types.StringValue(in.GetAccess()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var workflows []WorkflowsValue
	for _, in := range access.GetTaskSets() {
		workflows = append(workflows, WorkflowsValue{
			Access: types.StringValue(in.GetAccess()),
			Code:   convert.StrToType(in.Code.Get()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var tasks []TasksValue
	for _, in := range access.GetTasks() {
		tasks = append(tasks, TasksValue{
			Access: types.StringValue(in.GetAccess()),
			Code:   convert.StrToType(in.Code.Get()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var vdiPools []VdiPoolsValue
	for _, in := range access.GetVdiPools() {
		vdiPools = append(vdiPools, VdiPoolsValue{
			Access: types.StringValue(in.GetAccess()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	var clouds []CloudsValue
	for _, in := range access.GetZones() {
		clouds = append(clouds, CloudsValue{
			Access: types.StringValue(in.GetAccess()),
			Id:     types.Int64Value(in.GetId()),
			Name:   types.StringValue(in.GetName()),
			state:  attr.ValueStateKnown,
		})
	}

	// Create sets, ensuring empty slices result in empty sets rather than null
	var blueprintsSet types.Set
	if len(blueprints) == 0 {
		blueprintsSet = types.SetValueMust(BlueprintsValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		blueprintsSet, setDiags = types.SetValueFrom(ctx, BlueprintsValue{}.Type(ctx), blueprints)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var catalogItemTypesSet types.Set
	if len(catalogItemTypes) == 0 {
		catalogItemTypesSet = types.SetValueMust(CatalogItemTypesValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		catalogItemTypesSet, setDiags = types.SetValueFrom(ctx, CatalogItemTypesValue{}.Type(ctx), catalogItemTypes)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var featuresSet types.Set
	if len(features) == 0 {
		featuresSet = types.SetValueMust(FeaturesValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		featuresSet, setDiags = types.SetValueFrom(ctx, FeaturesValue{}.Type(ctx), features)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var instanceTypesSet types.Set
	if len(instanceTypes) == 0 {
		instanceTypesSet = types.SetValueMust(InstanceTypesValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		instanceTypesSet, setDiags = types.SetValueFrom(ctx, InstanceTypesValue{}.Type(ctx), instanceTypes)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var personasSet types.Set
	if len(personas) == 0 {
		personasSet = types.SetValueMust(PersonasValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		personasSet, setDiags = types.SetValueFrom(ctx, PersonasValue{}.Type(ctx), personas)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var reportTypesSet types.Set
	if len(reportTypes) == 0 {
		reportTypesSet = types.SetValueMust(ReportTypesValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		reportTypesSet, setDiags = types.SetValueFrom(ctx, ReportTypesValue{}.Type(ctx), reportTypes)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var groupsSet types.Set
	if len(groups) == 0 {
		groupsSet = types.SetValueMust(GroupsValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		groupsSet, setDiags = types.SetValueFrom(ctx, GroupsValue{}.Type(ctx), groups)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var workflowsSet types.Set
	if len(workflows) == 0 {
		workflowsSet = types.SetValueMust(WorkflowsValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		workflowsSet, setDiags = types.SetValueFrom(ctx, WorkflowsValue{}.Type(ctx), workflows)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var tasksSet types.Set
	if len(tasks) == 0 {
		tasksSet = types.SetValueMust(TasksValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		tasksSet, setDiags = types.SetValueFrom(ctx, TasksValue{}.Type(ctx), tasks)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var vdiPoolsSet types.Set
	if len(vdiPools) == 0 {
		vdiPoolsSet = types.SetValueMust(VdiPoolsValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		vdiPoolsSet, setDiags = types.SetValueFrom(ctx, VdiPoolsValue{}.Type(ctx), vdiPools)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	var cloudsSet types.Set
	if len(clouds) == 0 {
		cloudsSet = types.SetValueMust(CloudsValue{}.Type(ctx), []attr.Value{})
	} else {
		var setDiags diag.Diagnostics
		cloudsSet, setDiags = types.SetValueFrom(ctx, CloudsValue{}.Type(ctx), clouds)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return NewAccessValueNull(), diags
		}
	}

	return NewAccessValue(AccessValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"blueprints":         blueprintsSet,
		"catalog_item_types": catalogItemTypesSet,
		"features":           featuresSet,
		"instance_types":     instanceTypesSet,
		"personas":           personasSet,
		"report_types":       reportTypesSet,
		"groups":             groupsSet,
		"workflows":          workflowsSet,
		"tasks":              tasksSet,
		"vdi_pools":          vdiPoolsSet,
		"clouds":             cloudsSet,
	})
}

func getUser(
	ctx context.Context,
	data *UserModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	if !data.Id.IsNull() {
		return getUserByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Username.IsNull() {
		return getUserByUsername(ctx, data, apiClient)
	}

	diags := diag.Diagnostics{}
	diags.AddError(summary, consts.ErrorNoValidUserTerms)

	return diags
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data UserModel

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

	diags = getUser(ctx, &data, apiClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
