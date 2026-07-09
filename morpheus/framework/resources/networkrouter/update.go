// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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

	router := &sdk.UpdateNetworkRouterRequestNetworkRouter{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		router.Name = plan.Name.ValueStringPointer()
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		router.Enabled = plan.Enabled.ValueBoolPointer()
	}

	// Build config map for update
	// Only generic config supports update
	configMap := buildUpdateConfig(ctx, plan)
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
	resp.Diagnostics.Append(applyRouterPermissions(ctx, id, plan, client)...)
	if resp.Diagnostics.HasError() {
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

// applyRouterPermissions calls UpdateNetworkRouterPermissions when visibility or tenant_ids
// are set in the plan. A 403 response is an error because group-scoped routers do not
// support tenant permissions and the caller must resolve the conflict.
func applyRouterPermissions(
	ctx context.Context,
	id int64,
	plan NetworkRouterModel,
	client *sdk.APIClient,
) diag.Diagnostics {
	var diags diag.Diagnostics

	visSet := !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown()
	tenSet := !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown()
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
