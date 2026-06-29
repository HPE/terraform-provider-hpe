# AppStateWorkloadsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**SubRefName** | Pointer to **NullableString** |  | [optional] 
**StateDate** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**IacDrift** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AppStateWorkloadsInner{
    // Set fields directly
}
```

### SubRefName (Nullable)

Use the Nullable wrapper methods:
- `obj.SubRefName.IsSet()` — check if set
- `obj.SubRefName.Get()` — get the inner value (returns pointer)
- `obj.SubRefName.Set(&val)` — set the value
- `obj.SubRefName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


