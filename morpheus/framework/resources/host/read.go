// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/hosttypes"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

// Read implements resource.Resource.
func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state HostModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	id := state.Id.ValueInt64()

	getResp, hresp, err := client.HostsAPI.
		GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).
		Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, fmt.Sprintf("host %d not found, removing from state", id))
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"failed to read host",
			fmt.Sprintf("host %d: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if getResp == nil || getResp.Server == nil {
		resp.Diagnostics.AddError(
			"failed to read host",
			fmt.Sprintf("host %d: response or server is nil", id),
		)

		return
	}

	bareMetalCodes := bareMetalCodesOrWarn(ctx, client, &resp.Diagnostics, id)

	getHostAsState(ctx, getResp.Server, bareMetalCodes, &state)

	// Compare the appliance's hash of the stored SSH password with the baseline
	// recorded when the provider last sent it; a difference means it was
	// changed out of band and the next plan re-sends the configured value.
	drifted, driftDiags := wo.CheckDrift(ctx, resp.Private, getResp.Server)
	resp.Diagnostics.Append(driftDiags...)
	if len(drifted) > 0 {
		tflog.Info(ctx, "write-only value changed on the appliance", map[string]any{"attributes": drifted})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported host does not plan a change nobody made.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, HostResourceSchema(ctx), &resp.State)...,
	)
}

// getHostAsState maps the GetHost response server into the Terraform state model.
// Create-only inputs (config variants, passwords, delete_options,
// network_interfaces, volumes) are left as-is: they are not returned by the API.
func getHostAsState(
	_ context.Context,
	server *sdk.GetHost200ResponseServer,
	bareMetalCodes map[string]bool,
	state *HostModel,
) {
	state.Id = convert.Int64ToType(server.Id)
	state.Name = convert.StrToType(server.Name)
	state.Uuid = convert.StrToType(server.Uuid)
	state.ExternalId = convert.StrToType(server.ExternalId.Get())
	state.Status = convert.StrToType(server.Status)
	state.PowerState = convert.StrToType(server.PowerState)
	state.InternalIp = convert.StrToType(server.InternalIp.Get())
	state.ExternalIp = convert.StrToType(server.ExternalIp.Get())
	state.Hostname = convert.StrToType(server.Hostname)
	state.AgentInstalled = convert.BoolToType(server.AgentInstalled)
	state.AgentVersion = convert.StrToType(server.AgentVersion.Get())
	state.Managed = convert.BoolToType(server.Managed)
	state.Description = convert.StrToType(server.Description.Get())
	state.Enabled = convert.BoolToType(server.Enabled)
	state.Visibility = convert.StrToType(server.Visibility)

	state.DateCreated = convert.TimeToType(server.DateCreated)
	state.LastUpdated = convert.TimeToType(server.LastUpdated)

	// cloud / group.
	if server.ZoneId != nil {
		state.CloudId = convert.Int64ToType(server.ZoneId)
	}
	if server.SiteId != nil {
		state.GroupId = convert.Int64ToType(server.SiteId)
	}

	// plan (null, not unknown, when the response omits it).
	state.PlanId = convert.Int64ToType(nil)
	if server.Plan != nil {
		state.PlanId = convert.Int64ToType(server.Plan.Id.Get())
	}

	// compute server type.
	state.ComputeServerTypeId = convert.Int64ToType(nil)
	if server.ComputeServerType != nil {
		state.ComputeServerTypeId = convert.Int64ToType(server.ComputeServerType.Id)
		if server.ComputeServerType.Code != nil {
			state.ComputeServerType = convert.StrToType(server.ComputeServerType.Code)
		}
		// The server's rendering omits bareMetalHost; the server types list
		// (bareMetalCodes) says which types are bare metal.
		state.BareMetal = hosttypes.BareMetalValue(bareMetalCodes, server.ComputeServerType.Code)
	} else {
		state.BareMetal = hosttypes.BareMetalValue(bareMetalCodes, nil)
	}

	// parent host.
	if server.ParentServer != nil && server.ParentServer.Id != nil {
		state.ParentHostId = convert.Int64ToType(server.ParentServer.Id)
	} else {
		state.ParentHostId = convert.Int64ToType(nil)
	}

	// labels. A configured empty set is kept as an empty set (the API returns
	// none, which StrSliceToSet would map to null -- an inconsistent result).
	switch {
	case len(server.Labels) > 0:
		state.Labels = convert.StrSliceToSet(server.Labels)
	case !state.Labels.IsNull() && !state.Labels.IsUnknown() && len(state.Labels.Elements()) == 0:
		// keep the configured empty set
	default:
		state.Labels = convert.StrSliceToSet(nil)
	}

	// ssh access: the key pair and the config's userData are read back, so
	// drift is visible and an import recovers them. ssh_username is the
	// credential given for the add; Morpheus replaces it with the account its
	// own provisioning creates on the node (cloudinit, say), so a configured
	// value is kept as the managed one and the API's is read only when the
	// attribute is not set (an import, or never configured).
	if state.SshUsername.IsNull() || state.SshUsername.IsUnknown() {
		state.SshUsername = convert.StrToType(server.SshUsername.Get())
	}

	if server.SshKeyPair != nil && server.SshKeyPair.Id != nil {
		state.SshKeyPairId = convert.Int64ToType(server.SshKeyPair.Id)
	} else {
		state.SshKeyPairId = convert.Int64ToType(nil)
	}

	state.UserData = convert.StrToType(configUserData(server.Config))
}

// configUserData returns config.userData, or nil when it is absent or empty
// (Morpheus stores a cleared value as null and some types as "").
func configUserData(config *sdk.AddBaremetalHost200ResponseServerConfig) *string {
	if config == nil {
		return nil
	}

	v := config.UserData.Get()
	if v == nil || *v == "" {
		return nil
	}

	return v
}

// bareMetalCodesOrWarn returns the bare-metal server type codes, or nil with a
// warning when the server types list cannot be read: bare_metal is then null
// ("not known") rather than failing the whole operation (and, after a create,
// tainting a host that was provisioned fine).
func bareMetalCodesOrWarn(ctx context.Context, client *sdk.APIClient, diags *diag.Diagnostics, id int64) map[string]bool {
	codes, err := hosttypes.BareMetalCodes(ctx, client)
	if err != nil {
		diags.AddWarning(
			"bare_metal could not be determined",
			fmt.Sprintf("host %d: the server types list could not be read (%s); bare_metal is null "+
				"until the next refresh.", id, err.Error()),
		)

		return nil
	}

	return codes
}
