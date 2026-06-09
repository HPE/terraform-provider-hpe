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
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const readOperation = "read network router resource"

func getRouterAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan NetworkRouterModel,
) (NetworkRouterModel, diag.Diagnostics) {
	var state NetworkRouterModel
	var diags diag.Diagnostics
	importing := plan.Name.IsNull()
	resp, hresp, err := client.NetworksAPI.GetNetworkRouter(ctx, id).Execute()
	if err != nil || resp == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("network router %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	router := resp.NetworkRouter
	if router == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("network router %d GET returned no networkRouter payload", id),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(router.Id)
	state.Name = convert.StrToType(router.Name)
	state.Code = convert.StrToType(router.Code)
	state.Enabled = convert.BoolToType(router.Enabled)
	state.EnableBgp = convert.BoolToType(router.EnableBgp)

	// Preserve plan values for immutable fields
	state.GroupId = plan.GroupId
	state.CloudId = plan.CloudId
	state.NetworkIntegrationId = plan.NetworkIntegrationId

	switch {
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		// Use the type ID we set in config if it's a generic config.
		state.TypeId = plan.TypeId
		// Config: preserve from plan - we only care about what was set by user.
		// We don't want to read everything back from API as there may be stuff
		// that exists in remote `config` not set by user.
		state.Config = plan.Config
	case !plan.ConfigNsxtGatewayTier0.IsNull() && !plan.ConfigNsxtGatewayTier0.IsUnknown():
		// read type ID from API if using static config (it's not known at plan)
		state.TypeId = convert.Int64ToType(router.Type.Id)
		state.ConfigNsxtGatewayTier0 = plan.ConfigNsxtGatewayTier0
	case !plan.ConfigNsxtGatewayTier1.IsNull() && !plan.ConfigNsxtGatewayTier1.IsUnknown():
		// read type ID from API if using static config (it's not known at plan)
		state.TypeId = convert.Int64ToType(router.Type.Id)
		state.ConfigNsxtGatewayTier1 = plan.ConfigNsxtGatewayTier1
	}

	// Populate type_id from API response if importing (plan.Name is null)
	if importing {
		if router.Type != nil && router.Type.Id != nil {
			state.TypeId = convert.Int64ToType(router.Type.Id)
		}

		if router.Zone != nil && router.Zone.Id != nil {
			state.CloudId = convert.Int64ToType(router.Zone.Id)
		}

		if router.NetworkServer != nil && router.NetworkServer.Id != nil {
			state.NetworkIntegrationId = convert.Int64ToType(router.NetworkServer.Id)
		}

		// On import, site == null means shared group access
		if router.Site != nil {
			state.GroupId = convert.Int64ToType(router.Site.Id)
		} else {
			state.SharedGroupAccess = types.BoolValue(true)
		}
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkRouterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			readOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	newState, diags := getRouterAsState(ctx, id, client, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
