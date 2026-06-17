// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read backup data source"

// Error messages surfaced by the backup data source.
const (
	ErrorNoBackupFound      = `no backup found`
	ErrorMultipleBackups    = `multiple backups were returned`
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = BackupDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data BackupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	if err := getBackup(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// getBackup dispatches to the id or name lookup based on the configured search
// terms.
func getBackup(
	ctx context.Context,
	data *BackupModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getBackupByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getBackupByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

// getBackupByID fetches a single backup by ID and populates the state.
func getBackupByID(
	ctx context.Context,
	id int64,
	data *BackupModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.GetBackups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if result == nil || result.Backup == nil {
		return fmt.Errorf("backup %d GET returned no backup", id)
	}

	populateBackupState(data, result.Backup)

	return nil
}

// getBackupByName looks up a backup by name and, on a unique match, delegates
// to getBackupByID so that the full backup object is read consistently.
func getBackupByName(
	ctx context.Context,
	name string,
	data *BackupModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.ListBackups(ctx).Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup %q LIST failed: %s", name, errfmt.ErrMsg(err, hresp))
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
		return fmt.Errorf("backup %q has missing ID", name)
	}

	return getBackupByID(ctx, *matching[0].Id, data, apiClient)
}

// populateBackupState maps the SDK backup response onto the flattened model.
//
//nolint:funlen // mapping every backup attribute requires length
func populateBackupState(
	data *BackupModel,
	b *sdk.GetBackups200ResponseBackup,
) {
	data.Id = convert.Int64ToType(b.Id)
	data.Name = convert.StrToType(b.Name)
	data.ContainerId = convert.Int64ToType(b.ContainerId.Get())
	data.CronExpression = convert.StrToType(b.CronExpression.Get())
	data.DateCreated = timeToType(b.DateCreated)
	data.Enabled = convert.BoolToType(b.Enabled)
	data.LastStatus = convert.StrToType(b.LastStatus.Get())
	data.LastUpdated = timeToType(b.LastUpdated)
	data.Location = convert.StrToType(b.Location.Get())
	data.LocationType = convert.StrToType(b.LocationType)
	data.NextFire = timeToType(b.NextFire.Get())
	data.RetentionCount = convert.Int64ToType(b.RetentionCount.Get())
	data.TargetAll = convert.BoolToType(b.TargetAll.Get())
	data.TargetHost = convert.StrToType(b.TargetHost.Get())
	data.TargetName = convert.StrToType(b.TargetName.Get())
	data.TargetPassword = convert.StrToType(b.TargetPassword.Get())
	data.TargetPasswordHash = convert.StrToType(b.TargetPasswordHash.Get())
	data.TargetPath = convert.StrToType(b.TargetPath.Get())
	data.TargetPort = int32ToType(b.TargetPort.Get())
	data.TargetUsername = convert.StrToType(b.TargetUsername.Get())
	data.VolumePath = convert.StrToType(b.VolumePath.Get())

	data.BackupProvider = backupProviderValue(b.BackupProvider)
	data.BackupRepository = backupRepositoryValue(b.BackupRepository)
	data.BackupType = backupTypeValue(b.BackupType)
	data.Host = hostValue(b.Server)
	data.Instance = instanceValue(b.Instance)
	data.Job = jobValue(b.Job)
	data.LastResult = lastResultValue(b.LastResult)
	data.Schedule = scheduleValue(b.Schedule)
	data.Stats = statsValue(b.Stats)
	data.StorageProvider = storageProviderValue(b.StorageProvider)
}

// timeToType formats an optional time as an RFC3339 string value.
func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

// int32ToType converts an optional int32 to a Terraform Int64 value.
func int32ToType(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

func backupProviderValue(
	in *sdk.GetBackups200ResponseBackupBackupProvider,
) BackupProviderValue {
	if in == nil {
		return NewBackupProviderValueNull()
	}

	return BackupProviderValue{
		Code:  convert.StrToType(in.Code),
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func backupRepositoryValue(
	in *sdk.GetBackups200ResponseBackupBackupRepository,
) BackupRepositoryValue {
	if in == nil {
		return NewBackupRepositoryValueNull()
	}

	return BackupRepositoryValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func backupTypeValue(
	in *sdk.GetBackups200ResponseBackupBackupType,
) BackupTypeValue {
	if in == nil {
		return NewBackupTypeValueNull()
	}

	return BackupTypeValue{
		Code:  convert.StrToType(in.Code),
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func hostValue(
	in *sdk.GetBackups200ResponseBackupServer,
) HostValue {
	if in == nil {
		return NewHostValueNull()
	}

	return HostValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func instanceValue(
	in *sdk.GetBackups200ResponseBackupInstance,
) InstanceValue {
	if in == nil {
		return NewInstanceValueNull()
	}

	return InstanceValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func jobValue(
	in *sdk.GetBackups200ResponseBackupJob,
) JobValue {
	if in == nil {
		return NewJobValueNull()
	}

	return JobValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func lastResultValue(
	in *sdk.GetBackups200ResponseBackupLastResult,
) LastResultValue {
	if in == nil {
		return NewLastResultValueNull()
	}

	return LastResultValue{
		DateCreated: timeToType(in.DateCreated),
		Id:          convert.Int64ToType(in.Id),
		Status:      convert.StrToType(in.Status),
		state:       attr.ValueStateKnown,
	}
}

func scheduleValue(
	in *sdk.GetBackups200ResponseBackupSchedule,
) ScheduleValue {
	if in == nil {
		return NewScheduleValueNull()
	}

	return ScheduleValue{
		Cron:  convert.StrToType(in.Cron),
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func statsValue(
	in *sdk.GetBackups200ResponseBackupStats,
) StatsValue {
	if in == nil {
		return NewStatsValueNull()
	}

	return StatsValue{
		AvgSize:         convert.Int64ToType(in.AvgSize),
		FailRate:        types.Float64PointerValue(in.FailRate),
		Failed:          convert.Int64ToType(in.Failed),
		LastFiveResults: convert.StrSliceToSet(in.LastFiveResults),
		Success:         convert.Int64ToType(in.Success),
		SuccessRate:     types.Float64PointerValue(in.SuccessRate),
		TotalCompleted:  convert.Int64ToType(in.TotalCompleted),
		TotalSize:       convert.Int64ToType(in.TotalSize),
		state:           attr.ValueStateKnown,
	}
}

func storageProviderValue(
	in *sdk.GetBackups200ResponseBackupStorageProvider,
) StorageProviderValue {
	if in == nil {
		return NewStorageProviderValueNull()
	}

	return StorageProviderValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}
