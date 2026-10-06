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
	"github.com/hashicorp/terraform-plugin-log/tflog"

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

// deleteBusyStatuses are the instance statuses in which Morpheus refuses
// DELETE /api/instances/{id} (InstanceService.canDeleteInstance) unless the
// caller's force flag is honoured — and force is honoured only for users with
// the "Provisioning: Force Delete" permission. A destroy that lands while the
// instance is still finishing a resize is the common case: Delete waits for
// these to clear before asking.
var deleteBusyStatuses = []string{
	"provisioning",
	"restarting",
	"resizing",
	"cloning",
}

// deleteBusyWait bounds how long Delete waits for a busy instance to settle
// before deleting it anyway (Morpheus may still accept it if force is honoured,
// and the server fallback does not check status).
const deleteBusyWait = 10 * time.Minute

// deletePoll is the interval between status polls during a destroy.
const deletePoll = 5 * time.Second

// Delete implements resource.Resource.
//
// The instance is deleted through the instance API — the path the Morpheus UI
// uses, which tears the virtual machine down through the cloud's provisioning
// service and removes the records last. The server API (stop, then
// DELETE /api/servers/{id}?force=on&removeResources=on&removeInstances=on) is
// kept only as a fallback for the cases it was introduced for in #902: an
// instance Morpheus refuses to delete in its current state, or one whose
// removal fails and leaves a server behind. It is no longer the primary path
// because on Morpheus 8.1.x it removes the instance and server records but can
// leave the virtual machine in the hypervisor (MORPH-18513): its
// onDeleteServerInstances waits on the instance removal before queuing the VM
// deletion, and that wait times out.
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

	d := &instanceDeleter{client: client, id: id, timeout: deleteTimeout}

	// Phase 0: read the instance. Its servers are remembered so that, once the
	// instance is gone, Delete can confirm they went with it.
	instance, gone, err := d.getInstance(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete instance resource", err.Error())

		return
	}
	if gone {
		return
	}
	serverIDs, err := serverIDsOf(id, instance)
	if err != nil {
		resp.Diagnostics.AddError("delete instance resource", err.Error())

		return
	}

	// Phase 1: the instance API.
	if err := d.deleteViaInstanceAPI(ctx, instance); err != nil {
		var refused *instanceDeleteRefused
		if !errors.As(err, &refused) {
			resp.Diagnostics.AddError("delete instance resource", err.Error())

			return
		}
		// Morpheus refused, or the removal failed and stopped: fall back to the
		// server path, which is what the resource did unconditionally before.
		tflog.Warn(ctx, "instance delete did not complete; falling back to deleting its servers",
			map[string]any{"instance": id, "reason": refused.Error()})
		resp.Diagnostics.AddWarning(
			"instance delete fell back to the server API",
			fmt.Sprintf("instance %d: %s. Its server(s) are being force-deleted instead; "+
				"on Morpheus 8.1.x check the hypervisor afterwards, as that path can leave "+
				"the virtual machine behind.", id, refused.Error()),
		)
		for _, serverID := range serverIDs {
			if err := d.forceDeleteServer(ctx, serverID); err != nil {
				resp.Diagnostics.AddError("delete instance resource", err.Error())

				return
			}
		}
		// The server path removes the instance record itself; confirm.
		if err := d.waitForInstanceGone(ctx); err != nil {
			resp.Diagnostics.AddError("delete instance resource", err.Error())
		}

		return
	}

	// Phase 2: the instance is gone; its servers should be too. One that is
	// still being deprovisioned is waited for; one Morpheus has left behind is
	// force-deleted, as before.
	for _, serverID := range serverIDs {
		if err := d.ensureServerGone(ctx, serverID); err != nil {
			resp.Diagnostics.AddError("delete instance resource", err.Error())

			return
		}
	}
}

// instanceDeleteRefused marks an instance-API delete that Morpheus refused or
// that stopped in an error status: the cases the server fallback exists for.
type instanceDeleteRefused struct{ reason string }

func (e *instanceDeleteRefused) Error() string { return e.reason }

type instanceDeleter struct {
	client  *sdk.APIClient
	id      int64
	timeout time.Duration
}

// getInstance reads the instance; gone is true when it no longer exists.
func (d *instanceDeleter) getInstance(ctx context.Context) (*sdk.GetInstance200ResponseInstance, bool, error) {
	ret, hresp, err := d.client.InstancesAPI.GetInstance(ctx, d.id).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return nil, true, nil
		}

		return nil, false, fmt.Errorf("instance %d: GET failed %s", d.id, errfmt.ErrMsg(err, hresp))
	}
	if ret.Instance == nil {
		return nil, false, fmt.Errorf("instance %d: GET returned empty instance", d.id)
	}

	return ret.Instance, false, nil
}

// serverIDsOf lists the servers behind the instance's containers. At the
// moment the resource creates a single container, but the API allows more.
func serverIDsOf(id int64, instance *sdk.GetInstance200ResponseInstance) ([]int64, error) {
	if instance.ContainerDetails == nil {
		return nil, fmt.Errorf("instance %d: GET returned empty containers", id)
	}
	serverIDs := make([]int64, 0, len(instance.ContainerDetails))
	for _, container := range instance.ContainerDetails {
		if container.Server == nil {
			return nil, fmt.Errorf("instance %d: GET returned empty server in container", id)
		}
		if container.Server.Id == nil {
			return nil, fmt.Errorf("instance %d: GET returned empty server ID in container", id)
		}
		serverIDs = append(serverIDs, *container.Server.Id)
	}

	return serverIDs, nil
}

// deleteViaInstanceAPI waits for the instance to be deletable, deletes it, and
// waits for it to disappear. A refusal (4xx) or a removal that stops in an
// error status is returned as *instanceDeleteRefused so the caller can fall
// back; anything else is a hard error.
func (d *instanceDeleter) deleteViaInstanceAPI(ctx context.Context, instance *sdk.GetInstance200ResponseInstance) error {
	if err := d.waitWhileBusy(ctx, instance); err != nil {
		return err
	}

	result, hresp, err := d.client.InstancesAPI.DeleteInstance(ctx, d.id).Force("on").
		RemoveVolumes("on").ReleaseEIPs("on").ReleaseFloatingIps("on").Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		// A 4xx, or a 200 whose body says success:false, is Morpheus declining
		// (locked, busy, or a policy); the server fallback does not check status.
		var declined errfmt.DeleteRefusedError
		if errors.As(err, &declined) || (hresp != nil && hresp.StatusCode >= 400 && hresp.StatusCode < 500) {
			return &instanceDeleteRefused{reason: "Morpheus refused the instance delete: " + err.Error()}
		}

		return fmt.Errorf("instance %d: DELETE failed: %w", d.id, err)
	}

	if err := d.waitForInstanceGone(ctx); err != nil {
		// A removal that stopped in an error status (removalFailedError, from
		// checkDeleteStatusDone) is something the server path can still clean
		// up. A timeout or transport error cannot be helped.
		var failed *removalFailedError
		if errors.As(err, &failed) {
			return &instanceDeleteRefused{reason: "the instance removal failed: " + failed.Error()}
		}

		return err
	}

	return nil
}

// waitWhileBusy polls until the instance has left the statuses in which
// Morpheus refuses a delete, bounded by deleteBusyWait. Reaching the bound is
// not an error: the delete is attempted anyway.
func (d *instanceDeleter) waitWhileBusy(ctx context.Context, instance *sdk.GetInstance200ResponseInstance) error {
	status := func(i *sdk.GetInstance200ResponseInstance) string {
		if i == nil || i.Status == nil {
			return ""
		}

		return *i.Status
	}
	if !slices.Contains(deleteBusyStatuses, status(instance)) {
		return nil
	}
	tflog.Info(ctx, "instance is busy; waiting before delete",
		map[string]any{"instance": d.id, "status": status(instance)})

	wait := min(deleteBusyWait, d.timeout)
	poll := func() (struct{}, error) {
		current, gone, err := d.getInstance(ctx)
		if err != nil {
			return struct{}{}, backoff.Permanent(err)
		}
		if gone || !slices.Contains(deleteBusyStatuses, status(current)) {
			return struct{}{}, nil
		}

		return struct{}{}, backoff.RetryAfter(int(deletePoll.Seconds()))
	}
	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(deletePoll)),
		backoff.WithMaxElapsedTime(wait))
	var permanent *backoff.PermanentError
	if errors.As(err, &permanent) {
		return permanent.Unwrap()
	}
	// Still busy after the bound: carry on and let Morpheus decide.
	return nil
}

// waitForInstanceGone polls the instance until it is gone. An error status is
// permanent (see checkDeleteStatusDone); anything else keeps polling until the
// delete timeout.
func (d *instanceDeleter) waitForInstanceGone(ctx context.Context) error {
	poll := func() (struct{}, error) {
		current, gone, err := d.getInstance(ctx)
		if err != nil {
			return struct{}{}, backoff.Permanent(err)
		}
		if gone {
			return struct{}{}, nil
		}
		if current.Status == nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", d.id))
		}

		return struct{}{}, checkDeleteStatusDone(*current.Status, current.StatusMessage.Get())
	}
	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(deletePoll)),
		backoff.WithMaxElapsedTime(d.timeout))
	if err != nil {
		return fmt.Errorf("instance %d: waiting for the delete to finish: %w", d.id, err)
	}

	return nil
}

// ensureServerGone confirms a server went away with its instance. A server still
// deprovisioning is waited for; one that remains in any other state has been
// left behind by the instance removal and is force-deleted, as the resource
// always did before.
func (d *instanceDeleter) ensureServerGone(ctx context.Context, serverID int64) error {
	getID := sdk.GetHostIdParameter{Int64: &serverID}
	ret, hresp, err := d.client.HostsAPI.GetHost(ctx, getID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("instance %d: GET server %d failed %s", d.id, serverID, errfmt.ErrMsg(err, hresp))
	}
	status := ""
	if ret.Server != nil && ret.Server.Status != nil {
		status = *ret.Server.Status
	}
	if status == "deprovisioning" || status == "removing" {
		tflog.Info(ctx, "server is still being removed; waiting",
			map[string]any{"instance": d.id, "server": serverID, "status": status})

		return d.waitForServerGone(ctx, serverID)
	}

	tflog.Warn(ctx, "server left behind by the instance removal; force-deleting it",
		map[string]any{"instance": d.id, "server": serverID, "status": status})

	return d.forceDeleteServer(ctx, serverID)
}

// forceDeleteServer is the pre-MORPH-18513 destroy: stop the server, then
// DELETE /api/servers/{id}?force=on&removeResources=on&removeInstances=on and
// wait for it to disappear.
func (d *instanceDeleter) forceDeleteServer(ctx context.Context, serverID int64) error {
	// Stop the server if it is not already stopped, otherwise it cannot be deleted
	stopID := sdk.UpdateHostIdParameter{Int64: &serverID}
	_, hresp, err := d.client.HostsAPI.StopHost(ctx, stopID).Execute()
	if err != nil || (hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusConflict) {
		if hresp == nil || hresp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("instance %d: failed to stop server %d before delete %s",
				d.id, serverID, errfmt.ErrMsg(err, hresp))
		}

		return nil // already gone
	}

	if err := d.waitForServerStopped(ctx, serverID); err != nil {
		return err
	}

	result, hresp, err := d.client.HostsAPI.RemoveHost(ctx, stopID).Force("on").
		RemoveResources("on").RemoveInstances("on").Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		return fmt.Errorf("instance %d: DELETE server %d failed: %w", d.id, serverID, err)
	}

	return d.waitForServerGone(ctx, serverID)
}

// waitForServerStopped polls the server until it reports stopped; an error
// status is permanent.
func (d *instanceDeleter) waitForServerStopped(ctx context.Context, serverID int64) error {
	getID := sdk.GetHostIdParameter{Int64: &serverID}
	poll := func() (struct{}, error) {
		ret, hresp, err := d.client.HostsAPI.GetHost(ctx, getID).Execute()
		if err != nil {
			if hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}

			return struct{}{}, backoff.Permanent(err)
		}
		if ret.Server == nil || ret.Server.Status == nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty status", d.id, serverID))
		}

		return struct{}{}, checkStatusDone(*ret.Server.Status, []string{"provisioned", "stopped"}, StopErrorStatuses)
	}
	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(deletePoll)),
		backoff.WithMaxElapsedTime(d.timeout))
	if err != nil {
		return fmt.Errorf("instance %d: waiting for server %d to stop: %w", d.id, serverID, err)
	}

	return nil
}

// waitForServerGone polls the server until it is gone; an error status is
// permanent.
func (d *instanceDeleter) waitForServerGone(ctx context.Context, serverID int64) error {
	getID := sdk.GetHostIdParameter{Int64: &serverID}
	poll := func() (struct{}, error) {
		ret, hresp, err := d.client.HostsAPI.GetHost(ctx, getID).Execute()
		if err != nil {
			if hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}

			return struct{}{}, backoff.Permanent(err)
		}
		if ret.Server == nil || ret.Server.Status == nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("instance %d: GET server %d returned empty status", d.id, serverID))
		}

		return struct{}{}, checkStatusDone(*ret.Server.Status, nil, StopErrorStatuses)
	}
	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(deletePoll)),
		backoff.WithMaxElapsedTime(d.timeout))
	if err != nil {
		return fmt.Errorf("instance %d: waiting for server %d to be deleted: %w", d.id, serverID, err)
	}

	return nil
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

	return backoff.Permanent(&removalFailedError{reason: reason})
}

// removalFailedError is the terminal error for an instance whose removal
// Morpheus accepted but did not complete. backoff.Retry returns it unwrapped
// from backoff.Permanent, so callers match on this type, not on PermanentError.
type removalFailedError struct{ reason string }

func (e *removalFailedError) Error() string { return e.reason }
