// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read load balancer pool data source"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_pool"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerPoolDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerPoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("failed to create client: %s", err.Error()))

		return
	}

	state, err := getLoadBalancerPool(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func getLoadBalancerPool(
	ctx context.Context,
	config LoadBalancerPoolModel,
	client *sdk.APIClient,
) (*LoadBalancerPoolModel, error) {
	loadBalancerID := config.LoadBalancerId.ValueInt64()

	if !config.Id.IsNull() {
		return getLoadBalancerPoolByID(ctx, loadBalancerID, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getLoadBalancerPoolByName(ctx, loadBalancerID, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getLoadBalancerPoolByID(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (*LoadBalancerPoolModel, error) {
	resp, hresp, err := client.LoadBalancersAPI.
		GetLoadBalancerPool(ctx, loadBalancerID, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d pool %d GET failed: %s",
			loadBalancerID, id, errfmt.ErrMsg(err, hresp),
		)
	}

	p := resp.GetLoadBalancerPool()
	state := populateLoadBalancerPoolState(ctx, loadBalancerID, &p)

	return state, nil
}

func getLoadBalancerPoolByName(
	ctx context.Context,
	loadBalancerID int64,
	name string,
	client *sdk.APIClient,
) (*LoadBalancerPoolModel, error) {
	pools, hresp, err := client.LoadBalancersAPI.
		ListLoadBalancerPools(ctx, loadBalancerID).
		Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d pool list failed: %s",
			loadBalancerID, errfmt.ErrMsg(err, hresp),
		)
	}

	var matching []sdk.ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInner
	for _, p := range pools.GetLoadBalancerPools() {
		if pName, ok := p.GetNameOk(); ok && *pName == name {
			matching = append(matching, p)
		}
	}

	if len(matching) == 0 {
		return nil, fmt.Errorf(
			"load balancer %d pool with name %q not found",
			loadBalancerID, name,
		)
	}

	if len(matching) > 1 {
		var ids []string
		for _, p := range matching {
			if id, ok := p.GetIdOk(); ok {
				ids = append(ids, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple load balancer pools found with name %q. IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(ids, ", "),
		)
	}

	id, ok := matching[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf(
			"load balancer %d pool with name %q has missing ID",
			loadBalancerID, name,
		)
	}

	return getLoadBalancerPoolByID(ctx, loadBalancerID, *id, client)
}

//nolint:funlen,cyclop // mapping all fields requires length
func populateLoadBalancerPoolState(
	ctx context.Context,
	loadBalancerID int64,
	p *sdk.GetLoadBalancerPool200ResponseLoadBalancerPool,
) *LoadBalancerPoolModel {
	state := &LoadBalancerPoolModel{}

	// Core fields
	state.Id = convert.Int64ToType(p.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(p.Name)
	state.VipBalance = convert.StrToType(p.VipBalance)
	state.MinActive = convert.Int64ToType(p.MinActive)

	// NullableString fields — existing (fixed with .IsSet() guard)
	if p.Description.IsSet() {
		state.Description = convert.StrToType(p.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if p.VipSticky.IsSet() {
		state.VipSticky = convert.StrToType(p.VipSticky.Get())
	} else {
		state.VipSticky = types.StringNull()
	}

	if p.VipClientIpMode.IsSet() {
		state.VipClientIpMode = convert.StrToType(p.VipClientIpMode.Get())
	} else {
		state.VipClientIpMode = types.StringNull()
	}

	if p.Port.IsSet() {
		state.Port = convert.StrToType(p.Port.Get())
	} else {
		state.Port = types.StringNull()
	}

	// NullableString fields — new
	if p.AllowNat.IsSet() {
		state.AllowNat = convert.StrToType(p.AllowNat.Get())
	} else {
		state.AllowNat = types.StringNull()
	}

	if p.AllowSnat.IsSet() {
		state.AllowSnat = convert.StrToType(p.AllowSnat.Get())
	} else {
		state.AllowSnat = types.StringNull()
	}

	if p.Category.IsSet() {
		state.Category = convert.StrToType(p.Category.Get())
	} else {
		state.Category = types.StringNull()
	}

	if p.CreatedBy.IsSet() {
		state.CreatedBy = convert.StrToType(p.CreatedBy.Get())
	} else {
		state.CreatedBy = types.StringNull()
	}

	if p.DownAction.IsSet() {
		state.DownAction = convert.StrToType(p.DownAction.Get())
	} else {
		state.DownAction = types.StringNull()
	}

	if p.MaxQueueDepth.IsSet() {
		state.MaxQueueDepth = convert.StrToType(p.MaxQueueDepth.Get())
	} else {
		state.MaxQueueDepth = types.StringNull()
	}

	if p.MaxQueueTime.IsSet() {
		state.MaxQueueTime = convert.StrToType(p.MaxQueueTime.Get())
	} else {
		state.MaxQueueTime = types.StringNull()
	}

	if p.MinInService.IsSet() {
		state.MinInService = convert.StrToType(p.MinInService.Get())
	} else {
		state.MinInService = types.StringNull()
	}

	if p.MinUpAction.IsSet() {
		state.MinUpAction = convert.StrToType(p.MinUpAction.Get())
	} else {
		state.MinUpAction = types.StringNull()
	}

	if p.MinUpMonitor.IsSet() {
		state.MinUpMonitor = convert.StrToType(p.MinUpMonitor.Get())
	} else {
		state.MinUpMonitor = types.StringNull()
	}

	if p.PortType.IsSet() {
		state.PortType = convert.StrToType(p.PortType.Get())
	} else {
		state.PortType = types.StringNull()
	}

	if p.RampTime.IsSet() {
		state.RampTime = convert.StrToType(p.RampTime.Get())
	} else {
		state.RampTime = types.StringNull()
	}

	if p.VipServerIpMode.IsSet() {
		state.VipServerIpMode = convert.StrToType(p.VipServerIpMode.Get())
	} else {
		state.VipServerIpMode = types.StringNull()
	}

	// *string fields
	state.InternalId = convert.StrToType(p.InternalId)
	state.ExternalId = convert.StrToType(p.ExternalId)
	state.Status = convert.StrToType(p.Status)
	state.Visibility = convert.StrToType(p.Visibility)

	// *int64 fields
	state.NumberActive = convert.Int64ToType(p.NumberActive)
	state.NumberInService = convert.Int64ToType(p.NumberInService)
	state.HealthScore = convert.Int64ToType(p.HealthScore)
	state.PerformanceScore = convert.Int64ToType(p.PerformanceScore)
	state.HealthPenalty = convert.Int64ToType(p.HealthPenalty)
	state.SecurityPenalty = convert.Int64ToType(p.SecurityPenalty)
	state.ErrorPenalty = convert.Int64ToType(p.ErrorPenalty)

	// *bool field
	state.Enabled = convert.BoolToType(p.Enabled)

	// *time.Time fields
	if dc, ok := p.GetDateCreatedOk(); ok && dc != nil {
		state.DateCreated = types.StringValue(dc.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if lu, ok := p.GetLastUpdatedOk(); ok && lu != nil {
		state.LastUpdated = types.StringValue(lu.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// LoadBalancer — nested object with id, name, ip, type.
	if p.LoadBalancer != nil {
		lb := p.LoadBalancer

		typeVal := types.ObjectNull(TypeValue{}.AttributeTypes(ctx))

		if lbType, ok := lb.GetTypeOk(); ok && lbType != nil {
			tv, d := NewTypeValue(
				TypeValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"code": convert.StrToType(lbType.Code),
					"id":   convert.Int64ToType(lbType.Id),
					"name": convert.StrToType(lbType.Name),
				},
			)
			if !d.HasError() {
				tvObj, tvDiags := tv.ToObjectValue(ctx)
				if !tvDiags.HasError() {
					typeVal = tvObj
				}
			}
		}

		lbVal, d := NewLoadBalancerValue(
			LoadBalancerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(lb.Id),
				"ip":   convert.StrToType(lb.Ip),
				"name": convert.StrToType(lb.Name),
				"type": typeVal,
			},
		)
		if !d.HasError() {
			state.LoadBalancer = lbVal
		} else {
			state.LoadBalancer = NewLoadBalancerValueNull()
		}
	} else {
		state.LoadBalancer = NewLoadBalancerValueNull()
	}

	// Monitors — nested set of {id, name}.
	if len(p.Monitors) > 0 {
		monitorVals := make([]attr.Value, 0, len(p.Monitors))
		for _, m := range p.Monitors {
			mv := MonitorsValue{
				Id:    convert.Int64ToType(m.Id),
				Name:  convert.StrToType(m.Name),
				state: attr.ValueStateKnown,
			}

			objVal, diags := mv.ToObjectValue(ctx)
			if !diags.HasError() {
				monitorVals = append(monitorVals, objVal)
			}
		}

		setVal, diags := types.SetValue(
			MonitorsType{types.ObjectType{AttrTypes: MonitorsValue{}.AttributeTypes(ctx)}},
			monitorVals,
		)
		if !diags.HasError() {
			state.Monitors = setVal
		} else {
			state.Monitors = types.SetNull(
				MonitorsType{types.ObjectType{AttrTypes: MonitorsValue{}.AttributeTypes(ctx)}},
			)
		}
	} else {
		state.Monitors = types.SetNull(
			MonitorsType{types.ObjectType{AttrTypes: MonitorsValue{}.AttributeTypes(ctx)}},
		)
	}

	// Nodes — nested set of {id, name}.
	if len(p.Nodes) > 0 {
		nodeVals := make([]attr.Value, 0, len(p.Nodes))
		for _, n := range p.Nodes {
			nv := NodesValue{
				Id:    convert.Int64ToType(n.Id),
				Name:  convert.StrToType(n.Name),
				state: attr.ValueStateKnown,
			}

			objVal, diags := nv.ToObjectValue(ctx)
			if !diags.HasError() {
				nodeVals = append(nodeVals, objVal)
			}
		}

		setVal, diags := types.SetValue(
			NodesType{types.ObjectType{AttrTypes: NodesValue{}.AttributeTypes(ctx)}},
			nodeVals,
		)
		if !diags.HasError() {
			state.Nodes = setVal
		} else {
			state.Nodes = types.SetNull(
				NodesType{types.ObjectType{AttrTypes: NodesValue{}.AttributeTypes(ctx)}},
			)
		}
	} else {
		state.Nodes = types.SetNull(
			NodesType{types.ObjectType{AttrTypes: NodesValue{}.AttributeTypes(ctx)}},
		)
	}

	// Members — empty nested set (schema has no attributes inside MembersValue).
	if len(p.Members) > 0 {
		memberVals := make([]attr.Value, 0, len(p.Members))
		for range p.Members {
			mv := MembersValue{state: attr.ValueStateKnown}

			objVal, diags := mv.ToObjectValue(ctx)
			if !diags.HasError() {
				memberVals = append(memberVals, objVal)
			}
		}

		setVal, diags := types.SetValue(
			MembersType{types.ObjectType{AttrTypes: MembersValue{}.AttributeTypes(ctx)}},
			memberVals,
		)
		if !diags.HasError() {
			state.Members = setVal
		} else {
			state.Members = types.SetNull(
				MembersType{types.ObjectType{AttrTypes: MembersValue{}.AttributeTypes(ctx)}},
			)
		}
	} else {
		state.Members = types.SetNull(
			MembersType{types.ObjectType{AttrTypes: MembersValue{}.AttributeTypes(ctx)}},
		)
	}

	// Config — NSX-T typed vs generic dynamic.
	lbTypeCode := ""
	if p.LoadBalancer != nil {
		if lbType, ok := p.LoadBalancer.GetTypeOk(); ok && lbType != nil {
			if code, ok := lbType.GetCodeOk(); ok && code != nil {
				lbTypeCode = *code
			}
		}
	}

	configMap := p.GetConfig()

	switch lbTypeCode {
	case "nsx-t":
		state.Config = types.DynamicNull()

		if configMap != nil {
			activeMonitorPaths := types.Int64Null()
			if v, ok := configMap["activeMonitorPaths"].(float64); ok {
				activeMonitorPaths = types.Int64Value(int64(v))
			}

			memberGroupVal := NewMemberGroupValueNull()
			if mgMap, ok := configMap["memberGroup"].(map[string]interface{}); ok {
				mgPath := types.StringNull()
				if v, ok := mgMap["path"].(string); ok {
					mgPath = types.StringValue(v)
				}

				mgIpRevisionFilter := types.StringNull()
				if v, ok := mgMap["ipRevisionFilter"].(string); ok {
					mgIpRevisionFilter = types.StringValue(v)
				}

				mgMaxIpListSize := types.Int64Null()
				if v, ok := mgMap["maxIpListSize"].(float64); ok {
					mgMaxIpListSize = types.Int64Value(int64(v))
				}

				mgPort := types.Int64Null()
				if v, ok := mgMap["port"].(float64); ok {
					mgPort = types.Int64Value(int64(v))
				}

				mgVal, d := NewMemberGroupValue(
					MemberGroupValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"ip_revision_filter": mgIpRevisionFilter,
						"max_ip_list_size":   mgMaxIpListSize,
						"path":               mgPath,
						"port":               mgPort,
					},
				)
				if !d.HasError() {
					memberGroupVal = mgVal
				}
			}

			memberGroupObjVal, d := memberGroupVal.ToObjectValue(ctx)
			if d.HasError() {
				memberGroupObjVal = types.ObjectNull(MemberGroupValue{}.AttributeTypes(ctx))
			}

			passiveMonitorPath := types.Int64Null()
			if v, ok := configMap["passiveMonitorPath"].(float64); ok {
				passiveMonitorPath = types.Int64Value(int64(v))
			}

			snatIpAddresses := types.ListNull(types.StringType)
			if v, ok := configMap["snatIpAddresses"].([]interface{}); ok {
				strVals := make([]attr.Value, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						strVals = append(strVals, types.StringValue(s))
					}
				}

				listVal, d := types.ListValue(types.StringType, strVals)
				if !d.HasError() {
					snatIpAddresses = listVal
				}
			}

			snatTranslationType := types.StringNull()
			if v, ok := configMap["snatTranslationType"].(string); ok {
				snatTranslationType = types.StringValue(v)
			}

			tcpMultiplexing := types.BoolNull()
			if v, ok := configMap["tcpMultiplexing"].(bool); ok {
				tcpMultiplexing = types.BoolValue(v)
			}

			tcpMultiplexingNumber := types.Int64Null()
			if v, ok := configMap["tcpMultiplexingNumber"].(float64); ok {
				tcpMultiplexingNumber = types.Int64Value(int64(v))
			}

			nsxtVal, d := NewConfigNsxtValue(
				ConfigNsxtValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"active_monitor_paths":            activeMonitorPaths,
					"member_group":          memberGroupObjVal,
					"passive_monitor_path":            passiveMonitorPath,
					"snat_ip_addresses":               snatIpAddresses,
					"snat_translation_type":           snatTranslationType,
					"tcp_multiplexing":                tcpMultiplexing,
					"tcp_multiplexing_number":         tcpMultiplexingNumber,
				},
			)
			if !d.HasError() {
				state.ConfigNsxt = nsxtVal
			} else {
				state.ConfigNsxt = NewConfigNsxtValueNull()
			}
		} else {
			state.ConfigNsxt = NewConfigNsxtValueNull()
		}

	default:
		state.ConfigNsxt = NewConfigNsxtValueNull()

		if configMap != nil {
			dyn, err := convert.MapToDynamic(ctx, configMap)
			if err == nil {
				state.Config = dyn
			} else {
				state.Config = types.DynamicNull()
			}
		} else {
			state.Config = types.DynamicNull()
		}
	}

	return state
}
