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

	monitor := &sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitor{}
	val := plan.Name.ValueString()
	monitor.Name = &val

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		val := plan.Description.ValueString()
		monitor.Description = &val
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
	if lb := lbResp.LoadBalancer; lb != nil && lb.Type != nil {
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

		val := monitorType
		monitor.MonitorType = &val
	}

	if !plan.MonitorInterval.IsNull() && !plan.MonitorInterval.IsUnknown() {
		val := plan.MonitorInterval.ValueInt64()
		monitor.MonitorInterval = &val
	}

	if !plan.MonitorTimeout.IsNull() && !plan.MonitorTimeout.IsUnknown() {
		val := plan.MonitorTimeout.ValueInt64()
		monitor.MonitorTimeout = &val
	}

	if !plan.SendData.IsNull() && !plan.SendData.IsUnknown() {
		val := plan.SendData.ValueString()
		monitor.SendData.Set(&val)
	}

	if !plan.SendVersion.IsNull() && !plan.SendVersion.IsUnknown() {
		val := plan.SendVersion.ValueString()
		monitor.SendVersion.Set(&val)
	}

	if !plan.SendType.IsNull() && !plan.SendType.IsUnknown() {
		val := plan.SendType.ValueString()
		monitor.SendType.Set(&val)
	}

	if !plan.ReceiveData.IsNull() && !plan.ReceiveData.IsUnknown() {
		val := plan.ReceiveData.ValueString()
		monitor.ReceiveData.Set(&val)
	}

	if !plan.ReceiveCode.IsNull() && !plan.ReceiveCode.IsUnknown() {
		val := plan.ReceiveCode.ValueString()
		monitor.ReceiveCode.Set(&val)
	}

	if !plan.MonitorUsername.IsNull() && !plan.MonitorUsername.IsUnknown() {
		val := plan.MonitorUsername.ValueString()
		monitor.MonitorUsername.Set(&val)
	}

	if !config.MonitorPasswordWo.IsNull() && !config.MonitorPasswordWo.IsUnknown() {
		val := config.MonitorPasswordWo.ValueString()
		monitor.MonitorPassword.Set(&val)
	}

	if !plan.MonitorDestination.IsNull() && !plan.MonitorDestination.IsUnknown() {
		val := plan.MonitorDestination.ValueString()
		monitor.MonitorDestination.Set(&val)
	}

	if !plan.FallCount.IsNull() && !plan.FallCount.IsUnknown() {
		val := plan.FallCount.ValueInt64()
		monitor.FallCount = &val
	}

	if !plan.RiseCount.IsNull() && !plan.RiseCount.IsUnknown() {
		val := plan.RiseCount.ValueInt64()
		monitor.RiseCount = &val
	}

	if !plan.AliasPort.IsNull() && !plan.AliasPort.IsUnknown() {
		val := plan.AliasPort.ValueInt64()
		monitor.AliasPort = &val
	}

	if !plan.DataLength.IsNull() && !plan.DataLength.IsUnknown() {
		val := plan.DataLength.ValueInt64()
		monitor.DataLength = &val
	}

	if !plan.MaxRetry.IsNull() && !plan.MaxRetry.IsUnknown() {
		val := plan.MaxRetry.ValueInt64()
		monitor.MaxRetry = &val
	}

	// TODO: extraConfig is a write-only String field on the domain object
	// (NetworkLoadBalancerMonitor.extraConfig). It is stored encrypted at rest
	// (EncryptedString). The NSX-V seed defines it as a textarea labeled
	// "Extension" visible for monitor types http, https, tcp, udp, dns, mssql,
	// and ldap. The field is NOT returned by the GET API (commented out in the
	// GSON response template). The state value is preserved from the plan on
	// read to avoid drift.
	if !plan.ExtraConfig.IsNull() && !plan.ExtraConfig.IsUnknown() {
		val := plan.ExtraConfig.ValueString()
		monitor.ExtraConfig.Set(&val)
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		sdkConfig := &sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}

		if !plan.Config.Monitor.IsNull() && !plan.Config.Monitor.IsUnknown() {
			attrs := plan.Config.Monitor.Attributes()
			if idAttr, ok := attrs["id"]; ok {
				if idVal, ok := idAttr.(basetypes.Int64Value); ok && !idVal.IsNull() && !idVal.IsUnknown() {
					configMonitor := &sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitorConfigMonitor{}
					val := idVal.ValueInt64()
					configMonitor.Id = &val
					sdkConfig.Monitor = configMonitor
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
				sdkConfig.MonitorConfig.Set(&s)
			}
		}

		monitor.Config = sdkConfig
	}

	createReq := &sdk.CreateLoadBalancerMonitorRequest{}
	createReq.LoadBalancerMonitor = monitor

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

	created := createResp.LoadBalancerMonitor
	if created == nil || created.Id == nil {
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
