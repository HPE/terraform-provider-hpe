# AddSystemRequestSystem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the system. | 
**Description** | Pointer to **string** | Optional description for the system. | [optional] 
**Type** | [**AddSystemRequestSystemType**](AddSystemRequestSystemType.md) |  | 
**Layout** | [**AddSystemRequestSystemLayout**](AddSystemRequestSystemLayout.md) |  | 
**Enabled** | Pointer to **bool** | Whether the system is enabled. | [optional] 
**ExternalId** | Pointer to **string** | External identifier for the system. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Arbitrary configuration data for the system. | [optional] 
**Components** | Pointer to [**[]AddSystemRequestSystemComponentsInner**](AddSystemRequestSystemComponentsInner.md) | Optional component payloads for create. Existing components are matched by &#x60;id&#x60; first, then by &#x60;typeCode&#x60; when exactly one existing component matches. New components can be created by supplying a valid &#x60;typeCode&#x60;.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSystemRequestSystem{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


