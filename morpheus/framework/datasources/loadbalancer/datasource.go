// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read load balancer data source"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerDataSourceSchema(ctx)
}

func getLoadBalancer(
	ctx context.Context,
	config LoadBalancerModel,
	client *sdk.APIClient,
) (*LoadBalancerModel, error) {
	if !config.Id.IsNull() {
		return getLoadBalancerByID(ctx, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getLoadBalancerByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getLoadBalancerByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*LoadBalancerModel, error) {
	resp, hresp, err := client.LoadBalancersAPI.GetLoadBalancer(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load balancer %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	lb, ok := resp.GetLoadBalancerOk()
	if !ok {
		return nil, fmt.Errorf("load balancer %d is nil", id)
	}

	state := &LoadBalancerModel{}
	populateLoadBalancerState(ctx, state, lb)

	return state, nil
}

func getLoadBalancerByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*LoadBalancerModel, error) {
	lbs, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).
		Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load balancer %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matching []sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner
	for _, lb := range lbs.GetLoadBalancers() {
		if lbName, ok := lb.GetNameOk(); ok && *lbName == name {
			matching = append(matching, lb)
		}
	}

	if len(matching) == 0 {
		return nil, fmt.Errorf("load balancer %s not found", name)
	}

	if len(matching) > 1 {
		var ids []string
		for _, lb := range matching {
			if id, ok := lb.GetIdOk(); ok {
				ids = append(ids, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple load balancers found with name %s. IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(ids, ", "),
		)
	}

	id, ok := matching[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("load balancer %s has missing ID", name)
	}

	return getLoadBalancerByID(ctx, *id, client)
}

//nolint:funlen,cyclop // mapping all fields requires length
func populateLoadBalancerState(
	ctx context.Context,
	state *LoadBalancerModel,
	lb *sdk.GetLoadBalancer200ResponseLoadBalancer,
) {
	state.Id = convert.Int64ToType(lb.Id)
	state.Name = convert.StrToType(lb.Name)
	state.Description = convert.StrToType(lb.Description)
	state.Visibility = convert.StrToType(lb.Visibility)
	state.Host = convert.StrToType(lb.Host)
	state.Ip = convert.StrToType(lb.Ip)
	state.Port = convert.Int64ToType(lb.Port)
	state.TenantId = convert.Int64ToType(lb.AccountId)
	state.Uuid = convert.StrToType(lb.Uuid)
	state.DateCreated = timeToType(lb.DateCreated)
	state.LastUpdated = timeToType(lb.LastUpdated)

	// Nullable typed fields
	state.Username = convert.StrToType(lb.Username.Get())
	state.InternalIp = convert.StrToType(lb.InternalIp.Get())
	state.ExternalIp = convert.StrToType(lb.ExternalIp.Get())
	state.ApiPort = convert.StrToType(lb.ApiPort.Get())
	state.AdminPort = convert.StrToType(lb.AdminPort.Get())
	state.SslEnabled = convert.BoolToType(lb.SslEnabled.Get())
	state.SslCert = convert.StrToType(lb.SslCert.Get())

	state.Enabled = convert.BoolToType(lb.Enabled)
	state.AllowVipEntry = convert.BoolToType(lb.AllowVipEntry)
	state.Password = convert.StrToType(lb.Password.Get())
	state.PasswordHash = convert.StrToType(lb.PasswordHash.Get())
	state.ExternalId = convert.StrToType(lb.ExternalId.Get())
	state.VirtualServiceName = convert.StrToType(lb.VirtualServiceName.Get())
	state.PoolName = convert.StrToType(lb.PoolName.Get())
	state.ServerName = convert.StrToType(lb.ServerName.Get())

	// Nested objects
	state.Cloud = mapCloud(lb.Cloud)
	state.Type = mapType(lb.Type)
	state.Owner = mapOwner(lb.Owner)
	state.Credential = mapCredential(ctx, lb.Credential)
	state.Permissions = mapPermissions(ctx, lb.ResourcePermission)
	state.Tenants = mapTenants(ctx, lb.Tenants)

	state.Config = mapEmptyConfig(lb.Config)
	state.InstancePrice = mapEmptyInstancePrice(lb.InstancePrice)
	state.VipPools = mapEmptyVipPools(ctx, lb.VipPools)
}

func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

func mapCloud(c *sdk.GetLoadBalancer200ResponseLoadBalancerCloud) CloudValue {
	if c == nil {
		return NewCloudValueNull()
	}

	return CloudValue{
		Id:    convert.Int64ToType(c.Id),
		Name:  convert.StrToType(c.Name),
		state: attr.ValueStateKnown,
	}
}

func mapType(t *sdk.GetLoadBalancer200ResponseLoadBalancerType) TypeValue {
	if t == nil {
		return NewTypeValueNull()
	}

	return TypeValue{
		Id:    convert.Int64ToType(t.Id),
		Code:  convert.StrToType(t.Code),
		Name:  convert.StrToType(t.Name),
		state: attr.ValueStateKnown,
	}
}

func mapOwner(o *sdk.GetLoadBalancer200ResponseLoadBalancerOwner) OwnerValue {
	if o == nil {
		return NewOwnerValueNull()
	}

	return OwnerValue{
		Id:    convert.Int64ToType(o.Id),
		Name:  convert.StrToType(o.Name),
		state: attr.ValueStateKnown,
	}
}

func mapCredential(
	_ context.Context,
	c *sdk.GetLoadBalancer200ResponseLoadBalancerCredential,
) CredentialValue {
	if c == nil {
		return NewCredentialValueNull()
	}

	return CredentialValue{
		Id:             convert.Int64ToType(c.Id.Get()),
		Name:           convert.StrToType(c.Name.Get()),
		CredentialType: convert.StrToType(c.Type),
		Types:          convert.StrSliceToSet(c.Types),
		state:          attr.ValueStateKnown,
	}
}

func mapPermissions(
	ctx context.Context,
	rp *sdk.GetLoadBalancer200ResponseLoadBalancerResourcePermission,
) PermissionsValue {
	if rp == nil {
		return NewPermissionsValueNull()
	}

	groupsSet := mapPermissionGroups(ctx, rp.Sites)
	plansSet := mapPermissionPlans(ctx, rp.Plans)

	var tenantId types.Int64
	if rp.Account != nil && rp.Account.Id != nil {
		tenantId = types.Int64Value(*rp.Account.Id)
	} else {
		tenantId = types.Int64Null()
	}

	return PermissionsValue{
		All:           convert.BoolToType(rp.All),
		AllPlans:      convert.BoolToType(rp.AllPlans),
		CanManage:     convert.BoolToType(rp.CanManage),
		DefaultStore:  convert.BoolToType(rp.DefaultStore),
		DefaultTarget: convert.BoolToType(rp.DefaultTarget),
		Groups:        groupsSet,
		Plans:         plansSet,
		TenantId:      tenantId,
		state:         attr.ValueStateKnown,
	}
}

func mapPermissionGroups(
	ctx context.Context,
	sites []sdk.GetLoadBalancer200ResponseLoadBalancerResourcePermissionSitesInner,
) types.Set {
	elemType := GroupsType{
		ObjectType: types.ObjectType{
			AttrTypes: GroupsValue{}.AttributeTypes(ctx),
		},
	}

	if len(sites) == 0 {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	vals := make([]attr.Value, 0, len(sites))
	for _, s := range sites {
		vals = append(vals, GroupsValue{
			Default: convert.BoolToType(s.Default),
			Id:      convert.Int64ToType(s.Id),
			Name:    convert.StrToType(s.Name),
			state:   attr.ValueStateKnown,
		})
	}

	set, diags := types.SetValueFrom(ctx, elemType, vals)
	if diags.HasError() {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	return set
}

func mapPermissionPlans(
	ctx context.Context,
	plans []sdk.GetLoadBalancer200ResponseLoadBalancerResourcePermissionPlansInner,
) types.Set {
	elemType := PlansType{
		ObjectType: types.ObjectType{
			AttrTypes: PlansValue{}.AttributeTypes(ctx),
		},
	}

	if len(plans) == 0 {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	vals := make([]attr.Value, 0, len(plans))
	for _, p := range plans {
		vals = append(vals, PlansValue{
			Default: convert.BoolToType(p.Default),
			Id:      convert.Int64ToType(p.Id),
			Name:    convert.StrToType(p.Name),
			state:   attr.ValueStateKnown,
		})
	}

	set, diags := types.SetValueFrom(ctx, elemType, vals)
	if diags.HasError() {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	return set
}

func mapTenants(
	ctx context.Context,
	tenants []sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner,
) types.Set {
	elemType := TenantsType{
		ObjectType: types.ObjectType{
			AttrTypes: TenantsValue{}.AttributeTypes(ctx),
		},
	}

	if len(tenants) == 0 {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	vals := make([]attr.Value, 0, len(tenants))
	for _, t := range tenants {
		vals = append(vals, TenantsValue{
			Id:    convert.Int64ToType(t.Id),
			Name:  convert.StrToType(t.Name),
			state: attr.ValueStateKnown,
		})
	}

	set, diags := types.SetValueFrom(ctx, elemType, vals)
	if diags.HasError() {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	return set
}

func mapEmptyConfig(config map[string]interface{}) ConfigValue {
	if config == nil {
		return NewConfigValueNull()
	}

	return ConfigValue{state: attr.ValueStateKnown}
}

func mapEmptyInstancePrice(instancePrice map[string]interface{}) InstancePriceValue {
	if instancePrice == nil {
		return NewInstancePriceValueNull()
	}

	return InstancePriceValue{state: attr.ValueStateKnown}
}

func mapEmptyVipPools(ctx context.Context, vipPools []map[string]interface{}) types.Set {
	elemType := VipPoolsType{
		ObjectType: types.ObjectType{
			AttrTypes: VipPoolsValue{}.AttributeTypes(ctx),
		},
	}

	if vipPools == nil {
		return types.SetNull(elemType)
	}

	vals := make([]attr.Value, 0, len(vipPools))
	for range vipPools {
		vals = append(vals, VipPoolsValue{state: attr.ValueStateKnown})
	}

	return types.SetValueMust(elemType, vals)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getLoadBalancer(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
