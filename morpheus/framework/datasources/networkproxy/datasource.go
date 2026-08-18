// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkproxy

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
	summary                     = "read network proxy data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorRunningPreApply        = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkProxyFound    = `no network proxy found`
	ErrorMultipleNetworkProxies = `multiple network proxies were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_proxy"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkProxyDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves the id of a Morpheus network proxy by name or id, " +
		"for use with network_proxy_id on hpe_morpheus_network."
	resp.Schema.MarkdownDescription = "Retrieves the id of a Morpheus network proxy by name or " +
		"id, for use with `network_proxy_id` on `hpe_morpheus_network`."
}

func networkProxyAsState(
	np *sdk.GetNetworkProxy200ResponseNetworkProxy,
) NetworkProxyModel {
	var accountID *int64
	if account := np.Account.Get(); account != nil {
		accountID = account.Id
	}

	var ownerID *int64
	if owner := np.Owner.Get(); owner != nil {
		ownerID = owner.Id
	}

	return NetworkProxyModel{
		Id:               convert.Int64ToType(np.Id),
		Name:             convert.StrToType(np.Name),
		ProxyHost:        convert.StrToType(np.ProxyHost),
		ProxyPort:        convert.Int64ToType(np.ProxyPort),
		ProxyDomain:      convert.StrToType(np.ProxyDomain),
		ProxyUser:        convert.StrToType(np.ProxyUser.Get()),
		ProxyWorkstation: convert.StrToType(np.ProxyWorkstation.Get()),
		AccountId:        convert.Int64ToType(accountID),
		OwnerId:          convert.Int64ToType(ownerID),
		Visibility:       convert.StrToType(np.Visibility),
	}
}

func getNetworkProxyByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkProxy200ResponseNetworkProxy, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkProxy(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network proxy %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkProxy == nil {
		return nil, fmt.Errorf("GET failed for network proxy %d: response missing networkProxy", id)
	}

	np := *r.NetworkProxy

	return &np, nil
}

func getNetworkProxyByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkProxy200ResponseNetworkProxy, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkProxies(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network proxy %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// The list endpoint types networkProxies as a generic interface{} (the
	// OpenAPI list schema is missing `type: array`), so use a JSON round-trip
	// for safe extraction of the id and name. This mirrors the network_type
	// data source.
	raw, marshalErr := json.Marshal(rs.NetworkProxies)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling network proxies: %w", marshalErr)
	}

	var networkProxies []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &networkProxies); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding network proxies: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, np := range networkProxies {
		if np.Name != nil && *np.Name == name {
			if np.Id != nil {
				matchedID = *np.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkProxyFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkProxies)
	}

	return getNetworkProxyByID(ctx, matchedID, apiClient)
}

func getNetworkProxy(
	ctx context.Context,
	config *NetworkProxyModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkProxy200ResponseNetworkProxy, error) {
	if !config.Id.IsNull() {
		return getNetworkProxyByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkProxyByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkProxyModel

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

	np, err := getNetworkProxy(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkProxyAsState(np)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
