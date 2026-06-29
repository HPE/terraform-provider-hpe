# AddStorageServersRequestStorageServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name | 
**Type** | **string** | The &#x60;Storage Type&#x60; Code or ID | 
**Description** | Pointer to **string** | description | [optional] 
**Enabled** | Pointer to **bool** | The enabled flag | [optional] [default to true]
**Config** | **map[string]interface{}** | Configuration object with parameters that vary by &#x60;type&#x60; | 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**ServiceHost** | Pointer to **string** | Storage server host | [optional] 
**ServiceUrl** | Pointer to **string** | Storage server URL | [optional] 
**ServiceUsername** | Pointer to **string** | Service username for authentication | [optional] 
**ServicePassword** | Pointer to **string** | Service password for authentication | [optional] 
**ServicePort** | Pointer to **int32** | Service port | [optional] [default to 22]
**Credential** | Pointer to [**AddStorageServersRequestStorageServerCredential**](AddStorageServersRequestStorageServerCredential.md) |  | [optional] 
**Tenants** | Pointer to [**[]AddStorageServersRequestStorageServerTenantsInner**](AddStorageServersRequestStorageServerTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddStorageServersRequestStorageServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


