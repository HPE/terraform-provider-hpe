# SystemUninitializedCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the system. | 
**Description** | Pointer to **string** | Optional description for the system. | [optional] 
**Type** | [**SystemUninitializedCreateType**](SystemUninitializedCreateType.md) |  | 
**Layout** | [**SystemUninitializedCreateLayout**](SystemUninitializedCreateLayout.md) |  | 
**Enabled** | Pointer to **bool** | Whether the system is enabled. | [optional] 
**ExternalId** | Pointer to **string** | External identifier for the system. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Arbitrary configuration data for the system. | [optional] 
**Components** | Pointer to [**[]SystemUninitializedCreateComponentsInner**](SystemUninitializedCreateComponentsInner.md) | Optional component payloads to enrich skeleton components at creation time. Each entry is matched to a layout component type by &#x60;typeCode&#x60;. Unmatched entries are ignored.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SystemUninitializedCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


