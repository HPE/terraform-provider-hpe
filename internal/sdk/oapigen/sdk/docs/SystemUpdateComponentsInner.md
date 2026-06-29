# SystemUpdateComponentsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Existing system component id. When supplied, endpoints that support component updates match by id first. | [optional] 
**TypeCode** | Pointer to **string** | The code of the component type this entry applies to. | [optional] 
**Name** | Pointer to **string** | Optional override for the component name. Defaults to the component type name. | [optional] 
**ExternalId** | Pointer to **string** | External identifier for the component. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Arbitrary configuration data for the component. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SystemUpdateComponentsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


