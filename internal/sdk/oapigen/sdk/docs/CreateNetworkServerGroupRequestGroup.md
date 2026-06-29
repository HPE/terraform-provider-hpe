# CreateNetworkServerGroupRequestGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**CreateNetworkServerGroupRequestGroupAccount**](CreateNetworkServerGroupRequestGroupAccount.md) |  | [optional] 
**Owner** | Pointer to [**CreateNetworkServerGroupRequestGroupOwner**](CreateNetworkServerGroupRequestGroupOwner.md) |  | [optional] 
**NetworkServer** | Pointer to [**CreateNetworkServerGroupRequestGroupNetworkServer**](CreateNetworkServerGroupRequestGroupNetworkServer.md) |  | [optional] 
**Permissions** | Pointer to [**CreateNetworkServerGroupRequestGroupPermissions**](CreateNetworkServerGroupRequestGroupPermissions.md) |  | [optional] 
**Tags** | Pointer to [**[]CreateNetworkServerGroupRequestGroupTagsInner**](CreateNetworkServerGroupRequestGroupTagsInner.md) |  | [optional] 
**Members** | Pointer to [**[]CreateNetworkServerGroupRequestGroupMembersInner**](CreateNetworkServerGroupRequestGroupMembersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkServerGroupRequestGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


