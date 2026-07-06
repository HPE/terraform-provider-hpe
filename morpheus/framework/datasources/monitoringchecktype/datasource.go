// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringchecktype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                       = "read monitoring check type data source"
	ErrorNoValidSearchTerms       = `no valid search terms - an id or name is required`
	ErrorRunningPreApply          = `Error running pre-apply plan: exit status 1`
	ErrorNoMonitoringCheckType    = `no monitoring check type found`
	ErrorMultipleMonitoringChecks = `multiple monitoring check types were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "monitoring_check_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = MonitoringCheckTypeDataSourceSchema(ctx)
}

func checkTypeAsState(
	ct *sdk.GetCheckTypes200ResponseCheckType,
) MonitoringCheckTypeModel {
	return MonitoringCheckTypeModel{
		Id:              convert.Int64ToType(ct.Id),
		Name:            convert.StrToType(ct.Name),
		Code:            convert.StrToType(ct.Code),
		DefaultInterval: convert.Int64ToType(ct.DefaultInterval),
		MetricName:      convert.StrToType(ct.MetricName),
		InUptime:        convert.BoolToType(ct.InUptime),
		CreateIncident:  convert.BoolToType(ct.CreateIncident),
		PushOnly:        convert.BoolToType(ct.PushOnly),
		TunnelSupported: convert.BoolToType(ct.TunnelSupported),
	}
}

func getCheckTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckTypes200ResponseCheckType, error) {
	r, hresp, err := apiClient.ChecksAPI.GetCheckTypes(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for monitoring check type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.CheckType == nil {
		return nil, fmt.Errorf("GET failed for monitoring check type %d: response missing checkType", id)
	}

	ct := *r.CheckType

	return &ct, nil
}

func getCheckTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckTypes200ResponseCheckType, error) {
	rs, hresp, err := apiClient.ChecksAPI.ListCheckTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for monitoring check type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.CheckTypes)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling check types: %w", marshalErr)
	}

	var checkTypes []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &checkTypes); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding check types: %w", unmarshalErr)
	}

	var matchedID int64

	var matchCount int

	for _, ct := range checkTypes {
		if ct.Name != nil && *ct.Name == name {
			if ct.Id != nil {
				matchedID = *ct.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoMonitoringCheckType)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleMonitoringChecks)
	}

	return getCheckTypeByID(ctx, matchedID, apiClient)
}

func getCheckType(
	ctx context.Context,
	config *MonitoringCheckTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckTypes200ResponseCheckType, error) {
	if !config.Id.IsNull() {
		return getCheckTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getCheckTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config MonitoringCheckTypeModel

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

	ct, err := getCheckType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := checkTypeAsState(ct)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
