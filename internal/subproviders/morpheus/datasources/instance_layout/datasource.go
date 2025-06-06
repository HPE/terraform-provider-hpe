// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instance_layout

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
)

const (
	summary                      = "read instance layout data source"
	ErrorNoInstanceLayoutFound   = `no instance layout found`
	ErrorNoValidSearchTerms      = `no valid search terms - an id or name is required`
	ErrorRunningPreApply         = `Error running pre-apply plan: exit status 1`
	ErrorMultipleInstanceLayouts = `multiple instance layouts were returned`
)

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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance_layout"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceLayoutDataSourceSchema(ctx)
}

func getInstanceLayoutByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner, error) {
	c, hresp, err := apiClient.LibraryAPI.GetLayout(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for instance layout %d", id)
	}

	layout := c.GetInstanceTypeLayout()

	return &layout, nil
}

func getInstanceLayoutByName(
	ctx context.Context,
	data InstanceLayoutModel,
	apiClient *sdk.APIClient,
) (*sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner, error) {
	name := data.Name.ValueString()

	req := apiClient.LibraryAPI.ListLayouts(ctx).Name(name)
	if !data.Version.IsNull() {
		req = req.Max(5000) // the api doesn't support filtering by version
	}

	ls, hresp, err := req.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for instance layout %s", name)
	}

	var layouts []sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner

	for _, l := range ls.InstanceTypeLayouts {
		if l.GetName() == name {
			layouts = append(layouts, l)
		}
	}

	if !data.Version.IsNull() {
		version := data.Version.ValueString()

		var filtered []sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner
		for _, l := range layouts {
			if l.GetInstanceVersion() == version {
				filtered = append(filtered, l)
			}
		}

		layouts = filtered
	}

	if len(layouts) == 1 {
		return &layouts[0], nil
	} else if len(layouts) > 1 {
		return nil, errors.New(ErrorMultipleInstanceLayouts)
	}

	return nil, errors.New(ErrorNoInstanceLayoutFound)
}

func getInstanceLayout(
	ctx context.Context,
	data InstanceLayoutModel,
	apiClient *sdk.APIClient,
) (*sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner, error) {
	if !data.Id.IsNull() {
		return getInstanceLayoutByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getInstanceLayoutByName(ctx, data, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data InstanceLayoutModel

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

	layout, err := getInstanceLayout(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	data.Id = convert.Int64ToType(layout.Id)
	data.Name = convert.StrToType(layout.Name)
	data.Code = convert.StrToType(layout.Code)
	data.Description = convert.StrToType(layout.Description)
	data.Version = convert.StrToType(layout.InstanceVersion)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
