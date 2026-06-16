// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data LoadBalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diags := getLoadBalancerPoolAsState(
		ctx, data.LoadBalancerId.ValueInt64(), data.Id.ValueInt64(), client,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Partition is write-only; preserve from prior state.
	state.Partition = data.Partition

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getLoadBalancerPoolAsState(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (LoadBalancerPoolModel, diag.Diagnostics) {
	var state LoadBalancerPoolModel
	var diags diag.Diagnostics

	poolResp, httpResp, err := client.LoadBalancersAPI.
		GetLoadBalancerPool(ctx, loadBalancerID, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading load balancer pool",
			fmt.Sprintf("load balancer pool %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	p := poolResp.LoadBalancerPool
	if p == nil {
		diags.AddError(
			"error reading load balancer pool",
			fmt.Sprintf("load balancer pool %d GET returned no pool", id),
		)

		return state, diags
	}

	// Core fields
	state.Id = convert.Int64ToType(p.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(p.Name)

	// NullableString fields
	if p.Description.IsSet() {
		state.Description = convert.StrToType(p.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

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

	if p.Port.IsSet() && p.Port.Get() != nil {
		portVal, parseErr := strconv.ParseInt(*p.Port.Get(), 10, 64)
		if parseErr == nil {
			state.Port = types.Int64Value(portVal)
		} else {
			state.Port = types.Int64Null()
		}
	} else {
		state.Port = types.Int64Null()
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

	if p.VipClientIpMode.IsSet() {
		state.VipClientIpMode = convert.StrToType(p.VipClientIpMode.Get())
	} else {
		state.VipClientIpMode = types.StringNull()
	}

	if p.VipServerIpMode.IsSet() {
		state.VipServerIpMode = convert.StrToType(p.VipServerIpMode.Get())
	} else {
		state.VipServerIpMode = types.StringNull()
	}

	if p.VipSticky.IsSet() {
		state.VipSticky = convert.StrToType(p.VipSticky.Get())
	} else {
		state.VipSticky = types.StringNull()
	}

	// *string fields
	state.InternalId = convert.StrToType(p.InternalId)
	state.ExternalId = convert.StrToType(p.ExternalId)
	state.Status = convert.StrToType(p.Status)
	state.Visibility = convert.StrToType(p.Visibility)
	state.VipBalance = convert.StrToType(p.VipBalance)

	// *int64 fields
	state.MinActive = convert.Int64ToType(p.MinActive)
	state.NumberActive = convert.Int64ToType(p.NumberActive)
	state.NumberInService = convert.Int64ToType(p.NumberInService)
	state.HealthScore = convert.Int64ToType(p.HealthScore)
	state.PerformanceScore = convert.Int64ToType(p.PerformanceScore)
	state.HealthPenalty = convert.Int64ToType(p.HealthPenalty)
	state.SecurityPenalty = convert.Int64ToType(p.SecurityPenalty)
	state.ErrorPenalty = convert.Int64ToType(p.ErrorPenalty)

	// *bool fields
	state.Enabled = convert.BoolToType(p.Enabled)

	// *time.Time fields
	if p.DateCreated != nil {
		state.DateCreated = types.StringValue(p.DateCreated.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if p.LastUpdated != nil {
		state.LastUpdated = types.StringValue(p.LastUpdated.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Nested LoadBalancer object
	lb, lbDiags := buildLoadBalancerValue(ctx, p.LoadBalancer)
	if diags.Append(lbDiags...); diags.HasError() {
		return state, diags
	}

	state.LoadBalancer = lb

	// Partition is write-only; not returned by the API.
	state.Partition = types.StringNull()

	// Config: parse the config map from the API response into the appropriate
	// typed block (config_nsxt) or dynamic fallback (config).
	configMap := p.Config

	lbTypeCode := ""
	if p.LoadBalancer != nil {
		if lbType := p.LoadBalancer.Type; lbType != nil {
			if code := lbType.Code; code != nil {
				lbTypeCode = *code
			}
		}
	}

	switch lbTypeCode {
	case "nsx-t":
		state.Config = types.DynamicNull()
		state.ConfigNsxt = parseConfigNsxt(ctx, configMap)

	default:
		state.ConfigNsxt = NewConfigNsxtValueNull()

		if len(configMap) > 0 {
			dynVal, err := convert.MapToDynamic(ctx, configMap)
			if err != nil {
				diags.AddError(
					"error reading load balancer pool config",
					fmt.Sprintf("failed to convert config map to dynamic: %s", err),
				)

				return state, diags
			}

			state.Config = dynVal
		} else {
			state.Config = types.DynamicNull()
		}
	}

	return state, diags
}

func buildLoadBalancerValue(
	ctx context.Context,
	lb *sdk.GetLoadBalancerPool200ResponseLoadBalancerPoolLoadBalancer,
) (LoadBalancerValue, diag.Diagnostics) {
	if lb == nil {
		return NewLoadBalancerValueNull(), nil
	}

	// Build the nested type object
	typeVal := types.ObjectNull(TypeValue{}.AttributeTypes(ctx))

	if lbType := lb.Type; lbType != nil {
		tv, d := NewTypeValue(
			TypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"code": convert.StrToType(lbType.Code),
				"id":   convert.Int64ToType(lbType.Id),
				"name": convert.StrToType(lbType.Name),
			},
		)
		if d.HasError() {
			return NewLoadBalancerValueNull(), d
		}

		var tvDiags diag.Diagnostics
		typeVal, tvDiags = tv.ToObjectValue(ctx)
		if tvDiags.HasError() {
			return NewLoadBalancerValueNull(), tvDiags
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
	if d.HasError() {
		return NewLoadBalancerValueNull(), d
	}

	return lbVal, nil
}

func parseConfigNsxt(ctx context.Context, configMap map[string]interface{}) ConfigNsxtValue {
	if configMap == nil {
		return NewConfigNsxtValueNull()
	}

	activeMonitorPaths := types.Int64Null()
	if v, ok := configMap["activeMonitorPaths"].(float64); ok {
		activeMonitorPaths = types.Int64Value(int64(v))
	}

	passiveMonitorPath := types.Int64Null()
	if v, ok := configMap["passiveMonitorPath"].(float64); ok {
		passiveMonitorPath = types.Int64Value(int64(v))
	}

	snatTranslationType := types.StringNull()
	if v, ok := configMap["snatTranslationType"].(string); ok {
		snatTranslationType = types.StringValue(v)
	}

	snatIpAddresses := types.SetNull(types.StringType)
	if v, ok := configMap["snatIpAddresses"].([]interface{}); ok {
		strVals := make([]attr.Value, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				strVals = append(strVals, types.StringValue(s))
			}
		}

		setVal, d := types.SetValue(types.StringType, strVals)
		if !d.HasError() {
			snatIpAddresses = setVal
		}
	}

	tcpMultiplexing := types.BoolNull()
	if v, ok := configMap["tcpMultiplexing"].(bool); ok {
		tcpMultiplexing = types.BoolValue(v)
	}

	tcpMultiplexingNumber := types.Int64Null()
	if v, ok := configMap["tcpMultiplexingNumber"].(float64); ok {
		tcpMultiplexingNumber = types.Int64Value(int64(v))
	}

	// Build the nested member_group object.
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

	nsxtVal, d := NewConfigNsxtValue(
		ConfigNsxtValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"active_monitor_paths":    activeMonitorPaths,
			"member_group":            memberGroupObjVal,
			"passive_monitor_path":    passiveMonitorPath,
			"snat_ip_addresses":       snatIpAddresses,
			"snat_translation_type":   snatTranslationType,
			"tcp_multiplexing":        tcpMultiplexing,
			"tcp_multiplexing_number": tcpMultiplexingNumber,
		},
	)
	if d.HasError() {
		return NewConfigNsxtValueNull()
	}

	return nsxtVal
}
