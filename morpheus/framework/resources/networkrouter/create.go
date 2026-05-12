// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const createOperation = "create network router resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	name := plan.Name.ValueString()

	router := sdk.NewCreateNetworkRouterRequestNetworkRouterWithDefaults()
	router.SetName(name)

	// Set type (required)
	switch {
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		routerType := sdk.NewCreateNetworkRouterRequestNetworkRouterType(plan.TypeId.ValueInt64())
		router.SetType(*routerType)
	case !plan.ConfigNsxtGatewayTier0.IsNull() && !plan.ConfigNsxtGatewayTier0.IsUnknown():
		typeId, err := typeIdFromCode(ctx, client, codeNSXTTier0Gateway)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"failed to find network type: "+err.Error(),
			)

			return
		}
		routerType := sdk.NewCreateNetworkRouterRequestNetworkRouterType(*typeId)
		router.SetType(*routerType)

	case !plan.ConfigNsxtGatewayTier1.IsNull() && !plan.ConfigNsxtGatewayTier1.IsUnknown():
		typeId, err := typeIdFromCode(ctx, client, codeNSXTTier1Gateway)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"failed to find network type: "+err.Error(),
			)

			return
		}
		routerType := sdk.NewCreateNetworkRouterRequestNetworkRouterType(*typeId)
		router.SetType(*routerType)
	}

	// Set site (group_id)
	groupID := plan.GroupId.ValueInt64()
	site := sdk.NewCreateNetworkRouterRequestNetworkRouterSite(
		sdk.Int64AsCreateNetworkRouterRequestNetworkRouterSiteId(&groupID),
	)
	router.SetSite(*site)

	// Set enabled
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		router.SetEnabled(plan.Enabled.ValueBool())
	}

	// Set zone (cloud_id) if provided
	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		zone := sdk.NewCreateNetworkRouterRequestNetworkRouterZoneWithDefaults()
		zone.SetId(plan.CloudId.ValueInt64())
		router.SetZone(*zone)
	}

	// Set networkServer (network_integration_id) if provided
	if !plan.NetworkIntegrationId.IsNull() && !plan.NetworkIntegrationId.IsUnknown() {
		ns := sdk.NewCreateNetworkRouterRequestNetworkRouterNetworkServerWithDefaults()
		ns.SetId(plan.NetworkIntegrationId.ValueInt64())
		router.SetNetworkServer(*ns)
	}

	// Set config from the dynamic config attribute or typed config blocks
	routerConfig, diags := buildRouterConfig(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if routerConfig != nil {
		router.SetConfig(*routerConfig)
	}

	createReq := sdk.NewCreateNetworkRouterRequestWithDefaults()
	createReq.SetNetworkRouter(*router)

	result, hresp, err := client.NetworksAPI.CreateNetworkRouter(ctx).
		CreateNetworkRouterRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("network router %s POST failed: %s",
				name, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if !result.IsSetId() || result.GetId() == 0 {
		resp.Diagnostics.AddError(
			createOperation,
			"network router "+name+": id is nil or zero",
		)

		return
	}

	id := result.GetId()
	plan.Id = types.Int64Value(id)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_router",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getRouterAsState(ctx, id, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"failed to read network router state",
			fmt.Sprintf("Network router %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	// Preserve user-specified config in state to avoid spurious diffs
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set network router state",
			fmt.Sprintf("Network router %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}

func buildRouterConfig(
	ctx context.Context,
	plan NetworkRouterModel,
) (*sdk.CreateNetworkRouterRequestNetworkRouterConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch {
	// Typed NSX Tier0 config
	case !plan.ConfigNsxtGatewayTier0.IsNull() && !plan.ConfigNsxtGatewayTier0.IsUnknown():
		cfg := nsxTier0Config(plan.ConfigNsxtGatewayTier0)

		return &cfg, diags
	// Typed NSX Tier1 config
	case !plan.ConfigNsxtGatewayTier1.IsNull() && !plan.ConfigNsxtGatewayTier1.IsUnknown():
		cfg := nsxTier1Config(plan.ConfigNsxtGatewayTier1)

		return &cfg, diags

	// Dynamic config is the fallback
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()

		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			diags.AddError(
				createOperation,
				"failed to convert config: "+err.Error(),
			)

			return nil, diags
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			diags.AddError(
				createOperation,
				"config must be a valid object/map",
			)

			return nil, diags
		}

		routerConfig := sdk.CreateNetworkRouterRequestNetworkRouterConfig{}
		routerConfig.MapmapOfStringAny = &configDataMap

		return &routerConfig, diags

	default:
		return nil, diags
	}
}

func typeIdFromCode(ctx context.Context, client *sdk.APIClient, code string) (*int64, error) {
	res, hresp, err := client.NetworksAPI.ListNetworkRouterTypes(ctx).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"network router types GET failed: %s",
			errfmt.ErrMsg(err, hresp))
	}

	types := res.GetNetworkRouterTypes()
	for _, t := range types {
		if t.GetCode() == code {
			if id, ok := t.GetIdOk(); ok {
				return id, nil
			} else {
				return nil, fmt.Errorf("Network router type id for code %s is nil", code)
			}
		}
	}

	return nil, fmt.Errorf("Could not find network router type for code %s.\nThe network integration for the type may not yet be configured on the Morpheus appliance.", code)
}
