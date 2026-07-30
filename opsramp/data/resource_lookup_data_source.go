// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface satisfaction
var _ datasource.DataSource = &resourceLookupDataSource{}

func NewResourceLookupDataSource() datasource.DataSource {
	return &resourceLookupDataSource{}
}

type resourceLookupDataSource struct {
	BaseData
}

// DS config model
type resourceLookupModel struct {
	Query  types.String `tfsdk:"query"`
	Exists types.Bool   `tfsdk:"exists"`
	ID     types.String `tfsdk:"id"`
}

func (d *resourceLookupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_lookup"
}

func (d *resourceLookupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Checks if a resource exists given a query string. Returns `exists` and the first match `id`.",
		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Required:    true,
				Description: "Opaque query string appended to the API (e.g., `name=my-name&type=app`).",
			},
			"exists": schema.BoolAttribute{
				Computed:    true,
				Description: "True if at least one resource matches.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the first matching resource (if any).",
			},
		},
	}
}

func (d *resourceLookupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data resourceLookupModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exists, id, _, err := d.apiClient.QueryResources(data.Query.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Query failed", err.Error())

		return
	}

	data.Exists = types.BoolValue(exists)
	if exists {
		data.ID = types.StringValue(id)
	} else {
		data.ID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
