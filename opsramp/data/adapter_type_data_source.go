// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &adapterTypeDataSource{}

func NewAdapterTypeDataSource() datasource.DataSource {
	return &adapterTypeDataSource{}
}

type adapterTypeDataSource struct {
	BaseData
}

type adapterTypeDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Id          types.Int64  `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	GroupName   types.String `tfsdk:"group_name"`
	Category    types.String `tfsdk:"category"`
}

func (d *adapterTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adapter_type"
}

func (d *adapterTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp monitoring adapter type by exact name. Use the `name` value as the " +
			"`application_type` on a `hpe_opsramp_template` metric group whose `adapter_type` is `APPLICATION`.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The exact adapter type name to look up (e.g. `RSE`, `HAProxy`).",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The numeric ID of the adapter type.",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name shown in the UI.",
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "The adapter type's group name, if any.",
			},
			"category": schema.StringAttribute{
				Computed:    true,
				Description: "The adapter type's category, if any.",
			},
		},
	}
}

func (d *adapterTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data adapterTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapterType, err := d.apiClient.FindAdapterTypeByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Adapter type lookup failed",
			fmt.Sprintf("Could not find adapter type %q: %s", data.Name.ValueString(), err.Error()),
		)

		return
	}

	data.Id = types.Int64Value(int64(adapterType.Id))
	data.DisplayName = types.StringValue(adapterType.DisplayName)
	data.GroupName = types.StringValue(adapterType.GroupName)
	data.Category = types.StringValue(adapterType.Category)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
