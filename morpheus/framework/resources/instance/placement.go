// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const placementStrategyAuto = "auto"

// placementServerIDs returns the node server ids of an instance, read from
// instance.containerDetails[].server.id. HVM placement is an action against
// each node's compute server (its server id), not the instance id.
func placementServerIDs(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
) ([]int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"host placement",
			fmt.Sprintf("instance %d: failed to read node servers: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return nil, diags
	}

	if resp.Instance == nil {
		diags.AddError(
			"host placement",
			fmt.Sprintf("instance %d: GET returned nil instance", instanceID),
		)

		return nil, diags
	}

	ids := make([]int64, 0, len(resp.Instance.ContainerDetails))
	for _, cont := range resp.Instance.ContainerDetails {
		if cont.Server != nil && cont.Server.Id != nil {
			ids = append(ids, *cont.Server.Id)
		}
	}

	return ids, diags
}

// buildManageHostPlacementRequest builds the ManageHostPlacement request body
// from a strategy and an optional preferred host id. For the "auto" strategy
// the preferred host is omitted so the appliance clears it; for failover and
// pinned the preferred host is sent when known.
func buildManageHostPlacementRequest(
	strategy string,
	preferredHostID *int64,
) sdk.ManageHostPlacementRequest {
	server := &sdk.ManageHostPlacementRequestServer{
		PlacementStrategy: sdk.PtrString(strategy),
	}

	if strategy != placementStrategyAuto && preferredHostID != nil {
		server.PreferredParentServer = &sdk.ManageHostPlacementRequestServerPreferredParentServer{
			Id: preferredHostID,
		}
	}

	return sdk.ManageHostPlacementRequest{Server: server}
}

// manageHostPlacement issues ManageHostPlacement for every node server id of an
// instance. The appliance returns 200 immediately and may dispatch an async
// live migration; it rejects a host that is not in the VM's cluster, not
// provisioned, or in maintenance, both via an HTTP error and via a failing
// envelope on the 200 body, so both are checked.
func manageHostPlacement(
	ctx context.Context,
	client *sdk.APIClient,
	serverIDs []int64,
	strategy string,
	preferredHostID *int64,
) diag.Diagnostics {
	var diags diag.Diagnostics

	body := buildManageHostPlacementRequest(strategy, preferredHostID)

	for _, serverID := range serverIDs {
		sid := serverID
		result, hresp, err := client.HostsAPI.
			ManageHostPlacement(ctx, sdk.UpdateHostIdParameter{Int64: &sid}).
			ManageHostPlacementRequest(body).
			Execute()
		if cerr := errfmt.CheckResponse(err, hresp); cerr != nil {
			diags.AddError(
				"host placement",
				fmt.Sprintf("server %d: ManageHostPlacement failed: %s", sid, cerr.Error()),
			)

			return diags
		}

		if result != nil && result.Success != nil && !*result.Success {
			diags.AddError(
				"host placement",
				fmt.Sprintf("server %d: the appliance refused the host placement "+
					"(the preferred host may not be in the VM's cluster, not "+
					"provisioned, or in maintenance)", sid),
			)

			return diags
		}

		// The 200 only means the request was accepted: when the preferred
		// host differs from the VM's current host the appliance dispatches an
		// asynchronous serverMove process (a live migration where possible).
		// That process is not correlated to this call in the response, so wait
		// on the authoritative signal instead -- the server's parentServer --
		// within the caller's context, which Create/Update bound with their
		// timeouts. A migration that never lands is reported rather than
		// silently recorded as done.
		if preferredHostID != nil {
			if err := waitForParentHost(ctx, client, sid, *preferredHostID); err != nil {
				diags.AddError(
					"host placement",
					fmt.Sprintf("server %d: %s", sid, err.Error()),
				)

				return diags
			}
		}
	}

	return diags
}

// placementPollInterval is how often the server is re-read while a migration
// is in flight.
const placementPollInterval = 10 * time.Second

// waitForParentHost polls the server until its parentServer is the preferred
// host, or the context's deadline passes. A migration can take minutes; the
// appliance refuses a second move while one is in progress, so polling is the
// only way to learn the outcome.
func waitForParentHost(ctx context.Context, client *sdk.APIClient, serverID, hostID int64) error {
	for {
		resp, hresp, err := client.HostsAPI.
			GetHost(ctx, sdk.GetHostIdParameter{Int64: &serverID}).
			Execute()
		if err == nil && hresp != nil && hresp.StatusCode == http.StatusOK &&
			resp != nil && resp.Server != nil &&
			resp.Server.ParentServer != nil && resp.Server.ParentServer.Id != nil &&
			*resp.Server.ParentServer.Id == hostID {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("the VM did not arrive on host %d before the timeout: %w "+
				"(the migration may still be running on the appliance; check the server's processes)",
				hostID, ctx.Err())
		case <-time.After(placementPollInterval):
		}
	}
}

// readHostPlacement reads placement_strategy and preferred_host_id from a node
// server via GET /api/servers/{id}. Returns (null, null) when the server has no
// placement set.
func readHostPlacement(
	ctx context.Context,
	client *sdk.APIClient,
	serverID int64,
) (types.String, types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, hresp, err := client.HostsAPI.
		GetHost(ctx, sdk.GetHostIdParameter{Int64: &serverID}).
		Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"host placement",
			fmt.Sprintf("server %d: GET failed: %s", serverID, errfmt.ErrMsg(err, hresp)),
		)

		return types.StringNull(), types.Int64Null(), diags
	}

	strategy, preferred := mapServerPlacement(resp.Server)

	return strategy, preferred, diags
}

// mapServerPlacement maps a server's placement fields to Terraform types,
// returning (null, null) for a nil server or unset fields.
func mapServerPlacement(server *sdk.GetHost200ResponseServer) (types.String, types.Int64) {
	if server == nil {
		return types.StringNull(), types.Int64Null()
	}

	strategy := convert.StrToType(server.PlacementStrategy.Get())

	preferred := types.Int64Null()
	if server.PreferredParentServer != nil && server.PreferredParentServer.Id != nil {
		preferred = convert.Int64ToType(server.PreferredParentServer.Id)
	}

	return strategy, preferred
}

// applyHostPlacement reads an instance's node server ids and applies the given
// placement to each. serverIDs are resolved from the instance when nil.
func applyHostPlacement(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	strategy string,
	preferredHostID *int64,
) diag.Diagnostics {
	serverIDs, diags := placementServerIDs(ctx, client, instanceID)
	if diags.HasError() {
		return diags
	}

	if len(serverIDs) == 0 {
		diags.AddError(
			"host placement",
			fmt.Sprintf("instance %d: no node servers found to place", instanceID),
		)

		return diags
	}

	return manageHostPlacement(ctx, client, serverIDs, strategy, preferredHostID)
}

// plannedPlacement reports whether the plan requests an HVM host placement and,
// if so, the strategy and preferred host id to send. Placement is requested
// when config_hvm is set and either placement_strategy or preferred_host_id is
// known and non-null. When only preferred_host_id is given the strategy
// defaults to "failover" (the validation layer requires one of the two to
// accompany a host, and a bare host with no strategy is rejected there).
func plannedPlacement(plan InstanceModel) (string, *int64, bool) {
	if plan.ConfigHvm.IsNull() || plan.ConfigHvm.IsUnknown() {
		return "", nil, false
	}

	ps := plan.ConfigHvm.PlacementStrategy
	ph := plan.ConfigHvm.PreferredHostId

	strategySet := !ps.IsNull() && !ps.IsUnknown()
	preferredSet := !ph.IsNull() && !ph.IsUnknown()

	if !strategySet && !preferredSet {
		return "", nil, false
	}

	strategy := placementStrategyAuto
	if strategySet {
		strategy = ps.ValueString()
	}

	var preferred *int64
	if preferredSet {
		preferred = ph.ValueInt64Pointer()
	}

	return strategy, preferred, true
}

// placementChanged reports whether the plan's placement differs from state.
func placementChanged(plan, state InstanceModel) bool {
	if plan.ConfigHvm.IsNull() || plan.ConfigHvm.IsUnknown() {
		return false
	}

	strategyChanged := !plan.ConfigHvm.PlacementStrategy.IsUnknown() &&
		!plan.ConfigHvm.PlacementStrategy.Equal(state.ConfigHvm.PlacementStrategy)
	preferredChanged := !plan.ConfigHvm.PreferredHostId.IsUnknown() &&
		!plan.ConfigHvm.PreferredHostId.Equal(state.ConfigHvm.PreferredHostId)

	return strategyChanged || preferredChanged
}

// applyPlacementUpdate issues ManageHostPlacement for each node server when the
// plan's placement differs from state. When the resulting strategy is "auto"
// the preferred host is omitted, so the appliance clears any host preference.
func applyPlacementUpdate(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	plan, state InstanceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if !placementChanged(plan, state) {
		return diags
	}

	// Resolve the effective strategy: the planned one when set, else "auto".
	strategy := placementStrategyAuto
	ps := plan.ConfigHvm.PlacementStrategy
	if !ps.IsNull() && !ps.IsUnknown() {
		strategy = ps.ValueString()
	}

	var preferred *int64
	ph := plan.ConfigHvm.PreferredHostId
	if strategy != placementStrategyAuto && !ph.IsNull() && !ph.IsUnknown() {
		preferred = ph.ValueInt64Pointer()
	}

	return applyHostPlacement(ctx, client, instanceID, strategy, preferred)
}

// hostPlacement is the placement read from an instance's node servers. Read is
// false when the servers could not be read; the caller then keeps its prior
// values (the read is auxiliary, and must not fail or empty the instance).
type hostPlacement struct {
	Strategy  types.String
	Preferred types.Int64
	Read      bool
}

// refreshHostPlacement reads placement_strategy and preferred_host_id from the
// node compute servers of an already-fetched instance. Placement is written to
// every node, so every node is read: state takes the first node's values and a
// warning names any node that disagrees, so drift on one node of a multi-node
// instance is visible. Returns null values (Read=true) when the instance has
// no node server. A failed server GET returns Read=false with a warning, not an
// error: a 403/404 on it must not fail a refresh of an instance that does not
// use placement.
func refreshHostPlacement(
	ctx context.Context,
	client *sdk.APIClient,
	instance sdk.GetInstance200ResponseInstance,
) (hostPlacement, diag.Diagnostics) {
	var diags diag.Diagnostics

	strategy, preferred := types.StringNull(), types.Int64Null()
	first := true

	for _, cont := range instance.ContainerDetails {
		if cont.Server == nil || cont.Server.Id == nil {
			continue
		}

		sid := *cont.Server.Id

		resp, hresp, err := client.HostsAPI.
			GetHost(ctx, sdk.GetHostIdParameter{Int64: &sid}).
			Execute()
		if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK || resp == nil {
			diags.AddWarning(
				"host placement not read",
				fmt.Sprintf("server %d: GET failed: %s; placement_strategy and preferred_host_id "+
					"keep their previous values", sid, errfmt.ErrMsg(err, hresp)),
			)

			return hostPlacement{Strategy: types.StringNull(), Preferred: types.Int64Null()}, diags
		}

		s, p := mapServerPlacement(resp.Server)

		if first {
			strategy, preferred = s, p
			first = false

			continue
		}

		if !s.Equal(strategy) || !p.Equal(preferred) {
			diags.AddWarning(
				"host placement differs between nodes",
				fmt.Sprintf("server %d has placement_strategy %s / preferred_host_id %s, but the "+
					"instance's first node has %s / %s; state follows the first node. Apply to "+
					"re-send the configured placement to every node.",
					sid, s.String(), p.String(), strategy.String(), preferred.String()),
			)
		}
	}

	return hostPlacement{Strategy: strategy, Preferred: preferred, Read: true}, diags
}
