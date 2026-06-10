// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read load balancer virtual server data source"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_virtual_server"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerVirtualServerDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerVirtualServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("failed to create client: %s", err.Error()))

		return
	}

	state, err := getLoadBalancerVirtualServer(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func getLoadBalancerVirtualServer(
	ctx context.Context,
	config LoadBalancerVirtualServerModel,
	client *sdk.APIClient,
) (*LoadBalancerVirtualServerModel, error) {
	loadBalancerID := config.LoadBalancerId.ValueInt64()

	if !config.Id.IsNull() {
		return getVirtualServerByID(ctx, loadBalancerID, config.Id.ValueInt64(), client)
	}

	if !config.VipName.IsNull() {
		return getVirtualServerByName(ctx, loadBalancerID, config.VipName.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or vip_name must be specified")
}

func getVirtualServerByID(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (*LoadBalancerVirtualServerModel, error) {
	resp, hresp, err := client.LoadBalancersAPI.
		GetLoadBalancerVirtualServer(ctx, loadBalancerID, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d virtual server %d GET failed: %s",
			loadBalancerID, id, errfmt.ErrMsg(err, hresp),
		)
	}

	if resp.LoadBalancerInstance == nil {
		return nil, fmt.Errorf(
			"load balancer %d virtual server %d GET returned no loadBalancerInstance",
			loadBalancerID, id,
		)
	}

	state := populateVirtualServerState(ctx, loadBalancerID, resp.LoadBalancerInstance)

	return state, nil
}

func getVirtualServerByName(
	ctx context.Context,
	loadBalancerID int64,
	vipName string,
	client *sdk.APIClient,
) (*LoadBalancerVirtualServerModel, error) {
	list, hresp, err := client.LoadBalancersAPI.
		ListLoadBalancerVirtualServers(ctx, loadBalancerID).
		VipName(vipName).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d virtual server list failed: %s",
			loadBalancerID, errfmt.ErrMsg(err, hresp),
		)
	}

	var matching []sdk.ListLoadBalancerVirtualServers200ResponseAllOfLoadBalancerInstancesInner
	for _, vs := range list.LoadBalancerInstances {
		if vs.VipName != nil && *vs.VipName == vipName {
			matching = append(matching, vs)
		}
	}

	if len(matching) == 0 {
		return nil, fmt.Errorf(
			"load balancer %d virtual server with vip_name %q not found",
			loadBalancerID, vipName,
		)
	}

	if len(matching) > 1 {
		var ids []string
		for _, vs := range matching {
			if vs.Id != nil {
				ids = append(ids, fmt.Sprintf("%d", *vs.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple virtual servers found with vip_name %q. IDs: %s. "+
				"Please specify an ID instead",
			vipName,
			strings.Join(ids, ", "),
		)
	}

	id := matching[0].Id
	if id == nil {
		return nil, fmt.Errorf(
			"load balancer %d virtual server with vip_name %q has missing ID",
			loadBalancerID, vipName,
		)
	}

	return getVirtualServerByID(ctx, loadBalancerID, *id, client)
}

//nolint:funlen,cyclop // mapping all fields requires length
func populateVirtualServerState(
	ctx context.Context,
	loadBalancerID int64,
	vs *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstance,
) *LoadBalancerVirtualServerModel {
	state := &LoadBalancerVirtualServerModel{}

	state.Id = convert.Int64ToType(vs.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Description = convert.StrToType(vs.Description.Get())
	state.VipName = convert.StrToType(vs.VipName)
	state.VipAddress = convert.StrToType(vs.VipAddress)
	state.VipPort = convert.Int64ToType(vs.VipPort)
	state.VipProtocol = convert.StrToType(vs.VipProtocol)
	state.VipHostname = convert.StrToType(vs.VipHostname.Get())

	// Computed scalar fields
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
	if vs.DateCreated != nil {
		state.DateCreated = types.StringValue(vs.DateCreated.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if vs.LastUpdated != nil {
		state.LastUpdated = types.StringValue(vs.LastUpdated.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Pool — nested {id, name} object.
	if vs.Pool != nil {
		state.Pool = PoolValue{
			Id:    convert.Int64ToType(vs.Pool.Id),
			Name:  convert.StrToType(vs.Pool.Name),
			state: attr.ValueStateKnown,
		}
	} else {
		state.Pool = NewPoolValueNull()
	}

	// SSL cert — nested {id, name} object.
	if vs.SslCert != nil {
		state.SslCert = SslCertValue{
			Id:    convert.Int64ToType(vs.SslCert.Id),
			Name:  convert.StrToType(vs.SslCert.Name),
			state: attr.ValueStateKnown,
		}
	} else {
		state.SslCert = NewSslCertValueNull()
	}

	// SSL server cert — nested {id, name} object.
	if vs.SslServerCert != nil {
		state.SslServerCert = SslServerCertValue{
			Id:    convert.Int64ToType(vs.SslServerCert.Id),
			Name:  convert.StrToType(vs.SslServerCert.Name),
			state: attr.ValueStateKnown,
		}
	} else {
		state.SslServerCert = NewSslServerCertValueNull()
	}

	// VIP pool — empty nested object (the schema has no attributes inside).
	if vs.VipPool != nil {
		state.VipPool = VipPoolValue{state: attr.ValueStateKnown}
	} else {
		state.VipPool = NewVipPoolValueNull()
	}

	// Load balancer — nested object with id, name, ip, type.
	if vs.LoadBalancer != nil {
		lb := vs.LoadBalancer

		typeVal := types.ObjectNull(TypeValue{}.AttributeTypes(ctx))

		if lb.Type != nil {
			tv, d := NewTypeValue(
				TypeValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"code": convert.StrToType(lb.Type.Code),
					"id":   convert.Int64ToType(lb.Type.Id),
					"name": convert.StrToType(lb.Type.Name),
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

	// Config — the generic config map returned as a dynamic value.
	// For NSX-T, build a typed config_nsxt and extract pool_id from config.
	lbTypeCode := ""
	if vs.LoadBalancer != nil && vs.LoadBalancer.Type != nil && vs.LoadBalancer.Type.Code != nil {
		lbTypeCode = *vs.LoadBalancer.Type.Code
	}

	configMap := vs.Config

	switch lbTypeCode {
	case "nsx-t":
		state.Config = types.DynamicNull()

		if configMap != nil {
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
			if !d.HasError() {
				state.ConfigNsxt = nsxtVal
			} else {
				state.ConfigNsxt = NewConfigNsxtValueNull()
			}

			// For NSX-T, pool_id lives inside config as a string.
			if state.PoolId.IsNull() {
				if v, ok := configMap["pool"].(string); ok && v != "" {
					if pid, err := strconv.ParseInt(v, 10, 64); err == nil {
						state.PoolId = types.Int64Value(pid)
					}
				}
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

	// Pool_id — from the top-level pool object if present.
	if vs.Pool != nil && vs.Pool.Id != nil {
		state.PoolId = types.Int64Value(*vs.Pool.Id)
	}
	// If still null, it may have been set from configMap above for NSX-T.

	return state
}
