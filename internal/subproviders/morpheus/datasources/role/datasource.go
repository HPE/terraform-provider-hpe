// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package role

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/role/consts"
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

func getRoleByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.ListRoles200ResponseAllOfRolesInner, error) {
	r, hresp, err := apiClient.RolesAPI.GetRole(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for role %d", id)
	}

	role := r.GetRole()

	return &role, nil
}

func getRoleByName(
	ctx context.Context,
	data RoleModel,
	apiClient *sdk.APIClient,
) (*sdk.ListRoles200ResponseAllOfRolesInner, error) {
	name := data.Name.ValueString()

	req := apiClient.RolesAPI.ListRoles(ctx).Authority(name)

	rs, hresp, err := req.Execute()
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
		return &roles[0], nil
	} else if len(roles) > 1 {
		return nil, errors.New(consts.ErrorMultipleRoles)
	}

	return nil, errors.New(consts.ErrorNoRoleFound)
}

func getRole(
	ctx context.Context,
	data RoleModel,
	apiClient *sdk.APIClient,
) (*sdk.ListRoles200ResponseAllOfRolesInner, error) {
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

	data.Id = convert.Int64ToType(role.Id)
	data.Name = convert.StrToType(role.Name)
	data.Description = convert.StrToType(role.Description.Get())
	data.LandingUrl = convert.StrToType(role.LandingUrl.Get())
	data.Multitenant = convert.BoolToType(role.Multitenant)
	data.MultitenantLocked = convert.BoolToType(role.MultitenantLocked)
	data.RoleType = convert.StrToType(role.RoleType)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
