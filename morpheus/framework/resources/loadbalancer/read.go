// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func getLoadBalancerAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan LoadBalancerModel,
) (LoadBalancerModel, error) {
	var state LoadBalancerModel

	lb, hresp, err := client.LoadBalancersAPI.GetLoadBalancer(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return state, fmt.Errorf(
			"load balancer %d GET failed: %s", id, errfmt.ErrMsg(err, hresp),
		)
	}

	data := lb.GetLoadBalancer()

	if data.Cloud == nil {
		return state, fmt.Errorf("load balancer %d cloud id not found", id)
	}

	if data.Type == nil {
		return state, fmt.Errorf("load balancer %d type not found", id)
	}

	state.Id = convert.Int64ToType(data.Id)
	state.Name = convert.StrToType(data.Name)
	state.Description = convert.StrToType(data.Description)
	state.Visibility = convert.StrToType(data.Visibility)

	// cloud_id is Optional (not Computed), so preserve configured value from plan/state
	// to avoid post-apply inconsistencies when API returns an implicit/default cloud.
	state.CloudId = plan.CloudId

	// The API does not return group or network_server_id on load balancers,
	// so these must be preserved from plan/state. After import they will be null.
	state.GroupId = plan.GroupId
	state.NetworkServerId = plan.NetworkServerId
	state.TypeCode = plan.TypeCode

	// Set config based on the load balancer type.
	// For HAProxy LBs, parse the API config map into the typed config_haproxy attribute.
	// For NSX-T LBs, parse the API config map into the typed config_nsxt attribute.
	// For other types, leave both null — the caller preserves config from the plan.
	isHAProxy := data.Type.Code != nil && *data.Type.Code == typeCodeHAProxy
	isNSXT := data.Type.Code != nil && *data.Type.Code == typeCodeNSXT

	switch {
	case isHAProxy:
		haproxyCfg, err := parseHAProxyConfig(ctx, data.GetConfig())
		if err != nil {
			return state, fmt.Errorf("failed to parse HAProxy config: %w", err)
		}

		state.ConfigHaproxy = haproxyCfg
		state.ConfigNsxt = NewConfigNsxtValueNull()
		state.Config = types.DynamicNull()
	case isNSXT:
		nsxtCfg, err := parseNsxtConfig(ctx, data.GetConfig())
		if err != nil {
			return state, fmt.Errorf("failed to parse NSX-T config: %w", err)
		}

		state.ConfigNsxt = nsxtCfg
		state.ConfigHaproxy = NewConfigHaproxyValueNull()
		state.Config = types.DynamicNull()
	default:
		state.ConfigHaproxy = NewConfigHaproxyValueNull()
		state.ConfigNsxt = NewConfigNsxtValueNull()

		state.Config, err = convert.MapToDynamic(ctx, data.GetConfig())
		if err != nil {
			return state, fmt.Errorf("failed to convert generic config: %w", err)
		}
	}

	// Tenants
	tenants, d := convert.ToSetType(
		ctx,
		data.Tenants,
		func(
			in sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner,
		) TenantsValue {
			return TenantsValue{
				Id:    types.Int64Value(*in.Id),
				Name:  types.StringValue(*in.Name),
				state: attr.ValueStateKnown,
			}
		},
	)
	if d.HasError() {
		return state, fmt.Errorf("failed to convert tenants: %s", d.Errors())
	}

	state.Tenants = tenants

	// Resource permissions
	resourcePermission, ok := data.GetResourcePermissionOk()
	if ok && resourcePermission != nil {
		perms, err := convertResourcePermissions(ctx, resourcePermission)
		if err != nil {
			return state, fmt.Errorf("failed to convert resource permissions: %w", err)
		}

		// The API does not return per-group (site) assignments in the GET response,
		// so preserve groups from plan/state. After import they will be null.
		if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() &&
			!plan.Permissions.Groups.IsNull() && !plan.Permissions.Groups.IsUnknown() {
			perms.Groups = plan.Permissions.Groups
		}

		state.Permissions = perms
	} else {
		state.Permissions = NewPermissionsValueNull()
	}

	return state, nil
}

func convertResourcePermissions(
	ctx context.Context,
	resourcePermission *sdk.GetLoadBalancer200ResponseLoadBalancerResourcePermission,
) (PermissionsValue, error) {
	var groupIDValues []attr.Value

	groups, ok := resourcePermission.GetSitesOk()
	if ok {
		for _, group := range groups {
			if group.Id != nil {
				groupIDValues = append(
					groupIDValues, types.Int64Value(*group.Id),
				)
			}
		}
	}

	var groupIDsList attr.Value
	if len(groupIDValues) > 0 {
		l, d := types.ListValue(types.Int64Type, groupIDValues)
		if d.HasError() {
			return PermissionsValue{}, fmt.Errorf("failed to build groups list: %s", d.Errors())
		}

		groupIDsList = l
	} else {
		groupIDsList = types.ListNull(types.Int64Type)
	}

	perms, d := NewPermissionsValue(
		PermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"all": types.BoolValue(
				resourcePermission.All != nil &&
					*resourcePermission.All,
			),
			"groups": groupIDsList,
		},
	)
	if d.HasError() {
		return PermissionsValue{}, fmt.Errorf("failed to build permissions value: %s", d.Errors())
	}

	return perms, nil
}

func parseHAProxyConfig(
	ctx context.Context,
	configMap map[string]interface{},
) (ConfigHaproxyValue, error) {
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
	if d.HasError() {
		return ConfigHaproxyValue{}, fmt.Errorf("failed to build config_haproxy value: %s", d.Errors())
	}

	return cfg, nil
}

func parseNsxtConfig(
	ctx context.Context,
	configMap map[string]interface{},
) (ConfigNsxtValue, error) {
	var adminState bool
	var logLevel string
	var size string
	var tier1 string

	if adminStateRaw, ok := configMap["adminState"]; ok {
		if v, ok := adminStateRaw.(bool); ok {
			adminState = v
		}
	}

	if logLevelRaw, ok := configMap["loglevel"]; ok {
		if v, ok := logLevelRaw.(string); ok {
			logLevel = v
		}
	}

	if sizeRaw, ok := configMap["size"]; ok {
		if v, ok := sizeRaw.(string); ok {
			size = v
		}
	}

	if tier1Raw, ok := configMap["tier1"]; ok {
		if v, ok := tier1Raw.(string); ok {
			tier1 = v
		}
	}

	cfg, d := NewConfigNsxtValue(
		ConfigNsxtValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"admin_state":   types.BoolValue(adminState),
			"log_level":     types.StringValue(logLevel),
			"size":          types.StringValue(size),
			"tier1_gateway": types.StringValue(tier1),
		},
	)
	if d.HasError() {
		return ConfigNsxtValue{}, fmt.Errorf("failed to build config_nsxt value: %s", d.Errors())
	}

	return cfg, nil
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

	state, err := getLoadBalancerAsState(ctx, id, client, plan)
	if err != nil {
		resp.Diagnostics.AddError("read load balancer resource", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
