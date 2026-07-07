// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data BackupInstanceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	if err := getBackupInstance(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// getBackupInstance dispatches to the id or name lookup based on the
// configured search terms.
func getBackupInstance(
	ctx context.Context,
	data *BackupInstanceModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getBackupInstanceByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getBackupInstanceByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

// getBackupInstanceByID fetches a single backup by ID and populates the state.
func getBackupInstanceByID(
	ctx context.Context,
	id int64,
	data *BackupInstanceModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.GetBackups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup instance %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if result == nil || result.Backup == nil {
		return fmt.Errorf("backup instance %d GET returned no backup", id)
	}

	populateBackupInstanceState(data, result.Backup)

	return nil
}

// getBackupInstanceByName looks up a backup by name and, on a unique match,
// delegates to getBackupInstanceByID so the full backup object is read
// consistently.
func getBackupInstanceByName(
	ctx context.Context,
	name string,
	data *BackupInstanceModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.ListBackups(ctx).Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup instance %q LIST failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	if result == nil {
		return errors.New(ErrorNoBackupFound)
	}

	var matching []sdk.ListBackups200ResponseAllOfBackupsInner
	for _, b := range result.Backups {
		if b.Name != nil && *b.Name == name {
			matching = append(matching, b)
		}
	}

	if len(matching) == 0 {
		return errors.New(ErrorNoBackupFound)
	}

	if len(matching) > 1 {
		return errors.New(ErrorMultipleBackups)
	}

	if matching[0].Id == nil {
		return fmt.Errorf("backup instance %q has missing ID", name)
	}

	return getBackupInstanceByID(ctx, *matching[0].Id, data, apiClient)
}

// populateBackupInstanceState maps the SDK backup response onto the model.
func populateBackupInstanceState(
	data *BackupInstanceModel,
	b *sdk.GetBackups200ResponseBackup,
) {
	data.Id = convert.Int64ToType(b.Id)
	data.Name = convert.StrToType(b.Name)
	data.Enabled = convert.BoolToType(b.Enabled)
	data.LocationType = convert.StrToType(b.LocationType)
	data.DateCreated = convert.TimeToType(b.DateCreated)
	data.LastUpdated = convert.TimeToType(b.LastUpdated)

	// container_id — nullable
	if b.ContainerId.IsSet() {
		data.ContainerId = convert.Int64ToType(b.ContainerId.Get())
	} else {
		data.ContainerId = types.Int64Null()
	}

	// instance_id — nested Instance object
	if b.Instance != nil && b.Instance.Id != nil {
		data.InstanceId = convert.Int64ToType(b.Instance.Id)
	} else {
		data.InstanceId = types.Int64Null()
	}

	// job_id — nested Job object
	if b.Job != nil && b.Job.Id != nil {
		data.JobId = convert.Int64ToType(b.Job.Id)
	} else {
		data.JobId = types.Int64Null()
	}

	// storage_provider_id — nested StorageProvider object
	if b.StorageProvider != nil && b.StorageProvider.Id != nil {
		data.StorageProviderId = convert.Int64ToType(b.StorageProvider.Id)
	} else {
		data.StorageProviderId = types.Int64Null()
	}

	// Nullable scalar fields — guarded with .IsSet()
	if b.Location.IsSet() {
		data.Location = convert.StrToType(b.Location.Get())
	} else {
		data.Location = types.StringNull()
	}

	if b.CronExpression.IsSet() {
		data.CronExpression = convert.StrToType(b.CronExpression.Get())
	} else {
		data.CronExpression = types.StringNull()
	}

	if b.RetentionCount.IsSet() {
		data.RetentionCount = convert.Int64ToType(b.RetentionCount.Get())
	} else {
		data.RetentionCount = types.Int64Null()
	}

	if b.TargetAll.IsSet() {
		data.TargetAll = convert.BoolToType(b.TargetAll.Get())
	} else {
		data.TargetAll = types.BoolNull()
	}

	if b.TargetHost.IsSet() {
		data.TargetHost = convert.StrToType(b.TargetHost.Get())
	} else {
		data.TargetHost = types.StringNull()
	}

	if b.TargetName.IsSet() {
		data.TargetName = convert.StrToType(b.TargetName.Get())
	} else {
		data.TargetName = types.StringNull()
	}

	if b.TargetPath.IsSet() {
		data.TargetPath = convert.StrToType(b.TargetPath.Get())
	} else {
		data.TargetPath = types.StringNull()
	}

	if b.TargetPort.IsSet() {
		p := b.TargetPort.Get()
		if p != nil {
			v := int64(*p)
			data.TargetPort = convert.Int64ToType(&v)
		} else {
			data.TargetPort = types.Int64Null()
		}
	} else {
		data.TargetPort = types.Int64Null()
	}

	if b.TargetUsername.IsSet() {
		data.TargetUsername = convert.StrToType(b.TargetUsername.Get())
	} else {
		data.TargetUsername = types.StringNull()
	}

	if b.VolumePath.IsSet() {
		data.VolumePath = convert.StrToType(b.VolumePath.Get())
	} else {
		data.VolumePath = types.StringNull()
	}

	// backup_provider — nested {id, code, name}
	data.BackupProvider = backupProviderValue(b.BackupProvider)

	// backup_repository — nested {id, name}
	data.BackupRepository = backupRepositoryValue(b.BackupRepository)

	// backup_type — nested {id, code, name}
	data.BackupType = backupTypeValue(b.BackupType)

	// schedule — nested {id, cron, name}
	data.Schedule = scheduleValue(b.Schedule)
}

func backupProviderValue(in *sdk.GetBackups200ResponseBackupBackupProvider) BackupProviderValue {
	if in == nil {
		return NewBackupProviderValueNull()
	}

	return BackupProviderValue{
		Id:    convert.Int64ToType(in.Id),
		Code:  convert.StrToType(in.Code),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func backupRepositoryValue(in *sdk.GetBackups200ResponseBackupBackupRepository) BackupRepositoryValue {
	if in == nil {
		return NewBackupRepositoryValueNull()
	}

	return BackupRepositoryValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func backupTypeValue(in *sdk.GetBackups200ResponseBackupBackupType) BackupTypeValue {
	if in == nil {
		return NewBackupTypeValueNull()
	}

	return BackupTypeValue{
		Id:    convert.Int64ToType(in.Id),
		Code:  convert.StrToType(in.Code),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func scheduleValue(in *sdk.GetBackups200ResponseBackupSchedule) ScheduleValue {
	if in == nil {
		return NewScheduleValueNull()
	}

	return ScheduleValue{
		Id:    convert.Int64ToType(in.Id),
		Cron:  convert.StrToType(in.Cron),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}
