# ProvisioningLicensesUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**LicenseVersion** | Pointer to **string** | License Version | [optional] 
**Copies** | Pointer to **int64** | Copies - The number of times the key can be used. | [optional] [default to 1]
**Description** | Pointer to **string** | Description | [optional] 
**VirtualImages** | Pointer to **[]int64** | Virtual Images - Array of Virtual Image IDs to associate with license. | [optional] 
**Tenants** | Pointer to **[]int64** | Tenants - Array of tenants that are allowed to use the key. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ProvisioningLicensesUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


