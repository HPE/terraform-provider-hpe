// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var DeleteErrorStatuses = []string{
	"denied",
	"cancelled",
	"failed",
	"stopped",
	"suspended",
	"restoring",
}

// Delete implements resource.Resource.
func (g *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data InstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	deleteTimeout, diags := data.Timeouts.Delete(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	id := data.Id
	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	deleteReq := client.InstancesAPI.DeleteInstance(ctx, id.ValueInt64()).Force("on").
		RemoveVolumes("on").ReleaseEIPs("on").ReleaseFloatingIps("on")
	_, hresp, err := deleteReq.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	waitForDeleted := func() (*sdk.GetInstance200Response, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, data.Id.ValueInt64()).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusNotFound {
				return nil, backoff.Permanent(err)
			}
		}

		// 404 status code counts as a successful delete
		if hresp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		// Get instance
		instance, ok := resp.GetInstanceOk()
		if !ok || instance == nil {
			return nil, backoff.Permanent(fmt.Errorf("instance %d: GET returned empty instance", id))
		}

		// Get status
		status, ok := instance.GetStatusOk()
		if !ok || status == nil {
			return nil, backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", id))
		}

		return resp, checkStatusDone(
			*status,
			nil,
			DeleteErrorStatuses,
		)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(deleteTimeout),
	); err != nil {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: DELETE failed ", id)+err.Error(),
		)
	}
}
