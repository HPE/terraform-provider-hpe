// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &managementProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &managementProfileDataSource{}
)

func NewManagementProfileDataSource() datasource.DataSource {
	return &managementProfileDataSource{}
}

type managementProfileDataSource struct {
	BaseData
}

type managementProfileDataSourceModel struct {
	Client      types.String `tfsdk:"client"`
	Name        types.String `tfsdk:"name"`
	Uuid        types.String `tfsdk:"uuid"`
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

func (d *managementProfileDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_management_profile"
}

func (d *managementProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp management profile by name and returns its identifiers and attributes.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The management profile name to look up.",
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique UUID of the management profile.",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The numeric ID of the management profile.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Summary describing the management profile.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The management profile type (for example, Gateway).",
			},
		},
	}
}

func (d *managementProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data managementProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)
	profile, err := d.apiClient.FindManagementProfileByName(tenantID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Management profile lookup failed", fmt.Sprintf("Could not find management profile '%s': %s", data.Name.ValueString(), err.Error()))

		return
	}

	data.Name = types.StringValue(profile.Name)
	data.Uuid = types.StringValue(profile.Uuid)
	data.ID = types.Int64Value(int64(profile.Id))
	data.Description = types.StringValue(profile.Description)
	data.Type = types.StringValue(profile.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
