# GetNetworkServerGroup200ResponseGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**GetNetworkServerGroup200ResponseGroupAccount**](GetNetworkServerGroup200ResponseGroupAccount.md) |  | [optional] 
**Owner** | Pointer to [**GetNetworkServerGroup200ResponseGroupOwner**](GetNetworkServerGroup200ResponseGroupOwner.md) |  | [optional] 
**NetworkServer** | Pointer to [**GetNetworkServerGroup200ResponseGroupNetworkServer**](GetNetworkServerGroup200ResponseGroupNetworkServer.md) |  | [optional] 
**Permissions** | Pointer to [**GetNetworkServerGroup200ResponseGroupPermissions**](GetNetworkServerGroup200ResponseGroupPermissions.md) |  | [optional] 
**Tags** | Pointer to [**[]GetNetworkServerGroup200ResponseGroupTagsInner**](GetNetworkServerGroup200ResponseGroupTagsInner.md) |  | [optional] 
**Members** | Pointer to [**[]GetNetworkServerGroup200ResponseGroupMembersInner**](GetNetworkServerGroup200ResponseGroupMembersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkServerGroup200ResponseGroup{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


