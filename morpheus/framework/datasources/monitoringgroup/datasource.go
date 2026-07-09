// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringgroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read monitoring group data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoGroupFound       = `no monitoring group found`
	ErrorMultipleGroups     = `multiple monitoring groups were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "monitoring_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = MonitoringGroupDataSourceSchema(ctx)
}

func monitoringGroupAsState(
	grp *sdk.GetCheckGroups200ResponseCheckGroup,
) MonitoringGroupModel {
	state := MonitoringGroupModel{
		Id:          convert.Int64ToType(grp.Id),
		Name:        convert.StrToType(grp.Name),
		Description: convert.StrToType(grp.Description.Get()),
		InUptime:    convert.BoolToType(grp.InUptime),
		MinHappy:    convert.Int64ToType(grp.MinHappy),
		Severity:    convert.StrToType(grp.Severity),
	}

	state.DateCreated = convert.TimeToType(grp.DateCreated)
	state.LastUpdated = convert.TimeToType(grp.LastUpdated)

	// check_type nested object {id, code, name, metric_name}
	state.CheckType = NewCheckTypeValueNull()
	if grp.CheckType != nil {
		state.CheckType = CheckTypeValue{
			Id:         convert.Int64ToType(grp.CheckType.Id),
			Code:       convert.StrToType(grp.CheckType.Code),
			Name:       convert.StrToType(grp.CheckType.Name),
			MetricName: convert.StrToType(grp.CheckType.MetricName),
			state:      attr.ValueStateKnown,
		}
	}

	return state
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckGroups200ResponseCheckGroup, error) {
	r, hresp, err := apiClient.ChecksAPI.GetCheckGroups(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring group %d: %s", id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.CheckGroup == nil {
		return nil, fmt.Errorf(
			"GET failed for monitoring group %d: response missing checkGroup", id,
		)
	}

	grp := *r.CheckGroup

	return &grp, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckGroups200ResponseCheckGroup, error) {
	rs, hresp, err := apiClient.ChecksAPI.ListCheckGroups(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for monitoring group %s: %s", name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matched []sdk.ListCheckGroups200ResponseAllOfCheckGroupsInner

	for _, o := range rs.CheckGroups {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoGroupFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleGroups)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for monitoring group %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getMonitoringGroup(
	ctx context.Context,
	config *MonitoringGroupModel,
	apiClient *sdk.APIClient,
) (*sdk.GetCheckGroups200ResponseCheckGroup, error) {
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
	var config MonitoringGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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

	grp, err := getMonitoringGroup(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := monitoringGroupAsState(grp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
