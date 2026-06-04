// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktype

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_type"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkTypeDataSourceSchema(ctx)
}

func getNetworkType(
	ctx context.Context,
	config NetworkTypeModel,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	if !config.Id.IsNull() {
		return getNetworkTypeByID(ctx, config.Id.ValueInt64(), client)
	}
	if !config.Name.IsNull() {
		return getNetworkTypeByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getNetworkTypeByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	response, hresp, err := client.NetworksAPI.GetNetworkType(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network type %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	state := &NetworkTypeModel{}

	nt, ok := response.GetNetworkTypeOk()
	if !ok {
		return nil, fmt.Errorf("network type %d is nil", id)
	}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(nt.Name)
	state.Code = convert.StrToType(nt.Code)
	state.Description = convert.StrToType(nt.Description.Get())
	state.Category = convert.StrToType(nt.Category.Get())
	state.Creatable = convert.BoolToType(nt.Creatable)

	return state, nil
}

func getNetworkTypeByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	response, hresp, err := client.NetworksAPI.ListNetworkTypes(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network type %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingNetworkTypes []sdk.ListNetworkTypes200ResponseAllOfNetworkTypesInner
	for _, networkType := range response.GetNetworkTypes() {
		if networkTypeName, ok := networkType.GetNameOk(); ok && *networkTypeName == name {
			matchingNetworkTypes = append(matchingNetworkTypes, networkType)
		}
	}

	if len(matchingNetworkTypes) == 0 {
		return nil, fmt.Errorf("network type %s not found", name)
	}

	if len(matchingNetworkTypes) > 1 {
		var networkTypeIDs []string
		for _, n := range matchingNetworkTypes {
			if id, ok := n.GetIdOk(); ok {
				networkTypeIDs = append(networkTypeIDs, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple network types found with name %s. Network type IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(networkTypeIDs, ", "),
		)
	}

	id, ok := matchingNetworkTypes[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("network type %s has missing ID", name)
	}

	return getNetworkTypeByID(ctx, *id, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkTypeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network type data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getNetworkType(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network type data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
