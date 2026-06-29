# GetResourcePoolGroups200ResponseResourcePoolGroupsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Mode** | Pointer to **string** | Pool selection mode. Valid values are &#x60;roundrobin&#x60; or &#x60;availablecapacity&#x60;. | [optional] 
**Pools** | Pointer to **[]int64** | Array of Resource Pool IDs | [optional] 
**Tenants** | Pointer to [**[]GetResourcePoolGroups200ResponseResourcePoolGroupsInnerTenantsInner**](GetResourcePoolGroups200ResponseResourcePoolGroupsInnerTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**GetResourcePoolGroups200ResponseResourcePoolGroupsInnerResourcePermission**](GetResourcePoolGroups200ResponseResourcePoolGroupsInnerResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetResourcePoolGroups200ResponseResourcePoolGroupsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


