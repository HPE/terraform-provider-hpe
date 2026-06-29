# UpdateNetworkServerRequestNetworkServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**ServiceUrl** | Pointer to **string** | URL | [optional] 
**ServiceUsername** | Pointer to **string** | Username | [optional] 
**ServicePassword** | Pointer to **string** | Password | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**Tenants** | Pointer to [**[]NSXNetworkServerUpdateTenantsInner**](NSXNetworkServerUpdateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateNetworkServerRequestNetworkServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


