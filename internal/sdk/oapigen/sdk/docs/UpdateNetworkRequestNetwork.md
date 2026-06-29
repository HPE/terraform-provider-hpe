# UpdateNetworkRequestNetwork

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DisplayName** | Pointer to **string** | Display Name | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**Cidr** | Pointer to **string** | CIDR Network | [optional] 
**Gateway** | Pointer to **string** | Network Gateway | [optional] 
**DnsPrimary** | Pointer to **string** | Primary DNS Server | [optional] 
**DnsSecondary** | Pointer to **string** | Secondary DNS Server | [optional] 
**VlanId** | Pointer to **int64** |  | [optional] 
**SwitchId** | Pointer to **string** | Network switch identifier | [optional] 
**Pool** | Pointer to **NullableInt64** | Network Pool ID | [optional] 
**ZonePool** | Pointer to [**UpdateNetworkRequestNetworkZonePool**](UpdateNetworkRequestNetworkZonePool.md) |  | [optional] 
**AllowStaticOverride** | Pointer to **bool** | Allow IP Override | [optional] 
**AssignPublicIp** | Pointer to **bool** | Assign Public IP | [optional] 
**Active** | Pointer to **bool** | Activate (true) or disable (false) the network | [optional] 
**DhcpServer** | Pointer to **bool** | DHCP Server enabled network | [optional] 
**NetworkDomain** | Pointer to [**UpdateNetworkRequestNetworkNetworkDomain**](UpdateNetworkRequestNetworkNetworkDomain.md) |  | [optional] 
**SearchDomains** | Pointer to **string** | Search Domains | [optional] 
**NetworkProxy** | Pointer to [**UpdateNetworkRequestNetworkNetworkProxy**](UpdateNetworkRequestNetworkNetworkProxy.md) |  | [optional] 
**ApplianceUrlProxyBypass** | Pointer to **bool** | Bypass Proxy for Appliance URL | [optional] 
**NoProxy** | Pointer to **NullableString** | Comma-separated list of ip addresses or name servers to exclude proxy traversal for. Typically locally routable servers are excluded. | [optional] 
**Visibility** | Pointer to **string** | Visibility, private or public. | [optional] [default to "private"]
**Config** | Pointer to **map[string]interface{}** | Configuration object. Settings vary by type. | [optional] 
**Tenants** | Pointer to [**[]UpdateNetworkRequestNetworkTenantsInner**](UpdateNetworkRequestNetworkTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**UpdateNetworkRequestNetworkResourcePermissions**](UpdateNetworkRequestNetworkResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateNetworkRequestNetwork{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value
### NoProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.NoProxy.IsSet()` — check if set
- `obj.NoProxy.Get()` — get the inner value (returns pointer)
- `obj.NoProxy.Set(&val)` — set the value
- `obj.NoProxy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


