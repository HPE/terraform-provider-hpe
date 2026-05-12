// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState, config LoadBalancerMonitorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	monitor := sdk.NewUpdateLoadBalancerMonitorRequestLoadBalancerMonitorWithDefaults()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		monitor.SetName(plan.Name.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		monitor.SetDescription(plan.Description.ValueString())
	}

	loadBalancerID := currentState.LoadBalancerId.ValueInt64()

	if !plan.MonitorType.IsNull() && !plan.MonitorType.IsUnknown() {
		monitorType := plan.MonitorType.ValueString()

		lbResp, httpResp, err := client.LoadBalancersAPI.
			GetLoadBalancer(ctx, loadBalancerID).Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				"error reading load balancer",
				fmt.Sprintf("load balancer %d GET failed: %s", loadBalancerID,
					errfmt.ErrMsg(err, httpResp)),
			)

			return
		}

		lbTypeCode := ""
		if lb := lbResp.GetLoadBalancer(); lb.Type != nil {
			if code := lb.Type.Code; code != nil {
				lbTypeCode = *code
			}
		}

		switch lbTypeCode {
		case LBTypeNsxT:
			if mapped, ok := nsxtMonitorTypes[monitorType]; ok {
				monitorType = mapped
			} else {
				resp.Diagnostics.AddError(
					"invalid monitor_type for NSX-T load balancer",
					fmt.Sprintf("monitor_type %q is not valid for NSX-T; valid values: http, https, icmp, passive, tcp, udp",
						plan.MonitorType.ValueString()),
				)

				return
			}
		case LBTypeNsxV:
			if _, ok := nsxvMonitorTypes[monitorType]; !ok {
				resp.Diagnostics.AddError(
					"invalid monitor_type for NSX-V load balancer",
					fmt.Sprintf("monitor_type %q is not valid for NSX-V; valid values: dns, http, https, ldap, mssql, tcp, udp",
						plan.MonitorType.ValueString()),
				)

				return
			}
		}

		monitor.SetMonitorType(monitorType)
	}

	if !plan.MonitorInterval.IsNull() && !plan.MonitorInterval.IsUnknown() {
		monitor.SetMonitorInterval(plan.MonitorInterval.ValueInt64())
	}

	if !plan.MonitorTimeout.IsNull() && !plan.MonitorTimeout.IsUnknown() {
		monitor.SetMonitorTimeout(plan.MonitorTimeout.ValueInt64())
	}

	if !plan.SendData.IsNull() && !plan.SendData.IsUnknown() {
		monitor.SetSendData(plan.SendData.ValueString())
	}

	if !plan.SendVersion.IsNull() && !plan.SendVersion.IsUnknown() {
		monitor.SetSendVersion(plan.SendVersion.ValueString())
	}

	if !plan.SendType.IsNull() && !plan.SendType.IsUnknown() {
		monitor.SetSendType(plan.SendType.ValueString())
	}

	if !plan.ReceiveData.IsNull() && !plan.ReceiveData.IsUnknown() {
		monitor.SetReceiveData(plan.ReceiveData.ValueString())
	}

	if !plan.ReceiveCode.IsNull() && !plan.ReceiveCode.IsUnknown() {
		monitor.SetReceiveCode(plan.ReceiveCode.ValueString())
	}

	if !plan.MonitorUsername.IsNull() && !plan.MonitorUsername.IsUnknown() {
		monitor.SetMonitorUsername(plan.MonitorUsername.ValueString())
	}

	if !plan.MonitorPasswordWoVersion.Equal(currentState.MonitorPasswordWoVersion) {
		if config.MonitorPasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update load balancer monitor resource",
				"'monitor_password_wo_version' changed, but 'monitor_password_wo' is not set",
			)

			return
		}

		monitor.SetMonitorPassword(config.MonitorPasswordWo.ValueString())
	}

	if !plan.MonitorDestination.IsNull() && !plan.MonitorDestination.IsUnknown() {
		monitor.SetMonitorDestination(plan.MonitorDestination.ValueString())
	}

	if !plan.FallCount.IsNull() && !plan.FallCount.IsUnknown() {
		monitor.SetFallCount(plan.FallCount.ValueInt64())
	}

	if !plan.RiseCount.IsNull() && !plan.RiseCount.IsUnknown() {
		monitor.SetRiseCount(plan.RiseCount.ValueInt64())
	}

	if !plan.AliasPort.IsNull() && !plan.AliasPort.IsUnknown() {
		monitor.SetAliasPort(plan.AliasPort.ValueInt64())
	}

	if !plan.DataLength.IsNull() && !plan.DataLength.IsUnknown() {
		monitor.SetDataLength(plan.DataLength.ValueInt64())
	}

	if !plan.MaxRetry.IsNull() && !plan.MaxRetry.IsUnknown() {
		monitor.SetMaxRetry(plan.MaxRetry.ValueInt64())
	}

	if !plan.ExtraConfig.IsNull() && !plan.ExtraConfig.IsUnknown() {
		monitor.SetExtraConfig(plan.ExtraConfig.ValueString())
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		sdkConfig := sdk.NewUpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfigWithDefaults()

		if !plan.Config.Monitor.IsNull() && !plan.Config.Monitor.IsUnknown() {
			attrs := plan.Config.Monitor.Attributes()
			if idAttr, ok := attrs["id"]; ok {
				if idVal, ok := idAttr.(basetypes.Int64Value); ok && !idVal.IsNull() && !idVal.IsUnknown() {
					configMonitor := sdk.NewUpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfigMonitorWithDefaults()
					configMonitor.SetId(idVal.ValueInt64())
					sdkConfig.SetMonitor(*configMonitor)
				}
			}
		}

		if !plan.Config.MonitorConfig.IsNull() && !plan.Config.MonitorConfig.IsUnknown() {
			mcAny, err := convert.ValueToAny(ctx, plan.Config.MonitorConfig.UnderlyingValue())
			if err != nil {
				resp.Diagnostics.AddError(
					"error converting monitor_config",
					fmt.Sprintf("failed to convert monitor_config: %s", err),
				)

				return
			}

			if s, ok := mcAny.(string); ok {
				sdkConfig.SetMonitorConfig(s)
			}
		}

		monitor.SetConfig(*sdkConfig)
	}

	updateReq := sdk.NewUpdateLoadBalancerMonitorRequestWithDefaults()
	updateReq.SetLoadBalancerMonitor(*monitor)

	id := currentState.Id.ValueInt64()

	_, httpResp, err := client.LoadBalancersAPI.
		UpdateLoadBalancerMonitor(ctx, loadBalancerID, id).
		UpdateLoadBalancerMonitorRequest(*updateReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating load balancer monitor",
			"load balancer monitor "+plan.Name.ValueString()+" PUT failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, diags := getLoadBalancerMonitorAsState(
		ctx, loadBalancerID, id, client,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve write-only / sensitive values from plan.
	state.ExtraConfig = plan.ExtraConfig
	state.MonitorPasswordWoVersion = plan.MonitorPasswordWoVersion

	// Preserve plan config since the API may not return it verbatim.
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
