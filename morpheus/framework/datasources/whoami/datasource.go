// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package whoami

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	internalErrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read whoami data source"

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

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
	resp.TypeName = req.ProviderTypeName + "_" + "whoami"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = WhoamiDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data WhoamiModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	apiResp, hresp, apiErr := apiClient.AuthenticationAPI.Whoami(ctx).Execute()
	if apiResp == nil || apiErr != nil {
		resp.Diagnostics.AddError(summary,
			"GET failed for whoami: "+internalErrors.ErrMsg(apiErr, hresp))

		return
	}

	user := apiResp.User
	if user == nil {
		resp.Diagnostics.AddError(summary, "user object missing from whoami response")

		return
	}

	// Roles
	roleValues := []attr.Value{}
	if user.Roles != nil {
		for _, role := range user.Roles {
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
				resp.Diagnostics.Append(roleDiags...)

				return
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
			resp.Diagnostics.Append(roleDiags...)

			return
		}
	}

	// Permissions
	permValues := []attr.Value{}
	if apiResp.Permissions != nil {
		for _, perm := range apiResp.Permissions {
			permValue, permDiags := NewPermissionsValue(
				PermissionsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"access": convert.StrToType(perm.Access),
					"code":   convert.StrToType(perm.Code),
					"name":   convert.StrToType(perm.Name),
				},
			)
			if permDiags.HasError() {
				resp.Diagnostics.Append(permDiags...)

				return
			}
			permValues = append(permValues, permValue)
		}
	}

	var permSet types.Set
	if len(permValues) == 0 {
		permSet = types.SetValueMust(PermissionsType{
			ObjectType: types.ObjectType{
				AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
			},
		}, []attr.Value{})
	} else {
		var permDiags diag.Diagnostics
		permSet, permDiags = types.SetValueFrom(ctx, PermissionsType{
			ObjectType: types.ObjectType{
				AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
			},
		}, permValues)
		if permDiags.HasError() {
			resp.Diagnostics.Append(permDiags...)

			return
		}
	}

	// Tenant
	if tenant := user.Account; tenant != nil {
		tenantValue, tenantDiags := NewTenantValue(
			TenantValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(tenant.Id),
				"name": convert.StrToType(tenant.Name),
			},
		)
		if tenantDiags.HasError() {
			resp.Diagnostics.Append(tenantDiags...)

			return
		}
		data.Tenant = tenantValue
	} else {
		data.Tenant = NewTenantValueNull()
	}

	// Default Persona
	if defaultPersona := user.DefaultPersona; defaultPersona != nil {
		dpValue, dpDiags := NewDefaultPersonaValue(
			DefaultPersonaValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(defaultPersona.Id),
				"code": convert.StrToType(defaultPersona.Code),
				"name": convert.StrToType(defaultPersona.Name),
			},
		)
		if dpDiags.HasError() {
			resp.Diagnostics.Append(dpDiags...)

			return
		}
		data.DefaultPersona = dpValue
	} else {
		data.DefaultPersona = NewDefaultPersonaValueNull()
	}

	// Scalar fields from user
	data.Id = convert.Int64ToType(user.Id)
	data.TenantId = convert.Int64ToType(user.AccountId)
	data.Username = convert.StrToType(user.Username)
	data.DisplayName = convert.StrToType(user.DisplayName)
	data.Email = convert.StrToType(user.Email)
	data.FirstName = convert.StrToType(user.FirstName)
	data.LastName = convert.StrToType(user.LastName)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.ReceiveNotifications = convert.BoolToType(user.ReceiveNotifications)
	data.IsUsing2fa = convert.BoolToType(user.IsUsing2FA)
	data.AccountExpired = convert.BoolToType(user.AccountExpired)
	data.AccountLocked = convert.BoolToType(user.AccountLocked)
	data.PasswordExpired = convert.BoolToType(user.PasswordExpired)
	data.LoginCount = convert.Int64ToType(user.LoginCount)
	data.LoginAttempts = convert.Int64ToType(user.LoginAttempts)
	data.LinuxUsername = convert.StrToType(user.LinuxUsername.Get())
	data.WindowsUsername = convert.StrToType(user.WindowsUsername.Get())
	data.LinuxKeyPairId = convert.Int64ToType(user.LinuxKeyPairId.Get())

	// Time fields
	data.LastLoginDate = convert.TimeToType(user.LastLoginDate)
	data.DateCreated = convert.TimeToType(user.DateCreated)
	data.LastUpdated = convert.TimeToType(user.LastUpdated)

	// Root-level fields
	data.IsMasterAccount = convert.BoolToType(apiResp.IsMasterAccount)
	if apiResp.Appliance != nil {
		data.ApplianceBuildVersion = convert.StrToType(apiResp.Appliance.BuildVersion)
	} else {
		data.ApplianceBuildVersion = types.StringNull()
	}

	// Sets
	data.Roles = roleSet
	data.Permissions = permSet

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
