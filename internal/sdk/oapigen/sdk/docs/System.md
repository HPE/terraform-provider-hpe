# System

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | The name of the system. | [optional] 
**Description** | Pointer to **string** | Optional description for the system. | [optional] 
**Type** | Pointer to [**SystemType**](SystemType.md) |  | [optional] 
**Layout** | Pointer to [**SystemLayout**](SystemLayout.md) |  | [optional] 
**Status** | Pointer to **string** | The current status of the system. | [optional] 
**StatusMessage** | Pointer to **NullableString** | A message describing the current status. | [optional] 
**Enabled** | Pointer to **bool** | Whether the system is enabled. | [optional] 
**ExternalId** | Pointer to **string** | External ID of the system. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Arbitrary configuration data for the system. | [optional] 
**Components** | Pointer to [**[]SystemComponentsInner**](SystemComponentsInner.md) | Component instances belonging to this system. | [optional] 
**DateCreated** | Pointer to **time.Time** | The date the system was created. | [optional] 
**LastUpdated** | Pointer to **time.Time** | The date the system was last updated. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &System{
    // Set fields directly
}
```

### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


