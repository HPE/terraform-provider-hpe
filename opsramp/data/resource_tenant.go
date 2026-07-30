// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

// Ensure interface satisfaction
var _ datasource.DataSource = &resourceTenantDataSource{}

func NewResourceTenantDataSource() datasource.DataSource {
	return &resourceTenantDataSource{}
}

type resourceTenantDataSource struct {
	BaseData
}

// DS config model
type resourceTenantModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

func (d *resourceTenantDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (d *resourceTenantDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about the current authenticated tenant.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the tenant.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the tenant.",
			},
		},
	}
}

func (d *resourceTenantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data resourceTenantModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientInfo *client.ClientResponse
	var err error

	if d.apiClient.Scope == "MSP" {
		clientInfo, err = d.apiClient.GetTenantInfo(d.apiClient.TenantId)
	} else {
		clientInfo, err = d.apiClient.GetClientInfo(d.apiClient.TenantId)
	}

	if err != nil {
		resp.Diagnostics.AddError("Tenant information failed to retrieve", err.Error())

		return
	}

	data.Name = types.StringValue(clientInfo.Name)
	data.UUID = types.StringValue(clientInfo.UniqueId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
