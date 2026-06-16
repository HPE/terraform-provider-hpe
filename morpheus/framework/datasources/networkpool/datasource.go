// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                   = "read network pool data source"
	ErrorNoValidSearchTerms   = `no valid search terms - an id or name is required`
	ErrorRunningPreApply      = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkPoolFound   = `no network pool found`
	ErrorMultipleNetworkPools = `multiple network pools were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_pool"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkPoolDataSourceSchema(ctx)
}

func networkPoolAsState(
	pool *sdk.GetNetworkPool200ResponseNetworkPool,
) NetworkPoolModel {
	return NetworkPoolModel{
		Id:            convert.Int64ToType(pool.Id),
		Name:          convert.StrToType(pool.Name),
		Category:      convert.StrToType(pool.Category.Get()),
		Code:          convert.StrToType(pool.Code.Get()),
		DhcpIp:        convert.StrToType(pool.DhcpIp.Get()),
		DhcpServer:    convert.BoolToType(pool.DhcpServer),
		DnsDomain:     convert.StrToType(pool.DnsDomain.Get()),
		ExternalId:    convert.StrToType(pool.ExternalId.Get()),
		FreeCount:     convert.Int64ToType(pool.FreeCount),
		Gateway:       convert.StrToType(pool.Gateway.Get()),
		IpCount:       convert.Int64ToType(pool.IpCount),
		Netmask:       convert.StrToType(pool.Netmask.Get()),
		PoolEnabled:   convert.BoolToType(pool.PoolEnabled),
		SubnetAddress: convert.StrToType(pool.SubnetAddress.Get()),
	}
}

func getNetworkPoolByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPool200ResponseNetworkPool, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network pool %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkPool == nil {
		return nil, fmt.Errorf("GET failed for network pool %d: response missing networkPool", id)
	}

	pool := *r.NetworkPool

	return &pool, nil
}

func getNetworkPoolByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPool200ResponseNetworkPool, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkPools(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network pool %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// NetworkPools is interface{} in the SDK — use JSON round-trip to extract typed fields.
	raw, marshalErr := json.Marshal(rs.NetworkPools)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling network pools: %w", marshalErr)
	}

	var pools []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &pools); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding network pools: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, p := range pools {
		if p.Name != nil && *p.Name == name {
			if p.Id != nil {
				matchedID = *p.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkPoolFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkPools)
	}

	return getNetworkPoolByID(ctx, matchedID, apiClient)
}

func getNetworkPool(
	ctx context.Context,
	config *NetworkPoolModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPool200ResponseNetworkPool, error) {
	if !config.Id.IsNull() {
		return getNetworkPoolByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkPoolByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkPoolModel

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

	pool, err := getNetworkPool(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkPoolAsState(pool)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
