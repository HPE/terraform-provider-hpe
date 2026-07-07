// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read backup instance data source"

// Error messages surfaced by the backup instance data source.
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
	resp.TypeName = req.ProviderTypeName + "_" + "backup_instance"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = BackupInstanceDataSourceSchema(ctx)
}

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
	data.ContainerId = convert.Int64ToType(b.ContainerId.Get())

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
}
