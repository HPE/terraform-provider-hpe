// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkDhcpServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network dhcp server resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	serverId := float32(state.NetworkServerId.ValueInt64())

	tflog.Debug(ctx, fmt.Sprintf("Deleting network dhcp server %d", id))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkDhcpServer(ctx, id, serverId).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network dhcp server resource",
			fmt.Sprintf("network dhcp server %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}
