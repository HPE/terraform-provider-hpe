# SearchHitsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | ID | [optional] 
**Uuid** | Pointer to **string** | UUID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **NullableTime** |  | [optional] 
**Score** | Pointer to **NullableFloat32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SearchHitsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### DateCreated (Nullable)

Use the Nullable wrapper methods:
- `obj.DateCreated.IsSet()` — check if set
- `obj.DateCreated.Get()` — get the inner value (returns pointer)
- `obj.DateCreated.Set(&val)` — set the value
- `obj.DateCreated.Unset()` — clear the value
### Score (Nullable)

Use the Nullable wrapper methods:
- `obj.Score.IsSet()` — check if set
- `obj.Score.Get()` — get the inner value (returns pointer)
- `obj.Score.Set(&val)` — set the value
- `obj.Score.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


