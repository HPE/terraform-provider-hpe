# GetNetworkRouter200ResponseNetworkRouter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**RouterType** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**EnableBgp** | Pointer to **bool** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**ProviderId** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterType**](GetNetworkRouter200ResponseNetworkRouterType.md) |  | [optional] 
**NetworkServer** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterNetworkServer**](GetNetworkRouter200ResponseNetworkRouterNetworkServer.md) |  | [optional] 
**Zone** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterZone**](GetNetworkRouter200ResponseNetworkRouterZone.md) |  | [optional] 
**Instance** | Pointer to **NullableString** |  | [optional] 
**ExternalNetwork** | Pointer to **NullableString** |  | [optional] 
**Site** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterSite**](GetNetworkRouter200ResponseNetworkRouterSite.md) |  | [optional] 
**Interfaces** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterInterfacesInner**](GetNetworkRouter200ResponseNetworkRouterInterfacesInner.md) |  | [optional] 
**Firewall** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterFirewall**](GetNetworkRouter200ResponseNetworkRouterFirewall.md) |  | [optional] 
**Routes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Nats** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Permissions** | Pointer to [**GetNetworkRouter200ResponseNetworkRouterPermissions**](GetNetworkRouter200ResponseNetworkRouterPermissions.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouter200ResponseNetworkRouter{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value
### ExternalNetwork (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalNetwork.IsSet()` — check if set
- `obj.ExternalNetwork.Get()` — get the inner value (returns pointer)
- `obj.ExternalNetwork.Set(&val)` — set the value
- `obj.ExternalNetwork.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


