// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package monitoringcheck implements a data source for monitoring_check
package monitoringcheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                       = "read monitoring check data source"
	ErrorNoValidSearchTerms       = `no valid search terms - an id or name is required`
	ErrorRunningPreApply          = `Error running pre-apply plan: exit status 1`
	ErrorNoMonitoringCheckFound   = `no monitoring check found`
	ErrorMultipleMonitoringChecks = `multiple monitoring checks were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "monitoring_check"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = MonitoringCheckDataSourceSchema(ctx)
}

func monitoringCheckAsState(
	check *sdk.GetChecks200ResponseCheck,
) MonitoringCheckModel {
	state := MonitoringCheckModel{
		Id:          convert.Int64ToType(check.Id),
		Name:        convert.StrToType(check.Name),
		Active:      convert.BoolToType(check.Active),
		InUptime:    convert.BoolToType(check.InUptime),
		Severity:    convert.StrToType(check.Severity),
		Description: convert.StrToType(check.Description.Get()),
	}

	if check.CheckInterval.IsSet() {
		state.CheckInterval = convert.Int64ToType(check.CheckInterval.Get())
	} else {
		state.CheckInterval = types.Int64Null()
	}

	// Handle dynamic config via JSON round-trip (mirrors networkdhcpserver pattern)
	state.Config = types.DynamicNull()
	if check.Config != nil {
		raw, err := json.Marshal(check.Config)
		if err == nil && string(raw) != "{}" && string(raw) != "null" {
			state.Config = types.DynamicValue(types.StringValue(string(raw)))
		}
	}

	return state
}

func getCheckByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetChecks200ResponseCheck, error) {
	r, hresp, err := apiClient.ChecksAPI.GetChecks(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring check %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.Check == nil {
		return nil, fmt.Errorf("GET failed for monitoring check %d: empty response", id)
	}

	return r.Check, nil
}

func getCheckByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetChecks200ResponseCheck, error) {
	rs, hresp, err := apiClient.ChecksAPI.ListChecks(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring check %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matchedIDs []int64

	for i := range rs.Checks {
		if rs.Checks[i].Name == nil || *rs.Checks[i].Name != name {
			continue
		}
		if rs.Checks[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *rs.Checks[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoMonitoringCheckFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleMonitoringChecks)
	}

	return getCheckByID(ctx, matchedIDs[0], apiClient)
}

func getCheck(
	ctx context.Context,
	config *MonitoringCheckModel,
	apiClient *sdk.APIClient,
) (*sdk.GetChecks200ResponseCheck, error) {
	if !config.Id.IsNull() {
		return getCheckByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getCheckByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config MonitoringCheckModel

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

	check, err := getCheck(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := monitoringCheckAsState(check)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
