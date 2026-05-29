// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// dhcpServerReadPayload is a local struct for decoding the untyped GET
// response returned by the SDK (interface{}).
type dhcpServerReadPayload struct {
	Id              *int64                      `json:"id"`
	Name            *string                     `json:"name"`
	ServerIpAddress *string                     `json:"serverIpAddress"`
	LeaseTime       *int64                      `json:"leaseTime"`
	Config          json.RawMessage             `json:"config"`
	NetworkServer   *dhcpServerNetworkServerRef `json:"networkServer"`
}

type dhcpServerNetworkServerRef struct {
	Id *int64 `json:"id"`
}

func getNetworkDhcpServerAsState(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
	plan NetworkDhcpServerModel,
) (NetworkDhcpServerModel, diag.Diagnostics) {
	var state NetworkDhcpServerModel
	var diags diag.Diagnostics

	dhcpResp, hresp, err := client.NetworksAPI.
		GetNetworkDhcpServer(ctx, id, serverID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	raw := dhcpResp.GetNetworkDhcpServer()

	encoded, err := json.Marshal(raw)
	if err != nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to marshal response: %s",
				id, err,
			),
		)

		return state, diags
	}

	var payload dhcpServerReadPayload
	if err := json.Unmarshal(encoded, &payload); err != nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to unmarshal response: %s",
				id, err,
			),
		)

		return state, diags
	}

	if payload.Id == nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET response missing id", id),
		)

		return state, diags
	}

	state.Id = types.Int64Value(*payload.Id)
	state.Name = convert.StrToType(payload.Name)
	state.ServerIpAddress = convert.StrToType(payload.ServerIpAddress)

	// leaseTime may not be returned by the GET response;
	// preserve from plan / prior state when absent.
	if payload.LeaseTime != nil {
		state.LeaseTime = types.Int64Value(*payload.LeaseTime)
	} else {
		state.LeaseTime = plan.LeaseTime
	}

	if payload.NetworkServer != nil && payload.NetworkServer.Id != nil {
		state.NetworkIntegrationId = types.Int64Value(*payload.NetworkServer.Id)
	} else {
		// Some API versions do not include networkServer in this response.
		state.NetworkIntegrationId = plan.NetworkIntegrationId
	}

	configState, cfgDiags := resolveConfigState(ctx, id, payload.Config, plan)
	diags.Append(cfgDiags...)
	if diags.HasError() {
		return state, diags
	}

	state.Config = configState.config
	state.ConfigNsxt = configState.configNsxt

	return state, diags
}

type configResult struct {
	config     types.Dynamic
	configNsxt ConfigNsxtValue
}

// resolveConfigState determines whether the API response contains NSXT config
// or generic config and returns the appropriate Terraform state values.
//
// When the plan already indicates which variant the user wrote (config_nsxt vs
// config), we honour that. During import neither field is set, so we
// auto-detect by attempting an NSXT unmarshal first.
func resolveConfigState(
	ctx context.Context,
	id int64,
	rawConfig json.RawMessage,
	plan NetworkDhcpServerModel,
) (configResult, diag.Diagnostics) {
	var diags diag.Diagnostics

	planHasNsxt := !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown()
	planHasDynamic := !plan.Config.IsNull() && !plan.Config.IsUnknown()

	switch {
	case planHasNsxt:
		nsxt, nsxtDiags := buildNsxtConfigValue(ctx, id, rawConfig)
		diags.Append(nsxtDiags...)

		return configResult{
			config:     types.DynamicNull(),
			configNsxt: nsxt,
		}, diags

	case planHasDynamic:
		return configResult{
			config:     plan.Config,
			configNsxt: NewConfigNsxtValueNull(),
		}, diags

	default:
		// Import path: no plan context — auto-detect from API response.
		return detectConfigFromResponse(ctx, id, rawConfig)
	}
}

// buildNsxtConfigValue unmarshals raw config JSON into a ConfigNsxtValue.
func buildNsxtConfigValue(
	ctx context.Context,
	id int64,
	rawConfig json.RawMessage,
) (ConfigNsxtValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	edgeCluster := types.StringNull()
	activeEdgeNode := types.StringNull()
	standbyEdgeNode := types.StringNull()

	if len(rawConfig) > 0 {
		var nsxCfg sdk.NetworkDhcpServerConfigNSX
		if err := json.Unmarshal(rawConfig, &nsxCfg); err != nil {
			diags.AddWarning(
				"populate network dhcp server resource",
				fmt.Sprintf(
					"network dhcp server %d: failed to unmarshal NSXT config: %s",
					id, err,
				),
			)
		} else {
			edgeCluster = convert.StrToType(nsxCfg.EdgeCluster.Get())
			activeEdgeNode = convert.StrToType(
				nsxCfg.PreferredEdgeNode1.Get(),
			)
			standbyEdgeNode = convert.StrToType(
				nsxCfg.PreferredEdgeNode2.Get(),
			)
		}
	}

	nsxtValue, nsxtDiags := NewConfigNsxtValue(
		ConfigNsxtValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"edge_cluster":      edgeCluster,
			"active_edge_node":  activeEdgeNode,
			"standby_edge_node": standbyEdgeNode,
		},
	)
	diags.Append(nsxtDiags...)

	return nsxtValue, diags
}

// detectConfigFromResponse is used during import when there is no plan
// context. It tries NSXT config first; if any NSXT-specific field is present,
// it populates config_nsxt. Otherwise it falls back to a dynamic value for
// the generic config attribute.
func detectConfigFromResponse(
	ctx context.Context,
	id int64,
	rawConfig json.RawMessage,
) (configResult, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := configResult{
		config:     types.DynamicNull(),
		configNsxt: NewConfigNsxtValueNull(),
	}

	if len(rawConfig) == 0 {
		return result, diags
	}

	var nsxCfg sdk.NetworkDhcpServerConfigNSX
	if err := json.Unmarshal(rawConfig, &nsxCfg); err == nil && isNsxtConfig(&nsxCfg) {
		nsxt, nsxtDiags := buildNsxtConfigValue(ctx, id, rawConfig)
		diags.Append(nsxtDiags...)
		result.configNsxt = nsxt

		return result, diags
	}

	// Fall back to generic dynamic config.
	var configMap map[string]any
	if err := json.Unmarshal(rawConfig, &configMap); err != nil {
		diags.AddWarning(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to unmarshal config as map: %s",
				id, err,
			),
		)

		return result, diags
	}

	dynVal, err := convert.MapToDynamic(ctx, configMap)
	if err != nil {
		diags.AddWarning(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to convert config to dynamic value: %s",
				id, err,
			),
		)

		return result, diags
	}

	result.config = dynVal

	return result, diags
}

// isNsxtConfig returns true when at least one NSXT-specific field is present
// in the decoded config, distinguishing it from an arbitrary generic map.
func isNsxtConfig(cfg *sdk.NetworkDhcpServerConfigNSX) bool {
	return cfg.EdgeCluster.IsSet() ||
		cfg.PreferredEdgeNode1.IsSet() ||
		cfg.PreferredEdgeNode2.IsSet()
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkDhcpServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network dhcp server resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	serverID := plan.NetworkIntegrationId.ValueInt64()

	state, diags := getNetworkDhcpServerAsState(
		ctx, id, serverID, client, plan,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
