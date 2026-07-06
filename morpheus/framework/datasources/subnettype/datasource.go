// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnettype

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                  = "read subnet type data source"
	ErrorNoValidSearchTerms  = `no valid search terms - an id or name is required`
	ErrorRunningPreApply     = `Error running pre-apply plan: exit status 1`
	ErrorNoSubnetTypeFound   = `no subnet type found`
	ErrorMultipleSubnetTypes = `multiple subnet types were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

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
	resp.TypeName = req.ProviderTypeName + "_" + "subnet_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SubnetTypeDataSourceSchema(ctx)
}

// subnetTypeAsState maps a subnet type list item onto the Terraform state model.
// Subnet types are read from the list endpoint (/api/subnet-types); the list
// returns full objects, so no per-id fetch is required.
func subnetTypeAsState(
	t *sdk.ListSubnetTypes200ResponseAllOfSubnetTypesInner,
) SubnetTypeModel {
	return SubnetTypeModel{
		Id:                 convert.Int64ToType(t.Id),
		Name:               convert.StrToType(t.Name),
		Code:               convert.StrToType(t.Code),
		Description:        convert.StrToType(t.Description),
		Creatable:          convert.BoolToType(t.Creatable),
		Deletable:          convert.BoolToType(t.Deletable),
		DhcpServerEditable: convert.BoolToType(t.DhcpServerEditable),
		CanAssignPool:      convert.BoolToType(t.CanAssignPool),
	}
}

func getSubnetTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.ListSubnetTypes200ResponseAllOfSubnetTypesInner, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListSubnetTypes(ctx).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for subnet type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	for i := range rs.SubnetTypes {
		if rs.SubnetTypes[i].Id != nil && *rs.SubnetTypes[i].Id == id {
			return &rs.SubnetTypes[i], nil
		}
	}

	return nil, fmt.Errorf("%s with id %d", ErrorNoSubnetTypeFound, id)
}

func getSubnetTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.ListSubnetTypes200ResponseAllOfSubnetTypesInner, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListSubnetTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for subnet type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// The Name filter is a partial match; enforce an exact, unique match.
	var match *sdk.ListSubnetTypes200ResponseAllOfSubnetTypesInner

	var matchCount int

	for i := range rs.SubnetTypes {
		if rs.SubnetTypes[i].Name != nil && *rs.SubnetTypes[i].Name == name {
			match = &rs.SubnetTypes[i]
			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoSubnetTypeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleSubnetTypes)
	}

	return match, nil
}

func getSubnetType(
	ctx context.Context,
	config *SubnetTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.ListSubnetTypes200ResponseAllOfSubnetTypesInner, error) {
	if !config.Id.IsNull() {
		return getSubnetTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getSubnetTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SubnetTypeModel

	// Read config
	diags := req.Config.Get(ctx, &config)
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

	st, err := getSubnetType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := subnetTypeAsState(st)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
