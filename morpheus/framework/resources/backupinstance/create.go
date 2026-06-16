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

	// Compute the backup_type_code with the /api/backups/create endpoint if it's null in config.
	if plan.BackupTypeCode.IsNull() || plan.BackupTypeCode.IsUnknown() {
		validateBackupInstanceReq := createToValidateBackupInstance(backup.BackupInstance)
		result, httpResp, err := client.BackupsAPI.ValidateBackupCreate(ctx).ValidateBackupCreateRequest(sdk.ValidateBackupCreateRequest{
			Backup: sdk.ValidateBackupCreateRequestBackup{
				BackupInstance1: validateBackupInstanceReq,
			},
		}).Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

			return
		}

		// if this is non-nil, request was a success...
		if result.Backup == nil {
			resp.Diagnostics.AddError(createOperation, "Backup is nil in the create validation response")

			return
		}
		// ...and result.BackupTypes should be populated with a single computed code.
		if len(result.BackupTypes) != 1 {
			resp.Diagnostics.AddError(createOperation, "BackupTypes does not have 1 computed backup type")

			return
		}

		if result.BackupTypes[0].Code == nil {
			resp.Diagnostics.AddError(createOperation, "BackupTypes contains nil code")

			return
		}

		backup.BackupInstance.BackupType = *result.BackupTypes[0].Code

	} else if !plan.BackupTypeCode.IsNull() && !plan.BackupTypeCode.IsUnknown() {
		// otherwise, set it as normal
		backup.BackupInstance.BackupType = plan.BackupTypeCode.ValueString()
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

func createToValidateBackupInstance(source *sdk.BackupInstance) *sdk.BackupInstance1 {
	if source == nil {
		return nil
	}

	target := &sdk.BackupInstance1{
		LocationType: source.LocationType,
		Name:         source.Name,
		InstanceId:   source.InstanceId,
		ContainerId:  source.ContainerId,
		BackupType:   source.BackupType,
		JobAction:    source.JobAction,
	}

	if source.JobId != nil {
		target.JobId = source.JobId
	}

	if source.JobName != nil {
		target.JobName = source.JobName
	}

	if source.JobSchedule != nil {
		target.JobSchedule = source.JobSchedule
	}

	if source.RetentionCount != nil {
		target.RetentionCount = source.RetentionCount
	}

	if source.Target != nil {
		target.Target = source.Target
	}

	if source.BackupRepository != nil {
		target.BackupRepository = source.BackupRepository
	}

	if source.ProviderBackupType != nil {
		target.ProviderBackupType = source.ProviderBackupType
	}

	if source.BackupJob != nil {
		target.BackupJob = &sdk.BackupInstance1BackupJob{
			SyntheticFullEnabled:  source.BackupJob.SyntheticFullEnabled,
			SyntheticFullSchedule: source.BackupJob.SyntheticFullSchedule,
			AdditionalProperties:  source.BackupJob.AdditionalProperties,
		}
	}

	target.AdditionalProperties = source.AdditionalProperties

	return target
}
