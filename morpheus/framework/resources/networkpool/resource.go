package networkpool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &networkPoolResource{}
	_ resource.ResourceWithConfigure   = &networkPoolResource{}
	_ resource.ResourceWithImportState = &networkPoolResource{}
)

type networkPoolResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkPoolResource{}
}

func (r *networkPoolResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_pool"
}

func (r *networkPoolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkPoolResourceSchema(ctx)
}

func (r *networkPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolReq := sdk.CreateNetworkPoolRequestNetworkPool{
		Name: plan.Name.ValueStringPointer(),
		Type: &sdk.CreateNetworkPoolRequestNetworkPoolType{
			AdditionalProperties: map[string]interface{}{
				"code": plan.TypeCode.ValueString(),
			},
		},
	}

	additionalProps := map[string]interface{}{}
	if !plan.PoolEnabled.IsNull() && !plan.PoolEnabled.IsUnknown() {
		additionalProps["poolEnabled"] = plan.PoolEnabled.ValueBool()
	}
	if !plan.DnsDomain.IsNull() && !plan.DnsDomain.IsUnknown() {
		additionalProps["dnsDomain"] = plan.DnsDomain.ValueString()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		additionalProps["dhcpServer"] = plan.DhcpServer.ValueBool()
	}
	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		additionalProps["gateway"] = plan.Gateway.ValueString()
	}
	if !plan.Netmask.IsNull() && !plan.Netmask.IsUnknown() {
		additionalProps["netmask"] = plan.Netmask.ValueString()
	}
	if !plan.SubnetAddress.IsNull() && !plan.SubnetAddress.IsUnknown() {
		additionalProps["subnetAddress"] = plan.SubnetAddress.ValueString()
	}
	// Ranges: every element of ip_address_ranges, or the deprecated single
	// ip_ranges object. Exactly one is configured (schema validators).
	planned, diags := rangesFromList(ctx, plan.IpAddressRanges)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	if single, ok := singleRangeFromObject(plan.IpRanges); ok {
		planned = []ipRange{single}
	}
	poolReq.IpRanges = rangesToCreateRequest(planned)
	if len(additionalProps) > 0 {
		poolReq.AdditionalProperties = additionalProps
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkPool(ctx).
		CreateNetworkPoolRequest(sdk.CreateNetworkPoolRequest{
			NetworkPool: &poolReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.NetworkPool == nil || result.NetworkPool.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "NetworkPool ID is nil in the create response")

		return
	}

	id := *result.NetworkPool.Id

	readResult, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_pool",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.NetworkPool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	resp.Diagnostics.Append(mapReadResponseToModel(&plan, readResult.NetworkPool)...)

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool", "", err, httpResp)

		return
	}

	pool := result.NetworkPool
	if pool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	resp.Diagnostics.Append(mapReadResponseToModel(&state, pool)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan, state NetworkPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	poolReq := sdk.UpdateNetworkPoolRequestNetworkPool{
		Name: plan.Name.ValueStringPointer(),
	}

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	additionalProps := map[string]interface{}{}
	if !plan.PoolEnabled.IsNull() && !plan.PoolEnabled.IsUnknown() {
		additionalProps["poolEnabled"] = plan.PoolEnabled.ValueBool()
	}
	switch {
	case pending.Clears("dns_domain"):
		additionalProps["dnsDomain"] = nil
	case !plan.DnsDomain.IsNull() && !plan.DnsDomain.IsUnknown():
		additionalProps["dnsDomain"] = plan.DnsDomain.ValueString()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		additionalProps["dhcpServer"] = plan.DhcpServer.ValueBool()
	}
	switch {
	case pending.Clears("gateway"):
		additionalProps["gateway"] = nil
	case !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown():
		additionalProps["gateway"] = plan.Gateway.ValueString()
	}
	switch {
	case pending.Clears("netmask"):
		additionalProps["netmask"] = nil
	case !plan.Netmask.IsNull() && !plan.Netmask.IsUnknown():
		additionalProps["netmask"] = plan.Netmask.ValueString()
	}
	switch {
	case pending.Clears("subnet_address"):
		additionalProps["subnetAddress"] = nil
	case !plan.SubnetAddress.IsNull() && !plan.SubnetAddress.IsUnknown():
		additionalProps["subnetAddress"] = plan.SubnetAddress.ValueString()
	}
	// Ranges. Morpheus reconciles by id (see ranges.go): send each planned
	// range with the id carried in state so it is edited in place; a range
	// missing from the plan is deleted by the appliance.
	priorRanges, diags := rangesFromList(ctx, state.IpAddressRanges)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	planned, diags := rangesFromList(ctx, plan.IpAddressRanges)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	// What the appliance holds now: the managed ranges, and any range
	// ip_address_ranges cannot express (CIDR / IPv6), which must be sent back
	// unchanged or Morpheus deletes it.
	currentPool := getPool(ctx, client, id, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	current, unmanaged := rangesFromAPI(currentPool)

	// Which attribute is in use is decided by the CONFIGURATION, as in
	// carryRangeIDs: ip_ranges is computed and stays mirrored in state (and
	// so in the plan) after a move to ip_address_ranges.
	var config NetworkPoolModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}
	if single, ok := singleRangeFromObject(config.IpRanges); ok {
		// The deprecated single-range attribute would delete every other range
		// of a multi-range pool (or add one beside a CIDR range it cannot
		// express). Refuse rather than lose data.
		if total, ok := singleRangeUpdatable(current, unmanaged); !ok {
			resp.Diagnostics.AddError(
				"ip_ranges cannot manage a pool with several ranges",
				fmt.Sprintf("network pool %d has %d range(s) the deprecated ip_ranges attribute cannot manage as a single "+
					"start/end range; applying it would delete or duplicate them. Switch to ip_address_ranges, "+
					"listing every range.", id, total),
			)

			return
		}
		// The pool's one range is edited in place under its current id, taken
		// from the appliance rather than from prior state: state may not hold
		// ip_address_ranges yet (first apply after upgrading, -refresh=false).
		planned = []ipRange{single}
		if len(current) == 1 {
			planned[0].id = current[0].id
		}
	} else {
		planned = withIDsFromPrior(planned, priorRanges)
	}
	poolReq.IpRanges = append(rangesToUpdateRequest(planned), unmanagedToUpdateRequest(unmanaged)...)
	if len(additionalProps) > 0 {
		poolReq.AdditionalProperties = additionalProps
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkPool(ctx, id).
		UpdateNetworkPoolRequest(sdk.UpdateNetworkPoolRequest{
			NetworkPool: &poolReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusBadRequest && isReservationConflict(err.Error()) {
			resp.Diagnostics.AddError(
				"a removed IP range still has reserved addresses",
				"Morpheus refused to delete a range of network pool "+plan.Name.ValueString()+" because addresses in it are "+
					"still reserved (for example by hpe_morpheus_network_pool_ip resources or instance allocations). "+
					"Release them before removing the range from ip_address_ranges.",
			)

			return
		}
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_pool", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Re-read to get current state
	readResult, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool", "", err, httpResp)

		return
	}

	pool := readResult.NetworkPool
	if pool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	resp.Diagnostics.Append(mapReadResponseToModel(&plan, pool)...)

	// Morpheus keeps (and answers 200 for) a range it will not delete because
	// addresses in it are reserved. Record what the appliance holds, then
	// report it, so the plan keeps showing the removal until the
	// reservations are released.
	if apiRanges, _ := rangesFromAPI(pool); len(planned) > 0 {
		if kept := rangesNotRemoved(planned, apiRanges); len(kept) > 0 {
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			// The other attributes of this apply did go through: record them.
			resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
			for _, k := range kept {
				resp.Diagnostics.AddError(
					"a removed IP range still has reserved addresses",
					fmt.Sprintf("Morpheus kept range %s-%s of network pool %s because addresses in it are still reserved "+
						"(for example by hpe_morpheus_network_pool_ip resources or instance allocations). "+
						"Release them before removing the range from ip_address_ranges.", k.start, k.end, plan.Name.ValueString()),
				)
			}

			return
		}
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.DeleteNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_pool", "", err, httpResp)

		return
	}
}

func (r *networkPoolResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapReadResponseToModel(model *NetworkPoolModel, pool *sdk.GetNetworkPool200ResponseNetworkPool) diag.Diagnostics {
	var diags diag.Diagnostics

	if pool.Id != nil {
		model.Id = types.Int64Value(*pool.Id)
	}
	if pool.Name != nil {
		model.Name = types.StringValue(*pool.Name)
	}
	if t := pool.Type; t != nil && t.Code != nil {
		model.TypeCode = types.StringValue(*t.Code)
	}
	if pool.IpCount != nil {
		model.IpCount = types.Int64Value(*pool.IpCount)
	}
	if pool.FreeCount != nil {
		model.FreeCount = types.Int64Value(*pool.FreeCount)
	}
	if pool.PoolEnabled != nil {
		model.PoolEnabled = types.BoolValue(*pool.PoolEnabled)
	}
	if pool.DnsDomain.IsSet() && pool.DnsDomain.Get() != nil {
		model.DnsDomain = types.StringValue(*pool.DnsDomain.Get())
	} else {
		model.DnsDomain = types.StringNull()
	}
	if pool.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*pool.DhcpServer)
	}
	if pool.Gateway.IsSet() && pool.Gateway.Get() != nil {
		model.Gateway = types.StringValue(*pool.Gateway.Get())
	} else {
		model.Gateway = types.StringNull()
	}
	if pool.Netmask.IsSet() && pool.Netmask.Get() != nil {
		model.Netmask = types.StringValue(*pool.Netmask.Get())
	} else {
		model.Netmask = types.StringNull()
	}
	if pool.SubnetAddress.IsSet() && pool.SubnetAddress.Get() != nil {
		model.SubnetAddress = types.StringValue(*pool.SubnetAddress.Get())
	} else {
		model.SubnetAddress = types.StringNull()
	}
	// Ranges. ip_address_ranges always reflects the appliance (in the prior
	// state's order), so a configuration moving from ip_ranges to
	// ip_address_ranges with the same values plans no change. The deprecated
	// ip_ranges is only populated when the prior state used it, so the two
	// never appear set together.
	ctx := context.Background()
	apiRanges, unmanaged := rangesFromAPI(pool)
	for _, u := range unmanaged {
		diags.AddWarning("IP range not represented",
			fmt.Sprintf("network pool %s has %s without start and end addresses (a CIDR range); "+
				"ip_address_ranges cannot represent it. It is left as it is on every update.", model.Name.ValueString(), u))
	}
	prior, _ := rangesFromList(ctx, model.IpAddressRanges)
	ordered := orderLikePrior(apiRanges, prior)
	if ordered == nil {
		ordered = []ipRange{}
	}
	list, d := rangesToList(ctx, ordered)
	diags.Append(d...)
	model.IpAddressRanges = list

	if _, usedSingle := singleRangeFromObject(model.IpRanges); usedSingle && len(ordered) > 0 {
		model.IpRanges = singleRangeToObject(&ordered[0])
	} else {
		model.IpRanges = singleRangeToObject(nil)
	}

	return diags
}

// getPool fetches the pool. On failure it adds a diagnostic and returns nil;
// callers must check diags (or nil) before using the result.
func getPool(
	ctx context.Context, client *sdk.APIClient, id int64, diags *diag.Diagnostics,
) *sdk.GetNetworkPool200ResponseNetworkPool {
	result, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(diags, errfmt.OpRead, "network_pool", "", err, httpResp)

		return nil
	}
	if result.NetworkPool == nil {
		diags.AddError("API returned nil", "NetworkPool is nil in the response")

		return nil
	}

	return result.NetworkPool
}
