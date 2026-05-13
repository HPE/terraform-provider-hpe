// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read load balancer monitor data source"

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_monitor"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerMonitorDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerMonitorModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("failed to create client: %s", err.Error()))

		return
	}

	state, err := getLoadBalancerMonitor(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func getLoadBalancerMonitor(
	ctx context.Context,
	config LoadBalancerMonitorModel,
	client *sdk.APIClient,
) (*LoadBalancerMonitorModel, error) {
	loadBalancerID := config.LoadBalancerId.ValueInt64()

	if !config.Id.IsNull() {
		return getLoadBalancerMonitorByID(ctx, loadBalancerID, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getLoadBalancerMonitorByName(ctx, loadBalancerID, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getLoadBalancerMonitorByID(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (*LoadBalancerMonitorModel, error) {
	resp, hresp, err := client.LoadBalancersAPI.
		GetLoadBalancerMonitor(ctx, loadBalancerID, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load balancer monitor %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	m := resp.GetLoadBalancerMonitor()
	state := populateLoadBalancerMonitorState(ctx, loadBalancerID, &m)

	return state, nil
}

func getLoadBalancerMonitorByName(
	ctx context.Context,
	loadBalancerID int64,
	name string,
	client *sdk.APIClient,
) (*LoadBalancerMonitorModel, error) {
	monitors, hresp, err := client.LoadBalancersAPI.
		ListLoadBalancerMonitors(ctx, loadBalancerID).
		Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load balancer monitor %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matching []sdk.ListLoadBalancerMonitors200ResponseAllOfLoadBalancerMonitorsInner
	for _, m := range monitors.GetLoadBalancerMonitors() {
		if mName, ok := m.GetNameOk(); ok && *mName == name {
			matching = append(matching, m)
		}
	}

	if len(matching) == 0 {
		return nil, fmt.Errorf("load balancer monitor %s not found", name)
	}

	if len(matching) > 1 {
		var ids []string
		for _, m := range matching {
			if id, ok := m.GetIdOk(); ok {
				ids = append(ids, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple load balancer monitors found with name %s. IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(ids, ", "),
		)
	}

	id, ok := matching[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("load balancer monitor %s has missing ID", name)
	}

	return getLoadBalancerMonitorByID(ctx, loadBalancerID, *id, client)
}

//nolint:funlen,cyclop // mapping all fields requires length
func populateLoadBalancerMonitorState(
	ctx context.Context,
	loadBalancerID int64,
	m *sdk.GetLoadBalancerMonitor200ResponseLoadBalancerMonitor,
) *LoadBalancerMonitorModel {
	state := &LoadBalancerMonitorModel{}

	state.Id = convert.Int64ToType(m.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(m.Name)
	state.Description = convert.StrToType(m.Description)
	state.MonitorType = convert.StrToType(m.MonitorType)
	state.MonitorInterval = convert.Int64ToType(m.MonitorInterval)
	state.MonitorTimeout = convert.Int64ToType(m.MonitorTimeout)

	state.SendData = convert.StrToType(m.SendData.Get())
	state.SendVersion = convert.StrToType(m.SendVersion)
	state.SendType = convert.StrToType(m.SendType)
	state.ReceiveData = convert.StrToType(m.ReceiveData.Get())
	state.ReceiveCode = convert.StrToType(m.ReceiveCode)
	state.MonitorUsername = convert.StrToType(m.MonitorUsername.Get())
	state.MonitorDestination = convert.StrToType(m.MonitorDestination)
	state.FallCount = convert.Int64ToType(m.FallCount)
	state.RiseCount = convert.Int64ToType(m.RiseCount)
	state.AliasPort = convert.Int64ToType(m.AliasPort)
	state.DataLength = convert.Int64ToType(m.DataLength.Get())
	state.MaxRetry = convert.Int64ToType(m.MaxRetry)

	// Data-source-only fields
	state.Code = convert.StrToType(m.Code.Get())
	state.Category = convert.StrToType(m.Category.Get())
	state.Visibility = convert.StrToType(m.Visibility)
	state.Enabled = convert.BoolToType(m.Enabled)
	state.Status = convert.StrToType(m.Status)
	state.StatusMessage = convert.StrToType(m.StatusMessage.Get())
	state.ExternalId = convert.StrToType(m.ExternalId)
	state.InternalId = convert.StrToType(m.InternalId)

	if sd := m.StatusDate.Get(); sd != nil {
		state.StatusDate = types.StringValue(sd.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		state.StatusDate = types.StringNull()
	}

	// Config — structured nested object.
	if m.Config != nil {
		cfg := m.GetConfig()

		var monitorObjVal basetypes.ObjectValue
		if cfg.Monitor != nil {
			monitorInner := cfg.GetMonitor()
			monitorObjVal = types.ObjectValueMust(
				MonitorValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id": convert.Int64ToType(monitorInner.Id),
				},
			)
		} else {
			monitorObjVal = types.ObjectNull(
				MonitorValue{}.AttributeTypes(ctx),
			)
		}

		var monitorConfigDyn basetypes.DynamicValue
		if mcVal := cfg.MonitorConfig.Get(); mcVal != nil {
			monitorConfigDyn = types.DynamicValue(
				types.StringValue(*mcVal),
			)
		} else {
			monitorConfigDyn = types.DynamicNull()
		}

		state.Config = ConfigValue{
			Monitor:       monitorObjVal,
			MonitorConfig: monitorConfigDyn,
			state:         attr.ValueStateKnown,
		}
	} else {
		state.Config = NewConfigValueNull()
	}

	return state
}
