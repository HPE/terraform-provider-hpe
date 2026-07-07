// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func networkPoolServerAsState(
	server *sdk.GetNetworkPoolServer200ResponseNetworkPoolServer,
) NetworkPoolServerModel {
	state := NetworkPoolServerModel{
		Id:              convert.Int64ToType(server.Id),
		Name:            convert.StrToType(server.Name),
		Enabled:         convert.BoolToType(server.Enabled),
		Status:          convert.StrToType(server.Status),
		ServiceUrl:      convert.StrToType(server.ServiceUrl.Get()),
		ServiceMode:     convert.StrToType(server.ServiceMode.Get()),
		ServiceUsername: convert.StrToType(server.ServiceUsername.Get()),
		NetworkFilter:   convert.StrToType(server.NetworkFilter.Get()),
		ZoneFilter:      convert.StrToType(server.ZoneFilter.Get()),
		TenantMatch:     convert.StrToType(server.TenantMatch.Get()),
	}

	if server.IgnoreSsl.IsSet() {
		state.IgnoreSsl = convert.BoolToType(server.IgnoreSsl.Get())
	} else {
		state.IgnoreSsl = types.BoolNull()
	}

	if server.ServiceThrottleRate.IsSet() {
		state.ServiceThrottleRate = convert.Int64ToType(server.ServiceThrottleRate.Get())
	} else {
		state.ServiceThrottleRate = types.Int64Null()
	}

	// service_host (nullable)
	if server.ServiceHost.IsSet() {
		state.ServiceHost = convert.StrToType(server.ServiceHost.Get())
	} else {
		state.ServiceHost = types.StringNull()
	}

	// service_port (nullable int32 in SDK -> int64 in schema)
	if server.ServicePort.IsSet() {
		v := int64(*server.ServicePort.Get())
		state.ServicePort = convert.Int64ToType(&v)
	} else {
		state.ServicePort = types.Int64Null()
	}

	// status_date (nullable time)
	if server.StatusDate.IsSet() {
		t := server.StatusDate.Get().String()
		state.StatusDate = convert.StrToType(&t)
	} else {
		state.StatusDate = types.StringNull()
	}

	// date_created
	if server.DateCreated != nil {
		t := server.DateCreated.String()
		state.DateCreated = convert.StrToType(&t)
	} else {
		state.DateCreated = types.StringNull()
	}

	// last_updated
	if server.LastUpdated != nil {
		t := server.LastUpdated.String()
		state.LastUpdated = convert.StrToType(&t)
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Handle dynamic config via JSON round-trip (mirrors networkdhcpserver pattern)
	state.Config = types.DynamicNull()
	if server.Config != nil {
		raw, err := json.Marshal(server.Config)
		if err == nil && string(raw) != "{}" && string(raw) != "null" {
			state.Config = types.DynamicValue(types.StringValue(string(raw)))
		}
	}

	// type nested object {id, code, name}
	state.Type = NewTypeValueNull()
	if server.Type != nil {
		state.Type = TypeValue{
			Id:    convert.Int64ToType(server.Type.Id),
			Code:  convert.StrToType(server.Type.Code),
			Name:  convert.StrToType(server.Type.Name),
			state: attr.ValueStateKnown,
		}
	}

	return state
}

func getPoolServerByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServer200ResponseNetworkPoolServer, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network pool server %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.NetworkPoolServer == nil {
		return nil, fmt.Errorf("GET failed for network pool server %d: empty response", id)
	}

	return r.NetworkPoolServer, nil
}

func getPoolServerByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServer200ResponseNetworkPoolServer, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkPoolServers(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network pool server %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matchedIDs []int64

	for i := range rs.NetworkPoolServers {
		if rs.NetworkPoolServers[i].Name == nil || *rs.NetworkPoolServers[i].Name != name {
			continue
		}
		if rs.NetworkPoolServers[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *rs.NetworkPoolServers[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkPoolServerFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkPoolServer)
	}

	return getPoolServerByID(ctx, matchedIDs[0], apiClient)
}

func getPoolServer(
	ctx context.Context,
	config *NetworkPoolServerModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkPoolServer200ResponseNetworkPoolServer, error) {
	if !config.Id.IsNull() {
		return getPoolServerByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getPoolServerByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkPoolServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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

	server, err := getPoolServer(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkPoolServerAsState(server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
