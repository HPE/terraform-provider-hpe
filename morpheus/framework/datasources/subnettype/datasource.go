// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnettype

import (
	"context"
	"encoding/json"
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

func subnetTypeAsState(
	t *sdk.GetSubnetType200ResponseSubnetType,
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
) (*sdk.GetSubnetType200ResponseSubnetType, error) {
	r, hresp, err := apiClient.NetworksAPI.GetSubnetType(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for subnet type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.SubnetType == nil {
		return nil, fmt.Errorf("GET failed for subnet type %d: response missing subnetType", id)
	}

	st := *r.SubnetType

	return &st, nil
}

func getSubnetTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetSubnetType200ResponseSubnetType, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListSubnetTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for subnet type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.SubnetTypes)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling subnet types: %w", marshalErr)
	}

	var subnetTypes []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &subnetTypes); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding subnet types: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, st := range subnetTypes {
		if st.Name != nil && *st.Name == name {
			if st.Id != nil {
				matchedID = *st.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoSubnetTypeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleSubnetTypes)
	}

	return getSubnetTypeByID(ctx, matchedID, apiClient)
}

func getSubnetType(
	ctx context.Context,
	config *SubnetTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetSubnetType200ResponseSubnetType, error) {
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
