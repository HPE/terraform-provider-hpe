# ListSpecTemplates200ResponseAllOfSpecTemplatesInnerFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**ContentRef** | Pointer to **NullableString** |  | [optional] 
**ContentPath** | Pointer to **NullableString** |  | [optional] 
**Repository** | Pointer to **NullableString** |  | [optional] 
**Content** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListSpecTemplates200ResponseAllOfSpecTemplatesInnerFile{
    // Set fields directly
}
```

### ContentRef (Nullable)

Use the Nullable wrapper methods:
- `obj.ContentRef.IsSet()` — check if set
- `obj.ContentRef.Get()` — get the inner value (returns pointer)
- `obj.ContentRef.Set(&val)` — set the value
- `obj.ContentRef.Unset()` — clear the value
### ContentPath (Nullable)

Use the Nullable wrapper methods:
- `obj.ContentPath.IsSet()` — check if set
- `obj.ContentPath.Get()` — get the inner value (returns pointer)
- `obj.ContentPath.Set(&val)` — set the value
- `obj.ContentPath.Unset()` — clear the value
### Repository (Nullable)

Use the Nullable wrapper methods:
- `obj.Repository.IsSet()` — check if set
- `obj.Repository.Get()` — get the inner value (returns pointer)
- `obj.Repository.Set(&val)` — set the value
- `obj.Repository.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


