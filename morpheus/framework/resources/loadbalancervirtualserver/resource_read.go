// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var current LoadBalancerVirtualServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &current)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("read load balancer virtual server", "failed to create client: "+err.Error())

		return
	}

	lbID, err := loadBalancerIDFromInt64(current.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError("read load balancer virtual server", err.Error())

		return
	}

	id := current.Id.ValueInt64()

	state, lbTypeCode, configMap, diags := getVirtualServerAsState(ctx, lbID, id, client)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// If getVirtualServerAsState returns a zero-value ID, the resource was not found.
	if state.Id.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	// Config handling: on import the current state has no config (both are null)
	// so we must build it from the API response. On a normal read the current
	// state already has the correct config, so we preserve it — the API returns
	// config as a generic map that may not round-trip perfectly.
	switch lbTypeCode {
	case "nsx-t":
		if current.ConfigNsxt.IsNull() {
			if err := setConfigFromResponse(ctx, &state, lbTypeCode, configMap); err != nil {
				resp.Diagnostics.AddError("import load balancer virtual server", err.Error())

				return
			}
		} else {
			state.ConfigNsxt = current.ConfigNsxt
		}
	default:
		if current.Config.IsNull() {
			if err := setConfigFromResponse(ctx, &state, lbTypeCode, configMap); err != nil {
				resp.Diagnostics.AddError("import load balancer virtual server", err.Error())

				return
			}
		} else {
			state.Config = current.Config
		}
	}

	// VipPool is write-only; always preserve from current state.
	state.VipPool = current.VipPool

	// Pool ID: getVirtualServerAsState already tried the top-level pool object.
	// If still null, try extracting from the config map (NSX-T stores pool there).
	if state.PoolId.IsNull() && configMap != nil {
		if v, ok := configMap["pool"].(string); ok && v != "" {
			if pid, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
				state.PoolId = types.Int64Value(pid)
			}
		}
	}

	// Last resort: preserve pool_id from current state (e.g. API returned neither).
	if state.PoolId.IsNull() && !current.PoolId.IsNull() {
		state.PoolId = current.PoolId
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getVirtualServerAsState(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (LoadBalancerVirtualServerModel, string, map[string]interface{}, diag.Diagnostics) {
	var state LoadBalancerVirtualServerModel
	var diags diag.Diagnostics

	resp, hresp, err := client.LoadBalancersAPI.
		GetLoadBalancerVirtualServer(ctx, loadBalancerID, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			// Signal not-found with a null ID.
			state.Id = types.Int64Null()

			return state, "", nil, diags
		}

		diags.AddError(
			"error reading load balancer virtual server",
			fmt.Sprintf("load balancer %d virtual server %d GET failed: %s",
				loadBalancerID, id, errfmt.ErrMsg(err, hresp)),
		)

		return state, "", nil, diags
	}

	vs := resp.GetLoadBalancerInstance()

	// Extract the load balancer type code for config discrimination.
	lbTypeCode := ""
	if vs.LoadBalancer != nil {
		if lbType, ok := vs.LoadBalancer.GetTypeOk(); ok && lbType != nil {
			if code, ok := lbType.GetCodeOk(); ok && code != nil {
				lbTypeCode = *code
			}
		}
	}

	state.Id = convert.Int64ToType(vs.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Description = convert.StrToType(vs.Description.Get())
	state.VipName = convert.StrToType(vs.VipName)
	state.VipAddress = convert.StrToType(vs.VipAddress)
	state.VipPort = convert.Int64ToType(vs.VipPort)
	state.VipProtocol = convert.StrToType(vs.VipProtocol)
	state.VipHostname = convert.StrToType(vs.VipHostname.Get())

	// Computed fields
	state.Active = convert.BoolToType(vs.Active)
	state.BackendPort = convert.StrToType(vs.BackendPort.Get())
	state.ExternalAddress = convert.BoolToType(vs.ExternalAddress)
	state.ExternalId = convert.StrToType(vs.ExternalId)
	state.ExternalPortId = convert.StrToType(vs.ExternalPortId.Get())
	state.ExtraConfig = convert.StrToType(vs.ExtraConfig.Get())
	state.Instance = convert.StrToType(vs.Instance.Get())
	state.InternalId = convert.StrToType(vs.InternalId)
	state.NetworkId = convert.StrToType(vs.NetworkId.Get())
	state.PoolName = convert.StrToType(vs.PoolName.Get())
	state.Removing = convert.BoolToType(vs.Removing)
	state.ServerName = convert.StrToType(vs.ServerName.Get())
	state.ServiceAccess = convert.StrToType(vs.ServiceAccess.Get())
	state.ServicePort = convert.StrToType(vs.ServicePort.Get())
	state.SourceAddress = convert.StrToType(vs.SourceAddress.Get())
	state.SslEnabled = convert.StrToType(vs.SslEnabled.Get())
	state.SslMode = convert.StrToType(vs.SslMode.Get())
	state.SslRedirectMode = convert.StrToType(vs.SslRedirectMode.Get())
	state.Status = convert.StrToType(vs.Status)
	state.Sticky = convert.BoolToType(vs.Sticky)
	state.SubnetId = convert.StrToType(vs.SubnetId.Get())
	state.VipBalance = convert.StrToType(vs.VipBalance.Get())
	state.VipDirectAddress = convert.StrToType(vs.VipDirectAddress.Get())
	state.VipMode = convert.StrToType(vs.VipMode.Get())
	state.VipScheme = convert.StrToType(vs.VipScheme.Get())
	state.VipShared = convert.BoolToType(vs.VipShared)
	state.VipSource = convert.StrToType(vs.VipSource)
	state.VipStatus = convert.StrToType(vs.VipStatus)
	state.VipSticky = convert.StrToType(vs.VipSticky.Get())
	state.VipType = convert.StrToType(vs.VipType.Get())

	// Dates
	if dc, ok := vs.GetDateCreatedOk(); ok && dc != nil {
		state.DateCreated = types.StringValue(dc.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if lu, ok := vs.GetLastUpdatedOk(); ok && lu != nil {
		state.LastUpdated = types.StringValue(lu.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// SSL cert — GET returns an object {id, name}; schema expects Int64.
	if sslCert, ok := vs.GetSslCertOk(); ok && sslCert != nil {
		state.SslCert = convert.Int64ToType(sslCert.Id)
	} else {
		state.SslCert = types.Int64Null()
	}

	// SSL server cert — GET returns an object {id, name}; schema expects Int64.
	if sslServerCert, ok := vs.GetSslServerCertOk(); ok && sslServerCert != nil {
		state.SslServerCert = convert.Int64ToType(sslServerCert.Id)
	} else {
		state.SslServerCert = types.Int64Null()
	}

	// Pool — GET returns a nested object {id, name}; extract the ID for pool_id.
	if vs.Pool != nil && vs.Pool.Id != nil {
		state.PoolId = types.Int64Value(*vs.Pool.Id)
	} else {
		state.PoolId = types.Int64Null()
	}

	// Nested load_balancer object
	lb, d := buildLoadBalancerValue(ctx, vs.LoadBalancer)
	if diags.Append(d...); diags.HasError() {
		return state, lbTypeCode, nil, diags
	}

	state.LoadBalancer = lb

	return state, lbTypeCode, vs.GetConfig(), diags
}

func setConfigFromResponse(
	ctx context.Context,
	state *LoadBalancerVirtualServerModel,
	lbTypeCode string,
	configMap map[string]interface{},
) error {
	switch lbTypeCode {
	case "nsx-t":
		if configMap == nil {
			state.ConfigNsxt = NewConfigNsxtValueNull()

			return nil
		}

		appProfile := types.Int64Null()
		if v, ok := configMap["applicationProfile"].(float64); ok {
			appProfile = types.Int64Value(int64(v))
		}

		persistence := types.StringNull()
		if v, ok := configMap["persistence"].(string); ok {
			persistence = types.StringValue(v)
		}

		persistenceProfile := types.Int64Null()
		if v, ok := configMap["persistenceProfile"].(float64); ok {
			persistenceProfile = types.Int64Value(int64(v))
		}

		sslClientProfile := types.Int64Null()
		if v, ok := configMap["sslClientProfile"].(float64); ok {
			sslClientProfile = types.Int64Value(int64(v))
		}

		sslServerProfile := types.Int64Null()
		if v, ok := configMap["sslServerProfile"].(float64); ok {
			sslServerProfile = types.Int64Value(int64(v))
		}

		nsxtVal, d := NewConfigNsxtValue(
			ConfigNsxtValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"application_profile": appProfile,
				"persistence":         persistence,
				"persistence_profile": persistenceProfile,
				"ssl_client_profile":  sslClientProfile,
				"ssl_server_profile":  sslServerProfile,
			},
		)
		if d.HasError() {
			return fmt.Errorf("failed to build config_nsxt from API response")
		}

		state.ConfigNsxt = nsxtVal
	default:
		if configMap == nil {
			state.Config = types.DynamicNull()

			return nil
		}

		dyn, err := convert.MapToDynamic(ctx, configMap)
		if err != nil {
			return fmt.Errorf("failed to convert config from API response: %w", err)
		}

		state.Config = dyn
	}

	return nil
}

func buildLoadBalancerValue(
	ctx context.Context,
	lb *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstanceLoadBalancer,
) (LoadBalancerValue, diag.Diagnostics) {
	if lb == nil {
		return NewLoadBalancerValueNull(), nil
	}

	// Build the nested type object
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
