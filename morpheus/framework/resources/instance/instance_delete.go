// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// DeleteErrorStatuses are the instance statuses that mean a removal has failed
// and will not complete on its own.
//
// "stopped" and "suspended" are deliberately absent. Morpheus's container-to-
// instance status sync writes them onto an instance that is "removing" whenever
// all of its containers are in that state, so both are transient during a
// normal teardown; treating them as failures is the race MORPH-4733 reported.
// "warning" is present because it is the status Morpheus sets when the removal
// job itself fails, with the reason in statusMessage.
var DeleteErrorStatuses = []string{
	"denied",
	"cancelled",
	"failed",
	"warning",
	"restoring",
}

var StopErrorStatuses = []string{
	"denied",
	"cancelled",
	"failed",
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

	id := data.Id.ValueInt64()
	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	// Get instance
	ret, hresp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
	if err != nil {
		if hresp == nil || hresp.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: GET failed ", id)+errfmt.ErrMsg(err, hresp),
			)
		}

		return
	}

	instance := ret.Instance
	if instance == nil {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: GET returned empty instance", id),
		)

		return
	}

	// Get serverId(s).  At the moment we only support a single server per instance, but we'll loop just in case
	containers := instance.ContainerDetails
	if containers == nil {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: GET returned empty containers", id),
		)

		return
	}

	serverIds := make([]int64, 0)
	for _, container := range containers {
		server := container.Server
		if server == nil {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: GET returned empty server in container", id),
			)

			return
		}

		serverId := server.Id
		if serverId == nil {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: GET returned empty server ID in container", id),
			)

			return
		}

		serverIds = append(serverIds, *serverId)
	}

	// Stop the server(s) if they are not already stopped, otherwise we cannot delete them
	for _, serverId := range serverIds {
		stopId := sdk.UpdateHostIdParameter{
			Int64: &serverId,
		}
		stopReq := client.HostsAPI.StopHost(ctx, stopId)
		_, hresp, err := stopReq.Execute()
		if err != nil || (hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusConflict) {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: failed to stop server %d before delete ", id, serverId)+errfmt.ErrMsg(err, hresp),
			)

			return
		}

		getId := sdk.GetHostIdParameter{
			Int64: &serverId,
		}
		waitForStopped := func() (*sdk.GetHost200Response, error) {
			resp, hresp, err := client.HostsAPI.GetHost(ctx, getId).Execute()
			if err != nil {
				if hresp == nil || hresp.StatusCode != http.StatusNotFound {
					return nil, backoff.Permanent(err)
				}
			}

			server := resp.Server
			if server == nil {
				return nil, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty server", id, serverId))
			}

			if server.Status == nil {
				return nil, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty status", id, serverId))
			}

			return resp, checkStatusDone(
				*server.Status,
				[]string{"provisioned", "stopped"},
				StopErrorStatuses,
			)
		}

		if _, err := backoff.Retry(
			ctx,
			waitForStopped,
			backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
			backoff.WithMaxElapsedTime(deleteTimeout),
		); err != nil {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: failed to stop server %d before delete ", id, serverId)+err.Error(),
			)

			return
		}
	}

	// Delete the server(s) associated with the instance.
	for _, serverId := range serverIds {
		updateId := sdk.UpdateHostIdParameter{
			Int64: &serverId,
		}
		deleteServerReq := client.HostsAPI.RemoveHost(ctx, updateId).Force("on").
			RemoveResources("on").RemoveInstances("on")
		_, hresp, err := deleteServerReq.Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: DELETE server %d failed ", id, serverId)+errfmt.ErrMsg(err, hresp),
			)

			return
		}

		getId := sdk.GetHostIdParameter{
			Int64: &serverId,
		}
		waitForServerDeleted := func() (*sdk.GetHost200Response, error) {
			resp, hresp, err := client.HostsAPI.GetHost(ctx, getId).Execute()
			if err != nil {
				if hresp == nil || hresp.StatusCode != http.StatusNotFound {
					return nil, backoff.Permanent(err)
				}
			}

			// 404 status code counts as a successful delete
			if hresp.StatusCode == http.StatusNotFound {
				return nil, nil
			}

			server := resp.Server
			if server == nil {
				return nil, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty server", id, serverId))
			}

			if server.Status == nil {
				return nil, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty status", id, serverId))
			}

			return resp, checkStatusDone(
				*server.Status,
				nil,
				StopErrorStatuses,
			)
		}

		if _, err := backoff.Retry(
			ctx,
			waitForServerDeleted,
			backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
			backoff.WithMaxElapsedTime(deleteTimeout),
		); err != nil {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: DELETE server %d failed ", id, serverId)+err.Error(),
			)

			return
		}
	}

	// If we get here, all servers have been deleted successfully.  Now we can delete the instance itself.
	deleteReq := client.InstancesAPI.DeleteInstance(ctx, id).Force("on").
		RemoveVolumes("on").ReleaseEIPs("on").ReleaseFloatingIps("on")
	_, hresp, err = deleteReq.Execute()
	if err != nil {
		if hresp == nil || (hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusNotFound) {
			resp.Diagnostics.AddError(
				"delete instance resource",
				fmt.Sprintf("instance %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
			)

			return
		}
	}

	waitForDeleted := func() (*sdk.GetInstance200Response, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusNotFound {
				return nil, backoff.Permanent(err)
			}
		}

		// 404 status code counts as a successful delete
		if hresp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		instance := resp.Instance
		if instance == nil {
			return nil, backoff.Permanent(fmt.Errorf("instance %d: GET returned empty instance", id))
		}

		if instance.Status == nil {
			return nil, backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", id))
		}

		return resp, checkDeleteStatusDone(*instance.Status, instance.StatusMessage.Get())
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

// checkDeleteStatusDone classifies the status of an instance still present after
// its DELETE. Anything outside DeleteErrorStatuses keeps polling — the instance
// is expected to disappear, and the caller treats a 404 as success before
// reaching here. An error status is terminal, and the error carries Morpheus's
// statusMessage when there is one: for a failed removal that is the only place
// the reason is recorded ("Unable to remove instance: ...").
func checkDeleteStatusDone(status string, statusMessage *string) error {
	if !slices.Contains(DeleteErrorStatuses, status) {
		return backoff.RetryAfter(5)
	}

	reason := "reached error status: " + status
	if statusMessage != nil && strings.TrimSpace(*statusMessage) != "" {
		reason += " (" + strings.TrimSpace(*statusMessage) + ")"
	}

	return backoff.Permanent(errors.New(reason))
}
