// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read backup host data source"

// Error messages surfaced by the backup host data source.
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
	resp.TypeName = req.ProviderTypeName + "_" + "backup_host"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = BackupHostDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data BackupHostModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	if err := getBackupHost(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// getBackupHost dispatches to the id or name lookup based on the configured
// search terms.
func getBackupHost(
	ctx context.Context,
	data *BackupHostModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getBackupHostByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getBackupHostByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

// getBackupHostByID fetches a single backup by ID and populates the state.
func getBackupHostByID(
	ctx context.Context,
	id int64,
	data *BackupHostModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.GetBackups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup host %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if result == nil || result.Backup == nil {
		return fmt.Errorf("backup host %d GET returned no backup", id)
	}

	populateBackupHostState(data, result.Backup)

	return nil
}

// getBackupHostByName looks up a backup by name and, on a unique match,
// delegates to getBackupHostByID so the full backup object is read
// consistently.
func getBackupHostByName(
	ctx context.Context,
	name string,
	data *BackupHostModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.ListBackups(ctx).Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup host %q LIST failed: %s", name, errfmt.ErrMsg(err, hresp))
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
		return fmt.Errorf("backup host %q has missing ID", name)
	}

	return getBackupHostByID(ctx, *matching[0].Id, data, apiClient)
}

// populateBackupHostState maps the SDK backup response onto the model.
func populateBackupHostState(
	data *BackupHostModel,
	b *sdk.GetBackups200ResponseBackup,
) {
	data.Id = convert.Int64ToType(b.Id)
	data.Name = convert.StrToType(b.Name)
	data.Enabled = convert.BoolToType(b.Enabled)
	data.TargetPath = convert.StrToType(b.TargetPath.Get())

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

	// host — nested Server object
	data.Host = hostValue(b.Server)
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
