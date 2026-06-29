# ListNetworkServerGroups200ResponseAllOfGroupsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**ListNetworkServerGroups200ResponseAllOfGroupsInnerAccount**](ListNetworkServerGroups200ResponseAllOfGroupsInnerAccount.md) |  | [optional] 
**Owner** | Pointer to [**ListNetworkServerGroups200ResponseAllOfGroupsInnerOwner**](ListNetworkServerGroups200ResponseAllOfGroupsInnerOwner.md) |  | [optional] 
**NetworkServer** | Pointer to [**ListNetworkServerGroups200ResponseAllOfGroupsInnerNetworkServer**](ListNetworkServerGroups200ResponseAllOfGroupsInnerNetworkServer.md) |  | [optional] 
**Permissions** | Pointer to [**ListNetworkServerGroups200ResponseAllOfGroupsInnerPermissions**](ListNetworkServerGroups200ResponseAllOfGroupsInnerPermissions.md) |  | [optional] 
**Tags** | Pointer to [**[]ListNetworkServerGroups200ResponseAllOfGroupsInnerTagsInner**](ListNetworkServerGroups200ResponseAllOfGroupsInnerTagsInner.md) |  | [optional] 
**Members** | Pointer to [**[]ListNetworkServerGroups200ResponseAllOfGroupsInnerMembersInner**](ListNetworkServerGroups200ResponseAllOfGroupsInnerMembersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListNetworkServerGroups200ResponseAllOfGroupsInner{
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


