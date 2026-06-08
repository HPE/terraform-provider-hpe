// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_router_nat

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
			"network_router_nat",
		},
		"_",
	)
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkRouterNatResourceSchema(ctx)
}

// Create

const createOperation = "create network router nat resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterNatModel

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

	nat := sdk.CreateNetworkRouterNatRequestNetworkRouterNAT{}
	nat.Name = plan.Name.ValueString()

	createReq := sdk.CreateNetworkRouterNatRequest{
		NetworkRouterNAT: &nat,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterNat(ctx, routerID).
		CreateNetworkRouterNatRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d NAT POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if !result.Id.IsSet() || result.Id.Get() == nil {
		resp.Diagnostics.AddError("API returned nil", "ID is nil in the response")

		return
	}

	id := *result.Id.Get()
	plan.Id = types.Int64Value(id)

	state, pdiags := getNatAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read

const readOperation = "read network router nat resource"

func getNatAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterNatModel,
) (NetworkRouterNatModel, diag.Diagnostics) {
	var state NetworkRouterNatModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterNat(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("NAT %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, diags
	}

	nat := resp.NetworkRouterNAT
	if nat == nil {
		diags.AddError("API returned nil", "NetworkRouterNAT is nil in the response")

		return state, diags
	}

	if nat.Id != nil {
		state.Id = types.Int64Value(int64(*nat.Id))
	}

	state.RouterId = plan.RouterId

	if nat.Name != nil {
		state.Name = types.StringValue(*nat.Name)
	}

	if nat.Description != nil {
		state.Description = types.StringValue(*nat.Description)
	} else {
		state.Description = types.StringNull()
	}

	if nat.Enabled != nil {
		state.Enabled = types.BoolValue(*nat.Enabled)
	}

	if nat.SourceNetwork != nil {
		state.SourceNetwork = types.StringValue(*nat.SourceNetwork)
	} else {
		state.SourceNetwork = types.StringNull()
	}

	if nat.DestinationNetwork.IsSet() {
		state.DestinationNetwork = types.StringValue(*nat.DestinationNetwork.Get())
	} else {
		state.DestinationNetwork = types.StringNull()
	}

	if nat.TranslatedNetwork != nil {
		state.TranslatedNetwork = types.StringValue(*nat.TranslatedNetwork)
	} else {
		state.TranslatedNetwork = types.StringNull()
	}

	if nat.Priority != nil {
		state.Priority = types.Int64Value(int64(*nat.Priority))
	} else {
		state.Priority = types.Int64Null()
	}

	if nat.Protocol.IsSet() {
		state.Protocol = types.StringValue(*nat.Protocol.Get())
	} else {
		state.Protocol = types.StringNull()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterNatModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, diags := getNatAsState(ctx, plan.Id.ValueInt64(), plan.RouterId.ValueInt64(), client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkRouterNatModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("update network router nat resource", "failed to create client: "+err.Error())

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	nat := sdk.UpdateNetworkRouterNatRequestNetworkRouterNAT{}
	nat.Name = plan.Name.ValueString()

	updateReq := sdk.UpdateNetworkRouterNatRequest{
		NetworkRouterNAT: &nat,
	}

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterNat(ctx, id, routerID).
		UpdateNetworkRouterNatRequest(updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network router nat resource",
			fmt.Sprintf("NAT %d PUT failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, diags := getNatAsState(ctx, id, routerID, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkRouterNatModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete network router nat resource", "failed to create client: "+err.Error())

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting NAT %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterNat(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("NAT %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router nat resource",
			fmt.Sprintf("NAT %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp)),
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
			"import network router nat resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.nat_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router nat resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router nat resource",
			"provided nat_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("router_id"), routerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
