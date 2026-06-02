// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const updateOperation = "update network router resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			updateOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	router := sdk.NewUpdateNetworkRouterRequestNetworkRouterWithDefaults()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		router.Name = &name
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := plan.Enabled.ValueBool()
		router.Enabled = &enabled
	}

	// Build config map for update
	// Only generic config supports update
	configMap := buildUpdateConfig(ctx, plan)
	if configMap != nil {
		router.Config = configMap
	}

	updateReq := sdk.NewUpdateNetworkRouterRequestWithDefaults()
	updateReq.NetworkRouter = router

	_, hresp, err := client.NetworksAPI.UpdateNetworkRouter(ctx, id).
		UpdateNetworkRouterRequest(*updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("network router %d PUT failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	newState, diags := getRouterAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve user-specified config in state
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		newState.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func buildUpdateConfig(ctx context.Context, plan NetworkRouterModel) map[string]any {
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()

		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			return nil
		}

		configDataMap, ok := configMap.(map[string]any)
		if ok {
			return configDataMap
		}

		return nil
	}

	return nil
}
