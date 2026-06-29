# NetworkServerGroupCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**NetworkServerGroupCreateAccount**](NetworkServerGroupCreateAccount.md) |  | [optional] 
**Owner** | Pointer to [**NetworkServerGroupCreateOwner**](NetworkServerGroupCreateOwner.md) |  | [optional] 
**NetworkServer** | Pointer to [**NetworkServerGroupCreateNetworkServer**](NetworkServerGroupCreateNetworkServer.md) |  | [optional] 
**Permissions** | Pointer to [**NetworkServerGroupCreatePermissions**](NetworkServerGroupCreatePermissions.md) |  | [optional] 
**Tags** | Pointer to [**[]NetworkServerGroupCreateTagsInner**](NetworkServerGroupCreateTagsInner.md) |  | [optional] 
**Members** | Pointer to [**[]NetworkServerGroupCreateMembersInner**](NetworkServerGroupCreateMembersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkServerGroupCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


