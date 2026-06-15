// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *backupInstanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// instance_id is required
	instanceID := plan.InstanceId.ValueInt64()

	instance := sdk.BackupInstance{
		Name:         plan.Name.ValueString(),
		InstanceId:   instanceID,
		LocationType: "instance",
		JobAction:    "addTo",
		JobId:        plan.JobId.ValueInt64Pointer(),
	}
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		instance.Target = plan.StorageProviderId.ValueInt64Pointer()
	}

	// container_id is optional/computed. Use the user-supplied value when set;
	// otherwise resolve it from the instance's container list (the same thing
	// the API does on a dry-run create to /api/backups/create).
	if !plan.ContainerId.IsNull() && !plan.ContainerId.IsUnknown() {
		instance.ContainerId = plan.ContainerId.ValueInt64()
	} else {
		instResult, instHResp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if err != nil || instHResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, instHResp))

			return
		}

		if instResult.Instance == nil {
			resp.Diagnostics.AddError(
				createOperation,
				fmt.Sprintf("instance %d returned no instance", instanceID),
			)

			return
		}

		if len(instResult.Instance.Containers) == 0 {
			resp.Diagnostics.AddError(
				createOperation,
				fmt.Sprintf("instance %d has no containers", instanceID),
			)

			return
		}

		instance.ContainerId = instResult.Instance.Containers[0]
	}

	backup := sdk.AddBackupsRequestBackup{
		BackupInstance: &instance,
	}

	result, httpResp, err := client.BackupsAPI.AddBackups(ctx).AddBackupsRequest(sdk.AddBackupsRequest{
		Backup: backup,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	if result.Backup == nil || result.Backup.Id == nil {
		resp.Diagnostics.AddError(createOperation, "Backup ID is nil in the create response")

		return
	}
	id := *result.Backup.Id

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// The backup was created, but reading it back failed. Taint the resource
		// so the created backup is not leaked from Terraform's perspective.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "backup_instance",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
