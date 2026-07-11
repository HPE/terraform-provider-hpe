// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface satisfaction
var _ datasource.DataSource = &integrationInboundPropertiesDataSource{}

func NewIntegrationInboundPropertiesDataSource() datasource.DataSource {
	return &integrationInboundPropertiesDataSource{}
}

type integrationInboundPropertiesDataSource struct {
	BaseData
}

type integrationInboundPropertiesModel struct {
	IntegrationID types.String           `tfsdk:"integration_id"`
	Entity        types.String           `tfsdk:"entity"`
	Properties    []inboundPropertyModel `tfsdk:"properties"`
}

type inboundPropertyModel struct {
	Property     types.String `tfsdk:"property"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	PropertyType types.String `tfsdk:"property_type"`
	DisplayGroup types.String `tfsdk:"display_group"`
	Mapable      types.Bool   `tfsdk:"mapable"`
	ValueMapable types.Bool   `tfsdk:"value_mapable"`
	Parsable     types.Bool   `tfsdk:"parsable"`
}

func (d *integrationInboundPropertiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_inbound_properties"
}

func (d *integrationInboundPropertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the available inbound entity properties for an installed integration. These properties represent the valid values for opsramp_attribute in map_attributes blocks.",
		Attributes: map[string]schema.Attribute{
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the installed integration (e.g. INTG-...).",
			},
			"entity": schema.StringAttribute{
				Optional:    true,
				Description: "The entity type to retrieve properties for. Defaults to ALERT.",
			},
			"properties": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of available properties for attribute mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Computed:    true,
							Description: "The property identifier (e.g. alert.alertTime). Use this as the opsramp_attribute value.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The display name of the property.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A description of the property.",
						},
						"property_type": schema.StringAttribute{
							Computed:    true,
							Description: "The data type of the property (e.g. STRING).",
						},
						"display_group": schema.StringAttribute{
							Computed:    true,
							Description: "The display group (e.g. NATIVE_ATTRIBUTE).",
						},
						"mapable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the property can be mapped.",
						},
						"value_mapable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether specific values of this property can be mapped.",
						},
						"parsable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the property supports parsing operators.",
						},
					},
				},
			},
		},
	}
}

func (d *integrationInboundPropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")
		return
	}

	var data integrationInboundPropertiesModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity := "ALERT"
	if !data.Entity.IsNull() && data.Entity.ValueString() != "" {
		entity = data.Entity.ValueString()
	}

	tenantId := d.apiClient.TenantId
	integrationId := data.IntegrationID.ValueString()

	properties, err := d.apiClient.GetInboundEntityProperties(tenantId, integrationId, entity)
	if err != nil {
		resp.Diagnostics.AddError("Properties retrieval failed",
			fmt.Sprintf("Could not retrieve inbound properties for integration '%s': %s", integrationId, err.Error()))
		return
	}

	data.Entity = types.StringValue(entity)
	data.Properties = make([]inboundPropertyModel, len(properties))
	for i, p := range properties {
		data.Properties[i] = inboundPropertyModel{
			Property:     types.StringValue(p.Property),
			Name:         types.StringValue(p.Name),
			Description:  types.StringValue(p.Description),
			PropertyType: types.StringValue(p.PropertyType),
			DisplayGroup: types.StringValue(p.DisplayGroup),
			Mapable:      types.BoolValue(p.Mapable),
			ValueMapable: types.BoolValue(p.ValueMapable),
			Parsable:     types.BoolValue(p.Parsable),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
