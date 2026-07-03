// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolservertype

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
	summary                            = "read network pool server type data source"
	ErrorNoValidSearchTerms            = `no valid search terms - an id or name is required`
	ErrorRunningPreApply               = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkPoolServerTypeFound  = `no network pool server type found`
	ErrorMultipleNetworkPoolServerType = `multiple network pool server types were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_pool_server_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkPoolServerTypeDataSourceSchema(ctx)
}

func networkPoolServerTypeAsState(
	t *sdk.GetNetworkPoolServerType200ResponseNetworkPoolServerType,
) NetworkPoolServerTypeModel {
	return NetworkPoolServerTypeModel{
		Id:              convert.Int64ToType(t.Id),
		Name:            convert.StrToType(t.Name),
		Code:            convert.StrToType(t.Code),
		Description:     convert.StrToType(t.Description.Get()),
		IntegrationCode: convert.StrToType(t.IntegrationCode.Get()),
		PoolService:     convert.StrToType(t.PoolService.Get()),
		Enabled:         convert.BoolToType(t.Enabled),
		Selectable:      convert.BoolToType(t.Selectable),
	}
}

func getNetworkPoolServerTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServerType200ResponseNetworkPoolServerType, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkPoolServerType(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network pool server type %d: %s", id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.NetworkPoolServerType == nil {
		return nil, fmt.Errorf(
			"GET failed for network pool server type %d: response missing networkPoolServerType", id,
		)
	}

	t := *r.NetworkPoolServerType

	return &t, nil
}

func getNetworkPoolServerTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServerType200ResponseNetworkPoolServerType, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkPoolServerTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network pool server type %s: %s", name, providererrors.ErrMsg(err, hresp),
		)
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.NetworkPoolServerTypes)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling network pool server types: %w", marshalErr)
	}

	var networkPoolServerTypes []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &networkPoolServerTypes); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding network pool server types: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, t := range networkPoolServerTypes {
		if t.Name != nil && *t.Name == name {
			if t.Id != nil {
				matchedID = *t.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkPoolServerTypeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkPoolServerType)
	}

	return getNetworkPoolServerTypeByID(ctx, matchedID, apiClient)
}

func getNetworkPoolServerType(
	ctx context.Context,
	config *NetworkPoolServerTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServerType200ResponseNetworkPoolServerType, error) {
	if !config.Id.IsNull() {
		return getNetworkPoolServerTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkPoolServerTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkPoolServerTypeModel

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

	t, err := getNetworkPoolServerType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkPoolServerTypeAsState(t)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
