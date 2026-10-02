// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// templateScopeAPIValues maps the uppercase Terraform scope value to the casing
// the templates search endpoint expects.
var templateScopeAPIValues = map[string]string{
	"GLOBAL":           "Global",
	"SERVICE PROVIDER": "Service Provider",
	"PARTNER":          "Partner",
	"CLIENT":           "Client",
}

var _ datasource.DataSource = &templateDataSource{}

func NewTemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	BaseData
}

type templateDataSourceModel struct {
	Client        types.String `tfsdk:"client"`
	Name          types.String `tfsdk:"name"`
	Scope         types.String `tfsdk:"scope"`
	Id            types.String `tfsdk:"id"`
	Description   types.String `tfsdk:"description"`
	CollectorType types.String `tfsdk:"collector_type"`
	Status        types.String `tfsdk:"status"`
	Generation    types.Int64  `tfsdk:"generation"` // 2
	ResourceType  types.String `tfsdk:"resource_type"`
	Tags          types.String `tfsdk:"tags"`
}

func (d *templateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *templateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp monitoring template by exact name and returns its identifiers and summary attributes.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The template name to look up.",
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Template scope to search within. Valid values: `GLOBAL`, `SERVICE PROVIDER`, `PARTNER`, `CLIENT`. Defaults to `CLIENT`.",
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "SERVICE PROVIDER", "PARTNER", "CLIENT"),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the template, usable with `hpe_opsramp_template` import.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A description of the template.",
			},
			"collector_type": schema.StringAttribute{
				Computed:    true,
				Description: "The collector type (e.g. `OpsRamp Agent`, `OpsRamp Gateway`).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The template activity status (e.g. `ACTIVE`).",
			},
			"generation": schema.Int64Attribute{
				Computed:    true,
				Description: "The template generation (`1` or `2`).",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "The resource type this template applies to.",
			},
			"tags": schema.StringAttribute{
				Computed:    true,
				Description: "Tag(s) associated with the template.",
			},
		},
	}
}

func (d *templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data templateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)

	scope := data.Scope.ValueString()
	if scope == "" {
		scope = "CLIENT"
	}
	data.Scope = types.StringValue(scope)

	template, err := d.apiClient.FindTemplateByName(tenantID, data.Name.ValueString(), templateScopeAPIValues[scope])
	if err != nil {
		resp.Diagnostics.AddError(
			"Template lookup failed",
			fmt.Sprintf("Could not find template %q: %s", data.Name.ValueString(), err.Error()),
		)

		return
	}

	data.Id = types.StringValue(template.UniqueId)
	data.Description = types.StringValue(template.Description)
	data.CollectorType = types.StringValue(template.CollectorType)
	data.Status = types.StringValue(template.Status)
	data.Generation = types.Int64Value(int64(template.Generation))
	data.ResourceType = types.StringValue(template.ResourceType)
	data.Tags = types.StringValue(template.Tags)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
