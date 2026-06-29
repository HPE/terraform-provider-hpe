# TaskResponseFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **NullableInt64** |  | [optional] 
**SourceType** | Pointer to **NullableString** |  | [optional] 
**ContentRef** | Pointer to **NullableString** |  | [optional] 
**ContentPath** | Pointer to **NullableString** |  | [optional] 
**Repository** | Pointer to [**GetTasks200ResponseAllOfTaskFileRepository**](GetTasks200ResponseAllOfTaskFileRepository.md) |  | [optional] 
**Content** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &TaskResponseFile{
    // Set fields directly
}
```

### Id (Nullable)

Use the Nullable wrapper methods:
- `obj.Id.IsSet()` — check if set
- `obj.Id.Get()` — get the inner value (returns pointer)
- `obj.Id.Set(&val)` — set the value
- `obj.Id.Unset()` — clear the value
### SourceType (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceType.IsSet()` — check if set
- `obj.SourceType.Get()` — get the inner value (returns pointer)
- `obj.SourceType.Set(&val)` — set the value
- `obj.SourceType.Unset()` — clear the value
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
### Content (Nullable)

Use the Nullable wrapper methods:
- `obj.Content.IsSet()` — check if set
- `obj.Content.Get()` — get the inner value (returns pointer)
- `obj.Content.Set(&val)` — set the value
- `obj.Content.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


