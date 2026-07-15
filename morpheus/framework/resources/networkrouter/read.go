// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
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

		// cloud_id is only meaningful for router types that do not use a network
		// integration. When a networkServer is present the API derives the zone
		// from it, and a freshly-applied resource leaves cloud_id unset, so only
		// surface the zone as cloud_id when there is no integration - otherwise
		// import would not match apply.
		if router.NetworkServer != nil && router.NetworkServer.Id != nil {
			state.NetworkIntegrationId = convert.Int64ToType(router.NetworkServer.Id)
		} else if router.Zone != nil && router.Zone.Id != nil {
			state.CloudId = convert.Int64ToType(router.Zone.Id)
		}

		// On import, site == null means shared group access
		if router.Site != nil {
			state.GroupId = convert.Int64ToType(router.Site.Id)
		} else {
			state.SharedGroupAccess = types.BoolValue(true)
		}

		// Hydrate the config from the API. NSX-T gateway types use their typed
		// config block (matching how they are applied); any other type uses the
		// generic dynamic config.
		var typeCode string
		if router.Type != nil && router.Type.Code != nil {
			typeCode = *router.Type.Code
		}
		diags.Append(hydrateImportedConfig(ctx, typeCode, router.Config, &state)...)
	}

	// Preserve tenant_ids from prior state on normal refresh. The API may
	// silently drop IDs that don't exist in the environment. On import there
	// is no prior state, so we read from the API response instead.
	state.TenantIds = plan.TenantIds

	// Read permissions: visibility + tenant_ids
	if router.Permissions != nil {
		if router.Permissions.Visibility != nil {
			state.Visibility = types.StringValue(*router.Permissions.Visibility)
		}
	}

	// tenant_ids and visibility are computed and must be known after apply.
	// On create the plan values are unknown (they conflict with the required
	// group_id, so a user can never set them), and on import there is no prior
	// state. In both cases resolve from the API response; a group-scoped router
	// has no tenant permissions, which yields an empty set.
	if importing || state.TenantIds.IsUnknown() {
		var tenantIDs []int64
		if router.Permissions != nil && router.Permissions.TenantPermissions != nil {
			tenantIDs = router.Permissions.TenantPermissions.Accounts
		}
		setVal, setDiags := types.SetValueFrom(ctx, types.Int64Type, tenantIDs)
		diags.Append(setDiags...)
		state.TenantIds = setVal
	}

	// Guard against visibility remaining unknown after apply when the API
	// response carries no permissions block.
	if state.Visibility.IsUnknown() {
		state.Visibility = types.StringNull()
	}

	return state, diags
}

// hydrateImportedConfig populates the config representation on import. The typed
// NSX-T blocks are used for the NSX-T gateway types (so import matches how they
// are applied); any other type uses the generic dynamic config.
func hydrateImportedConfig(
	ctx context.Context,
	typeCode string,
	config map[string]interface{},
	state *NetworkRouterModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	switch typeCode {
	case codeNSXTTier0Gateway:
		cfg, d := tier0ConfigFromMap(ctx, config)
		diags.Append(d...)
		state.ConfigNsxtGatewayTier0 = cfg
	case codeNSXTTier1Gateway:
		cfg, d := tier1ConfigFromMap(ctx, config)
		diags.Append(d...)
		state.ConfigNsxtGatewayTier1 = cfg
	default:
		if len(config) > 0 {
			dyn, err := convert.MapToDynamic(ctx, config)
			if err != nil {
				diags.AddError(readOperation,
					"failed to map network router config: "+err.Error())
			} else {
				state.Config = dyn
			}
		}
	}

	return diags
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
