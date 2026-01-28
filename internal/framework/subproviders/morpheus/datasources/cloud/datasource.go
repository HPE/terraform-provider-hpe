// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/cloud/consts"
)

const summary = "read cloud data source"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_cloud"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CloudDataSourceSchema(ctx)
}

func getCloudByID(
	ctx context.Context,
	id int64,
	data *CloudModel,
	apiClient *sdk.APIClient,
) error {
	c, hresp, err := apiClient.CloudsAPI.GetClouds(ctx, id).Execute()
	if c == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for cloud %d", id)
	}

	cloud := c.GetZone()

	data.Id = convert.Int64ToType(cloud.Id)
	data.Name = convert.StrToType(cloud.Name)
	data.Code = convert.StrToType(cloud.Code)
	data.CostingMode = convert.StrToType(cloud.CostingMode.Get())
	data.ExternalId = convert.StrToType(cloud.ExternalId.Get())
	data.GuidanceMode = convert.StrToType(cloud.GuidanceMode.Get())
	data.InventoryLevel = convert.StrToType(cloud.InventoryLevel)
	data.Labels = convert.StrSliceToSet(cloud.Labels)
	data.Location = convert.StrToType(cloud.Location.Get())
	data.TimeZone = convert.StrToType(cloud.Timezone.Get())

	var groupIDs []int64
	for _, g := range cloud.Groups {
		groupIDs = append(groupIDs, (*g.Id))
	}
	data.GroupIds = convert.Int64SliceToSet(groupIDs)

	return nil
}

func getCloudByName(
	ctx context.Context,
	data *CloudModel,
	apiClient *sdk.APIClient,
) error {
	name := data.Name.ValueString()

	req := apiClient.CloudsAPI.ListClouds(ctx).Name(name)

	cs, hresp, err := req.Execute()
	if cs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for cloud %s", name)
	}

	clouds := sdk.NewListClouds200Response().Zones

	for _, c := range cs.Zones {
		tflog.Warn(ctx, fmt.Sprintf("found cloud: %s", c.GetName()))
		if c.GetName() == name {
			clouds = append(clouds, c)
		}
	}

	if len(clouds) > 1 {
		return errors.New(consts.ErrorMultipleClouds)
	} else if len(clouds) == 0 {
		return errors.New(consts.ErrorNoCloudFound)
	}

	cloud := clouds[0]

	data.Id = convert.Int64ToType(cloud.Id)
	data.Name = convert.StrToType(cloud.Name)
	data.Code = convert.StrToType(cloud.Code)
	data.CostingMode = convert.StrToType(cloud.CostingMode.Get())
	data.ExternalId = convert.StrToType(cloud.ExternalId.Get())
	data.GuidanceMode = convert.StrToType(cloud.GuidanceMode.Get())
	data.InventoryLevel = convert.StrToType(cloud.InventoryLevel)
	data.Labels = convert.StrSliceToSet(cloud.Labels)
	data.Location = convert.StrToType(cloud.Location.Get())
	data.TimeZone = convert.StrToType(cloud.Timezone.Get())

	var groupIDs []int64
	for _, g := range cloud.Groups {
		groupIDs = append(groupIDs, (*g.Id))
	}
	data.GroupIds = convert.Int64SliceToSet(groupIDs)

	return nil
}

func getCloud(
	ctx context.Context,
	data *CloudModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getCloudByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Name.IsNull() {
		return getCloudByName(ctx, data, apiClient)
	}

	return errors.New(consts.ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data CloudModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	if err := getCloud(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
