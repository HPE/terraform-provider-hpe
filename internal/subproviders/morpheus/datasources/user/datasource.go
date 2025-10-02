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
	if access, accessOk := user.GetAccessOk(); accessOk {
		buildAppTemplates := func() types.Set {
			appTemplatesValues := []attr.Value{}
			for _, in := range access.GetAppTemplates() {
				appTemplatesValues = append(appTemplatesValues,
					NewBlueprintsValueMust(BlueprintsValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(appTemplatesValues) == 0 {
				return types.SetNull(BlueprintsValue{}.Type(ctx))
			}
			set, _ := types.SetValue(BlueprintsValue{}.Type(ctx), appTemplatesValues)

			return set
		}
		appTemplatesSet := buildAppTemplates()

		buildCatalogItemTypes := func() types.Set {
			catalogItemTypesValues := []attr.Value{}
			for _, in := range access.GetCatalogItemTypes() {
				catalogItemTypesValues = append(catalogItemTypesValues,
					NewCatalogItemTypesValueMust(CatalogItemTypesValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(catalogItemTypesValues) == 0 {
				return types.SetNull(CatalogItemTypesValue{}.Type(ctx))
			}
			set, _ := types.SetValue(CatalogItemTypesValue{}.Type(ctx), catalogItemTypesValues)

			return set
		}
		catalogItemTypesSet := buildCatalogItemTypes()

		buildFeatures := func() types.Set {
			featureValues := []attr.Value{}
			for _, in := range access.GetFeatures() {
				featureValues = append(featureValues,
					NewFeaturesValueMust(FeaturesValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access":       types.StringValue(in.GetAccess()),
							"code":         types.StringValue(in.GetCode()),
							"name":         types.StringValue(in.GetName()),
							"sub_category": types.StringValue(in.GetSubCategory()),
						}),
				)
			}
			if len(featureValues) == 0 {
				return types.SetNull(FeaturesValue{}.Type(ctx))
			}
			set, _ := types.SetValue(FeaturesValue{}.Type(ctx), featureValues)

			return set
		}
		featuresSet := buildFeatures()

		buildInstanceTypes := func() types.Set {
			instanceTypesValues := []attr.Value{}
			for _, in := range access.GetInstanceTypes() {
				instanceTypesValues = append(instanceTypesValues,
					NewInstanceTypesValueMust(InstanceTypesValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"code":   types.StringValue(in.GetCode()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(instanceTypesValues) == 0 {
				return types.SetNull(InstanceTypesValue{}.Type(ctx))
			}
			set, _ := types.SetValue(InstanceTypesValue{}.Type(ctx), instanceTypesValues)

			return set
		}
		instanceTypesSet := buildInstanceTypes()

		buildPersonas := func() types.Set {
			personasValues := []attr.Value{}
			for _, in := range access.GetPersonas() {
				personasValues = append(personasValues,
					NewPersonasValueMust(PersonasValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"code":   types.StringValue(in.GetCode()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(personasValues) == 0 {
				return types.SetNull(PersonasValue{}.Type(ctx))
			}
			set, _ := types.SetValue(PersonasValue{}.Type(ctx), personasValues)

			return set
		}
		personasSet := buildPersonas()

		buildReportTypes := func() types.Set {
			reportTypeValues := []attr.Value{}
			for _, in := range access.GetReportTypes() {
				reportTypeValues = append(reportTypeValues,
					NewReportTypesValueMust(ReportTypesValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"code":   types.StringValue(in.GetCode()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(reportTypeValues) == 0 {
				return types.SetNull(ReportTypesValue{}.Type(ctx))
			}
			set, _ := types.SetValue(ReportTypesValue{}.Type(ctx), reportTypeValues)

			return set
		}
		reportTypesSet := buildReportTypes()

		buildSites := func() types.Set {
			sitesValues := []attr.Value{}
			for _, in := range access.GetSites() {
				sitesValues = append(sitesValues, NewGroupsValueMust(
					GroupsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"access": types.StringValue(in.GetAccess()),
						"id":     types.Int64Value(in.GetId()),
						"name":   types.StringValue(in.GetName()),
					}),
				)
			}
			if len(sitesValues) == 0 {
				return types.SetNull(GroupsValue{}.Type(ctx))
			}
			set, _ := types.SetValue(GroupsValue{}.Type(ctx), sitesValues)

			return set
		}
		sitesSet := buildSites()

		buildTaskSets := func() types.Set {
			taskSetsValue := []attr.Value{}
			for _, in := range access.GetTaskSets() {
				taskSetsValue = append(taskSetsValue,
					NewWorkflowsValueMust(WorkflowsValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"code":   convert.StrToType(in.Code.Get()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(taskSetsValue) == 0 {
				return types.SetNull(WorkflowsValue{}.Type(ctx))
			}
			set, _ := types.SetValue(WorkflowsValue{}.Type(ctx), taskSetsValue)

			return set
		}
		taskSetsSet := buildTaskSets()

		buildTasks := func() types.Set {
			tasksValues := []attr.Value{}
			for _, in := range access.GetTasks() {
				tasksValues = append(tasksValues, NewTasksValueMust(
					TasksValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"access": types.StringValue(in.GetAccess()),
						"code":   convert.StrToType(in.Code.Get()),
						"id":     types.Int64Value(in.GetId()),
						"name":   types.StringValue(in.GetName()),
					}),
				)
			}
			if len(tasksValues) == 0 {
				return types.SetNull(TasksValue{}.Type(ctx))
			}
			set, _ := types.SetValue(TasksValue{}.Type(ctx), tasksValues)

			return set
		}
		tasksSet := buildTasks()

		buildVdiPools := func() types.Set {
			vdiPoolsValues := []attr.Value{}
			for _, in := range access.GetVdiPools() {
				vdiPoolsValues = append(vdiPoolsValues,
					NewVdiPoolsValueMust(VdiPoolsValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(vdiPoolsValues) == 0 {
				return types.SetNull(VdiPoolsValue{}.Type(ctx))
			}
			set, _ := types.SetValue(VdiPoolsValue{}.Type(ctx), vdiPoolsValues)

			return set
		}
		vdiPoolsSet := buildVdiPools()

		buildZones := func() types.Set {
			zonesValues := []attr.Value{}
			for _, in := range access.GetZones() {
				zonesValues = append(zonesValues,
					NewCloudsValueMust(CloudsValue{}.AttributeTypes(ctx),
						map[string]attr.Value{
							"access": types.StringValue(in.GetAccess()),
							"id":     types.Int64Value(in.GetId()),
							"name":   types.StringValue(in.GetName()),
						}),
				)
			}
			if len(zonesValues) == 0 {
				return types.SetNull(CloudsValue{}.Type(ctx))
			}
			set, _ := types.SetValue(CloudsValue{}.Type(ctx), zonesValues)

			return set
		}
		zonesSet := buildZones()

		attributes := map[string]attr.Value{
			"app_templates":      appTemplatesSet,
			"catalog_item_types": catalogItemTypesSet,
			"features":           featuresSet,
			"instance_types":     instanceTypesSet,
			"personas":           personasSet,
			"report_types":       reportTypesSet,
			"sites":              sitesSet,
			"task_sets":          taskSetsSet,
			"tasks":              tasksSet,
			"vdi_pools":          vdiPoolsSet,
			"zones":              zonesSet,
		}

		accessValue, accessDiags := NewAccessValue(AccessValue{}.AttributeTypes(ctx), attributes)
		if accessDiags.HasError() {
			diags.Append(accessDiags...)

			return diags
		}
		data.Access = accessValue
	} else {
		data.Access = NewAccessValueNull()
	}

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
	data.AccountId = convert.Int64ToType(user.AccountId)
	data.AccountLocked = convert.BoolToType(user.AccountLocked)
	data.DateCreated = convert.TimeToStrType(user.DateCreated)
	data.DisplayName = convert.StrToType(user.DisplayName)
	data.Email = convert.StrToType(user.Email)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.FirstName = convert.StrToType(user.FirstName)
	data.Id = convert.Int64ToType(user.Id)
	data.IsUsing2fa = convert.BoolToType(user.IsUsing2FA)
	data.LastLoginDate = convert.TimeToStrType(user.LastLoginDate)
	data.LastName = convert.StrToType(user.LastName)
	data.LastUpdated = convert.TimeToStrType(user.LastUpdated)
	data.LinuxKeyPairId = convert.Int64ToType(user.LinuxKeyPairId.Get())
	data.LinuxPassword = convert.StrToType(user.LinuxPassword.Get())
	data.LinuxUsername = convert.StrToType(user.LinuxUsername.Get())
	data.LoginAttempts = convert.Int64ToType(user.LoginAttempts)
	data.LoginCount = convert.Int64ToType(user.LoginCount)
	data.PasswordExpired = convert.BoolToType(user.PasswordExpired)
	data.ReceiveNotifications = convert.BoolToType(user.ReceiveNotifications)
	data.Roles = roleSet
	data.Username = convert.StrToType(user.Username)
	data.WindowsPassword = convert.StrToType(user.WindowsPassword.Get())
	data.WindowsUsername = convert.StrToType(user.WindowsUsername.Get())

	return diags
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
