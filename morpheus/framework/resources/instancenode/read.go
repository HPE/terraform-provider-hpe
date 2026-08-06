// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/containerip"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Read implements resource.Resource.
func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state instanceNodeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	instanceID := state.InstanceID.ValueInt64()
	containerID := state.ContainerID.ValueInt64()

	err = refreshNodeState(ctx, client, instanceID, containerID, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read instance node",
			fmt.Sprintf("instance %d container %d: %s",
				instanceID, containerID, err.Error()),
		)

		return
	}

	// If the container is gone, remove from state.
	if state.ContainerID.IsNull() {
		tflog.Info(ctx, fmt.Sprintf(
			"container %d not found on instance %d, removing from state",
			containerID, instanceID))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// refreshNodeState reads the instance and populates state fields for the
// given container. If the container is not found, state.ContainerID is set
// to null to signal removal.
func refreshNodeState(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	containerID int64,
	state *instanceNodeModel,
) error {
	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			state.ContainerID = types.Int64Null()

			return nil
		}

		return fmt.Errorf("GetInstance: %s", errfmt.ErrMsg(err, hresp))
	}

	if getResp == nil || getResp.Instance == nil {
		return fmt.Errorf("instance %d: response is nil", instanceID)
	}

	for i := range getResp.Instance.ContainerDetails {
		cd := &getResp.Instance.ContainerDetails[i]
		if cd.Id != nil && *cd.Id == containerID {
			populateNodeMetadata(state, cd)

			return nil
		}
	}

	// Container not found — signal removal.
	state.ContainerID = types.Int64Null()

	return nil
}

// populateNodeMetadata sets state fields from the container detail.
func populateNodeMetadata(state *instanceNodeModel, cd *sdk.InstanceContainer2) {
	state.ContainerID = convert.Int64ToType(cd.Id)

	if cd.Server != nil {
		state.ServerID = convert.Int64ToType(cd.Server.Id)
	} else {
		state.ServerID = types.Int64Null()
	}

	if cd.Ip != nil && containerip.Ready(*cd.Ip) {
		state.IPAddress = types.StringValue(*cd.Ip)
	} else {
		state.IPAddress = types.StringNull()
	}

	state.Hostname = convert.StrToType(cd.Hostname)
	state.InternalIP = convert.StrToType(cd.InternalIp)
	state.ExternalFQDN = convert.StrToType(cd.ExternalFqdn)
	state.MacAddress = resolvePrimaryMacAddress(cd.Server)
}

// resolvePrimaryMacAddress returns the MAC address of the server's primary
// network interface. It prefers the interface flagged as primary; if none is
// flagged, it falls back to the first interface. Returns null if the server
// has no interfaces or is nil.
func resolvePrimaryMacAddress(server *sdk.InstanceContainerServer2) types.String {
	if server == nil || len(server.Interfaces) == 0 {
		return types.StringNull()
	}

	// Prefer the interface flagged as primary.
	for i := range server.Interfaces {
		iface := &server.Interfaces[i]
		if iface.PrimaryInterface != nil && *iface.PrimaryInterface {
			return convert.StrToType(iface.MacAddress)
		}
	}

	// Fallback: first interface.
	return convert.StrToType(server.Interfaces[0].MacAddress)
}
