// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_router_route

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = strings.Join(
		[]string{
			req.ProviderTypeName,
			constants.SubProviderName,
			"network_router_route",
		},
		"_",
	)
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkRouterRouteResourceSchema(ctx)
}

// Create

const createOperation = "create network router route resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterRouteModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(createOperation, "failed to create client: "+err.Error())

		return
	}

	routerID := plan.RouterId.ValueInt64()

	route := sdk.CreateNetworkRouterRouteRequestNetworkRoute{}
	route.Source = plan.Source.ValueString()
	route.Destination = plan.Destination.ValueString()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		route.Name = &name
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		route.Description = &description
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := plan.Enabled.ValueBool()
		route.Enabled = &enabled
	}

	if !plan.DefaultRoute.IsNull() && !plan.DefaultRoute.IsUnknown() {
		defaultRoute := plan.DefaultRoute.ValueBool()
		route.DefaultRoute = &defaultRoute
	}

	if !plan.NetworkMtu.IsNull() && !plan.NetworkMtu.IsUnknown() {
		networkMtu := float32(plan.NetworkMtu.ValueFloat64())
		route.NetworkMtu = &networkMtu
	}

	createReq := sdk.CreateNetworkRouterRouteRequest{
		NetworkRoute: &route,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterRoute(ctx, routerID).
		CreateNetworkRouterRouteRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d route POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	id := *result.Id.Get()
	plan.Id = types.Int64Value(id)

	state, pdiags := getRouteAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read

const readOperation = "read network router route resource"

func getRouteAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterRouteModel,
) (NetworkRouterRouteModel, diag.Diagnostics) {
	var state NetworkRouterRouteModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterRoute(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("route %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, diags
	}

	route := resp.NetworkRoute

	if route != nil && route.Id != nil {
		state.Id = types.Int64Value(*route.Id)
	}

	state.RouterId = plan.RouterId

	if route != nil && route.Name != nil {
		state.Name = types.StringValue(*route.Name)
	} else {
		state.Name = types.StringNull()
	}

	if route != nil && route.Description.IsSet() {
		state.Description = types.StringValue(*route.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if route != nil && route.Source != nil {
		state.Source = types.StringValue(*route.Source)
	}

	if route != nil && route.Destination != nil {
		state.Destination = types.StringValue(*route.Destination)
	}

	if route != nil && route.DefaultRoute != nil {
		state.DefaultRoute = types.BoolValue(*route.DefaultRoute)
	}

	if route != nil && route.Enabled != nil {
		state.Enabled = types.BoolValue(*route.Enabled)
	}

	if route != nil && route.NetworkMtu.IsSet() {
		state.NetworkMtu = types.Float64Value(float64(*route.NetworkMtu.Get()))
	} else {
		state.NetworkMtu = types.Float64Null()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterRouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, diags := getRouteAsState(ctx, plan.Id.ValueInt64(), plan.RouterId.ValueInt64(), client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update — not supported, all mutable fields use RequiresReplace

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// All attributes use RequiresReplace; this should never be called.
	resp.Diagnostics.AddError(
		"update network router route resource",
		"update is not supported — all changes require replacement",
	)
}

// Delete

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkRouterRouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete network router route resource", "failed to create client: "+err.Error())

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting route %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterRoute(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Route %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router route resource",
			fmt.Sprintf("route %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

// Import

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.route_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided route_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("router_id"), routerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
