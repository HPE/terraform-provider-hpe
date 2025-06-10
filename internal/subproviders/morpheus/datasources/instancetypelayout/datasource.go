// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name.tf.tmpl Name "Example name"
//go:generate go run ../../../../../cmd/render example-name-version.tf.tmpl Name "Example name" Version "1.2.3"

package instancetypelayout

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
	summary                          = "read instance type layout data source"
	ErrorNoInstanceTypeLayoutFound   = `no instance type layout found`
	ErrorNoValidSearchTerms          = `no valid search terms - an id or name is required`
	ErrorRunningPreApply             = `Error running pre-apply plan: exit status 1`
	ErrorMultipleInstanceTypeLayouts = `multiple instance type layouts were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance_type_layout"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceTypeLayoutDataSourceSchema(ctx)
}

func getInstanceTypeLayoutByID(
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

func getInstanceTypeLayoutByName(
	ctx context.Context,
	data InstanceTypeLayoutModel,
	apiClient *sdk.APIClient,
) (*sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner, error) {
	name := data.Name.ValueString()

	// Sort by descending display order (sortOrder)
	// https://docs.morpheusdata.com/en/latest/library/blueprints/layouts.html?highlight=high-to-low
	req := apiClient.LibraryAPI.ListLayouts(ctx).Name(name).Sort("sortOrder").Direction("desc")
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

	// We return the first layout which should have the highest display order (sortOrder)
	if len(layouts) > 0 {
		return &layouts[0], nil
	}

	return nil, errors.New(ErrorNoInstanceTypeLayoutFound)
}

func getInstanceTypeLayout(
	ctx context.Context,
	data InstanceTypeLayoutModel,
	apiClient *sdk.APIClient,
) (*sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner, error) {
	if !data.Id.IsNull() {
		return getInstanceTypeLayoutByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getInstanceTypeLayoutByName(ctx, data, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data InstanceTypeLayoutModel

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

	layout, err := getInstanceTypeLayout(ctx, data, apiClient)
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
	data.SortOrder = convert.Int64ToType(layout.SortOrder)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
