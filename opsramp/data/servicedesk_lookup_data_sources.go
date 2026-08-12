// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

var (
	_ datasource.DataSource              = &serviceDeskUrgencyDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceDeskUrgencyDataSource{}
	_ datasource.DataSource              = &serviceDeskBusinessImpactDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceDeskBusinessImpactDataSource{}
	_ datasource.DataSource              = &serviceDeskCategoryDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceDeskCategoryDataSource{}
)

type serviceDeskLookupModel struct {
	Client types.String `tfsdk:"client"`
	Name   types.String `tfsdk:"name"`
	ID     types.String `tfsdk:"id"`
}

type serviceDeskCategoryLookupModel struct {
	Client     types.String `tfsdk:"client"`
	Name       types.String `tfsdk:"name"`
	TicketType types.String `tfsdk:"ticket_type"`
	ID         types.String `tfsdk:"id"`
}

type serviceDeskUrgencyDataSource struct {
	BaseData
}

type serviceDeskBusinessImpactDataSource struct {
	BaseData
}

type serviceDeskCategoryDataSource struct {
	BaseData
}

func NewServiceDeskUrgencyDataSource() datasource.DataSource {
	return &serviceDeskUrgencyDataSource{}
}

func NewServiceDeskBusinessImpactDataSource() datasource.DataSource {
	return &serviceDeskBusinessImpactDataSource{}
}

func NewServiceDeskCategoryDataSource() datasource.DataSource {
	return &serviceDeskCategoryDataSource{}
}

func resolveLookupTenantID(apiClient *client.OpsRampClient, configuredClient types.String) string {
	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := apiClient.TenantId
	if !configuredClient.IsNull() && configuredClient.ValueString() != "" {
		tenantId = configuredClient.ValueString()
	}

	return tenantId
}

func serviceDeskLookupSchema(description string) schema.Schema {
	return schema.Schema{
		Description: description,
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The exact object name to look up.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The string ID of the matching object.",
			},
		},
	}
}

func (d *serviceDeskUrgencyDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_urgency"
}

func (d *serviceDeskUrgencyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = serviceDeskLookupSchema("Looks up an OpsRamp service desk urgency by name and returns its string ID.")
}

func (d *serviceDeskUrgencyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data serviceDeskLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)
	urgency, err := d.apiClient.FindServiceDeskUrgencyByName(tenantID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Service desk urgency lookup failed", err.Error())

		return
	}

	data.ID = types.StringValue(urgency.Id)
	data.Name = types.StringValue(urgency.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *serviceDeskBusinessImpactDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_business_impact"
}

func (d *serviceDeskBusinessImpactDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = serviceDeskLookupSchema("Looks up an OpsRamp service desk business impact by name and returns its string ID.")
}

func (d *serviceDeskBusinessImpactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data serviceDeskLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)
	businessImpact, err := d.apiClient.FindServiceDeskBusinessImpactByName(tenantID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Service desk business impact lookup failed", err.Error())

		return
	}

	data.ID = types.StringValue(businessImpact.Id)
	data.Name = types.StringValue(businessImpact.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *serviceDeskCategoryDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_category"
}

func (d *serviceDeskCategoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp service desk category by name and returns its string ID. An optional client and ticket_type can be supplied to narrow the lookup.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The exact category name to look up.",
			},
			"ticket_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional ticket type to disambiguate category lookup. Accepts values like `incidents`, `problems`, or `serviceRequests`.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The string ID of the matching category.",
			},
		},
	}
}

func (d *serviceDeskCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data serviceDeskCategoryLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)
	category, err := d.apiClient.FindServiceDeskCategoryByName(tenantID, data.Name.ValueString(), data.TicketType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Service desk category lookup failed", err.Error())

		return
	}

	data.ID = types.StringValue(category.Id)
	data.Name = types.StringValue(category.Name)
	data.TicketType = types.StringValue(category.TicketType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
