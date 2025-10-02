// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/user/consts"
	internalErrors "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
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

	roleSet, roleDiags := types.SetValueFrom(ctx, RolesType{
		ObjectType: types.ObjectType{
			AttrTypes: RolesValue{}.AttributeTypes(ctx),
		},
	}, roleValues)
	if roleDiags.HasError() {
		diags.Append(roleDiags...)

		return diags
	}

	// Handle Access complex object
	access, accessDiags := getAccessAsState(ctx, user)
	if accessDiags.HasError() {
		diags.Append(accessDiags...)

		return diags
	}
	data.Access = access

	if account, accountOk := user.GetAccountOk(); accountOk {
		accountValue, accountDiags := NewAccountValue(
			AccountValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(account.Id),
				"name": convert.StrToType(account.Name),
			},
		)
		if accountDiags.HasError() {
			diags.Append(accountDiags...)

			return diags
		}
		data.Account = accountValue
	} else {
		data.Account = NewAccountValueNull()
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
	data.TenantId = convert.Int64ToType(user.AccountId)
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
	user *sdk.ListUsers200ResponseAllOfUsersInner,
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

	blueprintsSet, diags := types.SetValueFrom(ctx, BlueprintsValue{}.Type(ctx), blueprints)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	catalogItemTypesSet, diags := types.SetValueFrom(ctx, CatalogItemTypesValue{}.Type(ctx), catalogItemTypes)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	featuresSet, diags := types.SetValueFrom(ctx, FeaturesValue{}.Type(ctx), features)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	instanceTypesSet, diags := types.SetValueFrom(ctx, InstanceTypesValue{}.Type(ctx), instanceTypes)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	personasSet, diags := types.SetValueFrom(ctx, PersonasValue{}.Type(ctx), personas)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	reportTypesSet, diags := types.SetValueFrom(ctx, ReportTypesValue{}.Type(ctx), reportTypes)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	groupsSet, diags := types.SetValueFrom(ctx, GroupsValue{}.Type(ctx), groups)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	workflowsSet, diags := types.SetValueFrom(ctx, WorkflowsValue{}.Type(ctx), workflows)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	tasksSet, diags := types.SetValueFrom(ctx, TasksValue{}.Type(ctx), tasks)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	vdiPoolsSet, diags := types.SetValueFrom(ctx, VdiPoolsValue{}.Type(ctx), vdiPools)
	if diags.HasError() {
		return NewAccessValueNull(), diags
	}
	cloudsSet, diags := types.SetValueFrom(ctx, CloudsValue{}.Type(ctx), clouds)
	if diags.HasError() {
		return NewAccessValueNull(), diags
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
