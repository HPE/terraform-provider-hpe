// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolservertype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func networkPoolServerTypeAsState(
	t *sdk.GetNetworkPoolServerType200ResponseNetworkPoolServerType,
) NetworkPoolServerTypeModel {
	state := NetworkPoolServerTypeModel{
		Id:         convert.Int64ToType(t.Id),
		Name:       convert.StrToType(t.Name),
		Code:       convert.StrToType(t.Code),
		Enabled:    convert.BoolToType(t.Enabled),
		Selectable: convert.BoolToType(t.Selectable),
	}

	if t.Description.IsSet() {
		state.Description = convert.StrToType(t.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if t.IntegrationCode.IsSet() {
		state.IntegrationCode = convert.StrToType(t.IntegrationCode.Get())
	} else {
		state.IntegrationCode = types.StringNull()
	}

	if t.PoolService.IsSet() {
		state.PoolService = convert.StrToType(t.PoolService.Get())
	} else {
		state.PoolService = types.StringNull()
	}

	return state
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
