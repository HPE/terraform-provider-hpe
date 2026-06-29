# GetresourcePoolGroup200ResponseResourcePoolGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Mode** | Pointer to **string** | Pool selection mode. Valid values are &#x60;roundrobin&#x60; or &#x60;availablecapacity&#x60;. | [optional] 
**Pools** | Pointer to **[]int64** | Array of Resource Pool IDs | [optional] 
**Tenants** | Pointer to [**[]GetresourcePoolGroup200ResponseResourcePoolGroupTenantsInner**](GetresourcePoolGroup200ResponseResourcePoolGroupTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**GetresourcePoolGroup200ResponseResourcePoolGroupResourcePermission**](GetresourcePoolGroup200ResponseResourcePoolGroupResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetresourcePoolGroup200ResponseResourcePoolGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


