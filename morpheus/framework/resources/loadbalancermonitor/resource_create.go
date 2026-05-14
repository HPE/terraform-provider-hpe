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

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, config LoadBalancerMonitorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	monitor := sdk.NewCreateLoadBalancerMonitorRequestLoadBalancerMonitorWithDefaults()
	monitor.SetName(plan.Name.ValueString())

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		monitor.SetDescription(plan.Description.ValueString())
	}

	loadBalancerID := plan.LoadBalancerId.ValueInt64()

	// Look up the load balancer type so we can map monitor_type to the
	// correct LB-specific value.
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

	if !plan.MonitorType.IsNull() && !plan.MonitorType.IsUnknown() {
		monitorType := plan.MonitorType.ValueString()

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

	if !config.MonitorPasswordWo.IsNull() && !config.MonitorPasswordWo.IsUnknown() {
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

	// TODO: extraConfig is a write-only String field on the domain object
	// (NetworkLoadBalancerMonitor.extraConfig). It is stored encrypted at rest
	// (EncryptedString). The NSX-V seed defines it as a textarea labeled
	// "Extension" visible for monitor types http, https, tcp, udp, dns, mssql,
	// and ldap. The field is NOT returned by the GET API (commented out in the
	// GSON response template). The state value is preserved from the plan on
	// read to avoid drift.
	if !plan.ExtraConfig.IsNull() && !plan.ExtraConfig.IsUnknown() {
		monitor.SetExtraConfig(plan.ExtraConfig.ValueString())
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		sdkConfig := sdk.NewCreateLoadBalancerMonitorRequestLoadBalancerMonitorConfigWithDefaults()

		if !plan.Config.Monitor.IsNull() && !plan.Config.Monitor.IsUnknown() {
			attrs := plan.Config.Monitor.Attributes()
			if idAttr, ok := attrs["id"]; ok {
				if idVal, ok := idAttr.(basetypes.Int64Value); ok && !idVal.IsNull() && !idVal.IsUnknown() {
					configMonitor := sdk.NewCreateLoadBalancerMonitorRequestLoadBalancerMonitorConfigMonitorWithDefaults()
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

	createReq := sdk.NewCreateLoadBalancerMonitorRequestWithDefaults()
	createReq.SetLoadBalancerMonitor(*monitor)

	createResp, httpResp, err := client.LoadBalancersAPI.
		CreateLoadBalancerMonitor(ctx, loadBalancerID).
		CreateLoadBalancerMonitorRequest(*createReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating load balancer monitor",
			"load balancer monitor "+plan.Name.ValueString()+" POST failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	created := createResp.GetLoadBalancerMonitor()
	if created.Id == nil {
		resp.Diagnostics.AddError(
			"error creating load balancer monitor",
			"load balancer monitor "+plan.Name.ValueString()+": id is nil",
		)

		return
	}

	state, diags := getLoadBalancerMonitorAsState(
		ctx, loadBalancerID, *created.Id, client,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve the canonical monitor_type from plan (the API returns
	// the LB-specific value, e.g. LBHttpMonitorProfile for NSX-T).
	state.MonitorType = plan.MonitorType

	// Preserve write-only / sensitive values from plan.
	state.ExtraConfig = plan.ExtraConfig
	state.MonitorPasswordWoVersion = plan.MonitorPasswordWoVersion

	// Preserve plan config since the API may not return it verbatim.
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
