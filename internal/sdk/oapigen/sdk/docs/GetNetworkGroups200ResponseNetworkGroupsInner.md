# GetNetworkGroups200ResponseNetworkGroupsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Networks** | Pointer to **[]int64** |  | [optional] 
**Subnets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tenants** | Pointer to [**[]GetNetworkGroups200ResponseNetworkGroupsInnerTenantsInner**](GetNetworkGroups200ResponseNetworkGroupsInnerTenantsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkGroups200ResponseNetworkGroupsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


