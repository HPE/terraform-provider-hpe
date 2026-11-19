// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Update implements resource.Resource. Only the in-place-updatable fields are
// sent; everything else is RequiresReplace in the schema.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan HostModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state HostModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config HostModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	id := plan.Id.ValueInt64()

	server := sdk.UpdateHostRequestServer{}

	// Attributes removed from the configuration, to be cleared: sent when
	// known, sent empty when being cleared, not sent otherwise. The API binds
	// a cleared string as null and an empty labels list as no labels.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		server.Name = plan.Name.ValueStringPointer()
	}

	server.Description = pending.String("description", plan.Description)

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		server.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		server.Visibility = plan.Visibility.ValueStringPointer()
	}

	labels, sendLabels, ldiags := pending.Strings(ctx, "labels", plan.Labels)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case sendLabels && len(labels) == 0:
		// The SDK omits an empty slice; the handler clears the labels when the
		// key is present with an empty list, so it is sent explicitly.
		if server.AdditionalProperties == nil {
			server.AdditionalProperties = map[string]interface{}{}
		}

		server.AdditionalProperties["labels"] = []string{}
	case sendLabels:
		server.Labels = labels
	}

	server.SshUsername = pending.String("ssh_username", plan.SshUsername)

	switch {
	case pending.Clears("ssh_key_pair_id"):
		// A sshKeyPair object without an id unassigns the key pair.
		server.SshKeyPair = &sdk.UpdateHostRequestServerSshKeyPair{}
	case !plan.SshKeyPairId.IsNull() && !plan.SshKeyPairId.IsUnknown():
		server.SshKeyPair = &sdk.UpdateHostRequestServerSshKeyPair{
			Id: plan.SshKeyPairId.ValueInt64Pointer(),
		}
	}

	// ssh_password_wo: re-sent when ModifyPlan marked write_only_revision
	// unknown (its digest changed, or the appliance's stored hash did). See
	// write_only.go.
	resend := wodigest.Resend(ctx, req.Plan)
	if resend {
		sshPwd, _, _ := writeOnlyFromConfig(config)
		if sshPwd != nil {
			server.SshPassword.Set(sshPwd)
		}
	}

	// user_data -> config.userData: a cleared value is sent as null.
	serverConfig := &sdk.UpdateHostRequestServerConfig{}
	if clearable.NullableString(pending, "user_data", plan.UserData, &serverConfig.UserData) {
		server.Config = serverConfig
	}

	updateReq := sdk.UpdateHostRequest{Server: &server}

	_, hresp, err := client.HostsAPI.
		UpdateHost(ctx, sdk.UpdateHostIdParameter{Int64: &id}).
		UpdateHostRequest(updateReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update host",
			fmt.Sprintf("host %d: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	getResp, hresp, err := client.HostsAPI.
		GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).
		Execute()
	if err != nil || getResp == nil || getResp.Server == nil {
		resp.Diagnostics.AddError(
			"failed to read host after update",
			fmt.Sprintf("host %d: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	bareMetalCodes := bareMetalCodesOrWarn(ctx, client, &resp.Diagnostics, id)

	getHostAsState(ctx, getResp.Server, bareMetalCodes, &plan)

	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if resend {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}

	// Record the digests of the configured write-only values and re-baseline
	// the appliance-side hash from the read-back server.
	resp.Diagnostics.Append(wo.StoreObserved(ctx, resp.Private, req.Config, getResp.Server)...)

	// Which clearable attributes the configuration set, for the next plan.
	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
