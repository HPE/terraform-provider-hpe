// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"context"

	"github.com/HPE/terraform-provider-opsramp/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface satisfaction
var _ datasource.DataSource = &dataRoleSource{}
var _ datasource.DataSourceWithConfigure = &dataRoleSource{}

func NewDataRoleSource() datasource.DataSource {
	return &dataRoleSource{}
}

type dataRoleSource struct {
	apiClient *client.OpsRampClient
}

func (d *dataRoleSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// DS config model
type dataRoleModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

func (d *dataRoleSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp Role by name. Returns its numeric ID and unique UUID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The numeric ID of the role.",
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the role.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the role to look up.",
			},
		},
	}
}

func (d *dataRoleSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpsRampClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.OpsRampClient",
		)
		return
	}

	d.apiClient = client
}

func (d *dataRoleSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")
		return
	}

	var data dataRoleModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := d.apiClient.TenantId

	role, err := d.apiClient.FindRoleByName(tenantId, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Role retrieve failed", err.Error())
		return
	}

	data.ID = types.Int64Value(int64(role.Id))
	data.Name = types.StringValue(role.Name)
	data.UUID = types.StringValue(role.UniqueId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
