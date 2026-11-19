// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

const readOperation = "read network router resource"

func getRouterAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan NetworkRouterModel,
) (NetworkRouterModel, bool, diag.Diagnostics) {
	var state NetworkRouterModel
	var diags diag.Diagnostics
	importing := plan.Name.IsNull()
	resp, hresp, err := client.NetworksAPI.GetNetworkRouter(ctx, id).Execute()
	if errfmt.IsNotFound(hresp) {
		return state, false, diags
	}
	if err != nil || resp == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("network router %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, false, diags
	}

	router := resp.NetworkRouter
	if router == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("network router %d GET returned no networkRouter payload", id),
		)

		return state, false, diags
	}

	state.Id = convert.Int64ToType(router.Id)
	state.Name = convert.StrToType(router.Name)
	state.Code = convert.StrToType(router.Code)
	state.Enabled = convert.BoolToType(router.Enabled)
	state.EnableBgp = convert.BoolToType(router.EnableBgp)
	state.ProviderId = convert.StrToType(router.ProviderId)

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
		state.TypeId = routerTypeID(router, plan.TypeId)
		state.ConfigNsxtGatewayTier0 = plan.ConfigNsxtGatewayTier0
	case !plan.ConfigNsxtGatewayTier1.IsNull() && !plan.ConfigNsxtGatewayTier1.IsUnknown():
		// read type ID from API if using static config (it's not known at plan)
		state.TypeId = routerTypeID(router, plan.TypeId)
		state.ConfigNsxtGatewayTier1 = plan.ConfigNsxtGatewayTier1
	case !plan.ConfigJuniperSdn.IsNull() && !plan.ConfigJuniperSdn.IsUnknown():
		// The SDN blocks update in place, so out-of-band changes must show as
		// drift: read them back from the API rather than echoing the plan.
		state.TypeId = routerTypeID(router, plan.TypeId)
		cfg, d := juniperSdnConfigFromMap(ctx, router.Config)
		diags.Append(d...)
		state.ConfigJuniperSdn = cfg
	case !plan.ConfigJuniperBgp.IsNull() && !plan.ConfigJuniperBgp.IsUnknown():
		state.TypeId = routerTypeID(router, plan.TypeId)
		cfg, d := juniperBgpConfigFromMap(ctx, router.Config)
		diags.Append(d...)
		// md5_key_wo is write-only and reads back as null (the API never
		// returns it either); wodigest tracks changes to it.
		state.ConfigJuniperBgp = cfg
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
		} else if zone := router.Zone.Get(); zone != nil && zone.Id != nil {
			state.CloudId = convert.Int64ToType(zone.Id)
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

	// write_only_revision is the provider's own marker; Create and Update set
	// it after this read, a refresh keeps it.
	state.WriteOnlyRevision = plan.WriteOnlyRevision

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

	return state, true, diags
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
	case codeJuniperSdnRouter:
		cfg, d := juniperSdnConfigFromMap(ctx, config)
		diags.Append(d...)
		state.ConfigJuniperSdn = cfg
	case codeJuniperBgpSdnRouter:
		cfg, d := juniperBgpConfigFromMap(ctx, config)
		diags.Append(d...)
		state.ConfigJuniperBgp = cfg
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

	newState, found, diags := getRouterAsState(ctx, id, client, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The router was deleted out of band (Morpheus returned 404); remove it
	// from state so the next plan recreates it instead of erroring.
	if !found {
		tflog.Info(ctx, fmt.Sprintf("network router %d not found, removing from state", id))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, NetworkRouterResourceSchema(ctx), &resp.State)...,
	)
}

// routerTypeID is the router's type id as the API returned it, or the prior
// value when the response carries no type: type is an optional pointer, and a
// typed block's type id is only known from the appliance. An unknown prior —
// the create plan's, since type_id is computed — becomes null, which can be
// written to state.
func routerTypeID(router *sdk.GetNetworkRouter200ResponseNetworkRouter, prior types.Int64) types.Int64 {
	if router.Type != nil && router.Type.Id != nil {
		return convert.Int64ToType(router.Type.Id)
	}

	if prior.IsUnknown() {
		return types.Int64Null()
	}

	return prior
}
