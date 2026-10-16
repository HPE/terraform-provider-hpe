// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Pool server type codes with a typed request variant in the API.
const (
	typeCodeInfoblox   = "infoblox"
	typeCodeBluecat    = "bluecat"
	typeCodePhpIPAM    = "phpipam"
	typeCodeSolarWinds = "solarwindsipam"
)

// The create and update requests are oneOf unions with a typed variant per
// pool server type. The plan is first mapped onto the Infoblox variant, which
// is the superset of every variant's fields, and then converted to the variant
// for the type code so that each type's body carries exactly its own fields:
// phpIPAM's required config.appId is a typed field of its variant, and
// Infoblox-only fields are not sent to other types.
//
// Types without a variant in the API (plugin-provided types such as
// EfficientIP, NetBox and Micetro, and the built-in Morpheus pool) use the
// Infoblox superset unchanged; the appliance selects the type from the "type"
// code and accepts the common fields.

// createRequest returns the create union for the type code.
func createRequest(
	typeCode string, s *sdk.InfobloxNetworkPoolServer, appID *string,
) sdk.CreateNetworkPoolServerRequestNetworkPoolServer {
	var inv *string
	if s.Config != nil {
		inv = s.Config.InventoryExisting
	}

	switch typeCode {
	case typeCodePhpIPAM:
		v := &sdk.PhpIPAMNetworkPoolServer{
			Type: s.Type, Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl, NetworkFilter: s.NetworkFilter,
			Config: sdk.PhpIPAMNetworkPoolServerConfig{InventoryExisting: inv},
		}
		if appID != nil {
			v.Config.AppId = *appID
		}
		if s.Credential != nil {
			v.Credential = &sdk.PhpIPAMNetworkPoolServerCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.CreateNetworkPoolServerRequestNetworkPoolServer{PhpIPAMNetworkPoolServer: v}

	case typeCodeBluecat:
		v := &sdk.BluecatNetworkPoolServer{
			Type: s.Type, Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl, NetworkFilter: s.NetworkFilter,
		}
		if inv != nil {
			v.Config = &sdk.BluecatNetworkPoolServerConfig{InventoryExisting: inv}
		}
		if s.Credential != nil {
			v.Credential = &sdk.BluecatNetworkPoolServerCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.CreateNetworkPoolServerRequestNetworkPoolServer{BluecatNetworkPoolServer: v}

	case typeCodeSolarWinds:
		v := &sdk.SolarWindsNetworkPoolServer{
			Type: s.Type, Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl,
		}
		if inv != nil {
			v.Config = &sdk.SolarWindsNetworkPoolServerConfig{InventoryExisting: inv}
		}
		if s.Credential != nil {
			v.Credential = &sdk.SolarWindsNetworkPoolServerCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.CreateNetworkPoolServerRequestNetworkPoolServer{SolarWindsNetworkPoolServer: v}

	default:
		return sdk.CreateNetworkPoolServerRequestNetworkPoolServer{InfobloxNetworkPoolServer: s}
	}
}

// updateRequest returns the update union for the type code.
func updateRequest(
	typeCode string, s *sdk.InfobloxNetworkPoolServerUpdate, appID *string,
) sdk.UpdateNetworkPoolServerRequestNetworkPoolServer {
	var inv *string
	if s.Config != nil {
		inv = s.Config.InventoryExisting
	}

	switch typeCode {
	case typeCodePhpIPAM:
		v := &sdk.PhpIPAMNetworkPoolServerUpdate{
			Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl, NetworkFilter: s.NetworkFilter,
		}
		if inv != nil || appID != nil {
			v.Config = &sdk.PhpIPAMNetworkPoolServerUpdateConfig{InventoryExisting: inv, AppId: appID}
		}
		if s.Credential != nil {
			v.Credential = &sdk.PhpIPAMNetworkPoolServerUpdateCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{PhpIPAMNetworkPoolServerUpdate: v}

	case typeCodeBluecat:
		v := &sdk.BluecatNetworkPoolServerUpdate{
			Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl, NetworkFilter: s.NetworkFilter,
		}
		if inv != nil {
			v.Config = &sdk.BluecatNetworkPoolServerUpdateConfig{InventoryExisting: inv}
		}
		if s.Credential != nil {
			v.Credential = &sdk.BluecatNetworkPoolServerUpdateCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{BluecatNetworkPoolServerUpdate: v}

	case typeCodeSolarWinds:
		v := &sdk.SolarWindsNetworkPoolServerUpdate{
			Name: s.Name, Enabled: s.Enabled,
			ServiceUrl: s.ServiceUrl, ServiceUsername: s.ServiceUsername, ServicePassword: s.ServicePassword,
			ServiceThrottleRate: s.ServiceThrottleRate, IgnoreSsl: s.IgnoreSsl,
		}
		if inv != nil {
			v.Config = &sdk.SolarWindsNetworkPoolServerUpdateConfig{InventoryExisting: inv}
		}
		if s.Credential != nil {
			v.Credential = &sdk.SolarWindsNetworkPoolServerUpdateCredential{Type: s.Credential.Type, Id: s.Credential.Id}
		}

		return sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{SolarWindsNetworkPoolServerUpdate: v}

	default:
		return sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{InfobloxNetworkPoolServerUpdate: s}
	}
}
