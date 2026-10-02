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

// metricScopeAPIValues maps the uppercase Terraform scope value to the casing
// the metrics list endpoint expects.
var metricScopeAPIValues = map[string]string{
	"GLOBAL":           "global",
	"SERVICE PROVIDER": "service provider",
	"PARTNER":          "partner",
	"CLIENT":           "client",
}

var _ datasource.DataSource = &metricDataSource{}

func NewMetricDataSource() datasource.DataSource {
	return &metricDataSource{}
}

type metricDataSource struct {
	BaseData
}

type metricDataSourceModel struct {
	Client       types.String `tfsdk:"client"`
	MetricName   types.String `tfsdk:"metric_name"`
	Scope        types.String `tfsdk:"scope"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	Unit         types.String `tfsdk:"unit"`
	UnitLabel    types.String `tfsdk:"unit_label"`
	GroupName    types.String `tfsdk:"group_name"`
	Status       types.String `tfsdk:"status"`
	Type         types.String `tfsdk:"type"`
	ResourceType types.String `tfsdk:"resource_type"`
}

func (d *metricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (d *metricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp metric by its exact metric name, for use when authoring `hpe_opsramp_template` metric groups.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"metric_name": schema.StringAttribute{
				Required:    true,
				Description: "The exact metric name to look up (e.g. `system_linux_services_status`).",
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Metric scope to search within. Valid values: `GLOBAL`, `SERVICE PROVIDER`, `PARTNER`, `CLIENT`. Defaults to `CLIENT`.",
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "SERVICE PROVIDER", "PARTNER", "CLIENT"),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name shown in the UI.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A description of the metric.",
			},
			"unit": schema.StringAttribute{
				Computed:    true,
				Description: "The metric unit (e.g. `%`, `ms`).",
			},
			"unit_label": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable unit label.",
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "The metric group name this metric belongs to.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The metric status (e.g. `active`).",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The metric data type (e.g. `Counter`, `Gauge`).",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "The resource type this metric applies to.",
			},
		},
	}
}

func (d *metricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data metricDataSourceModel
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

	metric, err := d.apiClient.FindMetricByName(tenantID, data.MetricName.ValueString(), metricScopeAPIValues[scope])
	if err != nil {
		resp.Diagnostics.AddError(
			"Metric lookup failed",
			fmt.Sprintf("Could not find metric %q: %s", data.MetricName.ValueString(), err.Error()),
		)

		return
	}

	data.DisplayName = types.StringValue(metric.DisplayName)
	data.Description = types.StringValue(metric.Description)
	data.Unit = types.StringValue(metric.Unit)
	data.UnitLabel = types.StringValue(metric.UnitLabel)
	data.GroupName = types.StringValue(metric.GroupName)
	data.Status = types.StringValue(metric.Status)
	data.Type = types.StringValue(metric.Type)
	data.ResourceType = types.StringValue(metric.ResourceType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
