// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringalert

import (
	"context"
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
	summary                 = "read monitoring alert data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoAlertFound       = `no monitoring alert found`
	ErrorMultipleAlerts     = `multiple monitoring alerts were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "monitoring_alert"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = MonitoringAlertDataSourceSchema(ctx)
}

func monitoringAlertAsState(
	alert *sdk.GetAlerts200ResponseAllOfAlert,
) MonitoringAlertModel {
	return MonitoringAlertModel{
		Id:          convert.Int64ToType(alert.Id),
		Name:        convert.StrToType(alert.Name),
		Active:      convert.BoolToType(alert.Active),
		AllApps:     convert.BoolToType(alert.AllApps),
		AllChecks:   convert.BoolToType(alert.AllChecks),
		AllGroups:   convert.BoolToType(alert.AllGroups),
		MinDuration: convert.Int64ToType(alert.MinDuration),
		MinSeverity: convert.StrToType(alert.MinSeverity),
		DateCreated: convert.TimeToType(alert.DateCreated),
		LastUpdated: convert.TimeToType(alert.LastUpdated),
	}
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetAlerts200ResponseAllOfAlert, error) {
	r, hresp, err := apiClient.AlertsAPI.GetAlerts(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring alert %d: %s", id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.Alert == nil {
		return nil, fmt.Errorf(
			"GET failed for monitoring alert %d: response missing alert", id,
		)
	}

	alert := *r.Alert

	return &alert, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetAlerts200ResponseAllOfAlert, error) {
	rs, hresp, err := apiClient.AlertsAPI.ListAlerts(ctx).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring alert %s: %s", name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matched []sdk.ListAlerts200ResponseAllOfAlertsInner

	for _, o := range rs.Alerts {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoAlertFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleAlerts)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for monitoring alert %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getMonitoringAlert(
	ctx context.Context,
	config *MonitoringAlertModel,
	apiClient *sdk.APIClient,
) (*sdk.GetAlerts200ResponseAllOfAlert, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config MonitoringAlertModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	alert, err := getMonitoringAlert(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := monitoringAlertAsState(alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
