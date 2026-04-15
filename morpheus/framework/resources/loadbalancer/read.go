// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// isNotFound returns true when the HTTP response indicates the resource
// no longer exists on the server (404 Not Found).
func isNotFound(hresp *http.Response) bool {
	return hresp != nil && hresp.StatusCode == http.StatusNotFound
}

func getLoadBalancerAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan LoadBalancerModel,
) (LoadBalancerModel, bool, diag.Diagnostics) {
	var state LoadBalancerModel
	var diags diag.Diagnostics

	lb, hresp, err := client.LoadBalancersAPI.GetLoadBalancer(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		if isNotFound(hresp) {
			return state, true, diags
		}

		diags.AddError(
			"populate load balancer resource",
			fmt.Sprintf("load balancer %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, false, diags
	}

	data := lb.GetLoadBalancer()

	isHAProxy := data.Type != nil && data.Type.Code != nil && *data.Type.Code == typeCodeHAProxy

	state.Id = convert.Int64ToType(data.Id)
	state.Name = convert.StrToType(data.Name)
	state.Description = convert.StrToType(data.Description)
	state.Visibility = convert.StrToType(data.Visibility)

	switch {
	case isHAProxy:
		state.TypeCode = types.StringNull()
	default:
		if data.Type != nil && data.Type.Code != nil {
			state.TypeCode = types.StringValue(*data.Type.Code)
		}
	}

	if data.Cloud != nil && data.Cloud.Id != nil {
		state.CloudId = convert.Int64ToType(data.Cloud.Id)
	} else {
		state.CloudId = types.Int64Null()
	}

	// The API does not return group or network_server_id on load balancers,
	// so these must be preserved from plan/state. After import they will be null.
	state.GroupId = plan.GroupId
	state.NetworkServerId = plan.NetworkServerId

	// Set config based on the load balancer type.
	// For HAProxy LBs, parse the API config map into the typed config_haproxy attribute.
	// For other types, leave both null — the caller preserves config from the plan.
	switch {
	case isHAProxy:
		haproxyCfg, d := parseHAProxyConfig(ctx, data.GetConfig())
		diags.Append(d...)
		if diags.HasError() {
			return state, false, diags
		}

		state.ConfigHaproxy = haproxyCfg
		state.Config = types.DynamicNull()
	default:
		state.Config = types.DynamicNull()
		state.ConfigHaproxy = NewConfigHaproxyValueNull()
	}

	// Tenants
	state.Tenants = types.SetNull(TenantsType{
		ObjectType: types.ObjectType{
			AttrTypes: TenantsValue{}.AttributeTypes(ctx),
		},
	})
	if len(data.Tenants) > 0 {
		var tenantValues []attr.Value
		for _, tenant := range data.Tenants {
			tv, d := NewTenantsValue(
				TenantsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.Int64ToType(tenant.Id),
					"name": convert.StrToType(tenant.Name),
				},
			)
			diags.Append(d...)
			if diags.HasError() {
				return state, false, diags
			}

			tenantValues = append(tenantValues, tv)
		}

		if len(tenantValues) > 0 {
			tenantSet, d := types.SetValue(
				TenantsType{
					ObjectType: types.ObjectType{
						AttrTypes: TenantsValue{}.AttributeTypes(ctx),
					},
				},
				tenantValues,
			)
			diags.Append(d...)
			if diags.HasError() {
				return state, false, diags
			}

			state.Tenants = tenantSet
		}
	}

	// Resource permissions
	resourcePermission, ok := data.GetResourcePermissionOk()
	if ok && resourcePermission != nil {
		perms, d := convertResourcePermissions(ctx, resourcePermission)
		diags.Append(d...)
		if diags.HasError() {
			return state, false, diags
		}

		state.Permissions = perms
	} else {
		state.Permissions = NewPermissionsValueNull()
	}

	return state, false, diags
}

func convertResourcePermissions(
	ctx context.Context,
	resourcePermission *sdk.GetLoadBalancer200ResponseLoadBalancerResourcePermission,
) (PermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupValues []attr.Value
	sites, ok := resourcePermission.GetSitesOk()
	if ok {
		for _, site := range sites {
			if site.Id != nil {
				groupValues = append(
					groupValues, types.Int64Value(*site.Id),
				)
			}
		}
	}

	var groupIDsSet attr.Value
	if len(groupValues) > 0 {
		s, d := types.SetValue(types.Int64Type, groupValues)
		diags.Append(d...)
		if diags.HasError() {
			return PermissionsValue{}, diags
		}

		groupIDsSet = s
	} else {
		groupIDsSet = types.SetNull(types.Int64Type)
	}

	perms, d := NewPermissionsValue(
		PermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"all": types.BoolValue(
				resourcePermission.All != nil &&
					*resourcePermission.All,
			),
			"groups": groupIDsSet,
		},
	)
	diags.Append(d...)

	return perms, diags
}

func parseHAProxyConfig(
	ctx context.Context,
	configMap map[string]interface{},
) (ConfigHaproxyValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var planID int64

	if planRaw, ok := configMap["plan"]; ok {
		if planMap, ok := planRaw.(map[string]interface{}); ok {
			if idRaw, ok := planMap["id"]; ok {
				switch v := idRaw.(type) {
				case float64:
					planID = int64(v)
				case int64:
					planID = v
				}
			}
		}
	}

	var pool string

	if poolRaw, ok := configMap["pool"]; ok {
		if poolMap, ok := poolRaw.(map[string]interface{}); ok {
			if idRaw, ok := poolMap["id"]; ok {
				if v, ok := idRaw.(string); ok {
					pool = v
				}
			}
		}
	}

	cfg, d := NewConfigHaproxyValue(
		ConfigHaproxyValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"plan_id": types.Int64Value(planID),
			"pool":    types.StringValue(pool),
		},
	)
	diags.Append(d...)

	return cfg, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan LoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read load balancer resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, notFound, diags := getLoadBalancerAsState(ctx, id, client, plan)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case !plan.ConfigHaproxy.IsNull() && !plan.ConfigHaproxy.IsUnknown():
		state.ConfigHaproxy = plan.ConfigHaproxy
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
