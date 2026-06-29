# SystemUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | New name for the system. | [optional] 
**Description** | Pointer to **string** | New description for the system. | [optional] 
**Enabled** | Pointer to **bool** | Enable or disable the system. | [optional] 
**ExternalId** | Pointer to **string** | Override the external identifier. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Override the system configuration data. | [optional] 
**Status** | Pointer to **string** | Override the system status. | [optional] 
**StatusMessage** | Pointer to **string** | Override the status message. | [optional] 
**Layout** | Pointer to [**SystemUpdateLayout**](SystemUpdateLayout.md) |  | [optional] 
**Components** | Pointer to [**[]SystemUpdateComponentsInner**](SystemUpdateComponentsInner.md) | Optional authoritative component payloads. Components are matched by &#x60;id&#x60; first, then by &#x60;typeCode&#x60; when exactly one existing component matches. New components can be created by supplying a valid &#x60;typeCode&#x60;. When this field is present, omitted existing components are removed from the system.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SystemUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


