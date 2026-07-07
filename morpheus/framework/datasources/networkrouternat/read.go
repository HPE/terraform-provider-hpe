// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func natAsState(
	nat *sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT,
	routerID int64,
) NetworkRouterNatModel {
	// Id is *int32 in SDK — convert to *int64
	var idPtr *int64
	if nat.Id != nil {
		i := int64(*nat.Id)
		idPtr = &i
	}

	// Priority is *int32 in SDK — convert to *int64
	var priorityPtr *int64
	if nat.Priority != nil {
		p := int64(*nat.Priority)
		priorityPtr = &p
	}

	state := NetworkRouterNatModel{
		Id:                 convert.Int64ToType(idPtr),
		RouterId:           types.Int64Value(routerID),
		Name:               convert.StrToType(nat.Name),
		Action:             convert.StrToType(nat.Action),
		Description:        convert.StrToType(nat.Description),
		Enabled:            convert.BoolToType(nat.Enabled),
		SourceNetwork:      convert.StrToType(nat.SourceNetwork),
		DestinationNetwork: convert.StrToType(nat.DestinationNetwork.Get()),
		TranslatedNetwork:  convert.StrToType(nat.TranslatedNetwork),
		SourcePorts:        convert.StrToType(nat.SourcePorts.Get()),
		DestinationPorts:   convert.StrToType(nat.DestinationPorts.Get()),
		TranslatedPorts:    convert.StrToType(nat.TranslatedPorts.Get()),
		Priority:           convert.Int64ToType(priorityPtr),
		Protocol:           convert.StrToType(nat.Protocol.Get()),
		ExternalId:         convert.StrToType(nat.ExternalId),
		ProviderId:         convert.StrToType(nat.ProviderId),
		SyncSource:         convert.StrToType(nat.SyncSource),
		DateCreated:        convert.TimeToType(nat.DateCreated),
		LastUpdated:        convert.TimeToType(nat.LastUpdated),
	}

	// RefType — nullable
	if nat.RefType.IsSet() {
		state.RefType = convert.StrToType(nat.RefType.Get())
	} else {
		state.RefType = types.StringNull()
	}

	// RefId — nullable
	if nat.RefId.IsSet() {
		state.RefId = convert.StrToType(nat.RefId.Get())
	} else {
		state.RefId = types.StringNull()
	}

	// InternalId — nullable
	if nat.InternalId.IsSet() {
		state.InternalId = convert.StrToType(nat.InternalId.Get())
	} else {
		state.InternalId = types.StringNull()
	}

	// MatchIpv6DestinationPrefix — nullable
	if nat.MatchIpv6DestinationPrefix.IsSet() {
		state.MatchIpv6destinationPrefix = convert.StrToType(nat.MatchIpv6DestinationPrefix.Get())
	} else {
		state.MatchIpv6destinationPrefix = types.StringNull()
	}

	// TranslatedIpv4SourcePrefix — nullable
	if nat.TranslatedIpv4SourcePrefix.IsSet() {
		state.TranslatedIpv4sourcePrefix = convert.StrToType(nat.TranslatedIpv4SourcePrefix.Get())
	} else {
		state.TranslatedIpv4sourcePrefix = types.StringNull()
	}

	return state
}

func getNatByID(
	ctx context.Context,
	id int64,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterNat(
		ctx, id, routerID,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router nat %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.NetworkRouterNAT, nil
}

func getNatByName(
	ctx context.Context,
	name string,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRoutersNats(
		ctx, routerID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router nats with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.NetworkRouterNATs
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterNatFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, int64(*items[i].Id))
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterNatFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouterNats)
	}

	return getNatByID(ctx, matchedIDs[0], routerID, apiClient)
}

func getNat(
	ctx context.Context,
	config *NetworkRouterNatModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	routerID := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getNatByID(ctx, config.Id.ValueInt64(), routerID, apiClient)
	} else if !config.Name.IsNull() {
		return getNatByName(ctx, config.Name.ValueString(), routerID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterNatModel

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

	nat, err := getNat(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	routerID := config.RouterId.ValueInt64()
	state := natAsState(nat, routerID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
