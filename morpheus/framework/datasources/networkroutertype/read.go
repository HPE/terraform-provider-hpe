// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkroutertype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func networkRouterTypeAsState(
	t *sdk.GetNetworkRouterType200ResponseNetworkRouterType,
) NetworkRouterTypeModel {
	return NetworkRouterTypeModel{
		Id:               convert.Int64ToType(t.Id),
		Name:             convert.StrToType(t.Name),
		Code:             convert.StrToType(t.Code),
		Description:      convert.StrToType(t.Description),
		Enabled:          convert.BoolToType(t.Enabled),
		Creatable:        convert.BoolToType(t.Creatable),
		Selectable:       convert.BoolToType(t.Selectable),
		HasFirewall:      convert.BoolToType(t.HasFirewall),
		HasDhcp:          convert.BoolToType(t.HasDhcp),
		HasRouting:       convert.BoolToType(t.HasRouting),
		HasNetworkServer: convert.BoolToType(t.HasNetworkServer),
	}
}

func getNetworkRouterTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterType200ResponseNetworkRouterType, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterType(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network router type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkRouterType == nil {
		return nil, fmt.Errorf("GET failed for network router type %d: response missing networkRouterType", id)
	}

	nt := *r.NetworkRouterType

	return &nt, nil
}

func getNetworkRouterTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterType200ResponseNetworkRouterType, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkRouterTypes(ctx).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network router type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.NetworkRouterTypes)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling network router types: %w", marshalErr)
	}

	var networkRouterTypes []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &networkRouterTypes); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding network router types: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, nrt := range networkRouterTypes {
		if nrt.Name != nil && *nrt.Name == name {
			if nrt.Id != nil {
				matchedID = *nrt.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkRouterTypeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouterType)
	}

	return getNetworkRouterTypeByID(ctx, matchedID, apiClient)
}

func getNetworkRouterType(
	ctx context.Context,
	config *NetworkRouterTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterType200ResponseNetworkRouterType, error) {
	if !config.Id.IsNull() {
		return getNetworkRouterTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkRouterTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterTypeModel

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

	nt, err := getNetworkRouterType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkRouterTypeAsState(nt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
