// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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
	monitor.Name = plan.Name.ValueStringPointer()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		monitor.Description.Set(plan.Description.ValueStringPointer())
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
		monitor.MonitorInterval = plan.MonitorInterval.ValueInt64Pointer()
	}

	if !plan.MonitorTimeout.IsNull() && !plan.MonitorTimeout.IsUnknown() {
		monitor.MonitorTimeout = plan.MonitorTimeout.ValueInt64Pointer()
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
		monitor.FallCount = plan.FallCount.ValueInt64Pointer()
	}

	if !plan.RiseCount.IsNull() && !plan.RiseCount.IsUnknown() {
		monitor.RiseCount = plan.RiseCount.ValueInt64Pointer()
	}

	if !plan.AliasPort.IsNull() && !plan.AliasPort.IsUnknown() {
		monitor.AliasPort = plan.AliasPort.ValueInt64Pointer()
	}

	if !plan.DataLength.IsNull() && !plan.DataLength.IsUnknown() {
		monitor.DataLength = plan.DataLength.ValueInt64Pointer()
	}

	if !plan.MaxRetry.IsNull() && !plan.MaxRetry.IsUnknown() {
		monitor.MaxRetry = plan.MaxRetry.ValueInt64Pointer()
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
		sdkConfig.HttpMethod, sdkConfig.HttpUri, sdkConfig.HttpReturnCode, sdkConfig.HttpPort = httpHealthCheckConfig(plan.Config)

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
	if monitor.Config != nil {
		createReq.AdditionalProperties = topLevelConfig(monitor.Config.ToMap())
	}

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

	// Store the configured-value digests now the create POST succeeded, before
	// the read-back's early returns -- a read-back failure must not discard
	// them and make the next apply see a spurious change.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, apiObj, found, diags := getLoadBalancerMonitorAsState(
		ctx, loadBalancerID, *created.Id, client,
	)
	if !found && !diags.HasError() {
		diags.AddError("read load balancer monitor", "monitor not found immediately after create")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "loadbalancermonitor",
			ResourceID:   *created.Id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// Preserve the canonical monitor_type from plan (the API returns
	// the LB-specific value, e.g. LBHttpMonitorProfile for NSX-T).
	state.MonitorType = plan.MonitorType

	// Preserve write-only / sensitive values from plan.
	state.ExtraConfig = plan.ExtraConfig
	state.MonitorPasswordWoVersion = plan.MonitorPasswordWoVersion
	state.WriteOnlyRevision = wodigest.NewRevision()
	// Record the appliance-side drift baseline now the read-back succeeded;
	// the configured-value digests were already stored above after the POST.
	resp.Diagnostics.Append(wo.Observe(ctx, resp.Private, apiObj)...)

	// Keep the planned config values the API may not echo verbatim, but take
	// the read-back for anything the plan left unknown (monitor,
	// monitor_config), so no unknown survives into state.
	state.Config = overlayPlanConfig(ctx, plan.Config, state.Config)

	// Restore fields the API accepts but does not return.
	preserveUnreturnedFields(&state, plan)

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// httpHealthCheckConfig returns the HTTP health-check fields of the config
// block — the SDN load balancer type's monitor option types, kept in the
// monitor's config map — as the nullable SDK values the requests take. A
// field that is not set is left unsent (the zero NullableX).
// topLevelConfig wraps the monitor's config block as the request's top-level
// `config` key. The appliance builds the monitor's stored config from the
// request's top-level `config` (one entry per option type of the load balancer
// type with fieldContext "config"), not from `loadBalancerMonitor.config`,
// which it ignores. The block is therefore sent at both places: under the
// monitor, where the API contract documents it, and at the top level, where
// the appliance reads it. Without the top-level copy the SDN plugin's http*
// option types are stored empty, GET returns `config: {}` and an import has
// no `config.http_*` to read back (MORPH-21017).
func topLevelConfig(cfg map[string]interface{}, err error) map[string]interface{} {
	if err != nil || len(cfg) == 0 {
		return nil
	}

	return map[string]interface{}{"config": cfg}
}

func httpHealthCheckConfig(cfg ConfigValue) (method, uri, code sdk.NullableString, port sdk.NullableInt64) {
	if !cfg.HttpMethod.IsNull() && !cfg.HttpMethod.IsUnknown() {
		method.Set(cfg.HttpMethod.ValueStringPointer())
	}

	if !cfg.HttpUri.IsNull() && !cfg.HttpUri.IsUnknown() {
		uri.Set(cfg.HttpUri.ValueStringPointer())
	}

	if !cfg.HttpReturnCode.IsNull() && !cfg.HttpReturnCode.IsUnknown() {
		code.Set(cfg.HttpReturnCode.ValueStringPointer())
	}

	if !cfg.HttpPort.IsNull() && !cfg.HttpPort.IsUnknown() {
		port.Set(cfg.HttpPort.ValueInt64Pointer())
	}

	return method, uri, code, port
}

// httpHealthCheckUpdateConfig is httpHealthCheckConfig for an update: a field
// the practitioner has removed from the configuration (null in the plan, a
// value in the prior state) is sent as an explicit null so the monitor's
// recorded value is cleared, rather than left to whatever the appliance does
// with an omitted key. The attributes are Optional without Computed, so the
// framework plans the null itself and no clearable record is needed.
//
// What the SDN controller then does is the plugin's: `http_uri` falls back to
// `/` and `http_port` is dropped from the health-check URL, while
// `http_method` and `http_return_code` keep their current value on the
// controller (the health check cannot exist without them) even though the
// Morpheus record is cleared.
func httpHealthCheckUpdateConfig(
	plan, prior ConfigValue,
) (method, uri, code sdk.NullableString, port sdk.NullableInt64) {
	method, uri, code, port = httpHealthCheckConfig(plan)

	if prior.IsNull() || prior.IsUnknown() {
		return method, uri, code, port
	}

	if plan.HttpMethod.IsNull() && !prior.HttpMethod.IsNull() {
		method.Set(nil)
	}

	if plan.HttpUri.IsNull() && !prior.HttpUri.IsNull() {
		uri.Set(nil)
	}

	if plan.HttpReturnCode.IsNull() && !prior.HttpReturnCode.IsNull() {
		code.Set(nil)
	}

	if plan.HttpPort.IsNull() && !prior.HttpPort.IsNull() {
		port.Set(nil)
	}

	return method, uri, code, port
}

// httpHealthCheckClearConfig is the update body for a config block removed
// as a whole: every HTTP health-check value the prior state carried is sent
// as an explicit null, exactly as removing the fields one at a time would.
// It returns nil when the prior block held none of them, so no config is sent.
func httpHealthCheckClearConfig(prior ConfigValue) *sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig {
	if prior.IsNull() || prior.IsUnknown() {
		return nil
	}

	cfg := &sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	cleared := false

	if !prior.HttpMethod.IsNull() {
		cfg.HttpMethod.Set(nil)
		cleared = true
	}

	if !prior.HttpUri.IsNull() {
		cfg.HttpUri.Set(nil)
		cleared = true
	}

	if !prior.HttpReturnCode.IsNull() {
		cfg.HttpReturnCode.Set(nil)
		cleared = true
	}

	if !prior.HttpPort.IsNull() {
		cfg.HttpPort.Set(nil)
		cleared = true
	}

	if !cleared {
		return nil
	}

	return cfg
}
