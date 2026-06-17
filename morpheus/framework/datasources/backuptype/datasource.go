// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuptype

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read backup type data source"

// Error messages surfaced by the backup type data source.
const (
	ErrorNoBackupTypeFound   = `no backup type found`
	ErrorMultipleBackupTypes = `multiple backup types were returned`
	ErrorNoValidSearchTerms  = `no valid search terms - an id or name is required`
	ErrorRunningPreApply     = `Error running pre-apply plan: exit status 1`
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = BackupTypeDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data BackupTypeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	if err := getBackupType(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// getBackupType dispatches to the id or name lookup based on the configured
// search terms.
func getBackupType(
	ctx context.Context,
	data *BackupTypeModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getBackupTypeByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getBackupTypeByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

// getBackupTypeByID fetches a single backup type by ID and populates the state.
func getBackupTypeByID(
	ctx context.Context,
	id int64,
	data *BackupTypeModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.GetBackupType(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup type %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if result == nil || result.BackupType == nil {
		return fmt.Errorf("backup type %d GET returned no backup type", id)
	}

	populateBackupTypeState(data, result.BackupType)

	return nil
}

// getBackupTypeByName looks up a backup type by name and, on a unique match,
// delegates to getBackupTypeByID so that the full object is read consistently.
func getBackupTypeByName(
	ctx context.Context,
	name string,
	data *BackupTypeModel,
	apiClient *sdk.APIClient,
) error {
	result, hresp, err := apiClient.BackupsAPI.ListBackupTypes(ctx).Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup type %q LIST failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	if result == nil {
		return errors.New(ErrorNoBackupTypeFound)
	}

	var matching []sdk.ListBackupTypes200ResponseAllOfBackupTypesInner
	for _, bt := range result.BackupTypes {
		if bt.Name != nil && *bt.Name == name {
			matching = append(matching, bt)
		}
	}

	if len(matching) == 0 {
		return errors.New(ErrorNoBackupTypeFound)
	}

	if len(matching) > 1 {
		return errors.New(ErrorMultipleBackupTypes)
	}

	if matching[0].Id == nil {
		return fmt.Errorf("backup type %q has missing ID", name)
	}

	return getBackupTypeByID(ctx, *matching[0].Id, data, apiClient)
}

// populateBackupTypeState maps the SDK backup type response onto the model.
//
//nolint:funlen // mapping every backup type attribute requires length
func populateBackupTypeState(
	data *BackupTypeModel,
	bt *sdk.GetBackupType200ResponseBackupType,
) {
	data.Active = convert.BoolToType(bt.Active)
	data.BackupFormat = convert.StrToType(bt.BackupFormat.Get())
	data.Code = convert.StrToType(bt.Code)
	data.ContainerCategory = convert.StrToType(bt.ContainerCategory.Get())
	data.ContainerFormat = convert.StrToType(bt.ContainerFormat.Get())
	data.ContainerType = convert.StrToType(bt.ContainerType.Get())
	data.CopyToStore = convert.BoolToType(bt.CopyToStore)
	data.DownloadEnabled = convert.BoolToType(bt.DownloadEnabled)
	data.DownloadFromStoreOnly = convert.BoolToType(bt.DownloadFromStoreOnly)
	data.HasCopyToStore = convert.BoolToType(bt.HasCopyToStore)
	data.HasStreamToStore = convert.BoolToType(bt.HasStreamToStore)
	data.Id = convert.Int64ToType(bt.Id)
	data.IsEmbedded = convert.BoolToType(bt.IsEmbedded)
	data.IsPlugin = convert.BoolToType(bt.IsPlugin.Get())
	data.Name = convert.StrToType(bt.Name)
	data.ProviderCode = convert.StrToType(bt.ProviderCode.Get())
	data.PruneResultsOnRestoreExisting = convert.BoolToType(bt.PruneResultsOnRestoreExisting)
	data.RestoreExistingEnabled = convert.BoolToType(bt.RestoreExistingEnabled)
	data.RestoreNewEnabled = convert.BoolToType(bt.RestoreNewEnabled)
	data.RestoreNewMode = convert.StrToType(bt.RestoreNewMode.Get())
	data.RestoreType = convert.StrToType(bt.RestoreType.Get())
	data.RestrictTargets = convert.BoolToType(bt.RestrictTargets)
	data.ViewSet = convert.StrToType(bt.ViewSet.Get())
}
