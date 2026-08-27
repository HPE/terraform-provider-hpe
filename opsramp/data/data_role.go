// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface satisfaction
var (
	_ datasource.DataSource              = &dataRoleSource{}
	_ datasource.DataSourceWithConfigure = &dataRoleSource{}
)

func NewDataRoleSource() datasource.DataSource {
	return &dataRoleSource{}
}

type dataRoleSource struct {
	BaseData
}

func (d *dataRoleSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// DS config model
type dataRoleModel struct {
	Client types.String `tfsdk:"client"`
	ID     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	UUID   types.String `tfsdk:"uuid"`
}

func (d *dataRoleSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp Role by name. Returns its numeric ID and unique UUID.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
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

	tenantId := resolveLookupTenantID(d.apiClient, data.Client)

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
