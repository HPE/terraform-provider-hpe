// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktype

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
	summary                   = "read network type data source"
	ErrorNoValidSearchTerms   = `no valid search terms - an id or name is required`
	ErrorRunningPreApply      = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkTypeFound   = `no network type found`
	ErrorMultipleNetworkTypes = `multiple network types were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkTypeDataSourceSchema(ctx)
}

func networkTypeAsState(
	nt *sdk.GetNetworkType200ResponseNetworkType,
) NetworkTypeModel {
	return NetworkTypeModel{
		Id:                     convert.Int64ToType(nt.Id),
		Name:                   convert.StrToType(nt.Name),
		Code:                   convert.StrToType(nt.Code),
		Description:            convert.StrToType(nt.Description.Get()),
		Category:               convert.StrToType(nt.Category.Get()),
		ExternalType:           convert.StrToType(nt.ExternalType.Get()),
		Creatable:              convert.BoolToType(nt.Creatable),
		Deletable:              convert.BoolToType(nt.Deletable),
		Overlay:                convert.BoolToType(nt.Overlay),
		NameEditable:           convert.BoolToType(nt.NameEditable),
		CidrRequired:           convert.BoolToType(nt.CidrRequired),
		CidrEditable:           convert.BoolToType(nt.CidrEditable),
		DhcpServerEditable:     convert.BoolToType(nt.DhcpServerEditable),
		DnsEditable:            convert.BoolToType(nt.DnsEditable),
		GatewayEditable:        convert.BoolToType(nt.GatewayEditable),
		VlanIdEditable:         convert.BoolToType(nt.VlanIdEditable),
		StaticOverrideEditable: convert.BoolToType(nt.StaticOverrideEditable),
		NetworkDomainEditable:  convert.BoolToType(nt.NetworkDomainEditable),
		CanAssignPool:          convert.BoolToType(nt.CanAssignPool),
		HasNetworkServer:       convert.BoolToType(nt.HasNetworkServer),
		HasCidr:                convert.BoolToType(nt.HasCidr),
		HasStaticRoutes:        convert.BoolToType(nt.HasStaticRoutes),
		HasFloatingIps:         convert.BoolToType(nt.HasFloatingIps),
	}
}

func getNetworkTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkType200ResponseNetworkType, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkType(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkType == nil {
		return nil, fmt.Errorf("GET failed for network type %d: response missing networkType", id)
	}

	nt := *r.NetworkType

	return &nt, nil
}

func getNetworkTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkType200ResponseNetworkType, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.NetworkTypes)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling network types: %w", marshalErr)
	}

	var networkTypes []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &networkTypes); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding network types: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, nt := range networkTypes {
		if nt.Name != nil && *nt.Name == name {
			if nt.Id != nil {
				matchedID = *nt.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkTypeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkTypes)
	}

	return getNetworkTypeByID(ctx, matchedID, apiClient)
}

func getNetworkType(
	ctx context.Context,
	config *NetworkTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkType200ResponseNetworkType, error) {
	if !config.Id.IsNull() {
		return getNetworkTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkTypeModel

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

	nt, err := getNetworkType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkTypeAsState(nt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
