# ServerBaremetalCreateServerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IloIpAddress** | Pointer to **string** |  | [optional] 
**IloUsername** | Pointer to **string** |  | [optional] 
**IloPassword** | Pointer to **string** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 
**ResourcePoolId** | Pointer to **int32** |  | [optional] 
**PreProvisioned** | Pointer to **bool** | Set to &#39;true&#39; if the server is pre-provisioned (brownfield) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ServerBaremetalCreateServerConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


