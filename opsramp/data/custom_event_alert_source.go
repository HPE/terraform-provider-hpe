// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

// Ensure interface satisfaction
var _ datasource.DataSource = &customEventAlertSourceDataSource{}

func NewCustomEventAlertSourceDataSource() datasource.DataSource {
	return &customEventAlertSourceDataSource{}
}

type customEventAlertSourceDataSource struct {
	BaseData
}

type customEventAlertSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	TechUID     types.String `tfsdk:"tech_uid"`
}

func (d *customEventAlertSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_event_alert_source"
}

func (d *customEventAlertSourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an available alert source by name for CUSTOM-EVENT integrations. Returns its numeric ID for use in opsramp_integration resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The numeric ID of the alert source.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the alert source to look up (e.g. 'Alien Vault').",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the alert source.",
			},
			"tech_uid": schema.StringAttribute{
				Computed:    true,
				Description: "The technical UID of the alert source (e.g. ALIENVAULT).",
			},
		},
	}
}

func (d *customEventAlertSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data customEventAlertSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := d.apiClient.TenantId

	sources, err := d.apiClient.GetAvailableAlertSources(tenantId, "CUSTOM-EVENT")
	if err != nil {
		resp.Diagnostics.AddError("Alert Source retrieval failed",
			fmt.Sprintf("Could not retrieve alert sources for CUSTOM-EVENT: %s", err.Error()))

		return
	}

	// Find the alert source by name (match on Name or DisplayName)
	searchName := data.Name.ValueString()
	var found *client.AlertSource
	for i := range sources {
		if sources[i].Name == searchName || sources[i].DisplayName == searchName {
			found = &sources[i]

			break
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Alert Source not found",
			fmt.Sprintf("No alert source named '%s' found for CUSTOM-EVENT", searchName))

		return
	}

	data.ID = types.Int64Value(int64(found.ID))
	data.Name = types.StringValue(found.Name)
	data.DisplayName = types.StringValue(found.DisplayName)
	data.TechUID = types.StringValue(found.TechUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
