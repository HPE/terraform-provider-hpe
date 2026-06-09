// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data LoadBalancerMonitorModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diags := getLoadBalancerMonitorAsState(
		ctx, data.LoadBalancerId.ValueInt64(), data.Id.ValueInt64(), client,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// MonitorType: on a normal read, preserve the canonical value from prior
	// state. On import the prior state value is null, so we reverse-map the
	// API value back to canonical form using the load balancer type.
	if !data.MonitorType.IsNull() {
		state.MonitorType = data.MonitorType
	} else if !state.MonitorType.IsNull() {
		lbResp, httpResp, err := client.LoadBalancersAPI.
			GetLoadBalancer(ctx, data.LoadBalancerId.ValueInt64()).Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				"error reading load balancer",
				fmt.Sprintf("load balancer %d GET failed: %s",
					data.LoadBalancerId.ValueInt64(),
					errfmt.ErrMsg(err, httpResp)),
			)

			return
		}

		lbTypeCode := ""
		if lb := lbResp.LoadBalancer; lb != nil && lb.Type != nil {
			if code := lb.Type.Code; code != nil {
				lbTypeCode = *code
			}
		}

		state.MonitorType = types.StringValue(
			canonicalMonitorType(state.MonitorType.ValueString(), lbTypeCode),
		)
	}

	// Preserve write-only values from prior state.
	// ExtraConfig is not returned by the API — on import it will be null.
	state.ExtraConfig = data.ExtraConfig
	state.MonitorPasswordWoVersion = data.MonitorPasswordWoVersion

	// Preserve plan config since the API may not return it verbatim.
	if !data.Config.IsNull() && !data.Config.IsUnknown() {
		state.Config = data.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getLoadBalancerMonitorAsState(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (LoadBalancerMonitorModel, diag.Diagnostics) {
	var state LoadBalancerMonitorModel
	var diags diag.Diagnostics

	monitorResp, httpResp, err := client.LoadBalancersAPI.
		GetLoadBalancerMonitor(ctx, loadBalancerID, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading load balancer monitor",
			fmt.Sprintf("load balancer monitor %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	m := monitorResp.LoadBalancerMonitor
	if m == nil {
		diags.AddError(
			"error reading load balancer monitor",
			fmt.Sprintf("load balancer monitor %d response did not contain a monitor", id),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(m.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(m.Name)
	state.Description = convert.StrToType(m.Description)
	state.MonitorType = convert.StrToType(m.MonitorType)
	state.MonitorInterval = convert.Int64ToType(m.MonitorInterval)
	state.MonitorTimeout = convert.Int64ToType(m.MonitorTimeout)

	// NullableString fields — use .Get() to obtain the *string.
	state.SendData = convert.StrToType(m.SendData.Get())
	state.SendVersion = convert.StrToType(m.SendVersion)
	state.SendType = convert.StrToType(m.SendType)
	state.ReceiveData = convert.StrToType(m.ReceiveData.Get())
	state.ReceiveCode = convert.StrToType(m.ReceiveCode)
	state.MonitorUsername = convert.StrToType(m.MonitorUsername.Get())
	// MonitorPasswordWo is write-only; not read from the API.
	// MonitorPasswordWoVersion is preserved from plan/state by callers.
	state.MonitorDestination = convert.StrToType(m.MonitorDestination)
	state.FallCount = convert.Int64ToType(m.FallCount)
	state.RiseCount = convert.Int64ToType(m.RiseCount)
	state.AliasPort = convert.Int64ToType(m.AliasPort)

	// DataLength is NullableInt64 in the SDK.
	state.DataLength = convert.Int64ToType(m.DataLength.Get())

	state.MaxRetry = convert.Int64ToType(m.MaxRetry)

	// ExtraConfig is write-only (not returned by the GET API).
	// It is set to null here; callers preserve it from plan or prior state.
	state.ExtraConfig = types.StringNull()

	// Config — structured nested object.
	if m.Config != nil {
		cfg := m.Config

		var monitorObjVal basetypes.ObjectValue
		if cfg.Monitor != nil {
			monitorInner := cfg.Monitor
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

	return state, diags
}
