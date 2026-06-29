# GetNetworkRouters200ResponseNetworkRoutersInner

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
**ProviderId** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerType**](GetNetworkRouters200ResponseNetworkRoutersInnerType.md) |  | [optional] 
**NetworkServer** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerNetworkServer**](GetNetworkRouters200ResponseNetworkRoutersInnerNetworkServer.md) |  | [optional] 
**Zone** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerZone**](GetNetworkRouters200ResponseNetworkRoutersInnerZone.md) |  | [optional] 
**Instance** | Pointer to **NullableString** |  | [optional] 
**ExternalNetwork** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerExternalNetwork**](GetNetworkRouters200ResponseNetworkRoutersInnerExternalNetwork.md) |  | [optional] 
**Site** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerSite**](GetNetworkRouters200ResponseNetworkRoutersInnerSite.md) |  | [optional] 
**Interfaces** | Pointer to [**[]GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInner**](GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInner.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouters200ResponseNetworkRoutersInner{
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
### ProviderId (Nullable)

Use the Nullable wrapper methods:
- `obj.ProviderId.IsSet()` — check if set
- `obj.ProviderId.Get()` — get the inner value (returns pointer)
- `obj.ProviderId.Set(&val)` — set the value
- `obj.ProviderId.Unset()` — clear the value
### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


