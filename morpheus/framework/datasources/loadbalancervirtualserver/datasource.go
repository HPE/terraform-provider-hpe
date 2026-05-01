// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package loadbalancervirtualserver implements a data source for load_balancer_virtual_server
package loadbalancervirtualserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                              = "read load balancer virtual server data source"
	ErrorNoValidSearchTerms              = `no valid search terms - an id or vip_name is required`
	ErrorRunningPreApply                 = `Error running pre-apply plan: exit status 1`
	ErrorNoLoadBalancerVirtualServer     = `no load balancer virtual server found`
	ErrorMultipleLoadBalancerVirtualServ = `multiple load balancer virtual servers were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_virtual_server"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerVirtualServerDataSourceSchema(ctx)
}

func loadBalancerIDFromNumber(n types.Number) (int64, error) {
	if n.IsNull() || n.IsUnknown() {
		return 0, fmt.Errorf("load_balancer_id is required")
	}

	bigF := n.ValueBigFloat()
	i64, _ := bigF.Int64()

	return i64, nil
}

//nolint:funlen,cyclop // mapping all fields requires length
func virtualServerAsState(
	ctx context.Context,
	vs *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstance,
	loadBalancerID int64,
) (LoadBalancerVirtualServerModel, error) {
	state := LoadBalancerVirtualServerModel{
		Id:               convert.Int64ToType(vs.Id),
		LoadBalancerId:   types.NumberValue(new(big.Float).SetInt64(loadBalancerID)),
		Description:      convert.StrToType(vs.Description.Get()),
		Active:           convert.BoolToType(vs.Active),
		BackendPort:      convert.StrToType(vs.BackendPort.Get()),
		ExternalAddress:  convert.BoolToType(vs.ExternalAddress),
		ExternalId:       convert.StrToType(vs.ExternalId),
		ExternalPortId:   convert.StrToType(vs.ExternalPortId.Get()),
		ExtraConfig:      convert.StrToType(vs.ExtraConfig.Get()),
		Instance:         convert.StrToType(vs.Instance.Get()),
		InternalId:       convert.StrToType(vs.InternalId),
		NetworkId:        convert.StrToType(vs.NetworkId.Get()),
		PoolName:         convert.StrToType(vs.PoolName.Get()),
		Removing:         convert.BoolToType(vs.Removing),
		ServerName:       convert.StrToType(vs.ServerName.Get()),
		ServiceAccess:    convert.StrToType(vs.ServiceAccess.Get()),
		ServicePort:      convert.StrToType(vs.ServicePort.Get()),
		SourceAddress:    convert.StrToType(vs.SourceAddress.Get()),
		SslEnabled:       convert.StrToType(vs.SslEnabled.Get()),
		SslMode:          convert.StrToType(vs.SslMode.Get()),
		SslRedirectMode:  convert.StrToType(vs.SslRedirectMode.Get()),
		Status:           convert.StrToType(vs.Status),
		Sticky:           convert.BoolToType(vs.Sticky),
		SubnetId:         convert.StrToType(vs.SubnetId.Get()),
		VipAddress:       convert.StrToType(vs.VipAddress),
		VipBalance:       convert.StrToType(vs.VipBalance.Get()),
		VipDirectAddress: convert.StrToType(vs.VipDirectAddress.Get()),
		VipHostname:      convert.StrToType(vs.VipHostname.Get()),
		VipMode:          convert.StrToType(vs.VipMode.Get()),
		VipName:          convert.StrToType(vs.VipName),
		VipPort:          convert.Int64ToType(vs.VipPort),
		VipProtocol:      convert.StrToType(vs.VipProtocol),
		VipScheme:        convert.StrToType(vs.VipScheme.Get()),
		VipShared:        convert.BoolToType(vs.VipShared),
		VipSource:        convert.StrToType(vs.VipSource),
		VipStatus:        convert.StrToType(vs.VipStatus),
		VipSticky:        convert.StrToType(vs.VipSticky.Get()),
		VipType:          convert.StrToType(vs.VipType.Get()),
	}

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

	// Config — always populate the generic config field.
	if vs.Config != nil {
		raw, err := json.Marshal(vs.Config)
		if err != nil {
			return LoadBalancerVirtualServerModel{},
				fmt.Errorf("error marshalling config: %w", err)
		}

		state.Config = types.DynamicValue(types.StringValue(string(raw)))
	} else {
		state.Config = types.DynamicNull()
	}

	// ConfigNsxt — populate when the load balancer type code is "nsx-t".
	lbTypeCode := getLBTypeCode(vs.LoadBalancer)
	if lbTypeCode == "nsx-t" {
		state.ConfigNsxt = buildConfigNsxt(ctx, vs.Config)
	} else {
		state.ConfigNsxt = NewConfigNsxtValueNull()
	}

	// SSL cert
	state.SslCert = mapSslCert(vs.SslCert)

	// Load balancer
	lb, err := buildLoadBalancerValue(ctx, vs.LoadBalancer)
	if err != nil {
		return LoadBalancerVirtualServerModel{}, err
	}

	state.LoadBalancer = lb

	return state, nil
}

func getVirtualServer(
	ctx context.Context,
	config *LoadBalancerVirtualServerModel,
	lbID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstance, error) {
	if !config.Id.IsNull() {
		return getVirtualServerByID(ctx, lbID, config.Id.ValueInt64(), apiClient)
	} else if !config.VipName.IsNull() {
		return getVirtualServerByVipName(ctx, lbID, config.VipName.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

func getVirtualServerByID(
	ctx context.Context,
	lbID int64,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstance, error) {
	r, hresp, err := apiClient.LoadBalancersAPI.
		GetLoadBalancerVirtualServer(ctx, lbID, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for load balancer %d virtual server %d: %s",
			lbID, id, errfmt.ErrMsg(err, hresp),
		)
	}

	vs := r.GetLoadBalancerInstance()

	return &vs, nil
}

func getVirtualServerByVipName(
	ctx context.Context,
	lbID int64,
	vipName string,
	apiClient *sdk.APIClient,
) (*sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstance, error) {
	rs, hresp, err := apiClient.LoadBalancersAPI.
		ListLoadBalancerVirtualServers(ctx, lbID).
		VipName(vipName).Max(10000).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for load balancer %d virtual server %q: %s",
			lbID, vipName, errfmt.ErrMsg(err, hresp),
		)
	}

	items := rs.GetLoadBalancerInstances()

	var matchedIDs []int64

	for i := range items {
		if items[i].GetVipName() == vipName {
			matchedIDs = append(matchedIDs, items[i].GetId())
		}
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoLoadBalancerVirtualServer)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleLoadBalancerVirtualServ)
	}

	return getVirtualServerByID(ctx, lbID, matchedIDs[0], apiClient)
}

func mapSslCert(
	cert *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstanceSslCert,
) SslCertValue {
	if cert == nil {
		return NewSslCertValueNull()
	}

	return SslCertValue{
		Id:    convert.Int64ToType(cert.Id),
		Name:  convert.StrToType(cert.Name),
		state: attr.ValueStateKnown,
	}
}

func buildLoadBalancerValue(
	ctx context.Context,
	lb *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstanceLoadBalancer,
) (LoadBalancerValue, error) {
	if lb == nil {
		return NewLoadBalancerValueNull(), nil
	}

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
			return NewLoadBalancerValueNull(),
				fmt.Errorf("error building load balancer type")
		}

		objVal, tvDiags := tv.ToObjectValue(ctx)
		if tvDiags.HasError() {
			return NewLoadBalancerValueNull(),
				fmt.Errorf("error converting load balancer type to object")
		}

		typeVal = objVal
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
		return NewLoadBalancerValueNull(),
			fmt.Errorf("error building load balancer value")
	}

	return lbVal, nil
}

// getLBTypeCode extracts the load balancer type code from the nested load balancer object.
func getLBTypeCode(
	lb *sdk.GetLoadBalancerVirtualServer200ResponseLoadBalancerInstanceLoadBalancer,
) string {
	if lb == nil {
		return ""
	}

	lbType, ok := lb.GetTypeOk()
	if !ok || lbType == nil {
		return ""
	}

	code, ok := lbType.GetCodeOk()
	if !ok || code == nil {
		return ""
	}

	return *code
}

// buildConfigNsxt extracts the applicationProfile from the config map and builds a ConfigNsxtValue.
func buildConfigNsxt(ctx context.Context, configMap map[string]interface{}) ConfigNsxtValue {
	if configMap == nil {
		return NewConfigNsxtValueNull()
	}

	appProfile, _ := configMap["applicationProfile"].(string)

	nsxtVal, d := NewConfigNsxtValue(
		ConfigNsxtValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"application_profile": types.StringValue(appProfile),
		},
	)
	if d.HasError() {
		return NewConfigNsxtValueNull()
	}

	return nsxtVal
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerVirtualServerModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	lbID, err := loadBalancerIDFromNumber(config.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	vs, err := getVirtualServer(ctx, &config, lbID, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state, err := virtualServerAsState(ctx, vs, lbID)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
