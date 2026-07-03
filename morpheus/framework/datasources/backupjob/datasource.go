// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupjob

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read backup job data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoBackupJobFound   = `no backup job found`
	ErrorMultipleBackupJobs = `multiple backup jobs were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "backup_job"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = BackupJobDataSourceSchema(ctx)
}

// timeToType formats a nullable timestamp as an RFC 3339 string.
func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

func backupJobAsState(
	j *sdk.GetBackupJobs200ResponseJob,
) BackupJobModel {
	return BackupJobModel{
		Id:             convert.Int64ToType(j.Id),
		Name:           convert.StrToType(j.Name),
		CronExpression: convert.StrToType(j.CronExpression.Get()),
		DateCreated:    timeToType(j.DateCreated),
		Enabled:        convert.BoolToType(j.Enabled),
		ExternalId:     convert.StrToType(j.ExternalId.Get()),
		LastUpdated:    timeToType(j.LastUpdated),
		NextFire:       timeToType(j.NextFire.Get()),
		RetentionCount: convert.Int64ToType(j.RetentionCount.Get()),
		Source:         convert.StrToType(j.Source),
		Visibility:     convert.StrToType(j.Visibility),
	}
}

func getBackupJobByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetBackupJobs200ResponseJob, error) {
	r, hresp, err := apiClient.BackupsAPI.GetBackupJobs(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for backup job %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.Job == nil {
		return nil, fmt.Errorf("GET failed for backup job %d: response missing job", id)
	}

	j := *r.Job

	return &j, nil
}

func getBackupJobByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetBackupJobs200ResponseJob, error) {
	rs, hresp, err := apiClient.BackupsAPI.ListBackupJobs(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for backup job %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.Jobs)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling backup jobs: %w", marshalErr)
	}

	var backupJobs []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &backupJobs); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding backup jobs: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, j := range backupJobs {
		if j.Name != nil && *j.Name == name {
			if j.Id != nil {
				matchedID = *j.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoBackupJobFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleBackupJobs)
	}

	return getBackupJobByID(ctx, matchedID, apiClient)
}

func getBackupJob(
	ctx context.Context,
	config *BackupJobModel,
	apiClient *sdk.APIClient,
) (*sdk.GetBackupJobs200ResponseJob, error) {
	if !config.Id.IsNull() {
		return getBackupJobByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getBackupJobByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config BackupJobModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	j, err := getBackupJob(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := backupJobAsState(j)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
