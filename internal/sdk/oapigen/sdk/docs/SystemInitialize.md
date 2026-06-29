# SystemInitialize

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Override the system name. | [optional] 
**Description** | Pointer to **string** | Override the system description. | [optional] 
**Enabled** | Pointer to **bool** | Enable or disable the system. | [optional] 
**ExternalId** | Pointer to **string** | Override the external identifier. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Override the system configuration data. | [optional] 
**Layout** | Pointer to [**SystemInitializeLayout**](SystemInitializeLayout.md) |  | [optional] 
**Components** | Pointer to [**[]SystemInitializeComponentsInner**](SystemInitializeComponentsInner.md) | Optional component overrides. Components are matched by &#x60;id&#x60; first, then by &#x60;typeCode&#x60; when exactly one existing component matches. Updates name, externalId, and config on matched components, and can create a new component when a valid &#x60;typeCode&#x60; is supplied with no existing match.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SystemInitialize{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


