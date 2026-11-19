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
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const updateOperation = "update network router resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config NetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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

	router := &sdk.UpdateNetworkRouterRequestNetworkRouter{}

	// Only send the name when it changed: the HPE Morpheus SDN router types
	// reject any update that carries the name field ("Field name cannot be
	// changed"), even with the current value, and no type needs it re-sent.
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		router.Name = plan.Name.ValueStringPointer()
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		router.Enabled = plan.Enabled.ValueBoolPointer()
	}

	// Write-only values are re-sent only when their digest says they changed.
	resend := wodigest.Resend(ctx, req.Plan)

	configMap := buildUpdateConfig(ctx, plan, configuredMd5Key(config), resend, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if configMap != nil {
		router.Config = configMap
	}

	updateReq := &sdk.UpdateNetworkRouterRequest{}
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

	// Apply permissions (visibility + tenant_ids).
	resp.Diagnostics.Append(applyRouterPermissions(ctx, id, plan, config, client)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, found, diags := getRouterAsState(ctx, id, client, plan)
	if !found && !diags.HasError() {
		diags.AddError("update network router resource", fmt.Sprintf("network router %d not found after update", id))
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve user-specified config in state
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		newState.Config = plan.Config
	}

	// Write-only digests: a re-send moves the revision; otherwise it is kept.
	if resend {
		newState.WriteOnlyRevision = wodigest.NewRevision()
	} else {
		newState.WriteOnlyRevision = state.WriteOnlyRevision
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func buildUpdateConfig(
	ctx context.Context, plan NetworkRouterModel, md5Key types.String, sendMd5Key bool, diags *diag.Diagnostics,
) map[string]any {
	// The SDN plugin's router fields are editable (resource_pool_id excepted:
	// it replaces); send them as the map the update request takes. The NSX-T
	// gateway blocks are create-only.
	switch {
	case !plan.ConfigJuniperSdn.IsNull() && !plan.ConfigJuniperSdn.IsUnknown():
		m, d := juniperSdnUpdateConfig(ctx, plan.ConfigJuniperSdn)
		diags.Append(d...)

		return m
	case !plan.ConfigJuniperBgp.IsNull() && !plan.ConfigJuniperBgp.IsUnknown():
		m, d := juniperBgpUpdateConfig(ctx, plan.ConfigJuniperBgp, md5Key, sendMd5Key)
		diags.Append(d...)

		return m
	}

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

// applyRouterPermissions calls UpdateNetworkRouterPermissions when visibility or tenant_ids
// are set in the plan. A 403 response is an error because group-scoped routers do not
// support tenant permissions and the caller must resolve the conflict.
func applyRouterPermissions(
	ctx context.Context,
	id int64,
	plan NetworkRouterModel,
	config NetworkRouterModel,
	client *sdk.APIClient,
) diag.Diagnostics {
	var diags diag.Diagnostics

	visSet := !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() &&
		!config.Visibility.IsNull() && !config.Visibility.IsUnknown()
	tenSet := !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() &&
		!config.TenantIds.IsNull() && !config.TenantIds.IsUnknown()
	if !visSet && !tenSet {
		return diags
	}

	perms := sdk.UpdateNetworkRouterPermissionsRequestPermissions{}

	if visSet {
		perms.Visibility = plan.Visibility.ValueStringPointer()
	}

	if tenSet {
		var ids []int64
		diags.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return diags
		}
		perms.TenantPermissions = &sdk.UpdateNetworkRouterPermissionsRequestPermissionsTenantPermissions{
			Accounts: ids,
		}
	}

	permReq := sdk.UpdateNetworkRouterPermissionsRequest{Permissions: &perms}

	_, hresp, err := client.NetworksAPI.UpdateNetworkRouterPermissions(ctx, id).
		UpdateNetworkRouterPermissionsRequest(permReq).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusForbidden {
			errfmt.DiagErrorf(&diags, errfmt.OpUpdate, "network router permissions",
				"network router %d permissions PUT returned 403"+
					" (group-scoped routers may not support tenant permissions): %s",
				id, errfmt.SafeErrMsg(err, hresp))

			return diags
		}

		errfmt.DiagErrorf(&diags, errfmt.OpUpdate, "network router permissions",
			"network router %d permissions PUT failed: %s", id, errfmt.SafeErrMsg(err, hresp))
	}

	return diags
}
