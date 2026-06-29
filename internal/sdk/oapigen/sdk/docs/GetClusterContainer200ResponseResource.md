# GetClusterContainer200ResponseResource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Title** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Metadata** | Pointer to **map[string]interface{}** |  | [optional] 
**Spec** | Pointer to **map[string]interface{}** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**RawSec** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterContainer200ResponseResource{
    // Set fields directly
}
```

### Spec (Nullable)

Use the Nullable wrapper methods:
- `obj.Spec.IsSet()` — check if set
- `obj.Spec.Get()` — get the inner value (returns pointer)
- `obj.Spec.Set(&val)` — set the value
- `obj.Spec.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### RawSec (Nullable)

Use the Nullable wrapper methods:
- `obj.RawSec.IsSet()` — check if set
- `obj.RawSec.Get()` — get the inner value (returns pointer)
- `obj.RawSec.Set(&val)` — set the value
- `obj.RawSec.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


