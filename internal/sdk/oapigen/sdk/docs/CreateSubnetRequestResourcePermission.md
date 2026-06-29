# CreateSubnetRequestResourcePermission

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**All** | Pointer to **bool** | Pass true to allow access all groups | [optional] 
**Sites** | Pointer to [**[]CreateSubnetRequestResourcePermissionSitesInner**](CreateSubnetRequestResourcePermissionSitesInner.md) | Array of groups ID objects that are allowed access | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**Plans** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateSubnetRequestResourcePermission{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


