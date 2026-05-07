// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// nsxTier0Config builds a typed NSX Tier-0 gateway config using the SDK struct.
func nsxTier0Config(cfg ConfigNsxGatewayTier0Value) sdk.CreateNetworkRouterRequestNetworkRouterConfig {
	t0 := sdk.NewCreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOfWithDefaults()

	if !cfg.HaMode.IsNull() && !cfg.HaMode.IsUnknown() {
		t0.SetHaMode(cfg.HaMode.ValueString())
	}
	if !cfg.EdgeCluster.IsNull() && !cfg.EdgeCluster.IsUnknown() {
		t0.SetEdgeCluster(cfg.EdgeCluster.ValueString())
	}
	if !cfg.FailOver.IsNull() && !cfg.FailOver.IsUnknown() {
		t0.SetFailOver(cfg.FailOver.ValueString())
	}
	if !cfg.IpManagementType.IsNull() && !cfg.IpManagementType.IsUnknown() {
		t0.SetIpManagementType(cfg.IpManagementType.ValueString())
	}
	if !cfg.IpServerId.IsNull() && !cfg.IpServerId.IsUnknown() {
		t0.SetIpServerId(cfg.IpServerId.ValueString())
	}
	if !cfg.LocalAsNum.IsNull() && !cfg.LocalAsNum.IsUnknown() {
		t0.SetLOCAL_AS_NUM(cfg.LocalAsNum.ValueString())
	}
	if !cfg.RestartMode.IsNull() && !cfg.RestartMode.IsUnknown() {
		t0.SetRESTART_MODE(cfg.RestartMode.ValueString())
	}
	if !cfg.RestartTime.IsNull() && !cfg.RestartTime.IsUnknown() {
		t0.SetRESTART_TIME(cfg.RestartTime.ValueInt64())
	}
	if !cfg.StaleRouteTime.IsNull() && !cfg.StaleRouteTime.IsUnknown() {
		t0.SetSTALE_ROUTE_TIME(cfg.StaleRouteTime.ValueInt64())
	}
	if !cfg.Ecmp.IsNull() && !cfg.Ecmp.IsUnknown() {
		t0.SetECMP(convert.BoolToStringOnOff(cfg.Ecmp.ValueBool()).ValueString())
	}
	if !cfg.MultipathRelax.IsNull() && !cfg.MultipathRelax.IsUnknown() {
		t0.SetMULTIPATH_RELAX(convert.BoolToStringOnOff(cfg.MultipathRelax.ValueBool()).ValueString())
	}
	if !cfg.Tier0Static.IsNull() && !cfg.Tier0Static.IsUnknown() {
		t0.SetTIER0STATIC(convert.BoolToStringOnOff(cfg.Tier0Static.ValueBool()).ValueString())
	}
	if !cfg.Tier0Nat.IsNull() && !cfg.Tier0Nat.IsUnknown() {
		t0.SetTIER0NAT(convert.BoolToStringOnOff(cfg.Tier0Nat.ValueBool()).ValueString())
	}
	if !cfg.Tier0IpsecLocalIp.IsNull() && !cfg.Tier0IpsecLocalIp.IsUnknown() {
		t0.SetTIER0IPSECLOCALIP(convert.BoolToStringOnOff(cfg.Tier0IpsecLocalIp.ValueBool()).ValueString())
	}
	if !cfg.Tier0DnsForwarderIp.IsNull() && !cfg.Tier0DnsForwarderIp.IsUnknown() {
		t0.SetTIER0DNSFORWARDERIP(convert.BoolToStringOnOff(cfg.Tier0DnsForwarderIp.ValueBool()).ValueString())
	}
	if !cfg.Tier0ServiceInterface.IsNull() && !cfg.Tier0ServiceInterface.IsUnknown() {
		t0.SetTIER0SERVICEINTERFACE(convert.BoolToStringOnOff(cfg.Tier0ServiceInterface.ValueBool()).ValueString())
	}
	if !cfg.Tier0ExternalInterface.IsNull() && !cfg.Tier0ExternalInterface.IsUnknown() {
		t0.SetTIER0EXTERNALINTERFACE(convert.BoolToStringOnOff(cfg.Tier0ExternalInterface.ValueBool()).ValueString())
	}
	if !cfg.Tier0LoopbackInterface.IsNull() && !cfg.Tier0LoopbackInterface.IsUnknown() {
		t0.SetTIER0LOOPBACKINTERFACE(convert.BoolToStringOnOff(cfg.Tier0LoopbackInterface.ValueBool()).ValueString())
	}
	if !cfg.Tier0Segment.IsNull() && !cfg.Tier0Segment.IsUnknown() {
		t0.SetTIER0SEGMENT(convert.BoolToStringOnOff(cfg.Tier0Segment.ValueBool()).ValueString())
	}
	if !cfg.Tier1DnsForwarderIp.IsNull() && !cfg.Tier1DnsForwarderIp.IsUnknown() {
		t0.SetTIER1DNSFORWARDERIP(convert.BoolToStringOnOff(cfg.Tier1DnsForwarderIp.ValueBool()).ValueString())
	}
	if !cfg.Tier1Static.IsNull() && !cfg.Tier1Static.IsUnknown() {
		t0.SetTIER1STATIC(convert.BoolToStringOnOff(cfg.Tier1Static.ValueBool()).ValueString())
	}
	if !cfg.Tier1LbVip.IsNull() && !cfg.Tier1LbVip.IsUnknown() {
		t0.SetTIER1LBVIP(convert.BoolToStringOnOff(cfg.Tier1LbVip.ValueBool()).ValueString())
	}
	if !cfg.Tier1Nat.IsNull() && !cfg.Tier1Nat.IsUnknown() {
		t0.SetTIER1NAT(convert.BoolToStringOnOff(cfg.Tier1Nat.ValueBool()).ValueString())
	}
	if !cfg.Tier1LbSnat.IsNull() && !cfg.Tier1LbSnat.IsUnknown() {
		t0.SetTIER1LBSNAT(convert.BoolToStringOnOff(cfg.Tier1LbSnat.ValueBool()).ValueString())
	}
	if !cfg.Tier1IpsecLocalEndpoint.IsNull() && !cfg.Tier1IpsecLocalEndpoint.IsUnknown() {
		t0.SetTIER1IPSECLOCALENDPOINT(convert.BoolToStringOnOff(cfg.Tier1IpsecLocalEndpoint.ValueBool()).ValueString())
	}
	if !cfg.Tier1ServiceInterface.IsNull() && !cfg.Tier1ServiceInterface.IsUnknown() {
		t0.SetTIER1SERVICEINTERFACE(convert.BoolToStringOnOff(cfg.Tier1ServiceInterface.ValueBool()).ValueString())
	}
	if !cfg.Tier1Segment.IsNull() && !cfg.Tier1Segment.IsUnknown() {
		t0.SetTIER1SEGMENT(convert.BoolToStringOnOff(cfg.Tier1Segment.ValueBool()).ValueString())
	}

	if !cfg.InterSrIbgp.IsNull() && !cfg.InterSrIbgp.IsUnknown() {
		t0.SetINTER_SR_IBGP(convert.BoolToStringOnOff(cfg.Tier1Segment.ValueBool()).ValueString())
	}

	anyOf := sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOfAsCreateNetworkRouterRequestNetworkRouterConfigAnyOf( //nolint:lll
		t0,
	)

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{
		CreateNetworkRouterRequestNetworkRouterConfigAnyOf: &anyOf,
	}
}

// nsxTier1Config builds a typed NSX Tier-1 gateway config using the SDK struct.
func nsxTier1Config(cfg ConfigNsxGatewayTier1Value) sdk.CreateNetworkRouterRequestNetworkRouterConfig {
	t1 := sdk.NewCreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf1WithDefaults()

	if !cfg.Tier0Gateway.IsNull() && !cfg.Tier0Gateway.IsUnknown() {
		t1.SetTier0Gateway(cfg.Tier0Gateway.ValueString())
	}
	if !cfg.EdgeCluster.IsNull() && !cfg.EdgeCluster.IsUnknown() {
		t1.SetEdgeCluster(cfg.EdgeCluster.ValueString())
	}
	if !cfg.FailOver.IsNull() && !cfg.FailOver.IsUnknown() {
		t1.SetFailOver(cfg.FailOver.ValueString())
	}
	if !cfg.IpManagementType.IsNull() && !cfg.IpManagementType.IsUnknown() {
		t1.SetIpManagementType(cfg.IpManagementType.ValueString())
	}
	if !cfg.IpServerId.IsNull() && !cfg.IpServerId.IsUnknown() {
		t1.SetIpServerId(cfg.IpServerId.ValueString())
	}
	if !cfg.Tier1Connected.IsNull() && !cfg.Tier1Connected.IsUnknown() {
		t1.SetTIER1CONNECTED(convert.BoolToStringOnOff(cfg.Tier1Connected.ValueBool()).ValueString())
	}
	if !cfg.Tier1Nat.IsNull() && !cfg.Tier1Nat.IsUnknown() {
		t1.SetTIER1NAT(convert.BoolToStringOnOff(cfg.Tier1Nat.ValueBool()).ValueString())
	}
	if !cfg.Tier1StaticRoutes.IsNull() && !cfg.Tier1StaticRoutes.IsUnknown() {
		t1.SetTIER1STATICROUTES(convert.BoolToStringOnOff(cfg.Tier1StaticRoutes.ValueBool()).ValueString())
	}
	if !cfg.Tier1LbVip.IsNull() && !cfg.Tier1LbVip.IsUnknown() {
		t1.SetTIER1LBVIP(convert.BoolToStringOnOff(cfg.Tier1LbVip.ValueBool()).ValueString())
	}
	if !cfg.Tier1LbSnat.IsNull() && !cfg.Tier1LbSnat.IsUnknown() {
		t1.SetTIER1LBSNAT(convert.BoolToStringOnOff(cfg.Tier1LbSnat.ValueBool()).ValueString())
	}
	if !cfg.Tier1DnsForwarderIp.IsNull() && !cfg.Tier1DnsForwarderIp.IsUnknown() {
		t1.SetTIER1DNSFORWARDERIP(convert.BoolToStringOnOff(cfg.Tier1DnsForwarderIp.ValueBool()).ValueString())
	}
	if !cfg.Tier1IpsecLocalEndpoint.IsNull() && !cfg.Tier1IpsecLocalEndpoint.IsUnknown() {
		t1.SetTIER1IPSECLOCALENDPOINT(convert.BoolToStringOnOff(cfg.Tier1IpsecLocalEndpoint.ValueBool()).ValueString())
	}

	anyOf := sdk.CreateNetworkRouterRequestNetworkRouterConfigAnyOfOneOf1AsCreateNetworkRouterRequestNetworkRouterConfigAnyOf( //nolint:lll
		t1,
	)

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{
		CreateNetworkRouterRequestNetworkRouterConfigAnyOf: &anyOf,
	}
}
