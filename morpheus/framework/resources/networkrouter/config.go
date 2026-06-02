// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// nsxTier0Config builds a typed NSX Tier-0 gateway config using the SDK struct.
func nsxTier0Config(cfg ConfigNsxtGatewayTier0Value) sdk.CreateNetworkRouterRequestNetworkRouterConfig {
	t0 := &sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf{}

	if !cfg.HaMode.IsNull() && !cfg.HaMode.IsUnknown() {
		haMode := cfg.HaMode.ValueString()
		t0.HaMode = &haMode
	}
	if !cfg.EdgeCluster.IsNull() && !cfg.EdgeCluster.IsUnknown() {
		edgeCluster := cfg.EdgeCluster.ValueString()
		t0.EdgeCluster = &edgeCluster
	}
	if !cfg.FailOver.IsNull() && !cfg.FailOver.IsUnknown() {
		failOver := cfg.FailOver.ValueString()
		t0.FailOver = &failOver
	}
	if !cfg.IpManagementType.IsNull() && !cfg.IpManagementType.IsUnknown() {
		ipManagementType := cfg.IpManagementType.ValueString()
		t0.IpManagementType = &ipManagementType
	}
	if !cfg.IpServerId.IsNull() && !cfg.IpServerId.IsUnknown() {
		ipServerId := cfg.IpServerId.ValueString()
		t0.IpServerId = &ipServerId
	}
	if !cfg.LocalAsNum.IsNull() && !cfg.LocalAsNum.IsUnknown() {
		localAsNum := cfg.LocalAsNum.ValueString()
		t0.LOCAL_AS_NUM = &localAsNum
	}
	if !cfg.RestartMode.IsNull() && !cfg.RestartMode.IsUnknown() {
		restartMode := cfg.RestartMode.ValueString()
		t0.RESTART_MODE = &restartMode
	}
	if !cfg.RestartTime.IsNull() && !cfg.RestartTime.IsUnknown() {
		restartTime := cfg.RestartTime.ValueInt64()
		t0.RESTART_TIME = &restartTime
	}
	if !cfg.StaleRouteTime.IsNull() && !cfg.StaleRouteTime.IsUnknown() {
		staleRouteTime := cfg.StaleRouteTime.ValueInt64()
		t0.STALE_ROUTE_TIME = &staleRouteTime
	}
	if !cfg.Ecmp.IsNull() && !cfg.Ecmp.IsUnknown() {
		ecmp := convert.BoolToStringOnOff(cfg.Ecmp.ValueBool()).ValueString()
		t0.ECMP = &ecmp
	}
	if !cfg.MultipathRelax.IsNull() && !cfg.MultipathRelax.IsUnknown() {
		multipathRelax := convert.BoolToStringOnOff(cfg.MultipathRelax.ValueBool()).ValueString()
		t0.MULTIPATH_RELAX = &multipathRelax
	}
	if !cfg.Tier0Static.IsNull() && !cfg.Tier0Static.IsUnknown() {
		tier0Static := convert.BoolToStringOnOff(cfg.Tier0Static.ValueBool()).ValueString()
		t0.TIER0STATIC = &tier0Static
	}
	if !cfg.Tier0Nat.IsNull() && !cfg.Tier0Nat.IsUnknown() {
		tier0Nat := convert.BoolToStringOnOff(cfg.Tier0Nat.ValueBool()).ValueString()
		t0.TIER0NAT = &tier0Nat
	}
	if !cfg.Tier0IpsecLocalIp.IsNull() && !cfg.Tier0IpsecLocalIp.IsUnknown() {
		tier0IpsecLocalIp := convert.BoolToStringOnOff(cfg.Tier0IpsecLocalIp.ValueBool()).ValueString()
		t0.TIER0IPSECLOCALIP = &tier0IpsecLocalIp
	}
	if !cfg.Tier0DnsForwarderIp.IsNull() && !cfg.Tier0DnsForwarderIp.IsUnknown() {
		tier0DnsForwarderIp := convert.BoolToStringOnOff(cfg.Tier0DnsForwarderIp.ValueBool()).ValueString()
		t0.TIER0DNSFORWARDERIP = &tier0DnsForwarderIp
	}
	if !cfg.Tier0ServiceInterface.IsNull() && !cfg.Tier0ServiceInterface.IsUnknown() {
		tier0ServiceInterface := convert.BoolToStringOnOff(cfg.Tier0ServiceInterface.ValueBool()).ValueString()
		t0.TIER0SERVICEINTERFACE = &tier0ServiceInterface
	}
	if !cfg.Tier0ExternalInterface.IsNull() && !cfg.Tier0ExternalInterface.IsUnknown() {
		tier0ExternalInterface := convert.BoolToStringOnOff(cfg.Tier0ExternalInterface.ValueBool()).ValueString()
		t0.TIER0EXTERNALINTERFACE = &tier0ExternalInterface
	}
	if !cfg.Tier0LoopbackInterface.IsNull() && !cfg.Tier0LoopbackInterface.IsUnknown() {
		tier0LoopbackInterface := convert.BoolToStringOnOff(cfg.Tier0LoopbackInterface.ValueBool()).ValueString()
		t0.TIER0LOOPBACKINTERFACE = &tier0LoopbackInterface
	}
	if !cfg.Tier0Segment.IsNull() && !cfg.Tier0Segment.IsUnknown() {
		tier0Segment := convert.BoolToStringOnOff(cfg.Tier0Segment.ValueBool()).ValueString()
		t0.TIER0SEGMENT = &tier0Segment
	}
	if !cfg.Tier1DnsForwarderIp.IsNull() && !cfg.Tier1DnsForwarderIp.IsUnknown() {
		tier1DnsForwarderIp := convert.BoolToStringOnOff(cfg.Tier1DnsForwarderIp.ValueBool()).ValueString()
		t0.TIER1DNSFORWARDERIP = &tier1DnsForwarderIp
	}
	if !cfg.Tier1Static.IsNull() && !cfg.Tier1Static.IsUnknown() {
		tier1Static := convert.BoolToStringOnOff(cfg.Tier1Static.ValueBool()).ValueString()
		t0.TIER1STATIC = &tier1Static
	}
	if !cfg.Tier1LbVip.IsNull() && !cfg.Tier1LbVip.IsUnknown() {
		tier1LbVip := convert.BoolToStringOnOff(cfg.Tier1LbVip.ValueBool()).ValueString()
		t0.TIER1LBVIP = &tier1LbVip
	}
	if !cfg.Tier1Nat.IsNull() && !cfg.Tier1Nat.IsUnknown() {
		tier1Nat := convert.BoolToStringOnOff(cfg.Tier1Nat.ValueBool()).ValueString()
		t0.TIER1NAT = &tier1Nat
	}
	if !cfg.Tier1LbSnat.IsNull() && !cfg.Tier1LbSnat.IsUnknown() {
		tier1LbSnat := convert.BoolToStringOnOff(cfg.Tier1LbSnat.ValueBool()).ValueString()
		t0.TIER1LBSNAT = &tier1LbSnat
	}
	if !cfg.Tier1IpsecLocalEndpoint.IsNull() && !cfg.Tier1IpsecLocalEndpoint.IsUnknown() {
		tier1IpsecLocalEndpoint := convert.BoolToStringOnOff(cfg.Tier1IpsecLocalEndpoint.ValueBool()).ValueString()
		t0.TIER1IPSECLOCALENDPOINT = &tier1IpsecLocalEndpoint
	}
	if !cfg.Tier1ServiceInterface.IsNull() && !cfg.Tier1ServiceInterface.IsUnknown() {
		tier1ServiceInterface := convert.BoolToStringOnOff(cfg.Tier1ServiceInterface.ValueBool()).ValueString()
		t0.TIER1SERVICEINTERFACE = &tier1ServiceInterface
	}
	if !cfg.Tier1Segment.IsNull() && !cfg.Tier1Segment.IsUnknown() {
		tier1Segment := convert.BoolToStringOnOff(cfg.Tier1Segment.ValueBool()).ValueString()
		t0.TIER1SEGMENT = &tier1Segment
	}

	if !cfg.InterSrIbgp.IsNull() && !cfg.InterSrIbgp.IsUnknown() {
		interSrIbgp := convert.BoolToStringOnOff(cfg.InterSrIbgp.ValueBool()).ValueString()
		t0.INTER_SR_IBGP = &interSrIbgp
	}

	anyOf := sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf(*t0)

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{
		CreateNetworkRouterRequestNetworkRouterConfigAnyOf: &sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOf{
			CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf: &anyOf,
		},
	}
}

// nsxTier1Config builds a typed NSX Tier-1 gateway config using the SDK struct.
func nsxTier1Config(cfg ConfigNsxtGatewayTier1Value) sdk.CreateNetworkRouterRequestNetworkRouterConfig {
	t1 := &sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf1{}

	if !cfg.Tier0Gateway.IsNull() && !cfg.Tier0Gateway.IsUnknown() {
		tier0Gateway := cfg.Tier0Gateway.ValueString()
		t1.Tier0Gateway = &tier0Gateway
	}
	if !cfg.EdgeCluster.IsNull() && !cfg.EdgeCluster.IsUnknown() {
		edgeCluster := cfg.EdgeCluster.ValueString()
		t1.EdgeCluster = &edgeCluster
	}
	if !cfg.FailOver.IsNull() && !cfg.FailOver.IsUnknown() {
		failOver := cfg.FailOver.ValueString()
		t1.FailOver = &failOver
	}
	if !cfg.IpManagementType.IsNull() && !cfg.IpManagementType.IsUnknown() {
		ipManagementType := cfg.IpManagementType.ValueString()
		t1.IpManagementType = &ipManagementType
	}
	if !cfg.IpServerId.IsNull() && !cfg.IpServerId.IsUnknown() {
		ipServerId := cfg.IpServerId.ValueString()
		t1.IpServerId = &ipServerId
	}
	if !cfg.Tier1Connected.IsNull() && !cfg.Tier1Connected.IsUnknown() {
		tier1Connected := convert.BoolToStringOnOff(cfg.Tier1Connected.ValueBool()).ValueString()
		t1.TIER1CONNECTED = &tier1Connected
	}
	if !cfg.Tier1Nat.IsNull() && !cfg.Tier1Nat.IsUnknown() {
		tier1Nat := convert.BoolToStringOnOff(cfg.Tier1Nat.ValueBool()).ValueString()
		t1.TIER1NAT = &tier1Nat
	}
	if !cfg.Tier1StaticRoutes.IsNull() && !cfg.Tier1StaticRoutes.IsUnknown() {
		tier1StaticRoutes := convert.BoolToStringOnOff(cfg.Tier1StaticRoutes.ValueBool()).ValueString()
		t1.TIER1STATICROUTES = &tier1StaticRoutes
	}
	if !cfg.Tier1LbVip.IsNull() && !cfg.Tier1LbVip.IsUnknown() {
		tier1LbVip := convert.BoolToStringOnOff(cfg.Tier1LbVip.ValueBool()).ValueString()
		t1.TIER1LBVIP = &tier1LbVip
	}
	if !cfg.Tier1LbSnat.IsNull() && !cfg.Tier1LbSnat.IsUnknown() {
		tier1LbSnat := convert.BoolToStringOnOff(cfg.Tier1LbSnat.ValueBool()).ValueString()
		t1.TIER1LBSNAT = &tier1LbSnat
	}
	if !cfg.Tier1DnsForwarderIp.IsNull() && !cfg.Tier1DnsForwarderIp.IsUnknown() {
		tier1DnsForwarderIp := convert.BoolToStringOnOff(cfg.Tier1DnsForwarderIp.ValueBool()).ValueString()
		t1.TIER1DNSFORWARDERIP = &tier1DnsForwarderIp
	}
	if !cfg.Tier1IpsecLocalEndpoint.IsNull() && !cfg.Tier1IpsecLocalEndpoint.IsUnknown() {
		tier1IpsecLocalEndpoint := convert.BoolToStringOnOff(cfg.Tier1IpsecLocalEndpoint.ValueBool()).ValueString()
		t1.TIER1IPSECLOCALENDPOINT = &tier1IpsecLocalEndpoint
	}

	anyOf := sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf1(*t1)

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{
		CreateNetworkRouterRequestNetworkRouterConfigAnyOf: &sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOf{
			CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf1: &anyOf,
		},
	}
}
